package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
)

type PelaksanaanRencanaAksiRepositoryImpl struct {
}

func NewPelaksanaanRencanaAksiRepositoryImpl() *PelaksanaanRencanaAksiRepositoryImpl {
	return &PelaksanaanRencanaAksiRepositoryImpl{}
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pelaksanaanRencanaAksi domain.PelaksanaanRencanaAksi) (domain.PelaksanaanRencanaAksi, error) {
	script := "INSERT INTO tb_pelaksanaan_rencana_aksi (id, rencana_aksi_id, bobot, bulan) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, pelaksanaanRencanaAksi.Id, pelaksanaanRencanaAksi.RencanaAksiId, pelaksanaanRencanaAksi.Bobot, pelaksanaanRencanaAksi.Bulan)
	if err != nil {
		return domain.PelaksanaanRencanaAksi{}, err
	}
	return pelaksanaanRencanaAksi, nil
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pelaksanaanRencanaAksi domain.PelaksanaanRencanaAksi) (domain.PelaksanaanRencanaAksi, error) {
	script := "UPDATE tb_pelaksanaan_rencana_aksi SET bobot = ?, bulan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pelaksanaanRencanaAksi.Bobot, pelaksanaanRencanaAksi.Bulan, pelaksanaanRencanaAksi.Id)
	if err != nil {
		return domain.PelaksanaanRencanaAksi{}, err
	}
	return pelaksanaanRencanaAksi, nil
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_pelaksanaan_rencana_aksi WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	return err
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.PelaksanaanRencanaAksi, error) {
	script := "SELECT id, rencana_aksi_id, bobot, bulan FROM tb_pelaksanaan_rencana_aksi WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PelaksanaanRencanaAksi{}, err
	}
	defer rows.Close()

	var pelaksanaanRencanaAksi domain.PelaksanaanRencanaAksi
	if rows.Next() {
		err := rows.Scan(&pelaksanaanRencanaAksi.Id, &pelaksanaanRencanaAksi.RencanaAksiId, &pelaksanaanRencanaAksi.Bobot, &pelaksanaanRencanaAksi.Bulan)
		if err != nil {
			return domain.PelaksanaanRencanaAksi{}, err
		}
		return pelaksanaanRencanaAksi, nil
	}
	return domain.PelaksanaanRencanaAksi{}, sql.ErrNoRows
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) FindByRencanaAksiId(ctx context.Context, tx *sql.Tx, rencanaAksiId string) ([]domain.PelaksanaanRencanaAksi, error) {
	script := "SELECT id, rencana_aksi_id, bobot, bulan FROM tb_pelaksanaan_rencana_aksi WHERE rencana_aksi_id = ?"
	rows, err := tx.QueryContext(ctx, script, rencanaAksiId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pelaksanaanRencanaAksiList []domain.PelaksanaanRencanaAksi
	for rows.Next() {
		var pelaksanaanRencanaAksi domain.PelaksanaanRencanaAksi
		err := rows.Scan(&pelaksanaanRencanaAksi.Id, &pelaksanaanRencanaAksi.RencanaAksiId, &pelaksanaanRencanaAksi.Bobot, &pelaksanaanRencanaAksi.Bulan)
		if err != nil {
			return nil, err
		}
		pelaksanaanRencanaAksiList = append(pelaksanaanRencanaAksiList, pelaksanaanRencanaAksi)
	}
	return pelaksanaanRencanaAksiList, nil
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) FindByRencanaAksiIdsBatch(ctx context.Context, tx *sql.Tx, rencanaAksiIds []string) ([]domain.PelaksanaanRencanaAksi, error) {
	if len(rencanaAksiIds) == 0 {
		return []domain.PelaksanaanRencanaAksi{}, errors.New("ids tidak boleh kosong")
	}

	baseQuery := `SELECT
            id,
            rencana_aksi_id,
            bobot,
            bulan
           FROM tb_pelaksanaan_rencana_aksi
           WHERE rencana_aksi_id IN (?)`
	query, args := helper.BuildInQueryString(baseQuery, rencanaAksiIds)

	rows, err := tx.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error querying rencana kinerja: %v", err)
	}
	defer rows.Close()

	var pelaksanaanRencanaAksiList []domain.PelaksanaanRencanaAksi
	for rows.Next() {
		var pelaksanaanRencanaAksi domain.PelaksanaanRencanaAksi
		err := rows.Scan(&pelaksanaanRencanaAksi.Id, &pelaksanaanRencanaAksi.RencanaAksiId, &pelaksanaanRencanaAksi.Bobot, &pelaksanaanRencanaAksi.Bulan)
		if err != nil {
			return nil, err
		}
		pelaksanaanRencanaAksiList = append(pelaksanaanRencanaAksiList, pelaksanaanRencanaAksi)
	}
	return pelaksanaanRencanaAksiList, nil
}

func (repository *PelaksanaanRencanaAksiRepositoryImpl) ExistsByRencanaAksiIdAndBulan(ctx context.Context, tx *sql.Tx, rencanaAksiId string, bulan int) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_pelaksanaan_rencana_aksi WHERE rencana_aksi_id = ? AND bulan = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rencanaAksiId, bulan).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
