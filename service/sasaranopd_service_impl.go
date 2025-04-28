package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strconv"
)

type SasaranOpdServiceImpl struct {
	sasaranOpdRepository      repository.SasaranOpdRepository
	opdRepository             repository.OpdRepository
	rencanaKinerjaRepository  repository.RencanaKinerjaRepository
	manualIndikatorRepository repository.ManualIKRepository
	pegawaiRepository         repository.PegawaiRepository
	pohonkinerjaRepository    repository.PohonKinerjaRepository
	DB                        *sql.DB
}

func NewSasaranOpdServiceImpl(
	sasaranOpdRepository repository.SasaranOpdRepository,
	opdRepository repository.OpdRepository,
	rencanaKinerjaRepository repository.RencanaKinerjaRepository,
	manualIndikatorRepository repository.ManualIKRepository,
	pegawaiRepository repository.PegawaiRepository,
	pohonkinerjaRepository repository.PohonKinerjaRepository,
	db *sql.DB) *SasaranOpdServiceImpl {
	return &SasaranOpdServiceImpl{
		sasaranOpdRepository:      sasaranOpdRepository,
		opdRepository:             opdRepository,
		rencanaKinerjaRepository:  rencanaKinerjaRepository,
		manualIndikatorRepository: manualIndikatorRepository,
		pegawaiRepository:         pegawaiRepository,
		pohonkinerjaRepository:    pohonkinerjaRepository,
		DB:                        db,
	}
}

func (service *SasaranOpdServiceImpl) FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranOpds, err := service.sasaranOpdRepository.FindAll(ctx, tx, KodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	var responses []sasaranopd.SasaranOpdResponse
	for _, sasaranOpd := range sasaranOpds {
		response := sasaranopd.SasaranOpdResponse{
			Id:             sasaranOpd.Id,
			IdPohon:        sasaranOpd.IdPohon,
			NamaPohon:      sasaranOpd.NamaPohon,
			JenisPohon:     sasaranOpd.JenisPohon,
			LevelPohon:     sasaranOpd.LevelPohon,
			TahunPohon:     sasaranOpd.TahunPohon,
			RencanaKinerja: make([]sasaranopd.RencanaKinerjaOpd, 0),
			Pelaksana:      make([]sasaranopd.PelaksanaOpdResponse, 0),
		}

		// Convert Pelaksana
		for _, pelaksana := range sasaranOpd.Pelaksana {
			response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
				Id:          pelaksana.Id,
				PegawaiId:   pelaksana.PegawaiId,
				Nip:         pelaksana.Nip,
				NamaPegawai: pelaksana.NamaPegawai,
			})
		}

		// Convert RencanaKinerja
		for _, rekin := range sasaranOpd.RencanaKinerja {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rekin.PegawaiId)
			if err != nil {
				return nil, err
			}
			rencanaKinerjaResponse := sasaranopd.RencanaKinerjaOpd{
				Id:                 rekin.Id,
				NamaRencanaKinerja: rekin.NamaRencanaKinerja,
				Nip:                rekin.PegawaiId,
				NamaPegawai:        pegawai.NamaPegawai,
				TahunAwal:          rekin.TahunAwal,
				TahunAkhir:         rekin.TahunAkhir,
				JenisPeriode:       rekin.JenisPeriode,
				Indikator:          make([]sasaranopd.IndikatorResponse, 0),
			}

			// Convert Indikator
			for _, indikator := range rekin.Indikator {
				indResponse := sasaranopd.IndikatorResponse{
					Id:        indikator.Id,
					Indikator: indikator.Indikator,
					Target:    make([]sasaranopd.TargetResponse, 0),
				}

				if indikator.ManualIK != nil {
					indResponse.ManualIK = &sasaranopd.ManualIKResponse{
						Formula:    indikator.ManualIK.Formula,
						SumberData: indikator.ManualIK.SumberData,
					}
				}

				for _, target := range indikator.Target {
					indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
						Tahun:  target.Tahun,
						Target: target.Target,
						Satuan: target.Satuan,
					})
				}

				rencanaKinerjaResponse.Indikator = append(rencanaKinerjaResponse.Indikator, indResponse)
			}

			response.RencanaKinerja = append(response.RencanaKinerja, rencanaKinerjaResponse)
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *SasaranOpdServiceImpl) FindByIdRencanaKinerja(ctx context.Context, idRencanaKinerja string) (*sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranOpd, err := service.sasaranOpdRepository.FindByIdRencanaKinerja(ctx, tx, idRencanaKinerja)
	if err != nil {
		return nil, err
	}

	response := &sasaranopd.SasaranOpdResponse{
		Id:             sasaranOpd.Id,
		IdPohon:        sasaranOpd.IdPohon,
		NamaPohon:      sasaranOpd.NamaPohon,
		JenisPohon:     sasaranOpd.JenisPohon,
		LevelPohon:     sasaranOpd.LevelPohon,
		TahunPohon:     sasaranOpd.TahunPohon,
		RencanaKinerja: make([]sasaranopd.RencanaKinerjaOpd, 0),
		Pelaksana:      make([]sasaranopd.PelaksanaOpdResponse, 0),
	}

	// Convert Pelaksana
	for _, pelaksana := range sasaranOpd.Pelaksana {
		response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
			Id:          pelaksana.Id,
			PegawaiId:   pelaksana.PegawaiId,
			Nip:         pelaksana.Nip,
			NamaPegawai: pelaksana.NamaPegawai,
		})
	}

	// Convert RencanaKinerja
	for _, rekin := range sasaranOpd.RencanaKinerja {
		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rekin.PegawaiId)
		if err != nil {
			return nil, err
		}
		rencanaKinerjaResponse := sasaranopd.RencanaKinerjaOpd{
			Id:                 rekin.Id,
			NamaRencanaKinerja: rekin.NamaRencanaKinerja,
			Nip:                rekin.PegawaiId,
			NamaPegawai:        pegawai.NamaPegawai,
			TahunAwal:          rekin.TahunAwal,
			TahunAkhir:         rekin.TahunAkhir,
			JenisPeriode:       rekin.JenisPeriode,
			Indikator:          make([]sasaranopd.IndikatorResponse, 0),
		}

		// Convert Indikator
		for _, indikator := range rekin.Indikator {
			indResponse := sasaranopd.IndikatorResponse{
				Id:         indikator.Id,
				Indikator:  indikator.Indikator,
				Formula:    indikator.ManualIK.Formula,
				SumberData: indikator.ManualIK.SumberData,
				Target:     make([]sasaranopd.TargetResponse, 0),
			}

			// if indikator.ManualIK != nil {
			// 	indResponse.ManualIK = &sasaranopd.ManualIKResponse{
			// 		Formula:    indikator.ManualIK.Formula,
			// 		SumberData: indikator.ManualIK.SumberData,
			// 	}
			// }

			for _, target := range indikator.Target {
				indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
					Tahun:  target.Tahun,
					Target: target.Target,
					Satuan: target.Satuan,
				})
			}

			rencanaKinerjaResponse.Indikator = append(rencanaKinerjaResponse.Indikator, indResponse)
		}

		response.RencanaKinerja = append(response.RencanaKinerja, rencanaKinerjaResponse)
	}

	return response, nil
}

func (service *SasaranOpdServiceImpl) FindIdPokinSasaran(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data pohon kinerja dengan sasaran
	pokin, err := service.sasaranOpdRepository.FindIdPokinSasaran(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Ambil data OPD jika ada
	var namaOpd string
	if pokin.KodeOpd != "" {
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err == nil {
			namaOpd = opd.NamaOpd
		}
	}

	// Konversi indikator dan target ke response
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range pokin.Indikator {
		var targetResponses []pohonkinerja.TargetResponse

		// Cari target yang ada di database
		var existingTarget *domain.Target
		for _, t := range ind.Target {
			if t.Id != fmt.Sprintf("TRG-%s-%s", ind.Id, pokin.Tahun) {
				existingTarget = &t
				break
			}
		}

		// Jika ada target di database, gunakan itu
		if existingTarget != nil {
			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              existingTarget.Id,
				IndikatorId:     existingTarget.IndikatorId,
				TargetIndikator: existingTarget.Target,
				SatuanIndikator: existingTarget.Satuan,
				TahunSasaran:    pokin.Tahun,
			})
		} else {
			// Jika tidak ada target di database, gunakan target default
			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              fmt.Sprintf("TRG-%s-%s", ind.Id, pokin.Tahun),
				IndikatorId:     ind.Id,
				TargetIndikator: "-",
				SatuanIndikator: "-",
				TahunSasaran:    pokin.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			IdPokin:       fmt.Sprint(pokin.Id),
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		})
	}

	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		NamaOpd:    namaOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Status:     pokin.Status,

		Indikator: indikatorResponses,
	}

	return response, nil
}
