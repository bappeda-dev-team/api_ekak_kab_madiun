package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type MatrixRenjaServiceImpl struct {
	MatrixRenjaRepository repository.MatrixRenjaRepository
	PeriodeRepository     repository.PeriodeRepository
	PegawaiRepository     repository.PegawaiRepository
	DB                    *sql.DB
}

func NewMatrixRenjaServiceImpl(
	matrixRenjaRepository repository.MatrixRenjaRepository,
	periodeRepository repository.PeriodeRepository,
	pegawaiRepository repository.PegawaiRepository,
	db *sql.DB,
) *MatrixRenjaServiceImpl {
	return &MatrixRenjaServiceImpl{
		MatrixRenjaRepository: matrixRenjaRepository,
		PeriodeRepository:     periodeRepository,
		PegawaiRepository:     pegawaiRepository,
		DB:                    db,
	}
}

func (service *MatrixRenjaServiceImpl) GetRenjaRanwal(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := service.MatrixRenjaRepository.GetRenjaRanwal(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	result := service.transformToResponseRanwal(data, kodeOpd, tahun)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *MatrixRenjaServiceImpl) transformToResponseRanwal(
	data []domain.SubKegiatanQuery,
	kodeOpd, tahun string,
) []programkegiatan.UrusanDetailResponse {

	if len(data) == 0 {
		return []programkegiatan.UrusanDetailResponse{}
	}

	// Helper: buat Anggaran slice (single tahun untuk ranwal)
	mkAnggaran := func(pagu int64) []programkegiatan.PaguAnggaranTotalResponse {
		return []programkegiatan.PaguAnggaranTotalResponse{
			{Tahun: tahun, PaguAnggaran: pagu},
		}
	}

	// ── Indikator collector ──
	type indEntry struct {
		resp      programkegiatan.IndikatorResponse
		targetSet map[string]struct{}
	}
	indikatorByKode := make(map[string]map[string]*indEntry)

	collectIndikator := func(item domain.SubKegiatanQuery) {
		if item.IndikatorId == "" {
			return
		}
		kode := item.IndikatorKode
		if indikatorByKode[kode] == nil {
			indikatorByKode[kode] = make(map[string]*indEntry)
		}
		ent, exists := indikatorByKode[kode][item.IndikatorId]
		if !exists {
			ent = &indEntry{
				resp: programkegiatan.IndikatorResponse{
					Id:           item.IndikatorId,
					Kode:         kode,
					KodeOpd:      kodeOpd,
					Indikator:    item.Indikator,
					Tahun:        tahun,
					PaguAnggaran: 0, // pagu di luar indikator
					StatusTarget: false,
					Target:       make([]programkegiatan.TargetResponse, 0),
				},
				targetSet: make(map[string]struct{}),
			}
			indikatorByKode[kode][item.IndikatorId] = ent
		}
		if item.TargetId != "" {
			if _, seen := ent.targetSet[item.TargetId]; !seen {
				ent.targetSet[item.TargetId] = struct{}{}
				ent.resp.Target = append(ent.resp.Target, programkegiatan.TargetResponse{
					Id:     item.TargetId,
					Target: item.Target,
					Satuan: item.Satuan,
				})
				ent.resp.StatusTarget = true
			}
		}
	}

	getIndikator := func(kode string) []programkegiatan.IndikatorResponse {
		m, ok := indikatorByKode[kode]
		if !ok {
			return []programkegiatan.IndikatorResponse{}
		}
		slice := make([]programkegiatan.IndikatorResponse, 0, len(m))
		for _, e := range m {
			slice = append(slice, e.resp)
		}
		return slice
	}

	// ── Metadata ──
	type subkegMeta struct{ nama, kodeKeg, pegawaiId, namaPegawai string }
	type kegMeta struct{ nama, kodePrg string }
	type prgMeta struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }

	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)

	paguSubkeg := make(map[string]int64) // sumber: tb_pagu jenis='ranwal'

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

	// ── PASS 1: satu iterasi ──
	for _, item := range data {
		collectIndikator(item)
		if item.KodeSubKegiatan == "" {
			continue
		}
		if _, ok := seenSubkeg[item.KodeSubKegiatan]; !ok {
			seenSubkeg[item.KodeSubKegiatan] = struct{}{}
			subkegData[item.KodeSubKegiatan] = subkegMeta{
				nama:        item.NamaSubKegiatan,
				kodeKeg:     item.KodeKegiatan,
				pegawaiId:   item.PegawaiId,
				namaPegawai: item.NamaPegawai,
			}
			paguSubkeg[item.KodeSubKegiatan] = item.PaguSubKegiatan // ← dari tb_pagu
			subkegByKeg[item.KodeKegiatan] = append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
		}
		if item.KodeKegiatan != "" {
			if _, ok := seenKeg[item.KodeKegiatan]; !ok {
				seenKeg[item.KodeKegiatan] = struct{}{}
				kegData[item.KodeKegiatan] = kegMeta{nama: item.NamaKegiatan, kodePrg: item.KodeProgram}
				kegByPrg[item.KodeProgram] = append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
			}
		}
		if item.KodeProgram != "" {
			if _, ok := seenPrg[item.KodeProgram]; !ok {
				seenPrg[item.KodeProgram] = struct{}{}
				prgData[item.KodeProgram] = prgMeta{nama: item.NamaProgram, kodeBidang: item.KodeBidangUrusan}
				prgByBidang[item.KodeBidangUrusan] = append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
			}
		}
		if item.KodeBidangUrusan != "" {
			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {
				seenBidang[item.KodeBidangUrusan] = struct{}{}
				bidangData[item.KodeBidangUrusan] = bidangMeta{nama: item.NamaBidangUrusan, kodeUrusan: item.KodeUrusan}
				bidangByUrusan[item.KodeUrusan] = append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
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

	// ── Agregasi pagu bottom-up: subkeg → keg → prg → bidang → urusan ──
	// O(N) total karena setiap node hanya diproses sekali
	paguKeg := make(map[string]int64)
	paguPrg := make(map[string]int64)
	paguBidang := make(map[string]int64)
	paguUrusan := make(map[string]int64)

	for kodeKeg, subkegList := range subkegByKeg {
		for _, ks := range subkegList {
			paguKeg[kodeKeg] += paguSubkeg[ks]
		}
	}
	for kodePrg, kegList := range kegByPrg {
		for _, kk := range kegList {
			paguPrg[kodePrg] += paguKeg[kk]
		}
	}
	for kodeBidang, prgList := range prgByBidang {
		for _, kp := range prgList {
			paguBidang[kodeBidang] += paguPrg[kp]
		}
	}
	for kodeUrusan, bidangList := range bidangByUrusan {
		for _, kb := range bidangList {
			paguUrusan[kodeUrusan] += paguBidang[kb]
		}
	}
	var paguTotal int64
	for _, p := range paguUrusan {
		paguTotal += p
	}

	// ── PASS 2: bangun hierarki ──
	urusanDetail := programkegiatan.UrusanDetailResponse{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
		PaguAnggaranTotal: []programkegiatan.PaguAnggaranTotalResponse{
			{Tahun: tahun, PaguAnggaran: paguTotal},
		},
		Urusan: make([]programkegiatan.UrusanResponse, 0),
	}

	for _, kodeUrusan := range urusanOrder {
		urusanResp := programkegiatan.UrusanResponse{
			Kode:         kodeUrusan,
			Nama:         urusanData[kodeUrusan],
			Jenis:        "urusans",
			Anggaran:     mkAnggaran(paguUrusan[kodeUrusan]),
			Indikator:    getIndikator(kodeUrusan),
			BidangUrusan: make([]programkegiatan.BidangUrusanResponse, 0),
		}
		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			bd := bidangData[kodeBidang]
			bidangResp := programkegiatan.BidangUrusanResponse{
				Kode:      kodeBidang,
				Nama:      bd.nama,
				Jenis:     "bidang_urusans",
				Anggaran:  mkAnggaran(paguBidang[kodeBidang]),
				Indikator: getIndikator(kodeBidang),
				Program:   make([]programkegiatan.ProgramResponse, 0),
			}
			for _, kodePrg := range prgByBidang[kodeBidang] {
				pd := prgData[kodePrg]
				prgResp := programkegiatan.ProgramResponse{
					Kode:      kodePrg,
					Nama:      pd.nama,
					Jenis:     "programs",
					Anggaran:  mkAnggaran(paguPrg[kodePrg]),
					Indikator: getIndikator(kodePrg),
					Kegiatan:  make([]programkegiatan.KegiatanResponse, 0),
				}
				for _, kodeKeg := range kegByPrg[kodePrg] {
					kd := kegData[kodeKeg]
					kegResp := programkegiatan.KegiatanResponse{
						Kode:        kodeKeg,
						Nama:        kd.nama,
						Jenis:       "kegiatans",
						Anggaran:    mkAnggaran(paguKeg[kodeKeg]),
						Indikator:   getIndikator(kodeKeg),
						SubKegiatan: make([]programkegiatan.SubKegiatanResponse, 0),
					}
					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]
						kegResp.SubKegiatan = append(kegResp.SubKegiatan, programkegiatan.SubKegiatanResponse{
							Kode:        kodeSubkeg,
							Nama:        sd.nama,
							Jenis:       "subkegiatans",
							Tahun:       tahun,
							PegawaiId:   sd.pegawaiId,
							NamaPegawai: sd.namaPegawai,
							Anggaran:    mkAnggaran(paguSubkeg[kodeSubkeg]),
							Indikator:   getIndikator(kodeSubkeg),
						})
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

func (service *MatrixRenjaServiceImpl) GetRenjaRankhir(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := service.MatrixRenjaRepository.GetRenjaRankhir(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	result := service.transformToResponseRankhir(data, kodeOpd, tahun)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// transformToResponseAkhir mengubah data flat menjadi hierarki terstruktur.
// Optimasi: O(n) single pass — tidak ada N+1 query, tidak ada linear search di dalam loop.
func (service *MatrixRenjaServiceImpl) transformToResponseRankhir(
	data []domain.SubKegiatanQuery,
	kodeOpd string,
	tahun string,
) []programkegiatan.UrusanDetailResponse {

	if len(data) == 0 {
		return []programkegiatan.UrusanDetailResponse{}
	}

	// Helper: buat slice Anggaran untuk satu tahun
	mkAnggaran := func(pagu int64) []programkegiatan.PaguAnggaranTotalResponse {
		return []programkegiatan.PaguAnggaranTotalResponse{
			{Tahun: tahun, PaguAnggaran: pagu},
		}
	}

	// -----------------------------------------------------------------------
	// Struktur bantu internal untuk membangun indikator dengan dedup target
	// -----------------------------------------------------------------------
	type indEntry struct {
		resp      programkegiatan.IndikatorResponse
		targetSet map[string]struct{}
	}

	indikatorByKode := make(map[string]map[string]*indEntry)

	collectIndikator := func(item domain.SubKegiatanQuery) {
		if item.IndikatorId == "" {
			return
		}
		kode := item.IndikatorKode
		if indikatorByKode[kode] == nil {
			indikatorByKode[kode] = make(map[string]*indEntry)
		}
		ent, exists := indikatorByKode[kode][item.IndikatorId]
		if !exists {
			ent = &indEntry{
				resp: programkegiatan.IndikatorResponse{
					Id:           item.IndikatorId,
					Kode:         kode,
					KodeOpd:      kodeOpd,
					Indikator:    item.Indikator,
					Tahun:        tahun,
					PaguAnggaran: 0, // pagu di luar indikator, ada di field Anggaran tiap level
					StatusTarget: false,
					Target:       make([]programkegiatan.TargetResponse, 0),
				},
				targetSet: make(map[string]struct{}),
			}
			indikatorByKode[kode][item.IndikatorId] = ent
		}
		if item.TargetId != "" {
			if _, seen := ent.targetSet[item.TargetId]; !seen {
				ent.targetSet[item.TargetId] = struct{}{}
				ent.resp.Target = append(ent.resp.Target, programkegiatan.TargetResponse{
					Id:     item.TargetId,
					Target: item.Target,
					Satuan: item.Satuan,
				})
				ent.resp.StatusTarget = true
			}
		}
	}

	getIndikator := func(kode string) []programkegiatan.IndikatorResponse {
		m, ok := indikatorByKode[kode]
		if !ok {
			return []programkegiatan.IndikatorResponse{}
		}
		slice := make([]programkegiatan.IndikatorResponse, 0, len(m))
		for _, ent := range m {
			slice = append(slice, ent.resp)
		}
		return slice
	}

	// -----------------------------------------------------------------------
	// Metadata tiap level hierarki
	// -----------------------------------------------------------------------
	type subkegMeta struct{ nama, kodeKeg string }
	type kegMeta struct{ nama, kodePrg string }
	type prgMeta struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }

	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)

	totalAnggaranBySubkeg := make(map[string]int64)

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
	// PASS 1: satu kali iterasi, kumpulkan semua data
	// -----------------------------------------------------------------------
	for _, item := range data {
		collectIndikator(item)

		if item.KodeSubKegiatan == "" {
			continue
		}

		if _, ok := seenSubkeg[item.KodeSubKegiatan]; !ok {
			seenSubkeg[item.KodeSubKegiatan] = struct{}{}
			subkegData[item.KodeSubKegiatan] = subkegMeta{
				nama:    item.NamaSubKegiatan,
				kodeKeg: item.KodeKegiatan,
			}
			totalAnggaranBySubkeg[item.KodeSubKegiatan] = item.TotalAnggaranSubKegiatan
			if item.KodeKegiatan != "" {
				subkegByKeg[item.KodeKegiatan] = append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
			}
		}

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

	// -----------------------------------------------------------------------
	// Agregasi pagu bottom-up: subkeg → keg → prg → bidang → urusan
	// O(N) total — setiap node diproses tepat satu kali
	// -----------------------------------------------------------------------
	paguKeg := make(map[string]int64)
	paguPrg := make(map[string]int64)
	paguBidang := make(map[string]int64)
	paguUrusan := make(map[string]int64)

	for kodeKeg, subkegList := range subkegByKeg {
		for _, ks := range subkegList {
			paguKeg[kodeKeg] += totalAnggaranBySubkeg[ks]
		}
	}
	for kodePrg, kegList := range kegByPrg {
		for _, kk := range kegList {
			paguPrg[kodePrg] += paguKeg[kk]
		}
	}
	for kodeBidang, prgList := range prgByBidang {
		for _, kp := range prgList {
			paguBidang[kodeBidang] += paguPrg[kp]
		}
	}
	for kodeUrusan, bidangList := range bidangByUrusan {
		for _, kb := range bidangList {
			paguUrusan[kodeUrusan] += paguBidang[kb]
		}
	}
	var paguTotal int64
	for _, p := range paguUrusan {
		paguTotal += p
	}

	// -----------------------------------------------------------------------
	// PASS 2: bangun response hierarki dari maps yang sudah terkumpul
	// -----------------------------------------------------------------------
	urusanDetail := programkegiatan.UrusanDetailResponse{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
		PaguAnggaranTotal: []programkegiatan.PaguAnggaranTotalResponse{
			{Tahun: tahun, PaguAnggaran: paguTotal},
		},
		Urusan: make([]programkegiatan.UrusanResponse, 0, len(urusanOrder)),
	}

	for _, kodeUrusan := range urusanOrder {
		urusanResp := programkegiatan.UrusanResponse{
			Kode:         kodeUrusan,
			Nama:         urusanData[kodeUrusan],
			Jenis:        "urusans",
			Anggaran:     mkAnggaran(paguUrusan[kodeUrusan]),
			Indikator:    getIndikator(kodeUrusan),
			BidangUrusan: make([]programkegiatan.BidangUrusanResponse, 0),
		}

		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			bd := bidangData[kodeBidang]
			bidangResp := programkegiatan.BidangUrusanResponse{
				Kode:      kodeBidang,
				Nama:      bd.nama,
				Jenis:     "bidang_urusans",
				Anggaran:  mkAnggaran(paguBidang[kodeBidang]),
				Indikator: getIndikator(kodeBidang),
				Program:   make([]programkegiatan.ProgramResponse, 0),
			}

			for _, kodePrg := range prgByBidang[kodeBidang] {
				pd := prgData[kodePrg]
				prgResp := programkegiatan.ProgramResponse{
					Kode:      kodePrg,
					Nama:      pd.nama,
					Jenis:     "programs",
					Anggaran:  mkAnggaran(paguPrg[kodePrg]),
					Indikator: getIndikator(kodePrg),
					Kegiatan:  make([]programkegiatan.KegiatanResponse, 0),
				}

				for _, kodeKeg := range kegByPrg[kodePrg] {
					kd := kegData[kodeKeg]
					kegResp := programkegiatan.KegiatanResponse{
						Kode:        kodeKeg,
						Nama:        kd.nama,
						Jenis:       "kegiatans",
						Anggaran:    mkAnggaran(paguKeg[kodeKeg]),
						Indikator:   getIndikator(kodeKeg),
						SubKegiatan: make([]programkegiatan.SubKegiatanResponse, 0),
					}

					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]
						kegResp.SubKegiatan = append(kegResp.SubKegiatan, programkegiatan.SubKegiatanResponse{
							Kode:      kodeSubkeg,
							Nama:      sd.nama,
							Jenis:     "subkegiatans",
							Tahun:     tahun,
							Anggaran:  mkAnggaran(totalAnggaranBySubkeg[kodeSubkeg]),
							Indikator: getIndikator(kodeSubkeg),
						})
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

func (service *MatrixRenjaServiceImpl) CreateOrUpdateTarget(ctx context.Context, request programkegiatan.TargetRenjaRequest) (programkegiatan.TargetResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.TargetResponse{}, err
	}
	defer tx.Rollback()

	var targetId string

	if request.Id == "" {
		// CREATE: generate ID baru
		randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		targetId = fmt.Sprintf("TRG-%s-%s", strings.ToUpper(request.Jenis), randomDigits)

		target := domain.Target{
			Id:          targetId,
			IndikatorId: request.IndikatorId,
			Target:      request.Target,
			Satuan:      request.Satuan,
			Jenis:       request.Jenis,
		}
		err = service.MatrixRenjaRepository.SaveTargetRenja(ctx, tx, target)
	} else {
		// UPDATE: gunakan ID yang dikirim FE
		targetId = request.Id
		target := domain.Target{
			Id:          request.Id,
			IndikatorId: request.IndikatorId,
			Target:      request.Target,
			Satuan:      request.Satuan,
			Jenis:       request.Jenis,
		}
		err = service.MatrixRenjaRepository.UpdateTargetRenja(ctx, tx, target)
	}

	if err != nil {
		return programkegiatan.TargetResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return programkegiatan.TargetResponse{}, err
	}

	return programkegiatan.TargetResponse{
		Id:     targetId,
		Target: request.Target,
		Satuan: request.Satuan,
	}, nil
}
