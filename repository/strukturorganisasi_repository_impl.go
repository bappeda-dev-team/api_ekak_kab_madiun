package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type StrukturOrganisasiRepositoryImpl struct {
}

func NewStrukturOrganisasiRepositoryImpl() *StrukturOrganisasiRepositoryImpl {
	return &StrukturOrganisasiRepositoryImpl{}
}

func (repository *StrukturOrganisasiRepositoryImpl) Create(
	ctx context.Context,
	tx *sql.Tx,
	so domain.StrukturOrganisasi,
) error {

	query := `
	INSERT INTO struktur_organisasi (
		nip_bawahan,
		nip_atasan,
		kode_opd,
		tahun
	) VALUES (
		?, ?, ?, ?
	)
	ON DUPLICATE KEY UPDATE
		nip_atasan = VALUES(nip_atasan),
		updated_at = CURRENT_TIMESTAMP
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		so.NipBawahan,
		so.NipAtasan,
		so.KodeOpd,
		so.Tahun,
	)
	if err != nil {
		return fmt.Errorf("insert struktur_organisasi failed: %w", err)
	}

	return nil
}

func (repository *StrukturOrganisasiRepositoryImpl) AtasanBawahanByKodeOpdTahun(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	tahun int,
) (map[string]string, error) {

	query := `
	SELECT
		nip_bawahan,
		nip_atasan
	FROM struktur_organisasi
	WHERE kode_opd = ?
	  AND tahun = ?
	`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("query struktur_organisasi failed: %w", err)
	}
	defer rows.Close()

	results := make(map[string]string)

	for rows.Next() {
		var bawahan, atasan string
		if err := rows.Scan(&bawahan, &atasan); err != nil {
			return nil, err
		}
		results[bawahan] = atasan
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
