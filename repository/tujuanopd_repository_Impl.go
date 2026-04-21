package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TujuanOpdRepositoryImpl struct {
}

func NewTujuanOpdRepositoryImpl() *TujuanOpdRepositoryImpl {
	return &TujuanOpdRepositoryImpl{}
}

func (repository *TujuanOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) (domain.TujuanOpd, error) {
	script := "INSERT INTO tb_tujuan_opd (kode_opd, kode_bidang_urusan, tujuan, periode_id, tahun_awal, tahun_akhir, jenis_periode) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		tujuanOpd.KodeOpd,
		tujuanOpd.KodeBidangUrusan,
		tujuanOpd.Tujuan,
		tujuanOpd.PeriodeId.Id,
		tujuanOpd.TahunAwal,
		tujuanOpd.TahunAkhir,
		tujuanOpd.JenisPeriode,
	)
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	tujuanOpd.Id = int(id)
	for _, indikator := range tujuanOpd.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator_matrix
            (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan, sumber_data, definisi_operasional, jenis)
            VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.KodeIndikator,
			id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData,
			indikator.DefinisiOperasional,
			indikator.Jenis,
		)
		if err != nil {
			return domain.TujuanOpd{}, err
		}
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.KodeIndikator, // FIX: pakai kode_indikator (VARCHAR), bukan indikator.Id
				target.Target,
				target.Satuan,
				target.Tahun,
			)
			if err != nil {
				return domain.TujuanOpd{}, err
			}
		}
	}
	return tujuanOpd, nil
}

func (repository *TujuanOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tujuanOpd domain.TujuanOpd) error {
	// Update tujuan OPD utama
	script := `UPDATE tb_tujuan_opd
        SET kode_opd = ?, kode_bidang_urusan = ?, tujuan = ?,
            periode_id = ?, tahun_awal = ?, tahun_akhir = ?, jenis_periode = ?
        WHERE id = ?`
	_, err := tx.ExecContext(ctx, script,
		tujuanOpd.KodeOpd,
		tujuanOpd.KodeBidangUrusan,
		tujuanOpd.Tujuan,
		tujuanOpd.PeriodeId.Id,
		tujuanOpd.TahunAwal,
		tujuanOpd.TahunAkhir,
		tujuanOpd.JenisPeriode,
		tujuanOpd.Id,
	)
	if err != nil {
		return err
	}
	// Hapus target lama (join via kode_indikator, bukan i.id)
	scriptDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator_matrix i
            ON t.indikator_id = i.kode_indikator  -- FIX: kode_indikator bukan i.id
        WHERE i.tujuan_opd_id = ?
    `
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, tujuanOpd.Id)
	if err != nil {
		return err
	}
	// Hapus indikator lama
	scriptDeleteIndikator := "DELETE FROM tb_indikator_matrix WHERE tujuan_opd_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, tujuanOpd.Id)
	if err != nil {
		return err
	}
	// Insert indikator dan target baru
	for _, indikator := range tujuanOpd.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator_matrix
            (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan, sumber_data, definisi_operasional, jenis)
            VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.KodeIndikator,
			tujuanOpd.Id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData,
			indikator.DefinisiOperasional,
			indikator.Jenis,
		)
		if err != nil {
			return err
		}
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.KodeIndikator, // FIX: pakai kode_indikator (VARCHAR)
				target.Target,
				target.Satuan,
				target.Tahun,
			)
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
        INNER JOIN tb_indikator_matrix i ON t.indikator_id = i.kode_indikator
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
	script := `
        SELECT 
            t.id, 
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
            t.tujuan, 
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode,
            i.id as indikator_id,
            i.kode_indikator,
            i.indikator,
            i.rumus_perhitungan, 
            i.definisi_operasional,
            i.sumber_data,
            COALESCE(tg.id, '') as target_id,
            COALESCE(tg.target, '') as target_value,
            COALESCE(tg.satuan, '') as satuan,
            COALESCE(tg.tahun, '') as tahun_target
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator_matrix i ON t.id = i.tujuan_opd_id
        LEFT JOIN tb_target tg ON i.kode_indikator = tg.indikator_id
        WHERE t.id = ?
        ORDER BY i.id ASC, tg.tahun ASC
    `

	rows, err := tx.QueryContext(ctx, script, tujuanOpdId)
	if err != nil {
		return domain.TujuanOpd{}, err
	}
	defer rows.Close()

	var tujuanOpd *domain.TujuanOpd
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			id                  int
			kodeOpd             string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwal           string
			tahunAkhir          string
			jenisPeriode        string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)

		err := rows.Scan(
			&id,
			&kodeOpd,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwal,
			&tahunAkhir,
			&jenisPeriode,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return domain.TujuanOpd{}, err
		}

		if tujuanOpd == nil {
			tujuanOpd = &domain.TujuanOpd{
				Id:               id,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwal,
				TahunAkhir:       tahunAkhir,
				JenisPeriode:     jenisPeriode,
				Indikator:        []domain.Indikator{},
			}
		}

		if indikatorId.Valid {
			if _, exists := indikatorMap[indikatorId.String]; !exists {
				indikatorMap[indikatorId.String] = &domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional,
					Target:              []domain.Target{},
				}
				tujuanOpd.Indikator = append(tujuanOpd.Indikator, *indikatorMap[indikatorId.String])
			}

			if targetId.Valid && tahunTarget.Valid {
				tahunTargetInt, _ := strconv.Atoi(tahunTarget.String)
				tahunAwalInt, _ := strconv.Atoi(tahunAwal)
				tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

				// Hanya tambahkan target jika tahunnya dalam range
				if tahunTargetInt >= tahunAwalInt && tahunTargetInt <= tahunAkhirInt {
					target := domain.Target{
						Id:          targetId.String,
						IndikatorId: kodeIndikator.String,
						Target:      targetValue.String,
						Satuan:      satuan.String,
						Tahun:       tahunTarget.String,
					}
					for i := range tujuanOpd.Indikator {
						if tujuanOpd.Indikator[i].Id == indikatorId.String {
							tujuanOpd.Indikator[i].Target = append(tujuanOpd.Indikator[i].Target, target)
							break
						}
					}
				}
			}
		}
	}

	if tujuanOpd == nil {
		return domain.TujuanOpd{}, fmt.Errorf("tujuan opd with id %d not found", tujuanOpdId)
	}

	// Generate target lengkap untuk setiap indikator
	for i := range tujuanOpd.Indikator {
		tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
		tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)

		// Buat map untuk target yang sudah ada
		existingTargets := make(map[string]domain.Target)
		for _, target := range tujuanOpd.Indikator[i].Target {
			existingTargets[target.Tahun] = target
		}

		// Reset dan generate ulang target
		var completeTargets []domain.Target
		for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
			tahunStr := strconv.Itoa(tahun)
			if target, exists := existingTargets[tahunStr]; exists {
				completeTargets = append(completeTargets, target)
			} else {
				completeTargets = append(completeTargets, domain.Target{
					Id:          "",
					IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
					Target:      "",
					Satuan:      "",
					Tahun:       tahunStr,
				})
			}
		}

		// Sort target berdasarkan tahun
		sort.Slice(completeTargets, func(i, j int) bool {
			tahunI, _ := strconv.Atoi(completeTargets[i].Tahun)
			tahunJ, _ := strconv.Atoi(completeTargets[j].Tahun)
			return tahunI < tahunJ
		})

		tujuanOpd.Indikator[i].Target = completeTargets
	}

	return *tujuanOpd, nil
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
        AND tahun = ?
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

func (repository *TujuanOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	scriptTujuan := `
        SELECT 
            t.id, 
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
            t.tujuan, 
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode,
            i.id as indikator_id,
            i.indikator,
            i.rumus_perhitungan, 
            i.sumber_data,
            tg.id as target_id,
            tg.target as target_value,
            tg.satuan,
            tg.tahun as tahun_target
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator i ON t.id = i.tujuan_opd_id
        LEFT JOIN tb_target tg ON i.id = tg.indikator_id
        WHERE t.kode_opd = ? 
        AND t.tahun_awal = ?
        AND t.tahun_akhir = ?
        AND t.jenis_periode = ?
        AND (tg.tahun IS NULL OR (CAST(tg.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)))
        ORDER BY t.id ASC, i.id ASC, tg.tahun ASC
    `

	rows, err := tx.QueryContext(ctx, scriptTujuan,
		kodeOpd,
		tahunAwal,
		tahunAkhir,
		jenisPeriode,
		tahunAwal,
		tahunAkhir,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			tujuanId         int
			kodeOpd          string
			kodeBidangUrusan string
			tujuan           string
			// periodeId        int
			tahunAwalData    string
			tahunAkhirData   string
			jenisPeriodeData string
			indikatorId      sql.NullString
			indikatorNama    sql.NullString
			rumusPerhitungan sql.NullString
			sumberData       sql.NullString
			targetId         sql.NullString
			targetValue      sql.NullString
			satuan           sql.NullString
			tahunTarget      sql.NullString
		)

		err := rows.Scan(
			&tujuanId,
			&kodeOpd,
			&kodeBidangUrusan,
			&tujuan,
			// &periodeId,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}

		// Buat atau ambil TujuanOpd
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
		}

		// Buat atau ambil Indikator jika ada
		if indikatorId.Valid {
			if _, exists := indikatorMap[indikatorId.String]; !exists {
				indikatorMap[indikatorId.String] = &domain.Indikator{
					Id:               indikatorId.String,
					Indikator:        indikatorNama.String,
					RumusPerhitungan: rumusPerhitungan,
					SumberData:       sumberData,
					Target:           []domain.Target{},
				}
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, *indikatorMap[indikatorId.String])
			}

			// Tambahkan target jika ada
			if targetId.Valid && tahunTarget.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				// Update target langsung ke indikator di map
				for i := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[i].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[i].Target = append(tujuanOpdMap[tujuanId].Indikator[i].Target, target)
						break
					}
				}
			}
		}
	}

	// Perbaikan pada bagian generate target
	var result []domain.TujuanOpd
	for _, tujuanOpd := range tujuanOpdMap {
		for i := range tujuanOpd.Indikator {
			tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)

			// Buat map untuk target yang sudah ada
			existingTargets := make(map[string]domain.Target)
			for _, target := range tujuanOpd.Indikator[i].Target {
				if target.Id != "" { // Hanya masukkan target yang valid
					existingTargets[target.Tahun] = target
				}
			}

			// Reset target array
			var completeTargets []domain.Target

			// Generate target untuk setiap tahun dalam range
			for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
				tahunStr := strconv.Itoa(tahun)
				if target, exists := existingTargets[tahunStr]; exists {
					// Gunakan target yang sudah ada
					completeTargets = append(completeTargets, target)
				} else {
					// Buat target kosong dengan tahun yang sesuai
					completeTargets = append(completeTargets, domain.Target{
						Id:          "",
						IndikatorId: tujuanOpd.Indikator[i].Id,
						Target:      "",
						Satuan:      "",
						Tahun:       tahunStr,
					})
				}
			}

			// Sort target berdasarkan tahun
			sort.Slice(completeTargets, func(i, j int) bool {
				tahunI, _ := strconv.Atoi(completeTargets[i].Tahun)
				tahunJ, _ := strconv.Atoi(completeTargets[j].Tahun)
				return tahunI < tahunJ
			})

			tujuanOpd.Indikator[i].Target = completeTargets
		}
		result = append(result, *tujuanOpd)
	}

	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	scriptTujuan := `
        SELECT 
            t.id, 
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
            t.tujuan, 
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode,
            i.id as indikator_id,
            i.indikator,
            i.rumus_perhitungan, 
            i.sumber_data,
            tg.id as target_id,
            tg.target as target_value,
            tg.satuan,
            tg.tahun as tahun_target
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator i ON t.id = i.tujuan_opd_id
        LEFT JOIN tb_target tg ON i.id = tg.indikator_id
        WHERE t.kode_opd = ? 
        AND CAST(? AS SIGNED) BETWEEN CAST(t.tahun_awal AS SIGNED) AND CAST(t.tahun_akhir AS SIGNED)
        AND t.jenis_periode = ?
        AND (tg.tahun IS NULL OR tg.tahun = ?)
        ORDER BY t.id ASC, i.id ASC, tg.tahun ASC
    `

	rows, err := tx.QueryContext(ctx, scriptTujuan,
		kodeOpd,
		tahun,
		jenisPeriode,
		tahun,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			tujuanId         int
			kodeOpd          string
			kodeBidangUrusan string
			tujuan           string
			tahunAwalData    string
			tahunAkhirData   string
			jenisPeriodeData string
			indikatorId      sql.NullString
			indikatorNama    sql.NullString
			rumusPerhitungan sql.NullString
			sumberData       sql.NullString
			targetId         sql.NullString
			targetValue      sql.NullString
			satuan           sql.NullString
			tahunTarget      sql.NullString
		)

		err := rows.Scan(
			&tujuanId,
			&kodeOpd,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}

		// Buat atau ambil TujuanOpd
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
		}

		// Buat atau ambil Indikator jika ada
		if indikatorId.Valid {
			if _, exists := indikatorMap[indikatorId.String]; !exists {
				indikatorMap[indikatorId.String] = &domain.Indikator{
					Id:               indikatorId.String,
					Indikator:        indikatorNama.String,
					RumusPerhitungan: rumusPerhitungan,
					SumberData:       sumberData,
					Target:           []domain.Target{},
				}
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, *indikatorMap[indikatorId.String])
			}

			// Tambahkan target jika ada
			if targetId.Valid && tahunTarget.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				// Update target langsung ke indikator di map
				for i := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[i].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[i].Target = append(tujuanOpdMap[tujuanId].Indikator[i].Target, target)
						break
					}
				}
			}
		}
	}

	var result []domain.TujuanOpd
	for _, tujuanOpd := range tujuanOpdMap {
		// Sort target berdasarkan tahun untuk setiap indikator
		for i := range tujuanOpd.Indikator {
			sort.Slice(tujuanOpd.Indikator[i].Target, func(x, y int) bool {
				tahunX, _ := strconv.Atoi(tujuanOpd.Indikator[i].Target[x].Tahun)
				tahunY, _ := strconv.Atoi(tujuanOpd.Indikator[i].Target[y].Tahun)
				return tahunX < tahunY
			})
		}
		result = append(result, *tujuanOpd)
	}

	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}

	return result, nil
}

// Perbaikan pada FindIndikatorByTujuanOpdId untuk menyertakan rumus_perhitungan dan sumber_data
func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanOpdId(ctx context.Context, tx *sql.Tx, tujuanOpdId int) ([]domain.Indikator, error) {
	script := `
        SELECT 
            id, 
            indikator,
            COALESCE(rumus_perhitungan, '') as rumus_perhitungan,
            COALESCE(sumber_data, '') as sumber_data
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
		var rumusPerhitungan, sumberData string

		err := rows.Scan(
			&indikator.Id,
			&indikator.Indikator,
			&rumusPerhitungan,
			&sumberData,
		)
		if err != nil {
			return nil, err
		}

		indikator.TujuanOpdId = tujuanOpdId
		indikator.RumusPerhitungan = sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""}
		indikator.SumberData = sql.NullString{String: sumberData, Valid: sumberData != ""}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *TujuanOpdRepositoryImpl) FindTujuanOpdForCascadingOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	script := `
		SELECT
			t.id,
			t.kode_opd,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode,
			t.kode_bidang_urusan
		FROM tb_tujuan_opd t
		INNER JOIN tb_periode p ON 
			t.tahun_awal = p.tahun_awal 
			AND t.tahun_akhir = p.tahun_akhir 
			AND t.jenis_periode = p.jenis_periode
		WHERE t.kode_opd = ?
		AND CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
		AND p.jenis_periode = ?
		ORDER BY t.id ASC
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun, jenisPeriode)
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
			&tujuanOpd.JenisPeriode,
			&tujuanOpd.KodeBidangUrusan,
		)
		if err != nil {
			return nil, err
		}
		tujuanOpds = append(tujuanOpds, tujuanOpd)
	}

	// Untuk setiap tujuan OPD, ambil indikatornya
	for i := range tujuanOpds {
		indikators, err := repository.FindIndikatorByTujuanOpdId(ctx, tx, tujuanOpds[i].Id)
		if err != nil {
			return nil, err
		}
		tujuanOpds[i].Indikator = indikators

		// Generate target untuk setiap indikator
		for j := range tujuanOpds[i].Indikator {
			tahunAwalInt, _ := strconv.Atoi(tujuanOpds[i].TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuanOpds[i].TahunAkhir)

			var targets []domain.Target
			for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
				targets = append(targets, domain.Target{
					Id:          "",
					IndikatorId: tujuanOpds[i].Indikator[j].Id,
					Target:      "",
					Satuan:      "",
					Tahun:       strconv.Itoa(tahun),
				})
			}
			tujuanOpds[i].Indikator[j].Target = targets
		}
	}

	return tujuanOpds, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorByTujuanOpdIdsBatch(ctx context.Context, tx *sql.Tx, tujuanOpdIds []int) (map[int][]domain.Indikator, error) {
	if len(tujuanOpdIds) == 0 {
		return make(map[int][]domain.Indikator), nil
	}

	// Build query dengan IN clause
	placeholders := make([]string, len(tujuanOpdIds))
	args := make([]interface{}, len(tujuanOpdIds))
	for i, id := range tujuanOpdIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT 
			id, 
			tujuan_opd_id,
			indikator,
			COALESCE(rumus_perhitungan, '') as rumus_perhitungan,
			COALESCE(sumber_data, '') as sumber_data
		FROM tb_indikator
		WHERE tujuan_opd_id IN (%s)
		ORDER BY tujuan_opd_id ASC, id ASC
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Group indikators by tujuan_opd_id
	result := make(map[int][]domain.Indikator)
	for rows.Next() {
		var indikator domain.Indikator
		var tujuanOpdId int
		err := rows.Scan(&indikator.Id, &tujuanOpdId, &indikator.Indikator, &indikator.RumusPerhitungan, &indikator.SumberData)
		if err != nil {
			return nil, err
		}
		result[tujuanOpdId] = append(result[tujuanOpdId], indikator)
	}

	return result, nil
}

// renstra new
func (repository *TujuanOpdRepositoryImpl) FindAllByPeriod(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahunAwal, tahunAkhir, jenisPeriode, jenisIndikator string,
) ([]domain.TujuanOpd, error) {
	jenisClause := ""
	var finalArgs []interface{}
	if jenisIndikator != "" {
		jenisClause = "AND im.jenis = ?"
		finalArgs = append(finalArgs, jenisIndikator) // 1: im.jenis = ?
	}
	finalArgs = append(finalArgs, tahunAwal, tahunAkhir)                        // 2,3: BETWEEN target
	finalArgs = append(finalArgs, kodeOpd, tahunAwal, tahunAkhir, jenisPeriode) // 4,5,6,7: WHERE tujuan
	query := fmt.Sprintf(`
        SELECT
            t.id,
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '')         AS kode_bidang_urusan,
            t.tujuan,
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode,
            im.id                                      AS indikator_id,
            im.kode_indikator                           AS kode_indikator,
            im.indikator,
            COALESCE(im.rumus_perhitungan, '')         AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')               AS sumber_data,
            COALESCE(im.definisi_operasional, '')      AS definisi_operasional,
            COALESCE(im.jenis, '')                     AS indikator_jenis,
            tg.id                                      AS target_id,
            tg.target                                  AS target_value,
            tg.satuan,
            tg.tahun                                   AS tahun_target
        FROM tb_tujuan_opd t
        LEFT JOIN tb_indikator_matrix im
            ON t.id = im.tujuan_opd_id %s
        LEFT JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND CAST(tg.tahun AS SIGNED) BETWEEN CAST(? AS SIGNED) AND CAST(? AS SIGNED)
        WHERE t.kode_opd      = ?
          AND t.tahun_awal    = ?
          AND t.tahun_akhir   = ?
          AND t.jenis_periode = ?
        ORDER BY t.id, im.id, tg.tahun
    `, jenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid && tahunTarget.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Renstra: generate slot target untuk setiap tahun dalam range
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			tahunAwalInt, _ := strconv.Atoi(tujuanOpd.TahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tujuanOpd.TahunAkhir)
			existingTargets := make(map[string]domain.Target)
			for _, t := range tujuanOpd.Indikator[i].Target {
				if t.Id != "" {
					existingTargets[t.Tahun] = t
				}
			}
			var completeTargets []domain.Target
			for year := tahunAwalInt; year <= tahunAkhirInt; year++ {
				yearStr := strconv.Itoa(year)
				if t, ok := existingTargets[yearStr]; ok {
					completeTargets = append(completeTargets, t)
				} else {
					completeTargets = append(completeTargets, domain.Target{
						Id:          "",
						IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
						Target:      "",
						Satuan:      "",
						Tahun:       yearStr,
					})
				}
			}
			tujuanOpd.Indikator[i].Target = completeTargets
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAllByTahun(
	ctx context.Context, tx *sql.Tx,
	kodeOpd, tahun, jenisPeriode, jenisIndikator string,
) ([]domain.TujuanOpd, error) {
	var finalArgs []interface{}
	// Args untuk subquery (dalam tanda kurung)
	finalArgs = append(finalArgs, tahun) // tg.tahun = ?
	jenisClause := ""
	if jenisIndikator != "" {
		jenisClause = "AND im.jenis = ?"
		finalArgs = append(finalArgs, jenisIndikator) // im.jenis = ?
	}
	// Args untuk WHERE tujuan_opd
	finalArgs = append(finalArgs, kodeOpd)      // t.kode_opd = ?
	finalArgs = append(finalArgs, jenisPeriode) // t.jenis_periode = ?
	finalArgs = append(finalArgs, tahun, tahun) // tahun_awal <= ? AND tahun_akhir >= ?
	query := fmt.Sprintf(`
    SELECT
        t.id,
        t.kode_opd,
        COALESCE(t.kode_bidang_urusan, '')      AS kode_bidang_urusan,
        t.tujuan,
        t.tahun_awal,
        t.tahun_akhir,
        t.jenis_periode,
        im_tg.indikator_id,
        im_tg.kode_indikator,
        im_tg.indikator,
        im_tg.rumus_perhitungan,
        im_tg.sumber_data,
        im_tg.definisi_operasional,
        im_tg.indikator_jenis,
        im_tg.target_id,
        im_tg.target_value,
        im_tg.satuan,
        im_tg.tahun_target
    FROM tb_tujuan_opd t
    LEFT JOIN (
        SELECT
            im.id                                  AS indikator_id,
            im.kode_indikator                      AS kode_indikator,
            im.tujuan_opd_id,
            COALESCE(im.indikator, '')             AS indikator,
            COALESCE(im.rumus_perhitungan, '')     AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')           AS sumber_data,
            COALESCE(im.definisi_operasional, '')  AS definisi_operasional,
            COALESCE(im.jenis, '')                 AS indikator_jenis,
            tg.id                                  AS target_id,
            tg.target                              AS target_value,
            tg.satuan,
            tg.tahun                               AS tahun_target
        FROM tb_indikator_matrix im
        INNER JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND tg.tahun = ?
        %s
    ) im_tg ON t.id = im_tg.tujuan_opd_id
    WHERE t.kode_opd     = ?
      AND t.jenis_periode = ?
      AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
      AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
    ORDER BY t.id, im_tg.indikator_id
`, jenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Ranwal / Rankhir: tepat 1 slot target per indikator untuk tahun yang diminta
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			if len(tujuanOpd.Indikator[i].Target) == 0 {
				// Belum ada data target di DB → buat slot kosong
				tujuanOpd.Indikator[i].Target = []domain.Target{{
					Id:          "",
					IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
					Target:      "",
					Satuan:      "",
					Tahun:       tahun,
				}}
			}
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindBidangUrusanBatch(
	ctx context.Context, tx *sql.Tx,
	kodeBidangUrusanList []string,
) (map[string]domainmaster.BidangUrusan, error) {
	result := make(map[string]domainmaster.BidangUrusan)
	if len(kodeBidangUrusanList) == 0 {
		return result, nil
	}
	// Deduplicate & filter kosong
	uniqueSet := make(map[string]struct{})
	for _, k := range kodeBidangUrusanList {
		if k != "" {
			uniqueSet[k] = struct{}{}
		}
	}
	if len(uniqueSet) == 0 {
		return result, nil
	}
	placeholders := make([]string, 0, len(uniqueSet))
	args := make([]interface{}, 0, len(uniqueSet))
	for k := range uniqueSet {
		placeholders = append(placeholders, "?")
		args = append(args, k)
	}
	// JOIN ke tb_urusan menggunakan digit pertama kode_bidang_urusan
	// Sama persis dengan pola FindByKodeBidangUrusan yang sudah jalan
	query := fmt.Sprintf(`
        SELECT
            COALESCE(bu.id, '')                  AS id,
            COALESCE(bu.kode_bidang_urusan, '')   AS kode_bidang_urusan,
            COALESCE(bu.nama_bidang_urusan, '')   AS nama_bidang_urusan,
            COALESCE(u.kode_urusan, '')           AS kode_urusan,
            COALESCE(u.nama_urusan, '')           AS nama_urusan
        FROM tb_bidang_urusan bu
        LEFT JOIN tb_urusan u
            ON LEFT(bu.kode_bidang_urusan, 1) = u.kode_urusan
        WHERE bu.kode_bidang_urusan IN (%s)
    `, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var bu domainmaster.BidangUrusan
		if err := rows.Scan(
			&bu.Id,
			&bu.KodeBidangUrusan,
			&bu.NamaBidangUrusan,
			&bu.KodeUrusan,
			&bu.NamaUrusan,
		); err != nil {
			return nil, err
		}
		result[bu.KodeBidangUrusan] = bu
	}
	return result, rows.Err()
}

func (r *TujuanOpdRepositoryImpl) CreateRenjaIndikator(
	ctx context.Context, tx *sql.Tx,
	tujuanOpdId int, indikators []domain.Indikator,
) error {
	for _, ind := range indikators {
		_, err := tx.ExecContext(ctx, `
            INSERT INTO tb_indikator_matrix
                (kode_indikator, tujuan_opd_id, indikator, rumus_perhitungan,
                 sumber_data, definisi_operasional, jenis)
            VALUES (?, ?, ?, ?, ?, ?, ?)`,
			ind.KodeIndikator, tujuanOpdId,
			ind.Indikator, ind.RumusPerhitungan.String,
			ind.SumberData.String, ind.DefinisiOperasional.String, ind.Jenis,
		)
		if err != nil {
			return err
		}
		// 1 target per indikator
		t := ind.Target[0]
		_, err = tx.ExecContext(ctx,
			"INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)",
			t.Id, ind.KodeIndikator, t.Target, t.Satuan, t.Tahun,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// UPDATE: hanya UPDATE, kode_indikator wajib ada
func (r *TujuanOpdRepositoryImpl) UpdateRenjaIndikator(
	ctx context.Context, tx *sql.Tx,
	indikators []domain.Indikator,
) error {
	for _, ind := range indikators {
		_, err := tx.ExecContext(ctx, `
            UPDATE tb_indikator_matrix
            SET indikator = ?, rumus_perhitungan = ?, sumber_data = ?,
                definisi_operasional = ?, jenis = ?
            WHERE kode_indikator = ?`,
			ind.Indikator, ind.RumusPerhitungan.String,
			ind.SumberData.String, ind.DefinisiOperasional.String,
			ind.Jenis, ind.KodeIndikator,
		)
		if err != nil {
			return err
		}
		// DELETE target lama yang sama tahunnya + INSERT baru
		t := ind.Target[0]
		_, err = tx.ExecContext(ctx,
			"DELETE FROM tb_target WHERE indikator_id = ? AND tahun = ?",
			ind.KodeIndikator, t.Tahun,
		)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			"INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)",
			t.Id, ind.KodeIndikator, t.Target, t.Satuan, t.Tahun,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TujuanOpdRepositoryImpl) DeleteIndikatorTargetRenja(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	// 1. Hapus Child terlebih dahulu (tb_target)
	_, err := tx.ExecContext(ctx, "DELETE FROM tb_target WHERE indikator_id = ?", indikatorId)
	if err != nil {
		return err
	}

	// 2. Hapus Parent (tb_indikator_matrix)
	_, err = tx.ExecContext(ctx, "DELETE FROM tb_indikator_matrix WHERE kode_indikator = ?", indikatorId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TujuanOpdRepositoryImpl) FindIndikatorByKodeIndikator(
	ctx context.Context, tx *sql.Tx, kodeIndikator string,
) (domain.Indikator, error) {
	row := tx.QueryRowContext(ctx, `
        SELECT kode_indikator,
               COALESCE(indikator, ''),
               COALESCE(rumus_perhitungan, ''),
               COALESCE(sumber_data, ''),
               COALESCE(definisi_operasional, ''),
               COALESCE(jenis, '')
        FROM tb_indikator_matrix
        WHERE kode_indikator = ?`,
		kodeIndikator,
	)
	var indikator domain.Indikator
	err := row.Scan(
		&indikator.KodeIndikator, // ← tidak scan id sama sekali
		&indikator.Indikator,
		&indikator.RumusPerhitungan,
		&indikator.SumberData,
		&indikator.DefinisiOperasional,
		&indikator.Jenis,
	)
	if err != nil {
		return domain.Indikator{}, err
	}
	return indikator, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAllByTahunForPokin(ctx context.Context, tx *sql.Tx, kodeOpd, tahun, jenisPeriode, jenisIndikator string) ([]domain.TujuanOpd, error) {
	var finalArgs []interface{}
	// Args untuk subquery (dalam tanda kurung)
	finalArgs = append(finalArgs, tahun) // tg.tahun = ?
	jenisClause := ""
	if jenisIndikator != "" {
		jenisClause = "AND im.jenis = ?"
		finalArgs = append(finalArgs, jenisIndikator) // im.jenis = ?
	}
	// Args untuk WHERE tujuan_opd
	finalArgs = append(finalArgs, kodeOpd)      // t.kode_opd = ?
	finalArgs = append(finalArgs, jenisPeriode) // t.jenis_periode = ?
	finalArgs = append(finalArgs, tahun, tahun) // tahun_awal <= ? AND tahun_akhir >= ?
	query := fmt.Sprintf(`
    SELECT
        t.id,
        t.kode_opd,
        COALESCE(t.kode_bidang_urusan, '')      AS kode_bidang_urusan,
        t.tujuan,
        t.tahun_awal,
        t.tahun_akhir,
        t.jenis_periode,
        im_tg.indikator_id,
        im_tg.kode_indikator,
        im_tg.indikator,
        im_tg.rumus_perhitungan,
        im_tg.sumber_data,
        im_tg.definisi_operasional,
        im_tg.indikator_jenis,
        im_tg.target_id,
        im_tg.target_value,
        im_tg.satuan,
        im_tg.tahun_target
    FROM tb_tujuan_opd t
    LEFT JOIN (
        SELECT
            im.id                                  AS indikator_id,
            im.kode_indikator                      AS kode_indikator,
            im.tujuan_opd_id,
            COALESCE(im.indikator, '')             AS indikator,
            COALESCE(im.rumus_perhitungan, '')     AS rumus_perhitungan,
            COALESCE(im.sumber_data, '')           AS sumber_data,
            COALESCE(im.definisi_operasional, '')  AS definisi_operasional,
            COALESCE(im.jenis, '')                 AS indikator_jenis,
            tg.id                                  AS target_id,
            tg.target                              AS target_value,
            tg.satuan,
            tg.tahun                               AS tahun_target
        FROM tb_indikator_matrix im
        LEFT JOIN tb_target tg
            ON im.kode_indikator = tg.indikator_id
            AND tg.tahun = ?
        %s
    ) im_tg ON t.id = im_tg.tujuan_opd_id
    WHERE t.kode_opd     = ?
      AND t.jenis_periode = ?
      AND CAST(t.tahun_awal  AS SIGNED) <= CAST(? AS SIGNED)
      AND CAST(t.tahun_akhir AS SIGNED) >= CAST(? AS SIGNED)
    ORDER BY t.id, im_tg.indikator_id
`, jenisClause)
	rows, err := tx.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tujuanOpdMap := make(map[int]*domain.TujuanOpd)
	indikatorSeen := make(map[string]bool) // key: "tujuanId-indikatorId"
	tujuanOrder := []int{}
	for rows.Next() {
		var (
			tujuanId            int
			kodeOpdData         string
			kodeBidangUrusan    string
			tujuan              string
			tahunAwalData       string
			tahunAkhirData      string
			jenisPeriodeData    string
			indikatorId         sql.NullString
			kodeIndikator       sql.NullString
			indikatorNama       sql.NullString
			rumusPerhitungan    sql.NullString
			sumberData          sql.NullString
			definisiOperasional sql.NullString // NEW
			indikatorJenis      sql.NullString // im.jenis
			targetId            sql.NullString
			targetValue         sql.NullString
			satuan              sql.NullString
			tahunTarget         sql.NullString
		)
		err := rows.Scan(
			&tujuanId,
			&kodeOpdData,
			&kodeBidangUrusan,
			&tujuan,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
			&indikatorId,
			&kodeIndikator,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&definisiOperasional, // NEW
			&indikatorJenis,
			&targetId,
			&targetValue,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpdData,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
			tujuanOrder = append(tujuanOrder, tujuanId)
		}
		if indikatorId.Valid {
			mapKey := fmt.Sprintf("%d-%s", tujuanId, indikatorId.String)
			if !indikatorSeen[mapKey] {
				indikatorSeen[mapKey] = true
				tujuanOpdMap[tujuanId].Indikator = append(tujuanOpdMap[tujuanId].Indikator, domain.Indikator{
					Id:                  indikatorId.String,
					KodeIndikator:       kodeIndikator.String,
					Indikator:           indikatorNama.String,
					RumusPerhitungan:    rumusPerhitungan,
					SumberData:          sumberData,
					DefinisiOperasional: definisiOperasional, // NEW
					Jenis:               indikatorJenis.String,
					TujuanOpdId:         tujuanId,
					Target:              []domain.Target{},
				})
			}
			if targetId.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: kodeIndikator.String,
					Target:      targetValue.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				for idx := range tujuanOpdMap[tujuanId].Indikator {
					if tujuanOpdMap[tujuanId].Indikator[idx].Id == indikatorId.String {
						tujuanOpdMap[tujuanId].Indikator[idx].Target = append(
							tujuanOpdMap[tujuanId].Indikator[idx].Target, target,
						)
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Ranwal / Rankhir: tepat 1 slot target per indikator untuk tahun yang diminta
	var result []domain.TujuanOpd
	for _, id := range tujuanOrder {
		tujuanOpd := tujuanOpdMap[id]
		for i := range tujuanOpd.Indikator {
			if len(tujuanOpd.Indikator[i].Target) == 0 {
				// Belum ada data target di DB → buat slot kosong
				tujuanOpd.Indikator[i].Target = []domain.Target{{
					Id:          "",
					IndikatorId: tujuanOpd.Indikator[i].KodeIndikator,
					Target:      "",
					Satuan:      "",
					Tahun:       tahun,
				}}
			}
		}
		result = append(result, *tujuanOpd)
	}
	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}
	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorTargetsRenstraByTujuanIds(
	ctx context.Context,
	tx *sql.Tx,
	tujuanIds []int,
) ([]domain.Indikator, error) {

	if len(tujuanIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// buat placeholder (?, ?, ?, ...)
	placeholders := make([]string, len(tujuanIds))
	args := make([]any, len(tujuanIds))

	for i, id := range tujuanIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT ind.kode_indikator, ind.indikator, ind.tujuan_opd_id,
                     ind.rumus_perhitungan, ind.sumber_data,
                     ind.definisi_operasional,
                     tar.id, tar.target, tar.satuan, tar.tahun, tar.jenis
		FROM tb_indikator_matrix ind
                LEFT JOIN tb_target tar ON tar.indikator_id = ind.kode_indikator
		WHERE ind.tujuan_opd_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId                 string
			indikatorName         string
			tujuanOpdId           int
			rumusPerhitunganNS    sql.NullString
			sumberDataNS          sql.NullString
			definisiOperasionalNS sql.NullString

			tarIdNS  sql.NullString
			targetNS sql.NullString
			satuanNS sql.NullString
			tahunNS  sql.NullString
			jenisNS  sql.NullString
		)

		err := rows.Scan(
			&indId,
			&indikatorName,
			&tujuanOpdId,
			&rumusPerhitunganNS,
			&sumberDataNS,
			&definisiOperasionalNS,
			&tarIdNS,
			&targetNS,
			&satuanNS,
			&tahunNS,
			&jenisNS,
		)
		if err != nil {
			return nil, err
		}

		// 🔹 ambil / buat indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:                  indId,
				KodeIndikator:       indId,
				Indikator:           indikatorName,
				TujuanOpdId:         tujuanOpdId,
				RumusPerhitungan:    rumusPerhitunganNS,
				DefinisiOperasional: definisiOperasionalNS,
				SumberData:          sumberDataNS,
				Target:              []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// 🔹 kalau ada target, append
		if tarIdNS.Valid {
			target := domain.Target{
				Id: tarIdNS.String,
			}

			if targetNS.Valid {
				target.Target = targetNS.String
			}
			if satuanNS.Valid {
				target.Satuan = satuanNS.String
			}
			if tahunNS.Valid {
				target.Tahun = tahunNS.String
			}
			if jenisNS.Valid {
				target.Jenis = jenisNS.String
			}

			ind.Target = append(ind.Target, target)
		}
	}

	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, v := range indikatorMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindIndikatorTargetsByTujuanIds(
	ctx context.Context,
	tx *sql.Tx,
	tujuanIds []int,
) ([]domain.Indikator, error) {

	if len(tujuanIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// buat placeholder (?, ?, ?, ...)
	placeholders := make([]string, len(tujuanIds))
	args := make([]any, len(tujuanIds))

	for i, id := range tujuanIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT ind.id, ind.indikator, ind.tujuan_opd_id,
                     ind.rumus_perhitungan, ind.sumber_data,
                     tar.id, tar.target, tar.satuan, tar.tahun, tar.jenis
		FROM tb_indikator ind
                LEFT JOIN tb_target tar ON tar.indikator_id = ind.id
		WHERE ind.tujuan_opd_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId         string
			indikatorName string
			tujuanOpdId   int

			rumusPerhitunganNS sql.NullString
			sumberDataNS       sql.NullString
			tarIdNS            sql.NullString
			targetNS           sql.NullString
			satuanNS           sql.NullString
			tahunNS            sql.NullString
			jenisNS            sql.NullString
		)

		err := rows.Scan(
			&indId,
			&indikatorName,
			&tujuanOpdId,
			&rumusPerhitunganNS,
			&sumberDataNS,
			&tarIdNS,
			&targetNS,
			&satuanNS,
			&tahunNS,
			&jenisNS,
		)
		if err != nil {
			return nil, err
		}

		// 🔹 ambil / buat indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:               indId,
				KodeIndikator:    indId,
				Indikator:        indikatorName,
				TujuanOpdId:      tujuanOpdId,
				RumusPerhitungan: rumusPerhitunganNS,
				SumberData:       sumberDataNS,
				Target:           []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// 🔹 kalau ada target, append
		if tarIdNS.Valid {
			target := domain.Target{
				Id: tarIdNS.String,
			}

			if targetNS.Valid {
				target.Target = targetNS.String
			}
			if satuanNS.Valid {
				target.Satuan = satuanNS.String
			}
			if tahunNS.Valid {
				target.Tahun = tahunNS.String
			}
			if jenisNS.Valid {
				target.Jenis = jenisNS.String
			}

			ind.Target = append(ind.Target, target)
		}
	}

	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, v := range indikatorMap {
		result = append(result, *v)
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindByIdOnly(
	ctx context.Context,
	tx *sql.Tx,
	tujuanOpdId int,
) (domain.TujuanOpd, error) {

	script := `
		SELECT
			t.id,
			t.kode_opd,
			COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
			t.tujuan,
			t.tahun_awal,
			t.tahun_akhir,
			t.jenis_periode
		FROM tb_tujuan_opd t
		WHERE t.id = ?
	`

	var result domain.TujuanOpd

	err := tx.QueryRowContext(ctx, script, tujuanOpdId).Scan(
		&result.Id,
		&result.KodeOpd,
		&result.KodeBidangUrusan,
		&result.Tujuan,
		&result.TahunAwal,
		&result.TahunAkhir,
		&result.JenisPeriode,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TujuanOpd{}, err
		}
		return domain.TujuanOpd{}, err
	}

	return result, nil
}

func (repository *TujuanOpdRepositoryImpl) FindAllOnly(ctx context.Context, tx *sql.Tx, kodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.TujuanOpd, error) {
	scriptTujuan := `
        SELECT
            t.id,
            t.kode_opd,
            COALESCE(t.kode_bidang_urusan, '') as kode_bidang_urusan,
            t.tujuan,
            t.tahun_awal,
            t.tahun_akhir,
            t.jenis_periode
        FROM tb_tujuan_opd t
        WHERE t.kode_opd = ?
        AND t.tahun_awal = ?
        AND t.tahun_akhir = ?
        AND t.jenis_periode = ?
        ORDER BY t.id ASC
    `

	rows, err := tx.QueryContext(ctx, scriptTujuan,
		kodeOpd,
		tahunAwal,
		tahunAkhir,
		jenisPeriode,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tujuanOpdMap := make(map[int]*domain.TujuanOpd)

	for rows.Next() {
		var (
			tujuanId         int
			kodeOpd          string
			kodeBidangUrusan string
			tujuan           string
			// periodeId        int
			tahunAwalData    string
			tahunAkhirData   string
			jenisPeriodeData string
		)

		err := rows.Scan(
			&tujuanId,
			&kodeOpd,
			&kodeBidangUrusan,
			&tujuan,
			// &periodeId,
			&tahunAwalData,
			&tahunAkhirData,
			&jenisPeriodeData,
		)
		if err != nil {
			return nil, err
		}

		// Buat atau ambil TujuanOpd
		if _, exists := tujuanOpdMap[tujuanId]; !exists {
			tujuanOpdMap[tujuanId] = &domain.TujuanOpd{
				Id:               tujuanId,
				KodeOpd:          kodeOpd,
				KodeBidangUrusan: kodeBidangUrusan,
				Tujuan:           tujuan,
				TahunAwal:        tahunAwalData,
				TahunAkhir:       tahunAkhirData,
				JenisPeriode:     jenisPeriodeData,
				Indikator:        []domain.Indikator{},
			}
		}

	}

	// Perbaikan pada bagian generate target
	var result []domain.TujuanOpd
	for _, tujuanOpd := range tujuanOpdMap {
		result = append(result, *tujuanOpd)
	}

	if len(result) == 0 {
		return make([]domain.TujuanOpd, 0), nil
	}

	return result, nil
}
