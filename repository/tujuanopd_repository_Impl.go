package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type TujuanOpdRepositoryImpl struct {
}

func NewTujuanOpdRepositoryImpl() *TujuanOpdRepositoryImpl {
	return &TujuanOpdRepositoryImpl{}
}

func (repository *TujuanOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) (domain.TujuanOpd, error) {
	script := "INSERT INTO tb_tujuan_opd (kode_opd, tujuan, rumus_perhitungan, sumber_data, tahun_awal, tahun_akhir) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, tujuanOpd.KodeOpd, tujuanOpd.Tujuan, tujuanOpd.RumusPerhitungan, tujuanOpd.SumberData, tujuanOpd.TahunAwal, tujuanOpd.TahunAkhir)
	if err != nil {
		return domain.TujuanOpd{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.TujuanOpd{}, err
	}

	tujuanOpd.Id = int(id)

	for _, indikator := range tujuanOpd.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, tujuan_opd_id, indikator) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator, indikator.Id, id, indikator.Indikator)
		if err != nil {
			return domain.TujuanOpd{}, err
		}

		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domain.TujuanOpd{}, err
			}
		}
	}

	return tujuanOpd, nil
}

func (repository *TujuanOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) error {
	// Update tujuan OPD
	script := "UPDATE tb_tujuan_opd SET kode_opd = ?, tujuan = ?, rumus_perhitungan = ?, sumber_data = ?, tahun_awal = ?, tahun_akhir = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script,
		tujuanOpd.KodeOpd,
		tujuanOpd.Tujuan,
		tujuanOpd.RumusPerhitungan,
		tujuanOpd.SumberData,
		tujuanOpd.TahunAwal,
		tujuanOpd.TahunAkhir,
		tujuanOpd.Id)
	if err != nil {
		return err
	}

	// Hapus indikator dan target lama
	scriptDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.tujuan_opd_id = ?
    `
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, tujuanOpd.Id)
	if err != nil {
		return err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE tujuan_opd_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, tujuanOpd.Id)
	if err != nil {
		return err
	}

	// Insert indikator dan target baru
	for _, indikator := range tujuanOpd.Indikator {
		// Insert indikator
		scriptIndikator := "INSERT INTO tb_indikator (id, tujuan_opd_id, indikator) VALUES (?, ?, ?)"
		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			tujuanOpd.Id,
			indikator.Indikator)
		if err != nil {
			return err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (repository *TujuanOpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, tujuanOpdId int) error {
	scriptDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.tujuan_opd_id = ?
    `
	_, err := tx.ExecContext(ctx, scriptDeleteTarget, tujuanOpdId)
	if err != nil {
		return err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE tujuan_opd_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, tujuanOpdId)
	if err != nil {
		return err
	}

	scriptDeleteTujuanOpd := "DELETE FROM tb_tujuan_opd WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteTujuanOpd, tujuanOpdId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *TujuanOpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, tujuanOpdId int) (domain.TujuanOpd, error) {
	script := `SELECT id, kode_opd, tujuan, rumus_perhitungan, sumber_data, tahun_awal, tahun_akhir 
               FROM tb_tujuan_opd WHERE id = ?`

	var tujuanOpd domain.TujuanOpd
	err := tx.QueryRowContext(ctx, script, tujuanOpdId).Scan(
		&tujuanOpd.Id,
		&tujuanOpd.KodeOpd,
		&tujuanOpd.Tujuan,
		&tujuanOpd.RumusPerhitungan,
		&tujuanOpd.SumberData,
		&tujuanOpd.TahunAwal,
		&tujuanOpd.TahunAkhir,
	)

	if err == sql.ErrNoRows {
		return domain.TujuanOpd{}, fmt.Errorf("tujuan opd with id %d not found", tujuanOpdId)
	}

	if err != nil {
		return domain.TujuanOpd{}, err
	}

	return tujuanOpd, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error) {
	script := `SELECT id, indikator 
               FROM tb_indikator 
               WHERE tujuan_opd_id = ?`

	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.Indikator)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string, tahun string) ([]domain.Target, error) {
	script := `
        SELECT id, target, satuan, tahun
        FROM tb_target 
        WHERE indikator_id = ?
        AND tahun <= ?
        ORDER BY tahun ASC
    `

	rows, err := tx.QueryContext(ctx, script, indikatorId, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(
			&target.Id,
			&target.Target,
			&target.Satuan,
			&target.Tahun,
		)
		if err != nil {
			return nil, err
		}
		target.IndikatorId = indikatorId
		targets = append(targets, target)
	}

	return targets, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.TujuanOpd, error) {
	script := `
        SELECT 
            t.id, 
            t.kode_opd, 
            t.tujuan, 
            t.rumus_perhitungan, 
            t.sumber_data, 
            t.tahun_awal, 
            t.tahun_akhir,
            COALESCE(i.id, '') as indikator_id,
            COALESCE(i.indikator, '') as indikator,
            COALESCE(tg.id, '') as target_id,
            COALESCE(tg.target, '') as target,
            COALESCE(tg.satuan, '') as satuan,
            COALESCE(tg.tahun, '') as target_tahun
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator i ON t.id = i.tujuan_opd_id
        LEFT JOIN tb_target tg ON i.id = tg.indikator_id 
        WHERE t.kode_opd = ? 
        AND ? BETWEEN t.tahun_awal AND t.tahun_akhir
        AND (tg.tahun <= ? OR tg.tahun IS NULL)
        ORDER BY t.id ASC
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tujuanOpds []domain.TujuanOpd
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			tujuanId         int
			kodeOpd          string
			tujuan           string
			rumusPerhitungan string
			sumberData       string
			tahunAwal        string
			tahunAkhir       string
			indikatorId      string
			indikator        string
			targetId         string
			target           string
			satuan           string
			targetTahun      string
		)

		err := rows.Scan(
			&tujuanId,
			&kodeOpd,
			&tujuan,
			&rumusPerhitungan,
			&sumberData,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikator,
			&targetId,
			&target,
			&satuan,
			&targetTahun,
		)
		if err != nil {
			return nil, err
		}

		// Cek dan tambahkan TujuanOpd jika belum ada
		tujuanOpd, exists := tujuanOpdMap[tujuanId]
		if !exists {
			tujuanOpd = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpd,
				Tujuan:           tujuan,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
				TahunAwal:        tahunAwal,
				TahunAkhir:       tahunAkhir,
				Indikator:        make([]domain.Indikator, 0),
			}
			tujuanOpdMap[tujuanId] = tujuanOpd
			tujuanOpds = append(tujuanOpds, *tujuanOpd)
		}

		// Jika ada indikator
		if indikatorId != "" {
			ind, exists := indikatorMap[indikatorId]
			if !exists {
				ind = &domain.Indikator{
					Id:          indikatorId,
					TujuanOpdId: tujuanId,
					Indikator:   indikator,
					Target:      make([]domain.Target, 0),
				}
				indikatorMap[indikatorId] = ind
				tujuanOpd.Indikator = append(tujuanOpd.Indikator, *ind)
			}

			// Jika ada target
			if targetId != "" {
				newTarget := domain.Target{
					Id:          targetId,
					IndikatorId: indikatorId,
					Target:      target,
					Satuan:      satuan,
					Tahun:       targetTahun,
				}
				ind.Target = append(ind.Target, newTarget)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Jika tidak ada data, kembalikan slice kosong
	if len(tujuanOpds) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}

	return tujuanOpds, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.TujuanOpd, error) {
	script := `
        SELECT 
            id, 
            kode_opd, 
            tujuan, 
            tahun_awal, 
            tahun_akhir
        FROM tb_tujuan_opd
        WHERE kode_opd = ?
        AND ? BETWEEN tahun_awal AND tahun_akhir
        ORDER BY id ASC
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tujuanOpds []domain.TujuanOpd
	for rows.Next() {
		var tujuanOpd domain.TujuanOpd
		err := rows.Scan(
			&tujuanOpd.Id,
			&tujuanOpd.KodeOpd,
			&tujuanOpd.Tujuan,
			&tujuanOpd.TahunAwal,
			&tujuanOpd.TahunAkhir,
		)
		if err != nil {
			return nil, err
		}
		tujuanOpds = append(tujuanOpds, tujuanOpd)
	}

	return tujuanOpds, nil
}

// Tambahkan fungsi untuk mengambil indikator
func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanOpdId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error) {
	script := `
        SELECT id, indikator 
        FROM tb_indikator
        WHERE tujuan_opd_id = ?
        ORDER BY id ASC
    `

	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(
			&indikator.Id,
			&indikator.Indikator,
		)
		if err != nil {
			return nil, err
		}
		indikator.TujuanOpdId = tujuanOpdId
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}
