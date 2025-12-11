package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/rincianbelanja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"
	"sort"
)

type RincianBelanjaServiceImpl struct {
	rincianBelanjaRepository repository.RincianBelanjaRepository
	pegawaiRepository        repository.PegawaiRepository
	DB                       *sql.DB
}

func NewRincianBelanjaServiceImpl(rincianBelanjaRepository repository.RincianBelanjaRepository, pegawaiRepository repository.PegawaiRepository, DB *sql.DB) *RincianBelanjaServiceImpl {
	return &RincianBelanjaServiceImpl{
		rincianBelanjaRepository: rincianBelanjaRepository,
		pegawaiRepository:        pegawaiRepository,
		DB:                       DB,
	}
}

func (service *RincianBelanjaServiceImpl) Create(ctx context.Context, request rincianbelanja.RincianBelanjaCreateRequest) (rincianbelanja.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.RenaksiId == "" {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("renaksi_id tidak boleh kosong")
	}
	if request.Anggaran < 0 {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("anggaran tidak boleh negatif")
	}

	// Konversi request ke domain model
	rincianBelanja := domain.RincianBelanja{
		RenaksiId: request.RenaksiId,
		Anggaran:  int64(request.Anggaran),
	}

	// Simpan ke database
	result, err := service.rincianBelanjaRepository.Create(ctx, tx, rincianBelanja)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}

	// Konversi domain model ke response
	response := rincianbelanja.RencanaAksiResponse{
		RenaksiId: result.RenaksiId,
		Renaksi:   result.Renaksi,
		Anggaran:  int(result.Anggaran),
	}

	return response, nil
}

func (service *RincianBelanjaServiceImpl) Update(ctx context.Context, request rincianbelanja.RincianBelanjaUpdateRequest) (rincianbelanja.RencanaAksiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	existing, err := service.rincianBelanjaRepository.FindByRenaksiId(ctx, tx, request.RenaksiId)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}
	if existing.RenaksiId == "" {
		return rincianbelanja.RencanaAksiResponse{}, errors.New("rincian belanja tidak ditemukan")
	}

	// Konversi request ke domain model
	rincianBelanja := domain.RincianBelanja{
		RenaksiId: request.RenaksiId,
		Anggaran:  int64(request.Anggaran),
	}

	// Update ke database
	result, err := service.rincianBelanjaRepository.Update(ctx, tx, rincianBelanja)
	if err != nil {
		return rincianbelanja.RencanaAksiResponse{}, err
	}

	// Konversi domain model ke response
	response := rincianbelanja.RencanaAksiResponse{
		RenaksiId: result.RenaksiId,
		Renaksi:   result.Renaksi,
		Anggaran:  int(result.Anggaran),
	}

	return response, nil
}

func (service *RincianBelanjaServiceImpl) FindRincianBelanjaAsn(ctx context.Context, pegawaiId string, tahun string) []rincianbelanja.RincianBelanjaAsnResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	rincianBelanjaList, err := service.rincianBelanjaRepository.FindRincianBelanjaAsn(ctx, tx, pegawaiId, tahun)
	if err != nil {
		panic(err)
	}

	pegawai, err := service.pegawaiRepository.FindByNip(ctx, tx, pegawaiId)
	if err != nil {
		panic(err)
	}

	var responses []rincianbelanja.RincianBelanjaAsnResponse
	for _, rb := range rincianBelanjaList {
		// Ambil indikator subkegiatan
		indikatorSubkegiatan, err := service.rincianBelanjaRepository.FindIndikatorSubkegiatanByKodeAndOpd(
			ctx,
			tx,
			rb.KodeSubkegiatan,
			pegawai.KodeOpd,
			tahun,
		)
		if err != nil {
			log.Printf("Error mengambil indikator subkegiatan: %v", err)
			continue
		}

		// Sort indikator subkegiatan berdasarkan ID
		sort.Slice(indikatorSubkegiatan, func(i, j int) bool {
			return indikatorSubkegiatan[i].Id < indikatorSubkegiatan[j].Id
		})

		// Konversi indikator subkegiatan ke response
		var indikatorSubkegiatanResponses []rincianbelanja.IndikatorResponse
		for _, ind := range indikatorSubkegiatan {
			// Sort target berdasarkan ID
			sort.Slice(ind.Target, func(i, j int) bool {
				return ind.Target[i].Id < ind.Target[j].Id
			})

			var targetResponses []rincianbelanja.TargetResponse
			for _, t := range ind.Target {
				targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
					Id:          t.Id,
					IndikatorId: t.IndikatorId,
					Target:      t.Target,
					Satuan:      t.Satuan,
				})
			}

			indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, rincianbelanja.IndikatorResponse{
				Id:              ind.Id,
				KodeSubkegiatan: ind.Kode,
				KodeOPD:         ind.KodeOpd,
				NamaIndikator:   ind.Indikator,
				Target:          targetResponses,
			})
		}

		var rencanaKinerjaResponses []rincianbelanja.RincianBelanjaResponse
		for _, rk := range rb.RencanaKinerja {
			var rencanaAksiResponses []rincianbelanja.RencanaAksiResponse
			var totalAnggaranRekin int = 0

			// Ambil indikator berdasarkan ID rencana kinerja
			indikators, err := service.rincianBelanjaRepository.FindIndikatorByRekinId(ctx, tx, rk.RencanaKinerjaId)
			if err != nil {
				log.Printf("Error mengambil indikator untuk rekin %s: %v", rk.RencanaKinerjaId, err)
				continue
			}

			// Sort indikator berdasarkan ID
			sort.Slice(indikators, func(i, j int) bool {
				return indikators[i].Id < indikators[j].Id
			})

			// Konversi indikator ke response
			var indikatorResponses []rincianbelanja.IndikatorResponse
			for _, ind := range indikators {
				// Sort target berdasarkan ID
				sort.Slice(ind.Target, func(i, j int) bool {
					return ind.Target[i].Id < ind.Target[j].Id
				})

				var targetResponses []rincianbelanja.TargetResponse
				for _, t := range ind.Target {
					targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
						Id:          t.Id,
						IndikatorId: t.IndikatorId,
						Target:      t.Target,
						Satuan:      t.Satuan,
					})
				}

				indikatorResponses = append(indikatorResponses, rincianbelanja.IndikatorResponse{
					Id:               ind.Id,
					RencanaKinerjaId: ind.RencanaKinerjaId,
					NamaIndikator:    ind.Indikator,
					Target:           targetResponses,
				})
			}

			// Sort rencana aksi berdasarkan ID
			if rk.RencanaAksi != nil {
				totalAnggaranRekin = 0
				for _, ra := range rk.RencanaAksi {
					rencanaAksiResponses = append(rencanaAksiResponses, rincianbelanja.RencanaAksiResponse{
						RenaksiId: ra.RenaksiId,
						Renaksi:   ra.Renaksi,
						Anggaran:  int(ra.Anggaran),
					})
					totalAnggaranRekin += int(ra.Anggaran)
				}
			}

			if rencanaAksiResponses == nil {
				rencanaAksiResponses = make([]rincianbelanja.RencanaAksiResponse, 0)
			}

			rencanaKinerjaResponses = append(rencanaKinerjaResponses, rincianbelanja.RincianBelanjaResponse{
				RencanaKinerjaId: rk.RencanaKinerjaId,
				RencanaKinerja:   rk.RencanaKinerja,
				Indikator:        indikatorResponses,
				TotalAnggaran:    totalAnggaranRekin,
				RencanaAksi:      rencanaAksiResponses,
			})
		}

		// Sort rencana kinerja responses berdasarkan ID
		sort.Slice(rencanaKinerjaResponses, func(i, j int) bool {
			return rencanaKinerjaResponses[i].RencanaKinerjaId < rencanaKinerjaResponses[j].RencanaKinerjaId
		})

		responses = append(responses, rincianbelanja.RincianBelanjaAsnResponse{
			PegawaiId:            rb.PegawaiId,
			NamaPegawai:          rb.NamaPegawai,
			KodeSubkegiatan:      rb.KodeSubkegiatan,
			NamaSubkegiatan:      rb.NamaSubkegiatan,
			IndikatorSubkegiatan: indikatorSubkegiatanResponses,
			TotalAnggaran:        rb.TotalAnggaran,
			RincianBelanja:       rencanaKinerjaResponses,
		})
	}

	return responses
}

func (service *RincianBelanjaServiceImpl) LaporanRincianBelanjaOpd(ctx context.Context, kodeOpd string, tahun string) ([]rincianbelanja.RincianBelanjaAsnResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	rincianBelanjaList, err := service.rincianBelanjaRepository.LaporanRincianBelanjaOpd(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Map untuk mengelompokkan berdasarkan kode subkegiatan
	subkegiatanMap := make(map[string]*rincianbelanja.RincianBelanjaAsnResponse)

	for _, rb := range rincianBelanjaList {
		// Ambil atau buat response subkegiatan
		subResponse, exists := subkegiatanMap[rb.KodeSubkegiatan]
		if !exists {
			// Ambil indikator subkegiatan
			indikatorSubkegiatan, err := service.rincianBelanjaRepository.FindIndikatorSubkegiatanByKodeAndOpd(
				ctx,
				tx,
				rb.KodeSubkegiatan,
				kodeOpd,
				tahun,
			)
			if err != nil {
				log.Printf("Error mengambil indikator subkegiatan: %v", err)
				continue
			}

			// Sort dan konversi indikator subkegiatan
			sort.Slice(indikatorSubkegiatan, func(i, j int) bool {
				return indikatorSubkegiatan[i].Id < indikatorSubkegiatan[j].Id
			})

			var indikatorSubkegiatanResponses []rincianbelanja.IndikatorResponse
			for _, ind := range indikatorSubkegiatan {
				sort.Slice(ind.Target, func(i, j int) bool {
					return ind.Target[i].Id < ind.Target[j].Id
				})

				var targetResponses []rincianbelanja.TargetResponse
				for _, t := range ind.Target {
					targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
						Id:          t.Id,
						IndikatorId: t.IndikatorId,
						Target:      t.Target,
						Satuan:      t.Satuan,
					})
				}

				indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, rincianbelanja.IndikatorResponse{
					Id:              ind.Id,
					KodeSubkegiatan: ind.Kode,
					KodeOPD:         ind.KodeOpd,
					NamaIndikator:   ind.Indikator,
					Target:          targetResponses,
				})
			}

			subResponse = &rincianbelanja.RincianBelanjaAsnResponse{
				KodeSubkegiatan:      rb.KodeSubkegiatan,
				NamaSubkegiatan:      rb.NamaSubkegiatan,
				IndikatorSubkegiatan: indikatorSubkegiatanResponses,
				TotalAnggaran:        0,
				RincianBelanja:       []rincianbelanja.RincianBelanjaResponse{},
			}
			subkegiatanMap[rb.KodeSubkegiatan] = subResponse
		}

		// Proses rencana kinerja
		for _, rk := range rb.RencanaKinerja {
			var rencanaAksiResponses []rincianbelanja.RencanaAksiResponse
			var totalAnggaranRekin int = 0

			// Ambil dan proses indikator rencana kinerja
			indikators, err := service.rincianBelanjaRepository.FindIndikatorByRekinId(ctx, tx, rk.RencanaKinerjaId)
			if err != nil {
				log.Printf("Error mengambil indikator untuk rekin %s: %v", rk.RencanaKinerjaId, err)
				continue
			}

			// Sort indikator berdasarkan ID
			sort.Slice(indikators, func(i, j int) bool {
				return indikators[i].Id < indikators[j].Id
			})

			// Konversi indikator ke response
			var indikatorResponses []rincianbelanja.IndikatorResponse
			for _, ind := range indikators {
				// Sort target berdasarkan ID
				sort.Slice(ind.Target, func(i, j int) bool {
					return ind.Target[i].Id < ind.Target[j].Id
				})

				var targetResponses []rincianbelanja.TargetResponse
				for _, t := range ind.Target {
					targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
						Id:          t.Id,
						IndikatorId: t.IndikatorId,
						Target:      t.Target,
						Satuan:      t.Satuan,
					})
				}

				indikatorResponses = append(indikatorResponses, rincianbelanja.IndikatorResponse{
					Id:               ind.Id,
					RencanaKinerjaId: ind.RencanaKinerjaId,
					NamaIndikator:    ind.Indikator,
					Target:           targetResponses,
				})
			}

			// Sort rencana aksi berdasarkan ID
			if rk.RencanaAksi != nil {
				for _, ra := range rk.RencanaAksi {
					rencanaAksiResponses = append(rencanaAksiResponses, rincianbelanja.RencanaAksiResponse{
						RenaksiId: ra.RenaksiId,
						Renaksi:   ra.Renaksi,
						Anggaran:  int(ra.Anggaran),
					})
					totalAnggaranRekin += int(ra.Anggaran)
				}
			}

			subResponse.RincianBelanja = append(subResponse.RincianBelanja, rincianbelanja.RincianBelanjaResponse{
				RencanaKinerjaId: rk.RencanaKinerjaId,
				RencanaKinerja:   rk.RencanaKinerja,
				PegawaiId:        rk.PegawaiId,
				NamaPegawai:      rk.NamaPegawai,
				Indikator:        indikatorResponses,
				TotalAnggaran:    totalAnggaranRekin,
				RencanaAksi:      rencanaAksiResponses,
			})
			subResponse.TotalAnggaran += totalAnggaranRekin
		}
	}

	// Convert map to slice dan sort berdasarkan kode subkegiatan
	var responses []rincianbelanja.RincianBelanjaAsnResponse
	for _, response := range subkegiatanMap {
		// Sort rincian belanja berdasarkan ID
		sort.Slice(response.RincianBelanja, func(i, j int) bool {
			return response.RincianBelanja[i].RencanaKinerjaId < response.RincianBelanja[j].RencanaKinerjaId
		})
		responses = append(responses, *response)
	}
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].KodeSubkegiatan < responses[j].KodeSubkegiatan
	})

	return responses, nil
}

func (service *RincianBelanjaServiceImpl) LaporanRincianBelanjaPegawai(ctx context.Context, pegawaiId string, tahun string) ([]rincianbelanja.RincianBelanjaAsnResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	rincianBelanjaList, err := service.rincianBelanjaRepository.LaporanRincianBelanjaPegawai(ctx, tx, pegawaiId, tahun)
	if err != nil {
		return nil, err
	}

	// Map untuk mengelompokkan berdasarkan kode OPD dan subkegiatan
	subkegiatanMap := make(map[string]*rincianbelanja.RincianBelanjaAsnResponse)

	for _, rb := range rincianBelanjaList {
		// Buat key unik untuk kombinasi OPD dan subkegiatan
		key := fmt.Sprintf("%s_%s", rb.KodeOpd, rb.KodeSubkegiatan)

		// Ambil atau buat response subkegiatan
		subResponse, exists := subkegiatanMap[key]
		if !exists {
			// Ambil indikator subkegiatan
			indikatorSubkegiatan, err := service.rincianBelanjaRepository.FindIndikatorSubkegiatanByKodeAndOpd(
				ctx,
				tx,
				rb.KodeSubkegiatan,
				rb.KodeOpd,
				tahun,
			)
			if err != nil {
				log.Printf("Error mengambil indikator subkegiatan: %v", err)
				continue
			}

			// Sort dan konversi indikator subkegiatan
			sort.Slice(indikatorSubkegiatan, func(i, j int) bool {
				return indikatorSubkegiatan[i].Id < indikatorSubkegiatan[j].Id
			})

			var indikatorSubkegiatanResponses []rincianbelanja.IndikatorResponse
			for _, ind := range indikatorSubkegiatan {
				sort.Slice(ind.Target, func(i, j int) bool {
					return ind.Target[i].Id < ind.Target[j].Id
				})

				var targetResponses []rincianbelanja.TargetResponse
				for _, t := range ind.Target {
					targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
						Id:          t.Id,
						IndikatorId: t.IndikatorId,
						Target:      t.Target,
						Satuan:      t.Satuan,
					})
				}

				indikatorSubkegiatanResponses = append(indikatorSubkegiatanResponses, rincianbelanja.IndikatorResponse{
					Id:              ind.Id,
					KodeSubkegiatan: ind.Kode,
					KodeOPD:         ind.KodeOpd,
					NamaIndikator:   ind.Indikator,
					Target:          targetResponses,
				})
			}

			subResponse = &rincianbelanja.RincianBelanjaAsnResponse{
				KodeOpd:              rb.KodeOpd,
				KodeSubkegiatan:      rb.KodeSubkegiatan,
				NamaSubkegiatan:      rb.NamaSubkegiatan,
				IndikatorSubkegiatan: indikatorSubkegiatanResponses,
				TotalAnggaran:        0,
				RincianBelanja:       []rincianbelanja.RincianBelanjaResponse{},
			}
			subkegiatanMap[key] = subResponse
		}

		// Proses rencana kinerja
		for _, rk := range rb.RencanaKinerja {
			var rencanaAksiResponses []rincianbelanja.RencanaAksiResponse
			var totalAnggaranRekin int = 0 // Sudah benar, di-reset untuk setiap rencana kinerja

			// Ambil dan proses indikator rencana kinerja
			indikators, err := service.rincianBelanjaRepository.FindIndikatorByRekinId(ctx, tx, rk.RencanaKinerjaId)
			if err != nil {
				log.Printf("Error mengambil indikator untuk rekin %s: %v", rk.RencanaKinerjaId, err)
				continue
			}

			// Sort dan konversi indikator
			var indikatorResponses []rincianbelanja.IndikatorResponse
			for _, ind := range indikators {
				sort.Slice(ind.Target, func(i, j int) bool {
					return ind.Target[i].Id < ind.Target[j].Id
				})

				var targetResponses []rincianbelanja.TargetResponse
				for _, t := range ind.Target {
					targetResponses = append(targetResponses, rincianbelanja.TargetResponse{
						Id:          t.Id,
						IndikatorId: t.IndikatorId,
						Target:      t.Target,
						Satuan:      t.Satuan,
					})
				}

				indikatorResponses = append(indikatorResponses, rincianbelanja.IndikatorResponse{
					Id:               ind.Id,
					RencanaKinerjaId: ind.RencanaKinerjaId,
					NamaIndikator:    ind.Indikator,
					Target:           targetResponses,
				})
			}

			// Proses rencana aksi - anggaran sudah di-SUM di query, langsung pakai
			if rk.RencanaAksi != nil {
				for _, ra := range rk.RencanaAksi {
					rencanaAksiResponses = append(rencanaAksiResponses, rincianbelanja.RencanaAksiResponse{
						RenaksiId: ra.RenaksiId,
						Renaksi:   ra.Renaksi,
						Anggaran:  int(ra.Anggaran), // Anggaran sudah benar dari query
					})
					totalAnggaranRekin += int(ra.Anggaran)
				}
			}

			subResponse.RincianBelanja = append(subResponse.RincianBelanja, rincianbelanja.RincianBelanjaResponse{
				RencanaKinerjaId: rk.RencanaKinerjaId,
				RencanaKinerja:   rk.RencanaKinerja,
				PegawaiId:        rk.PegawaiId,
				NamaPegawai:      rk.NamaPegawai,
				Indikator:        indikatorResponses,
				TotalAnggaran:    totalAnggaranRekin,
				RencanaAksi:      rencanaAksiResponses,
			})
			subResponse.TotalAnggaran += totalAnggaranRekin
		}

	}

	// Convert map to slice dan sort
	var responses []rincianbelanja.RincianBelanjaAsnResponse
	for _, response := range subkegiatanMap {
		responses = append(responses, *response)
	}
	for i := range responses {
		// Urutkan RincianBelanja sehingga rencana kinerja pegawai yang difilter muncul duluan
		sort.Slice(responses[i].RincianBelanja, func(x, y int) bool {
			// Jika salah satu adalah pegawai yang dicari, letakkan di atas
			isPegawaiX := responses[i].RincianBelanja[x].PegawaiId == pegawaiId
			isPegawaiY := responses[i].RincianBelanja[y].PegawaiId == pegawaiId

			if isPegawaiX != isPegawaiY {
				return isPegawaiX // true jika x adalah pegawai yang dicari
			}

			// Jika keduanya bukan pegawai yang dicari atau keduanya adalah pegawai yang dicari
			// urutkan berdasarkan RencanaKinerjaId
			return responses[i].RincianBelanja[x].RencanaKinerjaId < responses[i].RincianBelanja[y].RencanaKinerjaId
		})
	}

	// Sort berdasarkan kode OPD dan kode subkegiatan
	sort.Slice(responses, func(i, j int) bool {
		if responses[i].KodeOpd != responses[j].KodeOpd {
			return responses[i].KodeOpd < responses[j].KodeOpd
		}
		return responses[i].KodeSubkegiatan < responses[j].KodeSubkegiatan
	})

	return responses, nil
}
