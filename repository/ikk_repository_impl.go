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

func (repository *IkkRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error) {

	script := `
		INSERT INTO tb_ikk 
		(kode_bidang_urusan, kode_opd, jenis, tahun, keterangan) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		script,
		ikk.KodeBidangUrusan,
		ikk.KodeOpd,
		ikk.Jenis,
		ikk.Tahun,
		ikk.Keterangan,
	)
	if err != nil {
		return domain.Ikk{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Ikk{}, err
	}

	ikk.ID = int(id)

	// insert indikator
	for i, indikator := range ikk.Indikators {

		scriptIndikator := `
			INSERT INTO tb_indikator_ikk
			(id_ikk, kode_opd, kode_bidang_urusan, indikator, tahun)
			VALUES (?, ?, ?, ?, ?)
		`

		resultIndikator, err := tx.ExecContext(
			ctx,
			scriptIndikator,
			ikk.ID,
			ikk.KodeOpd,
			ikk.KodeBidangUrusan,
			indikator.Indikator,
			ikk.Tahun,
		)
		if err != nil {
			return ikk, err
		}

		indikatorID, err := resultIndikator.LastInsertId()
		if err != nil {
			return ikk, err
		}

		// set ID indikator ke struct
		ikk.Indikators[i].ID = int(indikatorID)

		// insert targets
		for j, target := range indikator.Targets {

			scriptTarget := `
				INSERT INTO tb_target_ikk
				(id_indikator, target, satuan, tahun)
				VALUES (?, ?, ?, ?)
			`

			resultTarget, err := tx.ExecContext(
				ctx,
				scriptTarget,
				indikatorID,
				target.Target,
				target.Satuan,
				ikk.Tahun,
			)
			if err != nil {
				return ikk, err
			}

			targetID, err := resultTarget.LastInsertId()
			if err != nil {
				return ikk, err
			}

			// set ID target ke struct
			ikk.Indikators[i].Targets[j].ID = int(targetID)
		}
	}

	return ikk, nil
}

func (repository *IkkRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, ikk domain.Ikk) (domain.Ikk, error) {

	// ================= UPDATE IKK =================
	query := `
		UPDATE tb_ikk
		SET
			kode_bidang_urusan = ?,
			kode_opd = ?,
			jenis = ?,
			tahun = ?,
			keterangan = ?
		WHERE id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		ikk.KodeBidangUrusan,
		ikk.KodeOpd,
		ikk.Jenis,
		ikk.Tahun,
		ikk.Keterangan,
		ikk.ID,
	)
	if err != nil {
		return domain.Ikk{}, err
	}

	// ================= HAPUS TARGET LAMA =================
	deleteTarget := `
		DELETE ti
		FROM tb_target_ikk ti
		INNER JOIN tb_indikator_ikk ii
			ON ii.id = ti.id_indikator
		WHERE ii.id_ikk = ?
	`

	_, err = tx.ExecContext(ctx, deleteTarget, ikk.ID)
	if err != nil {
		return domain.Ikk{}, err
	}

	// ================= HAPUS INDIKATOR LAMA =================
	deleteIndikator := `
		DELETE FROM tb_indikator_ikk
		WHERE id_ikk = ?
	`

	_, err = tx.ExecContext(ctx, deleteIndikator, ikk.ID)
	if err != nil {
		return domain.Ikk{}, err
	}

	// ================= INSERT ULANG INDIKATOR =================
	for i, indikator := range ikk.Indikators {

		queryIndikator := `
			INSERT INTO tb_indikator_ikk
			(id_ikk, kode_opd, kode_bidang_urusan, indikator, tahun)
			VALUES (?, ?, ?, ?, ?)
		`

		resultIndikator, err := tx.ExecContext(
			ctx,
			queryIndikator,
			ikk.ID,
			ikk.KodeOpd,
			ikk.KodeBidangUrusan,
			indikator.Indikator,
			ikk.Tahun,
		)
		if err != nil {
			return domain.Ikk{}, err
		}

		indikatorID, err := resultIndikator.LastInsertId()
		if err != nil {
			return domain.Ikk{}, err
		}

		ikk.Indikators[i].ID = int(indikatorID)

		// ================= INSERT TARGET =================
		for j, target := range indikator.Targets {

			queryTarget := `
				INSERT INTO tb_target_ikk
				(id_indikator, target, satuan, tahun)
				VALUES (?, ?, ?, ?)
			`

			resultTarget, err := tx.ExecContext(
				ctx,
				queryTarget,
				indikatorID,
				target.Target,
				target.Satuan,
				ikk.Tahun,
			)
			if err != nil {
				return domain.Ikk{}, err
			}

			targetID, err := resultTarget.LastInsertId()
			if err != nil {
				return domain.Ikk{}, err
			}

			ikk.Indikators[i].Targets[j].ID = int(targetID)
		}
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

func (repository *IkkRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Ikk, error) {

	// ================= IKK =================
	query := `
		SELECT
			id,
			kode_bidang_urusan,
			kode_opd,
			jenis,
			tahun,
			keterangan
		FROM tb_ikk
		WHERE id = ?
	`

	var result domain.Ikk

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.KodeBidangUrusan,
		&result.KodeOpd,
		&result.Jenis,
		&result.Tahun,
		&result.Keterangan,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Ikk{}, errors.New("ikk tidak ditemukan")
		}
		return domain.Ikk{}, err
	}

	result.Indikators = make([]domain.IndikatorIkk, 0)

	// ================= INDIKATOR =================
	queryIndikator := `
		SELECT
			id,
			indikator
		FROM tb_indikator_ikk
		WHERE id_ikk = ?
	`

	rowsIndikator, err := tx.QueryContext(
		ctx,
		queryIndikator,
		result.ID,
	)
	if err != nil {
		return domain.Ikk{}, err
	}
	defer rowsIndikator.Close()

	indikatorIDs := make([]int, 0)

	for rowsIndikator.Next() {

		var indikator domain.IndikatorIkk

		err := rowsIndikator.Scan(
			&indikator.ID,
			&indikator.Indikator,
		)
		if err != nil {
			return domain.Ikk{}, err
		}

		indikator.Targets = make([]domain.TargetIkk, 0)

		result.Indikators = append(result.Indikators, indikator)
		indikatorIDs = append(indikatorIDs, indikator.ID)
	}

	// ================= TARGET =================
	if len(indikatorIDs) > 0 {

		placeholders := ""

		for i := 0; i < len(indikatorIDs); i++ {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
		}

		queryTarget := `
			SELECT
				id,
				id_indikator,
				target,
				satuan
			FROM tb_target_ikk
			WHERE id_indikator IN (` + placeholders + `)
		`

		args := make([]interface{}, len(indikatorIDs))
		for i, v := range indikatorIDs {
			args[i] = v
		}

		rowsTarget, err := tx.QueryContext(
			ctx,
			queryTarget,
			args...,
		)
		if err != nil {
			return domain.Ikk{}, err
		}
		defer rowsTarget.Close()

		targetMap := make(map[int][]domain.TargetIkk)

		for rowsTarget.Next() {

			var target domain.TargetIkk
			var idIndikator int

			err := rowsTarget.Scan(
				&target.ID,
				&idIndikator,
				&target.Target,
				&target.Satuan,
			)
			if err != nil {
				return domain.Ikk{}, err
			}

			targetMap[idIndikator] = append(
				targetMap[idIndikator],
				target,
			)
		}

		// attach target ke indikator
		for i, indikator := range result.Indikators {
			result.Indikators[i].Targets = targetMap[indikator.ID]
		}
	}

	return result, nil
}

func (repository *IkkRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, jenis string, kodeOpd string) ([]domain.Ikk, error) {
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
		return []domain.Ikk{}, nil
	}

	// Membuat query dengan IN clause
	query := `SELECT ikk.id, 
			  ikk.kode_bidang_urusan, 
			  COALESCE(bu.nama_bidang_urusan, '') as nama_bidang_urusan,
			  ikk.kode_opd, 
			  COALESCE(od.nama_opd, '') as nama_opd,
			  ikk.jenis, 
			  ikk.nama_indikator, 
			  ikk.target, 
			  ikk.satuan, 
			  ikk.tahun, 
			  ikk.keterangan 
			  FROM tb_ikk ikk
			  LEFT JOIN tb_operasional_daerah od 
			  ON od.kode_opd = ?
			  LEFT JOIN tb_bidang_urusan bu
    		  ON bu.kode_bidang_urusan = ikk.kode_bidang_urusan
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

	var bidangUrusans []domain.Ikk
	for rows.Next() {
		bidangUrusan := domain.Ikk{}
		err := rows.Scan(&bidangUrusan.ID, &bidangUrusan.KodeBidangUrusan, &bidangUrusan.NamaBidangUrusan, &bidangUrusan.KodeOpd, &bidangUrusan.NamaOpd, &bidangUrusan.Jenis, &bidangUrusan.NamaIndikator, &bidangUrusan.Target, &bidangUrusan.Satuan, &bidangUrusan.Tahun, &bidangUrusan.Keterangan)
		if err != nil {
			return nil, err
		}
		bidangUrusans = append(bidangUrusans, bidangUrusan)
	}

	return bidangUrusans, nil
}


func (repository *IkkRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.Ikk, error) {

	query := `
		SELECT id, kode_opd, kode_bidang_urusan, jenis, tahun, keterangan
		FROM tb_ikk
	`

	args := make([]interface{}, 0)

	if kodeOpd != "" {
		query += " WHERE kode_opd = ?"
		args = append(args, kodeOpd)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ikkMap := make(map[int]*domain.Ikk)
	ikkIDs := make([]int, 0)

	for rows.Next() {
		var item domain.Ikk

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.KodeBidangUrusan,
			&item.Jenis,
			&item.Tahun,
			&item.Keterangan,
		)
		if err != nil {
			return nil, err
		}

		item.Indikators = make([]domain.IndikatorIkk, 0)

		copyItem := item
		ikkMap[item.ID] = &copyItem
		ikkIDs = append(ikkIDs, item.ID)
	}

	if len(ikkIDs) == 0 {
		return []domain.Ikk{}, nil
	}

	// ================= INDICATOR =================
	placeholders := makePlaceholders(len(ikkIDs))

	queryInd := `
		SELECT id, id_ikk, indikator
		FROM tb_indikator_ikk
		WHERE id_ikk IN (` + placeholders + `)
	`

	argsInd := make([]interface{}, len(ikkIDs))
	for i, v := range ikkIDs {
		argsInd[i] = v
	}

	rowsInd, err := tx.QueryContext(ctx, queryInd, argsInd...)
	if err != nil {
		return nil, err
	}
	defer rowsInd.Close()

	indikatorIDs := make([]int, 0)

	for rowsInd.Next() {
		var ind domain.IndikatorIkk
		var idIkk int

		err := rowsInd.Scan(
			&ind.ID,
			&idIkk,
			&ind.Indikator,
		)
		if err != nil {
			return nil, err
		}

		ind.Targets = make([]domain.TargetIkk, 0)

		ikkMap[idIkk].Indikators = append(ikkMap[idIkk].Indikators, ind)
		indikatorIDs = append(indikatorIDs, ind.ID)
	}

	// ================= TARGET =================
	if len(indikatorIDs) > 0 {

		placeholders = makePlaceholders(len(indikatorIDs))

		queryTarget := `
			SELECT id, id_indikator, target, satuan
			FROM tb_target_ikk
			WHERE id_indikator IN (` + placeholders + `)
		`

		argsTarget := make([]interface{}, len(indikatorIDs))
		for i, v := range indikatorIDs {
			argsTarget[i] = v
		}

		rowsT, err := tx.QueryContext(ctx, queryTarget, argsTarget...)
		if err != nil {
			return nil, err
		}
		defer rowsT.Close()

		targetMap := make(map[int][]domain.TargetIkk)

		for rowsT.Next() {
			var t domain.TargetIkk
			var idInd int

			err := rowsT.Scan(
				&t.ID,
				&idInd,
				&t.Target,
				&t.Satuan,
			)
			if err != nil {
				return nil, err
			}

			targetMap[idInd] = append(targetMap[idInd], t)
		}

		// attach target ke indikator
		for _, ikk := range ikkMap {
			for i, ind := range ikk.Indikators {
				ikk.Indikators[i].Targets = targetMap[ind.ID]
			}
		}
	}

	result := make([]domain.Ikk, 0, len(ikkMap))
	for _, v := range ikkMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *IkkRepositoryImpl) FindAllByJenisAndKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, jenis string) ([]domain.Ikk, error) {

	query := `
		SELECT id, kode_opd, kode_bidang_urusan, jenis, tahun, keterangan
		FROM tb_ikk
		WHERE 1=1
	`

	args := make([]interface{}, 0)

	if kodeOpd != "" {
		query += " AND kode_opd = ?"
		args = append(args, kodeOpd)
	}

	if jenis != "" {
		query += " AND jenis = ?"
		args = append(args, jenis)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ikkMap := make(map[int]*domain.Ikk)
	ikkIDs := make([]int, 0)

	for rows.Next() {
		var item domain.Ikk

		err := rows.Scan(
			&item.ID,
			&item.KodeOpd,
			&item.KodeBidangUrusan,
			&item.Jenis,
			&item.Tahun,
			&item.Keterangan,
		)
		if err != nil {
			return nil, err
		}

		item.Indikators = make([]domain.IndikatorIkk, 0)

		copyItem := item
		ikkMap[item.ID] = &copyItem
		ikkIDs = append(ikkIDs, item.ID)
	}

	if len(ikkIDs) == 0 {
		return []domain.Ikk{}, nil
	}

	// ================= INDICATOR =================
	placeholders := makePlaceholders(len(ikkIDs))

	queryInd := `
		SELECT id, id_ikk, indikator
		FROM tb_indikator_ikk
		WHERE id_ikk IN (` + placeholders + `)
	`

	argsInd := make([]interface{}, len(ikkIDs))
	for i, v := range ikkIDs {
		argsInd[i] = v
	}

	rowsInd, err := tx.QueryContext(ctx, queryInd, argsInd...)
	if err != nil {
		return nil, err
	}
	defer rowsInd.Close()

	indikatorIDs := make([]int, 0)

	for rowsInd.Next() {
		var ind domain.IndikatorIkk
		var idIkk int

		err := rowsInd.Scan(
			&ind.ID,
			&idIkk,
			&ind.Indikator,
		)
		if err != nil {
			return nil, err
		}

		ind.Targets = make([]domain.TargetIkk, 0)

		ikkMap[idIkk].Indikators = append(ikkMap[idIkk].Indikators, ind)
		indikatorIDs = append(indikatorIDs, ind.ID)
	}

	// ================= TARGET =================
	if len(indikatorIDs) > 0 {

		placeholders = makePlaceholders(len(indikatorIDs))

		queryTarget := `
			SELECT id, id_indikator, target, satuan
			FROM tb_target_ikk
			WHERE id_indikator IN (` + placeholders + `)
		`

		argsTarget := make([]interface{}, len(indikatorIDs))
		for i, v := range indikatorIDs {
			argsTarget[i] = v
		}

		rowsT, err := tx.QueryContext(ctx, queryTarget, argsTarget...)
		if err != nil {
			return nil, err
		}
		defer rowsT.Close()

		targetMap := make(map[int][]domain.TargetIkk)

		for rowsT.Next() {
			var t domain.TargetIkk
			var idInd int

			err := rowsT.Scan(
				&t.ID,
				&idInd,
				&t.Target,
				&t.Satuan,
			)
			if err != nil {
				return nil, err
			}

			targetMap[idInd] = append(targetMap[idInd], t)
		}

		for _, ikk := range ikkMap {
			for i, ind := range ikk.Indikators {
				ikk.Indikators[i].Targets = targetMap[ind.ID]
			}
		}
	}

	result := make([]domain.Ikk, 0, len(ikkMap))
	for _, v := range ikkMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *IkkRepositoryImpl) FindSelection(ctx context.Context, tx *sql.Tx) ([]domain.BidangUrusanSelection, error) {

	query := `SELECT
				bu.kode_bidang_urusan,
				COALESCE(bu.nama_bidang_urusan, '') AS nama_bidang_urusan,
				od.kode_opd,
				COALESCE(od.nama_opd, '') AS nama_opd
			FROM tb_bidang_urusan bu
			CROSS JOIN tb_operasional_daerah od
			ORDER BY bu.kode_bidang_urusan ASC, od.kode_opd ASC`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selections := make([]domain.BidangUrusanSelection, 0)

	for rows.Next() {
		var selection domain.BidangUrusanSelection

		err := rows.Scan(
			&selection.KodeBidangUrusan,
			&selection.NamaBidangUrusan,
			&selection.KodeOpd,
			&selection.NamaOpd,
		)

		if err != nil {
			return nil, err
		}

		selections = append(selections, selection)
	}

	return selections, nil
}

func (repository *IkkRepositoryImpl) FindSelectionByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domain.BidangUrusanSelection, error) {

	query := `SELECT
				bu.kode_bidang_urusan,
				COALESCE(bu.nama_bidang_urusan, '') AS nama_bidang_urusan,
				od.kode_opd,
				COALESCE(od.nama_opd, '') AS nama_opd
			FROM tb_bidang_urusan bu
			CROSS JOIN tb_operasional_daerah od
			WHERE od.kode_opd = ?`

	rows, err := tx.QueryContext(ctx, query, kodeOpd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selections := make([]domain.BidangUrusanSelection, 0)

	for rows.Next() {
		var selection domain.BidangUrusanSelection

		err := rows.Scan(
			&selection.KodeBidangUrusan,
			&selection.NamaBidangUrusan,
			&selection.KodeOpd,
			&selection.NamaOpd,
		)

		if err != nil {
			return nil, err
		}

		selections = append(selections, selection)
	}

	return selections, nil
}


