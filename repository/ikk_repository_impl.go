package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type IkkRepositoryImpl struct {
}

func NewIkkRepositoryImpl() *IkkRepositoryImpl {
	return &IkkRepositoryImpl{}
}

func (repository *IkkRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ikk domain.IndikatorIkk) (domain.IndikatorIkk, error) {
	script := "INSERT INTO tb_ikk (kode_bidang_urusan, jenis, nama_indikator, target, satuan, keterangan) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		ikk.KodeBidangUrusan,
		ikk.Jenis,
		ikk.NamaIndikator,
		ikk.Target,
		ikk.Satuan,
		ikk.Keterangan)
	if err != nil {
		return domain.IndikatorIkk{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.IndikatorIkk{}, err
	}
	ikk.ID = int(id)

	return ikk, nil
}

func (repository *IkkRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ikk domain.IndikatorIkk) (domain.IndikatorIkk, error) {
	script := "UPDATE tb_ikk SET kode_bidang_urusan = ?, jenis = ?, nama_indikator = ?, target = ?, satuan = ?, keterangan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, ikk.KodeBidangUrusan, ikk.Jenis, ikk.NamaIndikator, ikk.Target, ikk.Satuan, ikk.Keterangan, ikk.ID)
	if err != nil {
		return domain.IndikatorIkk{}, err
	}
	return ikk, nil
}

func (repository *IkkRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_ikk WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *IkkRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.IndikatorIkk, error) {
	script := "SELECT id, kode_bidang_urusan, jenis, nama_indikator, target, satuan, keterangan FROM tb_ikk WHERE id = ?"
	var ikk domain.IndikatorIkk
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&ikk.ID,
		&ikk.KodeBidangUrusan,
		&ikk.Jenis,
		&ikk.NamaIndikator,
		&ikk.Target,
		&ikk.Satuan,
		&ikk.Keterangan,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.IndikatorIkk{}, errors.New("ikk tidak ditemukan")
		}
		return domain.IndikatorIkk{}, err
	}
	return ikk, nil
}

func (repository *IkkRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, jenis string, kodeOpd string) ([]domain.IndikatorIkk, error) {
	// Memisahkan kode OPD untuk mendapatkan kode bidang urusan
	kodeBidangUrusans := make([]string, 0)

	// Format kode OPD: 1.01.2.22.0.00.01.0000
	// Kode bidang urusan terdiri dari 3 bagian: 1.01 | 2.22 | 0.00

	// Mengambil kode bidang urusan pertama (1.01)
	if len(kodeOpd) >= 4 {
		kode1 := kodeOpd[:4]
		if kode1 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode1)
		}
	}

	// Mengambil kode bidang urusan kedua (2.22)
	if len(kodeOpd) >= 9 {
		kode2 := kodeOpd[5:9]
		if kode2 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode2)
		}
	}

	// Mengambil kode bidang urusan ketiga (0.00)
	if len(kodeOpd) >= 14 {
		kode3 := kodeOpd[10:14]
		if kode3 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode3)
		}
	}

	// Jika tidak ada kode bidang urusan yang valid
	if len(kodeBidangUrusans) == 0 {
		return []domain.IndikatorIkk{}, nil
	}

	// Membuat query dengan IN clause
	query := `SELECT ikk.id, 
			  ikk.kode_bidang_urusan, 
			  COALESCE(od.nama_opd, '') as nama_opd,
			  ikk.jenis, 
			  ikk.nama_indikator, 
			  ikk.target, 
			  ikk.satuan, 
			  ikk.keterangan 
			  FROM tb_ikk ikk
			  LEFT JOIN tb_operasional_daerah od 
			  ON od.kode_opd = ?
			  WHERE ikk.kode_bidang_urusan IN (`

	// params := make([]interface{}, len(kodeBidangUrusans))
	params := make([]interface{}, 0)
	params = append(params, kodeOpd)
	for i := range kodeBidangUrusans {
		if i > 0 {
			query += ","
		}
		query += "?"
		// params[i] = kodeBidangUrusans[i]
		params = append(params, kodeBidangUrusans[i])
	}
	query += ")"

	query += " AND ikk.jenis = ?"
	params = append(params, jenis)
	
	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bidangUrusans []domain.IndikatorIkk
	for rows.Next() {
		bidangUrusan := domain.IndikatorIkk{}
		err := rows.Scan(&bidangUrusan.ID, &bidangUrusan.KodeBidangUrusan, &bidangUrusan.NamaOpd, &bidangUrusan.Jenis, &bidangUrusan.NamaIndikator, &bidangUrusan.Satuan, &bidangUrusan.Target, &bidangUrusan.Keterangan)
		if err != nil {
			return nil, err
		}
		bidangUrusans = append(bidangUrusans, bidangUrusan)
	}

	return bidangUrusans, nil
}


