package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type TargetProgramKegiatanServiceImpl struct {
	repo     repository.TargetProgramKegiatanRepository
	db       *sql.DB
	validate *validator.Validate
}

func NewTargetProgramKegiatanServiceImpl(
	repo repository.TargetProgramKegiatanRepository,
	db *sql.DB,
	validate *validator.Validate,
) *TargetProgramKegiatanServiceImpl {
	return &TargetProgramKegiatanServiceImpl{repo: repo, db: db, validate: validate}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func toTargetPKResponse(t domain.TargetProgramKegiatan) programkegiatan.TargetProgramKegiatanResponse {
	return programkegiatan.TargetProgramKegiatanResponse{
		Id:            t.Id,
		KodeIndikator: t.KodeIndikator,
		KodeOpd:       t.KodeOpd,
		Jenis:         t.Jenis,
		Target:        t.Target,
		Satuan:        t.Satuan,
		Tahun:         t.Tahun,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

// Create membuat target baru. Parameter jenis diisi oleh controller sesuai endpoint
// yang dipanggil (renstra / renja / rankhir / penetapan).
func (s *TargetProgramKegiatanServiceImpl) Create(ctx context.Context, jenis string, req programkegiatan.TargetProgramKegiatanCreateRequest) (programkegiatan.TargetProgramKegiatanResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	result, err := s.repo.Create(ctx, tx, domain.TargetProgramKegiatan{
		KodeIndikator: req.KodeIndikator,
		KodeOpd:       req.KodeOpd,
		Jenis:         jenis,
		Target:        req.Target,
		Satuan:        req.Satuan,
		Tahun:         req.Tahun,
	})
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}
	return toTargetPKResponse(result), nil
}

// ── Update ───────────────────────────────────────────────────────────────────

// Update hanya mengizinkan perubahan pada kolom target dan satuan.
// Kolom jenis dan tahun tidak dapat diubah.
func (s *TargetProgramKegiatanServiceImpl) Update(ctx context.Context, id int, req programkegiatan.TargetProgramKegiatanUpdateRequest) (programkegiatan.TargetProgramKegiatanResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	existing, err := s.repo.FindById(ctx, tx, id)
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}

	// Hanya target dan satuan yang diperbarui; jenis dan tahun tetap
	existing.Target = req.Target
	existing.Satuan = req.Satuan

	result, err := s.repo.Update(ctx, tx, existing)
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}
	return toTargetPKResponse(result), nil
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *TargetProgramKegiatanServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	if _, err := s.repo.FindById(ctx, tx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, tx, id)
}

// ── FindById ─────────────────────────────────────────────────────────────────

func (s *TargetProgramKegiatanServiceImpl) FindById(ctx context.Context, id int) (programkegiatan.TargetProgramKegiatanResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	result, err := s.repo.FindById(ctx, tx, id)
	if err != nil {
		return programkegiatan.TargetProgramKegiatanResponse{}, err
	}
	return toTargetPKResponse(result), nil
}

// ── FindByKodeIndikatorAndKodeOpd ────────────────────────────────────────────

func (s *TargetProgramKegiatanServiceImpl) FindByKodeIndikatorAndKodeOpd(ctx context.Context, kodeIndikator, kodeOpd string) ([]programkegiatan.TargetProgramKegiatanResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer helper.CommitOrRollback(tx)

	list, err := s.repo.FindByKodeIndikatorAndKodeOpd(ctx, tx, kodeIndikator, kodeOpd)
	if err != nil {
		return nil, err
	}
	result := make([]programkegiatan.TargetProgramKegiatanResponse, 0, len(list))
	for _, t := range list {
		result = append(result, toTargetPKResponse(t))
	}
	return result, nil
}
