package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/datamaster"
	"errors"
	"fmt"
	"strings"
)

type DataMasterRepositoryImpl struct {
}

func NewDataMasterRepositoryImpl() *DataMasterRepositoryImpl {
	return &DataMasterRepositoryImpl{}
}

func (repository *DataMasterRepositoryImpl) DataRBByTahun(ctx context.Context, tx *sql.Tx, tahun int) ([]datamaster.MasterRB, error) {

	query := `
		SELECT
			rb.id,
			rb.jenis_rb,
			rb.kegiatan_utama,
			rb.keterangan,
			rb.tahun_baseline,
			rb.tahun_next,

			ind.id AS indikator_id,
			ind.indikator,

			tar.id AS target_id,
			tar.tahun,
			tar.target,
			tar.realisasi,
			tar.satuan

		FROM datamaster_rb rb
		LEFT JOIN tb_indikator ind ON ind.id_rb = rb.id
		LEFT JOIN tb_target tar ON tar.indikator_id = ind.id
		WHERE rb.tahun_next = ? AND rb.is_active = 1
		ORDER BY rb.id, ind.id, tar.id;
	`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rb: %v", err)
	}
	defer rows.Close()

	rbMap := make(map[int]*datamaster.MasterRB)
	indikatorIndexMap := make(map[int]map[string]int)

	for rows.Next() {

		var (
			rb datamaster.MasterRB

			indikatorRb datamaster.IndikatorRB
			targetRb    datamaster.TargetRB

			// FROM datamaster_rb (beberapa bisa NULL: keterangan)
			keterangan sql.NullString

			// FROM tb_indikator (bisa NULL SEMUA)
			indikatorId       sql.NullString
			indikatorNullable sql.NullString

			// FROM tb_target (bisa NULL SEMUA)
			targetId          sql.NullString
			targetYear        sql.NullInt64
			targetValue       sql.NullInt64
			realisasiNullable sql.NullFloat64
			satuanNullable    sql.NullString
		)

		err := rows.Scan(
			&rb.Id,
			&rb.JenisRB,
			&rb.KegiatanUtama,
			&keterangan,
			&rb.TahunBaseline,
			&rb.TahunNext,

			&indikatorId,
			&indikatorNullable,

			&targetId,
			&targetYear,
			&targetValue,
			&realisasiNullable,
			&satuanNullable,
		)

		if err != nil {
			return nil, err
		}

		// =============================
		// Handle fields datamaster_rb
		// =============================
		if keterangan.Valid {
			rb.Keterangan = keterangan.String
		} else {
			rb.Keterangan = ""
		}

		// Tambahkan RB jika belum ada
		if rbMap[rb.Id] == nil {
			rb.Indikator = []datamaster.IndikatorRB{}
			rbMap[rb.Id] = &rb
			indikatorIndexMap[rb.Id] = make(map[string]int)
		}

		master := rbMap[rb.Id]

		// =============================
		// Jika indikator NULL → RB tetap muncul
		// =============================
		if !indikatorId.Valid {
			continue
		}

		indikatorRb.IdRB = rb.Id
		indikatorRb.IdIndikator = indikatorId.String

		if indikatorNullable.Valid {
			indikatorRb.Indikator = indikatorNullable.String
		} else {
			indikatorRb.Indikator = ""
		}

		indikatorRb.TargetRB = []datamaster.TargetRB{}

		// cek apakah indikator sudah ada
		idxMap := indikatorIndexMap[rb.Id]
		indIndex, exists := idxMap[indikatorRb.IdIndikator]
		if !exists {
			master.Indikator = append(master.Indikator, indikatorRb)
			indIndex = len(master.Indikator) - 1
			idxMap[indikatorRb.IdIndikator] = indIndex
		}

		// =============================
		// Jika target NULL → lanjut
		// =============================
		if !targetId.Valid {
			continue
		}

		targetRb.IdTarget = targetId.String
		targetRb.IdIndikator = indikatorId.String

		// Tahun target
		if targetYear.Valid {
			year := int(targetYear.Int64)
			if year == master.TahunBaseline {
				targetRb.TahunBaseline = year
			}
			if year == master.TahunNext {
				targetRb.TahunNext = year
			}
		}

		// Nilai target
		if targetValue.Valid {
			if targetYear.Int64 == int64(master.TahunBaseline) {
				targetRb.TargetBaseline = int(targetValue.Int64)
			}
			if targetYear.Int64 == int64(master.TahunNext) {
				targetRb.TargetNext = int(targetValue.Int64)
			}
		}

		// Realisasi
		if realisasiNullable.Valid {
			targetRb.RealisasiBaseline = float32(realisasiNullable.Float64)
		}

		// Satuan
		if satuanNullable.Valid {
			targetRb.SatuanBaseline = satuanNullable.String
			targetRb.SatuanNext = satuanNullable.String
		}

		// Masukkan target ke indikator
		master.Indikator[indIndex].TargetRB =
			append(master.Indikator[indIndex].TargetRB, targetRb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map → slice
	result := make([]datamaster.MasterRB, 0, len(rbMap))
	for _, v := range rbMap {
		result = append(result, *v)
	}

	if len(result) == 0 {
		return []datamaster.MasterRB{}, nil
	}

	return result, nil
}
func (r *DataMasterRepositoryImpl) InsertRB(ctx context.Context, tx *sql.Tx, req datamaster.MasterRB, userId int) (int64, error) {

	query := `
        INSERT INTO datamaster_rb (
            jenis_rb, kegiatan_utama, keterangan,
            tahun_baseline, tahun_next,
            last_updated_by, is_active
        ) VALUES (?, ?, ?, ?, ?, ?, 1)
    `

	res, err := tx.ExecContext(ctx, query,
		req.JenisRB,
		req.KegiatanUtama,
		req.Keterangan,
		req.TahunBaseline,
		req.TahunNext,
		userId,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *DataMasterRepositoryImpl) InsertIndikator(ctx context.Context, tx *sql.Tx, rbId int64, indikator datamaster.IndikatorRB) (string, error) {

	indikatorID := helper.GenerateID("IND-RB")

	query := `
        INSERT INTO tb_indikator (id, id_rb, indikator)
        VALUES (?, ?, ?)
    `

	_, err := tx.ExecContext(ctx, query,
		indikatorID,
		rbId,
		indikator.Indikator,
	)
	if err != nil {
		return "", err
	}

	return indikatorID, nil
}
func (r *DataMasterRepositoryImpl) InsertTarget(ctx context.Context, tx *sql.Tx, indikatorID string, t datamaster.TargetRB) error {

	query := `
        INSERT INTO tb_target (id, indikator_id, tahun, target, realisasi, satuan)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	// baseline row
	if t.TahunBaseline != 0 {
		targetID := helper.GenerateID("TRGT-RB")

		_, err := tx.ExecContext(ctx, query,
			targetID,
			indikatorID,
			t.TahunBaseline,
			t.TargetBaseline,
			t.RealisasiBaseline,
			t.SatuanBaseline,
		)
		if err != nil {
			return err
		}
	}

	// next row
	if t.TahunNext != 0 {
		targetID := helper.GenerateID("TRGT-RB")

		_, err := tx.ExecContext(ctx, query,
			targetID,
			indikatorID,
			t.TahunNext,
			t.TargetNext,
			nil, // realisasi next → NULL
			t.SatuanNext,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *DataMasterRepositoryImpl) UpdateRB(ctx context.Context, tx *sql.Tx, req datamaster.MasterRB, rbId int) error {
	query := `
		UPDATE datamaster_rb
		SET jenis_rb = ?, kegiatan_utama = ?, keterangan = ?,
		    tahun_baseline = ?, tahun_next = ?, last_updated_by = ?,
		    current_version = current_version + 1
		WHERE id = ?
	`

	_, err := tx.ExecContext(ctx, query,
		req.JenisRB,
		req.KegiatanUtama,
		req.Keterangan,
		req.TahunBaseline,
		req.TahunNext,
		req.LastUpdatedBy,
		rbId,
	)

	return err
}

func (r *DataMasterRepositoryImpl) DeleteAllIndikatorAndTargetByRB(ctx context.Context, tx *sql.Tx, rbId int) error {

	// Delete target by indikator_id (cascade manually)
	queryTarget := `
		DELETE t FROM tb_target t
		JOIN tb_indikator i ON t.indikator_id = i.id
		WHERE i.id_rb = ?
	`

	if _, err := tx.ExecContext(ctx, queryTarget, rbId); err != nil {
		return err
	}

	// Delete indikator
	queryIndikator := `DELETE FROM tb_indikator WHERE id_rb = ?`

	_, err := tx.ExecContext(ctx, queryIndikator, rbId)
	return err
}

func (r *DataMasterRepositoryImpl) FindRBById(ctx context.Context, tx *sql.Tx, rbId int) (datamaster.MasterRB, error) {
	query := `
		SELECT id, jenis_rb, kegiatan_utama, keterangan, tahun_baseline, tahun_next
		FROM datamaster_rb
		WHERE id = ? AND is_active = 1
	`

	row := tx.QueryRowContext(ctx, query, rbId)

	var rb datamaster.MasterRB

	err := row.Scan(
		&rb.Id,
		&rb.JenisRB,
		&rb.KegiatanUtama,
		&rb.Keterangan,
		&rb.TahunBaseline,
		&rb.TahunNext,
	)
	if err != nil {
		return datamaster.MasterRB{}, err
	}

	return rb, nil
}

func (repo *DataMasterRepositoryImpl) DeleteRB(ctx context.Context, tx *sql.Tx, rbId int) error {
	// delete indikator target
	err := repo.DeleteAllIndikatorAndTargetByRB(ctx, tx, rbId)
	if err != nil {
		return err
	}

	query := `DELETE FROM datamaster_rb WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, rbId)
	return err
}

func (repo *DataMasterRepositoryImpl) PokinByIdRBs(ctx context.Context, tx *sql.Tx, listIdRB []int) ([]datamaster.PokinIdRBTagging, error) {
	if len(listIdRB) == 0 {
		return []datamaster.PokinIdRBTagging{}, errors.New("ids tidak boleh kosong")
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(listIdRB))
	args := make([]any, len(listIdRB))
	for i, id := range listIdRB {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`SELECT ket.id_tagging, ket.kode_program_unggulan, rb.kegiatan_utama, tag.id_pokin, tag.nama_tagging, pokin.nama_pohon, pokin.kode_opd, pokin.jenis_pohon
		FROM tb_keterangan_tagging_program_unggulan ket
		JOIN tb_tagging_pokin tag ON ket.id_tagging = tag.id AND tag.nama_tagging = 'RB'
		JOIN datamaster_rb rb ON ket.kode_program_unggulan = rb.id
		JOIN tb_pohon_kinerja pokin ON tag.id_pokin = pokin.id
		WHERE ket.kode_program_unggulan IN (%s)`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []datamaster.PokinIdRBTagging{}, err
	}
	defer rows.Close()

	var result []datamaster.PokinIdRBTagging
	for rows.Next() {
		var pokinRB datamaster.PokinIdRBTagging
		err := rows.Scan(
			&pokinRB.IdTagging,
			&pokinRB.KodeRB,
			&pokinRB.KegiatanUtama,
			&pokinRB.IdPokin,
			&pokinRB.NamaTagging,
			&pokinRB.NamaPohon,
			&pokinRB.KodeOpd,
			&pokinRB.JenisPohon,
		)
		if err != nil {
			return []datamaster.PokinIdRBTagging{}, err
		}
		result = append(result, pokinRB)
	}

	if err = rows.Err(); err != nil {
		return []datamaster.PokinIdRBTagging{}, err
	}
	return result, nil
}
