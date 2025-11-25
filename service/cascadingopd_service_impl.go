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
	cascadingOpdRepository   repository.CascadingOpdRepository
	rincianBelanjaRepository repository.RincianBelanjaRepository
	rencanaAksiRepository    repository.RencanaAksiRepository
	urusanRepository         repository.UrusanRepository
	bidangUrusanRepository   repository.BidangUrusanRepository
	programRepository        repository.ProgramRepository
	kegiatanRepository       repository.KegiatanRepository
	subKegiatanRepository    repository.SubKegiatanRepository
}

func NewCascadingOpdServiceImpl(
	pohonKinerjaRepository repository.PohonKinerjaRepository,
	opdRepository repository.OpdRepository,
	pegawaiRepository repository.PegawaiRepository,
	tujuanOpdRepository repository.TujuanOpdRepository,
	rencanaKinerjaRepository repository.RencanaKinerjaRepository,
	DB *sql.DB,
	cascadingOpdRepository repository.CascadingOpdRepository,
	rincianBelanjaRepository repository.RincianBelanjaRepository,
	rencanaAksiRepository repository.RencanaAksiRepository,
	urusanRepository repository.UrusanRepository,
	bidangUrusanRepository repository.BidangUrusanRepository,
	programRepository repository.ProgramRepository,
	kegiatanRepository repository.KegiatanRepository,
	subKegiatanRepository repository.SubKegiatanRepository,
) *CascadingOpdServiceImpl {
	return &CascadingOpdServiceImpl{
		pohonKinerjaRepository:   pohonKinerjaRepository,
		opdRepository:            opdRepository,
		pegawaiRepository:        pegawaiRepository,
		tujuanOpdRepository:      tujuanOpdRepository,
		rencanaKinerjaRepository: rencanaKinerjaRepository,
		DB:                       DB,
		cascadingOpdRepository:   cascadingOpdRepository,
		rincianBelanjaRepository: rincianBelanjaRepository,
		rencanaAksiRepository:    rencanaAksiRepository,
		urusanRepository:         urusanRepository,
		bidangUrusanRepository:   bidangUrusanRepository,
		programRepository:        programRepository,
		kegiatanRepository:       kegiatanRepository,
		subKegiatanRepository:    subKegiatanRepository,
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

		// Ambil data rencana kinerja untuk semua level
		rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, p.Id)
		if err == nil {
			var validRencanaKinerja []domain.RencanaKinerja

			// Tambahkan log untuk debugging
			log.Printf("Processing rencana kinerja for pohon ID %d (Level %d)", p.Id, p.LevelPohon)

			// Ambil daftar pelaksana
			pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(p.Id))
			if err == nil {
				pelaksanaMap := make(map[string]*domainmaster.Pegawai)
				rekinMap := make(map[string]bool) // Ubah dari PegawaiId ke Id rencana kinerja untuk deduplication

				for _, pelaksana := range pelaksanaList {
					pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
					if err == nil {
						pelaksanaMap[pegawai.Nip] = &pegawai
					}
				}

				for _, rk := range rencanaKinerjaList {
					if pegawai, exists := pelaksanaMap[rk.PegawaiId]; exists {
						// Cek apakah rencana kinerja ini sudah ditambahkan (deduplication)
						if !rekinMap[rk.Id] {
							rk.NamaPegawai = pegawai.NamaPegawai
							validRencanaKinerja = append(validRencanaKinerja, rk)
							rekinMap[rk.Id] = true // Simpan ID rencana kinerja, bukan PegawaiId
						}
					}
				}

				// Track pegawai yang sudah memiliki rencana kinerja
				pegawaiWithRekinMap := make(map[string]bool)
				for _, rk := range validRencanaKinerja {
					pegawaiWithRekinMap[rk.PegawaiId] = true
				}

				// Tambahkan pelaksana yang belum memiliki rencana kinerja
				for _, pegawai := range pelaksanaMap {
					if !pegawaiWithRekinMap[pegawai.Nip] {
						validRencanaKinerja = append(validRencanaKinerja, domain.RencanaKinerja{
							IdPohon:     p.Id,
							NamaPohon:   p.NamaPohon,
							Tahun:       p.Tahun,
							PegawaiId:   pegawai.Nip,
							NamaPegawai: pegawai.NamaPegawai,
							Indikator:   nil,
						})
					}
				}
			}

			rencanaKinerjaMap[p.Id] = validRencanaKinerja
			log.Printf("Added %d valid rencana kinerja for pohon ID %d", len(validRencanaKinerja), p.Id)
		}
	}

	// Build response untuk strategic (level 4)
	if strategicList := pohonMap[4]; len(strategicList) > 0 {
		log.Printf("Processing %d strategic entries", len(strategicList))
		for _, strategicsByParent := range strategicList {
			sort.Slice(strategicsByParent, func(i, j int) bool {
				// Prioritaskan status "pokin dari pemda"
				if strategicsByParent[i].Status == "pokin dari pemda" && strategicsByParent[j].Status != "pokin dari pemda" {
					return true
				}
				if strategicsByParent[i].Status != "pokin dari pemda" && strategicsByParent[j].Status == "pokin dari pemda" {
					return false
				}
				return strategicsByParent[i].Id < strategicsByParent[j].Id
			})

			for _, strategic := range strategicsByParent {
				strategicResp := service.buildStrategicCascadingResponse(ctx, tx, pohonMap, strategic, indikatorMap, rencanaKinerjaMap)
				response.Strategics = append(response.Strategics, strategicResp)
			}
		}

		sort.Slice(response.Strategics, func(i, j int) bool {
			if response.Strategics[i].Status == "pokin dari pemda" && response.Strategics[j].Status != "pokin dari pemda" {
				return true
			}
			if response.Strategics[i].Status != "pokin dari pemda" && response.Strategics[j].Status == "pokin dari pemda" {
				return false
			}
			return response.Strategics[i].Id < response.Strategics[j].Id
		})
	}

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
			// Prioritaskan status "pokin dari pemda"
			if tacticalList[i].Status == "pokin dari pemda" && tacticalList[j].Status != "pokin dari pemda" {
				return true
			}
			if tacticalList[i].Status != "pokin dari pemda" && tacticalList[j].Status == "pokin dari pemda" {
				return false
			}
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
			// Prioritaskan status "pokin dari pemda"
			if operationalList[i].Status == "pokin dari pemda" && operationalList[j].Status != "pokin dari pemda" {
				return true
			}
			if operationalList[i].Status != "pokin dari pemda" && operationalList[j].Status == "pokin dari pemda" {
				return false
			}
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

	log.Printf("Building OperationalN response for ID: %d, Level: %d", operationalN.Id, operationalN.LevelPohon)

	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalNResponse

	// Proses rencana kinerja yang ada
	if rencanaKinerjaList, ok := rencanaKinerjaMap[operationalN.Id]; ok {
		log.Printf("Found %d rencana kinerja for OperationalN ID %d", len(rencanaKinerjaList), operationalN.Id)

		for _, rk := range rencanaKinerjaList {
			var indikatorResponses []pohonkinerja.IndikatorResponse

			// Ambil indikator jika rencana kinerja memiliki ID
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

			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalNResponse{
				Id:                 rk.Id,
				IdPohon:            operationalN.Id,
				NamaPohon:          operationalN.NamaPohon,
				NamaRencanaKinerja: rk.NamaRencanaKinerja,
				Tahun:              operationalN.Tahun,
				PegawaiId:          rk.PegawaiId,
				NamaPegawai:        rk.NamaPegawai,
				Indikator:          indikatorResponses,
			})
		}
	}

	// Buat response
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

	// Proses child nodes jika ada
	nextLevel := operationalN.LevelPohon + 1
	if childList := pohonMap[nextLevel][operationalN.Id]; len(childList) > 0 {
		var childs []pohonkinerja.OperationalNOpdCascadingResponse
		sort.Slice(childList, func(i, j int) bool {
			return childList[i].Id < childList[j].Id
		})

		for _, child := range childList {
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

// by rekin
func (service *CascadingOpdServiceImpl) FindByRekinPegawaiAndId(ctx context.Context, rekinId string) (pohonkinerja.CascadingRekinPegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Cari pohon kinerja berdasarkan rencana kinerja ID
	pokinRekin, err := service.cascadingOpdRepository.FindPokinByRekinId(ctx, tx, rekinId)
	if err != nil {
		log.Printf("Error: Pohon kinerja not found for rekin_id=%s: %v", rekinId, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak ditemukan")
	}

	// 2. Ambil data rencana kinerja
	rekin, err := service.rencanaKinerjaRepository.FindById(ctx, tx, rekinId, "", "")
	if err != nil {
		log.Printf("Error: Rencana kinerja not found: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak ditemukan")
	}

	// 3. Validasi pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rekin.PegawaiId)
	if err != nil {
		log.Printf("Error: Pegawai not found for NIP=%s: %v", rekin.PegawaiId, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("pegawai tidak ditemukan")
	}

	// 4. CHECK: Apakah pegawai adalah pelaksana di pohon ini?
	isPelaksana, err := service.isPegawaiPelaksanaPokin(ctx, tx, pokinRekin.Id, pegawai.Id)
	if err != nil {
		log.Printf("Error checking pelaksana: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}

	if !isPelaksana {
		log.Printf("Error: Pegawai ID %s (NIP %s) bukan pelaksana di pohon %d", pegawai.Id, pegawai.Nip, pokinRekin.Id)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak dipakai di cascading opd")
	}

	log.Printf("Pegawai ID %s is pelaksana of pokin %d", pegawai.Id, pokinRekin.Id)

	// 5. Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokinRekin.KodeOpd)
	if err != nil {
		log.Printf("Error: OPD not found for kode_opd=%s: %v", pokinRekin.KodeOpd, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// 6. Hitung total anggaran berdasarkan level - DENGAN FILTER PELAKSANA
	var totalAnggaran int64

	if pokinRekin.LevelPohon == 6 {
		totalAnggaran, err = service.calculateAnggaranForOperationalWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for operational: %v", err)
			totalAnggaran = 0
		}
	} else if pokinRekin.LevelPohon == 5 {
		totalAnggaran, err = service.calculateAnggaranForTacticalWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("ERROR: Failed to calculate anggaran for tactical: %v", err)
			totalAnggaran = 0
		} else {
			log.Printf("SUCCESS: Total anggaran for tactical ID %d = %d", pokinRekin.Id, totalAnggaran)
		}
	} else if pokinRekin.LevelPohon == 4 {
		totalAnggaran, err = service.calculateAnggaranForStrategicWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for strategic: %v", err)
			totalAnggaran = 0
		}
	} else {
		totalAnggaran = 0
	}

	// 7. Build response
	response := service.buildCascadingRekinResponse(ctx, tx, pokinRekin, opd, totalAnggaran)

	return response, nil
}

// Build response spesifik untuk cascading rekin pegawai
func (service *CascadingOpdServiceImpl) buildCascadingRekinResponse(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja,
	opd domainmaster.Opd,
	totalAnggaran int64) pohonkinerja.CascadingRekinPegawaiResponse {

	pokin.NamaOpd = opd.NamaOpd

	// Set parent sesuai level
	var parent *int
	if pokin.LevelPohon > 4 {
		parent = &pokin.Parent
	}

	// Inisialisasi response dasar
	response := pohonkinerja.CascadingRekinPegawaiResponse{
		Id:         pokin.Id,
		Parent:     parent,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		Keterangan: pokin.Keterangan,
		Status:     pokin.Status,
		PerangkatDaerah: opdmaster.OpdResponseForAll{
			KodeOpd: pokin.KodeOpd,
			NamaOpd: pokin.NamaOpd,
		},
		IsActive:     pokin.IsActive,
		PaguAnggaran: totalAnggaran,
	}

	// Ambil rencana kinerja di pohon ini
	response.RencanaKinerja = service.getRencanaKinerjaForRekin(ctx, tx, pokin)

	// Berdasarkan level, ambil program ATAU kegiatan/subkegiatan
	if pokin.LevelPohon == 4 || pokin.LevelPohon == 5 {
		// Level 4 & 5: Tampilkan program
		response.Program = service.getProgramForRekin(ctx, tx, pokin)
	} else if pokin.LevelPohon == 6 {
		// Level 6: Tampilkan kegiatan dan subkegiatan
		response.Kegiatan, response.SubKegiatan = service.getKegiatanSubkegiatanForRekin(ctx, tx, pokin)
	}

	return response
}

// Get rencana kinerja untuk response
func (service *CascadingOpdServiceImpl) getRencanaKinerjaForRekin(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja) []pohonkinerja.RencanaKinerjaResponse {

	var rencanaKinerjaList []pohonkinerja.RencanaKinerjaResponse

	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokin.Id)
	if err != nil {
		return rencanaKinerjaList
	}

	// Ambil pelaksana
	pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokin.Id))
	if err != nil {
		return rencanaKinerjaList
	}

	pelaksanaMap := make(map[string]*domainmaster.Pegawai)
	for _, pelaksana := range pelaksanaList {
		pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
		if err == nil {
			pelaksanaMap[pegawai.Nip] = &pegawai
		}
	}

	rekinMap := make(map[string]bool)

	for _, rk := range rekinList {

		if rekinMap[rk.Id] {
			log.Printf("[DEBUG] Skip duplicate rekin ID: %s", rk.Id)
			continue
		}

		if pegawai, exists := pelaksanaMap[rk.PegawaiId]; exists {
			var indikatorResponses []pohonkinerja.IndikatorResponse
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

			rencanaKinerjaList = append(rencanaKinerjaList, pohonkinerja.RencanaKinerjaResponse{
				Id:                 rk.Id,
				IdPohon:            pokin.Id,
				NamaPohon:          pokin.NamaPohon,
				NamaRencanaKinerja: rk.NamaRencanaKinerja,
				Tahun:              pokin.Tahun,
				PegawaiId:          rk.PegawaiId,
				NamaPegawai:        pegawai.NamaPegawai,
				Indikator:          indikatorResponses,
			})

			rekinMap[rk.Id] = true
		}
	}

	return rencanaKinerjaList
}

func (service *CascadingOpdServiceImpl) getKegiatanSubkegiatanForRekin(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja) ([]pohonkinerja.KegiatanCascadingRekinResponse, []pohonkinerja.SubKegiatanCascadingRekinResponse) {

	var kegiatanList []pohonkinerja.KegiatanCascadingRekinResponse
	var subkegiatanList []pohonkinerja.SubKegiatanCascadingRekinResponse

	// Ambil rencana kinerja dari pohon ini
	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokin.Id)
	if err != nil {
		log.Printf("Error getting rencana kinerja: %v", err)
		return kegiatanList, subkegiatanList
	}

	// Map untuk kegiatan dan subkegiatan unik
	kegiatanMap := make(map[string]string)
	subkegiatanMap := make(map[string]string)

	// Loop rencana kinerja dan kumpulkan kode kegiatan & subkegiatan
	for _, rk := range rekinList {
		// Kumpulkan kegiatan dari field KodeKegiatan
		if rk.KodeKegiatan != "" {
			if _, exists := kegiatanMap[rk.KodeKegiatan]; !exists {
				kegiatanMap[rk.KodeKegiatan] = rk.NamaKegiatan
			}
		}

		// Kumpulkan subkegiatan dari field KodeSubKegiatan
		if rk.KodeSubKegiatan != "" {
			if _, exists := subkegiatanMap[rk.KodeSubKegiatan]; !exists {
				subkegiatanMap[rk.KodeSubKegiatan] = rk.NamaSubKegiatan
			}
		}
	}

	// Build kegiatan list dengan indikator
	for kodeKegiatan, namaKegiatan := range kegiatanMap {
		var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse

		// Ambil indikator kegiatan (sama seperti di FindAll)
		indikatorKegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx, tx, kodeKegiatan, pokin.KodeOpd, pokin.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorKegiatan {
				targets, _ := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range targets {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		}

		kegiatanList = append(kegiatanList, pohonkinerja.KegiatanCascadingRekinResponse{
			KodeKegiatan: kodeKegiatan,
			NamaKegiatan: namaKegiatan,
			Indikator:    indikatorKegiatanResponses,
		})
	}

	// Build subkegiatan list dengan indikator
	for kodeSubkegiatan, namaSubkegiatan := range subkegiatanMap {
		var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse

		// Ambil indikator subkegiatan (sama seperti di FindAll)
		indikatorSubkegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx, tx, kodeSubkegiatan, pokin.KodeOpd, pokin.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorSubkegiatan {
				targets, _ := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range targets {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		}

		subkegiatanList = append(subkegiatanList, pohonkinerja.SubKegiatanCascadingRekinResponse{
			KodeSubkegiatan: kodeSubkegiatan,
			NamaSubkegiatan: namaSubkegiatan,
			Indikator:       indikatorSubkegiatanResponses,
		})
	}

	return kegiatanList, subkegiatanList
}

func (service *CascadingOpdServiceImpl) getProgramForRekin(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja) []pohonkinerja.ProgramCascadingRekinResponse {

	var programList []pohonkinerja.ProgramCascadingRekinResponse

	log.Printf("[DEBUG] Getting program for Level %d, Pokin ID %d", pokin.LevelPohon, pokin.Id)

	if pokin.LevelPohon != 4 && pokin.LevelPohon != 5 {
		log.Printf("[DEBUG] Level %d tidak memerlukan program", pokin.LevelPohon)
		return programList
	}

	// Map untuk menyimpan program unik
	programMap := make(map[string]string)

	// Ambil operational children (level 6) - SAMA SEPERTI FindAll
	var operationalIds []int
	var err error

	if pokin.LevelPohon == 5 {
		// Tactical: ambil operational children langsung
		operationalIds, err = service.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, pokin.Id)
		if err != nil {
			log.Printf("[ERROR] Failed to get operational children: %v", err)
			return programList
		}
		log.Printf("[DEBUG] Found %d operational children for tactical", len(operationalIds))
	} else if pokin.LevelPohon == 4 {
		// Strategic: ambil tactical children dulu, lalu operational dari setiap tactical
		tacticalIds, err := service.cascadingOpdRepository.FindTacticalChildrenByStrategicId(ctx, tx, pokin.Id)
		if err != nil {
			log.Printf("[ERROR] Failed to get tactical children: %v", err)
			return programList
		}
		log.Printf("[DEBUG] Found %d tactical children for strategic", len(tacticalIds))

		// Untuk setiap tactical, ambil operational children-nya
		for _, tacticalId := range tacticalIds {
			ops, err := service.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, tacticalId)
			if err == nil {
				operationalIds = append(operationalIds, ops...)
			}
		}
		log.Printf("[DEBUG] Total %d operational found from all tacticals", len(operationalIds))
	}

	// Untuk setiap operational, ambil rencana kinerja dan extract kode program
	// PERSIS SEPERTI DI FindAll line 488-507
	for _, operationalId := range operationalIds {
		rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, operationalId)
		if err == nil {
			for _, rk := range rencanaKinerjaList {
				if rk.KodeSubKegiatan != "" {
					segments := strings.Split(rk.KodeSubKegiatan, ".")
					if len(segments) >= 3 {
						kodeProgram := strings.Join(segments[:3], ".")
						if _, exists := programMap[kodeProgram]; !exists {
							program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
							if err == nil {
								programMap[kodeProgram] = program.NamaProgram
								log.Printf("[DEBUG] Found program: %s - %s", kodeProgram, program.NamaProgram)
							} else {
								log.Printf("[ERROR] Program not found for kode: %s, error: %v", kodeProgram, err)
							}
						}
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] Total unique programs found: %d", len(programMap))

	// Build program list dengan indikator - SAMA SEPERTI FindAll
	for kodeProgram, namaProgram := range programMap {
		var indikatorProgramResponses []pohonkinerja.IndikatorResponse

		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx, tx, kodeProgram, pokin.KodeOpd, pokin.Tahun,
		)
		if err == nil {
			log.Printf("[DEBUG] Found %d indikator for program %s", len(indikatorProgram), kodeProgram)
			for _, ind := range indikatorProgram {
				targets, _ := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range targets {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		} else {
			log.Printf("[WARN] No indikator found for program %s: %v", kodeProgram, err)
		}

		programList = append(programList, pohonkinerja.ProgramCascadingRekinResponse{
			KodeProgram: kodeProgram,
			NamaProgram: namaProgram,
			Indikator:   indikatorProgramResponses,
		})
	}

	log.Printf("[DEBUG] Returning %d programs", len(programList))
	return programList
}

// Logic: Hitung anggaran untuk level 6 (Operational)
func (service *CascadingOpdServiceImpl) calculateAnggaranForOperational(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	log.Printf("[DEBUG] Calculating anggaran for Operational ID: %d", pokinId)

	// Ambil total anggaran langsung dari repository
	totalAnggaran, err := service.cascadingOpdRepository.GetTotalAnggaranByPokinId(ctx, tx, pokinId)
	if err != nil {
		log.Printf("[ERROR] Failed to get anggaran for operational %d: %v", pokinId, err)
		return 0, err
	}

	log.Printf("[DEBUG] Operational ID %d total anggaran from DB: %d", pokinId, totalAnggaran)

	return totalAnggaran, nil
}

// Logic: Hitung anggaran untuk level 5 (Tactical) - sum dari semua operational di bawahnya
func (service *CascadingOpdServiceImpl) calculateAnggaranForTactical(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	var totalAnggaran int64 = 0

	log.Printf("[DEBUG] Calculating anggaran for Tactical ID: %d", pokinId)

	// Ambil semua operational children
	operationalIds, err := service.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, pokinId)
	if err != nil {
		log.Printf("[ERROR] Failed to find operational children for tactical %d: %v", pokinId, err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d operational children for tactical %d", len(operationalIds), pokinId)

	// Hitung anggaran untuk setiap operational
	for _, opId := range operationalIds {
		anggaran, err := service.calculateAnggaranForOperational(ctx, tx, opId)
		if err != nil {
			log.Printf("[ERROR] Failed to calculate anggaran for operational %d: %v", opId, err)
		} else {
			log.Printf("[DEBUG] Operational ID %d has anggaran: %d", opId, anggaran)
			totalAnggaran += anggaran
		}
	}

	log.Printf("[DEBUG] Total anggaran for Tactical ID %d: %d", pokinId, totalAnggaran)

	return totalAnggaran, nil
}

// Logic: Hitung anggaran untuk level 4 (Strategic) - sum dari semua tactical di bawahnya
func (service *CascadingOpdServiceImpl) calculateAnggaranForStrategic(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	var totalAnggaran int64 = 0

	// Ambil semua tactical children
	tacticalIds, err := service.cascadingOpdRepository.FindTacticalChildrenByStrategicId(ctx, tx, pokinId)
	if err != nil {
		return 0, err
	}

	// Hitung anggaran untuk setiap tactical
	for _, tactId := range tacticalIds {
		anggaran, err := service.calculateAnggaranForTactical(ctx, tx, tactId)
		if err == nil {
			totalAnggaran += anggaran
		}
	}

	return totalAnggaran, nil
}

// Logic: Collect programs untuk response flat
// func (service *CascadingOpdServiceImpl) collectProgramsForPokinFlat(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pokinId int,
// 	level int,
// 	programMap map[string]string) {

// 	var kodeList []string
// 	var err error

// 	// Jika level 4 atau 5, ambil dari children
// 	if level == 4 || level == 5 {
// 		kodeList, err = service.cascadingOpdRepository.FindKodeSubkegiatanFromChildren(ctx, tx, pokinId)
// 	} else if level == 6 {
// 		// Jika level 6, ambil dari pohon ini saja
// 		kodeList, err = service.cascadingOpdRepository.FindKodeSubkegiatanByPokinId(ctx, tx, pokinId)
// 	} else {
// 		return
// 	}

// 	if err != nil {
// 		log.Printf("Error collecting programs: %v", err)
// 		return
// 	}

// 	// Extract kode program dan ambil nama program
// 	for _, kodeSubkegiatan := range kodeList {
// 		segments := strings.Split(kodeSubkegiatan, ".")
// 		if len(segments) >= 3 {
// 			kodeProgram := strings.Join(segments[:3], ".")
// 			if _, exists := programMap[kodeProgram]; !exists {
// 				program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
// 				if err == nil {
// 					programMap[kodeProgram] = program.NamaProgram
// 				}
// 			}
// 		}
// 	}
// }

func (service *CascadingOpdServiceImpl) FindByIdPokin(ctx context.Context, pokinId int) (pohonkinerja.CascadingRekinPegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Cari pohon kinerja berdasarkan ID
	pokin, err := service.cascadingOpdRepository.FindPokinById(ctx, tx, pokinId)
	if err != nil {
		log.Printf("Error: Pohon kinerja not found for pokin_id=%d: %v", pokinId, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("pohon kinerja tidak ditemukan")
	}

	// 2. Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
	if err != nil {
		log.Printf("Error: OPD not found for kode_opd=%s: %v", pokin.KodeOpd, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// 3. Hitung total anggaran berdasarkan level - DENGAN FILTER PELAKSANA
	var totalAnggaran int64

	if pokin.LevelPohon == 6 {
		totalAnggaran, err = service.calculateAnggaranForOperationalWithPelaksana(ctx, tx, pokin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for operational: %v", err)
			totalAnggaran = 0
		}
	} else if pokin.LevelPohon == 5 {
		totalAnggaran, err = service.calculateAnggaranForTacticalWithPelaksana(ctx, tx, pokin.Id)
		if err != nil {
			log.Printf("ERROR: Failed to calculate anggaran for tactical: %v", err)
			totalAnggaran = 0
		} else {
			log.Printf("SUCCESS: Total anggaran for tactical ID %d = %d", pokin.Id, totalAnggaran)
		}
	} else if pokin.LevelPohon == 4 {
		totalAnggaran, err = service.calculateAnggaranForStrategicWithPelaksana(ctx, tx, pokin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for strategic: %v", err)
			totalAnggaran = 0
		}
	} else {
		totalAnggaran = 0
	}

	// 4. Build response
	response := service.buildCascadingRekinResponse(ctx, tx, pokin, opd, totalAnggaran)

	return response, nil
}

func (service *CascadingOpdServiceImpl) FindByNip(ctx context.Context, nip string, tahun string) ([]pohonkinerja.CascadingRekinPegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return []pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Validasi pegawai dan dapatkan ID pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, nip)
	if err != nil {
		log.Printf("Error: Pegawai not found for NIP=%s: %v", nip, err)
		return []pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("pegawai tidak ditemukan")
	}

	log.Printf("Found pegawai: ID=%s, NIP=%s, Nama=%s", pegawai.Id, pegawai.Nip, pegawai.NamaPegawai)

	// 2. Cari semua pohon kinerja yang memiliki rencana kinerja pegawai ini
	pokins, err := service.cascadingOpdRepository.FindPokinByNipAndTahun(ctx, tx, nip, tahun)
	if err != nil {
		log.Printf("Error: Failed to get pohon kinerja for NIP=%s: %v", nip, err)
		return []pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}

	if len(pokins) == 0 {
		log.Printf("No pohon kinerja found for NIP=%s, tahun=%s", nip, tahun)
		return []pohonkinerja.CascadingRekinPegawaiResponse{}, nil
	}

	log.Printf("Found %d pohon kinerja for NIP=%s", len(pokins), nip)

	// 3. Build response untuk setiap pohon kinerja
	var responses []pohonkinerja.CascadingRekinPegawaiResponse

	for _, pokin := range pokins {
		// CHECK: Apakah pegawai adalah pelaksana di pohon ini?
		isPelaksana, err := service.isPegawaiPelaksanaPokin(ctx, tx, pokin.Id, pegawai.Id)
		if err != nil {
			log.Printf("Error checking pelaksana for pokin %d: %v", pokin.Id, err)
			continue
		}

		if !isPelaksana {
			log.Printf("Skipping pohon ID %d - pegawai ID %s (NIP %s) bukan pelaksana", pokin.Id, pegawai.Id, nip)
			continue
		}

		log.Printf("Pohon ID %d - pegawai IS pelaksana, processing...", pokin.Id)

		// Validasi OPD
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err != nil {
			log.Printf("Warning: OPD not found for kode_opd=%s: %v", pokin.KodeOpd, err)
			continue
		}

		// Hitung total anggaran berdasarkan level
		var totalAnggaran int64

		if pokin.LevelPohon == 6 {
			totalAnggaran, err = service.calculateAnggaranForOperational(ctx, tx, pokin.Id)
			if err != nil {
				log.Printf("Warning: Failed to calculate anggaran for operational: %v", err)
				totalAnggaran = 0
			}
		} else if pokin.LevelPohon == 5 {
			totalAnggaran, err = service.calculateAnggaranForTactical(ctx, tx, pokin.Id)
			if err != nil {
				log.Printf("Warning: Failed to calculate anggaran for tactical: %v", err)
				totalAnggaran = 0
			}
		} else if pokin.LevelPohon == 4 {
			totalAnggaran, err = service.calculateAnggaranForStrategic(ctx, tx, pokin.Id)
			if err != nil {
				log.Printf("Warning: Failed to calculate anggaran for strategic: %v", err)
				totalAnggaran = 0
			}
		} else {
			totalAnggaran = 0
		}

		// Build response DENGAN FILTER NIP
		response := service.buildCascadingRekinResponseWithFilter(ctx, tx, pokin, opd, totalAnggaran, nip)

		// Double check: Pastikan rencana kinerja tidak kosong
		if len(response.RencanaKinerja) == 0 {
			log.Printf("Skipping pohon ID %d - no rencana kinerja after building response", pokin.Id)
			continue
		}

		responses = append(responses, response)
	}

	log.Printf("Returning %d pohon kinerja (after filtering) for NIP=%s", len(responses), nip)

	return responses, nil
}

// Check apakah pegawai adalah pelaksana di pohon kinerja ini
func (service *CascadingOpdServiceImpl) isPegawaiPelaksanaPokin(ctx context.Context, tx *sql.Tx, pokinId int, pegawaiId string) (bool, error) {
	// Ambil daftar pelaksana pokin
	pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokinId))
	if err != nil {
		return false, err
	}

	// Check apakah pegawaiId ada di list pelaksana
	for _, pelaksana := range pelaksanaList {
		if pelaksana.PegawaiId == pegawaiId {
			log.Printf("Found: Pegawai ID %s is pelaksana of pokin %d", pegawaiId, pokinId)
			return true, nil
		}
	}

	log.Printf("Not found: Pegawai ID %s is NOT pelaksana of pokin %d", pegawaiId, pokinId)
	return false, nil
}

// Build response dengan filter NIP untuk rencana kinerja, kegiatan, subkegiatan
func (service *CascadingOpdServiceImpl) buildCascadingRekinResponseWithFilter(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja,
	opd domainmaster.Opd,
	totalAnggaran int64,
	filterNip string) pohonkinerja.CascadingRekinPegawaiResponse {

	pokin.NamaOpd = opd.NamaOpd

	// Set parent sesuai level
	var parent *int
	if pokin.LevelPohon > 4 {
		parent = &pokin.Parent
	}

	// Inisialisasi response dasar
	response := pohonkinerja.CascadingRekinPegawaiResponse{
		Id:         pokin.Id,
		Parent:     parent,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		Keterangan: pokin.Keterangan,
		Status:     pokin.Status,
		PerangkatDaerah: opdmaster.OpdResponseForAll{
			KodeOpd: pokin.KodeOpd,
			NamaOpd: pokin.NamaOpd,
		},
		IsActive:     pokin.IsActive,
		PaguAnggaran: totalAnggaran,
	}

	// Ambil rencana kinerja dengan FILTER NIP
	response.RencanaKinerja = service.getRencanaKinerjaWithFilter(ctx, tx, pokin, filterNip)

	// Berdasarkan level, ambil program ATAU kegiatan/subkegiatan
	if pokin.LevelPohon == 4 || pokin.LevelPohon == 5 {
		// Level 4 & 5: Tampilkan program (semua program dari children)
		response.Program = service.getProgramForRekin(ctx, tx, pokin)
	} else if pokin.LevelPohon == 6 {
		// Level 6: Tampilkan kegiatan dan subkegiatan - HANYA dari rencana kinerja pegawai ini
		response.Kegiatan, response.SubKegiatan = service.getKegiatanSubkegiatanWithFilter(ctx, tx, pokin, filterNip)
	}

	return response
}

// Get rencana kinerja dengan filter NIP pegawai
func (service *CascadingOpdServiceImpl) getRencanaKinerjaWithFilter(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja,
	filterNip string) []pohonkinerja.RencanaKinerjaResponse {

	var rencanaKinerjaList []pohonkinerja.RencanaKinerjaResponse

	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokin.Id)
	if err != nil {
		return rencanaKinerjaList
	}

	// Loop dan FILTER hanya rencana kinerja milik pegawai ini
	for _, rk := range rekinList {
		// FILTER: Skip jika bukan rencana kinerja pegawai ini
		if rk.PegawaiId != filterNip {
			continue
		}

		// Ambil pegawai
		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rk.PegawaiId)
		if err != nil {
			log.Printf("Warning: Pegawai not found for NIP=%s", rk.PegawaiId)
			continue
		}

		var indikatorResponses []pohonkinerja.IndikatorResponse
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

		rencanaKinerjaList = append(rencanaKinerjaList, pohonkinerja.RencanaKinerjaResponse{
			Id:                 rk.Id,
			IdPohon:            pokin.Id,
			NamaPohon:          pokin.NamaPohon,
			NamaRencanaKinerja: rk.NamaRencanaKinerja,
			Tahun:              pokin.Tahun,
			PegawaiId:          rk.PegawaiId,
			NamaPegawai:        pegawai.NamaPegawai,
			Indikator:          indikatorResponses,
		})
	}

	return rencanaKinerjaList
}

// Get kegiatan dan subkegiatan dengan filter NIP pegawai
func (service *CascadingOpdServiceImpl) getKegiatanSubkegiatanWithFilter(
	ctx context.Context,
	tx *sql.Tx,
	pokin domain.PohonKinerja,
	filterNip string) ([]pohonkinerja.KegiatanCascadingRekinResponse, []pohonkinerja.SubKegiatanCascadingRekinResponse) {

	var kegiatanList []pohonkinerja.KegiatanCascadingRekinResponse
	var subkegiatanList []pohonkinerja.SubKegiatanCascadingRekinResponse

	// Ambil rencana kinerja dari pohon ini
	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokin.Id)
	if err != nil {
		log.Printf("Error getting rencana kinerja: %v", err)
		return kegiatanList, subkegiatanList
	}

	// Map untuk kegiatan dan subkegiatan unik
	kegiatanMap := make(map[string]string)
	subkegiatanMap := make(map[string]string)

	// Loop rencana kinerja dan kumpulkan kode kegiatan & subkegiatan - HANYA milik pegawai ini
	for _, rk := range rekinList {
		// FILTER: Skip jika bukan rencana kinerja pegawai ini
		if rk.PegawaiId != filterNip {
			continue
		}

		// Kumpulkan kegiatan dari field KodeKegiatan
		if rk.KodeKegiatan != "" {
			if _, exists := kegiatanMap[rk.KodeKegiatan]; !exists {
				kegiatanMap[rk.KodeKegiatan] = rk.NamaKegiatan
			}
		}

		// Kumpulkan subkegiatan dari field KodeSubKegiatan
		if rk.KodeSubKegiatan != "" {
			if _, exists := subkegiatanMap[rk.KodeSubKegiatan]; !exists {
				subkegiatanMap[rk.KodeSubKegiatan] = rk.NamaSubKegiatan
			}
		}
	}

	// Build kegiatan list dengan indikator
	for kodeKegiatan, namaKegiatan := range kegiatanMap {
		var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse

		indikatorKegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx, tx, kodeKegiatan, pokin.KodeOpd, pokin.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorKegiatan {
				targets, _ := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range targets {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		}

		kegiatanList = append(kegiatanList, pohonkinerja.KegiatanCascadingRekinResponse{
			KodeKegiatan: kodeKegiatan,
			NamaKegiatan: namaKegiatan,
			Indikator:    indikatorKegiatanResponses,
		})
	}

	// Build subkegiatan list dengan indikator
	for kodeSubkegiatan, namaSubkegiatan := range subkegiatanMap {
		var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse

		indikatorSubkegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
			ctx, tx, kodeSubkegiatan, pokin.KodeOpd, pokin.Tahun,
		)
		if err == nil {
			for _, ind := range indikatorSubkegiatan {
				targets, _ := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range targets {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					Kode:          ind.Kode,
					NamaIndikator: ind.Indikator,
					Target:        targetResponses,
				})
			}
		}

		subkegiatanList = append(subkegiatanList, pohonkinerja.SubKegiatanCascadingRekinResponse{
			KodeSubkegiatan: kodeSubkegiatan,
			NamaSubkegiatan: namaSubkegiatan,
			Indikator:       indikatorSubkegiatanResponses,
		})
	}

	return kegiatanList, subkegiatanList
}

// calculate anggaran
func (service *CascadingOpdServiceImpl) calculateAnggaranForOperationalWithPelaksana(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	log.Printf("[DEBUG] Calculating anggaran for Operational ID: %d (with pelaksana filter)", pokinId)

	var totalAnggaran int64 = 0

	// 1. Ambil semua rencana kinerja di pohon ini
	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokinId)
	if err != nil {
		log.Printf("[ERROR] Failed to get rencana kinerja: %v", err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d rencana kinerja for operational %d", len(rencanaKinerjaList), pokinId)

	// 2. Ambil daftar pelaksana - SAMA SEPERTI FindAll line 168
	pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokinId))
	if err != nil {
		log.Printf("[ERROR] Failed to get pelaksana: %v", err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d pelaksana for pokin %d", len(pelaksanaList), pokinId)

	// 3. Buat pelaksanaMap dengan NIP sebagai key - SAMA SEPERTI FindAll line 170-178
	pelaksanaMap := make(map[string]*domainmaster.Pegawai)
	for _, pelaksana := range pelaksanaList {
		pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
		if err == nil {
			pelaksanaMap[pegawai.Nip] = &pegawai
			log.Printf("[DEBUG] Pelaksana: ID=%s, NIP=%s, Nama=%s", pegawai.Id, pegawai.Nip, pegawai.NamaPegawai)
		}
	}

	// 4. Loop rencana kinerja dan hitung anggaran - HANYA jika pegawai adalah pelaksana
	for _, rk := range rencanaKinerjaList {
		// CHECK: Apakah pegawai ini ada di pelaksanaMap? - SAMA SEPERTI FindAll line 181
		if _, exists := pelaksanaMap[rk.PegawaiId]; !exists {
			log.Printf("[DEBUG] Skip rekin %s - pegawai %s bukan pelaksana", rk.Id, rk.PegawaiId)
			continue
		}

		log.Printf("[DEBUG] Processing rekin %s - pegawai %s IS pelaksana", rk.Id, rk.PegawaiId)

		// Hitung anggaran dari rencana kinerja ini - SAMA SEPERTI FindAll line 611-626
		var totalAnggaranRenkin int64 = 0
		if rk.Id != "" {
			rencanaAksiList, err := service.rencanaAksiRepository.FindAll(ctx, tx, rk.Id)
			if err == nil {
				for _, ra := range rencanaAksiList {
					rincianBelanja, err := service.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, ra.Id)
					if err == nil {
						totalAnggaranRenkin += rincianBelanja.Anggaran
					}
				}
			}
		}

		log.Printf("[DEBUG] Rekin %s anggaran: %d", rk.Id, totalAnggaranRenkin)
		totalAnggaran += totalAnggaranRenkin
	}

	log.Printf("[DEBUG] Operational ID %d total anggaran (pelaksana only): %d", pokinId, totalAnggaran)

	return totalAnggaran, nil
}

// Hitung anggaran level 5 DENGAN filter pelaksana
func (service *CascadingOpdServiceImpl) calculateAnggaranForTacticalWithPelaksana(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	var totalAnggaran int64 = 0

	log.Printf("[DEBUG] Calculating anggaran for Tactical ID: %d (with pelaksana filter)", pokinId)

	// Ambil semua operational children
	operationalIds, err := service.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, pokinId)
	if err != nil {
		log.Printf("[ERROR] Failed to find operational children for tactical %d: %v", pokinId, err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d operational children for tactical %d", len(operationalIds), pokinId)

	// Hitung anggaran untuk setiap operational - DENGAN FILTER PELAKSANA
	for _, opId := range operationalIds {
		anggaran, err := service.calculateAnggaranForOperationalWithPelaksana(ctx, tx, opId)
		if err != nil {
			log.Printf("[ERROR] Failed to calculate anggaran for operational %d: %v", opId, err)
		} else {
			log.Printf("[DEBUG] Operational ID %d has anggaran (pelaksana only): %d", opId, anggaran)
			totalAnggaran += anggaran
		}
	}

	log.Printf("[DEBUG] Total anggaran for Tactical ID %d (pelaksana only): %d", pokinId, totalAnggaran)

	return totalAnggaran, nil
}

// Hitung anggaran level 4 DENGAN filter pelaksana
func (service *CascadingOpdServiceImpl) calculateAnggaranForStrategicWithPelaksana(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	var totalAnggaran int64 = 0

	log.Printf("[DEBUG] Calculating anggaran for Strategic ID: %d (with pelaksana filter)", pokinId)

	// Ambil semua tactical children
	tacticalIds, err := service.cascadingOpdRepository.FindTacticalChildrenByStrategicId(ctx, tx, pokinId)
	if err != nil {
		log.Printf("[ERROR] Failed to find tactical children for strategic %d: %v", pokinId, err)
		return 0, err
	}

	log.Printf("[DEBUG] Found %d tactical children for strategic %d", len(tacticalIds), pokinId)

	// Hitung anggaran untuk setiap tactical - DENGAN FILTER PELAKSANA
	for _, tactId := range tacticalIds {
		anggaran, err := service.calculateAnggaranForTacticalWithPelaksana(ctx, tx, tactId)
		if err != nil {
			log.Printf("[ERROR] Failed to calculate anggaran for tactical %d: %v", tactId, err)
		} else {
			log.Printf("[DEBUG] Tactical ID %d has anggaran (pelaksana only): %d", tactId, anggaran)
			totalAnggaran += anggaran
		}
	}

	log.Printf("[DEBUG] Total anggaran for Strategic ID %d (pelaksana only): %d", pokinId, totalAnggaran)

	return totalAnggaran, nil
}

func (service *CascadingOpdServiceImpl) FindByMultipleRekinPegawai(ctx context.Context, request pohonkinerja.FindByMultipleRekinRequest) ([]pohonkinerja.CascadingRekinPegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if len(request.RekinIds) == 0 {
		return nil, errors.New("rekin_ids tidak boleh kosong")
	}

	var responses []pohonkinerja.CascadingRekinPegawaiResponse
	var errors []string

	// Loop untuk setiap rekin ID
	for _, rekinId := range request.RekinIds {
		if rekinId == "" {
			log.Printf("Warning: Skip empty rekin_id")
			continue
		}

		// Gunakan logika yang sama dengan FindByRekinPegawaiAndId
		response, err := service.processSingleRekin(ctx, tx, rekinId)
		if err != nil {
			log.Printf("Error processing rekin_id=%s: %v", rekinId, err)
			errors = append(errors, fmt.Sprintf("rekin_id %s: %s", rekinId, err.Error()))
			continue
		}

		responses = append(responses, response)
	}

	// Jika ada error dan tidak ada response yang berhasil, return error
	if len(responses) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("gagal memproses semua rencana kinerja: %v", errors)
	}

	// Log summary
	log.Printf("Successfully processed %d/%d rencana kinerja", len(responses), len(request.RekinIds))
	if len(errors) > 0 {
		log.Printf("Errors encountered: %v", errors)
	}

	return responses, nil
}

// Helper function untuk memproses single rekin (extracted from FindByRekinPegawaiAndId)
func (service *CascadingOpdServiceImpl) processSingleRekin(
	ctx context.Context,
	tx *sql.Tx,
	rekinId string,
) (pohonkinerja.CascadingRekinPegawaiResponse, error) {
	// 1. Cari pohon kinerja berdasarkan rencana kinerja ID
	pokinRekin, err := service.cascadingOpdRepository.FindPokinByRekinId(ctx, tx, rekinId)
	if err != nil {
		log.Printf("Error: Pohon kinerja not found for rekin_id=%s: %v", rekinId, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak ditemukan")
	}

	// 2. Ambil data rencana kinerja
	rekin, err := service.rencanaKinerjaRepository.FindById(ctx, tx, rekinId, "", "")
	if err != nil {
		log.Printf("Error: Rencana kinerja not found: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak ditemukan")
	}

	// 3. Validasi pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rekin.PegawaiId)
	if err != nil {
		log.Printf("Error: Pegawai not found for NIP=%s: %v", rekin.PegawaiId, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("pegawai tidak ditemukan")
	}

	// 4. CHECK: Apakah pegawai adalah pelaksana di pohon ini?
	isPelaksana, err := service.isPegawaiPelaksanaPokin(ctx, tx, pokinRekin.Id, pegawai.Id)
	if err != nil {
		log.Printf("Error checking pelaksana: %v", err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, err
	}

	if !isPelaksana {
		log.Printf("Error: Pegawai ID %s (NIP %s) bukan pelaksana di pohon %d", pegawai.Id, pegawai.Nip, pokinRekin.Id)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("rencana kinerja tidak dipakai di cascading opd")
	}

	log.Printf("Pegawai ID %s is pelaksana of pokin %d", pegawai.Id, pokinRekin.Id)

	// 5. Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokinRekin.KodeOpd)
	if err != nil {
		log.Printf("Error: OPD not found for kode_opd=%s: %v", pokinRekin.KodeOpd, err)
		return pohonkinerja.CascadingRekinPegawaiResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// 6. Hitung total anggaran berdasarkan level - DENGAN FILTER PELAKSANA
	var totalAnggaran int64

	if pokinRekin.LevelPohon == 6 {
		totalAnggaran, err = service.calculateAnggaranForOperationalWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for operational: %v", err)
			totalAnggaran = 0
		}
	} else if pokinRekin.LevelPohon == 5 {
		totalAnggaran, err = service.calculateAnggaranForTacticalWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("ERROR: Failed to calculate anggaran for tactical: %v", err)
			totalAnggaran = 0
		} else {
			log.Printf("SUCCESS: Total anggaran for tactical ID %d = %d", pokinRekin.Id, totalAnggaran)
		}
	} else if pokinRekin.LevelPohon == 4 {
		totalAnggaran, err = service.calculateAnggaranForStrategicWithPelaksana(ctx, tx, pokinRekin.Id)
		if err != nil {
			log.Printf("Warning: Failed to calculate anggaran for strategic: %v", err)
			totalAnggaran = 0
		}
	} else {
		totalAnggaran = 0
	}

	// 7. Build response
	response := service.buildCascadingRekinResponse(ctx, tx, pokinRekin, opd, totalAnggaran)

	return response, nil
}

func (service *CascadingOpdServiceImpl) MultiRekinDetails(
	ctx context.Context,
	request pohonkinerja.FindByMultipleRekinRequest,
) ([]pohonkinerja.DetailRekinResponse, error) {

	if len(request.RekinIds) == 0 {
		return nil, errors.New("rekin_ids tidak boleh kosong")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Ambil semua rekin detail
	detailRekins, err := service.rencanaKinerjaRepository.FindDetailRekins(ctx, tx, request.RekinIds)
	if err != nil {
		return nil, err
	}

	// 1.5 ambil indikator target
	indikatorMap, err := service.preloadIndikatorTargetRekin(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 2. Ambil anggaran
	anggaranMap, err := service.preloadAnggaran(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 3. Preload seluruh hierarki
	kegMap, subMap, progMap, bidMap, err := service.preloadAllHierarchies(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 4. Susun final response
	responses := service.fillResponseLoop(
		detailRekins,
		indikatorMap,
		anggaranMap,
		kegMap,
		subMap,
		progMap,
		bidMap,
	)

	return responses, nil
}

func (service *CascadingOpdServiceImpl) MultiRekinDetailsByOpdTahun(
	ctx context.Context,
	request pohonkinerja.MultiRekinDetailsByOpdAndTahunRequest,
) ([]pohonkinerja.DetailRekinResponse, error) {

	if request.KodeOpd == "" || request.Tahun == "" {
		return nil, errors.New("Kode opd atau tahun tidak ditemukan")
	}

	// if len(rekinIds) == 0 {
	// 	return nil, errors.New("rekin_ids tidak boleh kosong")
	// }

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Ambil semua rekin detail
	detailRekins, err := service.rencanaKinerjaRepository.FindDetailRekinsByOpdAndTahun(ctx, tx, request.KodeOpd, request.Tahun)
	if err != nil {
		return nil, err
	}

	// 1.5 ambil indikator target
	indikatorMap, err := service.preloadIndikatorTargetRekin(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 2. Ambil anggaran
	anggaranMap, err := service.preloadAnggaran(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 3. Preload seluruh hierarki
	kegMap, subMap, progMap, bidMap, err := service.preloadAllHierarchies(ctx, tx, detailRekins)
	if err != nil {
		return nil, err
	}

	// 4. Susun final response
	responses := service.fillResponseLoop(
		detailRekins,
		indikatorMap,
		anggaranMap,
		kegMap,
		subMap,
		progMap,
		bidMap,
	)

	return responses, nil
}

func (service *CascadingOpdServiceImpl) fillResponseLoop(
	rks []domain.DetailRekins,
	indikatorMap map[string][]domain.Indikator,
	anggaranMap map[string]int64,
	kegMap map[string]pohonkinerja.KegiatanRekinResponse,
	subMap map[string]pohonkinerja.SubKegiatanRekinResponse,
	progMap map[string][]pohonkinerja.ProgramRekinResponse,
	bidMap map[string][]pohonkinerja.BidangUrusanCascadingRekinResponse,
) []pohonkinerja.DetailRekinResponse {

	out := make([]pohonkinerja.DetailRekinResponse, len(rks))

	for i := range rks {
		rk := &rks[i]

		resp := pohonkinerja.DetailRekinResponse{
			Id:                 rk.Id,
			IdPohon:            rk.IdPohon,
			NamaRencanaKinerja: rk.NamaRencanaKinerja,
			Tahun:              rk.Tahun,
			PegawaiId:          rk.PegawaiId,
			NamaPegawai:        rk.NamaPegawai,
			KodeOpd:            rk.KodeOpd,
			LevelPohon:         rk.LevelPohon,
			PaguAnggaranRekin:  anggaranMap[rk.Id],
		}

		if inds, ok := indikatorMap[rk.Id]; ok {
			resp.Indikator = toIndikatorResponse(inds)
		}

		switch rk.LevelPohon {

		case 6:
			resp.Kegiatan = append(resp.Kegiatan, kegMap[rk.KodeSubKegiatan])
			resp.SubKegiatan = append(resp.SubKegiatan, subMap[rk.KodeSubKegiatan])

		case 5:
			resp.Program = append(resp.Program, progMap[rk.Id]...)

		case 4:
			resp.BidangUrusan = append(resp.BidangUrusan, bidMap[rk.Id]...)
		}

		out[i] = resp
	}
	return out
}

func (service *CascadingOpdServiceImpl) preloadIndikatorTargetRekin(
	ctx context.Context,
	tx *sql.Tx,
	rks []domain.DetailRekins,
) (map[string][]domain.Indikator, error) {

	allIds := make([]string, len(rks))
	for i := range rks {
		allIds[i] = rks[i].Id
	}

	list, err := service.rencanaKinerjaRepository.GetAllIndikatorTargetByRekinIds(ctx, tx, allIds)
	if err != nil {
		return nil, err
	}

	results := make(map[string][]domain.Indikator)

	for _, ind := range list {
		rkId := ind.RencanaKinerjaId
		results[rkId] = append(results[rkId], ind)
	}
	return results, nil
}

func (service *CascadingOpdServiceImpl) preloadAnggaran(
	ctx context.Context,
	tx *sql.Tx,
	rks []domain.DetailRekins,
) (map[string]int64, error) {

	// ---- Ambil anggaran mentah (level 6) ----
	rawAnggaran := make(map[string]int64)
	{
		allIds := make([]string, len(rks))
		for i := range rks {
			allIds[i] = rks[i].Id
		}

		list, err := service.rincianBelanjaRepository.FindAnggaranByRekinIds(ctx, tx, allIds)
		if err != nil {
			return nil, err
		}

		for _, a := range list {
			rawAnggaran[a.RekinId] = a.TotalAnggaran
		}
	}

	// ---- Siapkan map anak2 untuk roll-up ----
	level5map := make(map[string][]string) // rekin5 -> anak2 level6
	level4map := make(map[string][]string) // rekin4 -> anak2 level5

	for _, rk := range rks {

		switch rk.LevelPohon {

		case 5:
			// rekin level 5 anak level 6 langsung di table pohon_kinerja
			children, _ := service.pohonKinerjaRepository.FindChildPokins(ctx, tx, int64(rk.IdPohon))
			pokinIds := make([]int, len(children))
			for i := range children {
				pokinIds[i] = children[i].Id
			}

			rb, _ := service.rencanaKinerjaRepository.FindByPokinIds(ctx, tx, pokinIds)
			for _, child := range rb {
				level5map[rk.Id] = append(level5map[rk.Id], child.Id)
			}

		case 4:
			// rekin level 4 -> anak level 5
			children5, _ := service.pohonKinerjaRepository.FindChildPokins(ctx, tx, int64(rk.IdPohon))
			pokin5Ids := make([]int, len(children5))
			for i := range children5 {
				pokin5Ids[i] = children5[i].Id
			}

			rb5, _ := service.rencanaKinerjaRepository.FindByPokinIds(ctx, tx, pokin5Ids)
			for _, child := range rb5 {
				level4map[rk.Id] = append(level4map[rk.Id], child.Id)
			}
		}
	}

	// ---- ROLL UP ----
	final := make(map[string]int64)

	// level 6 fixed
	for id, val := range rawAnggaran {
		final[id] = val
	}

	// level 5 rollup
	for rekin5, children6 := range level5map {
		var total int64 = 0
		for _, r6 := range children6 {
			total += rawAnggaran[r6]
		}
		final[rekin5] = total
	}

	// level 4 rollup (sum dari total level 5)
	for rekin4, children5 := range level4map {
		var total int64 = 0
		for _, r5 := range children5 {
			total += final[r5]
		}
		final[rekin4] = total
	}

	return final, nil
}

func (service *CascadingOpdServiceImpl) preloadAllHierarchies(
	ctx context.Context,
	tx *sql.Tx,
	rks []domain.DetailRekins,
) (
	map[string]pohonkinerja.KegiatanRekinResponse, // Level 6
	map[string]pohonkinerja.SubKegiatanRekinResponse, // Level 6
	map[string][]pohonkinerja.ProgramRekinResponse, // Level 5
	map[string][]pohonkinerja.BidangUrusanCascadingRekinResponse, // Level 4
	error,
) {

	// Kumpulan kode_sub per level
	level6Subs := []string{}
	level5Subs := []string{}
	level4Subs := []string{}

	// Hindari duplikasi
	seen6 := map[string]struct{}{}
	seen5 := map[string]struct{}{}
	seen4 := map[string]struct{}{}

	// Rekin level 5 → kode_sub anak level 6
	mapLevel5ToSub := map[string][]string{}

	// Rekin level 4 → kode_sub anak level 6
	mapLevel4ToSub := map[string][]string{}

	// ------------- LOOP 1: Kumpulkan semua kode_sub -------------
	for _, rk := range rks {

		// ---------- LEVEL 6 ----------
		if rk.LevelPohon == 6 && rk.KodeSubKegiatan != "" {
			if _, ok := seen6[rk.KodeSubKegiatan]; !ok {
				seen6[rk.KodeSubKegiatan] = struct{}{}
				level6Subs = append(level6Subs, rk.KodeSubKegiatan)
			}
		}

		// ---------- LEVEL 5 ----------
		if rk.LevelPohon == 5 {

			children, err := service.pohonKinerjaRepository.FindChildPokins(ctx, tx, int64(rk.IdPohon))
			if err != nil {
				return nil, nil, nil, nil, err
			}

			ids := make([]int, len(children))
			for i := range children {
				ids[i] = children[i].Id
			}

			rb, err := service.rencanaKinerjaRepository.FindByPokinIds(ctx, tx, ids)
			if err != nil {
				return nil, nil, nil, nil, err
			}

			for _, child := range rb {
				if child.KodeSubKegiatan == "" {
					continue
				}

				mapLevel5ToSub[rk.Id] = append(mapLevel5ToSub[rk.Id], child.KodeSubKegiatan)

				if _, ok := seen5[child.KodeSubKegiatan]; !ok {
					seen5[child.KodeSubKegiatan] = struct{}{}
					level5Subs = append(level5Subs, child.KodeSubKegiatan)
				}
			}
		}

		// ---------- LEVEL 4 ----------
		if rk.LevelPohon == 4 {

			// Ambil anak level 5
			children5, err := service.pohonKinerjaRepository.FindChildPokins(ctx, tx, int64(rk.IdPohon))
			if err != nil {
				return nil, nil, nil, nil, err
			}

			ids5 := make([]int, len(children5))
			for i := range children5 {
				ids5[i] = children5[i].Id
			}

			// Ambil anak level 6
			children6, err := service.pohonKinerjaRepository.FindChildPokinsFromParentIds(ctx, tx, ids5)
			if err != nil {
				return nil, nil, nil, nil, err
			}

			ids6 := make([]int, len(children6))
			for i := range children6 {
				ids6[i] = children6[i].Id
			}

			rb, err := service.rencanaKinerjaRepository.FindByPokinIds(ctx, tx, ids6)
			if err != nil {
				return nil, nil, nil, nil, err
			}

			for _, child := range rb {
				if child.KodeSubKegiatan == "" {
					continue
				}

				mapLevel4ToSub[rk.Id] = append(mapLevel4ToSub[rk.Id], child.KodeSubKegiatan)

				if _, ok := seen4[child.KodeSubKegiatan]; !ok {
					seen4[child.KodeSubKegiatan] = struct{}{}
					level4Subs = append(level4Subs, child.KodeSubKegiatan)
				}
			}
		}
	}

	// -------- BATCH QUERY: Level 6 --------
	kegList, _ := service.kegiatanRepository.FindByKodeSubs(ctx, tx, level6Subs)
	subList, _ := service.subKegiatanRepository.FindByKodeSubs(ctx, tx, level6Subs)

	// -------- BATCH QUERY: Level 5 (program) --------
	progList, _ := service.programRepository.FindByKodeSubKegiatans(ctx, tx, level5Subs)

	// -------- BATCH QUERY: Level 4 (bidang urusan) --------
	bidangUrusanList, _ := service.bidangUrusanRepository.FindByKodeSubKegiatans(ctx, tx, level4Subs)

	// =======================
	// BUILD: Level 6
	// =======================
	kegMap := make(map[string]pohonkinerja.KegiatanRekinResponse)
	{
		seen := make(map[string]struct{})
		for _, k := range kegList {
			if _, ok := seen[k.KodeSubKegiatan]; ok {
				continue
			}
			seen[k.KodeSubKegiatan] = struct{}{}

			kegMap[k.KodeSubKegiatan] = pohonkinerja.KegiatanRekinResponse{
				KodeKegiatan: k.KodeKegiatan,
				NamaKegiatan: k.NamaKegiatan,
			}
		}
	}

	// =======================
	// BUILD: SubKegiatan
	// =======================
	subMap := make(map[string]pohonkinerja.SubKegiatanRekinResponse)
	{
		seen := make(map[string]struct{})
		for _, s := range subList {
			if _, ok := seen[s.KodeSubKegiatan]; ok {
				continue
			}
			seen[s.KodeSubKegiatan] = struct{}{}

			subMap[s.KodeSubKegiatan] = pohonkinerja.SubKegiatanRekinResponse{
				KodeSubkegiatan: s.KodeSubKegiatan,
				NamaSubkegiatan: s.NamaSubKegiatan,
			}
		}
	}

	// =======================
	// BUILD: Program (Level 5) PER Rekin ID
	// =======================
	progMap := make(map[string][]pohonkinerja.ProgramRekinResponse)

	for rekinId, kodeSubs := range mapLevel5ToSub {

		seenProgram := make(map[string]struct{})

		for _, kodeSub := range kodeSubs {
			for _, p := range progList {

				if p.KodeSubKegiatan != kodeSub {
					continue
				}

				// eliminate duplicate kode_program
				if _, exists := seenProgram[p.KodeProgram]; exists {
					continue
				}

				seenProgram[p.KodeProgram] = struct{}{}
				progMap[rekinId] = append(
					progMap[rekinId],
					pohonkinerja.ProgramRekinResponse{
						KodeProgram: p.KodeProgram,
						NamaProgram: p.NamaProgram,
					},
				)
			}
		}
	}

	// =======================
	// BUILD: Bidang Urusan (Level 4) PER Rekin ID
	// =======================
	bidangMap := make(map[string][]pohonkinerja.BidangUrusanCascadingRekinResponse)

	for rekinId, kodeSubs := range mapLevel4ToSub {

		seenBidang := make(map[string]struct{})

		for _, kodeSub := range kodeSubs {
			for _, b := range bidangUrusanList {

				if b.KodeSubKegiatan != kodeSub {
					continue
				}

				// eliminate duplicate bidang
				if _, exists := seenBidang[b.KodeBidangUrusan]; exists {
					continue
				}

				seenBidang[b.KodeBidangUrusan] = struct{}{}
				bidangMap[rekinId] = append(
					bidangMap[rekinId],
					pohonkinerja.BidangUrusanCascadingRekinResponse{
						KodeBidangUrusan: b.KodeBidangUrusan,
						NamaBidangUrusan: b.NamaBidangUrusan,
					},
				)
			}
		}
	}

	// FINAL RETURN (semua unique)
	return kegMap, subMap, progMap, bidangMap, nil
}

func toIndikatorResponse(list []domain.Indikator) []pohonkinerja.IndikatorRekinResponse {
	out := make([]pohonkinerja.IndikatorRekinResponse, 0, len(list))

	for _, ind := range list {

		// convert target
		targets := make([]pohonkinerja.TargetRekinResponse, 0, len(ind.Target))
		for _, t := range ind.Target {
			targets = append(targets, pohonkinerja.TargetRekinResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
				Tahun:           t.Tahun,
			})
		}

		// convert indikator
		out = append(out, pohonkinerja.IndikatorRekinResponse{
			Id:               ind.Id,
			RencanaKinerjaId: ind.RencanaKinerjaId,
			NamaIndikator:    ind.Indikator,
			Target:           targets,
		})
	}

	return out
}
