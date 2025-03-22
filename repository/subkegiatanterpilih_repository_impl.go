package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"

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

func (repository *SubKegiatanTerpilihRepositoryImpl) CreateRekin(ctx context.Context, tx *sql.Tx, idSubKegiatan string, rekinId string) error {
	// Validasi keberadaan subkegiatan di tb_subkegiatan
	checkSubkegiatanScript := "SELECT COUNT(*) FROM tb_subkegiatan WHERE id = ?"
	var subkegiatanCount int
	err := tx.QueryRowContext(ctx, checkSubkegiatanScript, idSubKegiatan).Scan(&subkegiatanCount)
	if err != nil {
		return fmt.Errorf("error saat memeriksa data subkegiatan: %v", err)
	}
	if subkegiatanCount == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan di tb_subkegiatan", idSubKegiatan)
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
	script := "INSERT INTO tb_subkegiatan_terpilih (id, subkegiatan_id, rekin_id) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, newId, idSubKegiatan, rekinId)
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
	script := "SELECT id, subkegiatan_id, rekin_id FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	rows, err := tx.QueryContext(ctx, script, rekinId)
	if err != nil {
		return nil, fmt.Errorf("error saat mengambil data subkegiatan terpilih: %v", err)
	}
	defer rows.Close()

	var result []domain.SubKegiatanTerpilih
	for rows.Next() {
		var subKegiatanTerpilih domain.SubKegiatanTerpilih
		err := rows.Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.SubkegiatanId, &subKegiatanTerpilih.RekinId)
		if err != nil {
			return nil, fmt.Errorf("error saat scanning data subkegiatan terpilih: %v", err)
		}
		result = append(result, subKegiatanTerpilih)
	}

	return result, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) CreateOPD(ctx context.Context, tx *sql.Tx, subkegiatanOpd domain.SubKegiatanOpd) (domain.SubKegiatanOpd, error) {
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
	script := "UPDATE tb_subkegiatan_opd SET kode_subkegiatan = ?, kode_opd = ?, tahun = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, subkegiatanOpd.KodeSubKegiatan, subkegiatanOpd.KodeOpd, subkegiatanOpd.Tahun, subkegiatanOpd.Id)
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
