package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"strings"
)

type UrusanRepositoryImpl struct {
}

func NewUrusanRepositoryImpl() *UrusanRepositoryImpl {
	return &UrusanRepositoryImpl{}
}

func (repository *UrusanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "INSERT INTO tb_urusan(id, kode_urusan, nama_urusan) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, urusan.Id, urusan.KodeUrusan, urusan.NamaUrusan)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "UPDATE tb_urusan SET kode_urusan = ?, nama_urusan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, urusan.KodeUrusan, urusan.NamaUrusan, urusan.Id)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM tb_urusan"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Urusan{}, err
	}

	defer rows.Close()

	var urusans []domainmaster.Urusan
	for rows.Next() {
		urusan := domainmaster.Urusan{}
		rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		urusans = append(urusans, urusan)
	}

	return urusans, nil
}

func (repository *UrusanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM tb_urusan WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domainmaster.Urusan{}, err
	}
	defer rows.Close()

	urusan := domainmaster.Urusan{}

	if rows.Next() {
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		if err != nil {
			return domainmaster.Urusan{}, err
		}
	} else {
		return domainmaster.Urusan{}, fmt.Errorf("urusan dengan id %s tidak ditemukan", id)
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_urusan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UrusanRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error) {
	kodeUrusans := make([]string, 0)

	if len(kodeOpd) >= 1 {
		kode1 := kodeOpd[0:1]
		if kode1 != "0" {
			kodeUrusans = append(kodeUrusans, kode1)
		}
	}

	if len(kodeOpd) >= 6 {
		kode2 := kodeOpd[5:6]
		if kode2 != "0" {
			kodeUrusans = append(kodeUrusans, kode2)
		}
	}

	if len(kodeOpd) >= 11 {
		kode3 := kodeOpd[10:11]
		if kode3 != "0" {
			kodeUrusans = append(kodeUrusans, kode3)
		}
	}

	if len(kodeUrusans) == 0 {
		return []domainmaster.Urusan{}, nil
	}

	query := "SELECT id, kode_urusan, nama_urusan FROM tb_urusan WHERE kode_urusan IN ("
	params := make([]interface{}, len(kodeUrusans))
	for i := range kodeUrusans {
		if i > 0 {
			query += ","
		}
		query += "?"
		params[i] = kodeUrusans[i]
	}
	query += ")"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urusans []domainmaster.Urusan
	for rows.Next() {
		urusan := domainmaster.Urusan{}
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan)
		if err != nil {
			return nil, err
		}
		urusans = append(urusans, urusan)
	}

	return urusans, nil
}

func (repository *UrusanRepositoryImpl) FindUrusanAndBidangByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error) {
	parts := strings.Split(kodeOpd, ".")
	var bidangUrusans []string
	var currentUrusans []string

	fmt.Printf("Processing kode OPD: %s\n", kodeOpd)
	fmt.Printf("Parts: %v\n", parts)

	// Ekstrak urusan dan bidang urusan (maksimal 3)
	bidangCount := 0
	for i := 0; i < len(parts)-1 && bidangCount < 3; i += 2 {
		urusan := parts[i]
		if i+1 < len(parts) {
			bidang := parts[i+1]

			// Jika urusan dan bidang valid
			if urusan != "0" && urusan != "00" && bidang != "00" {
				// Tambahkan urusan ke daftar urusan jika belum ada
				if !contains(currentUrusans, urusan) {
					currentUrusans = append(currentUrusans, urusan)
				}

				// Buat kode bidang urusan
				bidangUrusan := urusan + "." + bidang
				if !contains(bidangUrusans, bidangUrusan) {
					bidangUrusans = append(bidangUrusans, bidangUrusan)
					fmt.Printf("Added bidang urusan: %s\n", bidangUrusan)
					bidangCount++
				}
			}
		}
	}

	fmt.Printf("Extracted urusan: %v\n", currentUrusans)
	fmt.Printf("Extracted bidang urusan (max 3): %v\n", bidangUrusans)

	if len(bidangUrusans) == 0 {
		return []domainmaster.Urusan{}, nil
	}

	// Buat placeholders untuk IN clause
	placeholders := make([]string, len(bidangUrusans))
	for i := range bidangUrusans {
		placeholders[i] = "?"
	}

	// Query untuk MySQL dengan LIMIT 3
	query := fmt.Sprintf(`
        SELECT DISTINCT
            u.id,
            u.kode_urusan,
            u.nama_urusan,
            bu.kode_bidang_urusan,
            bu.nama_bidang_urusan
        FROM 
            tb_urusan u
            INNER JOIN tb_bidang_urusan bu ON u.kode_urusan = LEFT(bu.kode_bidang_urusan, 1)
        WHERE 
            bu.kode_bidang_urusan IN (%s)
        ORDER BY 
            FIELD(bu.kode_bidang_urusan, %s)
        LIMIT 3
    `, strings.Join(placeholders, ","), strings.Join(placeholders, ","))

	// Siapkan arguments untuk query (duplikat karena digunakan di dua tempat)
	args := make([]interface{}, len(bidangUrusans)*2)
	for i, v := range bidangUrusans {
		args[i] = v
		args[i+len(bidangUrusans)] = v
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var result []domainmaster.Urusan
	urusanMap := make(map[string]*domainmaster.Urusan)

	for rows.Next() {
		var (
			id, kodeUrusan, namaUrusan         string
			kodeBidangUrusan, namaBidangUrusan string
		)

		err := rows.Scan(&id, &kodeUrusan, &namaUrusan, &kodeBidangUrusan, &namaBidangUrusan)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		fmt.Printf("Found row: urusan=%s, bidang=%s\n", kodeUrusan, kodeBidangUrusan)

		// Cek apakah urusan sudah ada di map
		if existingUrusan, exists := urusanMap[kodeUrusan]; exists {
			// Tambahkan bidang urusan ke urusan yang sudah ada jika belum mencapai 3
			if len(existingUrusan.BidangUrusan) < 3 {
				existingUrusan.BidangUrusan = append(existingUrusan.BidangUrusan, domainmaster.BidangUrusan{
					KodeBidangUrusan: kodeBidangUrusan,
					NamaBidangUrusan: namaBidangUrusan,
					Tahun:            "",
				})
				fmt.Printf("Added bidang %s to existing urusan %s\n", kodeBidangUrusan, kodeUrusan)
			}
		} else {
			// Buat urusan baru
			newUrusan := &domainmaster.Urusan{
				Id:         id,
				KodeUrusan: kodeUrusan,
				NamaUrusan: namaUrusan,
				BidangUrusan: []domainmaster.BidangUrusan{
					{
						KodeBidangUrusan: kodeBidangUrusan,
						NamaBidangUrusan: namaBidangUrusan,
						Tahun:            "",
					},
				},
			}
			urusanMap[kodeUrusan] = newUrusan
			fmt.Printf("Created new urusan %s with bidang %s\n", kodeUrusan, kodeBidangUrusan)
		}
	}

	// Konversi map ke slice result
	for _, urusan := range urusanMap {
		result = append(result, *urusan)
	}

	return result, nil
}

// Helper function untuk mengecek apakah slice contains value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (repository *UrusanRepositoryImpl) FindByKodeUrusan(ctx context.Context, tx *sql.Tx, kodeUrusan string) (domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan FROM tb_urusan WHERE kode_urusan = ?"
	rows, err := tx.QueryContext(ctx, script, kodeUrusan)
	if err != nil {
		return domainmaster.Urusan{}, err
	}
	defer rows.Close()

	urusan := domainmaster.Urusan{}

	if rows.Next() {
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		if err != nil {
			return domainmaster.Urusan{}, err
		}
	} else {
		return domainmaster.Urusan{}, fmt.Errorf("urusan dengan id %s tidak ditemukan", kodeUrusan)
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domainmaster.Urusan, error) {
	if kodeSubKegiatan == "" {
		return domainmaster.Urusan{}, fmt.Errorf("Kode Subkegiatan tidak ditemukan")
	}
	kodeUrusan := kodeSubKegiatan[:1]
	script := "SELECT id, kode_urusan, nama_urusan FROM tb_urusan WHERE kode_urusan = ?"
	rows, err := tx.QueryContext(ctx, script, kodeUrusan)
	if err != nil {
		return domainmaster.Urusan{}, err
	}
	defer rows.Close()

	urusan := domainmaster.Urusan{}

	if rows.Next() {
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan)
		if err != nil {
			return domainmaster.Urusan{}, err
		}
	} else {
		return domainmaster.Urusan{}, fmt.Errorf("urusan dengan kode %s tidak ditemukan", kodeUrusan)
	}

	return urusan, nil
}
