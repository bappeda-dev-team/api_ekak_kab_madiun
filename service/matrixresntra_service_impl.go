package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
		resp programkegiatan.IndikatorMatrixResponse
	}
	indikatorByKodeTahun := make(map[string]map[string]map[string]*indEntry)
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
			}
			indikatorByKodeTahun[kode][th][item.IndikatorId] = ent
		} else if item.TargetId != "" && ent.resp.Target == "" {
			ent.resp.Target = item.Target
			ent.resp.Satuan = item.Satuan
		}
	}
	getIndikator := func(kode string) []programkegiatan.IndikatorMatrixResponse {
		tahunMap, ok := indikatorByKodeTahun[kode]
		if !ok {
			return []programkegiatan.IndikatorMatrixResponse{}
		}
		result := make([]programkegiatan.IndikatorMatrixResponse, 0)
		for _, th := range tahunRange {
			for _, ent := range tahunMap[th] {
				result = append(result, ent.resp)
			}
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
			Kode:         kodeUrusan,
			Nama:         urusanData[kodeUrusan],
			Jenis:        "urusans",
			Anggaran:     buildAnggaran(paguUrusan),
			Indikator:    getIndikator(kodeUrusan),
			BidangUrusan: make([]programkegiatan.BidangUrusanResponse, 0),
		}
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			paguBidang := sumPaguSubkeg(allSubkegByBidang(kodeBidang))
			bd := bidangData[kodeBidang]
			bidangResp := programkegiatan.BidangUrusanResponse{
				Kode:      kodeBidang,
				Nama:      bd.nama,
				Jenis:     "bidang_urusans",
				Anggaran:  buildAnggaran(paguBidang),
				Indikator: getIndikator(kodeBidang),
				Program:   make([]programkegiatan.ProgramResponse, 0),
			}
			for _, kodePrg := range prgByBidang[kodeBidang] {
				paguPrg := sumPaguSubkeg(allSubkegByPrg(kodePrg))
				pd := prgData[kodePrg]
				prgResp := programkegiatan.ProgramResponse{
					Kode:      kodePrg,
					Nama:      pd.nama,
					Jenis:     "programs",
					Anggaran:  buildAnggaran(paguPrg),
					Indikator: getIndikator(kodePrg),
					Kegiatan:  make([]programkegiatan.KegiatanResponse, 0),
				}
				for _, kodeKeg := range kegByPrg[kodePrg] {
					paguKeg := sumPaguSubkeg(allSubkegByKeg(kodeKeg))
					kd := kegData[kodeKeg]
					kegResp := programkegiatan.KegiatanResponse{
						Kode:        kodeKeg,
						Nama:        kd.nama,
						Jenis:       "kegiatans",
						Anggaran:    buildAnggaran(paguKeg),
						Indikator:   getIndikator(kodeKeg),
						SubKegiatan: make([]programkegiatan.SubKegiatanResponse, 0),
					}
					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]
						subkegResp := programkegiatan.SubKegiatanResponse{
							Kode:        kodeSubkeg,
							Nama:        sd.nama,
							Jenis:       "subkegiatans",
							PegawaiId:   sd.pegawaiId,
							NamaPegawai: sd.namaPegawai,
							Anggaran:    buildAnggaran(paguSubkegByTahun[kodeSubkeg]),
							Indikator:   getIndikator(kodeSubkeg),
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
	return programkegiatan.AnggaranRenstraResponse{
		KodeSubKegiatan: request.KodeSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
		Pagu:            request.Pagu,
	}, nil
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
		item.Target = tar.Target
		item.Satuan = tar.Satuan
	}
	item.IndikatorTahun = ind.Tahun
	return item
}
