package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/lockrenaksiopd"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrLockRenaksiOpdInvalidParameter = errors.New("parameter lock renaksi opd tidak valid")
	ErrLockRenaksiOpdNotFound         = errors.New("lock renaksi opd tidak ditemukan")
)

type LockRenaksiOpdServiceImpl struct {
	LockRenaksiOpdRepository repository.LockRenaksiOpdRepository
	RencanaAksiOpdRepository repository.RencanaAksiOpdRepository
	DB                       *sql.DB
	validator                *validator.Validate
}

func NewLockRenaksiOpdServiceImpl(
	lockRenaksiOpdRepository repository.LockRenaksiOpdRepository,
	rencanaAksiOpdRepository repository.RencanaAksiOpdRepository,
	db *sql.DB,
	validator *validator.Validate,
) *LockRenaksiOpdServiceImpl {
	return &LockRenaksiOpdServiceImpl{
		LockRenaksiOpdRepository: lockRenaksiOpdRepository,
		RencanaAksiOpdRepository: rencanaAksiOpdRepository,
		DB:                       db,
		validator:                validator,
	}
}

func (service *LockRenaksiOpdServiceImpl) Lock(ctx context.Context, kodeOpd, tahun string, request lockrenaksiopd.LockRenaksiOpdRequest) (lockrenaksiopd.LockRenaksiOpdResponse, error) {
	kodeOpd, tahun, err := validateLockRenaksiOpdParams(kodeOpd, tahun)
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	if err := service.validator.Struct(request); err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	lock, err := service.LockRenaksiOpdRepository.Lock(ctx, tx, domain.LockRenaksiOpd{
		KodeOpd:      kodeOpd,
		Tahun:        tahun,
		SasaranId:    request.SasaranId,
		RekinId:      request.RekinId,
		AksiKegiatan: request.AksiKegiatan,
		SubKegiatan:  request.SubKegiatan,
		Anggaran:     request.Anggaran,
		NamaPemilik:  request.NamaPemilik,
		Tw1:          request.Tw1,
		Tw2:          request.Tw2,
		Tw3:          request.Tw3,
		Tw4:          request.Tw4,
	})
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	return toLockRenaksiOpdResponse(lock), nil
}

func (service *LockRenaksiOpdServiceImpl) Unlock(ctx context.Context, kodeOpd, tahun string, id int) error {
	kodeOpd, tahun, err := validateLockRenaksiOpdParams(kodeOpd, tahun)
	if err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("%w: id_renaksiopd tidak valid", ErrLockRenaksiOpdInvalidParameter)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	lockKodeOpd, lockTahun, sasaranId, rekinId, err := service.RencanaAksiOpdRepository.FindLockContextById(ctx, tx, id)
	if err != nil {
		return err
	}
	if lockKodeOpd != kodeOpd || lockTahun != tahun {
		return fmt.Errorf("%w: id_renaksiopd %d tidak sesuai dengan kode_opd/tahun", ErrLockRenaksiOpdInvalidParameter, id)
	}

	if err := service.LockRenaksiOpdRepository.Unlock(ctx, tx, kodeOpd, tahun, sasaranId, rekinId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: id_renaksiopd %d", ErrLockRenaksiOpdNotFound, id)
		}
		return err
	}
	return nil
}

func (service *LockRenaksiOpdServiceImpl) FindById(ctx context.Context, kodeOpd, tahun string, id int) (lockrenaksiopd.LockRenaksiOpdResponse, error) {
	kodeOpd, tahun, err := validateLockRenaksiOpdParams(kodeOpd, tahun)
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	if id <= 0 {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, fmt.Errorf("%w: id_renaksiopd tidak valid", ErrLockRenaksiOpdInvalidParameter)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	lockKodeOpd, lockTahun, sasaranId, rekinId, err := service.RencanaAksiOpdRepository.FindLockContextById(ctx, tx, id)
	if err != nil {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	if lockKodeOpd != kodeOpd || lockTahun != tahun {
		return lockrenaksiopd.LockRenaksiOpdResponse{}, fmt.Errorf("%w: id_renaksiopd %d tidak sesuai dengan kode_opd/tahun", ErrLockRenaksiOpdInvalidParameter, id)
	}

	lock, err := service.LockRenaksiOpdRepository.FindByContext(ctx, tx, kodeOpd, tahun, sasaranId, rekinId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lockrenaksiopd.LockRenaksiOpdResponse{}, fmt.Errorf("%w: id_renaksiopd %d", ErrLockRenaksiOpdNotFound, id)
		}
		return lockrenaksiopd.LockRenaksiOpdResponse{}, err
	}
	return toLockRenaksiOpdResponse(lock), nil
}

func (service *LockRenaksiOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) ([]lockrenaksiopd.LockRenaksiOpdResponse, error) {
	kodeOpd, tahun, err := validateLockRenaksiOpdParams(kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	locks, err := service.LockRenaksiOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	result := make([]lockrenaksiopd.LockRenaksiOpdResponse, 0, len(locks))
	for _, lock := range locks {
		result = append(result, toLockRenaksiOpdResponse(lock))
	}
	return result, nil
}

func validateLockRenaksiOpdParams(kodeOpd, tahun string) (string, string, error) {
	kodeOpd = strings.TrimSpace(kodeOpd)
	tahun = strings.TrimSpace(tahun)
	if kodeOpd == "" {
		return "", "", fmt.Errorf("%w: kode_opd wajib diisi", ErrLockRenaksiOpdInvalidParameter)
	}
	if len(tahun) != 4 {
		return "", "", fmt.Errorf("%w: format tahun tidak valid", ErrLockRenaksiOpdInvalidParameter)
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return "", "", fmt.Errorf("%w: tahun harus berupa angka", ErrLockRenaksiOpdInvalidParameter)
	}
	return kodeOpd, tahun, nil
}

func toLockRenaksiOpdResponse(lock domain.LockRenaksiOpd) lockrenaksiopd.LockRenaksiOpdResponse {
	return lockrenaksiopd.LockRenaksiOpdResponse{
		Id:           lock.Id,
		KodeOpd:      lock.KodeOpd,
		Tahun:        lock.Tahun,
		SasaranId:    lock.SasaranId,
		RekinId:      lock.RekinId,
		AksiKegiatan: lock.AksiKegiatan,
		SubKegiatan:  lock.SubKegiatan,
		Anggaran:     lock.Anggaran,
		NamaPemilik:  lock.NamaPemilik,
		Tw1:          lock.Tw1,
		Tw2:          lock.Tw2,
		Tw3:          lock.Tw3,
		Tw4:          lock.Tw4,
		Locked:       true,
	}
}
