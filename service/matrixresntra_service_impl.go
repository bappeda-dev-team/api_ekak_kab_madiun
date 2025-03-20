package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"strconv"
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
					Indikator:    []programkegiatan.IndikatorResponse{},
					BidangUrusan: []programkegiatan.BidangUrusanResponse{},
				},
			}
			urusanMap[item.KodeUrusan] = urusan
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
				Indikator: []programkegiatan.IndikatorResponse{},
				Program:   []programkegiatan.ProgramResponse{},
			})
			bidangUrusan = &urusan.Urusan.BidangUrusan[len(urusan.Urusan.BidangUrusan)-1]
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
				Indikator: []programkegiatan.IndikatorResponse{},
				Kegiatan:  []programkegiatan.KegiatanResponse{},
			})
			program = &bidangUrusan.Program[len(bidangUrusan.Program)-1]
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
					Indikator:   []programkegiatan.IndikatorResponse{},
					SubKegiatan: []programkegiatan.SubKegiatanResponse{},
				})
				kegiatan = &program.Kegiatan[len(program.Kegiatan)-1]
			}
		}

		// Proses SubKegiatan
		if kegiatan != nil && item.KodeSubKegiatan != "" {
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
					Tahun:     item.TahunSubKegiatan,
					Indikator: []programkegiatan.IndikatorResponse{},
				})

			}
		}

		// Proses Indikator
		if item.IndikatorId != "" && item.IndikatorKodeOpd == kodeOpd {
			tahunIndikator, _ := strconv.Atoi(item.IndikatorTahun)
			tahunAwalInt, _ := strconv.Atoi(tahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

			// Hanya proses indikator yang masuk dalam range tahun
			if tahunIndikator >= tahunAwalInt && tahunIndikator <= tahunAkhirInt {
				indikator := programkegiatan.IndikatorResponse{
					Kode:      item.IndikatorKode,
					KodeOpd:   item.IndikatorKodeOpd,
					Indikator: item.Indikator,
					Tahun:     item.IndikatorTahun,
					Target:    []programkegiatan.TargetResponse{},
				}

				// Tambahkan target jika ada
				if item.Target != "" && item.Satuan != "" {
					indikator.Target = append(indikator.Target, programkegiatan.TargetResponse{
						Target: item.Target,
						Satuan: item.Satuan,
					})
				}

				// Tentukan level indikator dan tambahkan ke tempat yang sesuai
				switch item.IndikatorKode {
				case item.KodeUrusan:
					urusan.Urusan.Indikator = service.appendIndikator(urusan.Urusan.Indikator, indikator)
				case item.KodeBidangUrusan:
					if bidangUrusan != nil {
						bidangUrusan.Indikator = service.appendIndikator(bidangUrusan.Indikator, indikator)
					}
				case item.KodeProgram:
					if program != nil {
						program.Indikator = service.appendIndikator(program.Indikator, indikator)
					}
				case item.KodeKegiatan:
					if kegiatan != nil {
						kegiatan.Indikator = service.appendIndikator(kegiatan.Indikator, indikator)
					}
				case item.KodeSubKegiatan:
					if kegiatan != nil && len(kegiatan.SubKegiatan) > 0 {
						lastSubKegiatan := &kegiatan.SubKegiatan[len(kegiatan.SubKegiatan)-1]
						lastSubKegiatan.Indikator = service.appendIndikator(lastSubKegiatan.Indikator, indikator)
					}
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

func (service *MatrixRenstraServiceImpl) appendIndikator(existing []programkegiatan.IndikatorResponse, new programkegiatan.IndikatorResponse) []programkegiatan.IndikatorResponse {
	// Cek apakah indikator dengan kode dan tahun yang sama sudah ada
	indikatorExists := false
	for i := range existing {
		// Indikator dianggap sama jika kode dan tahun sama
		if existing[i].Kode == new.Kode && existing[i].Tahun == new.Tahun {
			indikatorExists = true
			// Update target jika ada target baru
			if len(new.Target) > 0 {
				targetExists := false
				for _, existingTarget := range existing[i].Target {
					if existingTarget.Target == new.Target[0].Target && existingTarget.Satuan == new.Target[0].Satuan {
						targetExists = true
						break
					}
				}
				if !targetExists {
					existing[i].Target = append(existing[i].Target, new.Target[0])
				}
			}
			break
		}
	}

	// Jika indikator belum ada, tambahkan sebagai indikator baru
	if !indikatorExists {
		return append(existing, new)
	}

	return existing
}
