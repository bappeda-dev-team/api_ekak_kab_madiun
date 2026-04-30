package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CrosscuttingOpdRepositoryImpl struct {
}

func NewCrosscuttingOpdRepositoryImpl() *CrosscuttingOpdRepositoryImpl {
	return &CrosscuttingOpdRepositoryImpl{}
}

func (repository *CrosscuttingOpdRepositoryImpl) CreateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error) {
	// Hanya insert ke tb_crosscutting
	scriptCrosscutting := `
        INSERT INTO tb_crosscutting (
            crosscutting_from, 
            crosscutting_to, 
            kode_opd, 
            keterangan_crosscutting, 
			tahun,
            status
        ) VALUES (?, ?, ?, ?, ?, ?)
    `
	result, err := tx.ExecContext(ctx, scriptCrosscutting,
		parentId,
		0,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status)
	if err != nil {
		return pokin, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return pokin, err
	}
	pokin.Id = int(newId)

	return pokin, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindAllCrosscutting(ctx context.Context, tx *sql.Tx, parentId int) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            c.id as id_crosscutting,
            CASE 
                WHEN c.crosscutting_to = 0 OR p.id IS NULL THEN c.id
                ELSE p.id 
            END as id,
            COALESCE(p.nama_pohon, '') as nama_pohon,
            COALESCE(p.parent, 0) as parent,
            COALESCE(p.jenis_pohon, '') as jenis_pohon,
            COALESCE(CAST(p.level_pohon AS SIGNED), 0) as level_pohon,
            c.kode_opd,
            c.keterangan_crosscutting as keterangan,
            c.tahun,
            c.status,
            COALESCE(p.pegawai_action, NULL) as pegawai_action,
            COALESCE(p.created_at, NOW()) as created_at
        FROM tb_crosscutting c
        LEFT JOIN tb_pohon_kinerja p ON p.id = c.crosscutting_to 
        WHERE c.crosscutting_from = ?
    `
	rows, err := tx.QueryContext(ctx, script, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PohonKinerja
	for rows.Next() {
		pokin := domain.PohonKinerja{}
		var pegawaiActionJSON sql.NullString
		err := rows.Scan(
			&pokin.IdCrosscutting,
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&pegawaiActionJSON,
			&pokin.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if pegawaiActionJSON.Valid {
			var pegawaiAction interface{}
			err = json.Unmarshal([]byte(pegawaiActionJSON.String), &pegawaiAction)
			if err != nil {
				return nil, err
			}
			pokin.PegawaiAction = pegawaiAction
		}

		result = append(result, pokin)
	}
	return result, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinIds []int) ([]domain.Indikator, error) {
	// Cek jika array kosong
	if len(pokinIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(pokinIds))
	for i := range pokinIds {
		placeholders[i] = "?"
	}

	query := "SELECT id, pokin_id, indikator FROM tb_indikator WHERE pokin_id IN (" + strings.Join(placeholders, ",") + ")"

	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		args[i] = id
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		var pokinId int
		err := rows.Scan(
			&indikator.Id,
			&pokinId,
			&indikator.Indikator,
		)
		if err != nil {
			return nil, err
		}
		indikator.PokinId = strconv.Itoa(pokinId)
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindTargetByIndikatorIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) ([]domain.Target, error) {
	// Cek jika array kosong
	if len(indikatorIds) == 0 {
		return []domain.Target{}, nil
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(indikatorIds))
	for i := range indikatorIds {
		placeholders[i] = "?"
	}

	query := "SELECT id, indikator_id, target, satuan FROM tb_target WHERE indikator_id IN (" + strings.Join(placeholders, ",") + ")"

	args := make([]interface{}, len(indikatorIds))
	for i, id := range indikatorIds {
		args[i] = id
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) UpdateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Cek status dan crosscutting_to dari tb_crosscutting
	var currentStatus string
	var crosscuttingTo int
	err := tx.QueryRowContext(ctx, `
        SELECT status, crosscutting_to 
        FROM tb_crosscutting 
        WHERE id = ?`, pokin.Id).Scan(&currentStatus, &crosscuttingTo)
	if err != nil {
		return pokin, err
	}

	// Update berdasarkan status
	if currentStatus == "crosscutting_menunggu" || currentStatus == "crosscutting_ditolak" {
		// Update kode_opd dan keterangan di tb_crosscutting
		scriptCross := `
            UPDATE tb_crosscutting 
            SET kode_opd = ?,
                keterangan_crosscutting = ?
            WHERE id = ?
        `
		_, err = tx.ExecContext(ctx, scriptCross,
			pokin.KodeOpd,
			pokin.Keterangan,
			pokin.Id)
		if err != nil {
			return pokin, err
		}
	} else if currentStatus == "crosscutting_disetujui" || currentStatus == "crosscutting_disetujui_existing" {
		// Update hanya keterangan di tb_crosscutting
		scriptCross := `
            UPDATE tb_crosscutting 
            SET keterangan_crosscutting = ?
            WHERE id = ?
        `
		_, err = tx.ExecContext(ctx, scriptCross,
			pokin.Keterangan,
			pokin.Id)
		if err != nil {
			return pokin, err
		}

		// Update keterangan di tb_pohon_kinerja jika ada crosscutting_to
		if crosscuttingTo > 0 {
			scriptPokin := `
                UPDATE tb_pohon_kinerja 
                SET keterangan_crosscutting = ?
                WHERE id = ?
            `
			_, err = tx.ExecContext(ctx, scriptPokin,
				pokin.Keterangan,
				crosscuttingTo)
			if err != nil {
				return pokin, err
			}
		}
	}

	pokin.Status = currentStatus
	return pokin, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) ValidateKodeOpdChange(ctx context.Context, tx *sql.Tx, id int) error {
	var status string
	err := tx.QueryRowContext(ctx, "SELECT status FROM tb_crosscutting WHERE crosscutting_to = ?", id).Scan(&status)
	if err != nil {
		return err
	}

	if status != "crosscutting_menunggu" {
		return errors.New("kode OPD hanya dapat diubah saat status crosscutting_menunggu")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) DeleteCrosscutting(ctx context.Context, tx *sql.Tx, pokinId int, nipPegawai string) error {
	// Validasi status
	var currentStatus string
	query := `
        SELECT status FROM tb_pohon_kinerja 
        WHERE id = ?
    `
	err := tx.QueryRowContext(ctx, query, pokinId).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "crosscutting_disetujui" {
		return errors.New("crosscutting hanya dapat dihapus saat status crosscutting_disetujui")
	}

	// // Buat pegawai_action
	// currentTime := time.Now()
	// pegawaiAction := map[string]interface{}{
	// 	"reject_by": nipPegawai,
	// 	"reject_at": currentTime,
	// }

	// pegawaiActionJSON, err := json.Marshal(pegawaiAction)
	// if err != nil {
	// 	return err
	// }

	// scriptUpdatePokin := `
	//     UPDATE tb_pohon_kinerja
	//     SET parent = 0,
	//         status = 'crosscutting_ditolak',
	//         pegawai_action = ?
	//     WHERE id = ?
	// `

	scriptUpdatePokin := `
	DELETE FROM tb_pohon_kinerja 
	WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, scriptUpdatePokin, pokinId)
	if err != nil {
		return err
	}

	// Update status di tb_crosscutting
	scriptUpdateCross := `
        UPDATE tb_crosscutting 
        SET status = 'crosscutting_ditolak'
        WHERE crosscutting_to = ?
    `
	result, err := tx.ExecContext(ctx, scriptUpdateCross, pokinId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated in tb_crosscutting, check crosscutting_to value")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) ApproveOrRejectCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int, request pohonkinerja.CrosscuttingApproveRequest) error {
	var currentStatus, keterangan, kodeOpd, tahun string
	var crosscuttingTo int
	err := tx.QueryRowContext(ctx, `
		SELECT status, keterangan_crosscutting, kode_opd, tahun, crosscutting_to
		FROM tb_crosscutting WHERE id = ?`, crosscuttingId).
		Scan(&currentStatus, &keterangan, &kodeOpd, &tahun, &crosscuttingTo)
	if err != nil {
		return fmt.Errorf("error getting crosscutting data: %w", err)
	}
	if currentStatus != "crosscutting_menunggu" && currentStatus != "crosscutting_ditolak" {
		return errors.New("crosscutting sudah disetujui")
	}
	currentTime := time.Now()
	var pegawaiAction map[string]interface{}
	if request.Approve {
		pegawaiAction = map[string]interface{}{"approve_by": request.NipPegawai, "approve_at": currentTime}
	} else {
		pegawaiAction = map[string]interface{}{"reject_by": request.NipPegawai, "reject_at": currentTime}
	}
	pegawaiActionJSON, err := json.Marshal(pegawaiAction)
	if err != nil {
		return fmt.Errorf("error marshaling pegawai action: %w", err)
	}
	if request.Approve {
		if request.CreateNew {
			// ── Logic 1: Buat pohon kinerja baru ──
			// Pohon kinerja BARU: status crosscutting_disetujui, TANPA keterangan_crosscutting
			// (keterangan ada di tb_crosscutting, bukan di pohon kinerja)
			scriptNewPokin := `
				INSERT INTO tb_pohon_kinerja (
					nama_pohon, parent, level_pohon, jenis_pohon,
					kode_opd, tahun, status, pegawai_action, keterangan
				) VALUES ('', ?, ?, ?, ?, ?, '', ?, '')
			`
			result, err := tx.ExecContext(ctx, scriptNewPokin,
				request.ParentId, request.LevelPohon, request.JenisPohon,
				kodeOpd, tahun, pegawaiActionJSON)
			if err != nil {
				return fmt.Errorf("error creating new pohon kinerja: %w", err)
			}
			newPokinId, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("error getting last insert id: %w", err)
			}
			// Update tb_crosscutting: status + crosscutting_to ke id pohon baru
			_, err = tx.ExecContext(ctx, `
				UPDATE tb_crosscutting
				SET status = 'crosscutting_disetujui', crosscutting_to = ?
				WHERE id = ?
			`, newPokinId, crosscuttingId)
			if err != nil {
				return fmt.Errorf("error updating crosscutting (create new): %w", err)
			}
		} else if request.UseExisting {
			// ── Logic 2: Gunakan pohon kinerja yang sudah ada ──
			// TIDAK ubah tb_pohon_kinerja sama sekali
			// Hanya update tb_crosscutting: status + crosscutting_to ke id yang sudah ada
			_, err = tx.ExecContext(ctx, `
				UPDATE tb_crosscutting
				SET status = 'crosscutting_disetujui_existing', crosscutting_to = ?
				WHERE id = ?
			`, request.ExistingId, crosscuttingId)
			if err != nil {
				return fmt.Errorf("error updating crosscutting (use existing): %w", err)
			}
		}
	} else {
		// ── Logic 3: Tolak / balikkan ke menunggu ──
		// if crosscuttingTo > 0 {
		// 	_, err = tx.ExecContext(ctx, `
		// 		UPDATE tb_pohon_kinerja
		// 		SET status = 'crosscutting_menunggu', pegawai_action = ?
		// 		WHERE id = ?
		// 	`, pegawaiActionJSON, crosscuttingTo)
		// 	if err != nil {
		// 		return fmt.Errorf("error reverting pohon kinerja status: %w", err)
		// 	}
		// }
		_, err = tx.ExecContext(ctx, `
			UPDATE tb_crosscutting SET status = 'crosscutting_menunggu' WHERE id = ?
		`, crosscuttingId)
		if err != nil {
			return fmt.Errorf("error reverting crosscutting status: %w", err)
		}
	}
	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) DeleteUnused(ctx context.Context, tx *sql.Tx, crosscuttingId int) error {
	// Cek apakah data dengan status yang sesuai ada
	checkQuery := `
        SELECT COUNT(id) 
        FROM tb_crosscutting
        WHERE id = ? 
        AND status IN ('crosscutting_menunggu', 'crosscutting_ditolak')
    `
	var count int
	err := tx.QueryRowContext(ctx, checkQuery, crosscuttingId).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("crosscutting tidak dapat dihapus karena status tidak sesuai atau data tidak ditemukan")
	}

	// Hapus data di tb_crosscutting
	deleteQuery := `
        DELETE FROM tb_crosscutting 
        WHERE id = ? 
        AND status IN ('crosscutting_menunggu', 'crosscutting_ditolak')
    `
	result, err := tx.ExecContext(ctx, deleteQuery, crosscuttingId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("gagal menghapus data crosscutting")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindPokinByCrosscuttingStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.Crosscutting, error) {
	script := `
        SELECT 
            c.id, 
            c.keterangan_crosscutting, 
            c.kode_opd, 
            c.tahun,
            c.status,
            COALESCE(p.kode_opd, '') as opd_pengirim
        FROM tb_crosscutting c
        LEFT JOIN tb_pohon_kinerja p ON c.crosscutting_from = p.id
        WHERE c.kode_opd = ? 
        AND c.tahun = ? 
        AND c.status IN ('crosscutting_menunggu', 'crosscutting_ditolak')
    `
	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crosscuttings []domain.Crosscutting
	for rows.Next() {
		var crosscutting domain.Crosscutting
		err := rows.Scan(
			&crosscutting.Id,
			&crosscutting.Keterangan,
			&crosscutting.KodeOpd,
			&crosscutting.Tahun,
			&crosscutting.Status,
			&crosscutting.OpdPengirim,
		)
		if err != nil {
			return nil, err
		}
		crosscuttings = append(crosscuttings, crosscutting)
	}
	return crosscuttings, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindOPDCrosscuttingFrom(ctx context.Context, tx *sql.Tx, crosscuttingTo int) (string, error) {
	script := `
        SELECT 
            CASE 
                WHEN c.crosscutting_to = 0 THEN ''
                ELSE p.kode_opd 
            END as kode_opd
        FROM tb_crosscutting c
        LEFT JOIN tb_pohon_kinerja p ON c.crosscutting_from = p.id
        WHERE c.crosscutting_to = ?
    `
	var kodeOpd string
	err := tx.QueryRowContext(ctx, script, crosscuttingTo).Scan(&kodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("crosscutting tidak ditemukan")
		}
		return "", err
	}
	return kodeOpd, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindCrosscuttingByPokinIdsBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.Crosscutting, error) {
	if len(pokinIds) == 0 {
		return map[int][]domain.Crosscutting{}, nil
	}
	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `
		SELECT 
			c.id,
			c.crosscutting_to,
			COALESCE(c.keterangan_crosscutting, '') AS keterangan_crosscutting,
			COALESCE(pk_from.kode_opd, '') AS kode_opd_asal,
			c.status
		FROM tb_crosscutting c
		LEFT JOIN tb_pohon_kinerja pk_from ON c.crosscutting_from = pk_from.id
		WHERE c.crosscutting_to IN (` + strings.Join(placeholders, ",") + `)
		AND c.status IN ('crosscutting_disetujui', 'crosscutting_disetujui_existing')
		ORDER BY c.crosscutting_to, c.id
	`
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindCrosscuttingByPokinIdsBatch: %w", err)
	}
	defer rows.Close()
	result := make(map[int][]domain.Crosscutting)
	for rows.Next() {
		var c domain.Crosscutting
		if err := rows.Scan(&c.Id, &c.CrosscuttingTo, &c.Keterangan, &c.OpdPengirim, &c.Status); err != nil {
			return nil, fmt.Errorf("scan crosscutting batch: %w", err)
		}
		result[c.CrosscuttingTo] = append(result[c.CrosscuttingTo], c)
	}
	return result, rows.Err()
}

func (repository *CrosscuttingOpdRepositoryImpl) FixPokinStatusAfterExistingUnlink(
	ctx context.Context,
	tx *sql.Tx,
	pokinId int,
) error {
	var jenisPohon string
	err := tx.QueryRowContext(ctx,
		`SELECT jenis_pohon FROM tb_pohon_kinerja WHERE id = ?`, pokinId,
	).Scan(&jenisPohon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // pohon tidak ditemukan, skip
		}
		return fmt.Errorf("FixPokinStatus: gagal ambil jenis_pohon id=%d: %w", pokinId, err)
	}
	var newStatus string
	switch jenisPohon {
	case "Strategic Pemda", "Tactical Pemda", "Operational Pemda":
		newStatus = "pokin dari pemda"
	default:
		newStatus = ""
	}
	_, err = tx.ExecContext(ctx,
		`UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?`, newStatus, pokinId,
	)
	if err != nil {
		return fmt.Errorf("FixPokinStatus: gagal update status id=%d: %w", pokinId, err)
	}
	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FixPokinStatusAfterExistingDelete(
	ctx context.Context,
	tx *sql.Tx,
	pokinId int,
) error {
	var jenisPohon string
	err := tx.QueryRowContext(ctx, `
		SELECT jenis_pohon FROM tb_pohon_kinerja WHERE id = ?`, pokinId).Scan(&jenisPohon)
	if err != nil {
		// Pohon tidak ditemukan → tidak perlu diproses
		return nil
	}
	var newStatus string
	switch jenisPohon {
	case "Strategic Pemda", "Tactical Pemda", "Operational Pemda":
		newStatus = "pokin dari pemda"
	default:
		newStatus = "" // String kosong untuk jenis tanpa "Pemda"
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?`, newStatus, pokinId); err != nil {
		return fmt.Errorf("gagal update status pokin id=%d: %w", pokinId, err)
	}
	return nil
}
