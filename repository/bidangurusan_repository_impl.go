package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"strings"
)

type BidangUrusanRepositoryImpl struct {
}

func NewBidangUrusanRepositoryImpl() *BidangUrusanRepositoryImpl {
	return &BidangUrusanRepositoryImpl{}
}

func (repository *BidangUrusanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan {
	script := "INSERT INTO tb_bidang_urusan (id, kode_bidang_urusan, nama_bidang_urusan, tahun) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, bidangurusan.Id, bidangurusan.KodeBidangUrusan, bidangurusan.NamaBidangUrusan, bidangurusan.Tahun)
	if err != nil {
		return bidangurusan
	}
	return bidangurusan
}

func (repository *BidangUrusanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan {
	script := "UPDATE tb_bidang_urusan SET kode_bidang_urusan = ?, nama_bidang_urusan = ?, tahun = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, bidangurusan.KodeBidangUrusan, bidangurusan.NamaBidangUrusan, bidangurusan.Tahun, bidangurusan.Id)
	if err != nil {
		return bidangurusan
	}
	return bidangurusan
}

func (repository *BidangUrusanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_bidang_urusan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *BidangUrusanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.BidangUrusan, error) {
	script := "SELECT id, kode_bidang_urusan, nama_bidang_urusan, tahun FROM tb_bidang_urusan WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domainmaster.BidangUrusan{}, err
	}
	defer rows.Close()

	bidangurusan := domainmaster.BidangUrusan{}
	if rows.Next() {
		rows.Scan(&bidangurusan.Id, &bidangurusan.KodeBidangUrusan, &bidangurusan.NamaBidangUrusan, &bidangurusan.Tahun)
	}
	return bidangurusan, nil
}

func (repository *BidangUrusanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.BidangUrusan, error) {
	script := "SELECT id, kode_bidang_urusan, nama_bidang_urusan, tahun FROM tb_bidang_urusan"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.BidangUrusan{}, err
	}
	defer rows.Close()

	var bidangurusans []domainmaster.BidangUrusan
	for rows.Next() {
		bidangurusan := domainmaster.BidangUrusan{}
		rows.Scan(&bidangurusan.Id, &bidangurusan.KodeBidangUrusan, &bidangurusan.NamaBidangUrusan, &bidangurusan.Tahun)
		bidangurusans = append(bidangurusans, bidangurusan)
	}
	return bidangurusans, nil
}

func (repository *BidangUrusanRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.BidangUrusan, error) {
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
		return []domainmaster.BidangUrusan{}, nil
	}

	// Membuat query dengan IN clause
	query := "SELECT id, kode_bidang_urusan, nama_bidang_urusan FROM tb_bidang_urusan WHERE kode_bidang_urusan IN ("
	params := make([]interface{}, len(kodeBidangUrusans))
	for i := range kodeBidangUrusans {
		if i > 0 {
			query += ","
		}
		query += "?"
		params[i] = kodeBidangUrusans[i]
	}
	query += ")"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bidangUrusans []domainmaster.BidangUrusan
	for rows.Next() {
		bidangUrusan := domainmaster.BidangUrusan{}
		err := rows.Scan(&bidangUrusan.Id, &bidangUrusan.KodeBidangUrusan, &bidangUrusan.NamaBidangUrusan)
		if err != nil {
			return nil, err
		}
		bidangUrusans = append(bidangUrusans, bidangUrusan)
	}

	return bidangUrusans, nil
}

// Tambahkan method baru
// func (repository *BidangUrusanRepositoryImpl) FindByKodeBidangUrusan(ctx context.Context, tx *sql.Tx, kodeBidangUrusan string) (domainmaster.BidangUrusan, error) {
// 	script := `
//         SELECT
//             bu.id,
//             bu.kode_bidang_urusan,
//             bu.nama_bidang_urusan,
//             u.nama_urusan
//         FROM
//             tb_bidang_urusan bu
//             INNER JOIN tb_urusan u ON LEFT(bu.kode_bidang_urusan, 1) = u.kode_urusan
//         WHERE
//             bu.kode_bidang_urusan = ?
//     `

// 	var bidangUrusan domainmaster.BidangUrusan
// 	err := tx.QueryRowContext(ctx, script, kodeBidangUrusan).Scan(
// 		&bidangUrusan.Id,
// 		&bidangUrusan.KodeBidangUrusan,
// 		&bidangUrusan.NamaBidangUrusan,
// 		&bidangUrusan.NamaUrusan,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return domainmaster.BidangUrusan{}, fmt.Errorf("bidang urusan dengan kode %s tidak ditemukan", kodeBidangUrusan)
// 		}
// 		return domainmaster.BidangUrusan{}, err
// 	}

// 	return bidangUrusan, nil
// }

func (repository *BidangUrusanRepositoryImpl) FindByKodeBidangUrusan(ctx context.Context, tx *sql.Tx, kodeBidangUrusan string) (domainmaster.BidangUrusan, error) {
	// Jika kodeBidangUrusan kosong, kembalikan objek default tanpa error
	if kodeBidangUrusan == "" {
		return domainmaster.BidangUrusan{
			Id:               "",
			KodeBidangUrusan: "",
			NamaBidangUrusan: "", // Memberikan label yang lebih bermakna
			NamaUrusan:       "", // Memberikan label yang lebih bermakna
		}, nil
	}

	script := `
        SELECT 
            COALESCE(bu.id, ''),
            COALESCE(bu.kode_bidang_urusan, ''),
            COALESCE(bu.nama_bidang_urusan, ''),
            COALESCE(u.nama_urusan, '')
        FROM 
            tb_bidang_urusan bu
            LEFT JOIN tb_urusan u ON LEFT(bu.kode_bidang_urusan, 1) = u.kode_urusan
        WHERE 
            bu.kode_bidang_urusan = ?
    `

	var bidangUrusan domainmaster.BidangUrusan
	err := tx.QueryRowContext(ctx, script, kodeBidangUrusan).Scan(
		&bidangUrusan.Id,
		&bidangUrusan.KodeBidangUrusan,
		&bidangUrusan.NamaBidangUrusan,
		&bidangUrusan.NamaUrusan,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domainmaster.BidangUrusan{
				Id:               "",
				KodeBidangUrusan: kodeBidangUrusan,
				NamaBidangUrusan: "",
				NamaUrusan:       "",
			}, nil
		}
		return domainmaster.BidangUrusan{}, err
	}

	return bidangUrusan, nil
}

func (repository *BidangUrusanRepositoryImpl) FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domainmaster.BidangUrusan, error) {
	if kodeSubKegiatan == "" {
		return domainmaster.BidangUrusan{}, fmt.Errorf("Kode Subkegiatan tidak ditemukan")
	}
	kodeBidangUrusan := kodeSubKegiatan[:4]
	script := "SELECT id, kode_bidang_urusan, nama_bidang_urusan FROM tb_bidang_urusan WHERE kode_bidang_urusan = ?"
	rows, err := tx.QueryContext(ctx, script, kodeBidangUrusan)
	if err != nil {
		return domainmaster.BidangUrusan{}, err
	}
	defer rows.Close()

	bidangurusan := domainmaster.BidangUrusan{}

	if rows.Next() {
		err := rows.Scan(&bidangurusan.Id, &bidangurusan.KodeBidangUrusan, &bidangurusan.NamaBidangUrusan)
		if err != nil {
			return domainmaster.BidangUrusan{}, err
		}
	} else {
		return domainmaster.BidangUrusan{}, fmt.Errorf("bidang urusan dengan kode %s tidak ditemukan", kodeBidangUrusan)
	}
	return bidangurusan, nil
}

func (repository *BidangUrusanRepositoryImpl) FindByKodeSubKegiatans(ctx context.Context, tx *sql.Tx, kodeSubKegiatans []string) (map[string]domainmaster.BidangUrusan, error) {

	resultMap := make(map[string]domainmaster.BidangUrusan)

	if len(kodeSubKegiatans) == 0 {
		return resultMap, nil
	}

	uniqBidangUrusan := make(map[string]struct{})
	kodeBidangUrusanList := make([]string, 0)

	subGroup := make(map[string][]string)

	for _, ks := range kodeSubKegiatans {
		if len(ks) < 7 {
			// Skip kode sub yang invalid
			continue
		}

		kodeBidangUrusan := ks[:4]

		subGroup[kodeBidangUrusan] = append(subGroup[kodeBidangUrusan], ks)

		if _, exists := uniqBidangUrusan[kodeBidangUrusan]; !exists {
			uniqBidangUrusan[kodeBidangUrusan] = struct{}{}
			kodeBidangUrusanList = append(kodeBidangUrusanList, kodeBidangUrusan)
		}
	}

	if len(kodeBidangUrusanList) == 0 {
		return resultMap, nil
	}

	placeholders := make([]string, len(kodeBidangUrusanList))
	args := make([]any, len(kodeBidangUrusanList))

	for i, kode := range kodeBidangUrusanList {
		placeholders[i] = "?"
		args[i] = kode
	}

	query := fmt.Sprintf(`
      SELECT kode_bidang_urusan, nama_bidang_urusan
      FROM tb_bidang_urusan
      WHERE kode_bidang_urusan IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bidangUrusanByKode := make(map[string]domainmaster.BidangUrusan)

	if rows.Next() {
		var bu domainmaster.BidangUrusan

		err := rows.Scan(&bu.KodeBidangUrusan, &bu.NamaBidangUrusan)
		if err != nil {
			return nil, err
		}
		bidangUrusanByKode[bu.KodeBidangUrusan] = bu
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Step 4. Mapping kembali ke kode_sub_kegiatan asli
	for kodeBidangUrusan, subs := range subGroup {
		base, ok := bidangUrusanByKode[kodeBidangUrusan]
		if !ok {
			continue // bidang urusan tidak ditemukan di DB
		}

		for _, kodeSub := range subs {
			resultMap[kodeSub] = domainmaster.BidangUrusan{
				KodeBidangUrusan: base.KodeBidangUrusan,
				NamaBidangUrusan: base.NamaBidangUrusan,
				KodeSubKegiatan:  kodeSub, // ★ simpan kode_sub asli
			}
		}
	}

	return resultMap, nil
}
