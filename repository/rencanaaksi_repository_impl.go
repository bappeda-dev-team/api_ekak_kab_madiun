package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
)

type RencanaAksiRepositoryImpl struct {
}

func NewRencanaAksiRepositoryImpl() *RencanaAksiRepositoryImpl {
	return &RencanaAksiRepositoryImpl{}
}

func (repository *RencanaAksiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error) {
	script := `
        SELECT urutan 
        FROM tb_rencana_aksi 
        WHERE rencana_kinerja_id = ? 
        ORDER BY urutan ASC
    `
	rows, err := tx.QueryContext(ctx, script, rencanaAksi.RencanaKinerjaId)
	if err != nil {
		return domain.RencanaAksi{}, fmt.Errorf("gagal mengambil data urutan: %v", err)
	}
	defer rows.Close()

	var urutanList []int
	for rows.Next() {
		var urutan int
		if err := rows.Scan(&urutan); err != nil {
			return domain.RencanaAksi{}, fmt.Errorf("gagal scan urutan: %v", err)
		}
		urutanList = append(urutanList, urutan)
	}

	// Cari urutan berurutan yang perlu diupdate
	var urutanToUpdate []int
	prevUrutan := -1
	for _, urutan := range urutanList {
		if urutan >= rencanaAksi.Urutan {
			if prevUrutan == -1 || urutan == prevUrutan+1 {
				urutanToUpdate = append(urutanToUpdate, urutan)
				prevUrutan = urutan
			} else {
				break // Keluar dari loop jika urutan tidak berurutan
			}
		}
	}

	// Update urutan dari yang terbesar
	for i := len(urutanToUpdate) - 1; i >= 0; i-- {
		scriptUpdate := `
            UPDATE tb_rencana_aksi 
            SET urutan = urutan + 1 
            WHERE rencana_kinerja_id = ? 
            AND urutan = ?
        `
		_, err := tx.ExecContext(ctx, scriptUpdate, rencanaAksi.RencanaKinerjaId, urutanToUpdate[i])
		if err != nil {
			return domain.RencanaAksi{}, fmt.Errorf("gagal mengupdate urutan: %v", err)
		}
	}

	// Insert data baru
	scriptInsert := "INSERT INTO tb_rencana_aksi (id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, scriptInsert, rencanaAksi.Id, rencanaAksi.RencanaKinerjaId, rencanaAksi.KodeOpd, rencanaAksi.Urutan, rencanaAksi.NamaRencanaAksi)
	if err != nil {
		return domain.RencanaAksi{}, err
	}

	return rencanaAksi, nil
}

func (repository *RencanaAksiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error) {
	// Dapatkan urutan lama
	var urutanLama int
	scriptGetOld := "SELECT urutan FROM tb_rencana_aksi WHERE id = ?"
	err := tx.QueryRowContext(ctx, scriptGetOld, rencanaAksi.Id).Scan(&urutanLama)
	if err != nil {
		return domain.RencanaAksi{}, fmt.Errorf("gagal mendapatkan urutan lama: %v", err)
	}

	if urutanLama != rencanaAksi.Urutan {
		// Dapatkan semua urutan yang ada
		script := `
            SELECT urutan 
            FROM tb_rencana_aksi 
            WHERE rencana_kinerja_id = ? 
            AND id != ?
            ORDER BY urutan ASC
        `
		rows, err := tx.QueryContext(ctx, script, rencanaAksi.RencanaKinerjaId, rencanaAksi.Id)
		if err != nil {
			return domain.RencanaAksi{}, fmt.Errorf("gagal mengambil data urutan: %v", err)
		}
		defer rows.Close()

		var urutanList []int
		for rows.Next() {
			var urutan int
			if err := rows.Scan(&urutan); err != nil {
				return domain.RencanaAksi{}, fmt.Errorf("gagal scan urutan: %v", err)
			}
			urutanList = append(urutanList, urutan)
		}

		// Cari urutan berurutan yang perlu diupdate
		var urutanToUpdate []int
		prevUrutan := -1
		for _, urutan := range urutanList {
			if urutan >= rencanaAksi.Urutan {
				if prevUrutan == -1 || urutan == prevUrutan+1 {
					urutanToUpdate = append(urutanToUpdate, urutan)
					prevUrutan = urutan
				} else {
					break // Keluar dari loop jika urutan tidak berurutan
				}
			}
		}

		// Update urutan dari yang terbesar
		for i := len(urutanToUpdate) - 1; i >= 0; i-- {
			scriptUpdate := `
                UPDATE tb_rencana_aksi 
                SET urutan = urutan + 1 
                WHERE rencana_kinerja_id = ? 
                AND urutan = ?
            `
			_, err := tx.ExecContext(ctx, scriptUpdate, rencanaAksi.RencanaKinerjaId, urutanToUpdate[i])
			if err != nil {
				return domain.RencanaAksi{}, fmt.Errorf("gagal mengupdate urutan: %v", err)
			}
		}
	}

	// Update data rencana aksi
	script := "UPDATE tb_rencana_aksi SET urutan = ?, nama_rencana_aksi = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, script, rencanaAksi.Urutan, rencanaAksi.NamaRencanaAksi, rencanaAksi.Id)
	if err != nil {
		return domain.RencanaAksi{}, err
	}

	return rencanaAksi, nil
}
func (repository *RencanaAksiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	// Hapus terlebih dahulu dari tb_pelaksanaan_rencana_aksi
	scriptPelaksanaan := "DELETE FROM tb_pelaksanaan_rencana_aksi WHERE rencana_aksi_id = ?"
	_, err := tx.ExecContext(ctx, scriptPelaksanaan, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data dari tb_pelaksanaan_rencana_aksi: %v", err)
	}

	// Kemudian hapus dari rencana_aksi
	scriptRencanaAksi := "DELETE FROM tb_rencana_aksi WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptRencanaAksi, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data dari rencana_aksi: %v", err)
	}

	return nil
}

func (repository *RencanaAksiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.RencanaAksi, error) {
	script := "SELECT id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi FROM tb_rencana_aksi WHERE id = ?"
	row := tx.QueryRowContext(ctx, script, id)
	var rencanaAksi domain.RencanaAksi
	err := row.Scan(&rencanaAksi.Id, &rencanaAksi.RencanaKinerjaId, &rencanaAksi.KodeOpd, &rencanaAksi.Urutan, &rencanaAksi.NamaRencanaAksi)
	if err != nil {
		return domain.RencanaAksi{}, err
	}
	return rencanaAksi, nil
}

func (repository *RencanaAksiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) ([]domain.RencanaAksi, error) {
	script := "SELECT id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi FROM tb_rencana_aksi"
	var args []any

	if rencanaKinerjaId != "" {
		script += " WHERE rencana_kinerja_id = ?"
		args = append(args, rencanaKinerjaId)
	}

	script += " ORDER BY urutan ASC"

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.RencanaAksi{}, err
	}
	defer rows.Close()

	var rencanaAksis []domain.RencanaAksi
	for rows.Next() {
		var rencanaAksi domain.RencanaAksi
		err := rows.Scan(&rencanaAksi.Id, &rencanaAksi.RencanaKinerjaId, &rencanaAksi.KodeOpd, &rencanaAksi.Urutan, &rencanaAksi.NamaRencanaAksi)
		if err != nil {
			return []domain.RencanaAksi{}, err
		}
		rencanaAksis = append(rencanaAksis, rencanaAksi)
	}

	return rencanaAksis, nil
}

func (repository *RencanaAksiRepositoryImpl) IsUrutanExistsForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_rencana_aksi WHERE rencana_kinerja_id = ? AND urutan = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId, urutan).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *RencanaAksiRepositoryImpl) IsUrutanExistsForRencanaKinerjaExcludingId(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int, excludeId string) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_rencana_aksi WHERE rencana_kinerja_id = ? AND urutan = ? AND id != ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId, urutan, excludeId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *RencanaAksiRepositoryImpl) GetTotalBobotForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) (int, error) {
	script := `
		SELECT COALESCE(SUM(pra.bobot), 0)
		FROM tb_rencana_aksi ra
		JOIN tb_pelaksanaan_rencana_aksi pra ON ra.id = pra.rencana_aksi_id
		WHERE ra.rencana_kinerja_id = ?
	`
	var totalBobot int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId).Scan(&totalBobot)
	return totalBobot, err
}

func (repository *RencanaAksiRepositoryImpl) FindAllBatch(ctx context.Context, tx *sql.Tx, rencanaKinerjaIds []string) ([]domain.RencanaAksi, error) {
	if len(rencanaKinerjaIds) == 0 {
		return []domain.RencanaAksi{}, errors.New("ids tidak boleh kosong")
	}

	baseQuery := `SELECT
                  id,
                  rencana_kinerja_id,
                  kode_opd, urutan,
                  nama_rencana_aksi
                 FROM tb_rencana_aksi
                 WHERE rencana_kinerja_id IN (?)
                 ORDER BY urutan ASC`

	query, args := helper.BuildInQueryString(baseQuery, rencanaKinerjaIds)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying rencana aksi batch: %v", err)
	}
	defer rows.Close()

	var rencanaAksis []domain.RencanaAksi
	for rows.Next() {
		var rencanaAksi domain.RencanaAksi
		err := rows.Scan(&rencanaAksi.Id, &rencanaAksi.RencanaKinerjaId, &rencanaAksi.KodeOpd, &rencanaAksi.Urutan, &rencanaAksi.NamaRencanaAksi)
		if err != nil {
			return []domain.RencanaAksi{}, err
		}
		rencanaAksis = append(rencanaAksis, rencanaAksi)
	}

	return rencanaAksis, nil
}
