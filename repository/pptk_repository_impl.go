package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type PptkRepositoryImpl struct {
}

func NewPptkRepositoryImpl() *PptkRepositoryImpl {
	return &PptkRepositoryImpl{}
}

func (repository *PptkRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error) {
	script := `INSERT INTO tb_pptk 
	(nip, kode_opd, tahun, kode_sub_kegiatan, nip_atasan, nonaktif_at) 
	VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pptk.Nip,
		pptk.KodeOpd,
		pptk.Tahun,
		pptk.KodeSubKegiatan,
		pptk.NipAtasan,
		pptk.NonAktifAt)
	if err != nil {
		return domain.Pptk{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Pptk{}, err
	}
	pptk.Id = int(id)

	return pptk, nil
}

func (repository *PptkRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error) {
	script := `UPDATE tb_pptk SET 
			   nip = ?, 
			   kode_opd = ?, 
			   tahun = ?, 
			   kode_sub_kegiatan = ?, 
			   nip_atasan = ?,
			   nonaktif_at = ?
			   WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, 
		pptk.Nip, 
		pptk.KodeOpd, 
		pptk.Tahun, 
		pptk.KodeSubKegiatan, 
		pptk.NipAtasan, 
		pptk.NonAktifAt, 
		pptk.Id)
	if err != nil {
		return domain.Pptk{}, err
	}
	return pptk, nil
}

func (repository *PptkRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_pptk WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PptkRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pptk, error) {
	script := `SELECT id, 
					  nip, 
					  kode_opd, 
					  tahun, 
					  kode_sub_kegiatan, 
					  nip_atasan, 
					  aktif_at, 
					  nonaktif_at 
					  FROM tb_pptk WHERE id = ?`
	var Pptk domain.Pptk
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&Pptk.Id,
		&Pptk.Nip,
		&Pptk.KodeOpd,
		&Pptk.Tahun,
		&Pptk.KodeSubKegiatan,
		&Pptk.NipAtasan,
		&Pptk.AktifAt,
		&Pptk.NonAktifAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Pptk{}, errors.New("pptk tidak ditemukan")
		}
		return domain.Pptk{}, err
	}
	return Pptk, nil
}

func (repository *PptkRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            id, 
			nip, 
			kode_opd, 
			tahun, 
			kode_sub_kegiatan, 
			nip_atasan, 
			aktif_at, 
			nonaktif_at
        FROM tb_pptk
        WHERE kode_opd = ? AND tahun = ?`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return []domain.Pptk{}, err
	}
	defer rows.Close()
	// log.Printf("TahunAWal: %s TAhun Akhir: %s", tahunAwal, tahunAkhir)

	var pptkList []domain.Pptk
	for rows.Next() {
		var pptk domain.Pptk
		err = rows.Scan(
			&pptk.Id,
			&pptk.Nip,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.AktifAt,
			&pptk.NonAktifAt,
		)
		if err != nil {
			return []domain.Pptk{}, err
		}
		pptkList = append(pptkList, pptk)
	}
	return pptkList, nil
}
func (repository *PptkRepositoryImpl) FindAllByNip(ctx context.Context, tx *sql.Tx, nip string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            id, 
			nip, 
			kode_opd, 
			tahun, 
			kode_sub_kegiatan, 
			nip_atasan, 
			aktif_at, 
			nonaktif_at
        FROM tb_pptk
        WHERE nip = ? AND tahun = ?`

	rows, err := tx.QueryContext(ctx, script, nip, tahun)
	if err != nil {
		return []domain.Pptk{}, err
	}
	defer rows.Close()
	// log.Printf("TahunAWal: %s TAhun Akhir: %s", tahunAwal, tahunAkhir)

	var pptkList []domain.Pptk
	for rows.Next() {
		var pptk domain.Pptk
		err = rows.Scan(
			&pptk.Id,
			&pptk.Nip,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.AktifAt,
			&pptk.NonAktifAt,
		)
		if err != nil {
			return []domain.Pptk{}, err
		}
		pptkList = append(pptkList, pptk)
	}
	return pptkList, nil
}