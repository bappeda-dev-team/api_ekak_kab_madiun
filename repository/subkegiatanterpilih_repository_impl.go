package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type SubKegiatanTerpilihRepositoryImpl struct {
}

func NewSubKegiatanTerpilihRepositoryImpl() *SubKegiatanTerpilihRepositoryImpl {
	return &SubKegiatanTerpilihRepositoryImpl{}
}

func (repository *SubKegiatanTerpilihRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error) {
	script := "UPDATE tb_rencana_kinerja SET kode_subkegiatan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, subKegiatanTerpilih.KodeSubKegiatan, subKegiatanTerpilih.Id)
	if err != nil {
		return subKegiatanTerpilih, err
	}

	return subKegiatanTerpilih, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) error {
	scriptDelete := "UPDATE tb_rencana_kinerja SET kode_subkegiatan = '' WHERE id = ? AND kode_subkegiatan = ?"
	_, err := tx.ExecContext(ctx, scriptDelete, id, kodeSubKegiatan)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindByIdAndKodeSubKegiatan(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) (domain.SubKegiatanTerpilih, error) {
	script := "SELECT id, kode_subkegiatan FROM tb_rencana_kinerja WHERE id = ? AND kode_subkegiatan = ?"
	var subKegiatanTerpilih domain.SubKegiatanTerpilih
	err := tx.QueryRowContext(ctx, script, id, kodeSubKegiatan).Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.KodeSubKegiatan)
	return subKegiatanTerpilih, err
}

func (repository *SubKegiatanTerpilihRepositoryImpl) CreateRekin(ctx context.Context, tx *sql.Tx, idSubKegiatan string, rekinId string, kodeSubKegiatan string) error {
	// Validasi keberadaan subkegiatan di tb_subkegiatan
	checkSubkegiatanScript := "SELECT COUNT(*) FROM tb_subkegiatan_opd WHERE kode_subkegiatan = ?"
	var subkegiatanCount int
	err := tx.QueryRowContext(ctx, checkSubkegiatanScript, kodeSubKegiatan).Scan(&subkegiatanCount)
	if err != nil {
		return fmt.Errorf("error saat memeriksa data subkegiatan: %v", err)
	}
	if subkegiatanCount == 0 {
		return fmt.Errorf("subkegiatan dengan kode %s belum dipilih opd", kodeSubKegiatan)
	}

	// Hapus data subkegiatan terpilih yang lama untuk rekin_id yang sama
	deleteScript := "DELETE FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	_, err = tx.ExecContext(ctx, deleteScript, rekinId)
	if err != nil {
		return fmt.Errorf("error saat menghapus data subkegiatan terpilih yang lama: %v", err)
	}

	// Generate UUID baru untuk primary key
	newId := uuid.New().String()

	// Insert data baru ke tb_subkegiatan_terpilih
	script := "INSERT INTO tb_subkegiatan_terpilih (id, subkegiatan_id, rekin_id, kode_subkegiatan) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, newId, idSubKegiatan, rekinId, kodeSubKegiatan)
	if err != nil {
		return fmt.Errorf("error saat menyimpan subkegiatan terpilih: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("gagal menyimpan subkegiatan terpilih")
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) DeleteSubKegiatanTerpilih(ctx context.Context, tx *sql.Tx, idSubKegiatan string) error {
	script := "DELETE FROM tb_subkegiatan_terpilih WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("error saat menghapus subkegiatan terpilih: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan", idSubKegiatan)
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.SubKegiatanTerpilih, error) {
	script := "SELECT id, subkegiatan_id, rekin_id, kode_subkegiatan FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	rows, err := tx.QueryContext(ctx, script, rekinId)
	if err != nil {
		return nil, fmt.Errorf("error saat mengambil data subkegiatan terpilih: %v", err)
	}
	defer rows.Close()

	var result []domain.SubKegiatanTerpilih
	for rows.Next() {
		var subKegiatanTerpilih domain.SubKegiatanTerpilih
		err := rows.Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.SubkegiatanId, &subKegiatanTerpilih.RekinId, &subKegiatanTerpilih.KodeSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("error saat scanning data subkegiatan terpilih: %v", err)
		}
		result = append(result, subKegiatanTerpilih)
	}

	return result, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) CreateOPD(ctx context.Context, tx *sql.Tx, subkegiatanOpd domain.SubKegiatanOpd) (domain.SubKegiatanOpd, error) {
	// Cek apakah kombinasi kode_subkegiatan dan kode_opd sudah ada
	checkScript := "SELECT COUNT(*) FROM tb_subkegiatan_opd WHERE kode_subkegiatan = ? AND kode_opd = ? AND tahun = ?"
	var count int
	err := tx.QueryRowContext(ctx, checkScript, subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun).Scan(&count)
	if err != nil {
		return domain.SubKegiatanOpd{}, fmt.Errorf("error saat memeriksa duplikasi: %v", err)
	}

	// Jika sudah ada kombinasi yang sama, kembalikan error
	if count > 0 {
		return domain.SubKegiatanOpd{}, fmt.Errorf("subkegiatan dengan kode %s sudah ada di OPD %s untuk tahun %s",
			subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun)
	}

	// Jika belum ada, lanjutkan dengan insert
	script := "INSERT INTO tb_subkegiatan_opd (id, kode_subkegiatan, kode_opd, tahun) VALUES (?,?,?,?)"
	result, err := tx.ExecContext(ctx, script, subkegiatanOpd.Id, subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun)
	if err != nil {
		return domain.SubKegiatanOpd{}, fmt.Errorf("error saat memilih subkegiatan opd: %v", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return subkegiatanOpd, err
	}
	subkegiatanOpd.Id = int(lastInsertId)

	return subkegiatanOpd, nil
}
func (repository *SubKegiatanTerpilihRepositoryImpl) UpdateOPD(ctx context.Context, tx *sql.Tx, subkegiatanOpd domain.SubKegiatanOpd) (domain.SubKegiatanOpd, error) {
	checkScript := "SELECT COUNT(*) FROM tb_subkegiatan_opd WHERE kode_subkegiatan = ? AND kode_opd = ? AND tahun = ?"
	var count int
	err := tx.QueryRowContext(ctx, checkScript, subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun).Scan(&count)
	if err != nil {
		return domain.SubKegiatanOpd{}, fmt.Errorf("error saat memeriksa duplikasi: %v", err)
	}

	// Jika sudah ada kombinasi yang sama, kembalikan error
	if count > 0 {
		return domain.SubKegiatanOpd{}, fmt.Errorf("subkegiatan dengan kode %s sudah ada di OPD %s untuk tahun %s",
			subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun)
	}

	script := "UPDATE tb_subkegiatan_opd SET kode_subkegiatan = ?, kode_opd = ?, tahun = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, script, subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun, subkegiatanOpd.Id)
	if err != nil {
		return domain.SubKegiatanOpd{}, fmt.Errorf("error saat mengupdate subkegiatan opd: %v", err)
	}
	return subkegiatanOpd, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindallOpd(ctx context.Context, tx *sql.Tx, kodeOpd, tahun *string) ([]domain.SubKegiatanOpd, error) {
	script := "SELECT id, kode_subkegiatan, kode_opd, tahun FROM tb_subkegiatan_opd WHERE 1=1"
	var params []interface{}

	if kodeOpd != nil {
		script += " AND kode_opd = ?"
		params = append(params, *kodeOpd)
	}

	if tahun != nil {
		script += " AND tahun = ?"
		params = append(params, *tahun)
	}
	script += " order by kode_subkegiatan asc"

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, fmt.Errorf("error saat mencari subkegiatan opd: %v", err)
	}

	defer rows.Close()

	var subkegiatanOPD []domain.SubKegiatanOpd
	for rows.Next() {
		var sub domain.SubKegiatanOpd
		err := rows.Scan(&sub.Id, &sub.KodeSubKegiatan, &sub.KodeOpd, &sub.Tahun)
		if err != nil {
			return nil, fmt.Errorf("error saat  mencari subkegiatan opd: %v", err)
		}
		subkegiatanOPD = append(subkegiatanOPD, sub)
	}

	return subkegiatanOPD, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.SubKegiatanOpd, error) {
	script := "SELECT id, kode_subkegiatan, kode_opd, tahun FROM tb_subkegiatan_opd WHERE id = ?"
	row := tx.QueryRowContext(ctx, script, id)

	var sub domain.SubKegiatanOpd
	err := row.Scan(&sub.Id, &sub.KodeSubKegiatan, &sub.KodeOpd, &sub.Tahun)
	if err != nil {
		return domain.SubKegiatanOpd{}, fmt.Errorf("error saat mencari usulan inovasi: %v", err)
	}
	return sub, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) DeleteSubOpd(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_subkegiatan_opd WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return fmt.Errorf("error saat menghapus subkegiatan opd: %v", err)
	}
	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindAllSubkegiatanByBidangUrusanOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.SubKegiatan, error) {
	// Ekstrak kode bidang urusan dari kode OPD
	var bidangUrusanCodes []string

	// Split kode OPD untuk mendapatkan bidang urusan
	// Format: 5.01.5.05.0.00.01.0000
	parts := strings.Split(kodeOpd, ".")
	if len(parts) >= 4 {
		// Ambil bidang urusan pertama (5.01)
		bidangUrusanCodes = append(bidangUrusanCodes, parts[0]+"."+parts[1])
		// Ambil bidang urusan kedua (5.05)
		bidangUrusanCodes = append(bidangUrusanCodes, parts[2]+"."+parts[3])
		// Ambil bidang urusan ketiga (0.00)
		bidangUrusanCodes = append(bidangUrusanCodes, parts[4]+"."+parts[5])
	}

	// Tambahkan kode X.XX yang harus ada di semua OPD
	bidangUrusanCodes = append(bidangUrusanCodes, "X.XX")

	// Buat query dengan UNION untuk menggabungkan hasil dari setiap bidang urusan
	var queries []string
	var params []interface{}

	for _, bidangUrusan := range bidangUrusanCodes {
		query := `
            SELECT DISTINCT s.kode_subkegiatan, s.nama_subkegiatan
            FROM tb_subkegiatan s
            WHERE s.kode_subkegiatan LIKE ?
        `
		// Gunakan pattern matching untuk mencari subkegiatan dengan kode bidang urusan yang sesuai
		// Contoh: untuk bidang urusan 5.01, cari subkegiatan yang dimulai dengan 5.01
		params = append(params, bidangUrusan+"%")
		queries = append(queries, query)
	}

	// Gabungkan semua query dengan UNION
	finalQuery := strings.Join(queries, " UNION ")
	finalQuery += " ORDER BY kode_subkegiatan ASC"

	// Eksekusi query
	rows, err := tx.QueryContext(ctx, finalQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("error saat mengambil data subkegiatan: %v", err)
	}
	defer rows.Close()

	var result []domain.SubKegiatan
	for rows.Next() {
		var subkegiatan domain.SubKegiatan
		err := rows.Scan(
			&subkegiatan.KodeSubKegiatan,
			&subkegiatan.NamaSubKegiatan,
		)
		if err != nil {
			return nil, fmt.Errorf("error saat scanning data subkegiatan: %v", err)
		}
		result = append(result, subkegiatan)
	}

	return result, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error) {
	script := "SELECT id, kode_subkegiatan, nama_subkegiatan FROM tb_subkegiatan_opd WHERE kode_subkegiatan = ?"
	row := tx.QueryRowContext(ctx, script, kodeSubKegiatan)

	var sub domain.SubKegiatan
	err := row.Scan(&sub.Id, &sub.KodeSubKegiatan, &sub.NamaSubKegiatan)
	if err != nil {
		return domain.SubKegiatan{}, fmt.Errorf("error saat mencari subkegiatan: %v", err)
	}
	return sub, nil
}
