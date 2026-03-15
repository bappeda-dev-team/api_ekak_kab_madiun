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
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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
	RedisClient              *redis.Client
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
	rencanaAksiRepository repository.RencanaAksiRepository,
	RedisClient *redis.Client) *CascadingOpdServiceImpl {
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
		RedisClient:              RedisClient,
	}
}

type ProgramKegiatanSubkegiatan struct {
	LevelPohon      int
	ParentPohon     int
	KodeProgram     string
	NamaProgram     string
	KodeKegiatan    string
	NamaKegiatan    string
	KodeSubKegiatan string
	NamaSubKegiatan string
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
	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdByTahun(ctx, tx, kodeOpd, tahun, "RPJMD")
	if err != nil {
		log.Printf("Warning: Failed to get tujuan OPD data: %v", err)
		return response, nil
	}

	// sudah lengkap plus bidang urusan
	response.TujuanOpd = pohonkinerja.MapTujuanOpdToResponseCascading(tujuanOpds)

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

	// Proses data pohon kinerja
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)
	rencanaKinerjaMap := make(map[int][]domain.RencanaKinerja)
	pohonIDs := make([]int, 0, len(pokins))
	pohonIDSet := make(map[int]struct{})

	// Kelompokkan data dan ambil data indikator & rencana kinerja
	maxLevel := 0
	for _, p := range pokins {
		if p.LevelPohon > maxLevel {
			maxLevel = p.LevelPohon
		}

		if _, ok := pohonIDSet[p.Id]; !ok {
			pohonIDSet[p.Id] = struct{}{}
			pohonIDs = append(pohonIDs, p.Id)
		}

		if pohonMap[p.LevelPohon] == nil {
			pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
		}

		p.NamaOpd = opd.NamaOpd
		pohonMap[p.LevelPohon][p.Parent] = append(
			pohonMap[p.LevelPohon][p.Parent],
			p,
		)
	}

	indikatorPohonMap, err := service.cascadingOpdRepository.FindIndikatorTargetByPokinIds(ctx, tx, pohonIDs)
	if err != nil {
		log.Printf("Error Indikator list")
		return response, nil
	}

	// to dto indikator pokin
	for pohonId, indikatorList := range indikatorPohonMap {
		var indikatorResponses []pohonkinerja.IndikatorResponse

		for _, ind := range indikatorList {
			var targetResponses []pohonkinerja.TargetResponse

			for _, tgt := range ind.Target {
				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
					Id:              tgt.Id,
					IndikatorId:     tgt.IndikatorId,
					TargetIndikator: tgt.Target,
					SatuanIndikator: tgt.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				NamaIndikator: ind.Indikator,
				Target:        targetResponses,
			})
		}

		indikatorMap[pohonId] = indikatorResponses
	}

	pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokinBatchForCascading(ctx, tx, pohonIDs)
	if err != nil {
		log.Printf("Error pelaksana list")
		return response, nil
	}

	pelaksanaByPohon := make(map[int]map[string]domain.PelaksanaPokin)
	for _, pel := range pelaksanaList {
		pohonId, err := strconv.Atoi(pel.PohonKinerjaId)
		if err != nil {
			log.Printf("Error convert pokin id from pelaksana")
			return response, nil
		}
		if pelaksanaByPohon[pohonId] == nil {
			pelaksanaByPohon[pohonId] = make(map[string]domain.PelaksanaPokin)
		}

		pelaksanaByPohon[pohonId][pel.Nip] = domain.PelaksanaPokin{
			Id:             pel.Id,
			PohonKinerjaId: pel.PohonKinerjaId,
			PegawaiId:      pel.PegawaiId,
			NamaPegawai:    pel.NamaPegawai,
			Nip:            pel.Nip,
		}
	}

	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindByPokinIds(ctx, tx, pohonIDs)
	if err != nil {
		log.Printf("Error rekin list: %v", err)
		return response, nil
	}
	rekinIds := make([]string, 0, len(rencanaKinerjaList))
	rekinByPohon := make(map[int][]domain.RencanaKinerja)
	for _, rk := range rencanaKinerjaList {
		rekinIds = append(rekinIds, rk.Id)
		rekinByPohon[rk.IdPohon] = append(
			rekinByPohon[rk.IdPohon],
			rk,
		)
	}
	indikatorByIdRekin, err := service.rencanaKinerjaRepository.IndikatorTargetSasaranByRekinIds(ctx, tx, rekinIds)
	if err != nil {
		log.Printf("Error indikator rekin list: %v", err)
		return response, nil
	}
	// id rekin -> pagu
	rincianBelanjaByIdRekin, err := service.rincianBelanjaRepository.TotalAnggaranByIdRekins(ctx, tx, rekinIds)
	if err != nil {
		log.Printf("Error rincian belanja rekin list: %v", err)
		return response, nil
	}

	// proses rencanaKinerjaMap
	for _, p := range pokins {
		var validRencanaKinerja []domain.RencanaKinerja

		pohonId := p.Id
		pelaksanaMap := pelaksanaByPohon[pohonId]
		rekinList := rekinByPohon[pohonId]

		rekinDedup := make(map[string]bool)
		pegawaiWithRekin := make(map[string]bool)

		// 1 cocokkan rencana kinerja dengan pelaksana
		for _, rk := range rekinList {
			if pel, ok := pelaksanaMap[rk.PegawaiId]; ok {
				if !rekinDedup[rk.Id] {
					rk.NamaPegawai = pel.NamaPegawai
					rk.Indikator = indikatorByIdRekin[rk.Id]
					rk.TotalAnggaran = rincianBelanjaByIdRekin[rk.Id]

					validRencanaKinerja = append(validRencanaKinerja, rk)

					rekinDedup[rk.Id] = true
					pegawaiWithRekin[rk.PegawaiId] = true
				}
			}
		}

		// 2 tambahkan pelaksana tanpa rencana kinerja
		for nip, pel := range pelaksanaMap {
			if !pegawaiWithRekin[nip] {
				validRencanaKinerja = append(validRencanaKinerja, domain.RencanaKinerja{
					IdPohon:     pohonId,
					NamaPohon:   p.NamaPohon,
					Tahun:       p.Tahun,
					PegawaiId:   nip,
					NamaPegawai: pel.NamaPegawai,
					Indikator:   nil,
				})
			}
		}

		rencanaKinerjaMap[pohonId] = validRencanaKinerja
	}
	// Map PKS from rekin
	// id pohon tactical (parent dari operational),
	// kodeSubKegiatan uniq, kode program, nama program
	pksMap := make(map[int]map[string]ProgramKegiatanSubkegiatan)
	// untuk indikator
	kodeRelasiSet := make(map[string]struct{})
	for _, rekins := range rekinByPohon {
		for _, rk := range rekins {
			pohonId := rk.ParentPohon
			// collect subkegiatan by rekin, key by pokin id
			if rk.KodeSubKegiatan != "" && rk.NamaSubKegiatan != "" {
				kodeSubKegiatan := rk.KodeSubKegiatan
				segments := strings.Split(kodeSubKegiatan, ".")
				if len(segments) >= 5 {
					if pksMap[pohonId] == nil {
						pksMap[pohonId] = make(map[string]ProgramKegiatanSubkegiatan)
					}

					kodeProgram := strings.Join(segments[:3], ".")
					kodeKegiatan := strings.Join(segments[:5], ".")
					pksMap[pohonId][kodeSubKegiatan] = ProgramKegiatanSubkegiatan{
						LevelPohon:      rk.LevelPohon,
						ParentPohon:     rk.ParentPohon,
						KodeProgram:     kodeProgram,
						NamaProgram:     rk.Program,
						KodeKegiatan:    kodeKegiatan,
						NamaKegiatan:    rk.NamaKegiatan,
						KodeSubKegiatan: kodeSubKegiatan,
						NamaSubKegiatan: rk.NamaSubKegiatan,
					}
					kodeRelasiSet[kodeProgram] = struct{}{}
					kodeRelasiSet[kodeKegiatan] = struct{}{}
					kodeRelasiSet[kodeSubKegiatan] = struct{}{}
				}
			}
		}
	}
	// flatten untuk batch query
	var kodeRelasiList []string
	for kode := range kodeRelasiSet {
		kodeRelasiList = append(kodeRelasiList, kode)
	}
	// cari indikator program, kegiatan, subekgatain by id pokin
	// di satu list
	indikatorRelasiMap, err := service.programRepository.FindIndikatorTargetByKodeRelasiBatch(ctx, tx, kodeRelasiList)
	if err != nil {
		return response, nil
	}

	start := time.Now()
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
				startBuildStrategic := time.Now()
				strategicResp := service.buildStrategicCascadingResponse(pohonMap, strategic, indikatorMap, rencanaKinerjaMap, pksMap, indikatorRelasiMap)
				response.Strategics = append(response.Strategics, strategicResp)
				log.Printf("buildStrategic %d took %v", strategic.Id, time.Since(startBuildStrategic))
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
	log.Printf("total buildStrategic took %v", time.Since(start))

	return response, nil
}

// VERSI OPTIMASI - Strategic
func (service *CascadingOpdServiceImpl) buildStrategicCascadingResponse(
	pohonMap map[int]map[int][]domain.PohonKinerja,
	strategic domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
	pksMap map[int]map[string]ProgramKegiatanSubkegiatan,
	indikatorRelasiMap map[string][]domain.Indikator,
) pohonkinerja.StrategicCascadingOpdResponse {

	// Proses rencana kinerja untuk strategic
	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
	// Proses pelaksana dari pohon kinerja
	if rencanaKinerjaList, ok := rencanaKinerjaMap[strategic.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			var indikatorResponses []pohonkinerja.IndikatorResponse

			for _, ind := range rk.Indikator {
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range ind.Target {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorResponses = append(indikatorResponses,
					pohonkinerja.IndikatorResponse{
						Id:            ind.Id,
						IdRekin:       rk.Id,
						NamaIndikator: ind.Indikator,
						Target:        targetResponses,
					},
				)
			}

			// Hanya ambil indikator jika rencana kinerja memiliki ID
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
		RencanaKinerja: rencanaKinerjaResponses,
		Indikator:      indikatorMap[strategic.Id],
	}

	// Build tactical responses dan hitung total pagu anggaran
	var totalPaguAnggaran int64 = 0
	var tacticals []pohonkinerja.TacticalCascadingOpdResponse
	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
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
			// FLAG 6
			tacticalResp := service.buildTacticalCascadingResponse(pohonMap, tactical, indikatorMap, rencanaKinerjaMap, pksMap, indikatorRelasiMap)
			tacticals = append(tacticals, tacticalResp)
			// Tambahkan pagu anggaran dari setiap tactical
			totalPaguAnggaran += tacticalResp.PaguAnggaran
		}
		strategicResp.Tacticals = tacticals
	}

	// Set pagu anggaran dari total pagu anggaran tactical
	strategicResp.PaguAnggaran = totalPaguAnggaran

	// set program strategic
	// Convert program map ke slice dan sort
	programMap := make(map[string]pohonkinerja.ProgramResponse)

	for _, tact := range tacticals {
		for _, prg := range pksMap[tact.Id] {

			// skip kalau program sudah pernah diproses
			if _, exists := programMap[prg.KodeProgram]; exists {
				continue
			}

			indikatorProgram := make([]pohonkinerja.IndikatorResponse, 0)

			for _, ind := range indikatorRelasiMap[prg.KodeProgram] {
				targetResp := make([]pohonkinerja.TargetResponse, 0, len(ind.Target))

				for _, tar := range ind.Target {
					targetResp = append(targetResp, pohonkinerja.TargetResponse{
						Id:              tar.Id,
						IndikatorId:     ind.Id,
						TargetIndikator: tar.Target,
						SatuanIndikator: tar.Satuan,
					})
				}

				indikatorProgram = append(indikatorProgram, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					NamaIndikator: ind.Indikator,
					Target:        targetResp,
				})
			}

			programMap[prg.KodeProgram] = pohonkinerja.ProgramResponse{
				KodeProgram: prg.KodeProgram,
				NamaProgram: prg.NamaProgram,
				Indikator:   indikatorProgram,
			}
		}
	}
	programList := make([]pohonkinerja.ProgramResponse, 0, len(programMap))
	for _, prg := range programMap {
		programList = append(programList, prg)
	}

	strategicResp.Program = programList

	return strategicResp
}

func (service *CascadingOpdServiceImpl) buildTacticalCascadingResponse(
	pohonMap map[int]map[int][]domain.PohonKinerja,
	tactical domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
	pksMap map[int]map[string]ProgramKegiatanSubkegiatan,
	indikatorRelasiMap map[string][]domain.Indikator,
) pohonkinerja.TacticalCascadingOpdResponse {

	log.Printf("Building tactical response for ID: %d, Level: %d", tactical.Id, tactical.LevelPohon)

	// Map untuk menyimpan program unik
	// programMap := make(map[string]string)

	// Proses rencana kinerja untuk tactical
	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
	if rencanaKinerjaList, ok := rencanaKinerjaMap[tactical.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			var indikatorResponses []pohonkinerja.IndikatorResponse

			for _, ind := range rk.Indikator {
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range ind.Target {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorResponses = append(indikatorResponses,
					pohonkinerja.IndikatorResponse{
						Id:            ind.Id,
						IdRekin:       rk.Id,
						NamaIndikator: ind.Indikator,
						Target:        targetResponses,
					},
				)
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
			operationalResp := service.buildOperationalCascadingResponse(pohonMap, operational, indikatorMap, rencanaKinerjaMap, indikatorRelasiMap)
			operationals = append(operationals, operationalResp)
			// Tambahkan total anggaran dari setiap operational
			totalPaguAnggaran += operationalResp.TotalAnggaran
		}
		tacticalResp.Operationals = operationals
	}

	// Set pagu anggaran dari total anggaran operational
	tacticalResp.PaguAnggaran = totalPaguAnggaran

	// set program tactical
	// Convert program map ke slice dan sort
	programMap := make(map[string]pohonkinerja.ProgramResponse)

	for _, prg := range pksMap[tacticalResp.Id] {
		// skip kalau program sudah pernah diproses
		if _, exists := programMap[prg.KodeProgram]; exists {
			continue
		}

		indikatorProgram := make([]pohonkinerja.IndikatorResponse, 0)

		for _, ind := range indikatorRelasiMap[prg.KodeProgram] {
			targetResp := make([]pohonkinerja.TargetResponse, 0, len(ind.Target))

			for _, tar := range ind.Target {
				targetResp = append(targetResp, pohonkinerja.TargetResponse{
					Id:              tar.Id,
					IndikatorId:     ind.Id,
					TargetIndikator: tar.Target,
					SatuanIndikator: tar.Satuan,
				})
			}

			indikatorProgram = append(indikatorProgram, pohonkinerja.IndikatorResponse{
				Id:            ind.Id,
				NamaIndikator: ind.Indikator,
				Target:        targetResp,
			})
		}

		programMap[prg.KodeProgram] = pohonkinerja.ProgramResponse{
			KodeProgram: prg.KodeProgram,
			NamaProgram: prg.NamaProgram,
			Indikator:   indikatorProgram,
		}
	}
	programList := make([]pohonkinerja.ProgramResponse, 0, len(programMap))
	for _, prg := range programMap {
		programList = append(programList, prg)
	}

	tacticalResp.Program = programList

	return tacticalResp
}

func (service *CascadingOpdServiceImpl) buildOperationalCascadingResponse(
	pohonMap map[int]map[int][]domain.PohonKinerja,
	operational domain.PohonKinerja,
	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
	indikatorRelasiMap map[string][]domain.Indikator,
) pohonkinerja.OperationalCascadingOpdResponse {

	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalResponse
	var totalAnggaranOperational int64 = 0
	if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
		for _, rk := range rencanaKinerjaList {
			totalAnggaranOperational += rk.TotalAnggaran

			// Indikator rencana kinerja
			var indikatorRekinResponses []pohonkinerja.IndikatorResponse
			for _, ind := range rk.Indikator {
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range ind.Target {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorRekinResponses = append(indikatorRekinResponses,
					pohonkinerja.IndikatorResponse{
						Id:            ind.Id,
						IdRekin:       rk.Id,
						NamaIndikator: ind.Indikator,
						Target:        targetResponses,
					},
				)
			}

			// Indikator subkegiatan
			var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse
			for _, ind := range indikatorRelasiMap[rk.KodeSubKegiatan] {
				targetResp := make([]pohonkinerja.TargetResponse, 0, len(ind.Target))

				for _, tar := range ind.Target {
					targetResp = append(targetResp, pohonkinerja.TargetResponse{
						Id:              tar.Id,
						IndikatorId:     ind.Id,
						TargetIndikator: tar.Target,
						SatuanIndikator: tar.Satuan,
					})
				}

				indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					NamaIndikator: ind.Indikator,
					Target:        targetResp,
				})
			}

			// Indikator kegiatan
			var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse
			for _, ind := range indikatorRelasiMap[rk.KodeKegiatan] {
				targetResp := make([]pohonkinerja.TargetResponse, 0, len(ind.Target))

				for _, tar := range ind.Target {
					targetResp = append(targetResp, pohonkinerja.TargetResponse{
						Id:              tar.Id,
						IndikatorId:     ind.Id,
						TargetIndikator: tar.Target,
						SatuanIndikator: tar.Satuan,
					})
				}

				indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
					Id:            ind.Id,
					NamaIndikator: ind.Indikator,
					Target:        targetResp,
				})
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
				Anggaran:             rk.TotalAnggaran,
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
			childResp := service.buildOperationalNCascadingResponse(pohonMap, opN, indikatorMap, rencanaKinerjaMap)
			childs = append(childs, childResp)
		}
		operationalResp.Childs = childs
	}

	return operationalResp
}

func (service *CascadingOpdServiceImpl) buildOperationalNCascadingResponse(
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

			for _, ind := range rk.Indikator {
				var targetResponses []pohonkinerja.TargetResponse
				for _, target := range ind.Target {
					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}
				indikatorResponses = append(indikatorResponses,
					pohonkinerja.IndikatorResponse{
						Id:            ind.Id,
						IdRekin:       rk.Id,
						NamaIndikator: ind.Indikator,
						Target:        targetResponses,
					},
				)
			}

			// Ambil indikator jika rencana kinerja memiliki ID

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

// func (service *CascadingOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.CascadingOpdResponse, error) {
// 	// Generate cache key
// 	cacheKey := helper.GenerateCacheKey(helper.CacheKeyCascadingOpdAll, kodeOpd, tahun)

// 	// Coba ambil dari cache terlebih dahulu
// 	var response pohonkinerja.CascadingOpdResponse
// 	err := helper.GetFromCache(ctx, service.RedisClient, cacheKey, &response)
// 	if err == nil {
// 		// Cache hit, return data dari cache
// 		log.Printf("Cache HIT untuk key: %s", cacheKey)
// 		return response, nil
// 	}

// 	// Cache miss, ambil data dari database
// 	log.Printf("Cache MISS untuk key: %s, mengambil dari database", cacheKey)

// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		log.Printf("Error starting transaction: %v", err)
// 		return pohonkinerja.CascadingOpdResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	// Validasi OPD
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
// 	if err != nil {
// 		log.Printf("Error: OPD not found for kode_opd=%s: %v", kodeOpd, err)
// 		return pohonkinerja.CascadingOpdResponse{}, errors.New("kode opd tidak ditemukan")
// 	}

// 	// Inisialisasi response dasar
// 	response = pohonkinerja.CascadingOpdResponse{
// 		KodeOpd:    kodeOpd,
// 		NamaOpd:    opd.NamaOpd,
// 		Tahun:      tahun,
// 		TujuanOpd:  make([]pohonkinerja.TujuanOpdCascadingResponse, 0),
// 		Strategics: make([]pohonkinerja.StrategicCascadingOpdResponse, 0),
// 	}

// 	// Ambil data tujuan OPD
// 	tujuanOpds, err := service.tujuanOpdRepository.FindTujuanOpdForCascadingOpd(ctx, tx, kodeOpd, tahun, "RPJMD")
// 	if err != nil {
// 		log.Printf("Warning: Failed to get tujuan OPD data: %v", err)
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.CascadingOpdCacheTTL)
// 		return response, nil
// 	}

// 	log.Printf("Processing %d tujuan OPD entries", len(tujuanOpds))

// 	// OPTIMASI: Batch fetch indikator untuk semua tujuan OPD
// 	var tujuanIds []int
// 	kodeBidangUrusanSet := make(map[string]bool)
// 	for _, tujuan := range tujuanOpds {
// 		tujuanIds = append(tujuanIds, tujuan.Id)
// 		if tujuan.KodeBidangUrusan != "" {
// 			kodeBidangUrusanSet[tujuan.KodeBidangUrusan] = true
// 		}
// 	}

// 	// Batch fetch indikator tujuan OPD (masih pakai loop karena belum ada batch method)
// 	indikatorTujuanMap := make(map[int][]domain.Indikator)
// 	for _, tujuanId := range tujuanIds {
// 		indikators, err := service.tujuanOpdRepository.FindIndikatorByTujuanOpdId(ctx, tx, tujuanId)
// 		if err == nil {
// 			indikatorTujuanMap[tujuanId] = indikators
// 		}
// 	}

// 	// OPTIMASI: Batch fetch bidang urusan
// 	var kodeBidangUrusanList []string
// 	for kode := range kodeBidangUrusanSet {
// 		kodeBidangUrusanList = append(kodeBidangUrusanList, kode)
// 	}

// 	bidangUrusanMap := make(map[string]domainmaster.BidangUrusan)
// 	for _, kode := range kodeBidangUrusanList {
// 		bidangUrusan, err := service.bidangUrusanRepository.FindByKodeBidangUrusan(ctx, tx, kode)
// 		if err == nil {
// 			bidangUrusanMap[kode] = bidangUrusan
// 		}
// 	}

// 	// Build tujuan opd response
// 	for _, tujuan := range tujuanOpds {
// 		indikators := indikatorTujuanMap[tujuan.Id]
// 		var indikatorResponses []pohonkinerja.IndikatorTujuanResponse
// 		for _, indikator := range indikators {
// 			indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorTujuanResponse{
// 				Indikator: indikator.Indikator,
// 			})
// 		}

// 		bidangUrusan := bidangUrusanMap[tujuan.KodeBidangUrusan]

// 		response.TujuanOpd = append(response.TujuanOpd, pohonkinerja.TujuanOpdCascadingResponse{
// 			Id:         tujuan.Id,
// 			KodeOpd:    tujuan.KodeOpd,
// 			Tujuan:     tujuan.Tujuan,
// 			KodeBidang: tujuan.KodeBidangUrusan,
// 			NamaBidang: bidangUrusan.NamaBidangUrusan,
// 			Indikator:  indikatorResponses,
// 		})
// 	}

// 	// Ambil data pohon kinerja
// 	pokins, err := service.cascadingOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
// 	if err != nil {
// 		log.Printf("Error getting pohon kinerja data: %v", err)
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.CascadingOpdCacheTTL)
// 		return response, nil
// 	}

// 	if len(pokins) == 0 {
// 		log.Printf("No pohon kinerja found for kodeOpd=%s, tahun=%s", kodeOpd, tahun)
// 		helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.CascadingOpdCacheTTL)
// 		return response, nil
// 	}

// 	log.Printf("Processing %d pohon kinerja entries", len(pokins))

// 	// OPTIMASI: Proses data pohon kinerja dengan batch queries
// 	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
// 	indikatorMap := make(map[int][]pohonkinerja.IndikatorResponse)
// 	rencanaKinerjaMap := make(map[int][]domain.RencanaKinerja)

// 	// Kumpulkan semua pokin IDs untuk batch queries
// 	var pokinIds []int
// 	maxLevel := 0
// 	for _, p := range pokins {
// 		if p.LevelPohon > maxLevel {
// 			maxLevel = p.LevelPohon
// 		}

// 		if pohonMap[p.LevelPohon] == nil {
// 			pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
// 		}

// 		p.NamaOpd = opd.NamaOpd
// 		pohonMap[p.LevelPohon][p.Parent] = append(
// 			pohonMap[p.LevelPohon][p.Parent],
// 			p,
// 		)
// 		pokinIds = append(pokinIds, p.Id)
// 	}

// 	// OPTIMASI: Batch fetch semua rencana kinerja untuk semua pokin
// 	// Perlu membuat method FindByPokinIdsBatch di repository
// 	rencanaKinerjaBatch := make(map[int][]domain.RencanaKinerja)
// 	for _, pokinId := range pokinIds {
// 		rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokinId)
// 		if err == nil {
// 			rencanaKinerjaBatch[pokinId] = rekinList
// 		}
// 	}

// 	// OPTIMASI: Batch fetch semua pelaksana
// 	pelaksanaBatch, err := service.pohonKinerjaRepository.FindPelaksanaPokinBatch(ctx, tx, pokinIds)
// 	if err != nil {
// 		log.Printf("Error batch fetching pelaksana: %v", err)
// 		pelaksanaBatch = make(map[int][]domain.PelaksanaPokin)
// 	}

// 	// OPTIMASI: Kumpulkan semua pegawai IDs untuk batch fetch
// 	// Perlu mapping: pelaksana.PegawaiId adalah id pegawai (bukan NIP)
// 	pegawaiIdsSet := make(map[string]bool)
// 	for _, pelaksanaList := range pelaksanaBatch {
// 		for _, pelaksana := range pelaksanaList {
// 			pegawaiIdsSet[pelaksana.PegawaiId] = true
// 		}
// 	}

// 	var pegawaiIds []string
// 	for id := range pegawaiIdsSet {
// 		pegawaiIds = append(pegawaiIds, id)
// 	}

// 	// Batch fetch semua pegawai (perlu method FindByIdsBatch)
// 	pegawaiMap := make(map[string]*domainmaster.Pegawai)
// 	for _, pegawaiId := range pegawaiIds {
// 		pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pegawaiId)
// 		if err == nil {
// 			pegawaiMap[pegawai.Id] = &pegawai
// 			pegawaiMap[pegawai.Nip] = &pegawai // Map by NIP juga untuk lookup rencana kinerja
// 		}
// 	}

// 	// Process rencana kinerja untuk setiap pokin menggunakan batch data
// 	for _, p := range pokins {
// 		rencanaKinerjaList := rencanaKinerjaBatch[p.Id]
// 		if len(rencanaKinerjaList) == 0 {
// 			continue
// 		}

// 		pelaksanaList := pelaksanaBatch[p.Id]
// 		pelaksanaMap := make(map[string]*domainmaster.Pegawai)
// 		rekinMap := make(map[string]bool)

// 		// Build pelaksana map dari batch pegawai
// 		// Perhatikan: pelaksana.PegawaiId adalah id pegawai, kita perlu map ke NIP
// 		for _, pelaksana := range pelaksanaList {
// 			if pegawai, exists := pegawaiMap[pelaksana.PegawaiId]; exists {
// 				pelaksanaMap[pegawai.Nip] = pegawai // Map by NIP untuk matching dengan rencana kinerja
// 			}
// 		}

// 		// Filter rencana kinerja yang sesuai dengan pelaksana
// 		var validRencanaKinerja []domain.RencanaKinerja
// 		for _, rk := range rencanaKinerjaList {
// 			// rk.PegawaiId adalah NIP, kita cek apakah NIP ini ada di pelaksanaMap
// 			if pegawai, exists := pelaksanaMap[rk.PegawaiId]; exists {
// 				if !rekinMap[rk.Id] {
// 					rk.NamaPegawai = pegawai.NamaPegawai
// 					validRencanaKinerja = append(validRencanaKinerja, rk)
// 					rekinMap[rk.Id] = true
// 				}
// 			}
// 		}

// 		// Track pegawai yang sudah memiliki rencana kinerja
// 		pegawaiWithRekinMap := make(map[string]bool)
// 		for _, rk := range validRencanaKinerja {
// 			pegawaiWithRekinMap[rk.PegawaiId] = true
// 		}

// 		// Tambahkan pelaksana yang belum memiliki rencana kinerja
// 		for _, pegawai := range pelaksanaMap {
// 			if !pegawaiWithRekinMap[pegawai.Nip] {
// 				validRencanaKinerja = append(validRencanaKinerja, domain.RencanaKinerja{
// 					IdPohon:     p.Id,
// 					NamaPohon:   p.NamaPohon,
// 					Tahun:       p.Tahun,
// 					PegawaiId:   pegawai.Nip,
// 					NamaPegawai: pegawai.NamaPegawai,
// 					Indikator:   nil,
// 				})
// 			}
// 		}

// 		if len(validRencanaKinerja) > 0 {
// 			rencanaKinerjaMap[p.Id] = validRencanaKinerja
// 		}
// 	}

// 	// OPTIMASI: Batch fetch indikator untuk pohon kinerja
// 	indikatorBatch, err := service.pohonKinerjaRepository.FindIndikatorByPokinIdsBatch(ctx, tx, pokinIds)
// 	if err == nil {
// 		// Kumpulkan semua indikator IDs untuk batch fetch target
// 		var allIndikatorIds []string
// 		for _, indikatorList := range indikatorBatch {
// 			for _, indikator := range indikatorList {
// 				allIndikatorIds = append(allIndikatorIds, indikator.Id)
// 			}
// 		}

// 		// Batch fetch semua target
// 		var targetBatch map[string][]domain.Target
// 		if len(allIndikatorIds) > 0 {
// 			targetBatch, err = service.pohonKinerjaRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorIds)
// 			if err != nil {
// 				log.Printf("Error batch fetching targets: %v", err)
// 				targetBatch = make(map[string][]domain.Target)
// 			}
// 		} else {
// 			targetBatch = make(map[string][]domain.Target)
// 		}

// 		// Build indikator map
// 		for pokinId, indikatorList := range indikatorBatch {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse
// 			for _, indikator := range indikatorList {
// 				targetList := targetBatch[indikator.Id]
// 				var targetResponses []pohonkinerja.TargetResponse
// 				for _, target := range targetList {
// 					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 						Id:              target.Id,
// 						IndikatorId:     target.IndikatorId,
// 						TargetIndikator: target.Target,
// 						SatuanIndikator: target.Satuan,
// 					})
// 				}

// 				indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 					Id:            indikator.Id,
// 					IdPokin:       fmt.Sprint(pokinId),
// 					NamaIndikator: indikator.Indikator,
// 					Target:        targetResponses,
// 				})
// 			}
// 			if len(indikatorResponses) > 0 {
// 				indikatorMap[pokinId] = indikatorResponses
// 			}
// 		}
// 	}

// 	// Build response untuk strategic (level 4)
// 	if strategicList := pohonMap[4]; len(strategicList) > 0 {
// 		log.Printf("Processing %d strategic entries", len(strategicList))
// 		for _, strategicsByParent := range strategicList {
// 			sort.Slice(strategicsByParent, func(i, j int) bool {
// 				// Prioritaskan status "pokin dari pemda"
// 				if strategicsByParent[i].Status == "pokin dari pemda" && strategicsByParent[j].Status != "pokin dari pemda" {
// 					return true
// 				}
// 				if strategicsByParent[i].Status != "pokin dari pemda" && strategicsByParent[j].Status == "pokin dari pemda" {
// 					return false
// 				}
// 				return strategicsByParent[i].Id < strategicsByParent[j].Id
// 			})

// 			for _, strategic := range strategicsByParent {
// 				strategicResp := service.buildStrategicCascadingResponseOptimized(ctx, tx, pohonMap, strategic, indikatorMap, rencanaKinerjaMap)
// 				response.Strategics = append(response.Strategics, strategicResp)
// 			}
// 		}

// 		sort.Slice(response.Strategics, func(i, j int) bool {
// 			if response.Strategics[i].Status == "pokin dari pemda" && response.Strategics[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if response.Strategics[i].Status != "pokin dari pemda" && response.Strategics[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return response.Strategics[i].Id < response.Strategics[j].Id
// 		})
// 	}

// 	// Simpan ke cache setelah berhasil mengambil data
// 	helper.SetToCache(ctx, service.RedisClient, cacheKey, response, helper.CascadingOpdCacheTTL)
// 	log.Printf("Data disimpan ke cache dengan key: %s", cacheKey)

// 	return response, nil
// }

// // VERSI OPTIMASI
// func (service *CascadingOpdServiceImpl) buildStrategicCascadingResponseOptimized(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	strategic domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.StrategicCascadingOpdResponse {

// 	// Proses rencana kinerja untuk strategic dengan batch queries
// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[strategic.Id]; ok {
// 		// OPTIMASI: Kumpulkan semua rekin IDs untuk batch fetch indikator
// 		var rekinIds []string
// 		for _, rk := range rencanaKinerjaList {
// 			if rk.Id != "" {
// 				rekinIds = append(rekinIds, rk.Id)
// 			}
// 		}

// 		// Batch fetch indikator rekin (perlu method FindIndikatorbyRekinIdsBatch)
// 		// Untuk sekarang masih pakai loop, tapi bisa dioptimasi lebih lanjut
// 		indikatorRekinMap := make(map[string][]domain.Indikator)
// 		for _, rekinId := range rekinIds {
// 			indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rekinId)
// 			if err == nil {
// 				indikatorRekinMap[rekinId] = indikators
// 			}
// 		}

// 		// Kumpulkan semua indikator IDs untuk batch fetch target
// 		var allIndikatorRekinIds []string
// 		for _, indikators := range indikatorRekinMap {
// 			for _, ind := range indikators {
// 				allIndikatorRekinIds = append(allIndikatorRekinIds, ind.Id)
// 			}
// 		}

// 		// Batch fetch semua target
// 		targetBatch := make(map[string][]domain.Target)
// 		if len(allIndikatorRekinIds) > 0 {
// 			var err error
// 			targetBatch, err = service.pohonKinerjaRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorRekinIds)
// 			if err != nil {
// 				log.Printf("Error batch fetching targets: %v", err)
// 			}
// 		}

// 		// Build rencana kinerja responses
// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse
// 			if rk.Id != "" {
// 				indikators := indikatorRekinMap[rk.Id]
// 				for _, ind := range indikators {
// 					targetList := targetBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						IdRekin:       rk.Id,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            strategic.Id,
// 				NamaPohon:          strategic.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              strategic.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	// OPTIMASI: Kumpulkan semua kode program untuk batch fetch
// 	programMap := make(map[string]string)
// 	if tacticalList, exists := pohonMap[5][strategic.Id]; exists {
// 		for _, tactical := range tacticalList {
// 			if operationalList, exists := pohonMap[6][tactical.Id]; exists {
// 				for _, operational := range operationalList {
// 					if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 						for _, rk := range rencanaKinerjaList {
// 							if rk.KodeSubKegiatan != "" {
// 								segments := strings.Split(rk.KodeSubKegiatan, ".")
// 								if len(segments) >= 3 {
// 									kodeProgram := strings.Join(segments[:3], ".")
// 									if _, exists := programMap[kodeProgram]; !exists {
// 										program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
// 										if err == nil {
// 											programMap[kodeProgram] = program.NamaProgram
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// Batch fetch indikator program dan target (perlu optimasi lebih lanjut)
// 	var programList []pohonkinerja.ProgramResponse
// 	for kodeProgram, namaProgram := range programMap {
// 		var indikatorProgramResponses []pohonkinerja.IndikatorResponse
// 		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 			ctx, tx, kodeProgram, strategic.KodeOpd, strategic.Tahun,
// 		)
// 		if err == nil {
// 			// Kumpulkan indikator IDs untuk batch fetch target
// 			var indikatorIds []string
// 			for _, ind := range indikatorProgram {
// 				indikatorIds = append(indikatorIds, ind.Id)
// 			}

// 			// Batch fetch target
// 			targetProgramBatch := make(map[string][]domain.Target)
// 			if len(indikatorIds) > 0 {
// 				// Gunakan FindTargetByIndikatorIdsBatch jika ada di cascadingOpdRepository
// 				// Untuk sekarang masih pakai loop
// 				for _, indikatorId := range indikatorIds {
// 					targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, indikatorId)
// 					if err == nil {
// 						targetProgramBatch[indikatorId] = targets
// 					}
// 				}
// 			}

// 			for _, ind := range indikatorProgram {
// 				targetList := targetProgramBatch[ind.Id]
// 				var targetResponses []pohonkinerja.TargetResponse
// 				for _, target := range targetList {
// 					targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 						Id:              target.Id,
// 						IndikatorId:     target.IndikatorId,
// 						TargetIndikator: target.Target,
// 						SatuanIndikator: target.Satuan,
// 					})
// 				}

// 				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
// 					Id:            ind.Id,
// 					Kode:          ind.Kode,
// 					NamaIndikator: ind.Indikator,
// 					Target:        targetResponses,
// 				})
// 			}
// 		}

// 		programList = append(programList, pohonkinerja.ProgramResponse{
// 			KodeProgram: kodeProgram,
// 			NamaProgram: namaProgram,
// 			Indikator:   indikatorProgramResponses,
// 		})
// 	}

// 	strategicResp := pohonkinerja.StrategicCascadingOpdResponse{
// 		Id:         strategic.Id,
// 		Parent:     nil,
// 		Strategi:   strategic.NamaPohon,
// 		JenisPohon: strategic.JenisPohon,
// 		LevelPohon: strategic.LevelPohon,
// 		Keterangan: strategic.Keterangan,
// 		Status:     strategic.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: strategic.KodeOpd,
// 			NamaOpd: strategic.NamaOpd,
// 		},
// 		IsActive:       strategic.IsActive,
// 		Program:        programList,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[strategic.Id],
// 	}

// 	// Build tactical responses dan hitung total pagu anggaran
// 	var totalPaguAnggaran int64 = 0
// 	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
// 		var tacticals []pohonkinerja.TacticalCascadingOpdResponse
// 		sort.Slice(tacticalList, func(i, j int) bool {
// 			if tacticalList[i].Status == "pokin dari pemda" && tacticalList[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if tacticalList[i].Status != "pokin dari pemda" && tacticalList[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return tacticalList[i].Id < tacticalList[j].Id
// 		})

// 		for _, tactical := range tacticalList {
// 			tacticalResp := service.buildTacticalCascadingResponseOptimized(ctx, tx, pohonMap, tactical, indikatorMap, rencanaKinerjaMap)
// 			tacticals = append(tacticals, tacticalResp)
// 			totalPaguAnggaran += tacticalResp.PaguAnggaran
// 		}
// 		strategicResp.Tacticals = tacticals
// 	}

// 	strategicResp.PaguAnggaran = totalPaguAnggaran
// 	return strategicResp
// }
// func (service *CascadingOpdServiceImpl) buildTacticalCascadingResponseOptimized(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	tactical domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
// 	indikatorRekinBatch map[string][]domain.Indikator,
// 	targetBatch map[string][]domain.Target) pohonkinerja.TacticalCascadingOpdResponse {

// 	// Proses rencana kinerja untuk tactical dengan batch data
// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[tactical.Id]; ok {
// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse
// 			if rk.Id != "" {
// 				// Gunakan batch data yang sudah di-fetch
// 				indikators := indikatorRekinBatch[rk.Id]
// 				for _, ind := range indikators {
// 					targetList := targetBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						IdRekin:       rk.Id,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            tactical.Id,
// 				NamaPohon:          tactical.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              tactical.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	// OPTIMASI: Kumpulkan semua kode program untuk batch fetch
// 	programMap := make(map[string]string)
// 	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
// 		for _, operational := range operationalList {
// 			if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 				for _, rk := range rencanaKinerjaList {
// 					if rk.KodeSubKegiatan != "" {
// 						segments := strings.Split(rk.KodeSubKegiatan, ".")
// 						if len(segments) >= 3 {
// 							kodeProgram := strings.Join(segments[:3], ".")
// 							if _, exists := programMap[kodeProgram]; !exists {
// 								program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
// 								if err == nil {
// 									programMap[kodeProgram] = program.NamaProgram
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// OPTIMASI: Batch fetch semua indikator program sekaligus
// 	var kodeProgramList []string
// 	for kode := range programMap {
// 		kodeProgramList = append(kodeProgramList, kode)
// 	}

// 	indikatorProgramBatch := make(map[string][]domain.Indikator)
// 	for _, kodeProgram := range kodeProgramList {
// 		indikators, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 			ctx, tx, kodeProgram, tactical.KodeOpd, tactical.Tahun,
// 		)
// 		if err == nil {
// 			indikatorProgramBatch[kodeProgram] = indikators
// 		}
// 	}

// 	// Kumpulkan semua indikator program IDs untuk batch fetch target
// 	var allIndikatorProgramIds []string
// 	for _, indikators := range indikatorProgramBatch {
// 		for _, ind := range indikators {
// 			allIndikatorProgramIds = append(allIndikatorProgramIds, ind.Id)
// 		}
// 	}

// 	// Batch fetch semua target program
// 	targetProgramBatch := make(map[string][]domain.Target)
// 	if len(allIndikatorProgramIds) > 0 {
// 		var err error
// 		targetProgramBatch, err = service.cascadingOpdRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorProgramIds)
// 		if err != nil {
// 			log.Printf("Error batch fetching target program: %v", err)
// 		}
// 	}

// 	// Build program list
// 	var programList []pohonkinerja.ProgramResponse
// 	for kodeProgram, namaProgram := range programMap {
// 		indikators := indikatorProgramBatch[kodeProgram]
// 		var indikatorProgramResponses []pohonkinerja.IndikatorResponse

// 		for _, ind := range indikators {
// 			targetList := targetProgramBatch[ind.Id]
// 			var targetResponses []pohonkinerja.TargetResponse
// 			for _, target := range targetList {
// 				targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 					Id:              target.Id,
// 					IndikatorId:     target.IndikatorId,
// 					TargetIndikator: target.Target,
// 					SatuanIndikator: target.Satuan,
// 				})
// 			}

// 			indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
// 				Id:            ind.Id,
// 				Kode:          ind.Kode,
// 				NamaIndikator: ind.Indikator,
// 				Target:        targetResponses,
// 			})
// 		}

// 		programList = append(programList, pohonkinerja.ProgramResponse{
// 			KodeProgram: kodeProgram,
// 			NamaProgram: namaProgram,
// 			Indikator:   indikatorProgramResponses,
// 		})
// 	}

// 	tacticalResp := pohonkinerja.TacticalCascadingOpdResponse{
// 		Id:         tactical.Id,
// 		Parent:     tactical.Parent,
// 		Strategi:   tactical.NamaPohon,
// 		JenisPohon: tactical.JenisPohon,
// 		LevelPohon: tactical.LevelPohon,
// 		Keterangan: tactical.Keterangan,
// 		Status:     tactical.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: tactical.KodeOpd,
// 			NamaOpd: tactical.NamaOpd,
// 		},
// 		IsActive:       tactical.IsActive,
// 		Program:        programList,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[tactical.Id],
// 	}

// 	// Build operational responses dan hitung total pagu anggaran
// 	var totalPaguAnggaran int64 = 0
// 	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
// 		var operationals []pohonkinerja.OperationalCascadingOpdResponse
// 		sort.Slice(operationalList, func(i, j int) bool {
// 			if operationalList[i].Status == "pokin dari pemda" && operationalList[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if operationalList[i].Status != "pokin dari pemda" && operationalList[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return operationalList[i].Id < operationalList[j].Id
// 		})

// 		for _, operational := range operationalList {
// 			operationalResp := service.buildOperationalCascadingResponseOptimized(ctx, tx, pohonMap, operational, indikatorMap, rencanaKinerjaMap, indikatorRekinBatch, targetBatch)
// 			operationals = append(operationals, operationalResp)
// 			totalPaguAnggaran += operationalResp.TotalAnggaran
// 		}
// 		tacticalResp.Operationals = operationals
// 	}

// 	tacticalResp.PaguAnggaran = totalPaguAnggaran
// 	return tacticalResp
// }

// // Versi optimized untuk operational dengan batch queries
// func (service *CascadingOpdServiceImpl) buildOperationalCascadingResponseOptimized(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	operational domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
// 	indikatorRekinBatch map[string][]domain.Indikator,
// 	targetBatch map[string][]domain.Target) pohonkinerja.OperationalCascadingOpdResponse {

// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalResponse
// 	var totalAnggaranOperational int64 = 0

// 	// OPTIMASI: Batch fetch anggaran untuk semua rencana kinerja
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 		// Kumpulkan semua rekin IDs untuk batch fetch rencana aksi
// 		var rekinIds []string
// 		for _, rk := range rencanaKinerjaList {
// 			if rk.Id != "" {
// 				rekinIds = append(rekinIds, rk.Id)
// 			}
// 		}

// 		// Batch fetch rencana aksi untuk semua rencana kinerja (perlu method FindAllBatch)
// 		// Untuk sekarang masih loop, tapi bisa dioptimasi lebih lanjut dengan batch method
// 		rencanaAksiMap := make(map[string][]domain.RencanaAksi)
// 		for _, rekinId := range rekinIds {
// 			rencanaAksiList, err := service.rencanaAksiRepository.FindAll(ctx, tx, rekinId)
// 			if err == nil {
// 				rencanaAksiMap[rekinId] = rencanaAksiList
// 			}
// 		}

// 		// Kumpulkan semua renaksi IDs untuk batch fetch anggaran
// 		var renaksiIds []string
// 		for _, rencanaAksiList := range rencanaAksiMap {
// 			for _, ra := range rencanaAksiList {
// 				renaksiIds = append(renaksiIds, ra.Id)
// 			}
// 		}

// 		// Batch fetch anggaran (perlu method FindAnggaranByRenaksiIdsBatch)
// 		// Untuk sekarang masih loop
// 		anggaranMap := make(map[string]int64)
// 		for _, renaksiId := range renaksiIds {
// 			rincianBelanja, err := service.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, renaksiId)
// 			if err == nil {
// 				anggaranMap[renaksiId] = rincianBelanja.Anggaran
// 			}
// 		}

// 		// Kumpulkan kode subkegiatan dan kegiatan untuk batch fetch indikator
// 		kodeSubkegiatanSet := make(map[string]bool)
// 		kodeKegiatanSet := make(map[string]bool)
// 		for _, rk := range rencanaKinerjaList {
// 			if rk.KodeSubKegiatan != "" {
// 				kodeSubkegiatanSet[rk.KodeSubKegiatan] = true
// 			}
// 			if rk.KodeKegiatan != "" {
// 				kodeKegiatanSet[rk.KodeKegiatan] = true
// 			}
// 		}

// 		// Batch fetch indikator subkegiatan dan kegiatan
// 		indikatorSubkegiatanBatch := make(map[string][]domain.Indikator)
// 		for kode := range kodeSubkegiatanSet {
// 			indikators, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 				ctx, tx, kode, operational.KodeOpd, operational.Tahun,
// 			)
// 			if err == nil {
// 				indikatorSubkegiatanBatch[kode] = indikators
// 			}
// 		}

// 		indikatorKegiatanBatch := make(map[string][]domain.Indikator)
// 		for kode := range kodeKegiatanSet {
// 			indikators, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 				ctx, tx, kode, operational.KodeOpd, operational.Tahun,
// 			)
// 			if err == nil {
// 				indikatorKegiatanBatch[kode] = indikators
// 			}
// 		}

// 		// Kumpulkan semua indikator IDs untuk batch fetch target
// 		var allIndikatorIds []string
// 		for _, indikators := range indikatorSubkegiatanBatch {
// 			for _, ind := range indikators {
// 				allIndikatorIds = append(allIndikatorIds, ind.Id)
// 			}
// 		}
// 		for _, indikators := range indikatorKegiatanBatch {
// 			for _, ind := range indikators {
// 				allIndikatorIds = append(allIndikatorIds, ind.Id)
// 			}
// 		}

// 		// Batch fetch semua target untuk indikator subkegiatan dan kegiatan
// 		targetSubkegiatanKegiatanBatch := make(map[string][]domain.Target)
// 		if len(allIndikatorIds) > 0 {
// 			var err error
// 			targetSubkegiatanKegiatanBatch, err = service.cascadingOpdRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorIds)
// 			if err != nil {
// 				log.Printf("Error batch fetching target subkegiatan/kegiatan: %v", err)
// 			}
// 		}

// 		// Build rencana kinerja responses
// 		for _, rk := range rencanaKinerjaList {
// 			var totalAnggaranRenkin int64 = 0
// 			if rk.Id != "" {
// 				rencanaAksiList := rencanaAksiMap[rk.Id]
// 				for _, ra := range rencanaAksiList {
// 					totalAnggaranRenkin += anggaranMap[ra.Id]
// 				}
// 			}
// 			totalAnggaranOperational += totalAnggaranRenkin

// 			// Indikator rencana kinerja (gunakan batch data)
// 			var indikatorRekinResponses []pohonkinerja.IndikatorResponse
// 			if rk.Id != "" {
// 				indikators := indikatorRekinBatch[rk.Id]
// 				for _, ind := range indikators {
// 					targetList := targetBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorRekinResponses = append(indikatorRekinResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						IdRekin:       rk.Id,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			// Indikator subkegiatan (gunakan batch data)
// 			var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse
// 			if rk.KodeSubKegiatan != "" {
// 				indikators := indikatorSubkegiatanBatch[rk.KodeSubKegiatan]
// 				for _, ind := range indikators {
// 					targetList := targetSubkegiatanKegiatanBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						Kode:          ind.Kode,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			// Indikator kegiatan (gunakan batch data)
// 			var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse
// 			if rk.KodeKegiatan != "" {
// 				indikators := indikatorKegiatanBatch[rk.KodeKegiatan]
// 				for _, ind := range indikators {
// 					targetList := targetSubkegiatanKegiatanBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						Kode:          ind.Kode,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalResponse{
// 				Id:                   rk.Id,
// 				IdPohon:              operational.Id,
// 				NamaPohon:            operational.NamaPohon,
// 				NamaRencanaKinerja:   rk.NamaRencanaKinerja,
// 				Tahun:                operational.Tahun,
// 				PegawaiId:            rk.PegawaiId,
// 				NamaPegawai:          rk.NamaPegawai,
// 				KodeSubkegiatan:      rk.KodeSubKegiatan,
// 				NamaSubkegiatan:      rk.NamaSubKegiatan,
// 				Anggaran:             totalAnggaranRenkin,
// 				IndikatorSubkegiatan: indikatorSubkegiatanResponses,
// 				KodeKegiatan:         rk.KodeKegiatan,
// 				NamaKegiatan:         rk.NamaKegiatan,
// 				IndikatorKegiatan:    indikatorKegiatanResponses,
// 				Indikator:            indikatorRekinResponses,
// 			})
// 		}
// 	}

// 	operationalResp := pohonkinerja.OperationalCascadingOpdResponse{
// 		Id:         operational.Id,
// 		Parent:     operational.Parent,
// 		Strategi:   operational.NamaPohon,
// 		JenisPohon: operational.JenisPohon,
// 		LevelPohon: operational.LevelPohon,
// 		Keterangan: operational.Keterangan,
// 		Status:     operational.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: operational.KodeOpd,
// 			NamaOpd: operational.NamaOpd,
// 		},
// 		IsActive:       operational.IsActive,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[operational.Id],
// 		TotalAnggaran:  totalAnggaranOperational,
// 	}

// 	// Build operational N responses jika ada
// 	nextLevel := operational.LevelPohon + 1
// 	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
// 		var childs []pohonkinerja.OperationalNOpdCascadingResponse
// 		sort.Slice(operationalNList, func(i, j int) bool {
// 			return operationalNList[i].Id < operationalNList[j].Id
// 		})

// 		for _, opN := range operationalNList {
// 			childResp := service.buildOperationalNCascadingResponseOptimized(ctx, tx, pohonMap, opN, indikatorMap, rencanaKinerjaMap, indikatorRekinBatch, targetBatch)
// 			childs = append(childs, childResp)
// 		}
// 		operationalResp.Childs = childs
// 	}

// 	return operationalResp
// }

// // Versi optimized untuk operational N dengan batch queries
// func (service *CascadingOpdServiceImpl) buildOperationalNCascadingResponseOptimized(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	operationalN domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja,
// 	indikatorRekinBatch map[string][]domain.Indikator,
// 	targetBatch map[string][]domain.Target) pohonkinerja.OperationalNOpdCascadingResponse {

// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalNResponse

// 	// Proses rencana kinerja dengan batch data
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[operationalN.Id]; ok {
// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse
// 			if rk.Id != "" {
// 				// Gunakan batch data yang sudah di-fetch
// 				indikators := indikatorRekinBatch[rk.Id]
// 				for _, ind := range indikators {
// 					targetList := targetBatch[ind.Id]
// 					var targetResponses []pohonkinerja.TargetResponse
// 					for _, target := range targetList {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}

// 					indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 						Id:            ind.Id,
// 						IdRekin:       rk.Id,
// 						NamaIndikator: ind.Indikator,
// 						Target:        targetResponses,
// 					})
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalNResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            operationalN.Id,
// 				NamaPohon:          operationalN.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              operationalN.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	operationalNResp := pohonkinerja.OperationalNOpdCascadingResponse{
// 		Id:         operationalN.Id,
// 		Parent:     operationalN.Parent,
// 		Strategi:   operationalN.NamaPohon,
// 		JenisPohon: operationalN.JenisPohon,
// 		LevelPohon: operationalN.LevelPohon,
// 		Keterangan: operationalN.Keterangan,
// 		Status:     operationalN.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: operationalN.KodeOpd,
// 			NamaOpd: operationalN.NamaOpd,
// 		},
// 		IsActive:       operationalN.IsActive,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[operationalN.Id],
// 	}

// 	// Proses child nodes jika ada
// 	nextLevel := operationalN.LevelPohon + 1
// 	if childList := pohonMap[nextLevel][operationalN.Id]; len(childList) > 0 {
// 		var childs []pohonkinerja.OperationalNOpdCascadingResponse
// 		sort.Slice(childList, func(i, j int) bool {
// 			return childList[i].Id < childList[j].Id
// 		})

// 		for _, child := range childList {
// 			childResp := service.buildOperationalNCascadingResponseOptimized(
// 				ctx,
// 				tx,
// 				pohonMap,
// 				child,
// 				indikatorMap,
// 				rencanaKinerjaMap,
// 				indikatorRekinBatch,
// 				targetBatch,
// 			)
// 			childs = append(childs, childResp)
// 		}
// 		operationalNResp.Childs = childs
// 	}

// 	return operationalNResp
// }

// //AKHIR VERSI OPTIMASI

// func (service *CascadingOpdServiceImpl) buildStrategicCascadingResponse(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	strategic domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.StrategicCascadingOpdResponse {

// 	// Proses rencana kinerja untuk strategic
// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
// 	// Proses pelaksana dari pohon kinerja
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[strategic.Id]; ok {
// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse

// 			// Hanya ambil indikator jika rencana kinerja memiliki ID
// 			if rk.Id != "" {
// 				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
// 				if err == nil {
// 					for _, ind := range indikatorRekin {
// 						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}
// 						indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							IdRekin:       rk.Id,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            strategic.Id,
// 				NamaPohon:          strategic.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              strategic.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	// Proses program dari level 6 melalui level 5
// 	programMap := make(map[string]string)
// 	if tacticalList, exists := pohonMap[5][strategic.Id]; exists {
// 		for _, tactical := range tacticalList {
// 			if operationalList, exists := pohonMap[6][tactical.Id]; exists {
// 				for _, operational := range operationalList {
// 					if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 						for _, rk := range rencanaKinerjaList {
// 							if rk.KodeSubKegiatan != "" {
// 								segments := strings.Split(rk.KodeSubKegiatan, ".")
// 								if len(segments) >= 3 {
// 									kodeProgram := strings.Join(segments[:3], ".")
// 									if _, exists := programMap[kodeProgram]; !exists {
// 										program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
// 										if err == nil {
// 											programMap[kodeProgram] = program.NamaProgram
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// Convert program map ke slice dan sort
// 	var programList []pohonkinerja.ProgramResponse
// 	for kodeProgram, namaProgram := range programMap {
// 		// Ambil indikator program
// 		var indikatorProgramResponses []pohonkinerja.IndikatorResponse
// 		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 			ctx,
// 			tx,
// 			kodeProgram,
// 			strategic.KodeOpd,
// 			strategic.Tahun,
// 		)
// 		if err == nil {
// 			for _, ind := range indikatorProgram {
// 				// Ambil target untuk indikator program
// 				targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 				var targetResponses []pohonkinerja.TargetResponse
// 				if err == nil {
// 					for _, target := range targets {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}
// 				}

// 				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
// 					Id:            ind.Id,
// 					Kode:          ind.Kode,
// 					NamaIndikator: ind.Indikator,
// 					Target:        targetResponses, // Menambahkan target ke indikator program
// 				})
// 			}
// 		}

// 		programList = append(programList, pohonkinerja.ProgramResponse{
// 			KodeProgram: kodeProgram,
// 			NamaProgram: namaProgram,
// 			Indikator:   indikatorProgramResponses,
// 		})
// 	}

// 	strategicResp := pohonkinerja.StrategicCascadingOpdResponse{
// 		Id:         strategic.Id,
// 		Parent:     nil,
// 		Strategi:   strategic.NamaPohon,
// 		JenisPohon: strategic.JenisPohon,
// 		LevelPohon: strategic.LevelPohon,
// 		Keterangan: strategic.Keterangan,
// 		Status:     strategic.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: strategic.KodeOpd,
// 			NamaOpd: strategic.NamaOpd,
// 		},
// 		IsActive:       strategic.IsActive,
// 		Program:        programList,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[strategic.Id],
// 	}

// 	// Build tactical responses dan hitung total pagu anggaran
// 	var totalPaguAnggaran int64 = 0
// 	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
// 		var tacticals []pohonkinerja.TacticalCascadingOpdResponse
// 		sort.Slice(tacticalList, func(i, j int) bool {
// 			// Prioritaskan status "pokin dari pemda"
// 			if tacticalList[i].Status == "pokin dari pemda" && tacticalList[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if tacticalList[i].Status != "pokin dari pemda" && tacticalList[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return tacticalList[i].Id < tacticalList[j].Id
// 		})

// 		for _, tactical := range tacticalList {
// 			tacticalResp := service.buildTacticalCascadingResponse(ctx, tx, pohonMap, tactical, indikatorMap, rencanaKinerjaMap)
// 			tacticals = append(tacticals, tacticalResp)
// 			// Tambahkan pagu anggaran dari setiap tactical
// 			totalPaguAnggaran += tacticalResp.PaguAnggaran
// 		}
// 		strategicResp.Tacticals = tacticals
// 	}

// 	// Set pagu anggaran dari total pagu anggaran tactical
// 	strategicResp.PaguAnggaran = totalPaguAnggaran

// 	return strategicResp
// }

// func (service *CascadingOpdServiceImpl) buildTacticalCascadingResponse(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	tactical domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.TacticalCascadingOpdResponse {

// 	log.Printf("Building tactical response for ID: %d, Level: %d", tactical.Id, tactical.LevelPohon)

// 	// Map untuk menyimpan program unik
// 	programMap := make(map[string]string)

// 	// Proses rencana kinerja untuk tactical
// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaResponse
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[tactical.Id]; ok {
// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse

// 			// Hanya ambil indikator jika rencana kinerja memiliki ID
// 			if rk.Id != "" {
// 				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
// 				if err == nil {
// 					for _, ind := range indikatorRekin {
// 						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}
// 						indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							IdRekin:       rk.Id,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            tactical.Id,
// 				NamaPohon:          tactical.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              tactical.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	// Proses program dari level operational
// 	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
// 		for _, operational := range operationalList {
// 			if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 				for _, rk := range rencanaKinerjaList {
// 					if rk.KodeSubKegiatan != "" {
// 						segments := strings.Split(rk.KodeSubKegiatan, ".")
// 						if len(segments) >= 3 {
// 							kodeProgram := strings.Join(segments[:3], ".")
// 							if _, exists := programMap[kodeProgram]; !exists {
// 								program, err := service.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
// 								if err == nil {
// 									programMap[kodeProgram] = program.NamaProgram
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// Convert program map ke slice dan sort
// 	var programList []pohonkinerja.ProgramResponse
// 	for kodeProgram, namaProgram := range programMap {
// 		// Ambil indikator program
// 		var indikatorProgramResponses []pohonkinerja.IndikatorResponse
// 		indikatorProgram, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 			ctx,
// 			tx,
// 			kodeProgram,
// 			tactical.KodeOpd,
// 			tactical.Tahun,
// 		)
// 		if err == nil {
// 			for _, ind := range indikatorProgram {
// 				// Ambil target untuk indikator program
// 				targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 				var targetResponses []pohonkinerja.TargetResponse
// 				if err == nil {
// 					for _, target := range targets {
// 						targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 							Id:              target.Id,
// 							IndikatorId:     target.IndikatorId,
// 							TargetIndikator: target.Target,
// 							SatuanIndikator: target.Satuan,
// 						})
// 					}
// 				}

// 				indikatorProgramResponses = append(indikatorProgramResponses, pohonkinerja.IndikatorResponse{
// 					Id:            ind.Id,
// 					Kode:          ind.Kode,
// 					NamaIndikator: ind.Indikator,
// 					Target:        targetResponses,
// 				})
// 			}
// 		}

// 		programList = append(programList, pohonkinerja.ProgramResponse{
// 			KodeProgram: kodeProgram,
// 			NamaProgram: namaProgram,
// 			Indikator:   indikatorProgramResponses,
// 		})
// 	}

// 	tacticalResp := pohonkinerja.TacticalCascadingOpdResponse{
// 		Id:         tactical.Id,
// 		Parent:     tactical.Parent,
// 		Strategi:   tactical.NamaPohon,
// 		JenisPohon: tactical.JenisPohon,
// 		LevelPohon: tactical.LevelPohon,
// 		Keterangan: tactical.Keterangan,
// 		Status:     tactical.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: tactical.KodeOpd,
// 			NamaOpd: tactical.NamaOpd,
// 		},
// 		IsActive:       tactical.IsActive,
// 		Program:        programList,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[tactical.Id],
// 	}

// 	// Build operational responses dan hitung total pagu anggaran
// 	var totalPaguAnggaran int64 = 0
// 	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
// 		var operationals []pohonkinerja.OperationalCascadingOpdResponse
// 		sort.Slice(operationalList, func(i, j int) bool {
// 			// Prioritaskan status "pokin dari pemda"
// 			if operationalList[i].Status == "pokin dari pemda" && operationalList[j].Status != "pokin dari pemda" {
// 				return true
// 			}
// 			if operationalList[i].Status != "pokin dari pemda" && operationalList[j].Status == "pokin dari pemda" {
// 				return false
// 			}
// 			return operationalList[i].Id < operationalList[j].Id
// 		})

// 		for _, operational := range operationalList {
// 			operationalResp := service.buildOperationalCascadingResponse(ctx, tx, pohonMap, operational, indikatorMap, rencanaKinerjaMap)
// 			operationals = append(operationals, operationalResp)
// 			// Tambahkan total anggaran dari setiap operational
// 			totalPaguAnggaran += operationalResp.TotalAnggaran
// 		}
// 		tacticalResp.Operationals = operationals
// 	}

// 	// Set pagu anggaran dari total anggaran operational
// 	tacticalResp.PaguAnggaran = totalPaguAnggaran

// 	return tacticalResp
// }

// func (service *CascadingOpdServiceImpl) buildOperationalCascadingResponse(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	operational domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.OperationalCascadingOpdResponse {

// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalResponse
// 	var totalAnggaranOperational int64 = 0
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[operational.Id]; ok {
// 		for _, rk := range rencanaKinerjaList {
// 			var totalAnggaranRenkin int64 = 0
// 			if rk.Id != "" {
// 				// Ambil semua rencana aksi untuk rencana kinerja ini
// 				rencanaAksiList, err := service.rencanaAksiRepository.FindAll(ctx, tx, rk.Id)
// 				if err == nil {
// 					for _, ra := range rencanaAksiList {
// 						// Ambil anggaran untuk setiap rencana aksi
// 						rincianBelanja, err := service.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, ra.Id)
// 						if err == nil {
// 							totalAnggaranRenkin += rincianBelanja.Anggaran
// 						}
// 					}
// 				}
// 			}

// 			totalAnggaranOperational += totalAnggaranRenkin

// 			// Indikator rencana kinerja
// 			var indikatorRekinResponses []pohonkinerja.IndikatorResponse
// 			if rk.Id != "" {
// 				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
// 				if err == nil {
// 					for _, ind := range indikatorRekin {
// 						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}
// 						indikatorRekinResponses = append(indikatorRekinResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							IdRekin:       rk.Id,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}

// 			// Indikator subkegiatan
// 			var indikatorSubkegiatanResponses []pohonkinerja.IndikatorResponse
// 			if rk.KodeSubKegiatan != "" {
// 				indikatorSubkegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 					ctx,
// 					tx,
// 					rk.KodeSubKegiatan,
// 					operational.KodeOpd,
// 					operational.Tahun,
// 				)
// 				if err == nil {
// 					for _, ind := range indikatorSubkegiatan {
// 						// Ambil target untuk indikator
// 						targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}

// 						indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							Kode:          ind.Kode,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}

// 			// Indikator kegiatan
// 			var indikatorKegiatanResponses []pohonkinerja.IndikatorResponse
// 			if rk.KodeKegiatan != "" {
// 				indikatorKegiatan, err := service.cascadingOpdRepository.FindByKodeAndOpdAndTahun(
// 					ctx,
// 					tx,
// 					rk.KodeKegiatan,
// 					operational.KodeOpd,
// 					operational.Tahun,
// 				)
// 				if err == nil {
// 					for _, ind := range indikatorKegiatan {
// 						// Ambil target untuk indikator
// 						targets, err := service.cascadingOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}

// 						indikatorKegiatanResponses = append(indikatorKegiatanResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							Kode:          ind.Kode,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}
// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalResponse{
// 				Id:                   rk.Id,
// 				IdPohon:              operational.Id,
// 				NamaPohon:            operational.NamaPohon,
// 				NamaRencanaKinerja:   rk.NamaRencanaKinerja,
// 				Tahun:                operational.Tahun,
// 				PegawaiId:            rk.PegawaiId,
// 				NamaPegawai:          rk.NamaPegawai,
// 				KodeSubkegiatan:      rk.KodeSubKegiatan,
// 				NamaSubkegiatan:      rk.NamaSubKegiatan,
// 				Anggaran:             totalAnggaranRenkin,
// 				IndikatorSubkegiatan: indikatorSubkegiatanResponses,
// 				KodeKegiatan:         rk.KodeKegiatan,
// 				NamaKegiatan:         rk.NamaKegiatan,
// 				IndikatorKegiatan:    indikatorKegiatanResponses,
// 				Indikator:            indikatorRekinResponses,
// 			})
// 		}
// 	}

// 	operationalResp := pohonkinerja.OperationalCascadingOpdResponse{
// 		Id:         operational.Id,
// 		Parent:     operational.Parent,
// 		Strategi:   operational.NamaPohon,
// 		JenisPohon: operational.JenisPohon,
// 		LevelPohon: operational.LevelPohon,
// 		Keterangan: operational.Keterangan,
// 		Status:     operational.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: operational.KodeOpd,
// 			NamaOpd: operational.NamaOpd,
// 		},
// 		IsActive:       operational.IsActive,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[operational.Id],
// 		TotalAnggaran:  totalAnggaranOperational,
// 	}

// 	// Build operational N responses jika ada
// 	nextLevel := operational.LevelPohon + 1
// 	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
// 		var childs []pohonkinerja.OperationalNOpdCascadingResponse
// 		sort.Slice(operationalNList, func(i, j int) bool {
// 			return operationalNList[i].Id < operationalNList[j].Id
// 		})

// 		for _, opN := range operationalNList {
// 			childResp := service.buildOperationalNCascadingResponse(ctx, tx, pohonMap, opN, indikatorMap, rencanaKinerjaMap)
// 			childs = append(childs, childResp)
// 		}
// 		operationalResp.Childs = childs
// 	}

// 	return operationalResp
// }

// func (service *CascadingOpdServiceImpl) buildOperationalNCascadingResponse(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	pohonMap map[int]map[int][]domain.PohonKinerja,
// 	operationalN domain.PohonKinerja,
// 	indikatorMap map[int][]pohonkinerja.IndikatorResponse,
// 	rencanaKinerjaMap map[int][]domain.RencanaKinerja) pohonkinerja.OperationalNOpdCascadingResponse {

// 	log.Printf("Building OperationalN response for ID: %d, Level: %d", operationalN.Id, operationalN.LevelPohon)

// 	var rencanaKinerjaResponses []pohonkinerja.RencanaKinerjaOperationalNResponse

// 	// Proses rencana kinerja yang ada
// 	if rencanaKinerjaList, ok := rencanaKinerjaMap[operationalN.Id]; ok {
// 		log.Printf("Found %d rencana kinerja for OperationalN ID %d", len(rencanaKinerjaList), operationalN.Id)

// 		for _, rk := range rencanaKinerjaList {
// 			var indikatorResponses []pohonkinerja.IndikatorResponse

// 			// Ambil indikator jika rencana kinerja memiliki ID
// 			if rk.Id != "" {
// 				indikatorRekin, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rk.Id)
// 				if err == nil {
// 					for _, ind := range indikatorRekin {
// 						targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
// 						var targetResponses []pohonkinerja.TargetResponse
// 						if err == nil {
// 							for _, target := range targets {
// 								targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
// 									Id:              target.Id,
// 									IndikatorId:     target.IndikatorId,
// 									TargetIndikator: target.Target,
// 									SatuanIndikator: target.Satuan,
// 								})
// 							}
// 						}
// 						indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
// 							Id:            ind.Id,
// 							IdRekin:       rk.Id,
// 							NamaIndikator: ind.Indikator,
// 							Target:        targetResponses,
// 						})
// 					}
// 				}
// 			}

// 			rencanaKinerjaResponses = append(rencanaKinerjaResponses, pohonkinerja.RencanaKinerjaOperationalNResponse{
// 				Id:                 rk.Id,
// 				IdPohon:            operationalN.Id,
// 				NamaPohon:          operationalN.NamaPohon,
// 				NamaRencanaKinerja: rk.NamaRencanaKinerja,
// 				Tahun:              operationalN.Tahun,
// 				PegawaiId:          rk.PegawaiId,
// 				NamaPegawai:        rk.NamaPegawai,
// 				Indikator:          indikatorResponses,
// 			})
// 		}
// 	}

// 	// Buat response
// 	operationalNResp := pohonkinerja.OperationalNOpdCascadingResponse{
// 		Id:         operationalN.Id,
// 		Parent:     operationalN.Parent,
// 		Strategi:   operationalN.NamaPohon,
// 		JenisPohon: operationalN.JenisPohon,
// 		LevelPohon: operationalN.LevelPohon,
// 		Keterangan: operationalN.Keterangan,
// 		Status:     operationalN.Status,
// 		KodeOpd: opdmaster.OpdResponseForAll{
// 			KodeOpd: operationalN.KodeOpd,
// 			NamaOpd: operationalN.NamaOpd,
// 		},
// 		IsActive:       operationalN.IsActive,
// 		RencanaKinerja: rencanaKinerjaResponses,
// 		Indikator:      indikatorMap[operationalN.Id],
// 	}

// 	// Proses child nodes jika ada
// 	nextLevel := operationalN.LevelPohon + 1
// 	if childList := pohonMap[nextLevel][operationalN.Id]; len(childList) > 0 {
// 		var childs []pohonkinerja.OperationalNOpdCascadingResponse
// 		sort.Slice(childList, func(i, j int) bool {
// 			return childList[i].Id < childList[j].Id
// 		})

// 		for _, child := range childList {
// 			childResp := service.buildOperationalNCascadingResponse(
// 				ctx,
// 				tx,
// 				pohonMap,
// 				child,
// 				indikatorMap,
// 				rencanaKinerjaMap,
// 			)
// 			childs = append(childs, childResp)
// 		}
// 		operationalNResp.Childs = childs
// 	}

// 	return operationalNResp
// }

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
