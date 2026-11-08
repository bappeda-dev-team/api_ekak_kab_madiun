package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/laporancrosscutting"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CrosscuttingOpdServiceImpl struct {
	CrosscuttingOpdRepository repository.CrosscuttingOpdRepository
	PohonKinerjaRepository    repository.PohonKinerjaRepository
	PegawaiRepository         repository.PegawaiRepository
	OpdRepository             repository.OpdRepository
	RekinRepository           repository.RencanaKinerjaRepository
	PkRepository              repository.PkRepository
	DB                        *sql.DB
}

func NewCrosscuttingOpdServiceImpl(crosscuttingOpdRepository repository.CrosscuttingOpdRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, pegawaiRepository repository.PegawaiRepository, opdRepository repository.OpdRepository, rekinRepository repository.RencanaKinerjaRepository, pkRepository repository.PkRepository, DB *sql.DB) *CrosscuttingOpdServiceImpl {
	return &CrosscuttingOpdServiceImpl{
		CrosscuttingOpdRepository: crosscuttingOpdRepository,
		PohonKinerjaRepository:    pohonKinerjaRepository,
		PegawaiRepository:         pegawaiRepository,
		OpdRepository:             opdRepository,
		RekinRepository:           rekinRepository,
		PkRepository:              pkRepository,
		DB:                        DB,
	}
}

func (service *CrosscuttingOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.CrosscuttingOpdCreateRequest, parentId int) (pohonkinerja.CrosscuttingDikirimResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingDikirimResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Konversi request ke domain
	pokin := domain.PohonKinerja{
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Status:     "crosscutting_menunggu",
	}

	for _, indikatorReq := range request.Indikator {
		uuid := uuid.New().String()[:6]
		indikatorId := "IND-CRSS-" + uuid

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: indikatorReq.NamaIndikator,
			Tahun:     request.Tahun,
		}

		// Konversi target dengan generate ID
		for _, targetReq := range indikatorReq.Target {

			targetId := "TRG-CRSS-" + uuid

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId, // Menggunakan ID indikator yang baru digenerate
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       request.Tahun,
			}
			indikator.Target = append(indikator.Target, target)
		}
		pokin.Indikator = append(pokin.Indikator, indikator)
	}

	result, err := service.CrosscuttingOpdRepository.CreateCrosscutting(ctx, tx, pokin, parentId)
	if err != nil {
		return pohonkinerja.CrosscuttingDikirimResponse{}, err
	}
	namaOpdTujuan := ""
	if opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, result.KodeOpd); err == nil {
		namaOpdTujuan = opd.NamaOpd
	}
	response := pohonkinerja.CrosscuttingDikirimResponse{
		IdCrosscutting:         result.Id,
		KeteranganCrosscutting: result.Keterangan,
		NamaPohonTujuan:        request.NamaPohon,
		KodeOpdTujuan:          result.KodeOpd,
		NamaOpdTujuan:          namaOpdTujuan,
		Status:                 result.Status,
	}
	return response, nil
}

func (service *CrosscuttingOpdServiceImpl) Update(ctx context.Context, request pohonkinerja.CrosscuttingOpdUpdateRequest) (pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Konversi request ke domain
	pokin := domain.PohonKinerja{
		Id:         request.Id,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
	}

	// Update data
	result, err := service.CrosscuttingOpdRepository.UpdateCrosscutting(ctx, tx, pokin)
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}

	// Konversi ke response
	response := pohonkinerja.CrosscuttingOpdResponse{
		Id:         result.Id,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Status:     result.Status,
	}

	return response, nil
}

func (service *CrosscuttingOpdServiceImpl) FindAllByParent(ctx context.Context, parentId int) ([]pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Get all pohon kinerja
	pokins, err := service.CrosscuttingOpdRepository.FindAllCrosscutting(ctx, tx, parentId)
	if err != nil {
		return nil, err
	}

	// Collect all pokin IDs
	pokinIds := make([]int, len(pokins))
	for i, pokin := range pokins {
		pokinIds[i] = pokin.Id
	}

	// Get all indikator
	indikators, err := service.CrosscuttingOpdRepository.FindIndikatorByPokinId(ctx, tx, pokinIds)
	if err != nil {
		return nil, err
	}

	// Collect all indikator IDs
	indikatorIds := make([]string, len(indikators))
	for i, ind := range indikators {
		indikatorIds[i] = ind.Id
	}

	// Get all targets
	targets, err := service.CrosscuttingOpdRepository.FindTargetByIndikatorIds(ctx, tx, indikatorIds)
	if err != nil {
		return nil, err
	}

	// Create maps for easy lookup
	indikatorMap := make(map[string][]domain.Indikator) // Ubah ke string sebagai key
	for _, ind := range indikators {
		indikatorMap[ind.PokinId] = append(indikatorMap[ind.PokinId], ind)
	}

	targetMap := make(map[string][]domain.Target)
	for _, target := range targets {
		targetMap[target.IndikatorId] = append(targetMap[target.IndikatorId], target)
	}

	// Build response
	var responses []pohonkinerja.CrosscuttingOpdResponse
	for _, pokin := range pokins {
		opdRepo, err := service.OpdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		helper.PanicIfError(err)

		response := pohonkinerja.CrosscuttingOpdResponse{
			IdCrosscutting: pokin.IdCrosscutting,
			Id:             pokin.Id,
			NamaPohon:      pokin.NamaPohon,
			JenisPohon:     pokin.JenisPohon,
			LevelPohon:     pokin.LevelPohon,
			KodeOpd:        pokin.KodeOpd,
			NamaOpd:        opdRepo.NamaOpd,
			Keterangan:     pokin.Keterangan,
			Tahun:          pokin.Tahun,
			Status:         pokin.Status,
			PegawaiAction:  pokin.PegawaiAction,
		}

		// Proses pegawai_action untuk mendapatkan nama pegawai
		if pokin.PegawaiAction != nil {
			pegawaiActionMap, ok := pokin.PegawaiAction.(map[string]interface{})
			if ok {
				if approveBy, exists := pegawaiActionMap["approve_by"].(string); exists && approveBy != "" {
					pegawai, err := service.PegawaiRepository.FindByNip(ctx, tx, approveBy)
					if err == nil {
						pegawaiActionMap["approve_name"] = pegawai.NamaPegawai
					}
				}
				if rejectBy, exists := pegawaiActionMap["reject_by"].(string); exists && rejectBy != "" {
					pegawai, err := service.PegawaiRepository.FindByNip(ctx, tx, rejectBy)
					if err == nil {
						pegawaiActionMap["rejec_name"] = pegawai.NamaPegawai
					}
				}
				response.PegawaiAction = pegawaiActionMap
			}
		}

		// Add indikator
		pokinIdStr := strconv.Itoa(pokin.Id) // Konversi ID pohon kinerja ke string
		for _, indikator := range indikatorMap[pokinIdStr] {
			indikatorResponse := pohonkinerja.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       indikator.PokinId, // Sudah string, tidak perlu konversi
				NamaIndikator: indikator.Indikator,
			}

			// Add target
			for _, target := range targetMap[indikator.Id] {
				targetResponse := pohonkinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				}
				indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
			}
			response.Indikator = append(response.Indikator, indikatorResponse)
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (service *CrosscuttingOpdServiceImpl) ApproveOrReject(ctx context.Context, crosscuttingId int, request pohonkinerja.CrosscuttingApproveRequest) (*pohonkinerja.CrosscuttingApproveResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.Approve {
		if !request.CreateNew && !request.UseExisting {
			return nil, errors.New("harus memilih create new atau use existing untuk approval")
		}
		if request.CreateNew && request.UseExisting {
			return nil, errors.New("tidak bisa memilih create new dan use existing sekaligus")
		}
		if request.UseExisting && request.ExistingId == 0 {
			return nil, errors.New("existing_id harus diisi jika menggunakan pohon kinerja yang sudah ada")
		}
	}

	err = service.CrosscuttingOpdRepository.ApproveOrRejectCrosscutting(ctx, tx, crosscuttingId, request)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	response := &pohonkinerja.CrosscuttingApproveResponse{
		Id:      crosscuttingId,
		Message: "Crosscutting berhasil diproses",
	}

	if request.Approve {
		if request.CreateNew {
			response.Status = "crosscutting_disetujui"
		} else {
			response.Status = "crosscutting_disetujui_existing"
		}
		response.ApprovedBy = &request.NipPegawai
		response.ApprovedAt = &currentTime
	} else {
		response.Status = "crosscutting_ditolak"
		response.RejectedBy = &request.NipPegawai
		response.RejectedAt = &currentTime
	}

	return response, nil
}

func (service *CrosscuttingOpdServiceImpl) Delete(ctx context.Context, crosscuttingId int, nipPegawai string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)
	return service.CrosscuttingOpdRepository.DeleteCrosscutting(ctx, tx, crosscuttingId, nipPegawai)
}

func (service *CrosscuttingOpdServiceImpl) DeleteUnused(ctx context.Context, crosscuttingId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.CrosscuttingOpdRepository.DeleteUnused(ctx, tx, crosscuttingId)
	if err != nil {
		return err
	}

	return nil
}

func (service *CrosscuttingOpdServiceImpl) FindPokinByCrosscuttingStatus(ctx context.Context, kodeOpd string, tahun string) ([]pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	crosscuttings, err := service.CrosscuttingOpdRepository.FindPokinByCrosscuttingStatus(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	var responses []pohonkinerja.CrosscuttingOpdResponse
	for _, crosscutting := range crosscuttings {
		// Dapatkan data OPD
		opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, crosscutting.KodeOpd)
		if err != nil {
			return nil, err
		}

		var namaOpdPengirim string
		if crosscutting.OpdPengirim != "" {
			opdPengirim, err := service.OpdRepository.FindByKodeOpd(ctx, tx, crosscutting.OpdPengirim)
			if err != nil {
				return nil, err
			}
			namaOpdPengirim = opdPengirim.NamaOpd
		}

		response := pohonkinerja.CrosscuttingOpdResponse{
			Id:              crosscutting.Id,
			Keterangan:      crosscutting.Keterangan,
			KodeOpd:         crosscutting.KodeOpd,
			NamaOpd:         opd.NamaOpd,
			Tahun:           crosscutting.Tahun,
			Status:          crosscutting.Status,
			NamaOpdPengirim: namaOpdPengirim,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (service *CrosscuttingOpdServiceImpl) FindOPDCrosscuttingFrom(ctx context.Context, crosscuttingTo int) (pohonkinerja.CrosscuttingFromResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingFromResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	kodeOpd, err := service.CrosscuttingOpdRepository.FindOPDCrosscuttingFrom(ctx, tx, crosscuttingTo)
	if err != nil {
		return pohonkinerja.CrosscuttingFromResponse{}, err
	}

	response := pohonkinerja.CrosscuttingFromResponse{
		KodeOpd: kodeOpd,
		NamaOpd: "", // Default kosong
	}

	// Hanya cari nama OPD jika kodeOpd tidak kosong
	if kodeOpd != "" {
		opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
		if err != nil {
			return pohonkinerja.CrosscuttingFromResponse{}, err
		}
		response.NamaOpd = opd.NamaOpd
	}

	return response, nil
}

// Plan A
func (service *CrosscuttingOpdServiceImpl) DeleteCrosscuttingDiterima(ctx context.Context, crosscuttingId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)
	return service.CrosscuttingOpdRepository.DeleteCrosscuttingDiterima(ctx, tx, crosscuttingId)
}

// Plan B
func (service *CrosscuttingOpdServiceImpl) UnlinkCrosscuttingDiterima(
	ctx context.Context, crosscuttingId int,
) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)
	return service.CrosscuttingOpdRepository.UnlinkCrosscuttingDiterima(ctx, tx, crosscuttingId)
}

func (service *CrosscuttingOpdServiceImpl) FindAllLaporanCrosscutting(ctx context.Context, kodeOpd, tahun string) (laporancrosscutting.PohonKinerjaCrosscuttingAllResponse, error) {
	startTime := time.Now()
	serviceName := "CrosscuttingOpdService.FindAllLaporanCrosscutting"

	tx, err := service.DB.Begin()
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// Inisialisasi response dasar
	response := laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{
		KodeOpd:    kodeOpd,
		NamaOpd:    opd.NamaOpd,
		Tahun:      tahun,
		PohonKinerja: make([]laporancrosscutting.PohonKinerjaCrosscuttingItemResponse, 0),
	}
	
	// Ambil data pohon kinerja
	pokins, err := service.PohonKinerjaRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil || len(pokins) == 0 {
		log.Printf("[%s] [END] [%s] totalResponseTime=%v (empty/error)",
			time.Now().Format("2006-01-02 15:04:05.000"), serviceName, time.Since(startTime))
		return response, nil
	}

	// Kumpulkan semua pokin IDs
	pokinIds := make([]int, 0, len(pokins))
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			pokinIds = append(pokinIds, p.Id)
		}
	}
	// pokin -> subkegiatan + pagu -> croscutting -> subkegiatan + pagu 

	// Batch fetch pelaksana
	pelaksanas, _ := service.PohonKinerjaRepository.FindPelaksanaPokinBatch(ctx, tx, pokinIds)
	pelaksanaMap := make(map[int][]laporancrosscutting.PelaksanaOpdResponse)
	for pokinId, pelaksanaList := range pelaksanas {
		pelaksanaResponses := make([]laporancrosscutting.PelaksanaOpdResponse, 0, len(pelaksanaList))
		for _, p := range pelaksanaList {
			pelaksanaResponses = append(pelaksanaResponses, laporancrosscutting.PelaksanaOpdResponse{
				Id:          p.Id,
				PegawaiId:   p.PegawaiId,
				NamaPegawai: p.NamaPegawai,
			})
		}
		pelaksanaMap[pokinId] = pelaksanaResponses
	}

	// Batch fetch rencana kinerja
	rekins, err := service.RekinRepository.FindByPokinIdsArray(ctx, tx, pokinIds)
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}

	rekinMap := make(map[int][]laporancrosscutting.RencanaKinerjaResponse)

	for _, p := range rekins {
		rekinMap[p.IdPohon] = append(rekinMap[p.IdPohon],
			laporancrosscutting.RencanaKinerjaResponse{
				Id:              p.Id,
				PegawaiId:       p.PegawaiId,
				NamaPegawai:     p.NamaPegawai,
				KodeOpd:         p.KodeOpd,
				KodeSubKegiatan: p.KodeSubKegiatan,
				NamaSubKegiatan: p.NamaSubKegiatan,
			},
		)
	}

	// Batch fetch indikator
	indikatorBatch, _ := service.PohonKinerjaRepository.FindIndikatorByPokinIdsBatch(ctx, tx, pokinIds)

	// Kumpulkan indikator IDs untuk batch fetch target
	var allIndikatorIds []string
	for _, indikatorList := range indikatorBatch {
		for _, indikator := range indikatorList {
			allIndikatorIds = append(allIndikatorIds, indikator.Id)
		}
	}

	// Batch fetch target
	targetBatch := make(map[string][]domain.Target)
	if len(allIndikatorIds) > 0 {
		targetBatch, _ = service.PohonKinerjaRepository.FindTargetByIndikatorIdsBatch(ctx, tx, allIndikatorIds)
	}

	// Build indikator map
	indikatorMap := make(map[int][]laporancrosscutting.IndikatorResponse)
	for pokinId, indikatorList := range indikatorBatch {
		indikatorResponses := make([]laporancrosscutting.IndikatorResponse, 0, len(indikatorList))
		for _, indikator := range indikatorList {
			targetList := targetBatch[indikator.Id]
			targetResponses := make([]laporancrosscutting.TargetResponse, 0, len(targetList))
			for _, target := range targetList {
				targetResponses = append(targetResponses, laporancrosscutting.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}
			indikatorResponses = append(indikatorResponses, laporancrosscutting.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       fmt.Sprint(pokinId),
				NamaIndikator: indikator.Indikator,
				Target:        targetResponses,
			})
		}
		indikatorMap[pokinId] = indikatorResponses
	}

	
	// Batch fetch tematik - HAPUS SEMUA LOGGING DEBUG
	//
	tematikCloneFromSet := make(map[int]struct{})
	for _, p := range pokins {
		if p.LevelPohon >= 4 && p.CloneFrom > 0 {
			tematikCloneFromSet[p.CloneFrom] = struct{}{}
		}
	}
	tematikCloneFromIds := make([]int, 0, len(tematikCloneFromSet))
	for id := range tematikCloneFromSet {
		tematikCloneFromIds = append(tematikCloneFromIds, id)
	}
	sort.Ints(tematikCloneFromIds)
	

	// Proses pohon ke map
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			if pohonMap[p.LevelPohon] == nil {
				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
			}
			p.NamaOpd = opd.NamaOpd
			pohonMap[p.LevelPohon][p.Parent] = append(pohonMap[p.LevelPohon][p.Parent], p)
		}
	}

	// Batch fetch crosscutting dari tb_crosscutting (by crosscutting_to = id pokin)
	crosscuttingBatch, _ := service.CrosscuttingOpdRepository.FindCrosscuttingByPokinIdsBatch(ctx, tx, pokinIds)
	// crosscuttingDikirimBatch, _ := service.CrosscuttingOpdRepository.FindCrosscuttingPenerimaByPokinIdsBatch(ctx, tx, pokinIds)
	crosscuttingDikirimBatch, err := service.CrosscuttingOpdRepository.FindCrosscuttingPenerimaByPokinIdsBatch(ctx, tx, pokinIds)

	if err != nil {
		log.Printf("ERROR CROSSCUTTING DIKIRIM: %v", err)
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}

	log.Printf("CROSSCUTTING DIKIRIM: %+v", crosscuttingDikirimBatch)

	// Kumpulkan ID pokin penerima dari crosscutting yang dikirim
	pokinPenerimaSet := make(map[int]struct{})

	for _, list := range crosscuttingDikirimBatch {
		for _, c := range list {
			if c.IdPokin > 0 {
				pokinPenerimaSet[c.IdPokin] = struct{}{}
			}
		}
	}

	pokinPenerimaIds := make([]int, 0, len(pokinPenerimaSet))
	for id := range pokinPenerimaSet {
		pokinPenerimaIds = append(pokinPenerimaIds, id)
	}
	

	// Batch fetch rencana kinerja penerima
	rekinDikirims, err := service.RekinRepository.FindByPokinIdsArray(ctx, tx, pokinPenerimaIds)
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}

	rekinDikirimMap := make(map[int][]laporancrosscutting.RencanaKinerjaResponse)

	for _, p := range rekinDikirims {
		rekinDikirimMap[p.IdPohon] = append(rekinDikirimMap[p.IdPohon],
			laporancrosscutting.RencanaKinerjaResponse{
				Id:              p.Id,
				PegawaiId:       p.PegawaiId,
				NamaPegawai:     p.NamaPegawai,
				KodeOpd:         p.KodeOpd,
				KodeSubKegiatan: p.KodeSubKegiatan,
				NamaSubKegiatan: p.NamaSubKegiatan,
			},
		)
	}

	// =========================
	// BATCH FETCH PAGU
	// =========================

	kodeSubKegiatanSet := make(map[string]struct{})

	// Ambil kode subkegiatan REKIN utama
	for _, rekinList := range rekinMap {
		for _, rekin := range rekinList {
			if rekin.KodeSubKegiatan != "" {
				kodeSubKegiatanSet[rekin.KodeSubKegiatan] = struct{}{}
			}
		}
	}

	// Ambil kode subkegiatan REKIN penerima
	for _, rekinList := range rekinDikirimMap {
		for _, rekin := range rekinList {
			if rekin.KodeSubKegiatan != "" {
				kodeSubKegiatanSet[rekin.KodeSubKegiatan] = struct{}{}
			}
		}
	}

	// Ubah set menjadi slice
	kodeSubKegiatans := make([]string, 0, len(kodeSubKegiatanSet))

	for kode := range kodeSubKegiatanSet {
		kodeSubKegiatans = append(kodeSubKegiatans, kode)
	}

	// Fetch pagu sekaligus
	paguMap, err := service.PkRepository.FindPaguPkByKodeSubkegiatans(ctx, tx, kodeSubKegiatans)
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}

	// Masukkan pagu ke REKIN utama
	for pokinId, rekinList := range rekinMap {
		for i := range rekinList {
			rekinList[i].Pagu = paguMap[rekinList[i].KodeSubKegiatan]
		}

		rekinMap[pokinId] = rekinList
	}

	// Masukkan pagu ke REKIN penerima
	for pokinId, rekinList := range rekinDikirimMap {
		for i := range rekinList {
			rekinList[i].Pagu = paguMap[rekinList[i].KodeSubKegiatan]
		}

		rekinDikirimMap[pokinId] = rekinList
	}

	pelaksanaDikirims, _ := service.PohonKinerjaRepository.FindPelaksanaPokinBatch(ctx, tx, pokinPenerimaIds)
	pelaksanaDikirimMap := make(map[int][]laporancrosscutting.PelaksanaOpdResponse)
	for pokinPenerima, pelaksanaList := range pelaksanaDikirims {
		pelaksanaResponses := make([]laporancrosscutting.PelaksanaOpdResponse, 0, len(pelaksanaList))
		for _, p := range pelaksanaList {
			pelaksanaResponses = append(pelaksanaResponses, laporancrosscutting.PelaksanaOpdResponse{
				Id:          p.Id,
				PegawaiId:   p.PegawaiId,
				NamaPegawai: p.NamaPegawai,
			})
		}
		pelaksanaDikirimMap[pokinPenerima] = pelaksanaResponses
	}

	// Ambil indikator khusus pokin penerima
	indikatorDikirimBatch, err := service.PohonKinerjaRepository.FindIndikatorByPokinIdsBatch(
		ctx, tx, pokinPenerimaIds,
	)
	if err != nil {
		return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
	}

	// Kumpulkan ID indikator penerima untuk mengambil target
	allIndikatorDikirimIds := make([]string, 0)

	for _, indikatorList := range indikatorDikirimBatch {
		for _, indikator := range indikatorList {
			allIndikatorDikirimIds = append(
				allIndikatorDikirimIds,
				indikator.Id,
			)
		}
	}

	// Ambil target khusus indikator penerima
	targetDikirimBatch := make(map[string][]domain.Target)

	if len(allIndikatorDikirimIds) > 0 {
		targetDikirimBatch, err = service.PohonKinerjaRepository.FindTargetByIndikatorIdsBatch(
			ctx, tx, allIndikatorDikirimIds,
		)
		if err != nil {
			return laporancrosscutting.PohonKinerjaCrosscuttingAllResponse{}, err
		}
	}

	// Bentuk map indikator penerima secara terpisah
	indikatorDikirimMap := make(
		map[int][]laporancrosscutting.IndikatorResponse,
	)

	for pokinIdPenerima, indikatorList := range indikatorDikirimBatch {
		indikatorResponses := make(
			[]laporancrosscutting.IndikatorResponse,
			0,
			len(indikatorList),
		)

		for _, indikator := range indikatorList {
			targetList := targetDikirimBatch[indikator.Id]

			targetResponses := make(
				[]laporancrosscutting.TargetResponse,
				0,
				len(targetList),
			)

			for _, target := range targetList {
				targetResponses = append(
					targetResponses,
					laporancrosscutting.TargetResponse{
						Id:              target.Id,
						IndikatorId:     target.IndikatorId,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					},
				)
			}

			indikatorResponses = append(
				indikatorResponses,
				laporancrosscutting.IndikatorResponse{
					Id:            indikator.Id,
					IdPokin:       fmt.Sprint(pokinIdPenerima),
					NamaIndikator: indikator.Indikator,
					Target:        targetResponses,
				},
			)
		}

		indikatorDikirimMap[pokinIdPenerima] = indikatorResponses
	}

	opdTujuanNamaMap := make(map[string]string)
	for _, list := range crosscuttingDikirimBatch {
		for _, c := range list {
			if c.KodeOpd != "" {
				opdTujuanNamaMap[c.KodeOpd] = ""
			}
		}
	}
	for kode := range opdTujuanNamaMap {
		o, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kode)
		if err == nil {
			opdTujuanNamaMap[kode] = o.NamaOpd
		}
	}

	crosscuttingDikirimMap := make(map[int][]laporancrosscutting.CrosscuttingDikirimResponse)
	for pokinId, list := range crosscuttingDikirimBatch {
		items := make([]laporancrosscutting.CrosscuttingDikirimResponse, 0, len(list))
		for _, c := range list {

			indikatorDikirim := indikatorDikirimMap[c.IdPokin]

			if indikatorDikirim == nil {
				indikatorDikirim = []laporancrosscutting.IndikatorResponse{}
			}

			pelaksanaDikirim := pelaksanaDikirimMap[c.IdPokin]

			if pelaksanaDikirim == nil {
				pelaksanaDikirim = []laporancrosscutting.PelaksanaOpdResponse{}
			}

			rekinDikirim := rekinDikirimMap[c.IdPokin]

			if rekinDikirim == nil {
				rekinDikirim = []laporancrosscutting.RencanaKinerjaResponse{}
			}

			items = append(items, laporancrosscutting.CrosscuttingDikirimResponse{
				IdCrosscutting:         c.Id,
				IdPokin:                c.IdPokin,
				KeteranganCrosscutting: c.Keterangan,
				NamaPohonTujuan:        c.NamaPohonAsal,
				KodeOpdTujuan:          c.KodeOpd,
				NamaOpdTujuan:          opdTujuanNamaMap[c.KodeOpd],
				Status:                 c.Status,
				// IdRekin:                c.IdRekin,
				IndikatorDikirim:       indikatorDikirim,
				PelaksanaDikirim:       pelaksanaDikirim,
				RekinDikirim:           rekinDikirim,
			})
		}
		crosscuttingDikirimMap[pokinId] = items
	}

	
	// Ambil kode OPD asal yang unik untuk fetch nama OPD sekali
	opdAsalSet := make(map[string]struct{})
	for _, list := range crosscuttingBatch {
		for _, c := range list {
			if c.OpdPengirim != "" {
				opdAsalSet[c.OpdPengirim] = struct{}{}
			}
		}
	}
	opdAsalNamaMap := make(map[string]string)
	for kode := range opdAsalSet {
		o, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kode)
		if err == nil {
			opdAsalNamaMap[kode] = o.NamaOpd
		}
	}
	// Build crosscutting map siap pakai
	crosscuttingStatusMap := make(map[int]string)
	crosscuttingMap := make(map[int][]laporancrosscutting.CrosscuttingPokinResponse)
	for pokinId, list := range crosscuttingBatch {
		if len(list) > 0 {
			// Satu status per pohon; contoh: ambil dari baris pertama (sesuaikan jika banyak baris berbeda status)
			crosscuttingStatusMap[pokinId] = list[0].Status
		}
		items := make([]laporancrosscutting.CrosscuttingPokinResponse, 0, len(list))
		for _, c := range list {
			items = append(items, laporancrosscutting.CrosscuttingPokinResponse{
				IdCrosscutting:         c.Id,
				KeteranganCrosscutting: c.Keterangan,
				KodeOpdAsal:            c.OpdPengirim,
				NamaPohonAsal:          c.NamaPohonAsal,
				NamaOpdAsal:            opdAsalNamaMap[c.OpdPengirim],
				Status:                 c.Status,
			})
		}
		crosscuttingMap[pokinId] = items
	}
	crosscuttingMap = make(map[int][]laporancrosscutting.CrosscuttingPokinResponse)
	for pokinId, list := range crosscuttingBatch {
		items := make([]laporancrosscutting.CrosscuttingPokinResponse, 0, len(list))
		for _, c := range list {
			items = append(items, laporancrosscutting.CrosscuttingPokinResponse{
				IdCrosscutting:         c.Id,
				KeteranganCrosscutting: c.Keterangan,
				KodeOpdAsal:            c.OpdPengirim,
				NamaOpdAsal:            opdAsalNamaMap[c.OpdPengirim],
				// Status:                 c.Status,
			})
		}
		crosscuttingMap[pokinId] = items
	}

	// Build response flat: hanya pohon yang memiliki crosscutting.
	for _, p := range pokins {
		if p.LevelPohon < 4 {
			continue
		}

		// crosscutting := crosscuttingMap[p.Id]
		crosscuttingDikirim := crosscuttingDikirimMap[p.Id]

		// Hanya tampilkan pohon yang memiliki crosscutting.
		if len(crosscuttingDikirim) == 0 {
			continue
		}

		p.NamaOpd = opd.NamaOpd

		indikator := indikatorMap[p.Id]
		if indikator == nil {
			indikator = []laporancrosscutting.IndikatorResponse{}
		}

		pelaksana := pelaksanaMap[p.Id]
		if pelaksana == nil {
			pelaksana = []laporancrosscutting.PelaksanaOpdResponse{}
		}

		rekin := rekinMap[p.Id]
		if rekin == nil {
			rekin = []laporancrosscutting.RencanaKinerjaResponse{}
		}

		// if crosscutting == nil {
		// 	crosscutting = []laporancrosscutting.CrosscuttingPokinResponse{}
		// }

		if crosscuttingDikirim == nil {
			crosscuttingDikirim = []laporancrosscutting.CrosscuttingDikirimResponse{}
		}

		response.PohonKinerja = append(
			response.PohonKinerja,
			laporancrosscutting.PohonKinerjaCrosscuttingItemResponse{
				PohonKinerja: laporancrosscutting.PohonKinerja{
					Id:          p.Id,
					Parent:      p.Parent,
					NamaPohon:   p.NamaPohon,
					KodeOpd:     p.KodeOpd,
					NamaOpd:     p.NamaOpd,
					Tahun:       p.Tahun,
					JenisPohon:  p.JenisPohon,
					LevelPohon:  p.LevelPohon,
				},
				// Crosscutting:        crosscutting,
				CrosscuttingDikirim: crosscuttingDikirim,
				Indikator:           indikator,
				Pelaksana:           pelaksana,
				RencanaKinerja:      rekin,
			},
		)
	}

	log.Printf(
		"[%s] [END] [%s] totalResponseTime=%v, crosscuttingCount=%d",
		time.Now().Format("2006-01-02 15:04:05.000"),
		serviceName,
		time.Since(startTime),
		len(response.PohonKinerja),
	)

	return response, nil
}
