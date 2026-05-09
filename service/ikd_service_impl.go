package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/ikd"
	"ekak_kabupaten_madiun/repository"
)

type IkdServiceImpl struct {
	IkdRepository repository.IkdRepository
	DB            *sql.DB
}

func NewIkdServiceImpl(ikdRepository repository.IkdRepository, db *sql.DB) *IkdServiceImpl {
	return &IkdServiceImpl{
		IkdRepository: ikdRepository,
		DB:            db,
	}
}

func (service *IkdServiceImpl) FindAll(
	ctx context.Context,
	kodeOpd string,
	tahun string,
	jenisPeriode string,
) ([]ikd.IkdResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return []ikd.IkdResponse{}, err
	}

	defer helper.CommitOrRollback(tx)

	ikdDetails, err := service.IkdRepository.FindAll(
		ctx,
		tx,
		kodeOpd,
		tahun,
		jenisPeriode,
	)

	if err != nil {
		return []ikd.IkdResponse{}, err
	}

	var responses []ikd.IkdResponse

	for _, data := range ikdDetails {

		// =========================
		// PELAKSANA
		// =========================
		pelaksanas := make([]ikd.PelaksanaResponse, 0)

		for _, pelaksana := range data.Pelaksana {

			pelaksanas = append(
				pelaksanas,
				ikd.PelaksanaResponse{
					Id:          pelaksana.Id,
					PegawaiId:   pelaksana.PegawaiId,
					Nip:         pelaksana.Nip,
					NamaPegawai: pelaksana.NamaPegawai,
				},
			)
		}

		// =========================
		// SASARAN OPD
		// =========================
		sasaranResponses := make([]ikd.SasaranOpdResponse, 0)

		for _, sasaran := range data.SasaranOpd {

			// =====================
			// INDIKATOR
			// =====================
			indikatorResponses := make([]ikd.IndikatorResponse, 0)

			for _, indikator := range sasaran.Indikator {

				// =================
				// TARGET
				// =================
				targetResponses := make([]ikd.TargetResponse, 0)

				for _, target := range indikator.Target {

					targetResponses = append(
						targetResponses,
						ikd.TargetResponse{
							Id:          target.Id,
							IndikatorId: target.IndikatorId,
							Tahun:       target.Tahun,
							Target:      target.Target,
							Satuan:      target.Satuan,
						},
					)
				}

				indikatorResponses = append(
					indikatorResponses,
					ikd.IndikatorResponse{
						Id:                indikator.Id,
						Indikator:         indikator.Indikator,
						RumusPerhitungan:  indikator.RumusPerhitungan.String,
						SumberData:        indikator.SumberData.String,
						Target:            targetResponses,
					},
				)
			}

			sasaranResponses = append(
				sasaranResponses,
				ikd.SasaranOpdResponse{
					Id:              sasaran.Id,
					IdPohon:         sasaran.IdPohon,
					NamaSasaranOpd:  sasaran.NamaSasaranOpd,
					IdTujuanOpd:     sasaran.IdTujuanOpd,
					NamaTujuanOpd:   sasaran.NamaTujuanOpd,
					TahunAwal:       sasaran.TahunAwal,
					TahunAkhir:      sasaran.TahunAkhir,
					JenisPeriode:    sasaran.JenisPeriode,
					Indikator:       indikatorResponses,
				},
			)
		}

		// =========================
		// PROGRAM OPD
		// =========================
		programResponses := make([]ikd.ProgramOpdResponse, 0)

		for _, program := range data.ProgramOpd {

			programResponses = append(
				programResponses,
				ikd.ProgramOpdResponse{
					Id:          program.Id,
					Parent:      program.Parent,
					NamaProgram: program.NamaProgram,
				},
			)
		}

		// =========================
		// RESPONSE
		// =========================
		responses = append(
			responses,
			ikd.IkdResponse{
				Id:                     data.Id,
				NamaPohon:              data.NamaPohon,
				Parent:                 data.Parent,
				JenisPohon:             data.JenisPohon,
				LevelPohon:             data.LevelPohon,
				KodeOpd:                data.KodeOpd,
				Keterangan:             data.Keterangan,
				KeteranganCrosscutting: data.KeteranganCrosscutting,
				Tahun:                  data.Tahun,
				Status:                 data.Status,
				IsActive:               data.IsActive,

				Pelaksana:  pelaksanas,
				SasaranOpd: sasaranResponses,
				ProgramOpd: programResponses,
			},
		)
	}

	if responses == nil {
		responses = make([]ikd.IkdResponse, 0)
	}

	return responses, nil
}