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
	DB                      *sql.DB
}

func NewMatrixRenstraServiceImpl(
	matrixRenstraRepository repository.MatrixRenstraRepository,
	periodeRepository repository.PeriodeRepository,
	db *sql.DB,
) *MatrixRenstraServiceImpl {
	return &MatrixRenstraServiceImpl{
		MatrixRenstraRepository: matrixRenstraRepository,
		PeriodeRepository:       periodeRepository,
		DB:                      db,
	}
}

func (service *MatrixRenstraServiceImpl) GetByKodeSubKegiatan(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string) ([]programkegiatan.UrusanDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get data
	data, err := service.MatrixRenstraRepository.GetByKodeSubKegiatan(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	// Transform data
	result := service.transformToResponse(data, kodeOpd, tahunAwal, tahunAkhir)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *MatrixRenstraServiceImpl) transformToResponse(data []domain.SubKegiatanQuery, kodeOpd string, tahunAwal string, tahunAkhir string) []programkegiatan.UrusanDetailResponse {
	// Helper function untuk membuat indikator
	createIndikator := func(kode string, tahun string, pagu int64, data []domain.SubKegiatanQuery) programkegiatan.IndikatorResponse {
		// Cari indikator yang sesuai dengan kode dan tahun
		for _, item := range data {
			if item.IndikatorKode == kode &&
				item.IndikatorKodeOpd == kodeOpd &&
				item.IndikatorTahun == tahun {
				return programkegiatan.IndikatorResponse{
					Id:           item.IndikatorId,
					Kode:         kode,
					KodeOpd:      kodeOpd,
					Indikator:    item.Indikator,
					Tahun:        tahun,
					PaguAnggaran: pagu,
					Target: []programkegiatan.TargetResponse{
						{
							Id:     item.TargetId,
							Target: item.Target,
							Satuan: item.Satuan,
						},
					},
				}
			}
		}

		// Jika tidak ada, buat indikator kosong
		return programkegiatan.IndikatorResponse{
			Kode:         kode,
			KodeOpd:      kodeOpd,
			Indikator:    "",
			Tahun:        tahun,
			PaguAnggaran: pagu,
			Target: []programkegiatan.TargetResponse{
				{
					Target: "",
					Satuan: "",
				},
			},
		}
	}

	// Generate tahun range
	tahunAwalInt, _ := strconv.Atoi(tahunAwal)
	tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)
	var tahunRange []string
	for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
		tahunRange = append(tahunRange, strconv.Itoa(tahun))
	}

	// Maps untuk menyimpan data
	urusanMap := make(map[string]*programkegiatan.UrusanDetailResponse)
	bidangUrusanMap := make(map[string]*programkegiatan.BidangUrusanResponse)
	programMap := make(map[string]*programkegiatan.ProgramResponse)
	kegiatanMap := make(map[string]*programkegiatan.KegiatanResponse)

	// Map untuk menyimpan pagu per subkegiatan per tahun
	paguSubKegiatanMap := make(map[string]map[string]int64) // map[kodeSubKegiatan][tahun]pagu

	// 1. Pertama, kumpulkan semua pagu dari subkegiatan
	for _, item := range data {
		if item.KodeSubKegiatan != "" && item.PaguAnggaran.Valid {
			if _, exists := paguSubKegiatanMap[item.KodeSubKegiatan]; !exists {
				paguSubKegiatanMap[item.KodeSubKegiatan] = make(map[string]int64)
			}
			paguSubKegiatanMap[item.KodeSubKegiatan][item.IndikatorTahun] = item.PaguAnggaran.Int64
		}
	}

	// 2. Proses subkegiatan
	for _, item := range data {
		if item.KodeSubKegiatan == "" {
			continue
		}

		// Inisialisasi kegiatan jika belum ada
		kegiatan, exists := kegiatanMap[item.KodeKegiatan]
		if !exists {
			kegiatan = &programkegiatan.KegiatanResponse{
				Kode:        item.KodeKegiatan,
				Nama:        item.NamaKegiatan,
				Jenis:       "kegiatans",
				Indikator:   make([]programkegiatan.IndikatorResponse, len(tahunRange)),
				SubKegiatan: []programkegiatan.SubKegiatanResponse{},
			}
			kegiatanMap[item.KodeKegiatan] = kegiatan
		}

		// Cek dan tambah subkegiatan
		var subKegiatan *programkegiatan.SubKegiatanResponse
		for i := range kegiatan.SubKegiatan {
			if kegiatan.SubKegiatan[i].Kode == item.KodeSubKegiatan {
				subKegiatan = &kegiatan.SubKegiatan[i]
				break
			}
		}

		if subKegiatan == nil {
			newSubKegiatan := programkegiatan.SubKegiatanResponse{
				Kode:      item.KodeSubKegiatan,
				Nama:      item.NamaSubKegiatan,
				Jenis:     "subkegiatans",
				Indikator: make([]programkegiatan.IndikatorResponse, len(tahunRange)),
			}

			// Inisialisasi indikator kosong untuk setiap tahun
			for i, tahun := range tahunRange {
				newSubKegiatan.Indikator[i] = programkegiatan.IndikatorResponse{
					Kode:         item.KodeSubKegiatan,
					KodeOpd:      kodeOpd,
					Indikator:    "",
					Tahun:        tahun,
					PaguAnggaran: paguSubKegiatanMap[item.KodeSubKegiatan][tahun],
					Target: []programkegiatan.TargetResponse{
						{
							Target: "",
							Satuan: "",
						},
					},
				}
			}

			kegiatan.SubKegiatan = append(kegiatan.SubKegiatan, newSubKegiatan)
			subKegiatan = &kegiatan.SubKegiatan[len(kegiatan.SubKegiatan)-1]
		}

		// Update indikator jika ada data
		if item.IndikatorId != "" &&
			item.IndikatorKode == item.KodeSubKegiatan &&
			item.IndikatorKodeOpd == kodeOpd {
			for i, ind := range subKegiatan.Indikator {
				if ind.Tahun == item.IndikatorTahun {
					subKegiatan.Indikator[i] = programkegiatan.IndikatorResponse{
						Id:           item.IndikatorId,
						Kode:         item.IndikatorKode,
						KodeOpd:      item.IndikatorKodeOpd,
						Indikator:    item.Indikator,
						Tahun:        item.IndikatorTahun,
						PaguAnggaran: paguSubKegiatanMap[item.KodeSubKegiatan][item.IndikatorTahun],
						Target: []programkegiatan.TargetResponse{
							{
								Id:     item.TargetId,
								Target: item.Target,
								Satuan: item.Satuan,
							},
						},
					}
					break
				}
			}
		}
	}

	// 3. Build struktur dari bawah ke atas dan hitung pagu
	// 3.1 Kelompokkan kegiatan ke program dan hitung pagu kegiatan
	for _, item := range data {
		if item.KodeProgram == "" {
			continue
		}

		program, exists := programMap[item.KodeProgram]
		if !exists {
			program = &programkegiatan.ProgramResponse{
				Kode:      item.KodeProgram,
				Nama:      item.NamaProgram,
				Jenis:     "programs",
				Indikator: make([]programkegiatan.IndikatorResponse, len(tahunRange)),
				Kegiatan:  []programkegiatan.KegiatanResponse{},
			}
			programMap[item.KodeProgram] = program
		}

		// Tambahkan kegiatan ke program
		if kegiatan, ok := kegiatanMap[item.KodeKegiatan]; ok {
			var exists bool
			for _, existingKegiatan := range program.Kegiatan {
				if existingKegiatan.Kode == kegiatan.Kode {
					exists = true
					break
				}
			}
			if !exists {
				// Hitung total pagu kegiatan dari subkegiatan
				paguKegiatan := make(map[string]int64)
				for _, subKegiatan := range kegiatan.SubKegiatan {
					for _, ind := range subKegiatan.Indikator {
						paguKegiatan[ind.Tahun] += ind.PaguAnggaran
					}
				}

				// Set indikator kegiatan dengan pagu yang sudah dihitung
				// Set indikator kegiatan dengan pagu yang sudah dihitung
				for i, tahun := range tahunRange {
					kegiatan.Indikator[i] = createIndikator(
						item.KodeKegiatan,
						tahun,
						paguKegiatan[tahun],
						data, // Kirim seluruh data
					)
				}

				program.Kegiatan = append(program.Kegiatan, *kegiatan)
			}
		}

		// Hitung total pagu program dari kegiatan
		paguProgram := make(map[string]int64)
		for _, keg := range program.Kegiatan {
			for _, ind := range keg.Indikator {
				paguProgram[ind.Tahun] += ind.PaguAnggaran
			}
		}

		// Set indikator program dengan pagu yang sudah dihitung
		for i, tahun := range tahunRange {
			program.Indikator[i] = createIndikator(
				item.KodeProgram,
				tahun,
				paguProgram[tahun],
				data, // Kirim seluruh data
			)
		}
	}

	// 3.2 Kelompokkan program ke bidang urusan dan hitung pagu bidang urusan
	for _, item := range data {
		if item.KodeBidangUrusan == "" {
			continue
		}

		bidangUrusan, exists := bidangUrusanMap[item.KodeBidangUrusan]
		if !exists {
			bidangUrusan = &programkegiatan.BidangUrusanResponse{
				Kode:      item.KodeBidangUrusan,
				Nama:      item.NamaBidangUrusan,
				Jenis:     "bidang_urusans",
				Indikator: make([]programkegiatan.IndikatorResponse, len(tahunRange)),
				Program:   []programkegiatan.ProgramResponse{},
			}
			bidangUrusanMap[item.KodeBidangUrusan] = bidangUrusan
		}

		// Tambahkan program ke bidang urusan
		if program, ok := programMap[item.KodeProgram]; ok {
			var exists bool
			for _, existingProgram := range bidangUrusan.Program {
				if existingProgram.Kode == program.Kode {
					exists = true
					break
				}
			}
			if !exists {
				bidangUrusan.Program = append(bidangUrusan.Program, *program)
			}
		}

		// Hitung total pagu bidang urusan dari program
		paguBidangUrusan := make(map[string]int64)
		for _, prog := range bidangUrusan.Program {
			for _, ind := range prog.Indikator {
				paguBidangUrusan[ind.Tahun] += ind.PaguAnggaran
			}
		}

		// Set indikator bidang urusan dengan pagu yang sudah dihitung
		for i, tahun := range tahunRange {
			bidangUrusan.Indikator[i] = createIndikator(
				item.KodeBidangUrusan,
				tahun,
				paguBidangUrusan[tahun],
				data, // Kirim seluruh data
			)
		}
	}

	// 3.3 Kelompokkan bidang urusan ke urusan dan hitung pagu urusan
	for _, item := range data {
		if item.KodeUrusan == "" {
			continue
		}

		urusan, exists := urusanMap[item.KodeUrusan]
		if !exists {
			urusan = &programkegiatan.UrusanDetailResponse{
				KodeOpd:           kodeOpd,
				TahunAwal:         tahunAwal,
				TahunAkhir:        tahunAkhir,
				PaguAnggaranTotal: make([]programkegiatan.PaguAnggaranTotalResponse, len(tahunRange)),
				Urusan: programkegiatan.UrusanResponse{
					Kode:         item.KodeUrusan,
					Nama:         item.NamaUrusan,
					Jenis:        "urusans",
					Indikator:    make([]programkegiatan.IndikatorResponse, len(tahunRange)),
					BidangUrusan: []programkegiatan.BidangUrusanResponse{},
				},
			}
			// Inisialisasi pagu total untuk setiap tahun
			for i, tahun := range tahunRange {
				urusan.PaguAnggaranTotal[i] = programkegiatan.PaguAnggaranTotalResponse{
					Tahun:        tahun,
					PaguAnggaran: 0,
				}
			}
			urusanMap[item.KodeUrusan] = urusan
		}

		// Tambahkan bidang urusan ke urusan
		if bidangUrusan, ok := bidangUrusanMap[item.KodeBidangUrusan]; ok {
			var exists bool
			for _, existingBidangUrusan := range urusan.Urusan.BidangUrusan {
				if existingBidangUrusan.Kode == bidangUrusan.Kode {
					exists = true
					break
				}
			}
			if !exists {
				urusan.Urusan.BidangUrusan = append(urusan.Urusan.BidangUrusan, *bidangUrusan)
			}
		}

		// Hitung total pagu urusan dari bidang urusan
		paguUrusan := make(map[string]int64)
		for _, bidangUrusan := range urusan.Urusan.BidangUrusan {
			for _, ind := range bidangUrusan.Indikator {
				paguUrusan[ind.Tahun] += ind.PaguAnggaran
			}
		}

		// Set indikator urusan dan pagu total dengan pagu yang sudah dihitung
		for i, tahun := range tahunRange {
			urusan.Urusan.Indikator[i] = createIndikator(
				item.KodeUrusan,
				tahun,
				paguUrusan[tahun],
				data, // Kirim seluruh data
			)
			urusan.PaguAnggaranTotal[i].PaguAnggaran = paguUrusan[tahun]
		}
	}

	// Convert map to slice untuk hasil akhir
	var result []programkegiatan.UrusanDetailResponse
	for _, urusan := range urusanMap {
		result = append(result, *urusan)
	}
	return result
}

// crud
func (service *MatrixRenstraServiceImpl) CreateIndikator(ctx context.Context, request programkegiatan.IndikatorRenstraCreateRequest) (programkegiatan.IndikatorResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}
	defer tx.Rollback()

	// Generate ID dengan format IND-MTRX-uuid7

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("IND-RNST-%s", randomDigits)

	// Simpan indikator
	indikator := domain.Indikator{
		Id:           uuId,
		Kode:         request.Kode,
		KodeOpd:      request.KodeOpd,
		Indikator:    request.Indikator,
		Tahun:        request.Tahun,
		PaguAnggaran: request.PaguAnggaran,
	}

	err = service.MatrixRenstraRepository.SaveIndikator(ctx, tx, indikator)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	// Simpan target
	uuIdTarget := fmt.Sprintf("TRG-RNST-%s", randomDigits)

	target := domain.Target{
		Id:          uuIdTarget,
		IndikatorId: indikator.Id,
		Target:      request.Target,
		Satuan:      request.Satuan,
	}

	err = service.MatrixRenstraRepository.SaveTarget(ctx, tx, target)
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return programkegiatan.IndikatorResponse{}, err
	}

	return programkegiatan.IndikatorResponse{
		Id:           indikator.Id,
		Kode:         request.Kode,
		KodeOpd:      request.KodeOpd,
		Indikator:    request.Indikator,
		Tahun:        request.Tahun,
		PaguAnggaran: request.PaguAnggaran,
		Target: []programkegiatan.TargetResponse{
			{
				Id:     target.Id,
				Target: request.Target,
				Satuan: request.Satuan,
			},
		},
	}, nil
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
		Id:           request.Id,
		Kode:         request.Kode,
		KodeOpd:      request.KodeOpd,
		Indikator:    request.Indikator,
		Tahun:        request.Tahun,
		PaguAnggaran: request.PaguAnggaran,
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
		Id:           request.Id,
		Kode:         request.Kode,
		KodeOpd:      request.KodeOpd,
		Indikator:    request.Indikator,
		Tahun:        request.Tahun,
		PaguAnggaran: request.PaguAnggaran,
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
		Id:           indikator.Id,
		Kode:         indikator.Kode,
		KodeOpd:      indikator.KodeOpd,
		Indikator:    indikator.Indikator,
		Tahun:        indikator.Tahun,
		PaguAnggaran: indikator.PaguAnggaran,
		Target:       make([]programkegiatan.TargetResponse, 0),
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
