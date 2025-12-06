package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/datamaster"
	"ekak_kabupaten_madiun/repository"
)

type DataMasterServiceImpl struct {
	DataMasterRepository      repository.DataMasterRepository
	RencanaKinerjaRepository  repository.RencanaKinerjaRepository
	CrosscuttingOpdRepository repository.CrosscuttingOpdRepository
	DB                        *sql.DB
}

func NewDataMasterServiceImpl(dataMasterRepository repository.DataMasterRepository, rencanaKinerjaRepository repository.RencanaKinerjaRepository, crosscuttingOpdRepository repository.CrosscuttingOpdRepository, DB *sql.DB) *DataMasterServiceImpl {
	return &DataMasterServiceImpl{
		DataMasterRepository:      dataMasterRepository,
		RencanaKinerjaRepository:  rencanaKinerjaRepository,
		CrosscuttingOpdRepository: crosscuttingOpdRepository,
		DB:                        DB,
	}
}

func (service *DataMasterServiceImpl) DataRBByTahun(ctx context.Context, tahunBase int) ([]datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	result, err := service.DataMasterRepository.DataRBByTahun(ctx, tx, tahunBase, nil)
	if err != nil {
		return nil, err
	}

	listIdRB := make([]int, 0, len(result))
	// Convert MasterRB → RBResponse
	responses := make([]datamaster.RBResponse, 0, len(result))
	respMap := make(map[int]*datamaster.RBResponse, len(result))

	for _, rb := range result {
		listIdRB = append(listIdRB, rb.Id)
		// Mapping MasterRB → RBResponse
		resp := datamaster.RBResponse{
			IdRB:          rb.Id,
			JenisRB:       rb.JenisRB,
			KegiatanUtama: rb.KegiatanUtama,
			Keterangan:    rb.Keterangan,
			TahunBaseline: rb.TahunBaseline,
			TahunNext:     rb.TahunNext,
			Indikator:     make([]datamaster.IndikatorRB, 0),
			SudahDiambil:  false,
		}

		// Mapping Indikator
		for _, ind := range rb.Indikator {
			indResp := datamaster.IndikatorRB{
				IdIndikator: ind.IdIndikator,
				IdRB:        ind.IdRB,
				Indikator:   ind.Indikator,
				TargetRB:    make([]datamaster.TargetRB, 0),
			}

			for _, tar := range ind.TargetRB {
				rbResp := datamaster.TargetRB{
					IdTarget:          tar.IdTarget,
					IdIndikator:       tar.IdIndikator,
					TahunBaseline:     tar.TahunBaseline,
					TargetBaseline:    tar.TargetBaseline,
					RealisasiBaseline: tar.RealisasiBaseline,
					SatuanBaseline:    tar.SatuanBaseline,
					TahunNext:         tar.TahunNext,
					TargetNext:        tar.TargetNext,
					SatuanNext:        tar.SatuanNext,
				}
				indResp.TargetRB = append(indResp.TargetRB, rbResp)
			}

			resp.Indikator = append(resp.Indikator, indResp)
		}

		responses = append(responses, resp)
		respMap[rb.Id] = &responses[len(responses)-1]
	}

	pokinRbRes, err := service.DataMasterRepository.PokinByIdRBs(ctx, tx, listIdRB)
	if err != nil {
		return nil, err
	}

	for _, pokin := range pokinRbRes {
		if r, ok := respMap[pokin.KodeRB]; ok {
			r.SudahDiambil = true
		}
	}

	return responses, nil
}

func (service *DataMasterServiceImpl) SaveRB(ctx context.Context, rb datamaster.RBRequest, userId int) (datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return datamaster.RBResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// === Convert RBRequest → MasterRB ===
	entity := datamaster.ConvertRBRequestToMaster(rb, userId)

	// === 1. INSERT MASTER RB ===
	rbID, err := service.DataMasterRepository.InsertRB(ctx, tx, entity, userId)
	if err != nil {
		return datamaster.RBResponse{}, err
	}

	// === Build Response RB ===
	response := datamaster.RBResponse{
		IdRB:          int(rbID),
		JenisRB:       entity.JenisRB,
		KegiatanUtama: entity.KegiatanUtama,
		Keterangan:    entity.Keterangan,
		TahunBaseline: entity.TahunBaseline,
		TahunNext:     entity.TahunNext,
		Indikator:     []datamaster.IndikatorRB{},
	}

	// === 2. INSERT INDIKATOR ===
	for _, ind := range entity.Indikator {

		// INSERT into DB
		indikatorID, err := service.DataMasterRepository.InsertIndikator(ctx, tx, rbID, ind)
		if err != nil {
			return datamaster.RBResponse{}, err
		}

		// siapkan indikator response
		indResp := datamaster.IndikatorRB{
			IdRB:        int(rbID),
			IdIndikator: indikatorID,
			Indikator:   ind.Indikator,
			TargetRB:    []datamaster.TargetRB{},
		}

		// === 3. INSERT TARGET ===
		for _, t := range ind.TargetRB {

			// simpan target
			if err := service.DataMasterRepository.InsertTarget(ctx, tx, indikatorID, t); err != nil {
				return datamaster.RBResponse{}, err
			}

			// tambahkan ke response
			indResp.TargetRB = append(indResp.TargetRB, datamaster.TargetRB{
				IdIndikator:       indikatorID,
				TahunBaseline:     t.TahunBaseline,
				TargetBaseline:    t.TargetBaseline,
				RealisasiBaseline: t.RealisasiBaseline,
				SatuanBaseline:    t.SatuanBaseline,
				TahunNext:         t.TahunNext,
				TargetNext:        t.TargetNext,
				SatuanNext:        t.SatuanNext,
			})
		}

		// masukkan indikator ke response
		response.Indikator = append(response.Indikator, indResp)
	}

	return response, nil
}

func (service *DataMasterServiceImpl) UpdateRB(ctx context.Context, rb datamaster.RBRequest, userId int, rbId int) (datamaster.RBResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return datamaster.RBResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// === 1. Cek apakah RB ada ===
	existingRB, err := service.DataMasterRepository.FindRBById(ctx, tx, rbId)
	if err != nil {
		return datamaster.RBResponse{}, errors.New("rb_not_found")
	}

	// === 2. Convert RBRequest → Entity ===
	entity := datamaster.ConvertRBRequestToMaster(rb, userId)
	entity.Id = existingRB.Id // penting

	// === 3. Update master RB ===
	if err := service.DataMasterRepository.UpdateRB(ctx, tx, entity, rbId); err != nil {
		return datamaster.RBResponse{}, err
	}

	// === 4. Hapus semua indikator & target lama ===
	if err := service.DataMasterRepository.DeleteAllIndikatorAndTargetByRB(ctx, tx, rbId); err != nil {
		return datamaster.RBResponse{}, err
	}

	// === 5. Insert indikator & target baru ===
	response := datamaster.RBResponse{
		IdRB:          rbId,
		JenisRB:       entity.JenisRB,
		KegiatanUtama: entity.KegiatanUtama,
		Keterangan:    entity.Keterangan,
		TahunBaseline: entity.TahunBaseline,
		TahunNext:     entity.TahunNext,
		Indikator:     []datamaster.IndikatorRB{},
	}

	for _, ind := range entity.Indikator {

		// Insert indikator baru
		indikatorID, err := service.DataMasterRepository.InsertIndikator(ctx, tx, int64(rbId), ind)
		if err != nil {
			return datamaster.RBResponse{}, err
		}

		indikatorRes := datamaster.IndikatorRB{
			IdRB:        rbId,
			IdIndikator: indikatorID,
			Indikator:   ind.Indikator,
			TargetRB:    []datamaster.TargetRB{},
		}

		for _, t := range ind.TargetRB {
			if err := service.DataMasterRepository.InsertTarget(ctx, tx, indikatorID, t); err != nil {
				return datamaster.RBResponse{}, err
			}

			indikatorRes.TargetRB = append(indikatorRes.TargetRB, datamaster.TargetRB{
				IdIndikator:       indikatorID,
				TahunBaseline:     t.TahunBaseline,
				TargetBaseline:    t.TargetBaseline,
				RealisasiBaseline: t.RealisasiBaseline,
				SatuanBaseline:    t.SatuanBaseline,
				TahunNext:         t.TahunNext,
				TargetNext:        t.TargetNext,
				SatuanNext:        t.SatuanNext,
			})
		}

		response.Indikator = append(response.Indikator, indikatorRes)
	}

	return response, nil
}

func (s *DataMasterServiceImpl) DeleteRB(ctx context.Context, rbId int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = s.DataMasterRepository.DeleteRB(ctx, tx, rbId)
	if err != nil {
		return err
	}
	return nil
}

func (service *DataMasterServiceImpl) FindByTahun(ctx context.Context, tahunBase int) ([]datamaster.RbResponseTahunan, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	result, err := service.DataMasterRepository.DataRBByTahun(ctx, tx, tahunBase, nil)
	if err != nil {
		return nil, err
	}

	// Convert MasterRB → RBResponse
	responses := make([]datamaster.RbResponseTahunan, 0, len(result))

	for _, rb := range result {
		// Mapping MasterRB → RBResponse
		resp := datamaster.RbResponseTahunan{
			IdRB:          rb.Id,
			JenisRB:       rb.JenisRB,
			KegiatanUtama: rb.KegiatanUtama,
			Keterangan:    rb.Keterangan,
			TahunBaseline: rb.TahunBaseline,
			TahunNext:     rb.TahunNext,
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (service *DataMasterServiceImpl) LaporanByTahun(ctx context.Context, tahunNext int, jenisRB string) ([]datamaster.RbLaporanTahunanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	result, err := service.DataMasterRepository.DataRBByTahun(ctx, tx, tahunNext, &jenisRB)
	if err != nil {
		return nil, err
	}

	// responses slice + map untuk lookup cepat
	responses := make([]datamaster.RbLaporanTahunanResponse, 0, len(result))
	respMap := make(map[int]*datamaster.RbLaporanTahunanResponse, len(result))
	listIdRB := make([]int, 0, len(result))

	for _, rb := range result {
		listIdRB = append(listIdRB, rb.Id)

		resp := datamaster.RbLaporanTahunanResponse{
			IdRB:          rb.Id,
			JenisRB:       rb.JenisRB,
			KegiatanUtama: rb.KegiatanUtama,
			Keterangan:    rb.Keterangan,
			TahunBaseline: rb.TahunBaseline,
			TahunNext:     rb.TahunNext,
			Indikator:     make([]datamaster.IndikatorRB, 0, len(rb.Indikator)),
			RencanaAksis:  make([]datamaster.RencanaAksiRB, 0),
		}

		for _, ind := range rb.Indikator {
			indResp := datamaster.IndikatorRB{
				IdIndikator: ind.IdIndikator,
				IdRB:        ind.IdRB,
				Indikator:   ind.Indikator,
				TargetRB:    make([]datamaster.TargetRB, 0, len(ind.TargetRB)),
			}

			for _, tar := range ind.TargetRB {
				tarResp := datamaster.TargetRB{
					IdTarget:          tar.IdTarget,
					IdIndikator:       tar.IdIndikator,
					TahunBaseline:     tar.TahunBaseline,
					TargetBaseline:    tar.TargetBaseline,
					RealisasiBaseline: tar.RealisasiBaseline,
					SatuanBaseline:    tar.SatuanBaseline,
					TahunNext:         tar.TahunNext,
					TargetNext:        tar.TargetNext,
					SatuanNext:        tar.SatuanNext,
				}
				indResp.TargetRB = append(indResp.TargetRB, tarResp)
			}

			resp.Indikator = append(resp.Indikator, indResp)
		}

		// append ke slice responses, lalu simpan pointer ke map
		responses = append(responses, resp)
		respMap[rb.Id] = &responses[len(responses)-1]
	}

	// find kebutuhan lewat pokin
	pokinRbRes, err := service.DataMasterRepository.PokinByIdRBs(ctx, tx, listIdRB)
	if err != nil {
		return nil, err
	}

	// buat map dari IdPokin -> KodeRB agar saat memproses rekin kita tahu RB targetnya
	pokinToRB := make(map[int]int, len(pokinRbRes)) // IdPokin -> KodeRB
	for _, pokin := range pokinRbRes {
		pokinToRB[pokin.IdPokin] = pokin.KodeRB
	}

	// list pokin id unique untuk query rekin
	listPokinIds := make([]int, 0, len(pokinRbRes))
	seen := make(map[int]bool)

	for _, pokin := range pokinRbRes {
		if !seen[pokin.IdPokin] {
			seen[pokin.IdPokin] = true
			listPokinIds = append(listPokinIds, pokin.IdPokin)
		}
	}

	// get the rekin by id pokin
	rekinRes, err := service.RencanaKinerjaRepository.FindByPokinIds(ctx, tx, listPokinIds, tahunNext)
	if err != nil {
		return nil, err
	}

	// get pokin yang ada di rekin saja
	uniquePohonIdsInRekin := make([]int, 0)
	seenPohon := make(map[int]struct{})
	for _, r := range rekinRes {
		if _, ok := seenPohon[r.IdPohon]; ok {
			continue
		}
		seenPohon[r.IdPohon] = struct{}{}
		uniquePohonIdsInRekin = append(uniquePohonIdsInRekin, r.IdPohon)
	}

	crossRows, err := service.CrosscuttingOpdRepository.FindCrosscuttingByPohonIdsFrom(ctx, tx, uniquePohonIdsInRekin)
	if err != nil {
		return nil, err
	}

	crossMap := make(map[int][]datamaster.OpdCrosscutting)
	for _, c := range crossRows {
		crossMap[c.CrosscuttingFrom] = append(crossMap[c.CrosscuttingFrom], datamaster.OpdCrosscutting{
			IdPohon:       c.CrosscuttingFrom,
			KodeOpd:       c.KodeOpd,
			NamaOpd:       c.OpdPengirim,
			NipPelaksana:  "",
			NamaPelaksana: "",
		})
	}

	// pokinToRekinMap:
	// map[int][]int   => pokinId -> list rekinId
	pokinToRekinSet := make(map[int]map[string]struct{})
	// tracker uniq rekin
	addedRA := make(map[int]map[string]struct{})

	// untuk tiap rekin, buat RencanaAksiRB dan append ke RB yang sesuai
	for _, rekin := range rekinRes {
		// append to pokinToRekinMap
		// buat cari crosscutting
		if _, ok := pokinToRekinSet[rekin.IdPohon]; !ok {
			pokinToRekinSet[rekin.IdPohon] = make(map[string]struct{})
		}
		// skip jika rekin.Id sudah pernah diproses untuk pokin ini (uniqueness)
		if _, seen := pokinToRekinSet[rekin.IdPohon][rekin.Id]; seen {
			continue
		}
		// tandai sebagai sudah diproses
		pokinToRekinSet[rekin.IdPohon][rekin.Id] = struct{}{}

		// bangun response RencanaAksiRB dari rekin
		ra := datamaster.RencanaAksiRB{
			IdRencanaAksi:   rekin.Id,
			RencanaAksi:     rekin.NamaRencanaKinerja,
			IndikatorOutput: make([]datamaster.IndikatorRencanaAksiRB, 0, len(rekin.Indikator)),
			Anggaran:        0,
			Realisasi:       0,
			Capaian:         "0%",
			OpdKoordinator:  rekin.NamaOpd,
			NipPelaksana:    rekin.PegawaiId,
			NamaPelaksana:   rekin.NamaPegawai,
			OpdCrosscutting: crossMap[rekin.IdPohon],
		}

		// indikator - target, biasa
		for _, ind := range rekin.Indikator {
			indResp := datamaster.IndikatorRencanaAksiRB{
				Indikator:       ind.Indikator,
				TargetIndikator: make([]datamaster.TargetIndikatorRencanaAksiRB, 0, len(ind.Target)),
			}

			for _, tar := range ind.Target {
				tarResp := datamaster.TargetIndikatorRencanaAksiRB{
					Target:    tar.Target,
					Realisasi: "0",
					Satuan:    tar.Satuan,
					Capaian:   "0%",
					Tahun:     tar.Tahun,
				}
				indResp.TargetIndikator = append(indResp.TargetIndikator, tarResp)
			}

			ra.IndikatorOutput = append(ra.IndikatorOutput, indResp)
		}

		// cari kodeRB lewat pokinToRB map
		// Asumsi: rekin.PokinId adalah field yang mengacu ke IdPokin
		kodeRB, ok := pokinToRB[rekin.IdPohon]
		if !ok {
			// jika tidak ditemukan, lewati (atau log) — tergantung kebijakan kamu
			continue
		}

		// checker agar tidak duplikat
		if _, ok := addedRA[kodeRB]; !ok {
			addedRA[kodeRB] = make(map[string]struct{})
		}

		if _, exists := addedRA[kodeRB][ra.IdRencanaAksi]; exists {
			// sudah ada, skip
			continue
		}

		// append rencana aksi ke RB yang sesuai (cek kalau RB ada di respMap)
		if resp, exists := respMap[kodeRB]; exists && resp != nil {
			resp.RencanaAksis = append(resp.RencanaAksis, ra)
			// tandai sudah ditambahkan
			addedRA[kodeRB][ra.IdRencanaAksi] = struct{}{}
		} else {
			log.Printf("warning: respMap tidak memiliki entry untuk kodeRB=%d (rekinId=%s)", kodeRB, rekin.Id)
		}

	}

	return responses, nil
}
