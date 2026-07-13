package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/sasaranpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const lockJenisSasaranPemda = "sasaran_pemda"

type SasaranPemdaServiceImpl struct {
	SasaranPemdaRepository  repository.SasaranPemdaRepository
	PeriodeRepository       repository.PeriodeRepository
	PohonKinerjaRepository  repository.PohonKinerjaRepository
	TujuanPemdaRepository   repository.TujuanPemdaRepository
	LockDataPemdaRepository repository.LockDataPemdaRepository
	DB                      *sql.DB
	Validator               *validator.Validate
}

func NewSasaranPemdaServiceImpl(
	sasaranPemdaRepository repository.SasaranPemdaRepository,
	periodeRepository repository.PeriodeRepository,
	pohonKinerjaRepository repository.PohonKinerjaRepository,
	tujuanPemdaRepository repository.TujuanPemdaRepository,
	lockDataPemdaRepository repository.LockDataPemdaRepository,
	db *sql.DB,
	validator *validator.Validate,
) *SasaranPemdaServiceImpl {
	return &SasaranPemdaServiceImpl{
		SasaranPemdaRepository:  sasaranPemdaRepository,
		PeriodeRepository:       periodeRepository,
		PohonKinerjaRepository:  pohonKinerjaRepository,
		TujuanPemdaRepository:   tujuanPemdaRepository,
		LockDataPemdaRepository: lockDataPemdaRepository,
		DB:                      db,
		Validator:               validator,
	}
}

// ── ID helper ────────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		id := rand.Intn(90000) + 10000
		if !s.SasaranPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}
func generateKodeIndikatorSasaran() string {
	return fmt.Sprintf("IND-SAS-PMD-%s-%s", time.Now().Format("2006"), uuid.New().String()[:5])
}

// ── Validasi ─────────────────────────────────────────────────────
func validateTargetValuesSasaranCreate(inds []sasaranpemda.IndikatorCreateRequest) error {
	for _, ind := range inds {
		for _, t := range ind.Target {
			if err := helper.ValidateTargetFloat(t.Target.Float64()); err != nil {
				return helper.ErrTargetIndikator(ind.Indikator, t.Tahun, err.Error())
			}
		}
	}
	return nil
}
func validateTargetValuesSasaranUpdate(inds []sasaranpemda.IndikatorUpdateRequest) error {
	for _, ind := range inds {
		for _, t := range ind.Target {
			if err := helper.ValidateTargetFloat(t.Target.Float64()); err != nil {
				return helper.ErrTargetIndikator(ind.Indikator, t.Tahun, err.Error())
			}
		}
	}
	return nil
}
func validateTargetTahunSasaran(inds []sasaranpemda.IndikatorCreateRequest, tahunAwal, tahunAkhir int) error {
	for _, ind := range inds {
		m := map[string]bool{}
		for _, t := range ind.Target {
			v, _ := strconv.Atoi(t.Tahun)
			if v < tahunAwal || v > tahunAkhir {
				return fmt.Errorf("tahun target %d harus dalam rentang %d-%d", v, tahunAwal, tahunAkhir)
			}
			if m[t.Tahun] {
				return fmt.Errorf("duplikasi tahun %s pada indikator %s", t.Tahun, ind.Indikator)
			}
			m[t.Tahun] = true
		}
	}
	return nil
}
func validateTargetTahunSasaranUpdate(inds []sasaranpemda.IndikatorUpdateRequest, tahunAwal, tahunAkhir int) error {
	for _, ind := range inds {
		m := map[string]bool{}
		for _, t := range ind.Target {
			v, _ := strconv.Atoi(t.Tahun)
			if v < tahunAwal || v > tahunAkhir {
				return fmt.Errorf("tahun target %d harus dalam rentang %d-%d", v, tahunAwal, tahunAkhir)
			}
			if m[t.Tahun] {
				return fmt.Errorf("duplikasi tahun %s pada indikator %s", t.Tahun, ind.Indikator)
			}
			m[t.Tahun] = true
		}
	}
	return nil
}

// ── Response builder ─────────────────────────────────────────────
func toTargetPemdaSlice(targets []domain.TargetPemda) []sasaranpemda.TargetResponse {
	resp := make([]sasaranpemda.TargetResponse, 0, len(targets))
	for _, t := range targets {
		resp = append(resp, sasaranpemda.TargetResponse{
			Id:     t.Id,
			Target: sasaranpemda.NewTargetDisplayFromString(t.Target),
			Satuan: t.Satuan, Tahun: t.Tahun, Jenis: t.Jenis,
		})
	}
	sort.Slice(resp, func(i, j int) bool {
		ti, _ := strconv.Atoi(resp[i].Tahun)
		tj, _ := strconv.Atoi(resp[j].Tahun)
		return ti < tj
	})
	return resp
}
func toIndikatorResponsesSasaran(inds []domain.IndikatorPemda) []sasaranpemda.IndikatorResponse {
	resp := make([]sasaranpemda.IndikatorResponse, 0, len(inds))
	for _, ind := range inds {
		resp = append(resp, sasaranpemda.IndikatorResponse{
			Id:                  ind.Id,
			KodeIndikator:       ind.KodeIndikator,
			Indikator:           ind.Indikator.String,
			RumusPerhitungan:    ind.RumusPerhitungan.String,
			SumberData:          ind.SumberData.String,
			DefinisiOperasional: ind.DefinisiOperasional.String,
			Target:              toTargetPemdaSlice(ind.Target),
		})
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i].Id < resp[j].Id })
	return resp
}

// ── Lock Guard ───────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) assertSasaranNotLocked(
	ctx context.Context, tx *sql.Tx, tahun string,
) error {
	locked, err := s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisSasaranPemda, tahun)
	if err != nil {
		return err
	}
	if locked {
		return fmt.Errorf("data sasaran pemda terkunci untuk tahun %s, modifikasi tidak diizinkan", tahun)
	}
	return nil
}
func (s *SasaranPemdaServiceImpl) assertTargetLayerAllowed(
	ctx context.Context, tx *sql.Tx, tahun, targetJenis string,
) error {
	locked, err := s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisSasaranPemda, tahun)
	if err != nil {
		return err
	}
	if !locked {
		return nil
	}
	switch strings.TrimSpace(targetJenis) {
	case "renstra", "rankhir":
		return nil
	case "penetapan":
		return fmt.Errorf("target jenis '%s' tahun %s tidak dapat diubah karena data sasaran pemda terkunci",
			targetJenis, tahun)
	default:
		return fmt.Errorf("jenis target '%s' tidak dikenali", targetJenis)
	}
}

// ── CREATE ───────────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) Create(
	ctx context.Context, request sasaranpemda.SasaranPemdaCreateRequest,
) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	pokinData, err := s.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus 1-3, saat ini: %d", pokinData.LevelPohon)
	}
	if !s.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda id %d tidak ditemukan", request.TujuanPemdaId)
	}
	periode, err := s.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}
	if err := s.assertSasaranNotLocked(ctx, tx, periode.TahunAwal); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)
	if err := validateTargetTahunSasaran(request.Indikator, tahunAwal, tahunAkhir); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	if err := validateTargetValuesSasaranCreate(request.Indikator); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	sp := domain.SasaranPemda{
		Id:            s.generateRandomId(ctx, tx),
		SubtemaId:     request.SubtemaId,
		TujuanPemdaId: request.TujuanPemdaId,
		SasaranPemda:  request.SasaranPemda,
		PeriodeId:     request.PeriodeId,
		TahunAwal:     periode.TahunAwal,
		TahunAkhir:    periode.TahunAkhir,
		JenisPeriode:  periode.JenisPeriode,
	}
	for _, indReq := range request.Indikator {
		kodeInd := generateKodeIndikatorSasaran()
		ind := domain.IndikatorPemda{
			// Id = 0: DB auto-increment
			KodeIndikator:    kodeInd,
			SasaranPemdaId:   sp.Id,
			Indikator:        sql.NullString{String: indReq.Indikator, Valid: true},
			RumusPerhitungan: sql.NullString{String: indReq.RumusPerhitungan, Valid: true},
			SumberData:       sql.NullString{String: indReq.SumberData, Valid: true},
			Jenis:            "renstra",
		}
		for _, tReq := range indReq.Target {
			targetStr, err := helper.TargetToDBString(tReq.Target.Float64())
			if err != nil {
				return sasaranpemda.SasaranPemdaResponse{}, helper.ErrTargetIndikator(indReq.Indikator, tReq.Tahun, err.Error())
			}
			ind.Target = append(ind.Target, domain.TargetPemda{
				KodeIndikator: kodeInd,
				Target:        targetStr,
				Satuan:        tReq.Satuan,
				Tahun:         tReq.Tahun,
				Jenis:         "renstra",
			})
		}
		sp.Indikator = append(sp.Indikator, ind)
	}
	result, err := s.SasaranPemdaRepository.Create(ctx, tx, sp)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal membuat sasaran pemda: %v", err)
	}
	// Reload untuk mendapatkan id DB yang sebenarnya
	reloaded, err := s.SasaranPemdaRepository.FindById(ctx, tx, result.Id)
	if err != nil {
		reloaded = result
	}
	return sasaranpemda.SasaranPemdaResponse{
		Id: reloaded.Id, TujuanPemdaId: reloaded.TujuanPemdaId,
		SubtemaId: reloaded.SubtemaId, NamaSubtema: pokinData.NamaPohon,
		SasaranPemda: reloaded.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			Id: periode.Id, TahunAwal: periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir, JenisPeriode: periode.JenisPeriode,
		},
		Indikator: toIndikatorResponsesSasaran(reloaded.Indikator),
	}, nil
}

// ── UPDATE ───────────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) Update(
	ctx context.Context, request sasaranpemda.SasaranPemdaUpdateRequest,
) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	existing, err := s.SasaranPemdaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	if err := s.assertSasaranNotLocked(ctx, tx, existing.Periode.TahunAwal); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	if !s.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda id %d tidak ditemukan", request.TujuanPemdaId)
	}
	pokinData, err := s.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus 1-3, saat ini: %d", pokinData.LevelPohon)
	}
	tahunAwal, _ := strconv.Atoi(existing.Periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(existing.Periode.TahunAkhir)
	if err := validateTargetTahunSasaranUpdate(request.Indikator, tahunAwal, tahunAkhir); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	if err := validateTargetValuesSasaranUpdate(request.Indikator); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	// Buat lookup existing indikator: DB id → IndikatorPemda
	existingIndMap := make(map[int]domain.IndikatorPemda)
	for _, ind := range existing.Indikator {
		existingIndMap[ind.Id] = ind
	}
	existing.TujuanPemdaId = request.TujuanPemdaId
	existing.SubtemaId = request.SubtemaId
	existing.SasaranPemda = request.SasaranPemda
	var indikators []domain.IndikatorPemda
	for _, indReq := range request.Indikator {
		ind := domain.IndikatorPemda{
			SasaranPemdaId:   existing.Id,
			Indikator:        sql.NullString{String: indReq.Indikator, Valid: true},
			RumusPerhitungan: sql.NullString{String: indReq.RumusPerhitungan, Valid: true},
			SumberData:       sql.NullString{String: indReq.SumberData, Valid: true},
			Jenis:            "renstra",
		}
		if indReq.IdIndikator > 0 {
			// UPDATE — kode_indikator wajib dari DB (tidak boleh diubah dari request)
			ex, ok := existingIndMap[indReq.IdIndikator]
			if !ok {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("indikator id %d tidak ditemukan", indReq.IdIndikator)
			}
			ind.Id = indReq.IdIndikator
			ind.KodeIndikator = ex.KodeIndikator
		} else {
			// INSERT baru — generate kode baru
			kode := strings.TrimSpace(indReq.KodeIndikator)
			if kode == "" {
				kode = generateKodeIndikatorSasaran()
			}
			ind.KodeIndikator = kode
		}
		for _, tReq := range indReq.Target {
			targetStr, err := helper.TargetToDBString(tReq.Target.Float64())
			if err != nil {
				return sasaranpemda.SasaranPemdaResponse{}, helper.ErrTargetIndikator(indReq.Indikator, tReq.Tahun, err.Error())
			}
			ind.Target = append(ind.Target, domain.TargetPemda{
				Id:            tReq.Id, // 0 = baru
				KodeIndikator: ind.KodeIndikator,
				Target:        targetStr,
				Satuan:        tReq.Satuan,
				Tahun:         tReq.Tahun,
				Jenis:         "renstra",
			})
		}
		indikators = append(indikators, ind)
	}
	existing.Indikator = indikators
	if _, err := s.SasaranPemdaRepository.Update(ctx, tx, existing); err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	// Reload agar id DB target/indikator baru tersimpan
	reloaded, err := s.SasaranPemdaRepository.FindById(ctx, tx, existing.Id)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	return sasaranpemda.SasaranPemdaResponse{
		Id: reloaded.Id, TujuanPemdaId: reloaded.TujuanPemdaId,
		SubtemaId: reloaded.SubtemaId, NamaSubtema: pokinData.NamaPohon,
		SasaranPemda: reloaded.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			Id: reloaded.PeriodeId, TahunAwal: reloaded.Periode.TahunAwal,
			TahunAkhir: reloaded.Periode.TahunAkhir, JenisPeriode: reloaded.Periode.JenisPeriode,
		},
		Indikator: toIndikatorResponsesSasaran(reloaded.Indikator),
	}, nil
}

// ── DELETE ───────────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	existing, err := s.SasaranPemdaRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if err := s.assertSasaranNotLocked(ctx, tx, existing.Periode.TahunAwal); err != nil {
		return err
	}
	return s.SasaranPemdaRepository.Delete(ctx, tx, id)
}

// ── FIND ─────────────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) FindById(ctx context.Context, id int) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	sp, err := s.SasaranPemdaRepository.FindById(ctx, tx, id)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	pokinData, err := s.PohonKinerjaRepository.FindById(ctx, tx, sp.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal ambil pohon kinerja: %v", err)
	}
	tujuanPemda, err := s.TujuanPemdaRepository.FindById(ctx, tx, sp.TujuanPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal ambil tujuan pemda: %v", err)
	}
	return sasaranpemda.SasaranPemdaResponse{
		Id: sp.Id, TujuanPemdaId: sp.TujuanPemdaId,
		TujuanPemda: tujuanPemda.TujuanPemda,
		SubtemaId:   sp.SubtemaId, NamaSubtema: pokinData.NamaPohon,
		SasaranPemda: sp.SasaranPemda, JenisPohon: sp.JenisPohon,
		Periode: sasaranpemda.PeriodeResponse{
			Id: sp.PeriodeId, TahunAwal: sp.Periode.TahunAwal,
			TahunAkhir: sp.Periode.TahunAkhir, JenisPeriode: sp.Periode.JenisPeriode,
		},
		Indikator: toIndikatorResponsesSasaran(sp.Indikator),
	}, nil
}
func (s *SasaranPemdaServiceImpl) FindAll(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := s.SasaranPemdaRepository.FindAll(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return []sasaranpemda.SasaranPemdaResponse{}, nil
	}
	// Batch load pohon kinerja — 1 query, bukan N
	subtemaIds := make([]int, 0, len(list))
	for _, sp := range list {
		subtemaIds = append(subtemaIds, sp.SubtemaId)
	}
	pokinMap, err := s.PohonKinerjaRepository.FindByIds(ctx, tx, subtemaIds)
	if err != nil {
		return nil, fmt.Errorf("gagal batch load pohon kinerja: %v", err)
	}
	resp := make([]sasaranpemda.SasaranPemdaResponse, 0, len(list))
	for _, sp := range list {
		resp = append(resp, sasaranpemda.SasaranPemdaResponse{
			Id:           sp.Id,
			SubtemaId:    sp.SubtemaId,
			NamaSubtema:  pokinMap[sp.SubtemaId].NamaPohon,
			SasaranPemda: sp.SasaranPemda,
			Indikator:    toIndikatorResponsesSasaran(sp.Indikator),
		})
	}
	return resp, nil
}
func (s *SasaranPemdaServiceImpl) FindAllWithPokin(
	ctx context.Context, tahunAwal, tahunAkhir, jenisPeriode string,
) ([]sasaranpemda.TematikResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	var count int
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_periode WHERE tahun_awal=? AND tahun_akhir=? AND jenis_periode=?`,
		tahunAwal, tahunAkhir, jenisPeriode,
	).Scan(&count); err != nil {
		return nil, fmt.Errorf("error validasi periode: %v", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("periode %s-%s %s tidak ditemukan", tahunAwal, tahunAkhir, jenisPeriode)
	}
	isLocked, _ := s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisSasaranPemda, tahunAwal)
	pokinData, err := s.SasaranPemdaRepository.FindAllWithPokin(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}
	result := make([]sasaranpemda.TematikResponse, 0, len(pokinData))
	for _, tematik := range pokinData {
		tematikResp := sasaranpemda.TematikResponse{
			TematikId: tematik.TematikId, NamaTematik: tematik.NamaTematik,
			Tahun: tematik.Tahun, IsLock: isLocked,
			Subtematik: []sasaranpemda.SubtematikResponse{},
		}
		for _, sub := range tematik.Subtematik {
			if sub.LevelPohon < 1 || sub.LevelPohon > 3 {
				continue
			}
			subResp := sasaranpemda.SubtematikResponse{
				SubtematikId: sub.SubtematikId, NamaSubtematik: sub.NamaSubtematik,
				JenisPohon: sub.JenisPohon, LevelPohon: sub.LevelPohon,
				Tahun: sub.Tahun, IsActive: sub.IsActive,
				SasaranPemda: []sasaranpemda.SasaranPemdaWithPokinResponse{},
			}
			for _, sasaran := range sub.SasaranPemdaList {
				sasaranResp := sasaranpemda.SasaranPemdaWithPokinResponse{
					IdSasaranPemda: sasaran.Id, SasaranPemda: sasaran.SasaranPemda,
					Periode: sasaranpemda.PeriodeResponse{
						TahunAwal: tahunAwal, TahunAkhir: tahunAkhir, JenisPeriode: jenisPeriode,
					},
					Indikator: []sasaranpemda.IndikatorSubtematikResponse{},
				}
				tAwal, _ := strconv.Atoi(tahunAwal)
				tAkhir, _ := strconv.Atoi(tahunAkhir)
				for _, ind := range sasaran.Indikator {
					indResp := sasaranpemda.IndikatorSubtematikResponse{
						Id:               ind.Id,
						KodeIndikator:    ind.KodeIndikator,
						Indikator:        ind.Indikator,
						RumusPerhitungan: ind.RumusPerhitungan.String,
						SumberData:       ind.SumberData.String,
						Target:           []sasaranpemda.TargetResponse{},
					}
					existingTarget := make(map[string]domain.TargetDetail)
					for _, t := range ind.Target {
						existingTarget[t.Tahun] = t
					}
					for y := tAwal; y <= tAkhir; y++ {
						yStr := strconv.Itoa(y)
						if t, ok := existingTarget[yStr]; ok {
							indResp.Target = append(indResp.Target, sasaranpemda.TargetResponse{
								Id:     t.Id,
								Target: sasaranpemda.NewTargetDisplayFromString(t.Target),
								Satuan: t.Satuan, Tahun: yStr,
							})
						} else {
							indResp.Target = append(indResp.Target, sasaranpemda.TargetResponse{
								Id:     0,
								Target: sasaranpemda.NewTargetDisplayFromString(""),
								Satuan: "", Tahun: yStr,
							})
						}
					}
					sasaranResp.Indikator = append(sasaranResp.Indikator, indResp)
				}
				subResp.SasaranPemda = append(subResp.SasaranPemda, sasaranResp)
			}
			tematikResp.Subtematik = append(tematikResp.Subtematik, subResp)
		}
		if len(tematikResp.Subtematik) > 0 {
			result = append(result, tematikResp)
		}
	}
	return result, nil
}

// ranwal
func hasRealTargetSasaran(t domain.TargetPemda) bool {
	raw := strings.TrimSpace(t.Target)
	return raw != "" && raw != "-"
}
func applyTargetOverrideSasaran(base, override []domain.SasaranPemda) []domain.SasaranPemda {
	type key struct {
		sasaranId     int
		kodeIndikator string
	}
	lookup := make(map[key]domain.TargetPemda)
	for _, sp := range override {
		for _, ind := range sp.Indikator {
			for _, tg := range ind.Target {
				if hasRealTargetSasaran(tg) {
					lookup[key{sp.Id, ind.KodeIndikator}] = tg
					break
				}
			}
		}
	}
	if len(lookup) == 0 {
		return base
	}
	result := make([]domain.SasaranPemda, len(base))
	for i, sp := range base {
		newSP := sp
		newSP.Indikator = make([]domain.IndikatorPemda, len(sp.Indikator))
		for j, ind := range sp.Indikator {
			newInd := ind
			if tg, ok := lookup[key{sp.Id, ind.KodeIndikator}]; ok {
				newInd.Target = []domain.TargetPemda{tg}
			} else {
				newInd.Target = make([]domain.TargetPemda, len(ind.Target))
				copy(newInd.Target, ind.Target)
			}
			newSP.Indikator[j] = newInd
		}
		result[i] = newSP
	}
	return result
}
func (s *SasaranPemdaServiceImpl) loadLayerRenstra(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string,
) ([]domain.SasaranPemda, error) {
	return s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
}
func (s *SasaranPemdaServiceImpl) loadLayerRanwal(
	ctx context.Context, tx *sql.Tx, tahun, jenisPeriode string,
) ([]domain.SasaranPemda, error) {
	base, err := s.loadLayerRenstra(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	ov, err := s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "ranwal")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return applyTargetOverrideSasaran(base, ov), nil
}
func (s *SasaranPemdaServiceImpl) buildSasaranPemdaListResponse(
	ctx context.Context, tx *sql.Tx, list []domain.SasaranPemda,
) ([]sasaranpemda.SasaranPemdaResponse, error) {
	responses := make([]sasaranpemda.SasaranPemdaResponse, 0, len(list))
	for _, sp := range list {
		pokin, err := s.PohonKinerjaRepository.FindById(ctx, tx, sp.SubtemaId)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil pohon kinerja: %v", err)
		}
		tujuan, err := s.TujuanPemdaRepository.FindById(ctx, tx, sp.TujuanPemdaId)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil tujuan pemda: %v", err)
		}
		responses = append(responses, sasaranpemda.SasaranPemdaResponse{
			Id: sp.Id, TujuanPemdaId: sp.TujuanPemdaId,
			TujuanPemda: tujuan.TujuanPemda,
			SubtemaId:   sp.SubtemaId, NamaSubtema: pokin.NamaPohon,
			SasaranPemda: sp.SasaranPemda,
			Periode: sasaranpemda.PeriodeResponse{
				Id:           sp.PeriodeId,
				TahunAwal:    sp.Periode.TahunAwal,
				TahunAkhir:   sp.Periode.TahunAkhir,
				JenisPeriode: sp.Periode.JenisPeriode,
			},
			Indikator: toIndikatorResponsesSasaran(sp.Indikator),
		})
	}
	return responses, nil
}

func (s *SasaranPemdaServiceImpl) FindSasaranPemdaRanwal(
	ctx context.Context, tahun, jenisPeriode string,
) ([]sasaranpemda.SasaranPemdaResponse, error) {
	if len(strings.TrimSpace(tahun)) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid, contoh: 2025")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	list, err := s.SasaranPemdaRepository.FindRanwalByTahun(ctx, tx, tahun, jenisPeriode)
	if err != nil {
		return nil, err
	}
	responses := make([]sasaranpemda.SasaranPemdaResponse, 0, len(list))
	for _, sp := range list {
		responses = append(responses, sasaranpemda.SasaranPemdaResponse{
			Id:            sp.Id,
			TujuanPemdaId: sp.TujuanPemdaId,
			TujuanPemda:   sp.TujuanPemdaText,
			SubtemaId:     sp.SubtemaId,
			NamaSubtema:   sp.NamaSubtema,
			SasaranPemda:  sp.SasaranPemda,
			Periode: sasaranpemda.PeriodeResponse{
				Id:           sp.PeriodeId,
				TahunAwal:    sp.Periode.TahunAwal,
				TahunAkhir:   sp.Periode.TahunAkhir,
				JenisPeriode: sp.Periode.JenisPeriode,
			},
			Indikator: toIndikatorResponsesSasaran(sp.Indikator),
		})
	}
	return responses, nil
}

// ── DUAL RANKHIR ─────────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) FindSasaranPemdaRankhirDual(
	ctx context.Context, tahun, jenisPeriode string,
) ([]sasaranpemda.SasaranPemdaRankhirDualResponse, error) {
	if strings.TrimSpace(tahun) == "" {
		return nil, fmt.Errorf("tahun tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	baseList, err := s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "renstra")
	if err != nil {
		return nil, err
	}
	rankhirList, err := s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil {
		return nil, err
	}
	type dualKey struct {
		sasaranId     int
		kodeIndikator string
	}
	rankhirMap := make(map[dualKey][]domain.TargetPemda)
	for _, sp := range rankhirList {
		for _, ind := range sp.Indikator {
			k := dualKey{sp.Id, ind.KodeIndikator}
			rankhirMap[k] = append(rankhirMap[k], ind.Target...)
		}
	}
	responses := make([]sasaranpemda.SasaranPemdaRankhirDualResponse, 0, len(baseList))
	for _, sp := range baseList {
		resp := sasaranpemda.SasaranPemdaRankhirDualResponse{
			Id:           sp.Id,
			SasaranPemda: sp.SasaranPemda,
			Periode: sasaranpemda.PeriodeResponse{
				TahunAwal:    sp.Periode.TahunAwal,
				TahunAkhir:   sp.Periode.TahunAkhir,
				JenisPeriode: sp.Periode.JenisPeriode,
			},
			Indikator: []sasaranpemda.IndikatorRankhirDualResponse{}, // [] bukan null
		}
		for _, ind := range sp.Indikator {
			k := dualKey{sp.Id, ind.KodeIndikator}
			resp.Indikator = append(resp.Indikator, sasaranpemda.IndikatorRankhirDualResponse{
				Id:                  ind.Id,
				KodeIndikator:       ind.KodeIndikator,
				Indikator:           ind.Indikator.String,
				RumusPerhitungan:    ind.RumusPerhitungan.String,
				SumberData:          ind.SumberData.String,
				DefinisiOperasional: ind.DefinisiOperasional.String,
				TargetRanwal:        toTargetPemdaSlice(ind.Target),
				TargetRankhir:       toTargetPemdaSlice(rankhirMap[k]),
			})
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// ── DUAL PENETAPAN ───────────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) FindSasaranPemdaPenetapanDual(
	ctx context.Context, tahun, jenisPeriode string,
) ([]sasaranpemda.SasaranPemdaPenetapanDualResponse, error) {
	if strings.TrimSpace(tahun) == "" {
		return nil, fmt.Errorf("tahun tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	rankhirList, err := s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "rankhir")
	if err != nil {
		return nil, err
	}
	penetapanList, err := s.SasaranPemdaRepository.FindAllByTahun(ctx, tx, tahun, jenisPeriode, "penetapan")
	if err != nil {
		return nil, err
	}
	type dualKey struct {
		sasaranId     int
		kodeIndikator string
	}
	penetapanMap := make(map[dualKey][]domain.TargetPemda)
	for _, sp := range penetapanList {
		for _, ind := range sp.Indikator {
			k := dualKey{sp.Id, ind.KodeIndikator}
			penetapanMap[k] = append(penetapanMap[k], ind.Target...)
		}
	}
	responses := make([]sasaranpemda.SasaranPemdaPenetapanDualResponse, 0, len(rankhirList))
	for _, sp := range rankhirList {
		resp := sasaranpemda.SasaranPemdaPenetapanDualResponse{
			Id:           sp.Id,
			SasaranPemda: sp.SasaranPemda,
			Periode: sasaranpemda.PeriodeResponse{
				TahunAwal:    sp.Periode.TahunAwal,
				TahunAkhir:   sp.Periode.TahunAkhir,
				JenisPeriode: sp.Periode.JenisPeriode,
			},
			Indikator: []sasaranpemda.IndikatorPenetapanDualResponse{}, // [] bukan null
		}
		for _, ind := range sp.Indikator {
			k := dualKey{sp.Id, ind.KodeIndikator}
			resp.Indikator = append(resp.Indikator, sasaranpemda.IndikatorPenetapanDualResponse{
				Id:                  ind.Id,
				KodeIndikator:       ind.KodeIndikator,
				Indikator:           ind.Indikator.String,
				RumusPerhitungan:    ind.RumusPerhitungan.String,
				SumberData:          ind.SumberData.String,
				DefinisiOperasional: ind.DefinisiOperasional.String,
				TargetRankhir:       toTargetPemdaSlice(ind.Target),
				TargetPenetapan:     toTargetPemdaSlice(penetapanMap[k]),
			})
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// ── CREATE TARGET LAYER ──────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) CreateTargetSasaranLayer(
	ctx context.Context, jenis string, req sasaranpemda.LayerTargetBatchRequest,
) ([]sasaranpemda.TargetResponse, error) {
	jenis = strings.ToLower(strings.TrimSpace(jenis))
	if jenis != "rankhir" && jenis != "penetapan" {
		return nil, fmt.Errorf("jenis harus 'rankhir' atau 'penetapan'")
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]sasaranpemda.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		kode := strings.TrimSpace(item.KodeIndikator)
		if kode == "" {
			return nil, fmt.Errorf("kode_indikator tidak boleh kosong")
		}
		if strings.TrimSpace(item.Tahun) == "" {
			return nil, fmt.Errorf("tahun tidak boleh kosong untuk kode_indikator %s", kode)
		}
		if err := s.assertTargetLayerAllowed(ctx, tx, item.Tahun, jenis); err != nil {
			return nil, err
		}
		if _, err := s.SasaranPemdaRepository.FindIndikatorByKode(ctx, tx, kode); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("indikator kode '%s' tidak ditemukan", kode)
			}
			return nil, err
		}
		targetStr, err := helper.TargetToDBString(item.Target.Float64())
		if err != nil {
			return nil, helper.ErrTargetLayer(kode, item.Tahun, err.Error())
		}
		saved, err := s.SasaranPemdaRepository.UpsertTargetPemda(ctx, tx, domain.TargetPemda{
			KodeIndikator: kode,
			Target:        targetStr,
			Satuan:        item.Satuan,
			Tahun:         strings.TrimSpace(item.Tahun),
			Jenis:         jenis,
		})
		if err != nil {
			return nil, err
		}
		responses = append(responses, sasaranpemda.TargetResponse{
			Id:     saved.Id,
			Target: sasaranpemda.NewTargetDisplayFromString(saved.Target),
			Satuan: saved.Satuan, Tahun: saved.Tahun, Jenis: saved.Jenis,
		})
	}
	return responses, nil
}

// ── UPDATE TARGET LAYER ──────────────────────────────────────────
func (s *SasaranPemdaServiceImpl) UpdateTargetSasaranLayer(
	ctx context.Context, jenis string, req sasaranpemda.LayerTargetUpdateBatchRequest,
) ([]sasaranpemda.TargetResponse, error) {
	jenis = strings.ToLower(strings.TrimSpace(jenis))
	if jenis != "rankhir" && jenis != "penetapan" {
		return nil, fmt.Errorf("jenis harus 'rankhir' atau 'penetapan'")
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("targets tidak boleh kosong")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	responses := make([]sasaranpemda.TargetResponse, 0, len(req.Targets))
	for _, item := range req.Targets {
		if item.Id == 0 {
			return nil, fmt.Errorf("id target wajib diisi (bukan 0)")
		}
		existing, err := s.SasaranPemdaRepository.FindTargetLayerById(ctx, tx, item.Id)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target id %d tidak ditemukan", item.Id)
		}
		if err != nil {
			return nil, err
		}
		if existing.Jenis != jenis {
			return nil, fmt.Errorf(
				"target id %d adalah jenis '%s', tidak bisa diupdate via endpoint '%s'",
				item.Id, existing.Jenis, jenis,
			)
		}
		if err := s.assertTargetLayerAllowed(ctx, tx, existing.Tahun, jenis); err != nil {
			return nil, err
		}
		targetStr, err := helper.TargetToDBString(item.Target.Float64())
		if err != nil {
			return nil, helper.ErrTargetLayer(existing.KodeIndikator, existing.Tahun, err.Error())
		}
		saved, err := s.SasaranPemdaRepository.UpdateTargetLayerById(ctx, tx, item.Id, targetStr, item.Satuan)
		if err != nil {
			return nil, err
		}
		responses = append(responses, sasaranpemda.TargetResponse{
			Id:     saved.Id,
			Target: sasaranpemda.NewTargetDisplayFromString(saved.Target),
			Satuan: saved.Satuan, Tahun: saved.Tahun, Jenis: saved.Jenis,
		})
	}
	return responses, nil
}

// ── LOCK / UNLOCK ────────────────────────────────────────────────
func toLockSasaranResponse(l domain.LockDataPemda, locked bool) sasaranpemda.LockDataPemdaResponse {
	return sasaranpemda.LockDataPemdaResponse{Id: l.Id, Jenis: l.Jenis, Tahun: l.Tahun, Locked: locked}
}
func (s *SasaranPemdaServiceImpl) LockSasaranPemda(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error) {
	tahun = strings.TrimSpace(tahun)
	if len(tahun) != 4 {
		return sasaranpemda.LockDataPemdaResponse{}, fmt.Errorf("format tahun tidak valid")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	if err := s.LockDataPemdaRepository.Lock(ctx, tx, lockJenisSasaranPemda, tahun); err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	lock, err := s.LockDataPemdaRepository.FindByJenisTahun(ctx, tx, lockJenisSasaranPemda, tahun)
	if err != nil {
		return sasaranpemda.LockDataPemdaResponse{Jenis: lockJenisSasaranPemda, Tahun: tahun, Locked: true}, nil
	}
	return toLockSasaranResponse(lock, true), nil
}
func (s *SasaranPemdaServiceImpl) UnlockSasaranPemda(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error) {
	tahun = strings.TrimSpace(tahun)
	if len(tahun) != 4 {
		return sasaranpemda.LockDataPemdaResponse{}, fmt.Errorf("format tahun tidak valid")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	if err := s.LockDataPemdaRepository.Unlock(ctx, tx, lockJenisSasaranPemda, tahun); err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	return sasaranpemda.LockDataPemdaResponse{Jenis: lockJenisSasaranPemda, Tahun: tahun, Locked: false}, nil
}
func (s *SasaranPemdaServiceImpl) IsSasaranPemdaLocked(ctx context.Context, tahun string) (sasaranpemda.LockDataPemdaResponse, error) {
	tahun = strings.TrimSpace(tahun)
	if len(tahun) != 4 {
		return sasaranpemda.LockDataPemdaResponse{}, fmt.Errorf("format tahun tidak valid")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	locked, err := s.LockDataPemdaRepository.IsLocked(ctx, tx, lockJenisSasaranPemda, tahun)
	if err != nil {
		return sasaranpemda.LockDataPemdaResponse{}, err
	}
	resp := sasaranpemda.LockDataPemdaResponse{Jenis: lockJenisSasaranPemda, Tahun: tahun, Locked: locked}
	if locked {
		if lock, err := s.LockDataPemdaRepository.FindByJenisTahun(ctx, tx, lockJenisSasaranPemda, tahun); err == nil {
			resp.Id = lock.Id
		}
	}
	return resp, nil
}
func (s *SasaranPemdaServiceImpl) FindAllLockSasaranPemda(ctx context.Context) ([]sasaranpemda.LockDataPemdaResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)
	locks, err := s.LockDataPemdaRepository.FindAllByJenis(ctx, tx, lockJenisSasaranPemda)
	if err != nil {
		return nil, err
	}
	result := make([]sasaranpemda.LockDataPemdaResponse, 0, len(locks))
	for _, l := range locks {
		result = append(result, toLockSasaranResponse(l, true))
	}
	return result, nil
}
