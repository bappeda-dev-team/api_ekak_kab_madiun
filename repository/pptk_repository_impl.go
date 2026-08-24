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

func (repository *PptkRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx,kodeSubkegiatan string, kodeOpd string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            tp.id, 
			tp.nip, 
			peg.nama, 
			tp.kode_opd, 
			tp.tahun, 
			tp.kode_sub_kegiatan, 
			tp.nip_atasan, 
			tpa.nama, 
			tp.aktif_at, 
			tp.nonaktif_at
        FROM tb_pptk tp
		LEFT JOIN tb_pegawai peg ON tp.nip = peg.nip
		LEFT JOIN tb_pegawai tpa ON tp.nip_atasan = tpa.nip
        WHERE tp.kode_sub_kegiatan = ? AND tp.kode_opd = ? AND tp.tahun = ?`

	rows, err := tx.QueryContext(ctx, script, kodeSubkegiatan, kodeOpd, tahun)
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
			&pptk.NamaPegawai,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.NamaAtasan,
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
func (repository *PptkRepositoryImpl) FindAllByNip(ctx context.Context, tx *sql.Tx, kodeSubkegiatan string, pegawaiId string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            tp.id, 
			tp.nip, 
			peg.nama, 
			tp.kode_opd, 
			tp.tahun, 
			tp.kode_sub_kegiatan, 
			tp.nip_atasan, 
			tpa.nama, 
			tp.aktif_at, 
			tp.nonaktif_at
        FROM tb_pptk tp
		LEFT JOIN tb_pegawai peg ON tp.nip = peg.nip
		LEFT JOIN tb_pegawai tpa ON tp.nip_atasan = tpa.nip
        WHERE tp.kode_sub_kegiatan = ? AND tp.nip = ? AND tp.tahun = ?`

	rows, err := tx.QueryContext(ctx, script, kodeSubkegiatan, pegawaiId, tahun)
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
			&pptk.NamaPegawai,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.NamaAtasan,
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