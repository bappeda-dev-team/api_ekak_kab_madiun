package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type MatrixRenstraServiceImpl struct {
	MatrixRenstraRepository repository.MatrixRenstraRepository
	PeriodeRepository       repository.PeriodeRepository
	PegawaiRepository       repository.PegawaiRepository
	DB                      *sql.DB
}

func NewMatrixRenstraServiceImpl(
	matrixRenstraRepository repository.MatrixRenstraRepository,
	periodeRepository repository.PeriodeRepository,
	pegawaiRepository repository.PegawaiRepository,
	db *sql.DB,
) *MatrixRenstraServiceImpl {
	return &MatrixRenstraServiceImpl{
		MatrixRenstraRepository: matrixRenstraRepository,
		PeriodeRepository:       periodeRepository,
		PegawaiRepository:       pegawaiRepository,
		DB:                      db,
	}
}

func (service *MatrixRenstraServiceImpl) GetByKodeSubKegiatan(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]programkegiatan.UrusanDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// cek indikator matrix renstra
	indRenstra, err := service.MatrixRenstraRepository.FindIndikatorRenstra(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	if len(indRenstra) == 0 {
		log.Printf("MATRIX RENSTRA KOSONG KODE OPD - %s | TAHUN %s - %s", kodeOpd, tahunAwal, tahunAkhir)
	}

	// fallback ind renstra kosong
	indLama, err := service.MatrixRenstraRepository.FindIndikatorLama(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	if len(indLama) == 0 {
		log.Printf("INDIKATOR LAMA KOSONG")
	}

	final := make(map[string]domain.Indikator)

	// base dari lama

	for _, ind := range indLama {
		final[ind.Kode] = ind
	}

	// override oleh renstra

	for _, ind := range indRenstra {
		final[ind.Kode] = ind
	}
	var indikatorGabungan []domain.Indikator

	for _, v := range final {
		indikatorGabungan = append(indikatorGabungan, v)
	}

	subMap := make(map[string]domain.Indikator)

	kegiatanMap := make(map[string]domain.Indikator)

	programMap := make(map[string]domain.Indikator)

	for _, ind := range indikatorGabungan { // hasil merge renstra + lama

		kode := ind.Kode
		switch {
		case len(kode) == 17: // subkegiatan
			subMap[kode] = ind
		case len(kode) == 12: // kegiatan
			kegiatanMap[kode] = ind
		default: // program
			programMap[kode] = ind
		}

	}

	data, err := service.MatrixRenstraRepository.GetByKodeSubKegiatan(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	// ubah data, inject semua indikator kelama kalau renstra kosong
	for i, item := range data {

		// skip kalau sudah ada indikator renstra
		if item.Indikator != "" {
			continue
		}
		if ind, ok := subMap[item.KodeSubKegiatan]; ok {
			log.Println("SUBKEGIATAN INJECT")
			data[i] = injectIndikator(item, ind)
			continue
		}
		if ind, ok := kegiatanMap[item.KodeKegiatan]; ok {
			log.Println("KEGIATAN INJECT")
			data[i] = injectIndikator(item, ind)
			continue
		}
		if ind, ok := programMap[item.KodeProgram]; ok {
			log.Println("PROGRAM INJECT")
			data[i] = injectIndikator(item, ind)
			continue
		}

	}

	result := service.transformToResponse(data, kodeOpd, tahunAwal, tahunAkhir)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func buildTahunRange(tahunAwal, tahunAkhir string) ([]string, error) {
	tahunAwalInt, errAwal := strconv.Atoi(tahunAwal)
	tahunAkhirInt, errAkhir := strconv.Atoi(tahunAkhir)
	if errAwal != nil || errAkhir != nil {
		return nil, fmt.Errorf("tahun_awal dan tahun_akhir harus berupa angka")
	}
	if tahunAkhirInt < tahunAwalInt {
		return nil, fmt.Errorf("tahun_akhir tidak boleh lebih kecil dari tahun_awal")
	}
	tahunRange := make([]string, 0, tahunAkhirInt-tahunAwalInt+1)
	for t := tahunAwalInt; t <= tahunAkhirInt; t++ {
		tahunRange = append(tahunRange, strconv.Itoa(t))
	}
	return tahunRange, nil
}

func fillTargetByTahunRange(tahunRange []string, targets []domain.Target) []programkegiatan.TargetResponse {
	byTahun := make(map[string]domain.Target, len(targets))
	for _, t := range targets {
		if t.Tahun == "" {
			continue
		}
		if _, exists := byTahun[t.Tahun]; !exists {
			byTahun[t.Tahun] = t
		}
	}
	result := make([]programkegiatan.TargetResponse, 0, len(tahunRange))
	for _, th := range tahunRange {
		if t, ok := byTahun[th]; ok {
			result = append(result, programkegiatan.TargetResponse{
				Id:          t.Id,
				IndikatorId: t.IndikatorId,
				Tahun:       th,
				Target:      t.Target,
				Satuan:      t.Satuan,
			})
			continue
		}
		result = append(result, programkegiatan.TargetResponse{
			Tahun:  th,
			Target: "-",
			Satuan: "-",
		})
	}
	return result
}

func (service *MatrixRenstraServiceImpl) GetByKodeSubKegiatanVersiKedua(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]programkegiatan.UrusanDetailV2Response, error) {
	tahunRange, err := buildTahunRange(tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := service.MatrixRenstraRepository.GetHierarchyAndPagu(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	indList, err := service.MatrixRenstraRepository.FindIndikatorRenstraPeriod(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	indByKode := make(map[string][]programkegiatan.IndikatorPeriodResponse, len(indList))
	for _, ind := range indList {
		indByKode[ind.Kode] = append(indByKode[ind.Kode], programkegiatan.IndikatorPeriodResponse{
			KodeIndikator: ind.KodeIndikator,
			Kode:          ind.Kode,
			KodeOpd:       ind.KodeOpd,
			Indikator:     ind.Indikator,
			Tahun:         ind.Tahun,
			Target:        fillTargetByTahunRange(tahunRange, ind.Target),
		})
	}
	getIndikator := func(kode string) []programkegiatan.IndikatorPeriodResponse {
		if list := indByKode[kode]; list != nil {
			return list
		}
		return []programkegiatan.IndikatorPeriodResponse{}
	}

	buildAnggaran := func(paguByTahun map[string]int64) []programkegiatan.PaguAnggaranTotalResponse {
		result := make([]programkegiatan.PaguAnggaranTotalResponse, 0, len(tahunRange))
		for _, th := range tahunRange {
			result = append(result, programkegiatan.PaguAnggaranTotalResponse{
				Tahun:        th,
				PaguAnggaran: paguByTahun[th],
			})
		}
		return result
	}
	type subkegMeta struct{ nama, namaPegawai, pegawaiId, kodeKeg string }
	type kegMeta struct{ nama, kodePrg string }
	type prgMeta struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }
	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)
	paguSubkegByTahun := make(map[string]map[string]int64)
	seenSubkeg := make(map[string]struct{})
	seenKeg := make(map[string]struct{})
	seenPrg := make(map[string]struct{})
	seenBidang := make(map[string]struct{})
	seenUrusan := make(map[string]struct{})
	subkegByKeg := make(map[string][]string)
	kegByPrg := make(map[string][]string)
	prgByBidang := make(map[string][]string)
	bidangByUrusan := make(map[string][]string)
	var urusanOrder []string
	for _, item := range data {
		if item.KodeSubKegiatan == "" {
			continue
		}
		if _, ok := seenSubkeg[item.KodeSubKegiatan]; !ok {
			seenSubkeg[item.KodeSubKegiatan] = struct{}{}
			subkegData[item.KodeSubKegiatan] = subkegMeta{
				nama:        item.NamaSubKegiatan,
				namaPegawai: item.NamaPegawai,
				pegawaiId:   item.PegawaiId,
				kodeKeg:     item.KodeKegiatan,
			}
			if item.KodeKegiatan != "" {
				subkegByKeg[item.KodeKegiatan] = append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
			}
		}
		if paguSubkegByTahun[item.KodeSubKegiatan] == nil {
			paguSubkegByTahun[item.KodeSubKegiatan] = make(map[string]int64)
		}
		paguSubkegByTahun[item.KodeSubKegiatan][item.TahunSubKegiatan] = item.PaguSubKegiatan
		if item.KodeKegiatan != "" {
			if _, ok := seenKeg[item.KodeKegiatan]; !ok {
				seenKeg[item.KodeKegiatan] = struct{}{}
				kegData[item.KodeKegiatan] = kegMeta{nama: item.NamaKegiatan, kodePrg: item.KodeProgram}
				if item.KodeProgram != "" {
					kegByPrg[item.KodeProgram] = append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
				}
			}
		}
		if item.KodeProgram != "" {
			if _, ok := seenPrg[item.KodeProgram]; !ok {
				seenPrg[item.KodeProgram] = struct{}{}
				prgData[item.KodeProgram] = prgMeta{nama: item.NamaProgram, kodeBidang: item.KodeBidangUrusan}
				if item.KodeBidangUrusan != "" {
					prgByBidang[item.KodeBidangUrusan] = append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
				}
			}
		}
		if item.KodeBidangUrusan != "" {
			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {
				seenBidang[item.KodeBidangUrusan] = struct{}{}
				bidangData[item.KodeBidangUrusan] = bidangMeta{nama: item.NamaBidangUrusan, kodeUrusan: item.KodeUrusan}
				if item.KodeUrusan != "" {
					bidangByUrusan[item.KodeUrusan] = append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
				}
			}
		}
		if item.KodeUrusan != "" {
			if _, ok := seenUrusan[item.KodeUrusan]; !ok {
				seenUrusan[item.KodeUrusan] = struct{}{}
				urusanData[item.KodeUrusan] = item.NamaUrusan
				urusanOrder = append(urusanOrder, item.KodeUrusan)
			}
		}
	}
	sumPaguSubkeg := func(kodeSubkegList []string) map[string]int64 {
		hasil := make(map[string]int64, len(tahunRange))
		for _, kodeSubkeg := range kodeSubkegList {
			for _, th := range tahunRange {
				hasil[th] += paguSubkegByTahun[kodeSubkeg][th]
			}
		}
		return hasil
	}
	allSubkegByKeg := func(kodeKeg string) []string { return subkegByKeg[kodeKeg] }
	allSubkegByPrg := func(kodePrg string) []string {
		var result []string
		for _, kodeKeg := range kegByPrg[kodePrg] {
			result = append(result, allSubkegByKeg(kodeKeg)...)
		}
		return result
	}
	allSubkegByBidang := func(kodeBidang string) []string {
		var result []string
		for _, kodePrg := range prgByBidang[kodeBidang] {
			result = append(result, allSubkegByPrg(kodePrg)...)
		}
		return result
	}
	allSubkegByUrusan := func(kodeUrusan string) []string {
		var result []string
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			result = append(result, allSubkegByBidang(kodeBidang)...)
		}
		return result
	}
	grandPaguByTahun := make(map[string]int64, len(tahunRange))
	for kodeSubkeg := range paguSubkegByTahun {
		for _, th := range tahunRange {
			grandPaguByTahun[th] += paguSubkegByTahun[kodeSubkeg][th]
		}
	}
	detail := programkegiatan.UrusanDetailV2Response{
		KodeOpd:           kodeOpd,
		TahunAwal:         tahunAwal,
		TahunAkhir:        tahunAkhir,
		PaguAnggaranTotal: buildAnggaran(grandPaguByTahun),
		Urusan:            make([]programkegiatan.UrusanV2Response, 0, len(urusanOrder)),
	}
	for _, kodeUrusan := range urusanOrder {
		urusanResp := programkegiatan.UrusanV2Response{
			Kode:         kodeUrusan,
			Nama:         urusanData[kodeUrusan],
			Jenis:        "urusans",
			Anggaran:     buildAnggaran(sumPaguSubkeg(allSubkegByUrusan(kodeUrusan))),
			Indikator:    getIndikator(kodeUrusan),
			BidangUrusan: make([]programkegiatan.BidangUrusanV2Response, 0),
		}
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			bd := bidangData[kodeBidang]
			bidangResp := programkegiatan.BidangUrusanV2Response{
				Kode:      kodeBidang,
				Nama:      bd.nama,
				Jenis:     "bidang_urusans",
				Anggaran:  buildAnggaran(sumPaguSubkeg(allSubkegByBidang(kodeBidang))),
				Indikator: getIndikator(kodeBidang),
				Program:   make([]programkegiatan.ProgramV2Response, 0),
			}
			for _, kodePrg := range prgByBidang[kodeBidang] {
				pd := prgData[kodePrg]
				prgResp := programkegiatan.ProgramV2Response{
					Kode:      kodePrg,
					Nama:      pd.nama,
					Jenis:     "programs",
					Anggaran:  buildAnggaran(sumPaguSubkeg(allSubkegByPrg(kodePrg))),
					Indikator: getIndikator(kodePrg),
					Kegiatan:  make([]programkegiatan.KegiatanV2Response, 0),
				}
				for _, kodeKeg := range kegByPrg[kodePrg] {
					kd := kegData[kodeKeg]
					kegResp := programkegiatan.KegiatanV2Response{
						Kode:        kodeKeg,
						Nama:        kd.nama,
						Jenis:       "kegiatans",
						Anggaran:    buildAnggaran(sumPaguSubkeg(allSubkegByKeg(kodeKeg))),
						Indikator:   getIndikator(kodeKeg),
						SubKegiatan: make([]programkegiatan.SubKegiatanV2Response, 0),
					}
					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]
						kegResp.SubKegiatan = append(kegResp.SubKegiatan, programkegiatan.SubKegiatanV2Response{
							Kode:        kodeSubkeg,
							Nama:        sd.nama,
							Jenis:       "subkegiatans",
							PegawaiId:   sd.pegawaiId,
							NamaPegawai: sd.namaPegawai,
							Anggaran:    buildAnggaran(paguSubkegByTahun[kodeSubkeg]),
							Indikator:   getIndikator(kodeSubkeg),
						})
					}
					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
				}
				bidangResp.Program = append(bidangResp.Program, prgResp)
			}
			urusanResp.BidangUrusan = append(urusanResp.BidangUrusan, bidangResp)
		}
		detail.Urusan = append(detail.Urusan, urusanResp)
	}
	return []programkegiatan.UrusanDetailV2Response{detail}, nil
}

func (service *MatrixRenstraServiceImpl) CreateIndikatorV2(ctx context.Context, requests []programkegiatan.IndikatorRenstraV2CreateRequest) ([]programkegiatan.IndikatorV2UpsertResponse, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("indikator tidak boleh kosong")
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	prefixCounter := make(map[string]int)
	responses := make([]programkegiatan.IndikatorV2UpsertResponse, 0, len(requests))
	for _, req := range requests {
		if strings.TrimSpace(req.Kode) == "" || strings.TrimSpace(req.KodeOpd) == "" {
			return nil, fmt.Errorf("kode dan kode_opd wajib diisi")
		}
		if strings.TrimSpace(req.Indikator) == "" {
			return nil, fmt.Errorf("nama indikator wajib diisi")
		}
		if len(req.Target) == 0 {
			return nil, fmt.Errorf("indikator %q harus memiliki minimal 1 target", req.Indikator)
		}
		tahunSeen := make(map[string]struct{}, len(req.Target))
		for i, t := range req.Target {
			th := strings.TrimSpace(t.Tahun)
			if th == "" {
				return nil, fmt.Errorf("target ke-%d pada indikator %q wajib memiliki tahun", i+1, req.Indikator)
			}
			if _, dup := tahunSeen[th]; dup {
				return nil, fmt.Errorf("tahun target %s duplikat pada indikator %q", th, req.Indikator)
			}
			tahunSeen[th] = struct{}{}
			if err := helper.ValidateTargetRawString(t.Target); err != nil {
				return nil, fmt.Errorf("target tahun %s pada indikator %q: %w", th, req.Indikator, err)
			}
		}

		kodeIndikator := strings.TrimSpace(req.KodeIndikator)
		if kodeIndikator == "" {
			prefix := fmt.Sprintf("RENS-%s", req.KodeOpd)
			if _, loaded := prefixCounter[prefix]; !loaded {
				count, err := service.MatrixRenstraRepository.CountKodeIndikatorByPrefix(ctx, tx, prefix)
				if err != nil {
					return nil, err
				}
				prefixCounter[prefix] = count
			}
			prefixCounter[prefix]++
			rnd, err := randomUint31()
			if err != nil {
				return nil, err
			}
			kodeIndikator = fmt.Sprintf("%s-%03d-%d", prefix, prefixCounter[prefix], rnd)
		} else {
			_, err := service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
		}

		ind := domain.Indikator{
			KodeIndikator: kodeIndikator,
			Kode:          req.Kode,
			KodeOpd:       req.KodeOpd,
			Indikator:     req.Indikator,
			Tahun:         "",
			Jenis:         "renstra",
		}
		if err := service.MatrixRenstraRepository.UpsertIndikator(ctx, tx, ind); err != nil {
			return nil, err
		}

		targetResp := make([]programkegiatan.TargetResponse, 0, len(req.Target))
		for _, t := range req.Target {
			th := strings.TrimSpace(t.Tahun)
			targetId := strings.TrimSpace(t.Id)
			if targetId == "" {
				existingTarget, findErr := service.MatrixRenstraRepository.FindTargetByIndikatorIdAndTahun(ctx, tx, kodeIndikator, th)
				if findErr != nil && findErr != sql.ErrNoRows {
					return nil, findErr
				}
				if findErr == nil && existingTarget.Id != "" {
					targetId = existingTarget.Id
				} else {
					targetId = fmt.Sprintf("TRG-RNST-%s-%s", kodeIndikator, th)
				}
			}
			targetJenis := strings.TrimSpace(t.Jenis)
			if targetJenis == "" {
				targetJenis = "renstra"
			}
			target := domain.Target{
				Id:          targetId,
				IndikatorId: kodeIndikator,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       th,
				Jenis:       targetJenis,
			}
			if err := service.MatrixRenstraRepository.UpsertTarget(ctx, tx, target); err != nil {
				return nil, err
			}
			targetResp = append(targetResp, programkegiatan.TargetResponse{
				Id:          targetId,
				IndikatorId: kodeIndikator,
				Tahun:       th,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Jenis:       targetJenis,
			})
		}

		responses = append(responses, programkegiatan.IndikatorV2UpsertResponse{
			KodeIndikator: kodeIndikator,
			Kode:          req.Kode,
			KodeOpd:       req.KodeOpd,
			Indikator:     req.Indikator,
			Jenis:         "renstra",
			Target:        targetResp,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	helper.PublishAuditAfterCommit(ctx, helper.CreateEventRequest{
		Action:     "CREATE",
		EntityType: "indikator_matrix_renstra",
		EntityID:   responses[0].KodeIndikator,
		After:      responses,
		Metadata: map[string]any{
			"count": len(responses),
		},
	})
	return responses, nil
}

func (service *MatrixRenstraServiceImpl) UpdateIndikatorRenstra(ctx context.Context, request programkegiatan.IndikatorRenstraUpdateRequest) (programkegiatan.IndikatorRenstraUpdateResponse, error) {
	kodeIndikator := strings.TrimSpace(request.KodeIndikator)
	indikator := strings.TrimSpace(request.Indikator)
	if kodeIndikator == "" {
		return programkegiatan.IndikatorRenstraUpdateResponse{}, fmt.Errorf("kode_indikator wajib diisi")
	}
	if indikator == "" {
		return programkegiatan.IndikatorRenstraUpdateResponse{}, fmt.Errorf("indikator wajib diisi")
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.IndikatorRenstraUpdateResponse{}, err
	}
	defer tx.Rollback()

	existing, err := service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return programkegiatan.IndikatorRenstraUpdateResponse{}, fmt.Errorf("indikator %s tidak ditemukan", kodeIndikator)
		}
		return programkegiatan.IndikatorRenstraUpdateResponse{}, err
	}

	before := programkegiatan.IndikatorRenstraUpdateResponse{
		KodeIndikator: existing.KodeIndikator,
		Kode:          existing.Kode,
		KodeOpd:       existing.KodeOpd,
		Indikator:     existing.Indikator,
		Tahun:         existing.Tahun,
		Jenis:         "renstra",
	}

	if err := service.MatrixRenstraRepository.UpdateIndikatorRenstra(ctx, tx, kodeIndikator, indikator); err != nil {
		if err == sql.ErrNoRows {
			return programkegiatan.IndikatorRenstraUpdateResponse{}, fmt.Errorf("indikator %s tidak ditemukan", kodeIndikator)
		}
		return programkegiatan.IndikatorRenstraUpdateResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return programkegiatan.IndikatorRenstraUpdateResponse{}, err
	}

	after := before
	after.Indikator = indikator

	helper.PublishAuditAfterCommit(ctx, helper.CreateEventRequest{
		Action:     "UPDATE",
		EntityType: "indikator_matrix_renstra",
		EntityID:   kodeIndikator,
		Before:     before,
		After:      after,
	})
	return after, nil
}

func (service *MatrixRenstraServiceImpl) UpsertTarget(ctx context.Context, request programkegiatan.TargetRenstraUpsertRequest) (programkegiatan.TargetResponse, error) {
	kodeIndikator := strings.TrimSpace(request.KodeIndikator)
	tahun := strings.TrimSpace(request.Tahun)
	if kodeIndikator == "" || tahun == "" {
		return programkegiatan.TargetResponse{}, fmt.Errorf("kode_indikator dan tahun wajib diisi")
	}
	if err := helper.ValidateTargetRawString(request.Target); err != nil {
		return programkegiatan.TargetResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.TargetResponse{}, err
	}
	defer tx.Rollback()

	ind, err := service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return programkegiatan.TargetResponse{}, fmt.Errorf("indikator %s tidak ditemukan", kodeIndikator)
		}
		return programkegiatan.TargetResponse{}, err
	}

	var before any
	targetId := strings.TrimSpace(request.Id)
	existing, findErr := service.MatrixRenstraRepository.FindTargetByIndikatorIdAndTahun(ctx, tx, kodeIndikator, tahun)
	action := "CREATE"
	if findErr != nil && findErr != sql.ErrNoRows {
		return programkegiatan.TargetResponse{}, findErr
	}
	if findErr == nil && existing.Id != "" {
		action = "UPDATE"
		before = existing
		if targetId == "" {
			targetId = existing.Id
		}
	}
	if targetId == "" {
		targetId = fmt.Sprintf("TRG-RNST-%s-%s", kodeIndikator, tahun)
	}

	targetJenis := "renstra"
	if findErr == nil && existing.Jenis != "" {
		targetJenis = existing.Jenis
	}

	target := domain.Target{
		Id:          targetId,
		IndikatorId: kodeIndikator,
		Target:      request.Target,
		Satuan:      request.Satuan,
		Tahun:       tahun,
		Jenis:       targetJenis,
	}
	if err := service.MatrixRenstraRepository.UpsertTarget(ctx, tx, target); err != nil {
		return programkegiatan.TargetResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return programkegiatan.TargetResponse{}, err
	}

	resp := programkegiatan.TargetResponse{
		Id:          targetId,
		IndikatorId: kodeIndikator,
		Tahun:       tahun,
		Target:      request.Target,
		Satuan:      request.Satuan,
		Jenis:       targetJenis,
	}
	helper.PublishAuditAfterCommit(ctx, helper.CreateEventRequest{
		Action:     action,
		EntityType: "target_renstra",
		EntityID:   targetId,
		Before:     before,
		After:      resp,
		Metadata: map[string]any{
			"kode_indikator": kodeIndikator,
			"kode":           ind.Kode,
			"tahun":          tahun,
		},
	})
	return resp, nil
}

// transformToResponse membangun hierarki dari data flat hasil query.
// Optimasi:
//   - Tidak ada N+1 query (NamaPegawai sudah dari JOIN di repository)
//   - Single pass untuk kumpulkan semua metadata + pagu + indikator
//   - Map-based deduplication (tidak ada linear search di dalam loop)
//   - Pagu dari tb_pagu (jenis='renstra') ditampilkan di luar indikator
func (service *MatrixRenstraServiceImpl) transformToResponse(
	data []domain.SubKegiatanQuery,
	kodeOpd, tahunAwal, tahunAkhir string,
) []programkegiatan.UrusanDetailResponse {
	if len(data) == 0 {
		return []programkegiatan.UrusanDetailResponse{}
	}
	tahunAwalInt, _ := strconv.Atoi(tahunAwal)
	tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)
	tahunRange := make([]string, 0, tahunAkhirInt-tahunAwalInt+1)
	for t := tahunAwalInt; t <= tahunAkhirInt; t++ {
		tahunRange = append(tahunRange, strconv.Itoa(t))
	}
	buildAnggaran := func(paguByTahun map[string]int64) []programkegiatan.PaguAnggaranTotalResponse {
		result := make([]programkegiatan.PaguAnggaranTotalResponse, 0, len(tahunRange))
		for _, th := range tahunRange {
			result = append(result, programkegiatan.PaguAnggaranTotalResponse{
				Tahun:        th,
				PaguAnggaran: paguByTahun[th],
			})
		}
		return result
	}
	// Indikator: IndikatorMatrixResponse dengan Target & Satuan flat
	type indEntry struct {
		resp    programkegiatan.IndikatorMatrixResponse
		targets []programkegiatan.TargetResponse
	}
	indikatorByKodeTahun := make(map[string]map[string]map[string]*indEntry)
	indikatorOrderByKodeTahun := make(map[string]map[string][]string)
	collectIndikator := func(item domain.SubKegiatanQuery) {
		if item.IndikatorId == "" {
			return
		}
		kode := item.IndikatorKode
		th := item.IndikatorTahun
		if indikatorByKodeTahun[kode] == nil {
			indikatorByKodeTahun[kode] = make(map[string]map[string]*indEntry)
		}
		if indikatorByKodeTahun[kode][th] == nil {
			indikatorByKodeTahun[kode][th] = make(map[string]*indEntry)
		}
		if indikatorOrderByKodeTahun[kode] == nil {
			indikatorOrderByKodeTahun[kode] = make(map[string][]string)
		}
		ent, exists := indikatorByKodeTahun[kode][th][item.IndikatorId]
		if !exists {
			ent = &indEntry{
				resp: programkegiatan.IndikatorMatrixResponse{
					KodeIndikator: item.IndikatorId,
					Kode:          kode,
					KodeOpd:       kodeOpd,
					Indikator:     item.Indikator,
					Tahun:         th,
					Target:        item.Target,
					Satuan:        item.Satuan,
				},
				targets: make([]programkegiatan.TargetResponse, 0),
			}
			indikatorByKodeTahun[kode][th][item.IndikatorId] = ent
			indikatorOrderByKodeTahun[kode][th] = append(indikatorOrderByKodeTahun[kode][th], item.IndikatorId)
		} else if item.TargetId != "" && ent.resp.Target == "" {
			ent.resp.Target = item.Target
			ent.resp.Satuan = item.Satuan
		}
		if item.TargetId != "" || item.Target != "" {
			already := false
			for _, t := range ent.targets {
				if t.Id == item.TargetId && item.TargetId != "" {
					already = true
					break
				}
				if item.TargetId == "" && t.IndikatorId == item.IndikatorId && t.Tahun == th && t.Target == item.Target {
					already = true
					break
				}
			}
			if !already {
				ent.targets = append(ent.targets, programkegiatan.TargetResponse{
					Id:          item.TargetId,
					IndikatorId: item.IndikatorId,
					Tahun:       th,
					Target:      item.Target,
					Satuan:      item.Satuan,
				})
			}
		}
	}
	getIndikator := func(kode string) []programkegiatan.IndikatorMatrixResponse {
		tahunMap, ok := indikatorByKodeTahun[kode]
		if !ok {
			return []programkegiatan.IndikatorMatrixResponse{}
		}
		result := make([]programkegiatan.IndikatorMatrixResponse, 0)
		for _, th := range tahunRange {
			for _, id := range indikatorOrderByKodeTahun[kode][th] {
				if ent := tahunMap[th][id]; ent != nil {
					result = append(result, ent.resp)
				}
			}
		}
		return result
	}
	// indikator_baseline: 1 indikator tahun awal (indikator pertama jika >1).
	// target_baseline: 1 slot per tahun (tahun_awal s.d. tahun_akhir),
	// diisi dari target indikator biasa yang sudah ada (indikator pertama tiap tahun).
	getIndikatorBaseline := func(kode string) []programkegiatan.IndikatorBaselineResponse {
		tahunMap, ok := indikatorByKodeTahun[kode]
		if !ok {
			return []programkegiatan.IndikatorBaselineResponse{}
		}
		orderAwal := indikatorOrderByKodeTahun[kode][tahunAwal]
		if len(orderAwal) == 0 {
			return []programkegiatan.IndikatorBaselineResponse{}
		}
		entAwal := tahunMap[tahunAwal][orderAwal[0]]
		if entAwal == nil {
			return []programkegiatan.IndikatorBaselineResponse{}
		}
		targets := make([]programkegiatan.TargetResponse, 0, len(tahunRange))
		for _, th := range tahunRange {
			tr := programkegiatan.TargetResponse{
				Tahun: th,
			}
			orderTh := indikatorOrderByKodeTahun[kode][th]
			if len(orderTh) > 0 {
				if entTh := tahunMap[th][orderTh[0]]; entTh != nil {
					tr.Id = entTh.resp.KodeIndikator
					tr.IndikatorId = entTh.resp.KodeIndikator
					tr.Target = entTh.resp.Target
					tr.Satuan = entTh.resp.Satuan
					if len(entTh.targets) > 0 {
						first := entTh.targets[0]
						if first.Id != "" {
							tr.Id = first.Id
						}
						if first.Target != "" {
							tr.Target = first.Target
							tr.Satuan = first.Satuan
						}
						if first.IndikatorId != "" {
							tr.IndikatorId = first.IndikatorId
						}
					}
				}
			}
			targets = append(targets, tr)
		}
		return []programkegiatan.IndikatorBaselineResponse{
			{
				KodeIndikator:  entAwal.resp.KodeIndikator,
				Kode:           entAwal.resp.Kode,
				KodeOpd:        entAwal.resp.KodeOpd,
				Indikator:      entAwal.resp.Indikator,
				Tahun:          entAwal.resp.Tahun,
				TargetBaseline: targets,
			},
		}
	}
	type subkegMeta struct{ nama, namaPegawai, pegawaiId, kodeKeg string }
	type kegMeta struct{ nama, kodePrg string }
	type prgMeta struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }
	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)
	paguSubkegByTahun := make(map[string]map[string]int64)
	seenSubkeg := make(map[string]struct{})
	seenKeg := make(map[string]struct{})
	seenPrg := make(map[string]struct{})
	seenBidang := make(map[string]struct{})
	seenUrusan := make(map[string]struct{})
	subkegByKeg := make(map[string][]string)
	kegByPrg := make(map[string][]string)
	prgByBidang := make(map[string][]string)
	bidangByUrusan := make(map[string][]string)
	var urusanOrder []string
	for _, item := range data {
		collectIndikator(item)
		if item.KodeSubKegiatan == "" {
			continue
		}
		if _, ok := seenSubkeg[item.KodeSubKegiatan]; !ok {
			seenSubkeg[item.KodeSubKegiatan] = struct{}{}
			subkegData[item.KodeSubKegiatan] = subkegMeta{
				nama:        item.NamaSubKegiatan,
				namaPegawai: item.NamaPegawai,
				pegawaiId:   item.PegawaiId,
				kodeKeg:     item.KodeKegiatan,
			}
			if item.KodeKegiatan != "" {
				subkegByKeg[item.KodeKegiatan] = append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
			}
		}
		if paguSubkegByTahun[item.KodeSubKegiatan] == nil {
			paguSubkegByTahun[item.KodeSubKegiatan] = make(map[string]int64)
		}
		paguSubkegByTahun[item.KodeSubKegiatan][item.TahunSubKegiatan] = item.PaguSubKegiatan
		if item.KodeKegiatan != "" {
			if _, ok := seenKeg[item.KodeKegiatan]; !ok {
				seenKeg[item.KodeKegiatan] = struct{}{}
				kegData[item.KodeKegiatan] = kegMeta{nama: item.NamaKegiatan, kodePrg: item.KodeProgram}
				if item.KodeProgram != "" {
					kegByPrg[item.KodeProgram] = append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
				}
			}
		}
		if item.KodeProgram != "" {
			if _, ok := seenPrg[item.KodeProgram]; !ok {
				seenPrg[item.KodeProgram] = struct{}{}
				prgData[item.KodeProgram] = prgMeta{nama: item.NamaProgram, kodeBidang: item.KodeBidangUrusan}
				if item.KodeBidangUrusan != "" {
					prgByBidang[item.KodeBidangUrusan] = append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
				}
			}
		}
		if item.KodeBidangUrusan != "" {
			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {
				seenBidang[item.KodeBidangUrusan] = struct{}{}
				bidangData[item.KodeBidangUrusan] = bidangMeta{nama: item.NamaBidangUrusan, kodeUrusan: item.KodeUrusan}
				if item.KodeUrusan != "" {
					bidangByUrusan[item.KodeUrusan] = append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
				}
			}
		}
		if item.KodeUrusan != "" {
			if _, ok := seenUrusan[item.KodeUrusan]; !ok {
				seenUrusan[item.KodeUrusan] = struct{}{}
				urusanData[item.KodeUrusan] = item.NamaUrusan
				urusanOrder = append(urusanOrder, item.KodeUrusan)
			}
		}
	}
	sumPaguSubkeg := func(kodeSubkegList []string) map[string]int64 {
		hasil := make(map[string]int64, len(tahunRange))
		for _, kodeSubkeg := range kodeSubkegList {
			for _, th := range tahunRange {
				hasil[th] += paguSubkegByTahun[kodeSubkeg][th]
			}
		}
		return hasil
	}
	allSubkegByKeg := func(kodeKeg string) []string {
		return subkegByKeg[kodeKeg]
	}
	allSubkegByPrg := func(kodePrg string) []string {
		var result []string
		for _, kodeKeg := range kegByPrg[kodePrg] {
			result = append(result, allSubkegByKeg(kodeKeg)...)
		}
		return result
	}
	allSubkegByBidang := func(kodeBidang string) []string {
		var result []string
		for _, kodePrg := range prgByBidang[kodeBidang] {
			result = append(result, allSubkegByPrg(kodePrg)...)
		}
		return result
	}
	allSubkegByUrusan := func(kodeUrusan string) []string {
		var result []string
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			result = append(result, allSubkegByBidang(kodeBidang)...)
		}
		return result
	}
	grandPaguByTahun := make(map[string]int64, len(tahunRange))
	for kodeSubkeg := range paguSubkegByTahun {
		for _, th := range tahunRange {
			grandPaguByTahun[th] += paguSubkegByTahun[kodeSubkeg][th]
		}
	}
	urusanDetail := programkegiatan.UrusanDetailResponse{
		KodeOpd:           kodeOpd,
		TahunAwal:         tahunAwal,
		TahunAkhir:        tahunAkhir,
		PaguAnggaranTotal: buildAnggaran(grandPaguByTahun),
		Urusan:            make([]programkegiatan.UrusanResponse, 0, len(urusanOrder)),
	}
	for _, kodeUrusan := range urusanOrder {
		paguUrusan := sumPaguSubkeg(allSubkegByUrusan(kodeUrusan))
		urusanResp := programkegiatan.UrusanResponse{
			Kode:              kodeUrusan,
			Nama:              urusanData[kodeUrusan],
			Jenis:             "urusans",
			Anggaran:          buildAnggaran(paguUrusan),
			Indikator:         getIndikator(kodeUrusan),
			IndikatorBaseline: getIndikatorBaseline(kodeUrusan),
			BidangUrusan:      make([]programkegiatan.BidangUrusanResponse, 0),
		}
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			paguBidang := sumPaguSubkeg(allSubkegByBidang(kodeBidang))
			bd := bidangData[kodeBidang]
			bidangResp := programkegiatan.BidangUrusanResponse{
				Kode:              kodeBidang,
				Nama:              bd.nama,
				Jenis:             "bidang_urusans",
				Anggaran:          buildAnggaran(paguBidang),
				Indikator:         getIndikator(kodeBidang),
				IndikatorBaseline: getIndikatorBaseline(kodeBidang),
				Program:           make([]programkegiatan.ProgramResponse, 0),
			}
			for _, kodePrg := range prgByBidang[kodeBidang] {
				paguPrg := sumPaguSubkeg(allSubkegByPrg(kodePrg))
				pd := prgData[kodePrg]
				prgResp := programkegiatan.ProgramResponse{
					Kode:              kodePrg,
					Nama:              pd.nama,
					Jenis:             "programs",
					Anggaran:          buildAnggaran(paguPrg),
					Indikator:         getIndikator(kodePrg),
					IndikatorBaseline: getIndikatorBaseline(kodePrg),
					Kegiatan:          make([]programkegiatan.KegiatanResponse, 0),
				}
				for _, kodeKeg := range kegByPrg[kodePrg] {
					paguKeg := sumPaguSubkeg(allSubkegByKeg(kodeKeg))
					kd := kegData[kodeKeg]
					kegResp := programkegiatan.KegiatanResponse{
						Kode:              kodeKeg,
						Nama:              kd.nama,
						Jenis:             "kegiatans",
						Anggaran:          buildAnggaran(paguKeg),
						Indikator:         getIndikator(kodeKeg),
						IndikatorBaseline: getIndikatorBaseline(kodeKeg),
						SubKegiatan:       make([]programkegiatan.SubKegiatanResponse, 0),
					}
					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]
						subkegResp := programkegiatan.SubKegiatanResponse{
							Kode:              kodeSubkeg,
							Nama:              sd.nama,
							Jenis:             "subkegiatans",
							PegawaiId:         sd.pegawaiId,
							NamaPegawai:       sd.namaPegawai,
							Anggaran:          buildAnggaran(paguSubkegByTahun[kodeSubkeg]),
							Indikator:         getIndikator(kodeSubkeg),
							IndikatorBaseline: getIndikatorBaseline(kodeSubkeg),
						}
						kegResp.SubKegiatan = append(kegResp.SubKegiatan, subkegResp)
					}
					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
				}
				bidangResp.Program = append(bidangResp.Program, prgResp)
			}
			urusanResp.BidangUrusan = append(urusanResp.BidangUrusan, bidangResp)
		}
		urusanDetail.Urusan = append(urusanDetail.Urusan, urusanResp)
	}
	return []programkegiatan.UrusanDetailResponse{urusanDetail}
}

// crud

func (service *MatrixRenstraServiceImpl) DeleteIndikator(ctx context.Context, kodeIndikator string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("indikator %s tidak ditemukan", kodeIndikator)
		}
		return err
	}
	if err = service.MatrixRenstraRepository.DeleteTargetByIndikatorId(ctx, tx, kodeIndikator); err != nil {
		return err
	}
	if err = service.MatrixRenstraRepository.DeleteIndikator(ctx, tx, kodeIndikator); err != nil {
		return err
	}
	return tx.Commit()
}

func (service *MatrixRenstraServiceImpl) FindIndikatorByKodeIndikator(ctx context.Context, kodeIndikator string) (programkegiatan.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}
	defer tx.Rollback()
	ind, err := service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
	if err != nil {
		if err == sql.ErrNoRows {
			return programkegiatan.IndikatorResponse{}, fmt.Errorf("indikator %s tidak ditemukan", kodeIndikator)
		}
		return programkegiatan.IndikatorResponse{}, err
	}
	resp := programkegiatan.IndikatorResponse{
		KodeIndikator: ind.KodeIndikator,
		Kode:          ind.Kode,
		KodeOpd:       ind.KodeOpd,
		Indikator:     ind.Indikator,
		Tahun:         ind.Tahun,
		Target:        make([]programkegiatan.TargetResponse, 0),
	}
	for _, t := range ind.Target {
		resp.Target = append(resp.Target, programkegiatan.TargetResponse{
			Id:     t.Id,
			Target: t.Target,
			Satuan: t.Satuan,
		})
	}
	_ = tx.Commit()
	return resp, nil
}

func (service *MatrixRenstraServiceImpl) UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenstraRequest) (programkegiatan.AnggaranRenstraResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.AnggaranRenstraResponse{}, err
	}
	defer tx.Rollback()
	err = service.MatrixRenstraRepository.UpsertAnggaran(
		ctx, tx,
		request.KodeSubKegiatan,
		request.KodeOpd,
		request.Tahun,
		request.Pagu,
	)
	if err != nil {
		return programkegiatan.AnggaranRenstraResponse{}, err
	}
	// ← TAMBAHKAN INI sebelum return
	if err = tx.Commit(); err != nil {
		return programkegiatan.AnggaranRenstraResponse{}, err
	}
	resp := programkegiatan.AnggaranRenstraResponse{
		KodeSubKegiatan: request.KodeSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
		Pagu:            request.Pagu,
	}
	helper.PublishAuditAfterCommit(ctx, helper.CreateEventRequest{
		Action:     "UPDATE",
		EntityType: "pagu_renstra",
		EntityID:   request.KodeSubKegiatan,
		After:      resp,
		Metadata: map[string]any{
			"kode_opd": request.KodeOpd,
			"tahun":    request.Tahun,
		},
	})
	return resp, nil
}

func (service *MatrixRenstraServiceImpl) UpsertBatchIndikator(ctx context.Context, requests []programkegiatan.IndikatorRenstraCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var responses []programkegiatan.IndikatorUpsertResponse
	prefixCounter := make(map[string]int)
	// Kumpulkan kode_indikator yang diproses per scope (kode+kodeOpd+tahun)
	// untuk keperluan delete-not-in-list di akhir
	type scopeKey struct{ kode, kodeOpd, tahun string }
	processedPerScope := make(map[scopeKey][]string)
	for _, req := range requests {
		scope := scopeKey{req.Kode, req.KodeOpd, req.Tahun}
		kodeIndikator := req.KodeIndikator
		existingTargetId := ""
		if kodeIndikator == "" {
			// CREATE: urutan + bilangan acak agar kode_indikator jarang bentrok (ON DUPLICATE KEY)
			prefix := fmt.Sprintf("RENS-%s-%s", req.KodeOpd, req.Tahun)
			if _, loaded := prefixCounter[prefix]; !loaded {
				count, err := service.MatrixRenstraRepository.CountKodeIndikatorByPrefix(ctx, tx, prefix)
				if err != nil {
					return nil, err
				}
				prefixCounter[prefix] = count
			}
			prefixCounter[prefix]++
			rnd, err := randomUint31()
			if err != nil {
				return nil, err
			}
			kodeIndikator = fmt.Sprintf("%s-%03d-%d", prefix, prefixCounter[prefix], rnd)
		} else {
			// UPDATE: ambil target.id lama dari DB
			existing, err := service.MatrixRenstraRepository.FindIndikatorByKodeIndikator(ctx, tx, kodeIndikator)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			if len(existing.Target) > 0 && existing.Target[0].Id != "" {
				existingTargetId = existing.Target[0].Id
			}
		}
		// Catat kode_indikator ini sebagai "keep" untuk scope-nya
		processedPerScope[scope] = append(processedPerScope[scope], kodeIndikator)
		// Upsert indikator
		ind := domain.Indikator{
			KodeIndikator: kodeIndikator,
			Kode:          req.Kode,
			KodeOpd:       req.KodeOpd,
			Indikator:     req.Indikator,
			Tahun:         req.Tahun,
			Jenis:         "renstra",
		}
		if err := service.MatrixRenstraRepository.UpsertIndikator(ctx, tx, ind); err != nil {
			return nil, err
		}
		// Upsert target
		targetId := existingTargetId
		if targetId == "" {
			targetId = fmt.Sprintf("TRG-RNST-%s", kodeIndikator)
		}
		target := domain.Target{
			Id:          targetId,
			IndikatorId: kodeIndikator,
			Target:      req.Target,
			Satuan:      req.Satuan,
			Tahun:       req.Tahun,
		}
		if err := service.MatrixRenstraRepository.UpsertTarget(ctx, tx, target); err != nil {
			return nil, err
		}
		responses = append(responses, programkegiatan.IndikatorUpsertResponse{
			KodeIndikator: kodeIndikator,
			Kode:          req.Kode,
			KodeOpd:       req.KodeOpd,
			Indikator:     req.Indikator,
			Tahun:         req.Tahun,
			Jenis:         "renstra",
			Target:        req.Target,
			Satuan:        req.Satuan,
		})
	}
	// ── SYNC: hapus indikator di DB yang tidak ada di request ──
	for scope, keepList := range processedPerScope {
		err := service.MatrixRenstraRepository.DeleteIndicatorsExcept(
			ctx, tx,
			scope.kode, scope.kodeOpd, scope.tahun,
			keepList,
		)
		if err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	if len(responses) > 0 {
		helper.PublishAuditAfterCommit(ctx, helper.CreateEventRequest{
			Action:     "UPDATE",
			EntityType: "indikator_matrix_renstra",
			EntityID:   responses[0].KodeIndikator,
			After:      responses,
			Metadata: map[string]any{
				"count": len(responses),
			},
		})
	}
	return responses, nil
}

func randomUint31() (uint32, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b[:]) & 0x7fffffff, nil
}

func injectIndikator(item domain.SubKegiatanQuery, ind domain.Indikator) domain.SubKegiatanQuery {
	item.IndikatorId = ind.KodeIndikator
	item.IndikatorKode = ind.Kode
	item.Indikator = ind.Indikator
	for _, tar := range ind.Target {
		item.TargetId = tar.Id
		item.Target = tar.Target
		item.Satuan = tar.Satuan
	}
	item.IndikatorTahun = ind.Tahun
	return item
}
