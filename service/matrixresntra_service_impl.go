package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
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
	urusanMap := make(map[string]*programkegiatan.UrusanDetailResponse)

	// Helper function untuk membuat indikator
	createIndikator := func(item domain.SubKegiatanQuery) programkegiatan.IndikatorResponse {
		return programkegiatan.IndikatorResponse{
			Id:        item.IndikatorId,
			Kode:      item.IndikatorKode,
			KodeOpd:   item.IndikatorKodeOpd,
			Indikator: item.Indikator,
			Tahun:     item.IndikatorTahun,
			Target: []programkegiatan.TargetResponse{
				{
					Id:          item.TargetId,
					IndikatorId: item.IndikatorId,
					Target:      item.Target,
					Satuan:      item.Satuan,
				},
			},
		}
	}

	// Helper function untuk validasi indikator
	shouldAddIndikator := func(item domain.SubKegiatanQuery, expectedKode string) bool {
		return item.IndikatorId != "" &&
			item.IndikatorKode == expectedKode &&
			item.IndikatorKodeOpd == kodeOpd
	}

	for _, item := range data {
		// Proses Urusan
		urusan, exists := urusanMap[item.KodeUrusan]
		if !exists {
			urusan = &programkegiatan.UrusanDetailResponse{
				KodeOpd:    kodeOpd,
				TahunAwal:  tahunAwal,
				TahunAkhir: tahunAkhir,
				Urusan: programkegiatan.UrusanResponse{
					Kode:         item.KodeUrusan,
					Nama:         item.NamaUrusan,
					Jenis:        "urusans",
					Indikator:    []programkegiatan.IndikatorResponse{},
					BidangUrusan: []programkegiatan.BidangUrusanResponse{},
				},
			}
			urusanMap[item.KodeUrusan] = urusan
		}

		// Proses indikator urusan
		if shouldAddIndikator(item, item.KodeUrusan) {
			urusan.Urusan.Indikator = service.appendIndikator(
				urusan.Urusan.Indikator,
				createIndikator(item),
			)
		}

		// Proses Bidang Urusan
		var bidangUrusan *programkegiatan.BidangUrusanResponse
		for i := range urusan.Urusan.BidangUrusan {
			if urusan.Urusan.BidangUrusan[i].Kode == item.KodeBidangUrusan {
				bidangUrusan = &urusan.Urusan.BidangUrusan[i]
				break
			}
		}
		if bidangUrusan == nil {
			urusan.Urusan.BidangUrusan = append(urusan.Urusan.BidangUrusan, programkegiatan.BidangUrusanResponse{
				Kode:      item.KodeBidangUrusan,
				Nama:      item.NamaBidangUrusan,
				Jenis:     "bidang_urusans",
				Indikator: []programkegiatan.IndikatorResponse{},
				Program:   []programkegiatan.ProgramResponse{},
			})
			bidangUrusan = &urusan.Urusan.BidangUrusan[len(urusan.Urusan.BidangUrusan)-1]
		}

		// Proses indikator bidang urusan
		if shouldAddIndikator(item, item.KodeBidangUrusan) {
			bidangUrusan.Indikator = service.appendIndikator(
				bidangUrusan.Indikator,
				createIndikator(item),
			)
		}

		// Proses Program
		var program *programkegiatan.ProgramResponse
		for i := range bidangUrusan.Program {
			if bidangUrusan.Program[i].Kode == item.KodeProgram {
				program = &bidangUrusan.Program[i]
				break
			}
		}
		if program == nil && item.KodeProgram != "" {
			bidangUrusan.Program = append(bidangUrusan.Program, programkegiatan.ProgramResponse{
				Kode:      item.KodeProgram,
				Nama:      item.NamaProgram,
				Jenis:     "programs",
				Indikator: []programkegiatan.IndikatorResponse{},
				Kegiatan:  []programkegiatan.KegiatanResponse{},
			})
			program = &bidangUrusan.Program[len(bidangUrusan.Program)-1]
		}

		// Proses indikator program
		if shouldAddIndikator(item, item.KodeProgram) {
			program.Indikator = service.appendIndikator(
				program.Indikator,
				createIndikator(item),
			)
		}

		// Proses Kegiatan
		var kegiatan *programkegiatan.KegiatanResponse
		if program != nil {
			for i := range program.Kegiatan {
				if program.Kegiatan[i].Kode == item.KodeKegiatan {
					kegiatan = &program.Kegiatan[i]
					break
				}
			}
			if kegiatan == nil && item.KodeKegiatan != "" {
				program.Kegiatan = append(program.Kegiatan, programkegiatan.KegiatanResponse{
					Kode:        item.KodeKegiatan,
					Nama:        item.NamaKegiatan,
					Jenis:       "kegiatans",
					Indikator:   []programkegiatan.IndikatorResponse{},
					SubKegiatan: []programkegiatan.SubKegiatanResponse{},
				})
				kegiatan = &program.Kegiatan[len(program.Kegiatan)-1]
			}

			// Proses indikator kegiatan
			if shouldAddIndikator(item, item.KodeKegiatan) {
				kegiatan.Indikator = service.appendIndikator(
					kegiatan.Indikator,
					createIndikator(item),
				)
			}

			// Proses SubKegiatan
			if item.KodeSubKegiatan != "" {
				var subKegiatan *programkegiatan.SubKegiatanResponse
				for i := range kegiatan.SubKegiatan {
					if kegiatan.SubKegiatan[i].Kode == item.KodeSubKegiatan {
						subKegiatan = &kegiatan.SubKegiatan[i]
						break
					}
				}
				if subKegiatan == nil {
					kegiatan.SubKegiatan = append(kegiatan.SubKegiatan, programkegiatan.SubKegiatanResponse{
						Kode:      item.KodeSubKegiatan,
						Nama:      item.NamaSubKegiatan,
						Jenis:     "subkegiatans",
						Tahun:     item.TahunSubKegiatan,
						Indikator: []programkegiatan.IndikatorResponse{},
					})
					subKegiatan = &kegiatan.SubKegiatan[len(kegiatan.SubKegiatan)-1]
				}

				// Proses indikator subkegiatan
				if shouldAddIndikator(item, item.KodeSubKegiatan) {
					subKegiatan.Indikator = service.appendIndikator(
						subKegiatan.Indikator,
						createIndikator(item),
					)
				}
			}
		}
	}

	// Convert map to slice
	var result []programkegiatan.UrusanDetailResponse
	for _, urusan := range urusanMap {
		result = append(result, *urusan)
	}
	return result
}

// Helper function untuk menambahkan indikator
func (service *MatrixRenstraServiceImpl) appendIndikator(existing []programkegiatan.IndikatorResponse, new programkegiatan.IndikatorResponse) []programkegiatan.IndikatorResponse {
	// Cek apakah indikator dengan kode dan tahun yang sama sudah ada
	for i := range existing {
		if existing[i].Kode == new.Kode &&
			existing[i].Tahun == new.Tahun &&
			existing[i].KodeOpd == new.KodeOpd {
			// Update target jika ada target baru
			if len(new.Target) > 0 {
				targetExists := false
				for _, existingTarget := range existing[i].Target {
					if existingTarget.Target == new.Target[0].Target &&
						existingTarget.Satuan == new.Target[0].Satuan {
						targetExists = true
						break
					}
				}
				if !targetExists {
					existing[i].Target = append(existing[i].Target, new.Target[0])
				}
			}
			return existing
		}
	}

	// Jika indikator belum ada, tambahkan sebagai indikator baru
	return append(existing, new)
}
