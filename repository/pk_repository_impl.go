package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type PkRepositoryImpl struct{}

func NewPkRepositoryImpl() *PkRepositoryImpl {
	return &PkRepositoryImpl{}
}

func (repository *PkRepositoryImpl) FindByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[int][]domain.PkOpd, error) {
	query := `
    SELECT pk.id,
           pk.kode_opd,
           pk.nama_opd,
           pk.level_pk,
           pk.nama_atasan,
           pk.nip_atasan,
           pk.id_rekin_atasan,
           pk.rekin_atasan,
           pk.nip_pemilik_pk,
           pk.nama_pemilik_pk,
           pk.id_rekin_pemilik_pk,
           pk.rekin_pemilik_pk,
           pk.tahun,
           pk.keterangan
    FROM pk_opd pk
    WHERE pk.kode_opd = ? AND pk.tahun = ?
    ORDER BY pk.level_pk
    `

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// group by level
	results := make(map[int][]domain.PkOpd)

	for rows.Next() {
		var pkOpd domain.PkOpd

		err := rows.Scan(
			&pkOpd.Id,
			&pkOpd.KodeOpd,
			&pkOpd.NamaOpd,
			&pkOpd.LevelPk,
			&pkOpd.NamaAtasan,
			&pkOpd.NipAtasan,
			&pkOpd.IdRekinAtasan,
			&pkOpd.RekinAtasan,
			&pkOpd.NipPemilikPk,
			&pkOpd.NamaPemilikPk,
			&pkOpd.IdRekinPemilikPk,
			&pkOpd.RekinPemilikPk,
			&pkOpd.Tahun,
			&pkOpd.Keterangan,
		)
		if err != nil {
			return nil, fmt.Errorf("scan pk_opd failed: %w", err)
		}
		results[pkOpd.LevelPk] = append(results[pkOpd.LevelPk], pkOpd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repository *PkRepositoryImpl) HubungkanRekin(
	ctx context.Context,
	tx *sql.Tx,
	pk domain.PkOpd,
) error {

	query := `
	INSERT INTO pk_opd (
		id,
		kode_opd,
		nama_opd,
		level_pk,
		nip_atasan,
		nama_atasan,
		id_rekin_atasan,
		rekin_atasan,
		nip_pemilik_pk,
		nama_pemilik_pk,
		id_rekin_pemilik_pk,
		rekin_pemilik_pk,
		tahun,
		keterangan
	) VALUES (
		UUID(),
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
	ON DUPLICATE KEY UPDATE
		level_pk           = VALUES(level_pk),
		nip_atasan         = VALUES(nip_atasan),
		nama_atasan        = VALUES(nama_atasan),
		id_rekin_atasan    = VALUES(id_rekin_atasan),
		rekin_atasan       = VALUES(rekin_atasan),
		keterangan         = VALUES(keterangan),
		updated_at         = CURRENT_TIMESTAMP
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		pk.KodeOpd,
		pk.NamaOpd,
		pk.LevelPk,
		pk.NipAtasan,
		pk.NamaAtasan,
		pk.IdRekinAtasan,
		pk.RekinAtasan,
		pk.NipPemilikPk,
		pk.NamaPemilikPk,
		pk.IdRekinPemilikPk,
		pk.RekinPemilikPk,
		pk.Tahun,
		pk.Keterangan,
	)

	if err != nil {
		return fmt.Errorf("hubungkan rekin failed: %w", err)
	}

	return nil
}
