package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TujuanPemdaServiceImpl struct {
	TujuanPemdaRepository   repository.TujuanPemdaRepository
	PeriodeRepository       repository.PeriodeRepository
	PohonKinerjaRepository  repository.PohonKinerjaRepository
	VisiPemdaRepository     repository.VisiPemdaRepository
	MisiPemdaRepository     repository.MisiPemdaRepository
	LockDataPemdaRepository repository.LockDataPemdaRepository
	DB                      *sql.DB
}

func NewTujuanPemdaServiceImpl(tujuanPemdaRepository repository.TujuanPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, visiPemdaRepository repository.VisiPemdaRepository, misiPemdaRepository repository.MisiPemdaRepository, lockDataPemdaRepository repository.LockDataPemdaRepository, DB *sql.DB) *TujuanPemdaServiceImpl {
	return &TujuanPemdaServiceImpl{
		TujuanPemdaRepository:   tujuanPemdaRepository,
		PeriodeRepository:       periodeRepository,
		PohonKinerjaRepository:  pohonKinerjaRepository,
		VisiPemdaRepository:     visiPemdaRepository,
		MisiPemdaRepository:     misiPemdaRepository,
		LockDataPemdaRepository: lockDataPemdaRepository,
		DB:                      DB,
	}
}

// ─────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		id := rand.Intn(90000) + 10000
		if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}
func generateKodeIndikatorPemda() string {
	return fmt.Sprintf("IND-TJN-PMD-%s-%s", time.Now().Format("2006"), uuid.New().String()[:5])
}
func defaultJenisPemda(j string) string {
	if strings.TrimSpace(j) == "" {
		return "renstra"
	}
	return strings.TrimSpace(j)
}

func parseTargetFloat(raw string) float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return 0
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return f
}

// toTargetResponse — domain TargetPemda → web TargetResponse (target = int)
func toTargetResponse(t domain.TargetPemda) tujuanpemda.TargetResponse {
	return tujuanpemda.TargetResponse{
		Id:     t.Id,
		Target: tujuanpemda.NewTargetDisplayFromString(t.Target),
		Satuan: t.Satuan,
		Tahun:  t.Tahun,
		Jenis:  t.Jenis,
	}
}

// emptyTargetResponse — placeholder tahun yang belum punya data
func emptyTargetResponse(tahun, jenis string) tujuanpemda.TargetResponse {
	return tujuanpemda.TargetResponse{
		Id:     0,
		Target: tujuanpemda.NewTargetDisplayFromString("-"),
		Satuan: "-",
		Tahun:  tahun,
		Jenis:  jenis,
	}
}

// toIndikatorResponse — domain IndikatorPemda → web IndikatorResponse
func toIndikatorResponse(ind domain.IndikatorPemda) tujuanpemda.IndikatorResponse {
	targets := make([]tujuanpemda.TargetResponse, 0, len(ind.Target))
	for _, t := range ind.Target {
		targets = append(targets, toTargetResponse(t))
	}
	return tujuanpemda.IndikatorResponse{
		Id:                  ind.Id,
		KodeIndikator:       ind.KodeIndikator,
		Indikator:           ind.Indikator.String,
		RumusPerhitungan:    ind.RumusPerhitungan.String,
		SumberData:          ind.SumberData.String,
		DefinisiOperasional: ind.DefinisiOperasional.String,
		Jenis:               ind.Jenis,
		Target:              targets,
	}
}
func buildIndikatorResponses(indikators []domain.IndikatorPemda) []tujuanpemda.IndikatorResponse {
	responses := make([]tujuanpemda.IndikatorResponse, 0, len(indikators))
	for _, ind := range indikators {
		sort.Slice(ind.Target, func(i, j int) bool {
			return ind.Target[i].Tahun < ind.Target[j].Tahun
		})
		responses = append(responses, toIndikatorResponse(ind))
	}
	return responses
}
func validateTargetTahun(indikators []tujuanpemda.IndikatorCreateRequest, tahunAwal, tahunAkhir int) error {
	for _, ind := range indikators {
		tahunMap := make(map[string]bool)
		for _, t := range ind.Target {
			tt, _ := strconv.Atoi(t.Tahun)
			if tt < tahunAwal || tt > tahunAkhir {
				return fmt.Errorf("tahun target %d harus dalam rentang %d-%d", tt, tahunAwal, tahunAkhir)
			}
			if tahunMap[t.Tahun] {
				return fmt.Errorf("duplikasi tahun %s pada indikator %s", t.Tahun, ind.Indikator)
			}
			tahunMap[t.Tahun] = true
		}
	}
	return nil
}
func validateTargetTahunUpdate(indikators []tujuanpemda.IndikatorUpdateRequest, tahunAwal, tahunAkhir int) error {
	for _, ind := range indikators {
		tahunMap := make(map[string]bool)
		for _, t := range ind.Target {
			tt, _ := strconv.Atoi(t.Tahun)
			if tt < tahunAwal || tt > tahunAkhir {
				return fmt.Errorf("tahun target %d harus dalam rentang %d-%d", tt, tahunAwal, tahunAkhir)
			}
			if tahunMap[t.Tahun] {
				return fmt.Errorf("duplikasi tahun %s pada indikator %s", t.Tahun, ind.Indikator)
			}
			tahunMap[t.Tahun] = true
		}
	}
	return nil
}

func validateTargetValue(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return nil // placeholder kosong, boleh
	}
	if _, err := strconv.Atoi(raw); err != nil {
		return fmt.Errorf("target '%s' harus berupa angka bulat", raw)
	}
	return nil
}

func targetToDBString(v float64) (string, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "", fmt.Errorf("target format tidak valid")
	}
	if v == math.Trunc(v) {
		return strconv.FormatInt(int64(v), 10), nil
	}
	return strconv.FormatFloat(v, 'f', -1, 64), nil
}

func validateTargetValuesCreate(indikators []tujuanpemda.IndikatorCreateRequest) error {
	for _, ind := range indikators {
		for _, t := range ind.Target {
			if _, err := targetToDBString(t.Target.Float64()); err != nil {
				return fmt.Errorf("indikator '%s' tahun %s: %w", ind.Indikator, t.Tahun, err)
			}
		}
	}
	return nil
}
func validateTargetValuesUpdate(indikators []tujuanpemda.IndikatorUpdateRequest) error {
	for _, ind := range indikators {
		for _, t := range ind.Target {
			if _, err := targetToDBString(t.Target.Float64()); err != nil {
				return fmt.Errorf("indikator '%s' tahun %s: %w", ind.Indikator, t.Tahun, err)
			}
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────
// CREATE
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) Create(
	ctx context.Context, request tujuanpemda.TujuanPemdaCreateRequest,
) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}
	if err = service.PohonKinerjaRepository.ValidatePokinLevel(ctx, tx, request.TematikId, 0, "tujuan pemda"); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)
	if err = validateTargetTahun(request.Indikator, tahunAwal, tahunAkhir); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	if err = validateTargetValuesCreate(request.Indikator); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, request.IdVisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, request.IdMisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	tp := domain.TujuanPemda{
		Id:                service.generateRandomId(ctx, tx),
		TujuanPemda:       request.TujuanPemda,
		IdVisi:            request.IdVisi,
		IdMisi:            request.IdMisi,
		TematikId:         request.TematikId,
		PeriodeId:         request.PeriodeId,
		TahunAwalPeriode:  periode.TahunAwal,
		TahunAkhirPeriode: periode.TahunAkhir,
		JenisPeriode:      periode.JenisPeriode,
	}
	tp, err = service.TujuanPemdaRepository.Create(ctx, tx, tp)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	for _, indReq := range request.Indikator {
		jenis := defaultJenisPemda(indReq.Jenis)
		kodeInd := generateKodeIndikatorPemda()
		ind := domain.IndikatorPemda{
			KodeIndikator:       kodeInd,
			TujuanPemdaId:       tp.Id,
			Indikator:           sql.NullString{String: indReq.Indikator, Valid: true},
			RumusPerhitungan:    sql.NullString{String: indReq.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: indReq.SumberData, Valid: true},
			DefinisiOperasional: sql.NullString{String: indReq.DefinisiOperasional, Valid: true},
			Jenis:               jenis,
		}
		if _, err = service.TujuanPemdaRepository.CreateIndikator(ctx, tx, ind); err != nil {
			return tujuanpemda.TujuanPemdaResponse{}, err
		}
		for _, tReq := range indReq.Target {
			targetStr, err := targetToDBString(tReq.Target.Float64())
			if err != nil {
				return tujuanpemda.TujuanPemdaResponse{}, err
			}
			tg := domain.TargetPemda{
				KodeIndikator: kodeInd,
				Target:        targetStr,
				Satuan:        tReq.Satuan,
				Tahun:         tReq.Tahun,
				Jenis:         defaultJenisPemda(tReq.Jenis),
			}
			if tg.Jenis == "" {
				tg.Jenis = jenis
			}
			if _, err = service.TujuanPemdaRepository.CreateTarget(ctx, tx, tg); err != nil {
				return tujuanpemda.TujuanPemdaResponse{}, err
			}
		}
	}
	result, err := service.TujuanPemdaRepository.FindById(ctx, tx, tp.Id)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	return tujuanpemda.TujuanPemdaResponse{
		Id:          result.Id,
		Visi:        visiPemda.Visi,
		Misi:        misiPemda.Misi,
		TujuanPemda: result.TujuanPemda,
		TematikId:   result.TematikId,
		NamaTematik: pokinData.NamaPohon,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    periode.TahunAwal,
			TahunAkhir:   periode.TahunAkhir,
			JenisPeriode: periode.JenisPeriode,
		},
		Indikator: buildIndikatorResponses(result.IndikatorPemda),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// UPDATE
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) Update(
	ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest,
) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	existing, err := service.TujuanPemdaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	// Periode tidak di-update di endpoint ini — pakai data existing untuk validasi tahun target
	// periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	// if err != nil {
	// 	return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	// }
	tahunAwal, _ := strconv.Atoi(existing.TahunAwalPeriode)
	tahunAkhir, _ := strconv.Atoi(existing.TahunAkhirPeriode)
	if err = validateTargetTahunUpdate(request.Indikator, tahunAwal, tahunAkhir); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	if err = validateTargetValuesUpdate(request.Indikator); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	existingIndMap := make(map[int]domain.IndikatorPemda)
	for _, ind := range existing.IndikatorPemda {
		existingIndMap[ind.Id] = ind
	}
	tp := domain.TujuanPemda{
		Id:          request.Id,
		TujuanPemda: request.TujuanPemda,
		IdVisi:      request.IdVisi,
		IdMisi:      request.IdMisi,
		TematikId:   request.TematikId,
		// PeriodeId, TahunAwalPeriode, TahunAkhirPeriode, JenisPeriode — tidak di-update, pertahankan dari existing
		PeriodeId:         existing.PeriodeId,
		TahunAwalPeriode:  existing.TahunAwalPeriode,
		TahunAkhirPeriode: existing.TahunAkhirPeriode,
		JenisPeriode:      existing.JenisPeriode,
		// PeriodeId:         request.PeriodeId,
		// TahunAwalPeriode:  periode.TahunAwal,
		// TahunAkhirPeriode: periode.TahunAkhir,
		// JenisPeriode:      periode.JenisPeriode,
	}
	for _, indReq := range request.Indikator {
		jenis := defaultJenisPemda(indReq.Jenis)
		kodeInd := strings.TrimSpace(indReq.KodeIndikator)
		ind := domain.IndikatorPemda{
			Id:                  indReq.IdIndikator,
			TujuanPemdaId:       request.Id,
			Indikator:           sql.NullString{String: indReq.Indikator, Valid: true},
			RumusPerhitungan:    sql.NullString{String: indReq.RumusPerhitungan, Valid: true},
			SumberData:          sql.NullString{String: indReq.SumberData, Valid: true},
			DefinisiOperasional: sql.NullString{String: indReq.DefinisiOperasional, Valid: true},
			Jenis:               jenis,
		}
		if indReq.IdIndikator > 0 {
			// UPDATE — pertahankan id & kode_indikator lama
			ex, ok := existingIndMap[indReq.IdIndikator]
			if !ok {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("indikator id %d tidak ditemukan", indReq.IdIndikator)
			}
			ind.KodeIndikator = ex.KodeIndikator
		} else {
			// INSERT baru
			if kodeInd == "" {
				kodeInd = generateKodeIndikatorPemda()
			}
			ind.KodeIndikator = kodeInd
		}
		for _, tReq := range indReq.Target {
			targetStr, err := targetToDBString(tReq.Target.Float64())
			if err != nil {
				return tujuanpemda.TujuanPemdaResponse{}, err
			}
			tgJenis := defaultJenisPemda(tReq.Jenis)
			if tgJenis == "" {
				tgJenis = jenis
			}
			ind.Target = append(ind.Target, domain.TargetPemda{
				Id:            tReq.Id,
				KodeIndikator: ind.KodeIndikator,
				Target:        targetStr,
				Satuan:        tReq.Satuan,
				Tahun:         tReq.Tahun,
				Jenis:         tgJenis,
			})
		}
		tp.IndikatorPemda = append(tp.IndikatorPemda, ind)
	}
	if _, err = service.TujuanPemdaRepository.Update(ctx, tx, tp); err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	result, err := service.TujuanPemdaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	return service.toTujuanPemdaResponse(ctx, tx, result)
}

// ─────────────────────────────────────────────────────────────────
// DELETE
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.TujuanPemdaRepository.Delete(ctx, tx, id)
}

// ─────────────────────────────────────────────────────────────────
// FIND BY ID
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) FindById(
	ctx context.Context, tujuanPemdaId int,
) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	result, err := service.TujuanPemdaRepository.FindById(ctx, tx, tujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	return service.toTujuanPemdaResponse(ctx, tx, result)
}

// ─────────────────────────────────────────────────────────────────
// FIND ALL
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) FindAll(
	ctx context.Context, tahun string, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := service.TujuanPemdaRepository.FindAll(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	responses := make([]tujuanpemda.TujuanPemdaResponse, 0, len(list))
	for _, tp := range list {
		pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tp.TematikId)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil pohon kinerja: %v", err)
		}
		responses = append(responses, tujuanpemda.TujuanPemdaResponse{
			Id:          tp.Id,
			TujuanPemda: tp.TujuanPemda,
			TematikId:   tp.TematikId,
			NamaTematik: pokinData.NamaPohon,
			Periode: tujuanpemda.PeriodeResponse{
				TahunAwal:    tp.Periode.TahunAwal,
				TahunAkhir:   tp.Periode.TahunAkhir,
				JenisPeriode: tp.Periode.JenisPeriode,
			},
		})
	}
	return responses, nil
}

// ─────────────────────────────────────────────────────────────────
// UPDATE PERIODE
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) UpdatePeriode(
	ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest,
) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.Id) {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("tujuan pemda id %d tidak ditemukan", request.Id)
	}
	if request.PeriodeId != 0 && !service.PeriodeRepository.IsIdExists(ctx, tx, request.PeriodeId) {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode id %d tidak ditemukan", request.PeriodeId)
	}
	result, err := service.TujuanPemdaRepository.UpdatePeriode(ctx, tx, domain.TujuanPemda{
		Id: request.Id, PeriodeId: request.PeriodeId,
	})
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, result.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	return tujuanpemda.TujuanPemdaResponse{
		Id:          result.Id,
		TujuanPemda: result.TujuanPemda,
		TematikId:   result.TematikId,
		NamaTematik: pokinData.NamaPohon,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  result.Periode.TahunAwal,
			TahunAkhir: result.Periode.TahunAkhir,
		},
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// FIND ALL WITH POKIN (semua jenis)
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) FindAllWithPokin(
	ctx context.Context, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := service.TujuanPemdaRepository.FindAllWithPokin(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}
	return service.buildPokinResponse(ctx, tx, list, tahunAwal, tahunAkhir, jenisPeriode)
}

// ─────────────────────────────────────────────────────────────────
// FIND ALL WITH POKIN RENSTRA (jenis = renstra, 5 tahunan)
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) FindAllWithPokinRenstra(
	ctx context.Context, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := service.TujuanPemdaRepository.FindAllWithPokinRenstra(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}
	return service.buildPokinResponse(ctx, tx, list, tahunAwal, tahunAkhir, jenisPeriode)
}

// ─────────────────────────────────────────────────────────────────
// FIND POKIN WITH PERIODE (tidak berubah)
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) FindPokinWithPeriode(
	ctx context.Context, pokinId int, jenisPeriode string,
) (tujuanpemda.PokinWithPeriodeResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	if err = service.PohonKinerjaRepository.ValidatePokinId(ctx, tx, pokinId); err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}
	if jenisPeriode == "" {
		return tujuanpemda.PokinWithPeriodeResponse{}, fmt.Errorf("jenis periode tidak boleh kosong")
	}
	pokin, periode, err := service.PohonKinerjaRepository.FindPokinWithPeriode(ctx, tx, pokinId, jenisPeriode)
	if err != nil {
		return tujuanpemda.PokinWithPeriodeResponse{}, err
	}
	response := tujuanpemda.PokinWithPeriodeResponse{
		Id:         pokin.Id,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		Tahun:      pokin.Tahun,
		Status:     pokin.Status,
		Periode: tujuanpemda.PokinPeriodeResponse{
			Id:         periode.Id,
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: []tujuanpemda.PokinIndikatorResponse{},
	}
	for _, ind := range pokin.Indikator {
		indResp := tujuanpemda.PokinIndikatorResponse{
			Id:               ind.Id,
			Indikator:        ind.Indikator,
			RumusPerhitungan: ind.RumusPerhitungan.String,
			SumberData:       ind.SumberData.String,
			Target:           []tujuanpemda.PokinTargetResponse{},
		}
		for _, t := range ind.Target {
			indResp.Target = append(indResp.Target, tujuanpemda.PokinTargetResponse{
				Id: t.Id, Target: t.Target, Satuan: t.Satuan, Tahun: t.Tahun,
			})
		}
		response.Indikator = append(response.Indikator, indResp)
	}
	return response, nil
}

// ─────────────────────────────────────────────────────────────────
// PRIVATE HELPERS
// ─────────────────────────────────────────────────────────────────
func (service *TujuanPemdaServiceImpl) toTujuanPemdaResponse(
	ctx context.Context, tx *sql.Tx, tp domain.TujuanPemda,
) (tujuanpemda.TujuanPemdaResponse, error) {
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tp.TematikId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal ambil pohon kinerja: %v", err)
	}
	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, tp.IdVisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, tp.IdMisi)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	return tujuanpemda.TujuanPemdaResponse{
		Id:          tp.Id,
		IdVisi:      tp.IdVisi,
		Visi:        visiPemda.Visi,
		IdMisi:      tp.IdMisi,
		Misi:        misiPemda.Misi,
		TujuanPemda: tp.TujuanPemda,
		TematikId:   tp.TematikId,
		NamaTematik: pokinData.NamaPohon,
		JenisPohon:  tp.JenisPohon,
		PeriodeId:   tp.PeriodeId,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    tp.Periode.TahunAwal,
			TahunAkhir:   tp.Periode.TahunAkhir,
			JenisPeriode: tp.Periode.JenisPeriode,
		},
		Indikator: buildIndikatorResponses(tp.IndikatorPemda),
	}, nil
}

func hasRealTarget(t domain.TargetPemda) bool {
	raw := strings.TrimSpace(t.Target)
	return raw != "" && raw != "-"
}

// buildIndikatorResponsesForPeriod — indikator + target 5 tahunan (placeholder jika kosong)
func buildIndikatorResponsesForPeriod(
	indikators []domain.IndikatorPemda,
	tahunAwal, tahunAkhir int,
) []tujuanpemda.IndikatorResponse {
	responses := make([]tujuanpemda.IndikatorResponse, 0, len(indikators))
	for _, ind := range indikators {
		if ind.KodeIndikator == "" && !ind.Indikator.Valid {
			continue
		}
		jenis := ind.Jenis
		if jenis == "" {
			jenis = "renstra"
		}
		targetMap := make(map[string]domain.TargetPemda)
		for _, t := range ind.Target {
			targetMap[strings.TrimSpace(t.Tahun)] = t
		}
		targetResponses := make([]tujuanpemda.TargetResponse, 0, tahunAkhir-tahunAwal+1)
		for y := tahunAwal; y <= tahunAkhir; y++ {
			yStr := strconv.Itoa(y)
			if t, ok := targetMap[yStr]; ok && hasRealTarget(t) {
				targetResponses = append(targetResponses, toTargetResponse(t))
			} else {
				targetResponses = append(targetResponses, emptyTargetResponse(yStr, jenis))
			}
		}
		responses = append(responses, tujuanpemda.IndikatorResponse{
			Id:                  ind.Id,
			KodeIndikator:       ind.KodeIndikator,
			Indikator:           ind.Indikator.String,
			RumusPerhitungan:    ind.RumusPerhitungan.String,
			SumberData:          ind.SumberData.String,
			DefinisiOperasional: ind.DefinisiOperasional.String,
			Jenis:               jenis,
			Target:              targetResponses,
		})
	}
	sort.Slice(responses, func(i, j int) bool {
		if responses[i].Id != responses[j].Id {
			return responses[i].Id < responses[j].Id
		}
		return responses[i].KodeIndikator < responses[j].KodeIndikator
	})
	return responses
}
func (service *TujuanPemdaServiceImpl) buildTujuanPemdaPokinResponse(
	ctx context.Context,
	tx *sql.Tx,
	tp domain.TujuanPemda,
	tahunAwal, tahunAkhir, jenisPeriode string,
	tahunAwalInt, tahunAkhirInt int,
) tujuanpemda.TujuanPemdaResponse {
	visiPemda, _ := service.VisiPemdaRepository.FindById(ctx, tx, tp.IdVisi)
	if visiPemda.Id == 0 {
		visiPemda.Visi = "Belum ada visi"
	}
	misiPemda, _ := service.MisiPemdaRepository.FindById(ctx, tx, tp.IdMisi)
	if misiPemda.Id == 0 {
		misiPemda.Misi = "Belum ada misi"
	}
	return tujuanpemda.TujuanPemdaResponse{
		Id:          tp.Id,
		IdVisi:      visiPemda.Id,
		Visi:        visiPemda.Visi,
		IdMisi:      misiPemda.Id,
		Misi:        misiPemda.Misi,
		TujuanPemda: tp.TujuanPemda,
		TematikId:   tp.TematikId,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:    tahunAwal,
			TahunAkhir:   tahunAkhir,
			JenisPeriode: jenisPeriode,
		},
		Indikator: buildIndikatorResponsesForPeriod(tp.IndikatorPemda, tahunAwalInt, tahunAkhirInt),
	}
}

func (service *TujuanPemdaServiceImpl) buildPokinResponse(
	ctx context.Context,
	tx *sql.Tx,
	list []domain.TujuanPemdaWithPokin,
	tahunAwal, tahunAkhir, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaWithPokinResponse, error) {
	tahunAwalInt, err := strconv.Atoi(tahunAwal)
	if err != nil {
		return nil, fmt.Errorf("format tahun awal tidak valid")
	}
	tahunAkhirInt, err := strconv.Atoi(tahunAkhir)
	if err != nil {
		return nil, fmt.Errorf("format tahun akhir tidak valid")
	}
	pokinMap := make(map[int]tujuanpemda.TujuanPemdaWithPokinResponse)
	for _, item := range list {
		pokinResp, exists := pokinMap[item.PokinId]
		if !exists {
			pokinResp = tujuanpemda.TujuanPemdaWithPokinResponse{
				PokinId:     item.PokinId,
				NamaPohon:   item.NamaPohon,
				JenisPohon:  item.JenisPohon,
				LevelPohon:  item.LevelPohon,
				Keterangan:  item.Keterangan,
				TahunPokin:  item.TahunPokin,
				IsActive:    item.IsActive,
				TujuanPemda: make([]tujuanpemda.TujuanPemdaResponse, 0),
			}
		}
		// deduplikasi tujuan per pokin (hindari append ganda)
		tujuanMap := make(map[int]tujuanpemda.TujuanPemdaResponse)
		for _, existing := range pokinResp.TujuanPemda {
			tujuanMap[existing.Id] = existing
		}
		for _, tp := range item.TujuanPemda {
			if tp.Id == 0 {
				continue
			}
			if _, already := tujuanMap[tp.Id]; already {
				continue
			}
			tujuanResp := service.buildTujuanPemdaPokinResponse(
				ctx, tx, tp,
				tahunAwal, tahunAkhir, jenisPeriode,
				tahunAwalInt, tahunAkhirInt,
			)
			tujuanMap[tp.Id] = tujuanResp
		}
		pokinResp.TujuanPemda = make([]tujuanpemda.TujuanPemdaResponse, 0, len(tujuanMap))
		for _, tpResp := range tujuanMap {
			pokinResp.TujuanPemda = append(pokinResp.TujuanPemda, tpResp)
		}
		sort.Slice(pokinResp.TujuanPemda, func(i, j int) bool {
			return pokinResp.TujuanPemda[i].Id < pokinResp.TujuanPemda[j].Id
		})
		pokinMap[item.PokinId] = pokinResp
	}
	result := make([]tujuanpemda.TujuanPemdaWithPokinResponse, 0, len(pokinMap))
	for _, pokinResp := range pokinMap {
		result = append(result, pokinResp)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PokinId < result[j].PokinId
	})
	return result, nil
}

// ═════════════════════════════════════════════════════════════════
// GUARD CHAIN layer RKPD — langsung []domain.TujuanPemda (tanpa pokin)
// Pola sama dengan TujuanOpd: renstra → ranwal → rankhir → penetapan
// ═════════════════════════════════════════════════════════════════
// applyTargetOverrideTujuanPemda menimpa target layer atas ke base.
// Key: tujuanId + kodeIndikator (1 tahun saja per indikator).
func applyTargetOverrideTujuanPemda(base, override []domain.TujuanPemda) []domain.TujuanPemda {
	type key struct {
		tujuanId      int
		kodeIndikator string
	}
	lookup := make(map[key]domain.TargetPemda)
	for _, tp := range override {
		for _, ind := range tp.IndikatorPemda {
			for _, tg := range ind.Target {
				if hasRealTarget(tg) {
					lookup[key{tp.Id, ind.KodeIndikator}] = tg
					break
				}
			}
		}
	}
	if len(lookup) == 0 {
		return base
	}
	result := make([]domain.TujuanPemda, len(base))
	for i, tp := range base {
		newTP := tp
		newTP.IndikatorPemda = make([]domain.IndikatorPemda, len(tp.IndikatorPemda))
		for j, ind := range tp.IndikatorPemda {
			newInd := ind
			if tg, ok := lookup[key{tp.Id, ind.KodeIndikator}]; ok {
				newInd.Target = []domain.TargetPemda{tg}
			} else {
				newInd.Target = make([]domain.TargetPemda, len(ind.Target))
				copy(newInd.Target, ind.Target)
			}
			newTP.IndikatorPemda[j] = newInd
		}
		result[i] = newTP
	}
	return result
}
func (s *TujuanPemdaServiceImpl) loadLayerRenstra(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string) ([]domain.TujuanPemda, error) {
	return s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
}
func (s *TujuanPemdaServiceImpl) loadLayerRanwal(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string) ([]domain.TujuanPemda, error) {
	base, err := s.loadLayerRenstra(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	ov, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "ranwal")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverrideTujuanPemda(base, ov), nil
}
func (s *TujuanPemdaServiceImpl) loadLayerRankhir(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string) ([]domain.TujuanPemda, error) {
	base, err := s.loadLayerRanwal(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	ov, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverrideTujuanPemda(base, ov), nil
}
func (s *TujuanPemdaServiceImpl) loadLayerPenetapan(ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string) ([]domain.TujuanPemda, error) {
	base, err := s.loadLayerRankhir(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	ov, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "penetapan")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverrideTujuanPemda(base, ov), nil
}

// buildTujuanPemdaListResponse — domain []TujuanPemda → web []TujuanPemdaResponse
func (s *TujuanPemdaServiceImpl) buildTujuanPemdaListResponse(
	ctx context.Context, tx *sql.Tx, list []domain.TujuanPemda,
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	responses := make([]tujuanpemda.TujuanPemdaResponse, 0, len(list))
	for _, tp := range list {
		resp, err := s.toTujuanPemdaResponse(ctx, tx, tp)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// Helper privat — 1 func untuk semua Find layer
func (s *TujuanPemdaServiceImpl) findByLayerTahun(
	ctx context.Context, tahun, jenisPeriode string,
	loader func(context.Context, *sql.Tx, string, string) ([]domain.TujuanPemda, error),
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	if len(strings.TrimSpace(tahun)) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid, contoh: 2025")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := loader(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	return s.buildTujuanPemdaListResponse(ctx, tx, list)
}
func (s *TujuanPemdaServiceImpl) FindTujuanPemdaRanwal(
	ctx context.Context, tahun, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	return s.findByLayerTahun(ctx, tahun, jenisPeriode, s.loadLayerRanwal)
}
func (s *TujuanPemdaServiceImpl) FindTujuanPemdaRankhir(
	ctx context.Context, tahun, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	return s.findByLayerTahun(ctx, tahun, jenisPeriode, s.loadLayerRankhir)
}
func (s *TujuanPemdaServiceImpl) FindTujuanPemdaPenetapan(
	ctx context.Context, tahun, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaResponse, error) {
	return s.findByLayerTahun(ctx, tahun, jenisPeriode, s.loadLayerPenetapan)
}

// ─────────────────────────────────────────────────────────────────
// FUNC BARU: UpsertTargetPemdaLayer
// POST /tujuan_pemda/:jenis/target/upsert
// Menyimpan/memperbarui target di layer ranwal, rankhir, atau penetapan.
// Tidak mengubah metadata indikator (nama, rumus, dll.) — hanya target.
// ─────────────────────────────────────────────────────────────────
func (s *TujuanPemdaServiceImpl) UpsertTargetPemdaLayer(ctx context.Context, jenis string, req tujuanpemda.LayerTargetBatchRequest) ([]tujuanpemda.TargetResponse, error) {
	jenis = strings.TrimSpace(jenis)
	if jenis != "ranwal" && jenis != "rankhir" && jenis != "penetapan" {
		return nil, fmt.Errorf("jenis layer tidak valid: '%s'. Gunakan: ranwal, rankhir, atau penetapan", jenis)
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]tujuanpemda.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		if strings.TrimSpace(item.KodeIndikator) == "" {
			return nil, fmt.Errorf("kode_indikator tidak boleh kosong")
		}
		if strings.TrimSpace(item.Tahun) == "" {
			return nil, fmt.Errorf("tahun tidak boleh kosong untuk kode_indikator %s", item.KodeIndikator)
		}
		// Pastikan indikator renstra-nya ada
		_, err := s.TujuanPemdaRepository.FindIndikatorPemdaByKode(ctx, tx, item.KodeIndikator)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("indikator '%s' tidak ditemukan di renstra", item.KodeIndikator)
		}
		if err != nil {
			return nil, err
		}
		targetStr, err := targetToDBString(item.Target.Float64())
		if err != nil {
			return nil, fmt.Errorf("nilai target tidak valid untuk %s tahun %s: %w", item.KodeIndikator, item.Tahun, err)
		}
		saved, err := s.TujuanPemdaRepository.UpsertTargetPemda(ctx, tx, domain.TargetPemda{
			KodeIndikator: item.KodeIndikator,
			Target:        targetStr,
			Satuan:        item.Satuan,
			Tahun:         strings.TrimSpace(item.Tahun),
			Jenis:         jenis,
		})
		if err != nil {
			return nil, err
		}
		responses = append(responses, toTargetResponse(saved))
	}
	return responses, nil
}

// ═════════════════════════════════════════════════════════════════
// OPSI B — DUAL TARGET (func baru, terpisah dari guard chain)
//
// Rankhir Dual  → target: ranwal + rankhir  (2 slot per indikator)
// Penetapan Dual→ target: rankhir + penetapan (2 slot per indikator)
// Tidak ada fallback antar jenis; slot kosong = placeholder "-"
// ═════════════════════════════════════════════════════════════════
type tujuanIndikatorKey struct {
	tujuanId      int
	kodeIndikator string
}

// fillTargetSlotForJenis — isi slot target jenis tertentu dari data DB.
func fillTargetSlotForJenis(
	base *[]domain.TujuanPemda,
	layerData []domain.TujuanPemda,
	jenis string,
) {
	lookup := make(map[tujuanIndikatorKey]domain.TargetPemda)
	for _, tp := range layerData {
		for _, ind := range tp.IndikatorPemda {
			for _, tg := range ind.Target {
				if strings.TrimSpace(tg.Jenis) == jenis && hasRealTarget(tg) {
					lookup[tujuanIndikatorKey{tp.Id, ind.KodeIndikator}] = tg
				}
			}
		}
	}
	if len(lookup) == 0 {
		return
	}
	for i := range *base {
		tp := &(*base)[i]
		for j := range tp.IndikatorPemda {
			ind := &tp.IndikatorPemda[j]
			k := tujuanIndikatorKey{tp.Id, ind.KodeIndikator}
			tg, ok := lookup[k]
			if !ok {
				continue
			}
			for t := range ind.Target {
				if ind.Target[t].Jenis == jenis {
					ind.Target[t] = tg
					break
				}
			}
		}
	}
}

// loadTujuanPemdaWithDualTargets — skeleton indikator dari renstra,
// lalu isi slot target per jenis yang diminta.
func (s *TujuanPemdaServiceImpl) loadTujuanPemdaWithDualTargets(
	ctx context.Context, tx *sql.Tx,
	tahun, jenisPeriode string,
	jenisList []string,
) ([]domain.TujuanPemda, error) {
	base, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	// Buat slot kosong per jenis (tanpa fallback)
	for i := range base {
		for j := range base[i].IndikatorPemda {
			kode := base[i].IndikatorPemda[j].KodeIndikator
			slots := make([]domain.TargetPemda, 0, len(jenisList))
			for _, jenis := range jenisList {
				slots = append(slots, domain.TargetPemda{
					Id:            0,
					KodeIndikator: kode,
					Target:        "-",
					Satuan:        "-",
					Tahun:         tahun,
					Jenis:         jenis,
				})
			}
			base[i].IndikatorPemda[j].Target = slots
		}
	}
	// Isi dari DB per jenis
	for _, jenis := range jenisList {
		layerData, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, jenis)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		fillTargetSlotForJenis(&base, layerData, jenis)
	}
	return base, nil
}

// loadLayerRankhirDual — ranwal + rankhir
func (s *TujuanPemdaServiceImpl) loadLayerRankhirDual(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string,
) ([]domain.TujuanPemda, error) {
	renstra, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	for i := range renstra {
		for j := range renstra[i].IndikatorPemda {
			kode := renstra[i].IndikatorPemda[j].KodeIndikator
			// Slot 1: ranwal = renstra (label "ranwal" untuk UI)
			var ranwalSlot domain.TargetPemda
			if len(renstra[i].IndikatorPemda[j].Target) > 0 && hasRealTarget(renstra[i].IndikatorPemda[j].Target[0]) {
				tg := renstra[i].IndikatorPemda[j].Target[0]
				ranwalSlot = domain.TargetPemda{
					Id:            tg.Id,
					KodeIndikator: kode,
					Target:        tg.Target,
					Satuan:        tg.Satuan,
					Tahun:         tg.Tahun,
					Jenis:         "ranwal",
				}
			} else {
				ranwalSlot = domain.TargetPemda{
					Id: 0, KodeIndikator: kode,
					Target: "-", Satuan: "-", Tahun: tahun, Jenis: "ranwal",
				}
			}
			// Slot 2: rankhir placeholder
			rankhirSlot := domain.TargetPemda{
				Id: 0, KodeIndikator: kode,
				Target: "-", Satuan: "-", Tahun: tahun, Jenis: "rankhir",
			}
			renstra[i].IndikatorPemda[j].Target = []domain.TargetPemda{ranwalSlot, rankhirSlot}
		}
	}
	// Isi slot rankhir dari DB
	rankhirData, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillTargetSlotForJenis(&renstra, rankhirData, "rankhir")
	return renstra, nil
}

// loadLayerPenetapanDual — rankhir + penetapan
func (s *TujuanPemdaServiceImpl) loadLayerPenetapanDual(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string,
) ([]domain.TujuanPemda, error) {
	base, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	for i := range base {
		for j := range base[i].IndikatorPemda {
			kode := base[i].IndikatorPemda[j].KodeIndikator
			base[i].IndikatorPemda[j].Target = []domain.TargetPemda{
				{Id: 0, KodeIndikator: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "rankhir"},
				{Id: 0, KodeIndikator: kode, Target: "-", Satuan: "-", Tahun: tahun, Jenis: "penetapan"},
			}
		}
	}
	rankhirData, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillTargetSlotForJenis(&base, rankhirData, "rankhir")
	penetapanData, err := s.TujuanPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "penetapan")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	fillTargetSlotForJenis(&base, penetapanData, "penetapan")
	return base, nil
}

func validateLayerJenis(jenis string) (string, error) {
	jenis = strings.TrimSpace(jenis)
	if jenis != "ranwal" && jenis != "rankhir" && jenis != "penetapan" {
		return "", fmt.Errorf("jenis layer tidak valid: '%s'. Gunakan: ranwal, rankhir, atau penetapan", jenis)
	}
	return jenis, nil
}

// CreateTargetPemdaLayer — POST create target layer ranwal/rankhir/penetapan
// Hanya insert baru. Gagal jika target sudah ada (kode_indikator + tahun + jenis).
func (s *TujuanPemdaServiceImpl) CreateTargetPemdaLayer(
	ctx context.Context, jenis string, req tujuanpemda.LayerTargetBatchRequest,
) ([]tujuanpemda.TargetResponse, error) {
	var err error
	jenis, err = validateLayerJenis(jenis)
	if err != nil {
		return nil, err
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]tujuanpemda.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		kode := strings.TrimSpace(item.KodeIndikator)
		tahun := strings.TrimSpace(item.Tahun)
		if kode == "" {
			return nil, fmt.Errorf("kode_indikator tidak boleh kosong")
		}
		if tahun == "" {
			return nil, fmt.Errorf("tahun tidak boleh kosong untuk kode_indikator %s", kode)
		}
		// Indikator renstra harus sudah ada
		_, err := s.TujuanPemdaRepository.FindIndikatorPemdaByKode(ctx, tx, kode)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("indikator '%s' tidak ditemukan di renstra", kode)
		}
		if err != nil {
			return nil, err
		}
		// Tolak jika target sudah ada
		exists, err := s.TujuanPemdaRepository.TargetPemdaExistsByKey(ctx, tx, kode, tahun, jenis)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf(
				"target untuk kode_indikator '%s' tahun %s jenis %s sudah ada, gunakan update",
				kode, tahun, jenis,
			)
		}
		targetStr, err := targetToDBString(item.Target.Float64())
		if err != nil {
			return nil, fmt.Errorf("nilai target tidak valid untuk %s tahun %s: %w", kode, tahun, err)
		}
		saved, err := s.TujuanPemdaRepository.CreateTarget(ctx, tx, domain.TargetPemda{
			KodeIndikator: kode,
			Target:        targetStr,
			Satuan:        item.Satuan,
			Tahun:         tahun,
			Jenis:         jenis,
		})
		if err != nil {
			return nil, err
		}
		responses = append(responses, toTargetResponse(saved))
	}
	return responses, nil
}

func (s *TujuanPemdaServiceImpl) UpdateTargetPemdaLayer(
	ctx context.Context, jenis string, req tujuanpemda.LayerTargetUpdateBatchRequest,
) ([]tujuanpemda.TargetResponse, error) {
	var err error
	jenis, err = validateLayerJenis(jenis)
	if err != nil {
		return nil, err
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]tujuanpemda.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		if item.Id == 0 {
			return nil, fmt.Errorf("id target wajib diisi")
		}
		// Ambil data existing dari DB
		existing, err := s.TujuanPemdaRepository.FindTargetPemdaById(ctx, tx, item.Id)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target id %d tidak ditemukan", item.Id)
		}
		if err != nil {
			return nil, err
		}
		// Pastikan jenis di DB cocok dengan endpoint — jenis tidak boleh diubah
		if existing.Jenis != jenis {
			return nil, fmt.Errorf(
				"target id %d adalah jenis '%s', tidak bisa diupdate via endpoint '%s'",
				item.Id, existing.Jenis, jenis,
			)
		}
		targetStr, err := targetToDBString(item.Target.Float64())
		if err != nil {
			return nil, fmt.Errorf("nilai target tidak valid untuk id %d: %w", item.Id, err)
		}
		saved, err := s.TujuanPemdaRepository.UpdateTargetPemda(ctx, tx, item.Id, targetStr, item.Satuan)
		if err != nil {
			return nil, err
		}
		responses = append(responses, toTargetResponse(saved))
	}
	return responses, nil
}

// ── Public: Opsi B ──────────────────────────────────────────────
func (s *TujuanPemdaServiceImpl) FindTujuanPemdaRankhirDual(
	ctx context.Context, tahun, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaRankhirDualResponse, error) {
	if len(strings.TrimSpace(tahun)) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid, contoh: 2025")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := s.loadLayerRankhirDual(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	responses := make([]tujuanpemda.TujuanPemdaRankhirDualResponse, 0, len(list))
	for _, tp := range list {
		resp, err := s.toTujuanPemdaRankhirDualResponse(ctx, tx, tp)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}
func (s *TujuanPemdaServiceImpl) FindTujuanPemdaPenetapanDual(
	ctx context.Context, tahun, jenisPeriode string,
) ([]tujuanpemda.TujuanPemdaPenetapanDualResponse, error) {
	if len(strings.TrimSpace(tahun)) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid, contoh: 2025")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := s.loadLayerPenetapanDual(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	responses := make([]tujuanpemda.TujuanPemdaPenetapanDualResponse, 0, len(list))
	for _, tp := range list {
		resp, err := s.toTujuanPemdaPenetapanDualResponse(ctx, tx, tp)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// lock pemda
const lockJenisTujuanPemda = "tujuan_pemda"

// isTahunInPeriode — cek apakah tahun lock masuk range periode renstra
func isTahunInPeriode(tahun, tahunAwal, tahunAkhir string) bool {
	t, _ := strconv.Atoi(strings.TrimSpace(tahun))
	ta, _ := strconv.Atoi(strings.TrimSpace(tahunAwal))
	tb, _ := strconv.Atoi(strings.TrimSpace(tahunAkhir))
	return t >= ta && t <= tb
}

// isPeriodeOverlapLock — apakah periode tujuan pemda overlap dengan tahun yang di-lock
func (s *TujuanPemdaServiceImpl) isPeriodeOverlapLock(
	ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir string,
) (bool, string, error) {
	locks, err := s.LockDataPemdaRepository.FindAllByJenis(ctx, tx, lockJenisTujuanPemda)
	if err != nil {
		return false, "", err
	}
	for _, lock := range locks {
		if isTahunInPeriode(lock.Tahun, tahunAwal, tahunAkhir) {
			return true, lock.Tahun, nil
		}
	}
	return false, "", nil
}

// assertIndikatorNotLocked — blok create/update/delete indikator & tujuan pemda
func (s *TujuanPemdaServiceImpl) assertIndikatorNotLocked(
	ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir string,
) error {
	locked, tahun, err := s.isPeriodeOverlapLock(ctx, tx, tahunAwal, tahunAkhir)
	if err != nil {
		return err
	}
	if locked {
		return fmt.Errorf(
			"data tujuan pemda terkunci untuk tahun %s (periode %s-%s). Indikator dan penghapusan tidak diizinkan",
			tahun, tahunAwal, tahunAkhir,
		)
	}
	return nil
}

// assertTargetLayerAllowedWhenLocked — cek izin ubah target per jenis layer
func (s *TujuanPemdaServiceImpl) assertTargetLayerAllowedWhenLocked(
	ctx context.Context, tx *sql.Tx, tahun, targetJenis string,
) error {
	locked, err := s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisTujuanPemda, tahun)
	if err != nil {
		return err
	}
	if !locked {
		return nil // belum lock → semua jenis boleh
	}
	switch strings.TrimSpace(targetJenis) {
	case "renstra", "rankhir":
		return nil // ✅ boleh diubah
	case "ranwal", "penetapan":
		return fmt.Errorf(
			"target jenis '%s' tahun %s tidak dapat diubah karena data tujuan pemda terkunci",
			targetJenis, tahun,
		)
	default:
		return fmt.Errorf("jenis target '%s' tidak dikenali", targetJenis)
	}
}

// helper
func toTargetDualResponse(t domain.TargetPemda) tujuanpemda.TargetDualResponse {
	return tujuanpemda.TargetDualResponse{
		Id:     t.Id,
		Target: tujuanpemda.NewTargetDisplayFromString(t.Target),
		Satuan: t.Satuan,
		Tahun:  t.Tahun,
	}
}
func toTargetDualSlice(t domain.TargetPemda) []tujuanpemda.TargetDualResponse {
	return []tujuanpemda.TargetDualResponse{toTargetDualResponse(t)}
}
func findTargetByJenis(ind domain.IndikatorPemda, jenis string) domain.TargetPemda {
	for _, t := range ind.Target {
		if t.Jenis == jenis {
			return t
		}
	}
	return domain.TargetPemda{}
}
func (s *TujuanPemdaServiceImpl) toTujuanPemdaRankhirDualResponse(
	ctx context.Context, tx *sql.Tx, tp domain.TujuanPemda,
) (tujuanpemda.TujuanPemdaRankhirDualResponse, error) {
	base, err := s.toTujuanPemdaResponse(ctx, tx, tp) // reuse header (visi, misi, tematik)
	if err != nil {
		return tujuanpemda.TujuanPemdaRankhirDualResponse{}, err
	}
	indikators := make([]tujuanpemda.IndikatorRankhirDualResponse, 0, len(tp.IndikatorPemda))
	for _, ind := range tp.IndikatorPemda {
		ranwal := findTargetByJenis(ind, "ranwal")
		rankhir := findTargetByJenis(ind, "rankhir")
		indikators = append(indikators, tujuanpemda.IndikatorRankhirDualResponse{
			Id:                  ind.Id,
			KodeIndikator:       ind.KodeIndikator,
			Indikator:           ind.Indikator.String,
			RumusPerhitungan:    ind.RumusPerhitungan.String,
			SumberData:          ind.SumberData.String,
			DefinisiOperasional: ind.DefinisiOperasional.String,
			Jenis:               ind.Jenis,
			TargetRanwal:        toTargetDualSlice(ranwal),
			TargetRankhir:       toTargetDualSlice(rankhir),
		})
	}
	return tujuanpemda.TujuanPemdaRankhirDualResponse{
		Id: base.Id, IdVisi: base.IdVisi, Visi: base.Visi,
		IdMisi: base.IdMisi, Misi: base.Misi,
		TujuanPemda: base.TujuanPemda, TematikId: base.TematikId,
		NamaTematik: base.NamaTematik, Periode: base.Periode,
		Indikator: indikators,
	}, nil
}
func (s *TujuanPemdaServiceImpl) toTujuanPemdaPenetapanDualResponse(
	ctx context.Context, tx *sql.Tx, tp domain.TujuanPemda,
) (tujuanpemda.TujuanPemdaPenetapanDualResponse, error) {
	base, err := s.toTujuanPemdaResponse(ctx, tx, tp)
	if err != nil {
		return tujuanpemda.TujuanPemdaPenetapanDualResponse{}, err
	}
	indikators := make([]tujuanpemda.IndikatorPenetapanDualResponse, 0, len(tp.IndikatorPemda))
	for _, ind := range tp.IndikatorPemda {
		rankhir := findTargetByJenis(ind, "rankhir")
		penetapan := findTargetByJenis(ind, "penetapan")
		indikators = append(indikators, tujuanpemda.IndikatorPenetapanDualResponse{
			Id:                  ind.Id,
			KodeIndikator:       ind.KodeIndikator,
			Indikator:           ind.Indikator.String,
			RumusPerhitungan:    ind.RumusPerhitungan.String,
			SumberData:          ind.SumberData.String,
			DefinisiOperasional: ind.DefinisiOperasional.String,
			Jenis:               ind.Jenis,
			TargetRankhir:       toTargetDualSlice(rankhir),
			TargetPenetapan:     toTargetDualSlice(penetapan),
		})
	}
	return tujuanpemda.TujuanPemdaPenetapanDualResponse{
		Id: base.Id, IdVisi: base.IdVisi, Visi: base.Visi,
		IdMisi: base.IdMisi, Misi: base.Misi,
		TujuanPemda: base.TujuanPemda, TematikId: base.TematikId,
		NamaTematik: base.NamaTematik, Periode: base.Periode,
		Indikator: indikators,
	}, nil
}

// ── Public API lock ─────────────────────────────────────────────
func (s *TujuanPemdaServiceImpl) LockTujuanPemda(ctx context.Context, tahun string) error {
	if len(strings.TrimSpace(tahun)) != 4 {
		return fmt.Errorf("format tahun tidak valid")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return s.LockDataPemdaRepository.Lock(ctx, tx, lockJenisTujuanPemda, tahun)
}
func (s *TujuanPemdaServiceImpl) UnlockTujuanPemda(ctx context.Context, tahun string) error {
	if len(strings.TrimSpace(tahun)) != 4 {
		return fmt.Errorf("format tahun tidak valid")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return s.LockDataPemdaRepository.Unlock(ctx, tx, lockJenisTujuanPemda, tahun)
}
func (s *TujuanPemdaServiceImpl) IsTujuanPemdaLocked(ctx context.Context, tahun string) (bool, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx)
	return s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisTujuanPemda, tahun)
}
func (s *TujuanPemdaServiceImpl) FindAllLockTujuanPemda(
	ctx context.Context,
) ([]tujuanpemda.LockDataPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	locks, err := s.LockDataPemdaRepository.FindAllByJenis(ctx, tx, lockJenisTujuanPemda)
	if err != nil {
		return nil, err
	}
	result := make([]tujuanpemda.LockDataPemdaResponse, 0, len(locks))
	for _, l := range locks {
		result = append(result, tujuanpemda.LockDataPemdaResponse{
			Id: l.Id, Jenis: l.Jenis, Tahun: l.Tahun, Locked: true,
		})
	}
	return result, nil
}
