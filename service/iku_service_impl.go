package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/iku"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type IkuServiceImpl struct {
	IkuRepository repository.IkuRepository
	DB            *sql.DB
}

func NewIkuServiceImpl(ikuRepository repository.IkuRepository, db *sql.DB) *IkuServiceImpl {
	return &IkuServiceImpl{
		IkuRepository: ikuRepository,
		DB:            db,
	}
}

func (service *IkuServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository dengan parameter baru
	indikatorTargets, err := service.IkuRepository.FindAll(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	// Transform ke response
	var responses []iku.IkuResponse
	for _, item := range indikatorTargets {
		var targetResponses []iku.TargetResponse
		for _, target := range item.Target {
			targetResponses = append(targetResponses, iku.TargetResponse{
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		responses = append(responses, iku.IkuResponse{
			IndikatorId: item.Id,
			Sumber:      item.Sumber,
			// IsActive:         item.IsActive,
			IkuActive:           item.IkuActive,
			Indikator:           item.Indikator,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			DefinisiOperasional: item.DefinisiOperasional.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			JenisPeriode:        item.JenisPeriode,
			Target:              targetResponses,
		})
	}

	return responses, nil
}

// ═════════════════════════════════════════════════════════════════
// V2 — IKU Pemda filter berdasarkan tematik.tahun
// ═════════════════════════════════════════════════════════════════

func hasRealIkuTarget(t domain.Target) bool {
	raw := strings.TrimSpace(t.Target)
	return raw != "" && raw != "-"
}

func emptyIkuTargetSlot(tahun string) iku.TargetResponse {
	return iku.TargetResponse{Target: "", Satuan: "", Tahun: tahun}
}

// singleIkuTargetForYear — 1 slot untuk tahun tematik; ranwal = potongan renstra (bukan jenis ranwal di DB).
func singleIkuTargetForYear(targets []domain.Target, tahun string) []iku.TargetResponse {
	for _, t := range targets {
		if t.Tahun == tahun && hasRealIkuTarget(t) {
			return []iku.TargetResponse{{Target: t.Target, Satuan: t.Satuan, Tahun: t.Tahun}}
		}
	}
	return []iku.TargetResponse{emptyIkuTargetSlot(tahun)}
}

func buildDomainTargetMap(items []domain.Indikator) map[string][]domain.Target {
	m := make(map[string][]domain.Target, len(items))
	for _, item := range items {
		m[item.Id] = item.Target
	}
	return m
}

func toIkuResponseSliceSingleYear(items []domain.Indikator, tahun string) []iku.IkuResponse {
	responses := make([]iku.IkuResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, iku.IkuResponse{
			IndikatorId:         item.Id,
			Sumber:              item.Sumber,
			IkuActive:           item.IkuActive,
			Indikator:           item.Indikator,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			DefinisiOperasional: item.DefinisiOperasional.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			JenisPeriode:        item.JenisPeriode,
			Target:              singleIkuTargetForYear(item.Target, tahun),
		})
	}
	return responses
}

// FindIkuPemdaRanwalV2 — potongan target renstra untuk tahun tematik (tanpa override jenis ranwal).
func (service *IkuServiceImpl) FindIkuPemdaRanwalV2(
	ctx context.Context, tahun, jenisPeriode string,
) ([]iku.IkuResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	items, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	return toIkuResponseSliceSingleYear(items, tahun), nil
}

// FindIkuPemdaRankhirDualV2 — target_ranwal = potongan renstra; target_rankhir = rankhir saja (kosong jika belum ada).
func (service *IkuServiceImpl) FindIkuPemdaRankhirDualV2(
	ctx context.Context, tahun, jenisPeriode string,
) ([]iku.IkuPemdaRankhirDualResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	baseItems, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	rankhirItems, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil {
		return nil, err
	}
	rankhirMap := buildDomainTargetMap(rankhirItems)

	responses := make([]iku.IkuPemdaRankhirDualResponse, 0, len(baseItems))
	for _, item := range baseItems {
		responses = append(responses, iku.IkuPemdaRankhirDualResponse{
			IndikatorId:         item.Id,
			Sumber:              item.Sumber,
			IkuActive:           item.IkuActive,
			Indikator:           item.Indikator,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			DefinisiOperasional: item.DefinisiOperasional.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			JenisPeriode:        item.JenisPeriode,
			TargetRanwal:        singleIkuTargetForYear(item.Target, tahun),
			TargetRankhir:       singleIkuTargetForYear(rankhirMap[item.Id], tahun),
		})
	}
	return responses, nil
}

// FindIkuPemdaPenetapanDualV2 — target_rankhir dan target_penetapan masing-masing dari layer-nya (tanpa saling isi).
func (service *IkuServiceImpl) FindIkuPemdaPenetapanDualV2(
	ctx context.Context, tahun, jenisPeriode string,
) ([]iku.IkuPemdaPenetapanDualResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	baseItems, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	rankhirItems, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil {
		return nil, err
	}
	penetapanItems, err := service.IkuRepository.FindAllPemdaByTematikTahun(ctx, tx, tahun, jenisPeriode, "penetapan")
	if err != nil {
		return nil, err
	}
	rankhirMap := buildDomainTargetMap(rankhirItems)
	penetapanMap := buildDomainTargetMap(penetapanItems)

	responses := make([]iku.IkuPemdaPenetapanDualResponse, 0, len(baseItems))
	for _, item := range baseItems {
		responses = append(responses, iku.IkuPemdaPenetapanDualResponse{
			IndikatorId:         item.Id,
			Sumber:              item.Sumber,
			IkuActive:           item.IkuActive,
			Indikator:           item.Indikator,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			DefinisiOperasional: item.DefinisiOperasional.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			JenisPeriode:        item.JenisPeriode,
			TargetRankhir:       singleIkuTargetForYear(rankhirMap[item.Id], tahun),
			TargetPenetapan:     singleIkuTargetForYear(penetapanMap[item.Id], tahun),
		})
	}
	return responses, nil
}

func (service *IkuServiceImpl) FindAllIkuOpd(ctx context.Context, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	indikators, err := service.getIndikatorWithFallback(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	var responses []iku.IkuOpdResponse
	for _, item := range indikators {
		var targetResponses []iku.TargetOpdResponse

		// Pastikan target terurut berdasarkan tahun
		sort.Slice(item.Target, func(i, j int) bool {
			tahunI, _ := strconv.Atoi(item.Target[i].Tahun)
			tahunJ, _ := strconv.Atoi(item.Target[j].Tahun)
			return tahunI < tahunJ
		})

		// Konversi semua target, termasuk yang kosong
		for _, target := range item.Target {
			targetResponses = append(targetResponses, iku.TargetOpdResponse{
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		responses = append(responses, iku.IkuOpdResponse{
			IndikatorId:         item.Id,
			AsalIku:             item.AsalIku,
			Indikator:           item.Indikator,
			IkuActive:           item.IkuActive,
			IsHide:              item.IsHide,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			DefinisiOperasional: item.DefinisiOperasional.String,
			JenisPeriode:        item.JenisPeriode,
			Target:              targetResponses,
		})
	}

	// Urutkan responses berdasarkan CreatedAt
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].CreatedAt.Before(responses[j].CreatedAt)
	})

	if len(responses) == 0 {
		responses = make([]iku.IkuOpdResponse, 0)
	}

	return responses, nil
}

func (service *IkuServiceImpl) UpdateIkuActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	// 1. coba update di tabel baru

	rows, err := service.IkuRepository.UpdateIkuOpdActive(ctx, tx, id, request.IsActive)
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	// 2. fallback ke tabel lama
	rows, err = service.IkuRepository.UpdateIkuActive(ctx, tx, id, request.IsActive)
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	// 3. benar-benar tidak ditemukan
	return errors.New("iku not found in both indikator and indikator_matrix")
}

func (service *IkuServiceImpl) UpdateIkuOpdActive(ctx context.Context, id string, request iku.IkuUpdateActiveRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	rows, err := service.IkuRepository.UpdateIkuOpdActive(ctx, tx, id, request.IsActive)
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	// 3. benar-benar tidak ditemukan
	return errors.New("iku not found")
}

func (service *IkuServiceImpl) FindAllIkuRenja(ctx context.Context, kodeOpd string, tahun string, jenisPeriode string, jenisIndikator string) ([]iku.IkuOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	indikators, err := service.IkuRepository.FindAllIkuRenja(ctx, tx, kodeOpd, tahun, jenisPeriode, jenisIndikator)
	if err != nil {
		return nil, err
	}
	var responses []iku.IkuOpdResponse
	for _, item := range indikators {
		var targetResponses []iku.TargetOpdResponse
		for _, target := range item.Target {
			targetResponses = append(targetResponses, iku.TargetOpdResponse{
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}
		responses = append(responses, iku.IkuOpdResponse{
			IndikatorId:         item.Id,
			AsalIku:             item.AsalIku,
			Indikator:           item.Indikator,
			IkuActive:           item.IkuActive,
			IsHide:              item.IsHide,
			RumusPerhitungan:    item.RumusPerhitungan.String,
			SumberData:          item.SumberData.String,
			CreatedAt:           item.CreatedAt,
			TahunAwal:           item.TahunAwal,
			TahunAkhir:          item.TahunAkhir,
			DefinisiOperasional: item.DefinisiOperasional.String,
			JenisPeriode:        item.JenisPeriode,
			Jenis:               jenisIndikator,
			Target:              targetResponses,
		})
	}
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].CreatedAt.Before(responses[j].CreatedAt)
	})
	if len(responses) == 0 {
		responses = make([]iku.IkuOpdResponse, 0)
	}
	return responses, nil
}

func (s *IkuServiceImpl) getIndikatorWithFallback(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	tahunAwal string,
	tahunAkhir string,
	jenisPeriode string,
) ([]domain.Indikator, error) {

	indikatorBaru, err := s.IkuRepository.FindAllIkuOpd(ctx, tx, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	indikatorLama, err := s.IkuRepository.
		FindAllIkuOpdOld(ctx, tx, kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	return mergeIndikator(indikatorBaru, indikatorLama), nil
}
