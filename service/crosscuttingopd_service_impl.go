package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CrosscuttingOpdServiceImpl struct {
	CrosscuttingOpdRepository repository.CrosscuttingOpdRepository
	PohonKinerjaRepository    repository.PohonKinerjaRepository
	PegawaiRepository         repository.PegawaiRepository
	OpdRepository             repository.OpdRepository
	DB                        *sql.DB
}

func NewCrosscuttingOpdServiceImpl(crosscuttingOpdRepository repository.CrosscuttingOpdRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, pegawaiRepository repository.PegawaiRepository, opdRepository repository.OpdRepository, DB *sql.DB) *CrosscuttingOpdServiceImpl {
	return &CrosscuttingOpdServiceImpl{
		CrosscuttingOpdRepository: crosscuttingOpdRepository,
		PohonKinerjaRepository:    pohonKinerjaRepository,
		PegawaiRepository:         pegawaiRepository,
		OpdRepository:             opdRepository,
		DB:                        DB,
	}
}

func (service *CrosscuttingOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.CrosscuttingOpdCreateRequest, parentId int) (pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
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

	// Konversi indikator untuk response
	for _, indikator := range result.Indikator {
		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            indikator.Id,
			NamaIndikator: indikator.Indikator,
		}

		// Konversi target untuk response
		for _, target := range indikator.Target {
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
