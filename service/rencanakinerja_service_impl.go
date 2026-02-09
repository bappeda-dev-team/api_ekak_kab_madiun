package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/permasalahan"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RencanaKinerjaServiceImpl struct {
	rencanaKinerjaRepository         repository.RencanaKinerjaRepository
	DB                               *sql.DB
	Validate                         *validator.Validate
	opdRepository                    repository.OpdRepository
	RencanaAksiRepository            repository.RencanaAksiRepository
	UsulanMusrebangRepository        repository.UsulanMusrebangRepository
	UsulanMandatoriRepository        repository.UsulanMandatoriRepository
	UsulanPokokPikiranRepository     repository.UsulanPokokPikiranRepository
	UsulanInisiatifRepository        repository.UsulanInisiatifRepository
	SubKegiatanRepository            repository.SubKegiatanRepository
	SubKegiatanTerpilihRepository    repository.SubKegiatanTerpilihRepository
	DasarHukumRepository             repository.DasarHukumRepository
	GambaranUmumRepository           repository.GambaranUmumRepository
	InovasiRepository                repository.InovasiRepository
	PelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository
	pegawaiRepository                repository.PegawaiRepository
	pohonKinerjaRepository           repository.PohonKinerjaRepository
	manualIKRepository               repository.ManualIKRepository
	permasalahanRekinRepository      repository.PermasalahanRekinRepository
	SubKegiatanService               *SubKegiatanServiceImpl
	PeriodeRepository                repository.PeriodeRepository
	SasaranOpdRepository             repository.SasaranOpdRepository
	CascadingOpdService              *CascadingOpdServiceImpl
	// TAMBAHKAN REPOSITORY BARU:
	cascadingOpdRepository   repository.CascadingOpdRepository
	programRepository        repository.ProgramRepository
	rincianBelanjaRepository repository.RincianBelanjaRepository
	rencanaAksiRepository    repository.RencanaAksiRepository
}

func NewRencanaKinerjaServiceImpl(rencanaKinerjaRepository repository.RencanaKinerjaRepository, DB *sql.DB, validate *validator.Validate, opdRepository repository.OpdRepository, usulanMusrebangRepository repository.UsulanMusrebangRepository, usulanMandatoriRepository repository.UsulanMandatoriRepository, usulanPokokPikiranRepository repository.UsulanPokokPikiranRepository, usulanInisiatifRepository repository.UsulanInisiatifRepository, subKegiatanRepository repository.SubKegiatanRepository, dasarHukumRepository repository.DasarHukumRepository, gambaranUmumRepository repository.GambaranUmumRepository, inovasiRepository repository.InovasiRepository, pelaksanaanRencanaAksiRepository repository.PelaksanaanRencanaAksiRepository, pegawaiRepository repository.PegawaiRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, manualIKRepository repository.ManualIKRepository, permasalahanRekinRepository repository.PermasalahanRekinRepository, subKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository, subKegiatanService *SubKegiatanServiceImpl, periodeRepository repository.PeriodeRepository, sasaranOpdRepository repository.SasaranOpdRepository, cascadingOpdService *CascadingOpdServiceImpl, cascadingOpdRepository repository.CascadingOpdRepository, programRepository repository.ProgramRepository, rincianBelanjaRepository repository.RincianBelanjaRepository, rencanaAksiRepository repository.RencanaAksiRepository) *RencanaKinerjaServiceImpl {
	return &RencanaKinerjaServiceImpl{
		rencanaKinerjaRepository:         rencanaKinerjaRepository,
		DB:                               DB,
		Validate:                         validate,
		opdRepository:                    opdRepository,
		RencanaAksiRepository:            rencanaAksiRepository,
		UsulanMusrebangRepository:        usulanMusrebangRepository,
		UsulanMandatoriRepository:        usulanMandatoriRepository,
		UsulanPokokPikiranRepository:     usulanPokokPikiranRepository,
		UsulanInisiatifRepository:        usulanInisiatifRepository,
		SubKegiatanRepository:            subKegiatanRepository,
		DasarHukumRepository:             dasarHukumRepository,
		GambaranUmumRepository:           gambaranUmumRepository,
		InovasiRepository:                inovasiRepository,
		PelaksanaanRencanaAksiRepository: pelaksanaanRencanaAksiRepository,
		pegawaiRepository:                pegawaiRepository,
		pohonKinerjaRepository:           pohonKinerjaRepository,
		manualIKRepository:               manualIKRepository,
		permasalahanRekinRepository:      permasalahanRekinRepository,
		SubKegiatanTerpilihRepository:    subKegiatanTerpilihRepository,
		SubKegiatanService:               subKegiatanService,
		PeriodeRepository:                periodeRepository,
		SasaranOpdRepository:             sasaranOpdRepository,
		CascadingOpdService:              cascadingOpdService,
		// TAMBAHKAN ASSIGNMENT BARU:
		cascadingOpdRepository:   cascadingOpdRepository,
		programRepository:        programRepository,
		rincianBelanjaRepository: rincianBelanjaRepository,
		rencanaAksiRepository:    rencanaAksiRepository,
	}
}

func (service *RencanaKinerjaServiceImpl) Create(ctx context.Context, request rencanakinerja.RencanaKinerjaCreateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Create RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Perbaikan pengecekan kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Kode OPD %s tidak ditemukan", request.KodeOpd)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak ditemukan", request.KodeOpd)
		}
		log.Printf("Gagal memeriksa kode OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memeriksa kode OPD: %v", err)
	}

	if opd.KodeOpd == "" {
		log.Printf("Kode OPD %s tidak valid", request.KodeOpd)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak valid", request.KodeOpd)
	}

	pegawais, err := service.pegawaiRepository.FindByNip(ctx, tx, request.PegawaiId)
	if err != nil {
		log.Printf("Gagal mengambil data pegawai: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	if pegawais.Id == "" {
		log.Printf("Pegawai dengan Nip %s tidak ditemukan", request.PegawaiId)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pegawai dengan Nip %s tidak ditemukan", request.PegawaiId)
	}

	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, request.IdPohon)
	if err != nil {
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	if pohon.Id == 0 {
		log.Printf("Pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
	}

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	year := time.Now().Year()
	customId := fmt.Sprintf("REKIN-PEG-%v-%v", year, randomDigits)

	rencanaKinerja := domain.RencanaKinerja{
		Id:                   customId,
		IdPohon:              request.IdPohon,
		NamaRencanaKinerja:   request.NamaRencanaKinerja,
		Tahun:                request.Tahun,
		StatusRencanaKinerja: request.StatusRencanaKinerja,
		Catatan:              request.Catatan,
		KodeOpd:              request.KodeOpd,
		PegawaiId:            pegawais.Nip,
		KodeSubKegiatan:      "",
		TahunAwal:            "",
		TahunAkhir:           "",
		JenisPeriode:         "",
		PeriodeId:            0,
		Indikator:            make([]domain.Indikator, len(request.Indikator)),
	}

	log.Printf("RencanaKinerja dibuat dengan ID: %s", customId)

	for i, indikatorRequest := range request.Indikator {
		indikatorRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		indikatorId := fmt.Sprintf("IND-REKIN-%s", indikatorRandomDigits)
		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.NamaIndikator,
			Tahun:            request.Tahun,
			Target:           make([]domain.Target, len(indikatorRequest.Target)),
			RencanaKinerjaId: rencanaKinerja.Id,
		}

		if indikator.Indikator == "" {
			log.Printf("Indikator kosong ditemukan: %+v", indikator)
		}

		log.Printf("Indikator dibuat: %+v", indikator)

		for j, targetRequest := range indikatorRequest.Target {
			targetRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			targetId := fmt.Sprintf("TRGT-IND-REKIN-%s", targetRandomDigits)
			target := domain.Target{
				Id:          targetId,
				Tahun:       request.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
			log.Printf("Target dibuat dengan ID: %s", targetId)
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Create")
	rencanaKinerja, err = service.rencanaKinerjaRepository.Create(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal menyimpan RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menyimpan RencanaKinerja: %v", err)
	}

	rencanaKinerja.NamaOpd = opd.NamaOpd
	rencanaKinerja.NamaPegawai = pegawais.NamaPegawai
	rencanaKinerja.NamaPohon = pohon.NamaPohon
	log.Println("RencanaKinerja berhasil disimpan")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) Update(ctx context.Context, request rencanakinerja.RencanaKinerjaUpdateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Update RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Kode OPD %s tidak ditemukan", request.KodeOpd)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak ditemukan", request.KodeOpd)
		}
		log.Printf("Gagal memeriksa kode OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memeriksa kode OPD: %v", err)
	}

	if opd.KodeOpd == "" {
		log.Printf("Kode OPD %s tidak valid", request.KodeOpd)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak valid", request.KodeOpd)
	}

	// Validasi Pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, request.PegawaiId)
	if err != nil {
		log.Printf("Gagal mengambil data pegawai: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	if pegawai.Id == "" {
		log.Printf("Pegawai dengan NIP %s tidak ditemukan", request.PegawaiId)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pegawai dengan NIP %s tidak ditemukan", request.PegawaiId)
	}

	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, request.IdPohon)
	if err != nil {
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	if pohon.Id == 0 {
		log.Printf("Pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
	}

	//
	var rencanaKinerja domain.RencanaKinerja
	if request.Id != "" {
		rencanaKinerja, err = service.rencanaKinerjaRepository.FindById(ctx, tx, request.Id, "", "")
		if err != nil {
			log.Printf("Gagal menemukan RencanaKinerja: %v", err)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan RencanaKinerja: %v", err)
		}
	} else {
		randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		rencanaKinerja.Id = fmt.Sprintf("REKIN-PEG-%s", randomDigits)
		log.Printf("Membuat RencanaKinerja baru dengan ID: %s", rencanaKinerja.Id)
	}

	rencanaKinerja.IdPohon = request.IdPohon
	rencanaKinerja.NamaRencanaKinerja = request.NamaRencanaKinerja
	rencanaKinerja.Tahun = request.Tahun
	rencanaKinerja.StatusRencanaKinerja = request.StatusRencanaKinerja
	rencanaKinerja.Catatan = request.Catatan
	rencanaKinerja.KodeOpd = request.KodeOpd
	rencanaKinerja.PegawaiId = request.PegawaiId

	rencanaKinerja.Indikator = make([]domain.Indikator, len(request.Indikator))
	for i, indikatorRequest := range request.Indikator {
		var indikatorId string
		if indikatorRequest.Id != "" {
			indikatorId = indikatorRequest.Id
		} else {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-REKIN-%s", randomDigits)
			log.Printf("Membuat Indikator baru dengan ID: %s", indikatorId)
		}

		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.Indikator,
			Tahun:            rencanaKinerja.Tahun,
			RencanaKinerjaId: rencanaKinerja.Id,
		}

		indikator.Target = make([]domain.Target, len(indikatorRequest.Target))
		for j, targetRequest := range indikatorRequest.Target {
			var targetId string
			if targetRequest.Id != "" {
				targetId = targetRequest.Id
			} else {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRGT-IND-REKIN-%s", randomDigits)
				log.Printf("Membuat Target baru dengan ID: %s", targetId)
			}

			target := domain.Target{
				Id:          targetId,
				Tahun:       rencanaKinerja.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Update")
	rencanaKinerja, err = service.rencanaKinerjaRepository.Update(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal memperbarui RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memperbarui RencanaKinerja: %v", err)
	}

	rencanaKinerja.NamaOpd = opd.NamaOpd
	rencanaKinerja.NamaPegawai = pegawai.NamaPegawai
	rencanaKinerja.NamaPohon = pohon.NamaPohon

	log.Println("RencanaKinerja berhasil diperbarui")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) FindAll(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses FindAll RencanaKinerja")

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	log.Printf("Mencari RencanaKinerja dengan pegawaiId: %s, kodeOPD: %s, tahun: %s", pegawaiId, kodeOPD, tahun)
	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindAll(ctx, tx, pegawaiId, kodeOPD, tahun)
	if err != nil {
		log.Printf("Gagal mencari RencanaKinerja: %v", err)
		return nil, fmt.Errorf("gagal mencari RencanaKinerja: %v", err)
	}
	log.Printf("Ditemukan %d RencanaKinerja", len(rencanaKinerjaList))

	// Batch query untuk OPD, Pegawai, dan Pohon Kinerja
	// Kumpulkan semua kode_opd, pegawai_id, dan id_pohon yang unik
	kodeOpdSet := make(map[string]bool)
	pegawaiIdSet := make(map[string]bool)
	pohonIdSet := make(map[int]bool)

	for _, rencana := range rencanaKinerjaList {
		if rencana.KodeOpd != "" {
			kodeOpdSet[rencana.KodeOpd] = true
		}
		if rencana.PegawaiId != "" {
			pegawaiIdSet[rencana.PegawaiId] = true
		}
		if rencana.IdPohon != 0 {
			pohonIdSet[rencana.IdPohon] = true
		}
	}

	// Batch query OPD
	opdMap := make(map[string]domainmaster.Opd)
	if len(kodeOpdSet) > 0 {
		kodeOpdList := make([]string, 0, len(kodeOpdSet))
		for kode := range kodeOpdSet {
			kodeOpdList = append(kodeOpdList, kode)
		}
		// Jika ada method FindByKodeOpds batch, gunakan itu
		// Jika tidak, tetap gunakan loop (tapi sudah di-optimize dengan map)
		for _, kode := range kodeOpdList {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kode)
			if err == nil {
				opdMap[kode] = opd
			}
		}
	}

	// Batch query Pegawai
	pegawaiMap := make(map[string]domainmaster.Pegawai)
	if len(pegawaiIdSet) > 0 {
		pegawaiIdList := make([]string, 0, len(pegawaiIdSet))
		for nip := range pegawaiIdSet {
			pegawaiIdList = append(pegawaiIdList, nip)
		}
		for _, nip := range pegawaiIdList {
			pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, nip)
			if err == nil {
				pegawaiMap[nip] = pegawai
			}
		}
	}

	// Batch query Pohon Kinerja
	pohonIdList := make([]int, 0, len(pohonIdSet))
	for id := range pohonIdSet {
		pohonIdList = append(pohonIdList, id)
	}
	pohonMap, err := service.pohonKinerjaRepository.FindByIds(ctx, tx, pohonIdList)
	if err != nil {
		log.Printf("Gagal batch query Pohon Kinerja: %v", err)
		// Fallback ke individual query jika batch gagal
		pohonMap = make(map[int]domain.PohonKinerja)
	}

	var responses []rencanakinerja.RencanaKinerjaResponse
	for _, rencana := range rencanaKinerjaList {
		log.Printf("Memproses RencanaKinerja dengan ID: %s", rencana.Id)

		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencana.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Gagal mencari Indikator: %v", err)
			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
		}

		var indikatorResponses []rencanakinerja.IndikatorResponse
		for _, indikator := range indikators {
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Gagal mencari Target: %v", err)
				return nil, fmt.Errorf("gagal mencari Target: %v", err)
			}

			var targetResponses []rencanakinerja.TargetResponse
			for _, target := range targets {
				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			exist, err := service.manualIKRepository.IsIndikatorExist(ctx, tx, indikator.Id)
			if err != nil {
				log.Printf("Gagal memeriksa keberadaan indikator: %v", err)
				return nil, fmt.Errorf("gagal memeriksa keberadaan indikator: %v", err)
			}

			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
				ManualIKExist:    exist,
			})
		}

		// Ambil dari map yang sudah di-load
		opd, opdExists := opdMap[rencana.KodeOpd]
		if !opdExists {
			log.Printf("OPD tidak ditemukan untuk kode: %s", rencana.KodeOpd)
			opd = domainmaster.Opd{} // Default empty
		}

		pegawai, pegawaiExists := pegawaiMap[rencana.PegawaiId]
		if !pegawaiExists {
			log.Printf("Pegawai tidak ditemukan untuk NIP: %s", rencana.PegawaiId)
			pegawai = domainmaster.Pegawai{} // Default empty
		}

		pohon, pohonExists := pohonMap[rencana.IdPohon]
		pohonFound := pohonExists

		if !pohonExists && rencana.IdPohon != 0 {
			// Fallback ke individual query jika tidak ada di map
			pohonData, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencana.IdPohon)
			if err == nil && pohonData.Id != 0 {
				pohon = pohonData
				pohonMap[rencana.IdPohon] = pohon
				pohonFound = true
			} else {
				// Pohon kinerja tidak ditemukan di database
				pohonFound = false
				log.Printf("Pohon Kinerja dengan ID %d tidak ditemukan di database", rencana.IdPohon)
			}
		}

		// Logika: Cek apakah perlu ubah pohon kinerja
		perluUbahPohonKinerja := false
		if tahun == "2026" && rencana.IdPohon != 0 {
			if !pohonFound {
				// Jika pohon kinerja tidak ditemukan di database, perlu ubah
				perluUbahPohonKinerja = true
			} else if pohon.Tahun != "" && pohon.Tahun != tahun {
				// Jika tahun pohon kinerja berbeda dengan tahun yang dicari, perlu ubah
				perluUbahPohonKinerja = true
			}
		}

		responses = append(responses, rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencana.Id,
			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
			Tahun:                rencana.Tahun,
			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
			Catatan:              rencana.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			PegawaiId:             rencana.PegawaiId,
			NamaPegawai:           pegawai.NamaPegawai,
			IdPohon:               rencana.IdPohon,
			NamaPohon:             pohon.NamaPohon,
			LevelPohon:            pohon.LevelPohon,
			Indikator:             indikatorResponses,
			PerluUbahPohonKinerja: perluUbahPohonKinerja,
		})
		log.Printf("RencanaKinerja Response ditambahkan untuk ID: %s", rencana.Id)
	}

	return responses, nil
}

// func (service *RencanaKinerjaServiceImpl) FindAll(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error) {
// 	log.Println("Memulai proses FindAll RencanaKinerja")

// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		log.Printf("Gagal memulai transaksi: %v", err)
// 		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	log.Printf("Mencari RencanaKinerja dengan pegawaiId: %s, kodeOPD: %s, tahun: %s", pegawaiId, kodeOPD, tahun)
// 	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindAll(ctx, tx, pegawaiId, kodeOPD, tahun)
// 	if err != nil {
// 		log.Printf("Gagal mencari RencanaKinerja: %v", err)
// 		return nil, fmt.Errorf("gagal mencari RencanaKinerja: %v", err)
// 	}
// 	log.Printf("Ditemukan %d RencanaKinerja", len(rencanaKinerjaList))

// 	var responses []rencanakinerja.RencanaKinerjaResponse
// 	for _, rencana := range rencanaKinerjaList {
// 		log.Printf("Memproses RencanaKinerja dengan ID: %s", rencana.Id)

// 		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencana.Id)
// 		if err != nil && err != sql.ErrNoRows {
// 			log.Printf("Gagal mencari Indikator: %v", err)
// 			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
// 		}

// 		var indikatorResponses []rencanakinerja.IndikatorResponse
// 		for _, indikator := range indikators {
// 			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
// 			if err != nil && err != sql.ErrNoRows {
// 				log.Printf("Gagal mencari Target: %v", err)
// 				return nil, fmt.Errorf("gagal mencari Target: %v", err)
// 			}

// 			var targetResponses []rencanakinerja.TargetResponse
// 			for _, target := range targets {
// 				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
// 					Id:              target.Id,
// 					IndikatorId:     target.IndikatorId,
// 					TargetIndikator: target.Target,
// 					SatuanIndikator: target.Satuan,
// 				})
// 			}

// 			exist, err := service.manualIKRepository.IsIndikatorExist(ctx, tx, indikator.Id)
// 			if err != nil {
// 				log.Printf("Gagal memeriksa keberadaan indikator: %v", err)
// 				return nil, fmt.Errorf("gagal memeriksa keberadaan indikator: %v", err)
// 			}

// 			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
// 				Id:               indikator.Id,
// 				RencanaKinerjaId: indikator.RencanaKinerjaId,
// 				NamaIndikator:    indikator.Indikator,
// 				Target:           targetResponses,
// 				ManualIKExist:    exist,
// 			})
// 		}

// 		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencana.KodeOpd)
// 		if err != nil {
// 			log.Printf("Gagal mencari OPD: %v", err)
// 			return nil, fmt.Errorf("gagal mencari OPD: %v", err)
// 		}

// 		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencana.PegawaiId)
// 		if err != nil {
// 			log.Printf("Gagal mencari Pegawai: %v", err)
// 			return nil, fmt.Errorf("gagal mencari Pegawai: %v", err)
// 		}

// 		pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencana.IdPohon)
// 		if err != nil {
// 			log.Printf("Gagal mencari Pohon Kinerja: %v", err)
// 			return nil, fmt.Errorf("gagal mencari Pohon Kinerja: %v", err)
// 		}

// 		responses = append(responses, rencanakinerja.RencanaKinerjaResponse{
// 			Id:                   rencana.Id,
// 			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
// 			Tahun:                rencana.Tahun,
// 			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
// 			Catatan:              rencana.Catatan,
// 			KodeOpd: opdmaster.OpdResponseForAll{
// 				KodeOpd: opd.KodeOpd,
// 				NamaOpd: opd.NamaOpd,
// 			},
// 			PegawaiId:   rencana.PegawaiId,
// 			NamaPegawai: pegawai.NamaPegawai,
// 			IdPohon:     rencana.IdPohon,
// 			NamaPohon:   pohon.NamaPohon,
// 			LevelPohon:  pohon.LevelPohon,
// 			Indikator:   indikatorResponses,
// 		})
// 		log.Printf("RencanaKinerja Response ditambahkan untuk ID: %s", rencana.Id)
// 	}

// 	return responses, nil
// }

func (service *RencanaKinerjaServiceImpl) FindById(ctx context.Context, id string, kodeOPD string, tahun string) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Printf("Mencari RencanaKinerja dengan ID: %s, KodeOPD: %s, Tahun: %s", id, kodeOPD, tahun)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	rencanaKinerja, err := service.rencanaKinerjaRepository.FindById(ctx, tx, id, kodeOPD, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("RencanaKinerja tidak ditemukan untuk ID: %s", id)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("rencana kinerja tidak ditemukan")
		}
		log.Printf("Gagal menemukan rencana kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan rencana kinerja: %v", err)
	}

	log.Printf("RencanaKinerja ditemukan: %+v", rencanaKinerja)

	indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencanaKinerja.Id)
	if err != nil {
		log.Printf("Gagal menemukan indikator: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan indikator: %v", err)
	}
	rencanaKinerja.Indikator = indikators

	log.Printf("Jumlah indikator ditemukan: %d", len(indikators))

	for i, indikator := range rencanaKinerja.Indikator {
		targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			log.Printf("Gagal menemukan target untuk indikator %s: %v", indikator.Id, err)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan target untuk indikator %s: %v", indikator.Id, err)
		}
		rencanaKinerja.Indikator[i].Target = targets
		log.Printf("Jumlah target ditemukan untuk indikator %s: %d", indikator.Id, len(targets))
	}

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerja.KodeOpd)
	if err != nil {
		log.Printf("Gagal mengambil data OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
	}

	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencanaKinerja.PegawaiId)
	if err != nil {
		log.Printf("Gagal mengambil data pegawai: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencanaKinerja.IdPohon)
	if err != nil {
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Set semua data yang diperlukan ke dalam rencanaKinerja
	rencanaKinerja.NamaOpd = opd.NamaOpd
	rencanaKinerja.NamaPegawai = pegawai.NamaPegawai
	rencanaKinerja.NamaPohon = pohon.NamaPohon

	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	rencanaKinerja, err := service.rencanaKinerjaRepository.FindById(ctx, tx, id, "", "")
	if err != nil {
		return err
	}

	return service.rencanaKinerjaRepository.Delete(ctx, tx, rencanaKinerja.Id)
}

func (service *RencanaKinerjaServiceImpl) FindAllRincianKak(ctx context.Context, pegawaiId string, rencanaKinerjaId string) ([]rencanakinerja.DataRincianKerja, error) {
	log.Println("Memulai proses FindAllRincianKak")

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua rencana kinerja
	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindAllRincianKak(ctx, tx, rencanaKinerjaId, pegawaiId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil rencana kinerja: %v", err)
	}

	var responses []rencanakinerja.DataRincianKerja
	for _, rencanaKinerja := range rencanaKinerjaList {
		// Ambil indikator untuk setiap rencana kinerja
		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencanaKinerja.Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("gagal mengambil indikator: %v", err)
		}

		// Proses indikator dan target
		var indikatorResponses []rencanakinerja.IndikatorResponse
		for _, indikator := range indikators {
			// Tambahkan pengambilan manual IK untuk setiap indikator
			manualIK, err := service.manualIKRepository.FindByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				log.Printf("Warning: gagal mengambil manual IK: %v", err)
			}

			// Filter output data yang true saja
			var outputData []string
			if manualIK.Kinerja {
				outputData = append(outputData, "kinerja")
			}
			if manualIK.Penduduk {
				outputData = append(outputData, "penduduk")
			}
			if manualIK.Spatial {
				outputData = append(outputData, "spatial")
			}
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("gagal mengambil target: %v", err)
			}

			var targetResponses []rencanakinerja.TargetResponse
			for _, target := range targets {
				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
				ManualIK: &rencanakinerja.DataOutput{
					OutputData: outputData,
				},
			})
		}

		// Setelah mengambil data OPD dan sebelum membuat response
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerja.KodeOpd)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data OPD: %v", err)
		}

		// Tambahkan untuk mengambil data pegawai
		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencanaKinerja.PegawaiId)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data pegawai: %v", err)
		}

		// Tambahkan untuk mengambil data pohon kinerja
		pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencanaKinerja.IdPohon)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}
		// Ambil data terkait untuk setiap rencana
		rencanaAksiList, err := service.RencanaAksiRepository.FindAll(ctx, tx, rencanaKinerja.Id)
		if err != nil {
			log.Printf("Warning: gagal mengambil rencana aksi: %v", err)
			rencanaAksiList = []domain.RencanaAksi{}
		}

		// Modifikasi bagian yang memproses rencana aksi
		var rencanaAksiResponses []rencanaaksi.RencanaAksiResponse
		bobotPerBulan := make([]int, 12)    // Array untuk menyimpan total per bulan
		bulanTerpakai := make(map[int]bool) // Map untuk melacak bulan yang digunakan

		for _, rencanaAksi := range rencanaAksiList {
			// Ambil data pelaksanaan untuk setiap rencana aksi
			pelaksanaanList, err := service.PelaksanaanRencanaAksiRepository.FindByRencanaAksiId(ctx, tx, rencanaAksi.Id)
			if err != nil {
				log.Printf("Warning: gagal mengambil pelaksanaan rencana aksi: %v", err)
				pelaksanaanList = []domain.PelaksanaanRencanaAksi{}
			}

			// Buat map untuk menyimpan data pelaksanaan per bulan
			pelaksanaanPerBulan := make(map[int]domain.PelaksanaanRencanaAksi)
			for _, pelaksanaan := range pelaksanaanList {
				pelaksanaanPerBulan[pelaksanaan.Bulan] = pelaksanaan
				if pelaksanaan.Bobot > 0 {
					bulanTerpakai[pelaksanaan.Bulan] = true // Menandai bulan yang digunakan
				}
			}

			// Buat slice pelaksanaan yang terurut untuk 12 bulan
			var pelaksanaanLengkap []domain.PelaksanaanRencanaAksi
			totalBobotRencanaAksi := 0

			for bulan := 1; bulan <= 12; bulan++ {
				if pelaksanaan, exists := pelaksanaanPerBulan[bulan]; exists {
					pelaksanaanLengkap = append(pelaksanaanLengkap, domain.PelaksanaanRencanaAksi{
						Id:            pelaksanaan.Id,
						RencanaAksiId: rencanaAksi.Id,
						Bulan:         bulan,
						Bobot:         pelaksanaan.Bobot,
					})
					totalBobotRencanaAksi += pelaksanaan.Bobot
					bobotPerBulan[bulan-1] += pelaksanaan.Bobot // Menambahkan ke total per bulan
				} else {
					pelaksanaanLengkap = append(pelaksanaanLengkap, domain.PelaksanaanRencanaAksi{
						Id:            "",
						RencanaAksiId: rencanaAksi.Id,
						Bulan:         bulan,
						Bobot:         0,
					})
				}
			}

			response := helper.ToRencanaAksiResponse(rencanaAksi, pelaksanaanLengkap)
			response.TotalBobotRencanaAksi = totalBobotRencanaAksi
			rencanaAksiResponses = append(rencanaAksiResponses, response)
		}

		// Konversi array bobotPerBulan ke slice BobotBulanan
		var totalPerBulanResponse []rencanaaksi.BobotBulanan
		totalKeseluruhan := 0

		// Hitung jumlah bulan unik yang digunakan
		bulanUnik := []int{}
		for bulan := range bulanTerpakai {
			bulanUnik = append(bulanUnik, bulan)
		}

		// Urutkan bulan-bulan yang digunakan
		sort.Ints(bulanUnik)

		for bulan := 1; bulan <= 12; bulan++ {
			bobot := bobotPerBulan[bulan-1]
			totalPerBulanResponse = append(totalPerBulanResponse, rencanaaksi.BobotBulanan{
				Bulan:      bulan,
				TotalBobot: bobot,
			})
			totalKeseluruhan += bobot
		}

		rencanaAksiTable := rencanaaksi.RencanaAksiTableResponse{
			RencanaAksi:      rencanaAksiResponses,
			TotalPerBulan:    totalPerBulanResponse,
			TotalKeseluruhan: totalKeseluruhan,
			WaktuDibutuhkan:  len(bulanUnik), // Jumlah bulan unik yang digunakan
		}

		// Modifikasi bagian subkegiatan
		subKegiatanTerpilihList, err := service.SubKegiatanTerpilihRepository.FindAll(ctx, tx, rencanaKinerja.Id)
		if err != nil {
			log.Printf("Warning: gagal mengambil data subkegiatan terpilih: %v", err)
			log.Printf("Kode subkegiatan: %v", subKegiatanTerpilihList)
			return nil, fmt.Errorf("gagal mengambil data subkegiatan terpilih: %v", err)
		}

		var subKegiatanResponses []subkegiatan.SubKegiatanResponse
		for _, st := range subKegiatanTerpilihList {
			// Menggunakan FindByKodeSubKegiatan alih-alih FindById
			subKegiatanDetail, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, st.KodeSubKegiatan)
			if err != nil {
				log.Printf("Warning: gagal mengambil detail subkegiatan: %v", err)
				continue
			}

			var indikatorResponses []subkegiatan.IndikatorResponse
			for _, indikator := range subKegiatanDetail.Indikator {
				var targetResponses []subkegiatan.TargetResponse
				for _, target := range indikator.Target {
					targetResponses = append(targetResponses, subkegiatan.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					})
				}

				indikatorResponses = append(indikatorResponses, subkegiatan.IndikatorResponse{
					Id:            indikator.Id,
					NamaIndikator: indikator.Indikator,
					Target:        targetResponses,
				})
			}

			subKegiatanResponses = append(subKegiatanResponses, subkegiatan.SubKegiatanResponse{
				SubKegiatanTerpilihId: st.Id,
				Id:                    subKegiatanDetail.Id,
				RekinId:               rencanaKinerja.Id,
				KodeSubKegiatan:       subKegiatanDetail.KodeSubKegiatan,
				NamaSubKegiatan:       subKegiatanDetail.NamaSubKegiatan,
				Indikator:             indikatorResponses,
			})
		}

		var isActive *bool // nil karena tidak perlu filter is_active
		var status *string

		usulanMusrebang, _ := service.UsulanMusrebangRepository.FindAll(ctx, tx, &rencanaKinerja.KodeOpd, isActive, &rencanaKinerja.Id, status)
		usulanMandatori, _ := service.UsulanMandatoriRepository.FindAll(ctx, tx, nil, &pegawaiId, nil, &rencanaKinerja.Id)
		usulanPokokPikiran, _ := service.UsulanPokokPikiranRepository.FindAll(ctx, tx, &rencanaKinerja.KodeOpd, isActive, &rencanaKinerja.Id, status)
		usulanInisiatif, _ := service.UsulanInisiatifRepository.FindAll(ctx, tx, &pegawaiId, nil, &rencanaKinerja.Id)
		dasarHukum, _ := service.DasarHukumRepository.FindAll(ctx, tx, rencanaKinerja.Id)
		gambaranUmum, _ := service.GambaranUmumRepository.FindAll(ctx, tx, rencanaKinerja.Id)
		inovasi, _ := service.InovasiRepository.FindAll(ctx, tx, rencanaKinerja.Id)

		// Gabungkan semua usulan
		var usulanGabungan []rencanakinerja.UsulanGabunganResponse

		// Proses usulan musrebang
		for _, um := range usulanMusrebang {
			usulanGabungan = append(usulanGabungan, rencanakinerja.UsulanGabunganResponse{
				Id:          um.Id,
				Usulan:      um.Usulan,
				Uraian:      um.Uraian,
				JenisUsulan: "usulan_musrebang",
				Tahun:       um.Tahun,
				RekinId:     um.RekinId,
				KodeOpd:     um.KodeOpd,
				IsActive:    um.IsActive,
				Status:      um.Status,
				Alamat:      um.Alamat,
			})
		}

		// Proses usulan pokok pikiran
		for _, up := range usulanPokokPikiran {
			usulanGabungan = append(usulanGabungan, rencanakinerja.UsulanGabunganResponse{
				Id:          up.Id,
				Usulan:      up.Usulan,
				Uraian:      up.Uraian,
				JenisUsulan: "usulan_pokok_pikiran",
				Tahun:       up.Tahun,
				RekinId:     up.RekinId,
				KodeOpd:     up.KodeOpd,
				IsActive:    up.IsActive,
				Status:      up.Status,
				Alamat:      up.Alamat,
			})
		}

		// Proses usulan mandatori
		for _, um := range usulanMandatori {
			usulanGabungan = append(usulanGabungan, rencanakinerja.UsulanGabunganResponse{
				Id:               um.Id,
				Usulan:           um.Usulan,
				Uraian:           um.Uraian,
				JenisUsulan:      "usulan_mandatori",
				Tahun:            um.Tahun,
				RekinId:          um.RekinId,
				PegawaiId:        um.PegawaiId,
				KodeOpd:          um.KodeOpd,
				IsActive:         um.IsActive,
				Status:           um.Status,
				PeraturanTerkait: um.PeraturanTerkait,
			})
		}

		// Proses usulan inisiatif
		for _, ui := range usulanInisiatif {
			usulanGabungan = append(usulanGabungan, rencanakinerja.UsulanGabunganResponse{
				Id:          ui.Id,
				Usulan:      ui.Usulan,
				Uraian:      ui.Uraian,
				JenisUsulan: "usulan_inisiatif",
				Tahun:       ui.Tahun,
				RekinId:     ui.RekinId,
				PegawaiId:   ui.PegawaiId,
				KodeOpd:     ui.KodeOpd,
				IsActive:    ui.IsActive,
				Status:      ui.Status,
				Manfaat:     ui.Manfaat,
			})
		}

		// Buat response untuk setiap rencana kinerja
		rencanaKinerjaResponse := rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencanaKinerja.Id,
			NamaRencanaKinerja:   rencanaKinerja.NamaRencanaKinerja,
			Tahun:                rencanaKinerja.Tahun,
			StatusRencanaKinerja: rencanaKinerja.StatusRencanaKinerja,
			Catatan:              rencanaKinerja.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			PegawaiId:   rencanaKinerja.PegawaiId,
			NamaPegawai: pegawai.NamaPegawai,
			IdPohon:     rencanaKinerja.IdPohon,
			NamaPohon:   pohon.NamaPohon,

			Indikator: indikatorResponses,
		}

		permasalahanRekin, err := service.permasalahanRekinRepository.FindAll(ctx, tx, &rencanaKinerja.Id)
		if err != nil {
			log.Printf("Warning: gagal mengambil permasalahan rekin: %v", err)
			permasalahanRekin = []domain.PermasalahanRekin{}
		}

		var permasalahanResponses []permasalahan.PermasalahanRekinResponse
		for _, p := range permasalahanRekin {
			permasalahanResponses = append(permasalahanResponses, permasalahan.PermasalahanRekinResponse{
				Id:                p.Id,
				RekinId:           p.RekinId,
				Permasalahan:      p.Permasalahan,
				PenyebabInternal:  p.PenyebabInternal,
				PenyebabEksternal: p.PenyebabEksternal,
				JenisPermasalahan: p.JenisPermasalahan,
			})
		}
		// Tambahkan ke responses
		responses = append(responses, rencanakinerja.DataRincianKerja{
			RencanaKinerja: rencanaKinerjaResponse,
			RencanaAksi:    rencanaAksiTable,
			Usulan:         usulanGabungan,
			DasarHukum:     helper.ToDasarHukumResponses(dasarHukum),
			SubKegiatan:    subKegiatanResponses,
			GambaranUmum:   helper.ToGambaranUmumResponses(gambaranUmum),
			Inovasi:        helper.ToInovasiResponses(inovasi),
			Permasalahan:   permasalahanResponses,
		})
	}

	return responses, nil
}

func (service *RencanaKinerjaServiceImpl) RekinsasaranOpd(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses RekinsasaranOpd")

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	log.Printf("Mencari RencanaKinerja dengan pegawaiId: %s, kodeOPD: %s, tahun: %s", pegawaiId, kodeOPD, tahun)
	rencanaKinerjaList, err := service.rencanaKinerjaRepository.RekinsasaranOpd(ctx, tx, pegawaiId, kodeOPD, tahun)
	if err != nil {
		log.Printf("Gagal mencari RencanaKinerja: %v", err)
		return nil, fmt.Errorf("gagal mencari RencanaKinerja: %v", err)
	}
	log.Printf("Ditemukan %d RencanaKinerja", len(rencanaKinerjaList))

	var responses []rencanakinerja.RencanaKinerjaResponse
	for _, rencana := range rencanaKinerjaList {
		log.Printf("Memproses RencanaKinerja dengan ID: %s", rencana.Id)

		tahunInt, _ := strconv.Atoi(tahun)
		tahunAwalInt, _ := strconv.Atoi(rencana.TahunAwal)
		tahunAkhirInt, _ := strconv.Atoi(rencana.TahunAkhir)

		if tahunInt < tahunAwalInt || tahunInt > tahunAkhirInt {
			continue // Skip jika tahun di luar range
		}

		indikators, err := service.rencanaKinerjaRepository.FindIndikatorSasaranbyRekinId(ctx, tx, rencana.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Gagal mencari Indikator: %v", err)
			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
		}

		var indikatorResponses []rencanakinerja.IndikatorResponse
		for _, indikator := range indikators {
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorIdAndTahun(ctx, tx, indikator.Id, tahun)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Gagal mencari Target: %v", err)
				return nil, fmt.Errorf("gagal mencari Target: %v", err)
			}

			var targetResponses []rencanakinerja.TargetResponse
			if len(targets) > 0 {
				for _, target := range targets {
					targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
						Tahun:           target.Tahun,
					})
				}
			} else {
				// Jika tidak ada target untuk tahun tersebut dan tahun dalam range, tambahkan target kosong
				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
					Id:              "",
					IndikatorId:     indikator.Id,
					TargetIndikator: "",
					SatuanIndikator: "",
					Tahun:           tahun,
				})
			}

			exist, err := service.manualIKRepository.IsIndikatorExist(ctx, tx, indikator.Id)
			if err != nil {
				log.Printf("Gagal memeriksa keberadaan indikator: %v", err)
				return nil, fmt.Errorf("gagal memeriksa keberadaan indikator: %v", err)
			}

			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
				ManualIKExist:    exist,
			})
		}

		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencana.KodeOpd)
		if err != nil {
			log.Printf("Gagal mencari OPD: %v", err)
			return nil, fmt.Errorf("gagal mencari OPD: %v", err)
		}

		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencana.PegawaiId)
		if err != nil {
			log.Printf("Gagal mencari Pegawai: %v", err)
			return nil, fmt.Errorf("gagal mencari Pegawai: %v", err)
		}

		pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencana.IdPohon)
		if err != nil {
			log.Printf("Gagal mencari Pohon Kinerja: %v", err)
			return nil, fmt.Errorf("gagal mencari Pohon Kinerja: %v", err)
		}

		responses = append(responses, rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencana.Id,
			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
			Catatan:              rencana.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			PegawaiId:   rencana.PegawaiId,
			NamaPegawai: pegawai.NamaPegawai,
			IdPohon:     rencana.IdPohon,
			NamaPohon:   pohon.NamaPohon,
			Indikator:   indikatorResponses,
		})
		log.Printf("RencanaKinerja Response ditambahkan untuk ID: %s", rencana.Id)
	}

	return responses, nil
}

func (service *RencanaKinerjaServiceImpl) CreateRekinLevel1(ctx context.Context, request rencanakinerja.RencanaKinerjaCreateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Create RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Perbaikan pengecekan kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Kode OPD %s tidak ditemukan", request.KodeOpd)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak ditemukan", request.KodeOpd)
		}
		log.Printf("Gagal memeriksa kode OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memeriksa kode OPD: %v", err)
	}

	if opd.KodeOpd == "" {
		log.Printf("Kode OPD %s tidak valid", request.KodeOpd)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak valid", request.KodeOpd)
	}

	pegawais, err := service.pegawaiRepository.FindByNip(ctx, tx, request.PegawaiId)
	if err != nil {
		log.Printf("Gagal mengambil data pegawai: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	if pegawais.Id == "" {
		log.Printf("Pegawai dengan Nip %s tidak ditemukan", request.PegawaiId)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pegawai dengan Nip %s tidak ditemukan", request.PegawaiId)
	}

	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, request.IdPohon)
	if err != nil {
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	if pohon.Id == 0 {
		log.Printf("Pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
	}

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	year := time.Now().Year()
	customId := fmt.Sprintf("REKIN-PEG-%v-%v", year, randomDigits)

	rencanaKinerja := domain.RencanaKinerja{
		Id:                   customId,
		IdPohon:              request.IdPohon,
		SasaranOpdId:         helper.EmptyIntIfNull(request.SasaranOpdId),
		NamaRencanaKinerja:   request.NamaRencanaKinerja,
		Tahun:                request.Tahun,
		KodeSubKegiatan:      "",
		StatusRencanaKinerja: request.StatusRencanaKinerja,
		Catatan:              request.Catatan,
		KodeOpd:              request.KodeOpd,
		PegawaiId:            pegawais.Nip,
		PeriodeId:            request.PeriodeId,
		TahunAwal:            "",
		TahunAkhir:           "",
		JenisPeriode:         "",
		Indikator:            make([]domain.Indikator, len(request.Indikator)),
	}

	log.Printf("RencanaKinerja dibuat dengan ID: %s", customId)

	for i, indikatorRequest := range request.Indikator {
		indikatorRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		indikatorId := fmt.Sprintf("IND-REKIN-%s", indikatorRandomDigits)
		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.NamaIndikator,
			Tahun:            "",
			Target:           make([]domain.Target, len(indikatorRequest.Target)),
			RencanaKinerjaId: rencanaKinerja.Id,
		}

		if indikator.Indikator == "" {
			log.Printf("Indikator kosong ditemukan: %+v", indikator)
		}

		log.Printf("Indikator dibuat: %+v", indikator)

		rencanaKinerja.Indikator[i] = indikator

		randomInt := rand.Intn(100000)

		manualIK := domain.ManualIK{
			Id:                  randomInt,
			IndikatorId:         indikatorId,
			Formula:             helper.EmptyStringIfNull(indikatorRequest.Formula),
			SumberData:          helper.EmptyStringIfNull(indikatorRequest.SumberData),
			Perspektif:          "",
			TujuanRekin:         "",
			Definisi:            "",
			KeyActivities:       "",
			JenisIndikator:      "",
			Kinerja:             false,
			Penduduk:            false,
			Spatial:             false,
			UnitPenanggungJawab: "",
			UnitPenyediaData:    "",
			JangkaWaktuAwal:     "",
			JangkaWaktuAkhir:    "",
			PeriodePelaporan:    "",
		}
		_, err := service.manualIKRepository.Create(ctx, tx, manualIK)
		if err != nil {
			log.Printf("Gagal membuat Manual IK: %v", err)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal membuat Manual IK: %v", err)
		}

		for j, targetRequest := range indikatorRequest.Target {
			targetRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			targetId := fmt.Sprintf("TRGT-IND-REKIN-%s", targetRandomDigits)
			target := domain.Target{
				Id:          targetId,
				Tahun:       targetRequest.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
			log.Printf("Target dibuat dengan ID: %s", targetId)
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Create")
	rencanaKinerja, err = service.rencanaKinerjaRepository.CreateRekinLevel1(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal menyimpan RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menyimpan RencanaKinerja: %v", err)
	}

	rencanaKinerja.NamaOpd = opd.NamaOpd
	rencanaKinerja.NamaPegawai = pegawais.NamaPegawai
	rencanaKinerja.NamaPohon = pohon.NamaPohon
	log.Println("RencanaKinerja berhasil disimpan")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) UpdateRekinLevel1(ctx context.Context, request rencanakinerja.RencanaKinerjaUpdateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Update RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Kode OPD %s tidak ditemukan", request.KodeOpd)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak ditemukan", request.KodeOpd)
		}
		log.Printf("Gagal memeriksa kode OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memeriksa kode OPD: %v", err)
	}

	if opd.KodeOpd == "" {
		log.Printf("Kode OPD %s tidak valid", request.KodeOpd)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak valid", request.KodeOpd)
	}

	// Validasi Pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, request.PegawaiId)
	if err != nil {
		log.Printf("Gagal mengambil data pegawai: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	if pegawai.Id == "" {
		log.Printf("Pegawai dengan NIP %s tidak ditemukan", request.PegawaiId)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pegawai dengan NIP %s tidak ditemukan", request.PegawaiId)
	}

	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, request.IdPohon)
	if err != nil {
		log.Printf("Gagal mengambil data pohon kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	if pohon.Id == 0 {
		log.Printf("Pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("pohon kinerja dengan ID %v tidak ditemukan", request.IdPohon)
	}

	var rencanaKinerja domain.RencanaKinerja

	if request.Id != "" {
		rencanaKinerja, err = service.rencanaKinerjaRepository.FindById(ctx, tx, request.Id, "", "")
		if err != nil {
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan RencanaKinerja: %v", err)
		}

		// Ambil semua indikator lama
		indikatorLamaList, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencanaKinerja.Id)
		if err != nil {
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil indikator lama: %v", err)
		}

		// Buat map untuk indikator baru
		mapIndikatorBaru := make(map[string]bool)
		for _, indReq := range request.Indikator {
			if indReq.Id != "" {
				mapIndikatorBaru[indReq.Id] = true
			}
		}

		// Hapus manual_ik untuk indikator yang tidak ada di request baru
		for _, indLama := range indikatorLamaList {
			if !mapIndikatorBaru[indLama.Id] {
				err := service.manualIKRepository.DeleteByIndikatorId(ctx, tx, indLama.Id)
				if err != nil {
					return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menghapus Manual IK: %v", err)
				}
			}
		}
	}

	rencanaKinerja.IdPohon = request.IdPohon
	rencanaKinerja.SasaranOpdId = helper.EmptyIntIfNull(request.SasaranOpdId)
	rencanaKinerja.NamaRencanaKinerja = request.NamaRencanaKinerja
	rencanaKinerja.Tahun = request.Tahun
	rencanaKinerja.StatusRencanaKinerja = request.StatusRencanaKinerja
	rencanaKinerja.Catatan = request.Catatan
	rencanaKinerja.KodeOpd = request.KodeOpd
	rencanaKinerja.PegawaiId = request.PegawaiId
	rencanaKinerja.PeriodeId = request.PeriodeId
	rencanaKinerja.TahunAwal = ""
	rencanaKinerja.TahunAkhir = ""
	rencanaKinerja.JenisPeriode = ""

	rencanaKinerja.Indikator = make([]domain.Indikator, len(request.Indikator))
	for i, indikatorRequest := range request.Indikator {
		var indikatorId string
		if indikatorRequest.Id != "" {
			indikatorId = indikatorRequest.Id
			// Cek manual_ik yang ada
			existingManualIK, err := service.manualIKRepository.FindByIndikatorId(ctx, tx, indikatorId)
			if err == nil && existingManualIK.Id != 0 {
				// Update manual_ik yang ada (hanya formula dan sumber data)
				manualIK := domain.ManualIK{
					Id:                  existingManualIK.Id,
					IndikatorId:         indikatorId,
					Formula:             helper.EmptyStringIfNull(indikatorRequest.Formula),
					SumberData:          helper.EmptyStringIfNull(indikatorRequest.SumberData),
					Perspektif:          existingManualIK.Perspektif,
					TujuanRekin:         existingManualIK.TujuanRekin,
					Definisi:            existingManualIK.Definisi,
					KeyActivities:       existingManualIK.KeyActivities,
					JenisIndikator:      existingManualIK.JenisIndikator,
					Kinerja:             existingManualIK.Kinerja,
					Penduduk:            existingManualIK.Penduduk,
					Spatial:             existingManualIK.Spatial,
					UnitPenanggungJawab: existingManualIK.UnitPenanggungJawab,
					UnitPenyediaData:    existingManualIK.UnitPenyediaData,
					JangkaWaktuAwal:     existingManualIK.JangkaWaktuAwal,
					JangkaWaktuAkhir:    existingManualIK.JangkaWaktuAkhir,
					PeriodePelaporan:    existingManualIK.PeriodePelaporan,
				}
				_, err := service.manualIKRepository.Update(ctx, tx, manualIK)
				if err != nil {
					return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal update Manual IK: %v", err)
				}
			} else {
				// Buat manual_ik baru jika belum ada
				randomDigits := rand.Intn(100000)
				manualIK := domain.ManualIK{
					Id:                  randomDigits,
					IndikatorId:         indikatorId,
					Formula:             helper.EmptyStringIfNull(indikatorRequest.Formula),
					SumberData:          helper.EmptyStringIfNull(indikatorRequest.SumberData),
					Perspektif:          "",
					TujuanRekin:         "",
					Definisi:            "",
					KeyActivities:       "",
					JenisIndikator:      "",
					Kinerja:             false,
					Penduduk:            false,
					Spatial:             false,
					UnitPenanggungJawab: "",
					UnitPenyediaData:    "",
					JangkaWaktuAwal:     "",
					JangkaWaktuAkhir:    "",
					PeriodePelaporan:    "",
				}
				_, err := service.manualIKRepository.Create(ctx, tx, manualIK)
				if err != nil {
					return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal membuat Manual IK: %v", err)
				}
			}
		} else {
			// Jika indikator baru, buat manual_ik baru
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-REKIN-%s", randomDigits)
			randomInt := rand.Intn(100000)
			manualIK := domain.ManualIK{
				Id:                  randomInt,
				IndikatorId:         indikatorId,
				Formula:             helper.EmptyStringIfNull(indikatorRequest.Formula),
				SumberData:          helper.EmptyStringIfNull(indikatorRequest.SumberData),
				Perspektif:          "",
				TujuanRekin:         "",
				Definisi:            "",
				KeyActivities:       "",
				JenisIndikator:      "",
				Kinerja:             false,
				Penduduk:            false,
				Spatial:             false,
				UnitPenanggungJawab: "",
				UnitPenyediaData:    "",
				JangkaWaktuAwal:     "",
				JangkaWaktuAkhir:    "",
				PeriodePelaporan:    "",
			}
			_, err := service.manualIKRepository.Create(ctx, tx, manualIK)
			if err != nil {
				return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal membuat Manual IK: %v", err)
			}
		}

		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.Indikator,
			RencanaKinerjaId: rencanaKinerja.Id,
			Tahun:            "",
		}

		indikator.Target = make([]domain.Target, len(indikatorRequest.Target))
		for j, targetRequest := range indikatorRequest.Target {
			var targetId string
			if targetRequest.Id != "" {
				targetId = targetRequest.Id
			} else {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRGT-IND-REKIN-%s", randomDigits)
				log.Printf("Membuat Target baru dengan ID: %s", targetId)
			}

			target := domain.Target{
				Id:          targetId,
				Tahun:       targetRequest.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Update")
	rencanaKinerja, err = service.rencanaKinerjaRepository.UpdateRekinLevel1(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal memperbarui RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memperbarui RencanaKinerja: %v", err)
	}

	rencanaKinerja.NamaOpd = opd.NamaOpd
	rencanaKinerja.NamaPegawai = pegawai.NamaPegawai
	rencanaKinerja.NamaPohon = pohon.NamaPohon

	log.Println("RencanaKinerja berhasil diperbarui")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) FindIdRekinLevel1(ctx context.Context, id string) (rencanakinerja.RencanaKinerjaLevel1Response, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data rencana kinerja
	rencanaKinerja, err := service.rencanaKinerjaRepository.FindIdRekinLevel1(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("rencana kinerja dengan ID %s tidak ditemukan", id)
		}
		return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal mengambil data rencana kinerja: %v", err)
	}

	// Ambil nama sasaran OPD
	if rencanaKinerja.SasaranOpdId != 0 {
		sasaranOpd, err := service.SasaranOpdRepository.FindByIdSasaran(ctx, tx, rencanaKinerja.SasaranOpdId)
		if err != nil {
			return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal mengambil data sasaran OPD: %v", err)
		}
		rencanaKinerja.NamaSasaranOpd = sasaranOpd.NamaSasaranOpd
		rencanaKinerja.TahunAwal = sasaranOpd.TahunAwal
		rencanaKinerja.TahunAkhir = sasaranOpd.TahunAkhir
		rencanaKinerja.JenisPeriode = sasaranOpd.JenisPeriode
	}

	// Ambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerja.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			rencanaKinerja.NamaOpd = "OPD tidak ditemukan"
		} else {
			return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
		}
	} else {
		rencanaKinerja.NamaOpd = opd.NamaOpd
	}

	// Ambil data pegawai
	if rencanaKinerja.PegawaiId != "" {
		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencanaKinerja.PegawaiId)
		if err != nil {
			if err == sql.ErrNoRows {
				rencanaKinerja.NamaPegawai = "Pegawai tidak ditemukan"
			} else {
				return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
			}
		} else {
			rencanaKinerja.NamaPegawai = pegawai.NamaPegawai
		}
	} else {
		rencanaKinerja.NamaPegawai = "Pegawai tidak ditentukan"
	}

	// Ambil data pohon kinerja
	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencanaKinerja.IdPohon)
	if err != nil {
		if err == sql.ErrNoRows {
			rencanaKinerja.NamaPohon = "Pohon Kinerja tidak ditemukan"
		} else {
			return rencanakinerja.RencanaKinerjaLevel1Response{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}
	} else {
		rencanaKinerja.NamaPohon = pohon.NamaPohon
	}

	// Konversi ke response
	var indikatorResponses []rencanakinerja.IndikatorResponseLevel1
	for _, indikator := range rencanaKinerja.Indikator {
		var targetResponses []rencanakinerja.TargetResponse
		for _, target := range indikator.Target {
			targetResponse := rencanakinerja.TargetResponse{
				Id:              target.Id,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
				Tahun:           target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := rencanakinerja.IndikatorResponseLevel1{
			Id:            indikator.Id,
			NamaIndikator: indikator.Indikator,
			Formula:       indikator.RumusPerhitungan.String,
			SumberData:    indikator.SumberData.String,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := rencanakinerja.RencanaKinerjaLevel1Response{
		Id:                   rencanaKinerja.Id,
		IdPohon:              rencanaKinerja.IdPohon,
		SasaranOpdId:         rencanaKinerja.SasaranOpdId,
		NamaSasaranOpd:       rencanaKinerja.NamaSasaranOpd,
		TahunAwal:            rencanaKinerja.TahunAwal,
		TahunAkhir:           rencanaKinerja.TahunAkhir,
		JenisPeriode:         rencanaKinerja.JenisPeriode,
		NamaRencanaKinerja:   rencanaKinerja.NamaRencanaKinerja,
		Tahun:                rencanaKinerja.Tahun,
		StatusRencanaKinerja: rencanaKinerja.StatusRencanaKinerja,
		Catatan:              rencanaKinerja.Catatan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: rencanaKinerja.KodeOpd,
			NamaOpd: rencanaKinerja.NamaOpd,
		},
		PegawaiId:   rencanaKinerja.PegawaiId,
		NamaPegawai: rencanaKinerja.NamaPegawai,
		NamaPohon:   rencanaKinerja.NamaPohon,
		Indikator:   indikatorResponses,
	}
	return response, nil
}

func (service *RencanaKinerjaServiceImpl) FindRekinLevel3(ctx context.Context, kodeOpd string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses FindRekinLevel3")

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindRekinLevel3(ctx, tx, kodeOpd, tahun)
	if err != nil {
		log.Printf("Gagal mencari RencanaKinerja Level 3: %v", err)
		return nil, fmt.Errorf("gagal mencari RencanaKinerja Level 3: %v", err)
	}

	var responses []rencanakinerja.RencanaKinerjaResponse
	for _, rencana := range rencanaKinerjaList {
		// Ambil data indikator
		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencana.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Gagal mencari Indikator: %v", err)
			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
		}

		var indikatorResponses []rencanakinerja.IndikatorResponse
		for _, indikator := range indikators {
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Gagal mencari Target: %v", err)
				return nil, fmt.Errorf("gagal mencari Target: %v", err)
			}

			var targetResponses []rencanakinerja.TargetResponse
			for _, target := range targets {
				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			exist, err := service.manualIKRepository.IsIndikatorExist(ctx, tx, indikator.Id)
			if err != nil {
				log.Printf("Gagal memeriksa keberadaan indikator: %v", err)
				return nil, fmt.Errorf("gagal memeriksa keberadaan indikator: %v", err)
			}

			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
				ManualIKExist:    exist,
			})
		}

		// Ambil data OPD
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencana.KodeOpd)
		if err != nil {
			log.Printf("Gagal mencari OPD: %v", err)
			return nil, fmt.Errorf("gagal mencari OPD: %v", err)
		}

		// Ambil data pegawai
		pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencana.PegawaiId)
		if err != nil {
			log.Printf("Gagal mencari Pegawai: %v", err)
			return nil, fmt.Errorf("gagal mencari Pegawai: %v", err)
		}

		// Ambil data pohon kinerja
		pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencana.IdPohon)
		if err != nil {
			log.Printf("Gagal mencari Pohon Kinerja: %v", err)
			return nil, fmt.Errorf("gagal mencari Pohon Kinerja: %v", err)
		}

		responses = append(responses, rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencana.Id,
			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
			Tahun:                rencana.Tahun,
			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
			Catatan:              rencana.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			PegawaiId:   rencana.PegawaiId,
			NamaPegawai: pegawai.NamaPegawai,
			IdPohon:     rencana.IdPohon,
			NamaPohon:   pohon.NamaPohon,

			Indikator: indikatorResponses,
		})
	}

	return responses, nil
}

// rekinatasan
func (service *RencanaKinerjaServiceImpl) FindRekinAtasan(ctx context.Context, rekinId string) (rencanakinerja.RekinAtasanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanakinerja.RekinAtasanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Validasi ID rencana kinerja
	err = service.rencanaKinerjaRepository.ValidateRekinId(ctx, tx, rekinId)
	if err != nil {
		return rencanakinerja.RekinAtasanResponse{}, err
	}

	// 2. Cari pohon kinerja dari rencana kinerja ini
	pokinChild, err := service.CascadingOpdService.cascadingOpdRepository.FindPokinByRekinId(ctx, tx, rekinId)
	if err != nil {
		log.Printf("Error: Pohon kinerja not found for rekin_id=%s: %v", rekinId, err)
		return rencanakinerja.RekinAtasanResponse{}, errors.New("pohon kinerja tidak ditemukan")
	}

	log.Printf("Found pohon kinerja: ID=%d, Level=%d, Parent=%d", pokinChild.Id, pokinChild.LevelPohon, pokinChild.Parent)

	// 3. Cari parent pohon kinerja
	if pokinChild.Parent == 0 {
		log.Printf("Pohon kinerja tidak memiliki parent (sudah di root)")
		return rencanakinerja.RekinAtasanResponse{}, errors.New("pohon kinerja tidak memiliki parent")
	}

	pokinParent, err := service.CascadingOpdService.cascadingOpdRepository.FindPokinById(ctx, tx, pokinChild.Parent)
	if err != nil {
		log.Printf("Error: Parent pohon kinerja not found: %v", err)
		return rencanakinerja.RekinAtasanResponse{}, errors.New("parent pohon kinerja tidak ditemukan")
	}

	log.Printf("Found parent pohon kinerja: ID=%d, Level=%d, Nama=%s", pokinParent.Id, pokinParent.LevelPohon, pokinParent.NamaPohon)

	// 4. Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokinParent.KodeOpd)
	if err != nil {
		log.Printf("Error: OPD not found for kode_opd=%s: %v", pokinParent.KodeOpd, err)
		return rencanakinerja.RekinAtasanResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// 5. Hitung pagu anggaran parent dari semua children-nya
	var paguAnggaran int64

	if pokinParent.LevelPohon == 5 {
		// Tactical: sum dari semua operational children
		paguAnggaran, err = service.CascadingOpdService.calculateAnggaranForTacticalWithPelaksana(ctx, tx, pokinParent.Id)
	} else if pokinParent.LevelPohon == 4 {
		// Strategic: sum dari semua tactical children
		paguAnggaran, err = service.CascadingOpdService.calculateAnggaranForStrategicWithPelaksana(ctx, tx, pokinParent.Id)
	} else if pokinParent.LevelPohon == 6 {
		// Operational: sum dari rencana kinerja di pohon ini
		paguAnggaran, err = service.CascadingOpdService.calculateAnggaranForOperationalWithPelaksana(ctx, tx, pokinParent.Id)
	} else {
		paguAnggaran = 0
	}

	if err != nil {
		log.Printf("Warning: Failed to calculate pagu anggaran: %v", err)
		paguAnggaran = 0
	}

	// 6. Ambil rencana kinerja atasan (pelaksana di pohon parent)
	var rekinAtasanList []rencanakinerja.RekinAtasanDetail
	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokinParent.Id)
	if err == nil {
		// Ambil pelaksana parent
		pelaksanaList, err := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokinParent.Id))
		if err == nil {
			pelaksanaMap := make(map[string]*domainmaster.Pegawai)
			for _, pelaksana := range pelaksanaList {
				pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
				if err == nil {
					pelaksanaMap[pegawai.Nip] = &pegawai
				}
			}

			// ✅ TAMBAHKAN MAP UNTUK DEDUPLICATION
			rekinMap := make(map[string]bool)

			// Filter rencana kinerja hanya yang pelaksana valid
			for _, rk := range rekinList {
				// ✅ SKIP JIKA SUDAH ADA (DEDUPLICATION)
				if rekinMap[rk.Id] {
					continue
				}

				if pegawai, exists := pelaksanaMap[rk.PegawaiId]; exists {
					rekinAtasanList = append(rekinAtasanList, rencanakinerja.RekinAtasanDetail{
						Id:                   rk.Id,
						NamaRencanaKinerja:   rk.NamaRencanaKinerja,
						IdPohon:              rk.IdPohon,
						Tahun:                rk.Tahun,
						StatusRencanaKinerja: rk.StatusRencanaKinerja,
						Catatan:              rk.Catatan,
						KodeOpd:              rk.KodeOpd,
						PegawaiId:            rk.PegawaiId,
						NamaPegawai:          pegawai.NamaPegawai,
					})

					// ✅ MARK SEBAGAI SUDAH DIPROSES
					rekinMap[rk.Id] = true
				}
			}
		}
	}

	// 7. Ambil Program/Kegiatan/Subkegiatan berdasarkan level parent
	var programList []rencanakinerja.ProgramAtasanResponse
	var kegiatanList []rencanakinerja.KegiatanAtasanResponse
	var subkegiatanList []rencanakinerja.SubKegiatanAtasanResponse

	if pokinParent.LevelPohon == 4 || pokinParent.LevelPohon == 5 {
		// Parent level 4 atau 5: Tampilkan program dari children
		programList = service.getProgramForRekinAtasan(ctx, tx, pokinParent)
	} else if pokinParent.LevelPohon == 6 {
		// Parent level 6: Tampilkan kegiatan dan subkegiatan
		kegiatanList, subkegiatanList = service.getKegiatanSubkegiatanForRekinAtasan(ctx, tx, pokinParent)
	}

	// 8. Build response
	response := rencanakinerja.RekinAtasanResponse{
		PokinParent: rencanakinerja.PokinParentInfo{
			Id:         pokinParent.Id,
			NamaPohon:  pokinParent.NamaPohon,
			LevelPohon: pokinParent.LevelPohon,
			KodeOpd:    pokinParent.KodeOpd,
			NamaOpd:    opd.NamaOpd,
		},
		RekinAtasan:       rekinAtasanList,
		Program:           programList,
		Kegiatan:          kegiatanList,
		Subkegiatan:       subkegiatanList,
		PaguAnggaranTotal: paguAnggaran,
	}

	return response, nil
}

// Helper: Get program untuk rekin atasan (level 4 & 5)
func (service *RencanaKinerjaServiceImpl) getProgramForRekinAtasan(
	ctx context.Context,
	tx *sql.Tx,
	pokinParent domain.PohonKinerja) []rencanakinerja.ProgramAtasanResponse {

	var programList []rencanakinerja.ProgramAtasanResponse

	// Ambil operational children
	var operationalIds []int
	var err error

	if pokinParent.LevelPohon == 5 {
		// Tactical: ambil operational children langsung
		operationalIds, err = service.CascadingOpdService.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, pokinParent.Id)
	} else if pokinParent.LevelPohon == 4 {
		// Strategic: ambil tactical children dulu, lalu operational dari setiap tactical
		tacticalIds, err := service.CascadingOpdService.cascadingOpdRepository.FindTacticalChildrenByStrategicId(ctx, tx, pokinParent.Id)
		if err == nil {
			for _, tacticalId := range tacticalIds {
				ops, err := service.CascadingOpdService.cascadingOpdRepository.FindOperationalChildrenByTacticalId(ctx, tx, tacticalId)
				if err == nil {
					operationalIds = append(operationalIds, ops...)
				}
			}
		}
	}

	if err != nil {
		log.Printf("Error getting operational children: %v", err)
		return programList
	}

	// Map untuk program unik
	programMap := make(map[string]*rencanakinerja.ProgramAtasanResponse)

	// Loop operational dan kumpulkan program
	for _, operationalId := range operationalIds {
		rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, operationalId)
		if err == nil {
			// Filter hanya pelaksana valid
			pelaksanaList, _ := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(operationalId))
			pelaksanaMap := make(map[string]bool)
			for _, pelaksana := range pelaksanaList {
				pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
				if err == nil {
					pelaksanaMap[pegawai.Nip] = true
				}
			}

			for _, rk := range rekinList {
				// Skip jika bukan pelaksana
				if !pelaksanaMap[rk.PegawaiId] {
					continue
				}

				if rk.KodeSubKegiatan != "" {
					segments := strings.Split(rk.KodeSubKegiatan, ".")
					if len(segments) >= 3 {
						kodeProgram := strings.Join(segments[:3], ".")

						if _, exists := programMap[kodeProgram]; !exists {
							program, err := service.CascadingOpdService.programRepository.FindByKodeProgram(ctx, tx, kodeProgram)
							if err == nil {
								programMap[kodeProgram] = &rencanakinerja.ProgramAtasanResponse{
									KodeProgram: kodeProgram,
									NamaProgram: program.NamaProgram,
									PaguProgram: 0,
								}
							}
						}

						// Hitung pagu untuk program ini
						if programData, exists := programMap[kodeProgram]; exists {
							var totalAnggaranRenkin int64 = 0
							if rk.Id != "" {
								rencanaAksiList, err := service.RencanaAksiRepository.FindAll(ctx, tx, rk.Id)
								if err == nil {
									for _, ra := range rencanaAksiList {
										rincianBelanja, err := service.CascadingOpdService.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, ra.Id)
										if err == nil {
											totalAnggaranRenkin += rincianBelanja.Anggaran
										}
									}
								}
							}
							programData.PaguProgram += totalAnggaranRenkin
						}
					}
				}
			}
		}
	}

	// Convert map ke slice
	for _, prog := range programMap {
		programList = append(programList, *prog)
	}

	return programList
}

// Helper: Get kegiatan dan subkegiatan untuk rekin atasan (level 6)
func (service *RencanaKinerjaServiceImpl) getKegiatanSubkegiatanForRekinAtasan(
	ctx context.Context,
	tx *sql.Tx,
	pokinParent domain.PohonKinerja) ([]rencanakinerja.KegiatanAtasanResponse, []rencanakinerja.SubKegiatanAtasanResponse) {

	var kegiatanList []rencanakinerja.KegiatanAtasanResponse
	var subkegiatanList []rencanakinerja.SubKegiatanAtasanResponse

	// Ambil rencana kinerja dari pohon parent
	rekinList, err := service.rencanaKinerjaRepository.FindByPokinId(ctx, tx, pokinParent.Id)
	if err != nil {
		log.Printf("Error getting rencana kinerja: %v", err)
		return kegiatanList, subkegiatanList
	}

	// Filter pelaksana valid
	pelaksanaList, _ := service.pohonKinerjaRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokinParent.Id))
	pelaksanaMap := make(map[string]bool)
	for _, pelaksana := range pelaksanaList {
		pegawai, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
		if err == nil {
			pelaksanaMap[pegawai.Nip] = true
		}
	}

	// Map untuk kegiatan dan subkegiatan
	kegiatanMap := make(map[string]*rencanakinerja.KegiatanAtasanResponse)
	subkegiatanMap := make(map[string]*rencanakinerja.SubKegiatanAtasanResponse)

	// Loop rencana kinerja
	for _, rk := range rekinList {
		// Skip jika bukan pelaksana
		if !pelaksanaMap[rk.PegawaiId] {
			continue
		}

		// Hitung anggaran dari rencana kinerja ini
		var totalAnggaranRenkin int64 = 0
		if rk.Id != "" {
			rencanaAksiList, err := service.RencanaAksiRepository.FindAll(ctx, tx, rk.Id)
			if err == nil {
				for _, ra := range rencanaAksiList {
					rincianBelanja, err := service.CascadingOpdService.rincianBelanjaRepository.FindAnggaranByRenaksiId(ctx, tx, ra.Id)
					if err == nil {
						totalAnggaranRenkin += rincianBelanja.Anggaran
					}
				}
			}
		}

		// Kumpulkan kegiatan
		if rk.KodeKegiatan != "" {
			if _, exists := kegiatanMap[rk.KodeKegiatan]; !exists {
				kegiatanMap[rk.KodeKegiatan] = &rencanakinerja.KegiatanAtasanResponse{
					KodeKegiatan: rk.KodeKegiatan,
					NamaKegiatan: rk.NamaKegiatan,
					PaguKegiatan: 0,
				}
			}
			kegiatanMap[rk.KodeKegiatan].PaguKegiatan += totalAnggaranRenkin
		}

		// Kumpulkan subkegiatan
		if rk.KodeSubKegiatan != "" {
			if _, exists := subkegiatanMap[rk.KodeSubKegiatan]; !exists {
				subkegiatanMap[rk.KodeSubKegiatan] = &rencanakinerja.SubKegiatanAtasanResponse{
					KodeSubkegiatan: rk.KodeSubKegiatan,
					NamaSubkegiatan: rk.NamaSubKegiatan,
					PaguSubkegiatan: 0,
				}
			}
			subkegiatanMap[rk.KodeSubKegiatan].PaguSubkegiatan += totalAnggaranRenkin
		}
	}

	// Convert map ke slice
	for _, keg := range kegiatanMap {
		kegiatanList = append(kegiatanList, *keg)
	}
	for _, sub := range subkegiatanMap {
		subkegiatanList = append(subkegiatanList, *sub)
	}

	return kegiatanList, subkegiatanList
}

func (service *RencanaKinerjaServiceImpl) CloneRencanaKinerja(ctx context.Context, rekinId string, tahunBaru string) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Printf("Memulai proses clone rencana kinerja ID: %s ke tahun: %s", rekinId, tahunBaru)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// 1. Validasi rencana kinerja exists
	err = service.rencanaKinerjaRepository.ValidateRekinId(ctx, tx, rekinId)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("rencana kinerja tidak ditemukan: %v", err)
	}

	// 2. Clone rencana kinerja utama
	newRekin, err := service.rencanaKinerjaRepository.CloneRencanaKinerja(ctx, tx, rekinId, tahunBaru)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal clone rencana kinerja: %v", err)
	}

	log.Printf("Rencana kinerja berhasil di-clone dengan ID baru: %s", newRekin.Id)

	// 3. Ambil indikator lama
	indikatorLama, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rekinId)
	if err != nil && err != sql.ErrNoRows {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil indikator lama: %v", err)
	}

	// 4. Clone indikator dan data terkait
	for _, indikator := range indikatorLama {
		// Generate ID untuk indikator baru
		randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		year := time.Now().Year()
		newIndikatorId := fmt.Sprintf("IND-REKIN-%v-%v", year, randomDigits)

		// Clone indikator menggunakan repository
		err = service.rencanaKinerjaRepository.CreateIndikatorClone(ctx, tx, newIndikatorId, newRekin.Id, indikator.Indikator, tahunBaru)
		if err != nil {
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal clone indikator: %v", err)
		}

		log.Printf("Indikator di-clone: %s -> %s", indikator.Id, newIndikatorId)

		// 4a. Clone Target
		err = service.rencanaKinerjaRepository.CloneTarget(ctx, tx, indikator.Id, newIndikatorId, tahunBaru)
		if err != nil {
			log.Printf("Warning: gagal clone target untuk indikator %s: %v", indikator.Id, err)
			// Jangan return error, karena mungkin memang tidak ada target
		}

		// 4b. Clone ManualIK
		err = service.manualIKRepository.CloneManualIK(ctx, tx, indikator.Id, newIndikatorId)
		if err != nil {
			log.Printf("Error clone manual IK untuk indikator %s: %v", indikator.Id, err)
			// Return error karena manual IK penting
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal clone manual IK untuk indikator %s: %v", indikator.Id, err)
		}
	}

	// 5. Clone Rencana Aksi (tanpa pelaksanaan)
	err = service.rencanaKinerjaRepository.CloneRencanaAksi(ctx, tx, rekinId, newRekin.Id)
	if err != nil {
		log.Printf("Warning: gagal clone rencana aksi: %v", err)
		// Tidak return error, karena mungkin memang tidak ada rencana aksi
	}

	// 6. Clone Dasar Hukum
	err = service.rencanaKinerjaRepository.CloneDasarHukum(ctx, tx, rekinId, newRekin.Id)
	if err != nil {
		log.Printf("Error clone dasar hukum: %v", err)
		// Return error karena dasar hukum penting
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal clone dasar hukum: %v", err)
	}

	// 7. Clone Gambaran Umum
	err = service.rencanaKinerjaRepository.CloneGambaranUmum(ctx, tx, rekinId, newRekin.Id)
	if err != nil {
		log.Printf("Warning: gagal clone gambaran umum: %v", err)
		// Tidak return error
	}

	// 8. Clone Inovasi
	err = service.rencanaKinerjaRepository.CloneInovasi(ctx, tx, rekinId, newRekin.Id)
	if err != nil {
		log.Printf("Warning: gagal clone inovasi: %v", err)
		// Tidak return error
	}

	// 9. Clone Permasalahan
	err = service.rencanaKinerjaRepository.ClonePermasalahan(ctx, tx, rekinId, newRekin.Id)
	if err != nil {
		log.Printf("Error clone permasalahan: %v", err)
		// Return error karena permasalahan penting
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal clone permasalahan: %v", err)
	}

	log.Printf("Proses clone selesai untuk rencana kinerja ID: %s", newRekin.Id)

	// 10. Ambil data lengkap rencana kinerja yang baru untuk response
	rencanaKinerjaBaru, err := service.rencanaKinerjaRepository.FindById(ctx, tx, newRekin.Id, newRekin.KodeOpd, tahunBaru)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data rencana kinerja baru: %v", err)
	}

	// Ambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerjaBaru.KodeOpd)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
	}

	// Ambil data pegawai
	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, rencanaKinerjaBaru.PegawaiId)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pegawai: %v", err)
	}

	// Ambil data pohon kinerja
	pohon, err := service.pohonKinerjaRepository.FindById(ctx, tx, rencanaKinerjaBaru.IdPohon)
	if err != nil {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Ambil indikator baru
	indikatorBaru, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, newRekin.Id)
	if err != nil && err != sql.ErrNoRows {
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil indikator baru: %v", err)
	}

	// Proses indikator dan target
	var indikatorResponses []rencanakinerja.IndikatorResponse
	for _, indikator := range indikatorBaru {
		// Ambil manual IK
		manualIK, err := service.manualIKRepository.FindByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			log.Printf("Warning: gagal mengambil manual IK: %v", err)
		}

		// Filter output data yang true saja
		var outputData []string
		if manualIK.Kinerja {
			outputData = append(outputData, "kinerja")
		}
		if manualIK.Penduduk {
			outputData = append(outputData, "penduduk")
		}
		if manualIK.Spatial {
			outputData = append(outputData, "spatial")
		}

		// Ambil target
		targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil && err != sql.ErrNoRows {
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil target: %v", err)
		}

		var targetResponses []rencanakinerja.TargetResponse
		for _, target := range targets {
			targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			})
		}

		indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
			Id:               indikator.Id,
			RencanaKinerjaId: indikator.RencanaKinerjaId,
			NamaIndikator:    indikator.Indikator,
			Target:           targetResponses,
			ManualIK: &rencanakinerja.DataOutput{
				OutputData: outputData,
			},
		})
	}

	// Buat response
	response := rencanakinerja.RencanaKinerjaResponse{
		Id:                   rencanaKinerjaBaru.Id,
		NamaRencanaKinerja:   rencanaKinerjaBaru.NamaRencanaKinerja,
		Tahun:                rencanaKinerjaBaru.Tahun,
		StatusRencanaKinerja: rencanaKinerjaBaru.StatusRencanaKinerja,
		Catatan:              rencanaKinerjaBaru.Catatan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: opd.KodeOpd,
			NamaOpd: opd.NamaOpd,
		},
		PegawaiId:   rencanaKinerjaBaru.PegawaiId,
		NamaPegawai: pegawai.NamaPegawai,
		IdPohon:     rencanaKinerjaBaru.IdPohon,
		NamaPohon:   pohon.NamaPohon,
		Indikator:   indikatorResponses,
	}

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) FindByFilter(ctx context.Context, filter domain.FilterParams) ([]rencanakinerja.RencanaKinerjaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	kodeOPD := filter["kode_opd"]
	tahun := filter["tahun"]

	rekinByFilters, err := service.rencanaKinerjaRepository.FindRekinByFilters(ctx, tx, filter)
	if err != nil {
		log.Printf("Gagal mencari RencanaKinerja: %v", err)
		return nil, fmt.Errorf("gagal mencari RencanaKinerja: %v", err)
	}

	response := []rencanakinerja.RencanaKinerjaResponse{}
	for _, rencana := range rekinByFilters {
		response = append(response, rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencana.Id,
			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
			Tahun:                tahun,
			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
			Catatan:              rencana.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: kodeOPD,
				NamaOpd: rencana.NamaOpd,
			},
			PegawaiId:     rencana.PegawaiId,
			NamaPegawai:   rencana.NamaPegawai,
			IdPohon:       rencana.IdPohon,
			IdParentPohon: rencana.ParentPohon,
			NamaPohon:     rencana.NamaPohon,
			LevelPohon:    rencana.LevelPohon,
			Indikator:     toIndikatorResponses(rencana.Indikator),
		})
	}

	return response, nil
}

func toIndikatorResponses(
	indikators []domain.Indikator,
) []rencanakinerja.IndikatorResponse {

	responses := make([]rencanakinerja.IndikatorResponse, 0, len(indikators))

	for _, indikator := range indikators {
		responses = append(responses, rencanakinerja.IndikatorResponse{
			Id:               indikator.Id,
			RencanaKinerjaId: indikator.RencanaKinerjaId,
			NamaIndikator:    indikator.Indikator,
			Target:           toTargetResponses(indikator.Target),
		})
	}

	return responses
}

func toTargetResponses(targets []domain.Target) []rencanakinerja.TargetResponse {
	responses := make([]rencanakinerja.TargetResponse, 0, len(targets))

	for _, t := range targets {
		responses = append(responses, rencanakinerja.TargetResponse{
			Id:              t.Id,
			IndikatorId:     t.IndikatorId,
			TargetIndikator: t.Target,
			SatuanIndikator: t.Satuan,
			Tahun:           t.Tahun,
		})
	}

	return responses
}
