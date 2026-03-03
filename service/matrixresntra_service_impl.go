package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strconv"

	"github.com/google/uuid"
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

	data, err := service.MatrixRenstraRepository.GetByKodeSubKegiatan(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
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

	// Build rentang tahun
	tahunAwalInt, _ := strconv.Atoi(tahunAwal)
	tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)
	tahunRange := make([]string, 0, tahunAkhirInt-tahunAwalInt+1)
	for t := tahunAwalInt; t <= tahunAkhirInt; t++ {
		tahunRange = append(tahunRange, strconv.Itoa(t))
	}

	// -----------------------------------------------------------------------
	// Helper: bangun []PaguAnggaranTotalResponse dari map[tahun]pagu
	// -----------------------------------------------------------------------
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

	// -----------------------------------------------------------------------
	// Struktur bantu untuk indikator (dengan dedup target)
	// -----------------------------------------------------------------------
	type indEntry struct {
		resp      programkegiatan.IndikatorResponse
		targetSet map[string]struct{}
	}

	// kode (subkeg/keg/prg/bidang/urusan) → tahun → indikatorId → *indEntry
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
				resp: programkegiatan.IndikatorResponse{
					Id:        item.IndikatorId,
					Kode:      kode,
					KodeOpd:   kodeOpd,
					Indikator: item.Indikator,
					Tahun:     th,
					Target:    make([]programkegiatan.TargetResponse, 0),
				},
				targetSet: make(map[string]struct{}),
			}
			indikatorByKodeTahun[kode][th][item.IndikatorId] = ent
		}
		if item.TargetId != "" {
			if _, seen := ent.targetSet[item.TargetId]; !seen {
				ent.targetSet[item.TargetId] = struct{}{}
				ent.resp.Target = append(ent.resp.Target, programkegiatan.TargetResponse{
					Id:     item.TargetId,
					Target: item.Target,
					Satuan: item.Satuan,
				})
			}
		}
	}

	// Ambil semua indikator untuk kode tertentu sebagai []IndikatorResponse
	// (satu entry per tahun, karena matrix renstra multi-tahun)
	getIndikator := func(kode string) []programkegiatan.IndikatorResponse {
		tahunMap, ok := indikatorByKodeTahun[kode]
		if !ok {
			return []programkegiatan.IndikatorResponse{}
		}
		result := make([]programkegiatan.IndikatorResponse, 0)
		for _, th := range tahunRange {
			for _, ent := range tahunMap[th] {
				result = append(result, ent.resp)
			}
		}
		return result
	}

	// -----------------------------------------------------------------------
	// Metadata tiap level hierarki
	// -----------------------------------------------------------------------
	type subkegMeta struct{ nama, namaPegawai, pegawaiId, kodeKeg string }
	type kegMeta struct{ nama, kodePrg string }
	type prgMeta struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }

	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string) // kodeUrusan → namaUrusan

	// pagu per subkegiatan per tahun (dari tb_pagu, jenis='renstra')
	// key: kodeSubkegiatan → tahun → pagu
	paguSubkegByTahun := make(map[string]map[string]int64)

	// Ordered, deduplicated children
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

	// -----------------------------------------------------------------------
	// PASS 1: satu kali iterasi — kumpulkan semua data
	// -----------------------------------------------------------------------
	for _, item := range data {
		collectIndikator(item)

		if item.KodeSubKegiatan == "" {
			continue
		}

		// Subkegiatan: catat satu kali per (subkegiatan)
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

		// Pagu subkegiatan per tahun (dari tb_pagu, sama nilainya per row — set saja, idempoten)
		if paguSubkegByTahun[item.KodeSubKegiatan] == nil {
			paguSubkegByTahun[item.KodeSubKegiatan] = make(map[string]int64)
		}
		// PaguSubKegiatan sudah di-COALESCE di query, selalu ada nilai (0 jika tidak ada di tb_pagu)
		paguSubkegByTahun[item.KodeSubKegiatan][item.TahunSubKegiatan] = item.PaguSubKegiatan

		// Kegiatan
		if item.KodeKegiatan != "" {
			if _, ok := seenKeg[item.KodeKegiatan]; !ok {
				seenKeg[item.KodeKegiatan] = struct{}{}
				kegData[item.KodeKegiatan] = kegMeta{nama: item.NamaKegiatan, kodePrg: item.KodeProgram}
				if item.KodeProgram != "" {
					kegByPrg[item.KodeProgram] = append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
				}
			}
		}

		// Program
		if item.KodeProgram != "" {
			if _, ok := seenPrg[item.KodeProgram]; !ok {
				seenPrg[item.KodeProgram] = struct{}{}
				prgData[item.KodeProgram] = prgMeta{nama: item.NamaProgram, kodeBidang: item.KodeBidangUrusan}
				if item.KodeBidangUrusan != "" {
					prgByBidang[item.KodeBidangUrusan] = append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
				}
			}
		}

		// Bidang Urusan
		if item.KodeBidangUrusan != "" {
			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {
				seenBidang[item.KodeBidangUrusan] = struct{}{}
				bidangData[item.KodeBidangUrusan] = bidangMeta{nama: item.NamaBidangUrusan, kodeUrusan: item.KodeUrusan}
				if item.KodeUrusan != "" {
					bidangByUrusan[item.KodeUrusan] = append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
				}
			}
		}

		// Urusan
		if item.KodeUrusan != "" {
			if _, ok := seenUrusan[item.KodeUrusan]; !ok {
				seenUrusan[item.KodeUrusan] = struct{}{}
				urusanData[item.KodeUrusan] = item.NamaUrusan
				urusanOrder = append(urusanOrder, item.KodeUrusan)
			}
		}
	}

	// -----------------------------------------------------------------------
	// Helper: hitung pagu per tahun untuk suatu level dari sum subkegiatan
	// -----------------------------------------------------------------------
	sumPaguSubkeg := func(kodeSubkegList []string) map[string]int64 {
		hasil := make(map[string]int64, len(tahunRange))
		for _, kodeSubkeg := range kodeSubkegList {
			for _, th := range tahunRange {
				hasil[th] += paguSubkegByTahun[kodeSubkeg][th]
			}
		}
		return hasil
	}

	// Kumpulkan semua kode subkegiatan per level untuk hitung pagu ke atas
	// kegiatan → program → bidang → urusan
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

	// Grand total pagu per tahun (sum semua subkegiatan)
	grandPaguByTahun := make(map[string]int64, len(tahunRange))
	for _, kodeSubkeg := range func() []string {
		var all []string
		for k := range paguSubkegByTahun {
			all = append(all, k)
		}
		return all
	}() {
		for _, th := range tahunRange {
			grandPaguByTahun[th] += paguSubkegByTahun[kodeSubkeg][th]
		}
	}

	// -----------------------------------------------------------------------
	// PASS 2: bangun response hierarki dari maps yang sudah terkumpul
	// -----------------------------------------------------------------------
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
func (service *MatrixRenstraServiceImpl) CreateIndikator(ctx context.Context, requests []programkegiatan.IndikatorRenstraCreateRequest) ([]programkegiatan.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var responses []programkegiatan.IndikatorResponse

	for _, request := range requests {
		randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		uuId := fmt.Sprintf("IND-RNST-%s", randomDigits)

		indikator := domain.Indikator{
			Id:        uuId,
			Kode:      request.Kode,
			KodeOpd:   request.KodeOpd,
			Indikator: request.Indikator,
			Tahun:     request.Tahun,
			// PaguAnggaran = 0, pagu dihandle UpsertAnggaran
		}

		err = service.MatrixRenstraRepository.SaveIndikator(ctx, tx, indikator)
		if err != nil {
			return nil, err
		}

		uuIdTarget := fmt.Sprintf("TRG-RNST-%s", randomDigits)
		target := domain.Target{
			Id:          uuIdTarget,
			IndikatorId: indikator.Id,
			Target:      request.Target,
			Satuan:      request.Satuan,
			// Jenis = 'renstra' diinject di repository SaveTarget
		}

		err = service.MatrixRenstraRepository.SaveTarget(ctx, tx, target)
		if err != nil {
			return nil, err
		}

		responses = append(responses, programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			Kode:      request.Kode,
			KodeOpd:   request.KodeOpd,
			Indikator: request.Indikator,
			Tahun:     request.Tahun,
			Target: []programkegiatan.TargetResponse{
				{Id: target.Id, Target: request.Target, Satuan: request.Satuan},
			},
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return responses, nil
}

func (service *MatrixRenstraServiceImpl) UpdateIndikator(ctx context.Context, request programkegiatan.UpdateIndikatorRequest) (programkegiatan.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}
	defer tx.Rollback()

	// Cek apakah indikator exists
	existingIndikator, err := service.MatrixRenstraRepository.FindIndikatorById(ctx, tx, request.Id)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	// Update indikator
	indikator := domain.Indikator{
		Id:        request.Id,
		Kode:      request.Kode,
		KodeOpd:   request.KodeOpd,
		Indikator: request.Indikator,
		Tahun:     request.Tahun,
	}

	err = service.MatrixRenstraRepository.UpdateIndikator(ctx, tx, indikator)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	// Update target
	target := domain.Target{
		Id:          existingIndikator.Target[0].Id, // Ambil ID target yang sudah ada
		IndikatorId: request.Id,
		Target:      request.Target,
		Satuan:      request.Satuan,
	}

	err = service.MatrixRenstraRepository.UpdateTarget(ctx, tx, target)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	return programkegiatan.IndikatorResponse{
		Id:        request.Id,
		Kode:      request.Kode,
		KodeOpd:   request.KodeOpd,
		Indikator: request.Indikator,
		Tahun:     request.Tahun,
		Target: []programkegiatan.TargetResponse{
			{
				Id:     target.Id,
				Target: request.Target,
				Satuan: request.Satuan,
			},
		},
	}, nil
}

func (service *MatrixRenstraServiceImpl) DeleteIndikator(ctx context.Context, indikatorId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Hapus target terlebih dahulu (karena foreign key)
	err = service.MatrixRenstraRepository.DeleteTargetByIndikatorId(ctx, tx, indikatorId)
	if err != nil {
		return err
	}

	// Hapus indikator
	err = service.MatrixRenstraRepository.DeleteIndikator(ctx, tx, indikatorId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (service *MatrixRenstraServiceImpl) FindIndikatorById(ctx context.Context, indikatorId string) (programkegiatan.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}
	defer tx.Rollback()

	// Cari indikator berdasarkan ID
	indikator, err := service.MatrixRenstraRepository.FindIndikatorById(ctx, tx, indikatorId)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	// Transform ke response
	response := programkegiatan.IndikatorResponse{
		Id:        indikator.Id,
		Kode:      indikator.Kode,
		KodeOpd:   indikator.KodeOpd,
		Indikator: indikator.Indikator,
		Tahun:     indikator.Tahun,
		// PaguAnggaran: indikator.PaguAnggaran,
		Target: make([]programkegiatan.TargetResponse, 0),
	}

	// Tambahkan target ke response
	if len(indikator.Target) > 0 {
		response.Target = append(response.Target, programkegiatan.TargetResponse{
			Id:     indikator.Target[0].Id,
			Target: indikator.Target[0].Target,
			Satuan: indikator.Target[0].Satuan,
		})
	}

	err = tx.Commit()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	return response, nil
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

	return programkegiatan.AnggaranRenstraResponse{
		KodeSubKegiatan: request.KodeSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
		Pagu:            request.Pagu,
	}, nil
}
