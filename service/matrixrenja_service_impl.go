package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
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

func (service *MatrixRenjaServiceImpl) GetByKodeOpdAndTahun(ctx context.Context, kodeOpd string, tahun string) ([]programkegiatan.UrusanDetailResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := service.MatrixRenjaRepository.GetByKodeOpdAndTahun(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	result := service.transformToResponse(ctx, tx, data, kodeOpd, tahun)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *MatrixRenjaServiceImpl) transformToResponse(ctx context.Context, tx *sql.Tx, data []domain.SubKegiatanQuery, kodeOpd string, tahun string) []programkegiatan.UrusanDetailResponse {
	// Helper function untuk membuat indikator
	createIndikator := func(kode string, data []domain.SubKegiatanQuery) programkegiatan.IndikatorResponse {
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
					PaguAnggaran: item.PaguAnggaran.Int64,
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
			PaguAnggaran: 0,
			Target: []programkegiatan.TargetResponse{
				{
					Target: "",
					Satuan: "",
				},
			},
		}
	}

	// Inisialisasi response utama
	urusanDetail := programkegiatan.UrusanDetailResponse{
		KodeOpd: kodeOpd,
		Tahun:   tahun,
		PaguAnggaranTotal: []programkegiatan.PaguAnggaranTotalResponse{
			{
				Tahun:        tahun,
				PaguAnggaran: 0, // akan diupdate nanti
			},
		},
		Urusan: []programkegiatan.UrusanResponse{},
	}

	// Maps untuk menyimpan data
	urusanMap := make(map[string]*programkegiatan.UrusanResponse)
	bidangUrusanMap := make(map[string]*programkegiatan.BidangUrusanResponse)
	programMap := make(map[string]*programkegiatan.ProgramResponse)
	kegiatanMap := make(map[string]*programkegiatan.KegiatanResponse)

	// Map untuk menyimpan pagu
	paguSubKegiatanMap := make(map[string]int64)
	var paguTotal int64 = 0

	// 1. Kumpulkan pagu dari subkegiatan
	for _, item := range data {
		if item.KodeSubKegiatan != "" &&
			item.IndikatorKode == item.KodeSubKegiatan &&
			item.IndikatorKodeOpd == kodeOpd &&
			item.PaguAnggaran.Valid {
			paguSubKegiatanMap[item.KodeSubKegiatan] = item.PaguAnggaran.Int64
			paguTotal += item.PaguAnggaran.Int64
		}
	}

	// Set pagu total
	urusanDetail.PaguAnggaranTotal[0].PaguAnggaran = paguTotal

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
				Indikator:   []programkegiatan.IndikatorResponse{},
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
			// pagu := paguSubKegiatanMap[item.KodeSubKegiatan]
			newSubKegiatan := programkegiatan.SubKegiatanResponse{
				Kode:        item.KodeSubKegiatan,
				Nama:        item.NamaSubKegiatan,
				Jenis:       "subkegiatans",
				PegawaiId:   item.PegawaiId,
				NamaPegawai: pegawai.NamaPegawai,
				Indikator: []programkegiatan.IndikatorResponse{
					createIndikator(item.KodeSubKegiatan, data),
				},
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
				Kode:  item.KodeProgram,
				Nama:  item.NamaProgram,
				Jenis: "programs",
				Indikator: []programkegiatan.IndikatorResponse{
					createIndikator(item.KodeProgram, data),
				},
				Kegiatan: []programkegiatan.KegiatanResponse{},
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
				kegiatan.Indikator = append(kegiatan.Indikator,
					createIndikator(item.KodeKegiatan, data))
				program.Kegiatan = append(program.Kegiatan, *kegiatan)
			}
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
				Kode:  item.KodeBidangUrusan,
				Nama:  item.NamaBidangUrusan,
				Jenis: "bidang_urusans",
				Indikator: []programkegiatan.IndikatorResponse{
					createIndikator(item.KodeBidangUrusan, data),
				},
				Program: []programkegiatan.ProgramResponse{},
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
	}

	// 3.3 Kelompokkan bidang urusan ke urusan
	for _, item := range data {
		if item.KodeUrusan == "" {
			continue
		}

		urusan, exists := urusanMap[item.KodeUrusan]
		if !exists {
			urusan = &programkegiatan.UrusanResponse{
				Kode:  item.KodeUrusan,
				Nama:  item.NamaUrusan,
				Jenis: "urusans",
				Indikator: []programkegiatan.IndikatorResponse{
					createIndikator(item.KodeUrusan, data),
				},
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
	}

	// Konversi map ke slice untuk hasil akhir
	for _, urusan := range urusanMap {
		urusanDetail.Urusan = append(urusanDetail.Urusan, *urusan)
	}

	return []programkegiatan.UrusanDetailResponse{urusanDetail}
}
