package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ekak_kabupaten_madiun/model/domain"
)

type CloneRecordRepositoryImpl struct {
}

func NewCloneRecordRepositoryImpl() *CloneRecordRepositoryImpl {
	return &CloneRecordRepositoryImpl{}
}

func (r *CloneRecordRepositoryImpl) Create(
	ctx context.Context,
	tx *sql.Tx,
	cloneRecord domain.CloneRecord,
) (domain.CloneRecord, error) {

	query := `
		INSERT INTO clone_record (
			kode_clone,
			kode_opd,
			tahun_asal,
			tahun_target,
			keterangan_tahun_clone,
			updated_by
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		cloneRecord.KodeClone,
		cloneRecord.KodeOpd,
		cloneRecord.TahunAsal,
		cloneRecord.TahunTarget,
		cloneRecord.KeteranganTahunClone,
		cloneRecord.UpdatedBy,
	)
	if err != nil {
		return cloneRecord, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return cloneRecord, err
	}

	cloneRecord.Id = int(lastInsertId)
	cloneRecord.CreatedAt = time.Now()
	cloneRecord.UpdatedAt = time.Now()

	return cloneRecord, nil
}

func (r *CloneRecordRepositoryImpl) GetCloneByKodeOpdTahunTujuan(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	tahunTarget string,
) (domain.CloneRecord, error) {

	query := `
		SELECT
			id,
			kode_clone,
			kode_opd,
			tahun_asal,
			tahun_target,
                        status,
			created_at,
			updated_at,
			updated_by
		FROM clone_record
		WHERE kode_opd = ?
		AND tahun_target = ?
                ORDER BY id DESC
		LIMIT 1
	`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahunTarget)
	if err != nil {
		return domain.CloneRecord{}, err
	}
	defer rows.Close()

	cloneRecord := domain.CloneRecord{}

	if rows.Next() {
		err := rows.Scan(
			&cloneRecord.Id,
			&cloneRecord.KodeClone,
			&cloneRecord.KodeOpd,
			&cloneRecord.TahunAsal,
			&cloneRecord.TahunTarget,
			&cloneRecord.Status,
			&cloneRecord.CreatedAt,
			&cloneRecord.UpdatedAt,
			&cloneRecord.UpdatedBy,
		)
		if err != nil {
			return cloneRecord, err
		}
		return cloneRecord, nil
	}

	return cloneRecord, errors.New("clone record not found")
}

func (r *CloneRecordRepositoryImpl) UpdateStatus(
	ctx context.Context,
	tx *sql.Tx,
	id int,
	status string,
	errMsg string,
) error {

	query := `
		UPDATE clone_record
		SET status = ?, error_message = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(ctx, query, status, errMsg, id)
	return err
}
