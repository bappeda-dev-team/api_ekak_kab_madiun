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

	// Get data
	data, err := service.MatrixRenstraRepository.GetByKodeSubKegiatan(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	// Transform data dengan mengirimkan ctx dan tx
	result := service.transformToResponse(ctx, tx, data, kodeOpd, tahunAwal, tahunAkhir)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

//perubahan

func (service *MatrixRenstraServiceImpl) transformToResponse(ctx context.Context, tx *sql.Tx, data []domain.SubKegiatanQuery, kodeOpd string, tahunAwal string, tahunAkhir string) []programkegiatan.UrusanDetailResponse {
	// Helper function untuk membuat indikator
	createIndikator := func(kode string, tahun string, pagu int64, data []domain.SubKegiatanQuery) programkegiatan.IndikatorResponse {
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

	// Inisialisasi response utama
	urusanDetail := programkegiatan.UrusanDetailResponse{
		KodeOpd:           kodeOpd,
		TahunAwal:         tahunAwal,
		TahunAkhir:        tahunAkhir,
		PaguAnggaranTotal: make([]programkegiatan.PaguAnggaranTotalResponse, len(tahunRange)),
		Urusan:            make([]programkegiatan.UrusanResponse, 0),
	}

	// Maps untuk menyimpan data
	urusanMap := make(map[string]*programkegiatan.UrusanResponse)
	bidangUrusanMap := make(map[string]*programkegiatan.BidangUrusanResponse)
	programMap := make(map[string]*programkegiatan.ProgramResponse)
	kegiatanMap := make(map[string]*programkegiatan.KegiatanResponse)

	// Map untuk menyimpan pagu per subkegiatan per tahun
	paguSubKegiatanMap := make(map[string]map[string]int64)
	paguGrandTotal := make(map[string]int64)

	// 1. Kumpulkan pagu HANYA dari subkegiatan
	for _, item := range data {
		if item.KodeSubKegiatan != "" &&
			item.IndikatorKode == item.KodeSubKegiatan &&
			item.IndikatorKodeOpd == kodeOpd &&
			item.PaguAnggaran.Valid {
			if _, exists := paguSubKegiatanMap[item.KodeSubKegiatan]; !exists {
				paguSubKegiatanMap[item.KodeSubKegiatan] = make(map[string]int64)
			}
			paguSubKegiatanMap[item.KodeSubKegiatan][item.IndikatorTahun] = item.PaguAnggaran.Int64
			paguGrandTotal[item.IndikatorTahun] += item.PaguAnggaran.Int64
		}
	}

	// Set pagu total
	for i, tahun := range tahunRange {
		urusanDetail.PaguAnggaranTotal[i] = programkegiatan.PaguAnggaranTotalResponse{
			Tahun:        tahun,
			PaguAnggaran: paguGrandTotal[tahun],
		}
	}

	// 2. Proses subkegiatan
	for _, item := range data {
		if item.KodeSubKegiatan == "" {
			continue
		}

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

		var subKegiatan *programkegiatan.SubKegiatanResponse
		for i := range kegiatan.SubKegiatan {
			if kegiatan.SubKegiatan[i].Kode == item.KodeSubKegiatan {
				subKegiatan = &kegiatan.SubKegiatan[i]
				break
			}
		}

		pegawai, _ := service.PegawaiRepository.FindByNip(ctx, tx, item.PegawaiId)

		if subKegiatan == nil {
			newSubKegiatan := programkegiatan.SubKegiatanResponse{
				Kode:        item.KodeSubKegiatan,
				Nama:        item.NamaSubKegiatan,
				Jenis:       "subkegiatans",
				PegawaiId:   item.PegawaiId,
				NamaPegawai: pegawai.NamaPegawai,
				Indikator:   make([]programkegiatan.IndikatorResponse, len(tahunRange)),
			}

			// Inisialisasi indikator untuk setiap tahun
			for i, tahun := range tahunRange {
				pagu := paguSubKegiatanMap[item.KodeSubKegiatan][tahun]
				newSubKegiatan.Indikator[i] = createIndikator(
					item.KodeSubKegiatan,
					tahun,
					pagu,
					data,
				)
			}

			kegiatan.SubKegiatan = append(kegiatan.SubKegiatan, newSubKegiatan)
		}
	}

	// 3. Build struktur dari bawah ke atas
	// 3.1 Kelompokkan kegiatan ke program
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

		if kegiatan, ok := kegiatanMap[item.KodeKegiatan]; ok {
			var exists bool
			for _, existingKegiatan := range program.Kegiatan {
				if existingKegiatan.Kode == kegiatan.Kode {
					exists = true
					break
				}
			}
			if !exists {
				// Hitung pagu kegiatan dari subkegiatan
				paguKegiatan := make(map[string]int64)
				for _, subKegiatan := range kegiatan.SubKegiatan {
					for _, tahun := range tahunRange {
						if pagu, exists := paguSubKegiatanMap[subKegiatan.Kode][tahun]; exists {
							paguKegiatan[tahun] += pagu
						}
					}
				}

				// Set indikator kegiatan
				for i, tahun := range tahunRange {
					kegiatan.Indikator[i] = createIndikator(
						item.KodeKegiatan,
						tahun,
						paguKegiatan[tahun],
						data,
					)
				}

				program.Kegiatan = append(program.Kegiatan, *kegiatan)
			}
		}

		// Hitung pagu program dari subkegiatan
		paguProgram := make(map[string]int64)
		for _, keg := range program.Kegiatan {
			for _, subKeg := range keg.SubKegiatan {
				for _, tahun := range tahunRange {
					if pagu, exists := paguSubKegiatanMap[subKeg.Kode][tahun]; exists {
						paguProgram[tahun] += pagu
					}
				}
			}
		}

		// Set indikator program
		for i, tahun := range tahunRange {
			program.Indikator[i] = createIndikator(
				item.KodeProgram,
				tahun,
				paguProgram[tahun],
				data,
			)
		}
	}

	// 3.2 Kelompokkan program ke bidang urusan
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

		// Hitung pagu bidang urusan dari subkegiatan
		paguBidangUrusan := make(map[string]int64)
		for _, prog := range bidangUrusan.Program {
			for _, keg := range prog.Kegiatan {
				for _, subKeg := range keg.SubKegiatan {
					for _, tahun := range tahunRange {
						if pagu, exists := paguSubKegiatanMap[subKeg.Kode][tahun]; exists {
							paguBidangUrusan[tahun] += pagu
						}
					}
				}
			}
		}

		// Set indikator bidang urusan
		for i, tahun := range tahunRange {
			bidangUrusan.Indikator[i] = createIndikator(
				item.KodeBidangUrusan,
				tahun,
				paguBidangUrusan[tahun],
				data,
			)
		}
	}

	// 3.3 Kelompokkan bidang urusan ke urusan
	for _, item := range data {
		if item.KodeUrusan == "" {
			continue
		}

		urusan, exists := urusanMap[item.KodeUrusan]
		if !exists {
			urusan = &programkegiatan.UrusanResponse{
				Kode:         item.KodeUrusan,
				Nama:         item.NamaUrusan,
				Jenis:        "urusans",
				Indikator:    make([]programkegiatan.IndikatorResponse, len(tahunRange)),
				BidangUrusan: []programkegiatan.BidangUrusanResponse{},
			}
			urusanMap[item.KodeUrusan] = urusan
		}

		if bidangUrusan, ok := bidangUrusanMap[item.KodeBidangUrusan]; ok {
			var exists bool
			for _, existingBidangUrusan := range urusan.BidangUrusan {
				if existingBidangUrusan.Kode == bidangUrusan.Kode {
					exists = true
					break
				}
			}
			if !exists {
				urusan.BidangUrusan = append(urusan.BidangUrusan, *bidangUrusan)
			}
		}

		// Hitung pagu urusan dari subkegiatan
		paguUrusan := make(map[string]int64)
		for _, bidangUrusan := range urusan.BidangUrusan {
			for _, prog := range bidangUrusan.Program {
				for _, keg := range prog.Kegiatan {
					for _, subKeg := range keg.SubKegiatan {
						for _, tahun := range tahunRange {
							if pagu, exists := paguSubKegiatanMap[subKeg.Kode][tahun]; exists {
								paguUrusan[tahun] += pagu
							}
						}
					}
				}
			}
		}

		// Set indikator urusan
		for i, tahun := range tahunRange {
			urusan.Indikator[i] = createIndikator(
				item.KodeUrusan,
				tahun,
				paguUrusan[tahun],
				data,
			)
		}
	}

	// Konversi map ke slice untuk hasil akhir
	for _, urusan := range urusanMap {
		urusanDetail.Urusan = append(urusanDetail.Urusan, *urusan)
	}

	return []programkegiatan.UrusanDetailResponse{urusanDetail}
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
