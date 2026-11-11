package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/model/web/renja"
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

func (service *MatrixRenjaServiceImpl) GetRenja(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.RanwalMatrixResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// Service yang menentukan jenis indikator dan jenis pagu
	data, err := service.MatrixRenjaRepository.GetRenja(ctx, tx, kodeOpd, tahun, jenisPagu)
	if err != nil {
		return nil, err
	}
	result := service.transformToResponse(data, kodeOpd, tahun)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
func (service *MatrixRenjaServiceImpl) GetRenjaRankhir(ctx context.Context, kodeOpd, tahun string) ([]renja.RankhirMatrixResponse, error) {
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

// transformToResponse: unified untuk ranwal & rankhir.
// Ranwal: pagu di PaguSubKegiatan | Rankhir: pagu di TotalAnggaranSubKegiatan
// Cukup penjumlahan keduanya karena hanya satu yang terisi per item.
func (service *MatrixRenjaServiceImpl) transformToResponse(data []domain.SubKegiatanQuery, kodeOpd, tahun string) []renja.UrusanDetailResponse {
	if len(data) == 0 {
		return []renja.UrusanDetailResponse{}
	}

	// ── Indikator collector dengan dedup target ──
	type indEntry struct {
		resp      renja.IndikatorMatrixResponse
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
				resp: renja.IndikatorMatrixResponse{
					KodeIndikator: item.IndikatorId,
					Kode:          kode,
					KodeOpd:       kodeOpd,
					Indikator:     item.Indikator,
					Tahun:         item.IndikatorTahun,
					StatusTarget:  false,
					Target:        item.Target,
					Satuan:        item.Satuan,
				},
				targetSet: make(map[string]struct{}),
			}

			indikatorByKode[kode][item.IndikatorId] = ent
		}

		if item.TargetId != "" {
			if _, seen := ent.targetSet[item.TargetId]; !seen {
				ent.targetSet[item.TargetId] = struct{}{}
				ent.resp.TargetId = item.TargetId
				ent.resp.Target = item.Target
				ent.resp.Satuan = item.Satuan
				ent.resp.StatusTarget = true
			}
		}
	}

	getIndikator := func(kode string) []renja.IndikatorMatrixResponse {
		m, ok := indikatorByKode[kode]

		if !ok {
			return []renja.IndikatorMatrixResponse{}
		}

		out := make([]renja.IndikatorMatrixResponse, 0, len(m))

		for _, e := range m {
			out = append(out, e.resp)
		}

		return out
	}

	// ── Metadata ──
	type subkegMeta struct {
		nama, kodeKeg, pegawaiId, namaPegawai string
	}

	type kegMeta struct {
		nama, kodePrg string
	}

	type prgMeta struct {
		nama, kodeBidang string
	}

	type bidangMeta struct {
		nama, kodeUrusan string
	}

	subkegData := make(map[string]subkegMeta)
	kegData := make(map[string]kegMeta)
	prgData := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)

	// pagu & jenis pagu
	paguSubkeg := make(map[string]int64)
	jenisPaguSubkeg := make(map[string]string)

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

	// ── PASS 1: single-pass collect ──
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

			// Ranwal: gunakan PaguRanwal flat (dari tb_pagu jenis='ranwal')
			paguSubkeg[item.KodeSubKegiatan] = item.PaguRanwal
			jenisPaguSubkeg[item.KodeSubKegiatan] = "ranwal"

			subkegByKeg[item.KodeKegiatan] =
				append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
		}

		if item.KodeKegiatan != "" {
			if _, ok := seenKeg[item.KodeKegiatan]; !ok {

				seenKeg[item.KodeKegiatan] = struct{}{}

				kegData[item.KodeKegiatan] = kegMeta{
					nama:    item.NamaKegiatan,
					kodePrg: item.KodeProgram,
				}

				kegByPrg[item.KodeProgram] =
					append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
			}
		}

		if item.KodeProgram != "" {
			if _, ok := seenPrg[item.KodeProgram]; !ok {

				seenPrg[item.KodeProgram] = struct{}{}

				prgData[item.KodeProgram] = prgMeta{
					nama:       item.NamaProgram,
					kodeBidang: item.KodeBidangUrusan,
				}

				prgByBidang[item.KodeBidangUrusan] =
					append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
			}
		}

		if item.KodeBidangUrusan != "" {
			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {

				seenBidang[item.KodeBidangUrusan] = struct{}{}

				bidangData[item.KodeBidangUrusan] = bidangMeta{
					nama:       item.NamaBidangUrusan,
					kodeUrusan: item.KodeUrusan,
				}

				bidangByUrusan[item.KodeUrusan] =
					append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
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

	// ── Agregasi pagu bottom-up ──
	paguKeg := make(map[string]int64)
	paguPrg := make(map[string]int64)
	paguBidang := make(map[string]int64)
	paguUrusan := make(map[string]int64)

	for kodeKeg, list := range subkegByKeg {
		for _, ks := range list {
			paguKeg[kodeKeg] += paguSubkeg[ks]
		}
	}

	for kodePrg, list := range kegByPrg {
		for _, kk := range list {
			paguPrg[kodePrg] += paguKeg[kk]
		}
	}

	for kodeBidang, list := range prgByBidang {
		for _, kp := range list {
			paguBidang[kodeBidang] += paguPrg[kp]
		}
	}

	for kodeUrusan, list := range bidangByUrusan {
		for _, kb := range list {
			paguUrusan[kodeUrusan] += paguBidang[kb]
		}
	}

	var paguTotal int64

	for _, p := range paguUrusan {
		paguTotal += p
	}

	// ── PASS 2: bangun hierarki ──
	urusanDetail := renja.UrusanDetailResponse{
		KodeOpd:    kodeOpd,
		Tahun:      tahun,
		PaguRanwal: paguTotal,
		Urusan:     make([]renja.UrusanResponse, 0),
	}

	for _, kodeUrusan := range urusanOrder {

		urusanResp := renja.UrusanResponse{
			Kode:         kodeUrusan,
			Nama:         urusanData[kodeUrusan],
			Jenis:        "urusans",
			PaguRanwal:   paguUrusan[kodeUrusan],
			Indikator:    getIndikator(kodeUrusan),
			BidangUrusan: make([]renja.BidangUrusanResponse, 0),
		}

		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {

			bd := bidangData[kodeBidang]

			bidangResp := renja.BidangUrusanResponse{
				Kode:       kodeBidang,
				Nama:       bd.nama,
				Jenis:      "bidang_urusans",
				PaguRanwal: paguBidang[kodeBidang],
				Indikator:  getIndikator(kodeBidang),
				Program:    make([]renja.ProgramResponse, 0),
			}

			for _, kodePrg := range prgByBidang[kodeBidang] {

				pd := prgData[kodePrg]

				prgResp := renja.ProgramResponse{
					Kode:       kodePrg,
					Nama:       pd.nama,
					Jenis:      "programs",
					PaguRanwal: paguPrg[kodePrg],
					Indikator:  getIndikator(kodePrg),
					Kegiatan:   make([]renja.KegiatanResponse, 0),
				}

				for _, kodeKeg := range kegByPrg[kodePrg] {

					kd := kegData[kodeKeg]

					kegResp := renja.KegiatanResponse{
						Kode:        kodeKeg,
						Nama:        kd.nama,
						Jenis:       "kegiatans",
						PaguRanwal:  paguKeg[kodeKeg],
						Indikator:   getIndikator(kodeKeg),
						SubKegiatan: make([]renja.SubKegiatanResponse, 0),
					}

					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {

						sd := subkegData[kodeSubkeg]

						kegResp.SubKegiatan = append(
							kegResp.SubKegiatan,
							renja.SubKegiatanResponse{
								Kode:        kodeSubkeg,
								Nama:        sd.nama,
								Jenis:       "subkegiatans",
								Tahun:       tahun,
								PegawaiId:   sd.pegawaiId,
								NamaPegawai: sd.namaPegawai,
								PaguRanwal:  paguSubkeg[kodeSubkeg],
								Indikator:   getIndikator(kodeSubkeg),
							},
						)
					}

					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
				}

				bidangResp.Program = append(bidangResp.Program, prgResp)
			}

			urusanResp.BidangUrusan =
				append(urusanResp.BidangUrusan, bidangResp)
		}

		urusanDetail.Urusan =
			append(urusanDetail.Urusan, urusanResp)
	}

	return []renja.UrusanDetailResponse{urusanDetail}
}

// transformToResponseAkhir mengubah data flat menjadi hierarki terstruktur.
// Optimasi: O(n) single pass — tidak ada N+1 query, tidak ada linear search di dalam loop.
// func (service *MatrixRenjaServiceImpl) transformToResponseRankhir(data []domain.SubKegiatanQuery,kodeOpd string,tahun string,) []renja.UrusanDetailResponse {

// 	if len(data) == 0 {
// 		return []renja.UrusanDetailResponse{}
// 	}

// 	// Helper: buat slice Anggaran untuk satu tahun
// 	mkAnggaran := func(pagu int64) []programkegiatan.PaguAnggaranTotalResponse {
// 		return []programkegiatan.PaguAnggaranTotalResponse{
// 			{Tahun: tahun, PaguAnggaran: pagu},
// 		}
// 	}

// 	// -----------------------------------------------------------------------
// 	// Struktur bantu internal untuk membangun indikator dengan dedup target
// 	// -----------------------------------------------------------------------
// 	type indEntry struct {
// 		resp      programkegiatan.IndikatorResponse
// 		targetSet map[string]struct{}
// 	}

// 	indikatorByKode := make(map[string]map[string]*indEntry)

// 	collectIndikator := func(item domain.SubKegiatanQuery) {
// 		if item.IndikatorId == "" {
// 			return
// 		}
// 		kode := item.IndikatorKode
// 		if indikatorByKode[kode] == nil {
// 			indikatorByKode[kode] = make(map[string]*indEntry)
// 		}
// 		ent, exists := indikatorByKode[kode][item.IndikatorId]
// 		if !exists {
// 			ent = &indEntry{
// 				resp: programkegiatan.IndikatorResponse{
// 					Id:           item.IndikatorId,
// 					Kode:         kode,
// 					KodeOpd:      kodeOpd,
// 					Indikator:    item.Indikator,
// 					Tahun:        tahun,
// 					PaguAnggaran: 0, // pagu di luar indikator, ada di field Anggaran tiap level
// 					StatusTarget: false,
// 					Target:       make([]programkegiatan.TargetResponse, 0),
// 				},
// 				targetSet: make(map[string]struct{}),
// 			}
// 			indikatorByKode[kode][item.IndikatorId] = ent
// 		}
// 		if item.TargetId != "" {
// 			if _, seen := ent.targetSet[item.TargetId]; !seen {
// 				ent.targetSet[item.TargetId] = struct{}{}
// 				ent.resp.Target = append(ent.resp.Target, programkegiatan.TargetResponse{
// 					Id:     item.TargetId,
// 					Target: item.Target,
// 					Satuan: item.Satuan,
// 				})
// 				ent.resp.StatusTarget = true
// 			}
// 		}
// 	}

// 	getIndikator := func(kode string) []programkegiatan.IndikatorResponse {
// 		m, ok := indikatorByKode[kode]
// 		if !ok {
// 			return []programkegiatan.IndikatorResponse{}
// 		}
// 		slice := make([]programkegiatan.IndikatorResponse, 0, len(m))
// 		for _, ent := range m {
// 			slice = append(slice, ent.resp)
// 		}
// 		return slice
// 	}

// 	// -----------------------------------------------------------------------
// 	// Metadata tiap level hierarki
// 	// -----------------------------------------------------------------------
// 	type subkegMeta struct{ nama, kodeKeg string }
// 	type kegMeta struct{ nama, kodePrg string }
// 	type prgMeta struct{ nama, kodeBidang string }
// 	type bidangMeta struct{ nama, kodeUrusan string }

// 	subkegData := make(map[string]subkegMeta)
// 	kegData := make(map[string]kegMeta)
// 	prgData := make(map[string]prgMeta)
// 	bidangData := make(map[string]bidangMeta)
// 	urusanData := make(map[string]string)

// 	totalAnggaranBySubkeg := make(map[string]int64)

// 	seenSubkeg := make(map[string]struct{})
// 	seenKeg := make(map[string]struct{})
// 	seenPrg := make(map[string]struct{})
// 	seenBidang := make(map[string]struct{})
// 	seenUrusan := make(map[string]struct{})

// 	subkegByKeg := make(map[string][]string)
// 	kegByPrg := make(map[string][]string)
// 	prgByBidang := make(map[string][]string)
// 	bidangByUrusan := make(map[string][]string)
// 	var urusanOrder []string

// 	// -----------------------------------------------------------------------
// 	// PASS 1: satu kali iterasi, kumpulkan semua data
// 	// -----------------------------------------------------------------------
// 	for _, item := range data {
// 		collectIndikator(item)

// 		if item.KodeSubKegiatan == "" {
// 			continue
// 		}

// 		if _, ok := seenSubkeg[item.KodeSubKegiatan]; !ok {
// 			seenSubkeg[item.KodeSubKegiatan] = struct{}{}
// 			subkegData[item.KodeSubKegiatan] = subkegMeta{
// 				nama:    item.NamaSubKegiatan,
// 				kodeKeg: item.KodeKegiatan,
// 			}
// 			totalAnggaranBySubkeg[item.KodeSubKegiatan] = item.TotalAnggaranSubKegiatan
// 			if item.KodeKegiatan != "" {
// 				subkegByKeg[item.KodeKegiatan] = append(subkegByKeg[item.KodeKegiatan], item.KodeSubKegiatan)
// 			}
// 		}

// 		if item.KodeKegiatan != "" {
// 			if _, ok := seenKeg[item.KodeKegiatan]; !ok {
// 				seenKeg[item.KodeKegiatan] = struct{}{}
// 				kegData[item.KodeKegiatan] = kegMeta{nama: item.NamaKegiatan, kodePrg: item.KodeProgram}
// 				if item.KodeProgram != "" {
// 					kegByPrg[item.KodeProgram] = append(kegByPrg[item.KodeProgram], item.KodeKegiatan)
// 				}
// 			}
// 		}

// 		if item.KodeProgram != "" {
// 			if _, ok := seenPrg[item.KodeProgram]; !ok {
// 				seenPrg[item.KodeProgram] = struct{}{}
// 				prgData[item.KodeProgram] = prgMeta{nama: item.NamaProgram, kodeBidang: item.KodeBidangUrusan}
// 				if item.KodeBidangUrusan != "" {
// 					prgByBidang[item.KodeBidangUrusan] = append(prgByBidang[item.KodeBidangUrusan], item.KodeProgram)
// 				}
// 			}
// 		}

// 		if item.KodeBidangUrusan != "" {
// 			if _, ok := seenBidang[item.KodeBidangUrusan]; !ok {
// 				seenBidang[item.KodeBidangUrusan] = struct{}{}
// 				bidangData[item.KodeBidangUrusan] = bidangMeta{nama: item.NamaBidangUrusan, kodeUrusan: item.KodeUrusan}
// 				if item.KodeUrusan != "" {
// 					bidangByUrusan[item.KodeUrusan] = append(bidangByUrusan[item.KodeUrusan], item.KodeBidangUrusan)
// 				}
// 			}
// 		}

// 		if item.KodeUrusan != "" {
// 			if _, ok := seenUrusan[item.KodeUrusan]; !ok {
// 				seenUrusan[item.KodeUrusan] = struct{}{}
// 				urusanData[item.KodeUrusan] = item.NamaUrusan
// 				urusanOrder = append(urusanOrder, item.KodeUrusan)
// 			}
// 		}
// 	}

// 	// -----------------------------------------------------------------------
// 	// Agregasi pagu bottom-up: subkeg → keg → prg → bidang → urusan
// 	// O(N) total — setiap node diproses tepat satu kali
// 	// -----------------------------------------------------------------------
// 	paguKeg := make(map[string]int64)
// 	paguPrg := make(map[string]int64)
// 	paguBidang := make(map[string]int64)
// 	paguUrusan := make(map[string]int64)

// 	for kodeKeg, subkegList := range subkegByKeg {
// 		for _, ks := range subkegList {
// 			paguKeg[kodeKeg] += totalAnggaranBySubkeg[ks]
// 		}
// 	}
// 	for kodePrg, kegList := range kegByPrg {
// 		for _, kk := range kegList {
// 			paguPrg[kodePrg] += paguKeg[kk]
// 		}
// 	}
// 	for kodeBidang, prgList := range prgByBidang {
// 		for _, kp := range prgList {
// 			paguBidang[kodeBidang] += paguPrg[kp]
// 		}
// 	}
// 	for kodeUrusan, bidangList := range bidangByUrusan {
// 		for _, kb := range bidangList {
// 			paguUrusan[kodeUrusan] += paguBidang[kb]
// 		}
// 	}
// 	var paguTotal int64
// 	for _, p := range paguUrusan {
// 		paguTotal += p
// 	}

// 	// -----------------------------------------------------------------------
// 	// PASS 2: bangun response hierarki dari maps yang sudah terkumpul
// 	// -----------------------------------------------------------------------
// 	urusanDetail := renja.UrusanDetailResponse{
// 		KodeOpd: kodeOpd,
// 		Tahun:   tahun,
// 		PaguAnggaranTotal: []programkegiatan.PaguAnggaranTotalResponse{
// 			{Tahun: tahun, PaguAnggaran: paguTotal},
// 		},
// 		Urusan: make([]renja.UrusanResponse, 0, len(urusanOrder)),
// 	}

// 	for _, kodeUrusan := range urusanOrder {
// 		urusanResp := renja.UrusanResponse{
// 			Kode:         kodeUrusan,
// 			Nama:         urusanData[kodeUrusan],
// 			Jenis:        "urusans",
// 			Anggaran:     mkAnggaran(paguUrusan[kodeUrusan]),
// 			Indikator:    getIndikator(kodeUrusan),
// 			BidangUrusan: make([]renja.BidangUrusanResponse, 0),
// 		}

// 		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
// 			bd := bidangData[kodeBidang]
// 			bidangResp := renja.BidangUrusanResponse{
// 				Kode:      kodeBidang,
// 				Nama:      bd.nama,
// 				Jenis:     "bidang_urusans",
// 				Anggaran:  mkAnggaran(paguBidang[kodeBidang]),
// 				Indikator: getIndikator(kodeBidang),
// 				Program:   make([]renja.ProgramResponse, 0),
// 			}

// 			for _, kodePrg := range prgByBidang[kodeBidang] {
// 				pd := prgData[kodePrg]
// 				prgResp := renja.ProgramResponse{
// 					Kode:      kodePrg,
// 					Nama:      pd.nama,
// 					Jenis:     "programs",
// 					Anggaran:  mkAnggaran(paguPrg[kodePrg]),
// 					Indikator: getIndikator(kodePrg),
// 					Kegiatan:  make([]renja.KegiatanResponse, 0),
// 				}

// 				for _, kodeKeg := range kegByPrg[kodePrg] {
// 					kd := kegData[kodeKeg]
// 					kegResp := renja.KegiatanResponse{
// 						Kode:        kodeKeg,
// 						Nama:        kd.nama,
// 						Jenis:       "kegiatans",
// 						Anggaran:    mkAnggaran(paguKeg[kodeKeg]),
// 						Indikator:   getIndikator(kodeKeg),
// 						SubKegiatan: make([]renja.SubKegiatanResponse, 0),
// 					}

// 					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
// 						sd := subkegData[kodeSubkeg]
// 						kegResp.SubKegiatan = append(kegResp.SubKegiatan, renja.SubKegiatanResponse{
// 							Kode:      kodeSubkeg,
// 							Nama:      sd.nama,
// 							Jenis:     "subkegiatans",
// 							Tahun:     tahun,
// 							Anggaran:  mkAnggaran(totalAnggaranBySubkeg[kodeSubkeg]),
// 							Indikator: getIndikator(kodeSubkeg),
// 						})
// 					}

// 					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
// 				}

// 				bidangResp.Program = append(bidangResp.Program, prgResp)
// 			}

// 			urusanResp.BidangUrusan = append(urusanResp.BidangUrusan, bidangResp)
// 		}

// 		urusanDetail.Urusan = append(urusanDetail.Urusan, urusanResp)
// 	}

// 	return []renja.UrusanDetailResponse{urusanDetail}
// }

func (service *MatrixRenjaServiceImpl) UpsertBatchIndikatorRenja(ctx context.Context, requests []programkegiatan.IndikatorRenjaCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	prefixBase := fmt.Sprintf("RENJA-RANKHIR-%s-%s-%s-", requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun)
	existingCount, err := service.MatrixRenjaRepository.CountIndikatorMatrixByScope(
		ctx, tx,
		requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun, requests[0].Jenis,
		prefixBase,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung urutan indikator: %w", err)
	}
	nextUrutan := existingCount + 1
	var (
		domainItems []domain.Indikator
		respItems   []programkegiatan.IndikatorUpsertResponse
		keepList    []string // ← TAMBAH INI
	)
	for _, item := range requests {
		var kodeIndikator string
		var targetId string
		if item.KodeIndikator == "" {
			kodeIndikator = fmt.Sprintf("%s%03d", prefixBase, nextUrutan)
			nextUrutan++
			targetId = fmt.Sprintf("TRG-%s-%05d", strings.ToUpper(requests[0].Jenis), uuid.New().ID()%100000)
		} else {
			kodeIndikator = item.KodeIndikator
			targetId = fmt.Sprintf("TRG-%s-%05d", strings.ToUpper(requests[0].Jenis), uuid.New().ID()%100000)
		}
		keepList = append(keepList, kodeIndikator) // ← TAMBAH INI
		domainItems = append(domainItems, domain.Indikator{
			KodeIndikator: kodeIndikator,
			Kode:          requests[0].Kode,
			KodeOpd:       requests[0].KodeOpd,
			Indikator:     item.Indikator,
			Tahun:         requests[0].Tahun,
			Jenis:         requests[0].Jenis,
			Target: []domain.Target{{
				Id:          targetId,
				IndikatorId: kodeIndikator,
				Tahun:       requests[0].Tahun,
				Target:      item.Target,
				Satuan:      item.Satuan,
				Jenis:       requests[0].Jenis,
			}},
		})
		respItems = append(respItems, programkegiatan.IndikatorUpsertResponse{
			KodeIndikator: kodeIndikator,
			Kode:          requests[0].Kode,
			KodeOpd:       requests[0].KodeOpd,
			Indikator:     item.Indikator,
			Jenis:         requests[0].Jenis,
			Tahun:         requests[0].Tahun,
			Target:        item.Target,
			Satuan:        item.Satuan,
		})
	}
	if err := service.MatrixRenjaRepository.UpsertBatchIndikatorRenja(ctx, tx, domainItems); err != nil {
		return nil, err
	}
	// ── SYNC: hapus indikator di DB yang tidak ada di request ──  ← TAMBAH INI
	if err := service.MatrixRenjaRepository.DeleteIndicatorsExcept(
		ctx, tx,
		requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun, requests[0].Jenis,
		keepList,
	); err != nil {
		return nil, fmt.Errorf("gagal menghapus indikator: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return respItems, nil
}

func (service *MatrixRenjaServiceImpl) UpsertBatchIndikatorRenjaPenetapan(ctx context.Context, requests []programkegiatan.IndikatorRenjaCreateRequest) ([]programkegiatan.IndikatorUpsertResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	prefixBase := fmt.Sprintf("RENJA-PENETAPAN-%s-%s-%s-", requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun)
	existingCount, err := service.MatrixRenjaRepository.CountIndikatorMatrixByScope(
		ctx, tx,
		requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun, requests[0].Jenis,
		prefixBase,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung urutan indikator: %w", err)
	}
	nextUrutan := existingCount + 1
	var (
		domainItems []domain.Indikator
		respItems   []programkegiatan.IndikatorUpsertResponse
		keepList    []string // ← TAMBAH INI
	)
	for _, item := range requests {
		var kodeIndikator string
		var targetId string
		if item.KodeIndikator == "" {
			kodeIndikator = fmt.Sprintf("%s%03d", prefixBase, nextUrutan)
			nextUrutan++
			targetId = fmt.Sprintf("TRG-%s-%05d", strings.ToUpper(requests[0].Jenis), uuid.New().ID()%100000)
		} else {
			kodeIndikator = item.KodeIndikator
			targetId = fmt.Sprintf("TRG-%s-%05d", strings.ToUpper(requests[0].Jenis), uuid.New().ID()%100000)
		}
		keepList = append(keepList, kodeIndikator) // ← TAMBAH INI
		domainItems = append(domainItems, domain.Indikator{
			KodeIndikator: kodeIndikator,
			Kode:          requests[0].Kode,
			KodeOpd:       requests[0].KodeOpd,
			Indikator:     item.Indikator,
			Tahun:         requests[0].Tahun,
			Jenis:         requests[0].Jenis,
			Target: []domain.Target{{
				Id:          targetId,
				IndikatorId: kodeIndikator,
				Tahun:       requests[0].Tahun,
				Target:      item.Target,
				Satuan:      item.Satuan,
				Jenis:       requests[0].Jenis,
			}},
		})
		respItems = append(respItems, programkegiatan.IndikatorUpsertResponse{
			KodeIndikator: kodeIndikator,
			Kode:          requests[0].Kode,
			KodeOpd:       requests[0].KodeOpd,
			Indikator:     item.Indikator,
			Jenis:         requests[0].Jenis,
			Tahun:         requests[0].Tahun,
			Target:        item.Target,
			Satuan:        item.Satuan,
		})
	}
	if err := service.MatrixRenjaRepository.UpsertBatchIndikatorRenja(ctx, tx, domainItems); err != nil {
		return nil, err
	}
	// ── SYNC: hapus indikator di DB yang tidak ada di request ──  ← TAMBAH INI
	if err := service.MatrixRenjaRepository.DeleteIndicatorsExcept(
		ctx, tx,
		requests[0].Kode, requests[0].KodeOpd, requests[0].Tahun, requests[0].Jenis,
		keepList,
	); err != nil {
		return nil, fmt.Errorf("gagal menghapus indikator: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return respItems, nil
}

func (service *MatrixRenjaServiceImpl) UpsertAnggaran(ctx context.Context, request programkegiatan.AnggaranRenjaRequest) (programkegiatan.AnggaranRenjaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.AnggaranRenjaResponse{}, err
	}
	defer tx.Rollback()
	err = service.MatrixRenjaRepository.UpsertAnggaran(
		ctx, tx,
		request.KodeSubKegiatan,
		request.KodeOpd,
		request.Tahun,
		request.Pagu,
	)
	if err != nil {
		return programkegiatan.AnggaranRenjaResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return programkegiatan.AnggaranRenjaResponse{}, err
	}
	return programkegiatan.AnggaranRenjaResponse{
		KodeSubKegiatan: request.KodeSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
		Pagu:            request.Pagu,
	}, nil
}

func (service *MatrixRenjaServiceImpl) GetRenjaPenetapan(ctx context.Context, kodeOpd, tahun, jenisPagu string) ([]renja.PenetapanMatrixResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := service.MatrixRenjaRepository.GetRenjaPenetapan(ctx, tx, kodeOpd, tahun, jenisPagu)
	if err != nil {
		return nil, err
	}
	result := service.transformToResponsePenetapan(data, kodeOpd, tahun)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

// transformToResponseRankhir — rankhir view.
// Indikator ditampilkan DUAL: IndikatorRanwal (target_ranwal) + IndikatorRankhir (target_rankhir).
// Pagu total ranwal: semua subkegiatan (termasuk sintetis dari renstra).
// Pagu total rankhir: hanya subkegiatan yang ada di rekin tahun ini (sintetis → 0).
func (service *MatrixRenjaServiceImpl) transformToResponseRankhir(data []domain.SubKegiatanQuery, kodeOpd, tahun string) []renja.UrusanDetailResponse {
	if len(data) == 0 {
		return []renja.UrusanDetailResponse{}
	}

	// ── Dual indikator collector (ranwal vs rankhir) ───────────────────────
	type indEntry struct {
		resp renja.IndikatorMatrixResponse
	}

	indRanwalByKode  := make(map[string]map[string]*indEntry)
	indRankhirByKode := make(map[string]map[string]*indEntry)

	addToMap := func(m map[string]map[string]*indEntry, item domain.SubKegiatanQuery, target, satuan, targetId string) {
		if item.IndikatorId == "" {
			return
		}
		kode := item.IndikatorKode
		if m[kode] == nil {
			m[kode] = make(map[string]*indEntry)
		}
		if _, exists := m[kode][item.IndikatorId]; !exists {
			m[kode][item.IndikatorId] = &indEntry{
				resp: renja.IndikatorMatrixResponse{
					KodeIndikator: item.IndikatorId,
					Kode:          kode,
					KodeOpd:       kodeOpd,
					Indikator:     item.Indikator,
					Tahun:         item.IndikatorTahun,
					StatusTarget:  targetId != "",
					Target:        target,
					Satuan:        satuan,
				},
			}
		}
	}

	collectIndikator := func(item domain.SubKegiatanQuery) {
		addToMap(indRanwalByKode, item, item.TargetRanwal, item.SatuanRanwal, item.TargetRanwalId)
		// Rankhir: hanya isi target jika ada baris tb_target jenis rankhir (ada id)
		rkTarget, rkSatuan := item.TargetRankhir, item.SatuanRankhir
		if item.TargetRankhirId == "" {
			rkTarget, rkSatuan = "", ""
		}
		addToMap(indRankhirByKode, item, rkTarget, rkSatuan, item.TargetRankhirId)
	}

	flatten := func(m map[string]map[string]*indEntry, kode string) []renja.IndikatorMatrixResponse {
		mp, ok := m[kode]
		if !ok {
			return []renja.IndikatorMatrixResponse{}
		}
		out := make([]renja.IndikatorMatrixResponse, 0, len(mp))
		for _, e := range mp {
			out = append(out, e.resp)
		}
		return out
	}
	getIndRanwal  := func(kode string) []renja.IndikatorMatrixResponse { return flatten(indRanwalByKode, kode) }
	getIndRankhir := func(kode string) []renja.IndikatorMatrixResponse { return flatten(indRankhirByKode, kode) }

	// ── Metadata structs ───────────────────────────────────────────────────
	type subkegMeta struct{ nama, kodeKeg, pegawaiId, namaPegawai, status string }
	type kegMeta    struct{ nama, kodePrg string }
	type prgMeta    struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }

	subkegData := make(map[string]subkegMeta)
	kegData    := make(map[string]kegMeta)
	prgData    := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)

	// pagu_ranwal : semua subkegiatan (termasuk sintetis renstra)
	// pagu_rankhir: HANYA subkegiatan yg ada di rekin (sintetis selalu 0 krn tidak ada rincian_belanja)
	paguRanwalSubkeg  := make(map[string]int64)
	paguRankhirSubkeg := make(map[string]int64)

	seenSubkeg  := make(map[string]struct{})
	seenKeg     := make(map[string]struct{})
	seenPrg     := make(map[string]struct{})
	seenBidang  := make(map[string]struct{})
	seenUrusan  := make(map[string]struct{})

	subkegByKeg    := make(map[string][]string)
	kegByPrg       := make(map[string][]string)
	prgByBidang    := make(map[string][]string)
	bidangByUrusan := make(map[string][]string)
	var urusanOrder []string

	// ── PASS 1: collect ────────────────────────────────────────────────────
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
				status:      item.StatusSubkegiatan,
			}
			// Ranwal: semua subkegiatan dari renstra
			paguRanwalSubkeg[item.KodeSubKegiatan] = item.PaguRanwal
			// Rankhir: hanya subkegiatan yg ada di rekin (sintetis memiliki StatusSubkegiatan != "")
			if item.StatusSubkegiatan == "" {
				paguRankhirSubkeg[item.KodeSubKegiatan] = item.PaguRankhir
			}
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

	// ── Agregasi pagu bottom-up (ranwal + rankhir terpisah) ───────────────
	paguRanwalKeg    := make(map[string]int64)
	paguRankhirKeg   := make(map[string]int64)
	paguRanwalPrg    := make(map[string]int64)
	paguRankhirPrg   := make(map[string]int64)
	paguRanwalBidang  := make(map[string]int64)
	paguRankhirBidang := make(map[string]int64)
	paguRanwalUrusan  := make(map[string]int64)
	paguRankhirUrusan := make(map[string]int64)

	for kk, list := range subkegByKeg {
		for _, ks := range list {
			paguRanwalKeg[kk]  += paguRanwalSubkeg[ks]
			paguRankhirKeg[kk] += paguRankhirSubkeg[ks]
		}
	}
	for kp, list := range kegByPrg {
		for _, kk := range list {
			paguRanwalPrg[kp]  += paguRanwalKeg[kk]
			paguRankhirPrg[kp] += paguRankhirKeg[kk]
		}
	}
	for kb, list := range prgByBidang {
		for _, kp := range list {
			paguRanwalBidang[kb]  += paguRanwalPrg[kp]
			paguRankhirBidang[kb] += paguRankhirPrg[kp]
		}
	}
	for ku, list := range bidangByUrusan {
		for _, kb := range list {
			paguRanwalUrusan[ku]  += paguRanwalBidang[kb]
			paguRankhirUrusan[ku] += paguRankhirBidang[kb]
		}
	}
	var totalRanwal, totalRankhir int64
	for _, p := range paguRanwalUrusan  { totalRanwal  += p }
	for _, p := range paguRankhirUrusan { totalRankhir += p }

	// ── PASS 2: bangun hierarki ────────────────────────────────────────────
	urusanDetail := renja.UrusanDetailResponse{
		KodeOpd:     kodeOpd,
		Tahun:       tahun,
		PaguRanwal:  totalRanwal,
		PaguRankhir: totalRankhir,
		Urusan:      make([]renja.UrusanResponse, 0),
	}

	for _, ku := range urusanOrder {
		urusanResp := renja.UrusanResponse{
			Kode:            ku,
			Nama:            urusanData[ku],
			Jenis:           "urusans",
			PaguRanwal:      paguRanwalUrusan[ku],
			PaguRankhir:     paguRankhirUrusan[ku],
			Indikator:       []renja.IndikatorMatrixResponse{},
			IndikatorRanwal:  getIndRanwal(ku),
			IndikatorRankhir: getIndRankhir(ku),
			BidangUrusan:    make([]renja.BidangUrusanResponse, 0),
		}
		for _, kb := range bidangByUrusan[ku] {
			bd := bidangData[kb]
			bidangResp := renja.BidangUrusanResponse{
				Kode:            kb,
				Nama:            bd.nama,
				Jenis:           "bidang_urusans",
				PaguRanwal:      paguRanwalBidang[kb],
				PaguRankhir:     paguRankhirBidang[kb],
				Indikator:       []renja.IndikatorMatrixResponse{},
				IndikatorRanwal:  getIndRanwal(kb),
				IndikatorRankhir: getIndRankhir(kb),
				Program:         make([]renja.ProgramResponse, 0),
			}
			for _, kp := range prgByBidang[kb] {
				pd := prgData[kp]
				prgResp := renja.ProgramResponse{
					Kode:            kp,
					Nama:            pd.nama,
					Jenis:           "programs",
					PaguRanwal:      paguRanwalPrg[kp],
					PaguRankhir:     paguRankhirPrg[kp],
					Indikator:       []renja.IndikatorMatrixResponse{},
					IndikatorRanwal:  getIndRanwal(kp),
					IndikatorRankhir: getIndRankhir(kp),
					Kegiatan:        make([]renja.KegiatanResponse, 0),
				}
				for _, kk := range kegByPrg[kp] {
					kd := kegData[kk]
					kegResp := renja.KegiatanResponse{
						Kode:            kk,
						Nama:            kd.nama,
						Jenis:           "kegiatans",
						PaguRanwal:      paguRanwalKeg[kk],
						PaguRankhir:     paguRankhirKeg[kk],
						Indikator:       []renja.IndikatorMatrixResponse{},
						IndikatorRanwal:  getIndRanwal(kk),
						IndikatorRankhir: getIndRankhir(kk),
						SubKegiatan:     make([]renja.SubKegiatanResponse, 0),
					}
					for _, ks := range subkegByKeg[kk] {
						sd := subkegData[ks]
						kegResp.SubKegiatan = append(kegResp.SubKegiatan,
							renja.SubKegiatanResponse{
								Kode:            ks,
								Nama:            sd.nama,
								Jenis:           "subkegiatans",
								Tahun:           tahun,
								PegawaiId:       sd.pegawaiId,
								NamaPegawai:     sd.namaPegawai,
								Status:          sd.status,
								PaguRanwal:      paguRanwalSubkeg[ks],
								PaguRankhir:     paguRankhirSubkeg[ks],
								Indikator:       []renja.IndikatorMatrixResponse{},
								IndikatorRanwal:  getIndRanwal(ks),
								IndikatorRankhir: getIndRankhir(ks),
							},
						)
					}
					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
				}
				bidangResp.Program = append(bidangResp.Program, prgResp)
			}
			urusanResp.BidangUrusan = append(urusanResp.BidangUrusan, bidangResp)
		}
		urusanDetail.Urusan = append(urusanDetail.Urusan, urusanResp)
	}
	return []renja.UrusanDetailResponse{urusanDetail}
}

// transformToResponsePenetapan membangun hierarki untuk GetRenjaPenetapan.
//
// Pagu ditampilkan FLAT: PaguRanwal + PaguRankhir per level.
// Indikator ditampilkan DUAL: IndikatorRankhir dan IndikatorPenetapan terpisah.
func (service *MatrixRenjaServiceImpl) transformToResponsePenetapan(data []domain.SubKegiatanQuery, kodeOpd, tahun string) []renja.UrusanDetailResponse {
	if len(data) == 0 {
		return []renja.UrusanDetailResponse{}
	}

	// ── Dual indikator collector (rankhir vs penetapan) ──────────────────────
	type indEntry struct {
		resp      renja.IndikatorMatrixResponse
		targetSet map[string]struct{}
	}

	// indRankhirByKode  : indikator jenis rankhir (atau fallback renstra)
	// indPenetapanByKode: indikator jenis penetapan saja
	indRankhirByKode   := make(map[string]map[string]*indEntry)
	indPenetapanByKode := make(map[string]map[string]*indEntry)

	addToMap := func(m map[string]map[string]*indEntry, item domain.SubKegiatanQuery) {
		if item.IndikatorId == "" {
			return
		}
		kode := item.IndikatorKode
		if m[kode] == nil {
			m[kode] = make(map[string]*indEntry)
		}
		ent, exists := m[kode][item.IndikatorId]
		if !exists {
			ent = &indEntry{
				resp: renja.IndikatorMatrixResponse{
					KodeIndikator: item.IndikatorId,
					Kode:          kode,
					KodeOpd:       kodeOpd,
					Indikator:     item.Indikator,
					Tahun:         item.IndikatorTahun,
					StatusTarget:  false,
					Target:        "",
					Satuan:        "",
				},
				targetSet: make(map[string]struct{}),
			}
			m[kode][item.IndikatorId] = ent
		}
		// Target hanya dari tb_target jenis rankhir / penetapan (harus ada id)
		if item.TargetId != "" {
			if _, seen := ent.targetSet[item.TargetId]; !seen {
				ent.targetSet[item.TargetId] = struct{}{}
				ent.resp.Target = item.Target
				ent.resp.Satuan = item.Satuan
				ent.resp.StatusTarget = true
			}
		}
	}

	collectIndikator := func(item domain.SubKegiatanQuery) {
		if item.JenisIndikator == "penetapan" {
			addToMap(indPenetapanByKode, item)
		} else {
			// "rankhir" atau kosong (baris tanpa indikator): masuk rankhir
			addToMap(indRankhirByKode, item)
		}
	}

	flatten := func(m map[string]map[string]*indEntry, kode string) []renja.IndikatorMatrixResponse {
		mp, ok := m[kode]
		if !ok {
			return []renja.IndikatorMatrixResponse{}
		}
		out := make([]renja.IndikatorMatrixResponse, 0, len(mp))
		for _, e := range mp {
			out = append(out, e.resp)
		}
		return out
	}

	getIndRankhir   := func(kode string) []renja.IndikatorMatrixResponse { return flatten(indRankhirByKode, kode) }
	getIndPenetapan := func(kode string) []renja.IndikatorMatrixResponse { return flatten(indPenetapanByKode, kode) }

	// ── Metadata structs ──────────────────────────────────────────────────────
	type subkegMeta struct{ nama, kodeKeg, pegawaiId, namaPegawai string }
	type kegMeta    struct{ nama, kodePrg string }
	type prgMeta    struct{ nama, kodeBidang string }
	type bidangMeta struct{ nama, kodeUrusan string }

	subkegData := make(map[string]subkegMeta)
	kegData    := make(map[string]kegMeta)
	prgData    := make(map[string]prgMeta)
	bidangData := make(map[string]bidangMeta)
	urusanData := make(map[string]string)

	// Flat pagu per sub kegiatan: rankhir (SUM belanja) + penetapan (tb_pagu jenis='penetapan')
	paguRankhirSubkeg  := make(map[string]int64)
	paguPenetapanSubkeg := make(map[string]int64)

	seenSubkeg := make(map[string]struct{})
	seenKeg    := make(map[string]struct{})
	seenPrg    := make(map[string]struct{})
	seenBidang := make(map[string]struct{})
	seenUrusan := make(map[string]struct{})

	subkegByKeg    := make(map[string][]string)
	kegByPrg       := make(map[string][]string)
	prgByBidang    := make(map[string][]string)
	bidangByUrusan := make(map[string][]string)

	var urusanOrder []string

	// ── PASS 1: single-pass collect ──────────────────────────────────────────
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
			// PaguRankhir  = dari SUM(rincian_belanja)
			// PaguPenetapan = dari tb_pagu jenis='penetapan'
			paguRankhirSubkeg[item.KodeSubKegiatan]  = item.PaguRankhir
			paguPenetapanSubkeg[item.KodeSubKegiatan] = item.PaguPenetapan
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

	// ── Agregasi flat-pagu bottom-up (rankhir + penetapan) ───────────────────
	paguRankhirKeg    := make(map[string]int64)
	paguPenetapanKeg  := make(map[string]int64)
	paguRankhirPrg    := make(map[string]int64)
	paguPenetapanPrg  := make(map[string]int64)
	paguRankhirBidang   := make(map[string]int64)
	paguPenetapanBidang := make(map[string]int64)
	paguRankhirUrusan   := make(map[string]int64)
	paguPenetapanUrusan := make(map[string]int64)

	for kodeKeg, list := range subkegByKeg {
		for _, ks := range list {
			paguRankhirKeg[kodeKeg]   += paguRankhirSubkeg[ks]
			paguPenetapanKeg[kodeKeg] += paguPenetapanSubkeg[ks]
		}
	}
	for kodePrg, list := range kegByPrg {
		for _, kk := range list {
			paguRankhirPrg[kodePrg]   += paguRankhirKeg[kk]
			paguPenetapanPrg[kodePrg] += paguPenetapanKeg[kk]
		}
	}
	for kodeBidang, list := range prgByBidang {
		for _, kp := range list {
			paguRankhirBidang[kodeBidang]   += paguRankhirPrg[kp]
			paguPenetapanBidang[kodeBidang] += paguPenetapanPrg[kp]
		}
	}
	for kodeUrusan, list := range bidangByUrusan {
		for _, kb := range list {
			paguRankhirUrusan[kodeUrusan]   += paguRankhirBidang[kb]
			paguPenetapanUrusan[kodeUrusan] += paguPenetapanBidang[kb]
		}
	}

	var totalRankhir, totalPenetapan int64
	for _, p := range paguRankhirUrusan   { totalRankhir   += p }
	for _, p := range paguPenetapanUrusan { totalPenetapan += p }

	// ── PASS 2: bangun hierarki dengan flat pagu (rankhir+penetapan) + dual indikator ──
	urusanDetail := renja.UrusanDetailResponse{
		KodeOpd:       kodeOpd,
		Tahun:         tahun,
		PaguRankhir:   totalRankhir,
		PaguPenetapan: totalPenetapan,
		Urusan:        make([]renja.UrusanResponse, 0),
	}

	for _, kodeUrusan := range urusanOrder {

		urusanResp := renja.UrusanResponse{
			Kode:               kodeUrusan,
			Nama:               urusanData[kodeUrusan],
			Jenis:              "urusans",
			PaguRankhir:        paguRankhirUrusan[kodeUrusan],
			PaguPenetapan:      paguPenetapanUrusan[kodeUrusan],
			Indikator:          []renja.IndikatorMatrixResponse{},
			IndikatorRankhir:   getIndRankhir(kodeUrusan),
			IndikatorPenetapan: getIndPenetapan(kodeUrusan),
			BidangUrusan:       make([]renja.BidangUrusanResponse, 0),
		}

		for _, kodeBidang := range bidangByUrusan[kodeUrusan] {
			bd := bidangData[kodeBidang]

			bidangResp := renja.BidangUrusanResponse{
				Kode:               kodeBidang,
				Nama:               bd.nama,
				Jenis:              "bidang_urusans",
				PaguRankhir:        paguRankhirBidang[kodeBidang],
				PaguPenetapan:      paguPenetapanBidang[kodeBidang],
				Indikator:          []renja.IndikatorMatrixResponse{},
				IndikatorRankhir:   getIndRankhir(kodeBidang),
				IndikatorPenetapan: getIndPenetapan(kodeBidang),
				Program:            make([]renja.ProgramResponse, 0),
			}

			for _, kodePrg := range prgByBidang[kodeBidang] {
				pd := prgData[kodePrg]

				prgResp := renja.ProgramResponse{
					Kode:               kodePrg,
					Nama:               pd.nama,
					Jenis:              "programs",
					PaguRankhir:        paguRankhirPrg[kodePrg],
					PaguPenetapan:      paguPenetapanPrg[kodePrg],
					Indikator:          []renja.IndikatorMatrixResponse{},
					IndikatorRankhir:   getIndRankhir(kodePrg),
					IndikatorPenetapan: getIndPenetapan(kodePrg),
					Kegiatan:           make([]renja.KegiatanResponse, 0),
				}

				for _, kodeKeg := range kegByPrg[kodePrg] {
					kd := kegData[kodeKeg]

					kegResp := renja.KegiatanResponse{
						Kode:               kodeKeg,
						Nama:               kd.nama,
						Jenis:              "kegiatans",
						PaguRankhir:        paguRankhirKeg[kodeKeg],
						PaguPenetapan:      paguPenetapanKeg[kodeKeg],
						Indikator:          []renja.IndikatorMatrixResponse{},
						IndikatorRankhir:   getIndRankhir(kodeKeg),
						IndikatorPenetapan: getIndPenetapan(kodeKeg),
						SubKegiatan:        make([]renja.SubKegiatanResponse, 0),
					}

					for _, kodeSubkeg := range subkegByKeg[kodeKeg] {
						sd := subkegData[kodeSubkeg]

						kegResp.SubKegiatan = append(kegResp.SubKegiatan,
							renja.SubKegiatanResponse{
								Kode:               kodeSubkeg,
								Nama:               sd.nama,
								Jenis:              "subkegiatans",
								Tahun:              tahun,
								PegawaiId:          sd.pegawaiId,
								NamaPegawai:        sd.namaPegawai,
								PaguRankhir:        paguRankhirSubkeg[kodeSubkeg],
								PaguPenetapan:      paguPenetapanSubkeg[kodeSubkeg],
								Indikator:          []renja.IndikatorMatrixResponse{},
								IndikatorRankhir:   getIndRankhir(kodeSubkeg),
								IndikatorPenetapan: getIndPenetapan(kodeSubkeg),
							},
						)
					}

					prgResp.Kegiatan = append(prgResp.Kegiatan, kegResp)
				}

				bidangResp.Program = append(bidangResp.Program, prgResp)
			}

			urusanResp.BidangUrusan = append(urusanResp.BidangUrusan, bidangResp)
		}

		urusanDetail.Urusan = append(urusanDetail.Urusan, urusanResp)
	}

	return []renja.UrusanDetailResponse{urusanDetail}
}

// ─────────────────────────────────────────────────────────────────────────────
// CRUD: target per layer (rankhir / penetapan)
// ─────────────────────────────────────────────────────────────────────────────

func (s *MatrixRenjaServiceImpl) UpsertTargetRenjaLayer(ctx context.Context, jenis string, req programkegiatan.TargetRenjaLayerRequest) (renja.TargetLayerResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return renja.TargetLayerResponse{}, err
	}
	defer tx.Rollback()
	if err := s.MatrixRenjaRepository.UpsertTargetRenjaLayer(ctx, tx, req.KodeIndikator, req.Tahun, req.Target, req.Satuan, jenis); err != nil {
		return renja.TargetLayerResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return renja.TargetLayerResponse{}, err
	}
	return renja.TargetLayerResponse{
		KodeIndikator: req.KodeIndikator,
		Tahun:         req.Tahun,
		Target:        req.Target,
		Satuan:        req.Satuan,
		Jenis:         jenis,
	}, nil
}

func (s *MatrixRenjaServiceImpl) DeleteTargetRenjaLayer(ctx context.Context, kodeIndikator, tahun, jenis string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.MatrixRenjaRepository.DeleteTargetRenjaLayer(ctx, tx, kodeIndikator, tahun, jenis); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *MatrixRenjaServiceImpl) FindTargetRenjaLayer(ctx context.Context, kodeIndikator, tahun, jenis string) (renja.TargetLayerResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return renja.TargetLayerResponse{}, err
	}
	defer tx.Rollback()
	t, err := s.MatrixRenjaRepository.FindTargetRenjaLayer(ctx, tx, kodeIndikator, tahun, jenis)
	if err != nil {
		return renja.TargetLayerResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return renja.TargetLayerResponse{}, err
	}
	return renja.TargetLayerResponse{
		KodeIndikator: t.IndikatorId,
		Tahun:         t.Tahun,
		Target:        t.Target,
		Satuan:        t.Satuan,
		Jenis:         t.Jenis,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// CRUD: pagu per jenis (ranwal / rankhir)
// ─────────────────────────────────────────────────────────────────────────────

func (s *MatrixRenjaServiceImpl) UpsertPaguRenjaByJenis(ctx context.Context, jenis string, req programkegiatan.PaguRenjaLayerRequest) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.MatrixRenjaRepository.UpsertPaguRenjaByJenis(ctx, tx, req.KodeSubKegiatan, req.KodeOpd, req.Tahun, jenis, req.Pagu); err != nil {
		return err
	}
	return tx.Commit()
}
