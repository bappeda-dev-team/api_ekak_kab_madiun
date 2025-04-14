package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

type CascadingOpdServiceImpl struct {
	pohonKinerjaRepository   repository.PohonKinerjaRepository
	opdRepository            repository.OpdRepository
	pegawaiRepository        repository.PegawaiRepository
	tujuanOpdRepository      repository.TujuanOpdRepository
	rencanaKinerjaRepository repository.RencanaKinerjaRepository
	DB                       *sql.DB
	programRepository        repository.ProgramRepository
	cascadingOpdRepository   repository.CascadingOpdRepository
	bidangUrusanRepository   repository.BidangUrusanRepository
	rincianBelanjaRepository repository.RincianBelanjaRepository
	rencanaAksiRepository    repository.RencanaAksiRepository
}

func NewCascadingOpdServiceImpl(
	pohonKinerjaRepository repository.PohonKinerjaRepository,
	opdRepository repository.OpdRepository,
	pegawaiRepository repository.PegawaiRepository,
	tujuanOpdRepository repository.TujuanOpdRepository,
	rencanaKinerjaRepository repository.RencanaKinerjaRepository,
	DB *sql.DB, programRepository repository.ProgramRepository,
	cascadingOpdRepository repository.CascadingOpdRepository,
	bidangUrusanRepository repository.BidangUrusanRepository,
	rincianBelanjaRepository repository.RincianBelanjaRepository,
	rencanaAksiRepository repository.RencanaAksiRepository) *CascadingOpdServiceImpl {
	return &CascadingOpdServiceImpl{
		pohonKinerjaRepository:   pohonKinerjaRepository,
		opdRepository:            opdRepository,
		pegawaiRepository:        pegawaiRepository,
		tujuanOpdRepository:      tujuanOpdRepository,
		rencanaKinerjaRepository: rencanaKinerjaRepository,
		DB:                       DB,
		programRepository:        programRepository,
		cascadingOpdRepository:   cascadingOpdRepository,
		bidangUrusanRepository:   bidangUrusanRepository,
		rincianBelanjaRepository: rincianBelanjaRepository,
		rencanaAksiRepository:    rencanaAksiRepository,
	}
}

func (service *CascadingOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.CascadingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return pohonkinerja.CascadingOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		log.Printf("Error: OPD not found for kode_opd=%s: %v", kodeOpd, err)
		return pohonkinerja.CascadingOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// Inisialisasi response dasar
	response := pohonkinerja.CascadingOpdResponse{
		KodeOpd:    kodeOpd,
		NamaOpd:    opd.NamaOpd,
		Tahun:      tahun,
		TujuanOpd:  make([]pohonkinerja.TujuanOpdCascadingResponse, 0),
		Strategics: make([]pohonkinerja.StrategicCascadingOpdResponse, 0),
	}

	// Ambil data tujuan OPD
	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdForCascadingOpd(ctx, tx, kodeOpd, tahun, "RPJMD")
	if err != nil {
		log.Printf("Warning: Failed to get tujuan OPD data: %v", err)
		return response, nil
	}

	log.Printf("Processing %d tujuan OPD entries", len(tujuanOpds))

	// Konversi tujuan OPD ke format response
	for _, tujuan := range tujuanOpds {
		indikators, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdId(ctx, tx, tujuan.Id)
		if err != nil {
			log.Printf("Warning: Failed to get indicators for tujuan ID %d: %v", tujuan.Id, err)
			continue
		}

		var indikatorResponses []pohonkinerja.IndikatorTujuanResponse
		for _, indikator := range indikators {
			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
				Indikator: indikator.Indikator,
			})
		}

		bidangUrusan, err := service.bidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, tujuan.KodeBidangUrusan)
		if err != nil {
			log.Printf("Warning: Failed to get bidang urusan data: %v", err)
			continue
		}

		response.TujuanOpd = append(response.TujuanOpd, pohonkinerja.TujuanOpdCascadingResponse{
			Id:         tujuan.Id,
			KodeOpd:    tujuan.KodeOpd,
			Tujuan:     tujuan.Tujuan,
			KodeBidang: tujuan.KodeBidangUrusan,
			NamaBidang: bidangUrusan.NamaBidangUrusan,
			Indikator:  indikatorResponses,
		})
	}

	// Ambil data pohon kinerja
	pokins, err := service.cascadingOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		log.Printf("Error getting pohon kinerja data: %v", err)
		return response, nil
	}

	if len(pokins) == 0 {
		log.Printf("No pohon kinerja found for kodeOpd=%s, tahun=%s", kodeOpd, tahun)
		return response, nil
	}

	log.Printf("Processing %d pohon kinerja entries", len(pokins))

	// Proses data pohon kinerja
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)
	rencanaKinerjaMap := make(map[int][]domain.RencanaKinerja)

	// Kelompokkan data dan ambil data indikator & rencana kinerja
	maxLevel := 0
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			if p.LevelPohon > maxLevel {
				maxLevel = p.LevelPohon
			}

			if pohonMap[p.LevelPohon] == nil {
				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
			}

			p.NamaOpd = opd.NamaOpd
			pohonMap[p.LevelPohon][p.Parent] = append(
				pohonMap[p.LevelPohon][p.Parent],
				p,
			)

			// Ambil data indikator dan target
			indikators, err := service.cascadingOpdRepository.FindIndikatorByPokinId(ctx, tx, fmt.Sprint(p.Id))
			if err == nil {
				var indikatorResponses []pohonkinerja.IndikatorResponse
				for _, indikator := range indikators {
					var targetResponses []pohonkinerja.TargetResponse
					for _, target := range indikator.Target {
						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
							Id:              target.Id,
							IndikatorId:     target.IndikatorId,
							TargetIndikator: target.Target,
							SatuanIndikator: target.Satuan,
						})
					}

					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
						Id:            indikator.Id,
						IdPokin:       fmt.Sprint(p.Id),
						NamaIndikator: indikator.Indikator,
						Target:        targetResponses,
					})
				}
				indikatorMap[p.Id] = indikatorResponses
				log.Printf("Added %d indicators for pohon ID %d", len(indikatorResponses), p.Id)
			} else {
				log.Printf("Warning: Failed to get indicators for pohon ID %d: %v", p.Id, err)
			}

			// Ambil data rencana kinerja
			rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, p.Id)
			if err == nil {
				var validRencanaKinerja []domain.RencanaKinerja

				// Ambil daftar pelaksana untuk pohon kinerja ini
				pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(p.Id))
				if err == nil {
					// Buat map untuk mempermudah pengecekan
					pelaksanaMap := make(map[string]*domainmaster.Pegawai)
					rekinMap := make(map[string]bool) // untuk tracking pegawai yang sudah punya rekin

					// Simpan data pegawai ke map
					for _, pelaksana := range pelaksanaList {
						pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
						if err == nil {
							pelaksanaMap[pegawai.Nip] = &pegawai
						}
					}

					// Proses rencana kinerja yang ada
					for _, rk := range rencanaKinerjaList {
						if pegawai, exists := pelaksanaMap[rk.PegawaiId]; exists {
							rk.NamaPegawai = pegawai.NamaPegawai
							validRencanaKinerja = append(validRencanaKinerja, rk)
							rekinMap[rk.PegawaiId] = true
						}
					}

					// Tambahkan pelaksana yang tidak memiliki rencana kinerja
					for _, pegawai := range pelaksanaMap {
						if !rekinMap[pegawai.Nip] {
							// Tambahkan rencana kinerja kosong, hanya dengan data pegawai
							validRencanaKinerja = append(validRencanaKinerja, domain.RencanaKinerja{
								IdPohon:     p.Id,
								NamaPohon:   p.NamaPohon,
								Tahun:       p.Tahun,
								PegawaiId:   pegawai.Nip,
								NamaPegawai: pegawai.NamaPegawai,
								// Field lain dibiarkan kosong
								Indikator: nil, // Pastikan indikator kosong
							})
						}
					}
				}

				rencanaKinerjaMap[p.Id] = validRencanaKinerja
			} else {
				log.Printf("Warning: Failed to get rencana kinerja for pohon ID %d: %v", p.Id, err)
			}
		}
	}

	// Build response untuk strategic (level 4)
	if strategicList := pohonMap[4]; len(strategicList) > 0 {
		log.Printf("Processing %d strategic entries", len(strategicList))
		for _, strategicsByParent := range strategicList {
			sort.Slice(strategicsByParent, func(i, j int) bool {
				return strategicsByParent[i].Id < strategicsByParent[j].Id
			})

			for _, strategic := range strategicsByParent {
				strategicResp := service.buildStrategicCascadingResponse(ctx, tx, pohonMap, strategic, indikatorMap, rencanaKinerjaMap)
				response.Strategics = append(response.Strategics, strategicResp)
			}
		}

		sort.Slice(response.Strategics, func(i, j int) bool {
			return response.Strategics[i].Id < response.Strategics[j].Id
		})
	}

	log.Printf("FindAll completed: Processed %d strategic entries", len(response.Strategics))
	return response, nil
}

func (service *CascadingOpdServiceImpl) buildStrategicCascadingResponse(
	ctx context.Context,
	tx *sql.Tx,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	strategic domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.StrategicCascadingOpdResponse {

	// Proses rencana kinerja untuk strategic
	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
	// Proses pelaksana dari pohon kinerja
	if rencanaKinerjaList, ok := rencanaKinerjaMap[strategic.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			var indikatorResponses []pohonkinerja.IndikatorResponse

			// Hanya ambil indikator jika rencana kinerja memiliki ID
			if rk.Id != "" {
				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
				if err == nil {
					for _, ind := range indikatorRekin {
						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
						var targetResponses []pohonkinerja.TargetResponse
						if err == nil {
							for _, target := range targets {
								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
									Id:              target.Id,
									IndikatorId:     target.IndikatorId,
									TargetIndikator: target.Target,
									SatuanIndikator: target.Satuan,
								})
							}
						}
						indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							IdRekin:       rk.Id,
							NamaIndikator: ind.Indikator,
							Target:        targetResponses,
						})
					}
				}
			}

			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
				Id:                 rk.Id,
				IdPohon:            strategic.Id,
				NamaPohon:          strategic.NamaPohon,
				NamaRencanaKinerja: rk.NamaRencanaKinerja,
				Tahun:              strategic.Tahun,
				PegawaiId:          rk.PegawaiId,
				NamaPegawai:        rk.NamaPegawai,
				Indikator:          indikatorResponses,
			})
		}
	}

	// Proses program dari level 6 melalui level 5
	programMap := make(map[string]string)
	if tacticalList, exists := pohonMap[5][strategic.Id]; exists {
		for _, tactical := range tacticalList {
			if operationalList, exists := pohonMap[6][tactical.Id]; exists {
				for _, operational := range operationalList {
					if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
						for _, rk := range rencanaKinerjaList {
							if rk.KodeSubKegiatan != "" {
								segments := strings.Split(rk.KodeSubKegiatan, ".")
								if len(segments) >= 3 {
									kodeProgram := strings.Join(segments[:3], ".")
									if _, exists := programMap[kodeProgram]; !exists {
										program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
										if err == nil {
											programMap[kodeProgram] = program.NamaProgram
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Convert program map ke slice dan sort
	var programList []pohonkinerja.ProgramResponse
	for kodeProgram, namaProgram := range programMap {
		// Ambil indikator program
		var indikatorProgramResponses []pohonkinerja.IndikatorResponse
		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx,
			tx,
			kodeProgram,
			strategic.KodeOpd,
			strategic.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorProgram {
				// Ambil target untuk indikator program
				targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				if err == nil {
					for _, target := range targets {
						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
							Id:              target.Id,
							IndikatorId:     target.IndikatorId,
							TargetIndikator: target.Target,
							SatuanIndikator: target.Satuan,
						})
					}
				}

				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses, // Menambahkan target ke indikator program
				})
			}
		}

		programList = append(programList, pohonkinerja.ProgramResponse{
			KodeProgram: kodeProgram,
			NamaProgram: namaProgram,
			Indikator:   indikatorProgramResponses,
		})
	}

	strategicResp := pohonkinerja.StrategicCascadingOpdResponse{
		Id:         strategic.Id,
		Parent:     nil,
		Strategi:   strategic.NamaPohon,
		JenisPohon: strategic.JenisPohon,
		LevelPohon: strategic.LevelPohon,
		Keterangan: strategic.Keterangan,
		Status:     strategic.Status,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		IsActive:       strategic.IsActive,
		Program:        programList,
		RencanaKinerja: rencanaKinerjaResponses,
		Indikator:      indikatorMap[strategic.Id],
	}

	// Build tactical responses dan hitung total pagu anggaran
	var totalPaguAnggaran int64 = 0
	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
		var tacticals []pohonkinerja.TacticalCascadingOpdResponse
		sort.Slice(tacticalList, func(i, j int) bool {
			return tacticalList[i].Id < tacticalList[j].Id
		})

		for _, tactical := range tacticalList {
			tacticalResp := service.buildTacticalCascadingResponse(ctx, tx, pohonMap, tactical, indikatorMap, rencanaKinerjaMap)
			tacticals = append(tacticals, tacticalResp)
			// Tambahkan pagu anggaran dari setiap tactical
			totalPaguAnggaran += tacticalResp.PaguAnggaran
		}
		strategicResp.Tacticals = tacticals
	}

	// Set pagu anggaran dari total pagu anggaran tactical
	strategicResp.PaguAnggaran = totalPaguAnggaran

	return strategicResp
}

func (service *CascadingOpdServiceImpl) buildTacticalCascadingResponse(
	ctx context.Context,
	tx *sql.Tx,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	tactical domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.TacticalCascadingOpdResponse {

	log.Printf("Building tactical response for ID: %d, Level: %d", tactical.Id, tactical.LevelPohon)

	// Map untuk menyimpan program unik
	programMap := make(map[string]string)

	// Proses rencana kinerja untuk tactical
	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
	if rencanaKinerjaList, ok := rencanaKinerjaMap[tactical.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			var indikatorResponses []pohonkinerja.IndikatorResponse

			// Hanya ambil indikator jika rencana kinerja memiliki ID
			if rk.Id != "" {
				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
				if err == nil {
					for _, ind := range indikatorRekin {
						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
						var targetResponses []pohonkinerja.TargetResponse
						if err == nil {
							for _, target := range targets {
								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
									Id:              target.Id,
									IndikatorId:     target.IndikatorId,
									TargetIndikator: target.Target,
									SatuanIndikator: target.Satuan,
								})
							}
						}
						indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							IdRekin:       rk.Id,
							NamaIndikator: ind.Indikator,
							Target:        targetResponses,
						})
					}
				}
			}

			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
				Id:                 rk.Id,
				IdPohon:            tactical.Id,
				NamaPohon:          tactical.NamaPohon,
				NamaRencanaKinerja: rk.NamaRencanaKinerja,
				Tahun:              tactical.Tahun,
				PegawaiId:          rk.PegawaiId,
				NamaPegawai:        rk.NamaPegawai,
				Indikator:          indikatorResponses,
			})
		}
	}

	// Proses program dari level operational
	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
		for _, operational := range operationalList {
			if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
				for _, rk := range rencanaKinerjaList {
					if rk.KodeSubKegiatan != "" {
						segments := strings.Split(rk.KodeSubKegiatan, ".")
						if len(segments) >= 3 {
							kodeProgram := strings.Join(segments[:3], ".")
							if _, exists := programMap[kodeProgram]; !exists {
								program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
								if err == nil {
									programMap[kodeProgram] = program.NamaProgram
								}
							}
						}
					}
				}
			}
		}
	}

	// Convert program map ke slice dan sort
	var programList []pohonkinerja.ProgramResponse
	for kodeProgram, namaProgram := range programMap {
		// Ambil indikator program
		var indikatorProgramResponses []pohonkinerja.IndikatorResponse
		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx,
			tx,
			kodeProgram,
			tactical.KodeOpd,
			tactical.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorProgram {
				// Ambil target untuk indikator program
				targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				if err == nil {
					for _, target := range targets {
						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
							Id:              target.Id,
							IndikatorId:     target.IndikatorId,
							TargetIndikator: target.Target,
							SatuanIndikator: target.Satuan,
						})
					}
				}

				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		}

		programList = append(programList, pohonkinerja.ProgramResponse{
			KodeProgram: kodeProgram,
			NamaProgram: namaProgram,
			Indikator:   indikatorProgramResponses,
		})
	}

	tacticalResp := pohonkinerja.TacticalCascadingOpdResponse{
		Id:         tactical.Id,
		Parent:     tactical.Parent,
		Strategi:   tactical.NamaPohon,
		JenisPohon: tactical.JenisPohon,
		LevelPohon: tactical.LevelPohon,
		Keterangan: tactical.Keterangan,
		Status:     tactical.Status,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		},
		IsActive:       tactical.IsActive,
		Program:        programList,
		RencanaKinerja: rencanaKinerjaResponses,
		Indikator:      indikatorMap[tactical.Id],
	}

	// Build operational responses dan hitung total pagu anggaran
	var totalPaguAnggaran int64 = 0
	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
		var operationals []pohonkinerja.OperationalCascadingOpdResponse
		sort.Slice(operationalList, func(i, j int) bool {
			return operationalList[i].Id < operationalList[j].Id
		})

		for _, operational := range operationalList {
			operationalResp := service.buildOperationalCascadingResponse(ctx, tx, pohonMap, operational, indikatorMap, rencanaKinerjaMap)
			operationals = append(operationals, operationalResp)
			// Tambahkan total anggaran dari setiap operational
			totalPaguAnggaran += operationalResp.TotalAnggaran
		}
		tacticalResp.Operationals = operationals
	}

	// Set pagu anggaran dari total anggaran operational
	tacticalResp.PaguAnggaran = totalPaguAnggaran

	return tacticalResp
}

func (service *CascadingOpdServiceImpl) buildOperationalCascadingResponse(
	ctx context.Context,
	tx *sql.Tx,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	operational domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.OperationalCascadingOpdResponse {

	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalResponse
	var totalAnggaranOperational int64 = 0
	if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			var totalAnggaranRenkin int64 = 0
			if rk.Id != "" {
				// Ambil semua rencana aksi untuk rencana kinerja ini
				rencanaAksiList, err := service.rencanaAksiRepository.FindAll(ctx, tx, rk.Id)
				if err == nil {
					for _, ra := range rencanaAksiList {
						// Ambil anggaran untuk setiap rencana aksi
						rincianBelanja, err := service.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, ra.Id)
						if err == nil {
							totalAnggaranRenkin += rincianBelanja.Anggaran
						}
					}
				}
			}

			totalAnggaranOperational += totalAnggaranRenkin

			// Indikator rencana kinerja
			var indikatorRekinResponses []pohonkinerja.IndikatorResponse
			if rk.Id != "" {
				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
				if err == nil {
					for _, ind := range indikatorRekin {
						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
						var targetResponses []pohonkinerja.TargetResponse
						if err == nil {
							for _, target := range targets {
								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
									Id:              target.Id,
									IndikatorId:     target.IndikatorId,
									TargetIndikator: target.Target,
									SatuanIndikator: target.Satuan,
								})
							}
						}
						indikatorRekinResponses = append(indikatorRekinResponses, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							IdRekin:       rk.Id,
							NamaIndikator: ind.Indikator,
							Target:        targetResponses,
						})
					}
				}
			}

			// Indikator subkegiatan
			var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse
			if rk.KodeSubKegiatan != "" {
				indikatorSubkegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
					ctx,
					tx,
					rk.KodeSubKegiatan,
					operational.KodeOpd,
					operational.Tahun,
				)
				if err == nil {
					for _, ind := range indikatorSubkegiatan {
						// Ambil target untuk indikator
						targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
						var targetResponses []pohonkinerja.TargetResponse
						if err == nil {
							for _, target := range targets {
								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
									Id:              target.Id,
									IndikatorId:     target.IndikatorId,
									TargetIndikator: target.Target,
									SatuanIndikator: target.Satuan,
								})
							}
						}

						indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							Kode:          ind.Kode,
							NamaIndikator: ind.Indikator,
							Target:        targetResponses,
						})
					}
				}
			}

			// Indikator kegiatan
			var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse
			if rk.KodeKegiatan != "" {
				indikatorKegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
					ctx,
					tx,
					rk.KodeKegiatan,
					operational.KodeOpd,
					operational.Tahun,
				)
				if err == nil {
					for _, ind := range indikatorKegiatan {
						// Ambil target untuk indikator
						targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
						var targetResponses []pohonkinerja.TargetResponse
						if err == nil {
							for _, target := range targets {
								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
									Id:              target.Id,
									IndikatorId:     target.IndikatorId,
									TargetIndikator: target.Target,
									SatuanIndikator: target.Satuan,
								})
							}
						}

						indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
							Id:            ind.Id,
							Kode:          ind.Kode,
							NamaIndikator: ind.Indikator,
							Target:        targetResponses,
						})
					}
				}
			}
			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalResponse{
				Id:                   rk.Id,
				IdPohon:              operational.Id,
				NamaPohon:            operational.NamaPohon,
				NamaRencanaKinerja:   rk.NamaRencanaKinerja,
				Tahun:                operational.Tahun,
				PegawaiId:            rk.PegawaiId,
				NamaPegawai:          rk.NamaPegawai,
				KodeSubkegiatan:      rk.KodeSubKegiatan,
				NamaSubkegiatan:      rk.NamaSubKegiatan,
				Anggaran:             totalAnggaranRenkin,
				IndikatorSubkegiatan: indikatorSubkegiatanResponses,
				KodeKegiatan:         rk.KodeKegiatan,
				NamaKegiatan:         rk.NamaKegiatan,
				IndikatorKegiatan:    indikatorKegiatanResponses,
				Indikator:            indikatorRekinResponses,
			})
		}
	}

	operationalResp := pohonkinerja.OperationalCascadingOpdResponse{
		Id:         operational.Id,
		Parent:     operational.Parent,
		Strategi:   operational.NamaPohon,
		JenisPohon: operational.JenisPohon,
		LevelPohon: operational.LevelPohon,
		Keterangan: operational.Keterangan,
		Status:     operational.Status,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		},
		IsActive:       operational.IsActive,
		RencanaKinerja: rencanaKinerjaResponses,
		Indikator:      indikatorMap[operational.Id],
		TotalAnggaran:  totalAnggaranOperational,
	}

	// Build operational N responses jika ada
	nextLevel := operational.LevelPohon + 1
	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdCascadingResponse
		sort.Slice(operationalNList, func(i, j int) bool {
			return operationalNList[i].Id < operationalNList[j].Id
		})

		for _, opN := range operationalNList {
			childResp := service.buildOperationalNCascadingResponse(ctx, tx, pohonMap, opN, indikatorMap, rencanaKinerjaMap)
			childs = append(childs, childResp)
		}
		operationalResp.Childs = childs
	}

	return operationalResp
}

func (service *CascadingOpdServiceImpl) buildOperationalNCascadingResponse(
	ctx context.Context,
	tx *sql.Tx,
	pohonMap map[int]map[int][]domain.PohonKinerja,
	operationalN domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.OperationalNOpdCascadingResponse {

	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalNResponse
	if operationalN.PelaksanaIds != "" {
		pelaksanaIds := strings.Split(operationalN.PelaksanaIds, ",")
		existingPelaksana := make(map[string]bool)

		// Proses rencana kinerja yang ada
		if rencanaKinerjaList, ok := rencanaKinerjaMap[operationalN.Id]; ok {
			for _, rk := range rencanaKinerjaList {
				if rk.PegawaiId != "" {
					existingPelaksana[rk.PegawaiId] = true
					var namaPegawai string
					pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rk.PegawaiId)
					if err == nil {
						namaPegawai = pegawai.NamaPegawai
					}

					// Ambil indikator hanya jika ada rencana kinerja
					var indikatorResponses []pohonkinerja.IndikatorResponse
					indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
					if err == nil {
						for _, ind := range indikatorRekin {
							targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
							var targetResponses []pohonkinerja.TargetResponse
							if err == nil {
								for _, target := range targets {
									targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
										Id:              target.Id,
										IndikatorId:     target.IndikatorId,
										TargetIndikator: target.Target,
										SatuanIndikator: target.Satuan,
									})
								}
							}
							indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
								Id:            ind.Id,
								IdRekin:       rk.Id,
								NamaIndikator: ind.Indikator,
								Target:        targetResponses,
							})
						}
					}

					rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalNResponse{
						Id:                 rk.Id,
						IdPohon:            operationalN.Id,
						NamaPohon:          operationalN.NamaPohon,
						NamaRencanaKinerja: rk.NamaRencanaKinerja,
						Tahun:              operationalN.Tahun,
						PegawaiId:          rk.PegawaiId,
						NamaPegawai:        namaPegawai,
						Indikator:          indikatorResponses,
					})
				}
			}
		}

		// Tambahkan pelaksana yang tidak memiliki rencana kinerja
		for _, pelaksanaId := range pelaksanaIds {
			if !existingPelaksana[pelaksanaId] {
				pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksanaId)
				if err == nil {
					rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalNResponse{
						IdPohon:     operationalN.Id,
						NamaPohon:   operationalN.NamaPohon,
						Tahun:       operationalN.Tahun,
						PegawaiId:   pegawai.Nip,
						NamaPegawai: pegawai.NamaPegawai,
						// Field lain dibiarkan kosong
						Indikator: nil, // Pastikan indikator kosong
					})
				}
			}
		}
	}

	operationalNResp := pohonkinerja.OperationalNOpdCascadingResponse{
		Id:         operationalN.Id,
		Parent:     operationalN.Parent,
		Strategi:   operationalN.NamaPohon,
		JenisPohon: operationalN.JenisPohon,
		LevelPohon: operationalN.LevelPohon,
		Keterangan: operationalN.Keterangan,
		Status:     operationalN.Status,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		},
		IsActive:       operationalN.IsActive,
		RencanaKinerja: rencanaKinerjaResponses,
		Indikator:      indikatorMap[operationalN.Id],
	}

	// Build child responses untuk level berikutnya (N+1)
	nextLevel := operationalN.LevelPohon + 1
	if childList := pohonMap[nextLevel][operationalN.Id]; len(childList) > 0 {
		var childs []pohonkinerja.OperationalNOpdCascadingResponse
		sort.Slice(childList, func(i, j int) bool {
			return childList[i].Id < childList[j].Id
		})

		for _, child := range childList {
			// Rekursif untuk membangun response child
			childResp := service.buildOperationalNCascadingResponse(
				ctx,
				tx,
				pohonMap,
				child,
				indikatorMap,
				rencanaKinerjaMap,
			)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}
