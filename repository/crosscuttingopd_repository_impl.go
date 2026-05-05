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

func (repository *CrosscuttingOpdRepositoryImpl) DeleteCrosscutting(
	ctx context.Context, tx *sql.Tx, crosscuttingId int, nipPegawai string,
) error {
	var ccFrom, ccTo int
	var ccStatus string
	err := tx.QueryRowContext(ctx,
		`SELECT COALESCE(crosscutting_from, 0), COALESCE(crosscutting_to, 0), status
		 FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
	).Scan(&ccFrom, &ccTo, &ccStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("crosscutting tidak ditemukan")
		}
		return fmt.Errorf("gagal ambil data crosscutting id=%d: %w", crosscuttingId, err)
	}
	_ = ccFrom
	// Status belum disetujui → hapus baris saja
	if ccStatus == "crosscutting_menunggu" || ccStatus == "crosscutting_ditolak" {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
		); err != nil {
			return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	if ccStatus != "crosscutting_disetujui" && ccStatus != "crosscutting_disetujui_existing" {
		return fmt.Errorf("crosscutting id=%d tidak dapat dihapus, status: %s", crosscuttingId, ccStatus)
	}
	// Tidak ada pohon tujuan → hapus baris saja
	if ccTo == 0 {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
		); err != nil {
			return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	// ── Cabang A: crosscutting_disetujui_existing ─────────────────────────────
	if ccStatus == "crosscutting_disetujui_existing" {
		var pokinStatus string
		errPokin := tx.QueryRowContext(ctx,
			`SELECT COALESCE(status, '') FROM tb_pohon_kinerja WHERE id = ?`, ccTo,
		).Scan(&pokinStatus)
		if errPokin == sql.ErrNoRows {
			// Pohon tidak ada → hapus baris kita
			_, err = tx.ExecContext(ctx, `DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId)
			return err
		}
		if errPokin != nil {
			return fmt.Errorf("gagal cek status pohon id=%d: %w", ccTo, errPokin)
		}
		switch pokinStatus {
		case "crosscutting_disetujui_existing":
			// Pohon "existing" → hapus baris kita saja, STOP, jangan cascade
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
			); err != nil {
				return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
			}
			return nil
		case "crosscutting_disetujui":
			// Pohon lahir dari crosscutting → cek ref lain
			var countOther int
			if err := tx.QueryRowContext(ctx,
				`SELECT COUNT(*) FROM tb_crosscutting
				 WHERE crosscutting_to = ? AND id != ?`, ccTo, crosscuttingId,
			).Scan(&countOther); err != nil {
				return fmt.Errorf("gagal hitung ref lain crosscutting_to=%d: %w", ccTo, err)
			}
			if countOther > 0 {
				// Masih ada ref lain → hapus baris kita saja
				if _, err := tx.ExecContext(ctx,
					`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
				); err != nil {
					return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
				}
				return nil
			}
			// Tidak ada ref lain → pohon orphan, hapus pohon + child
			return repository.cleanOutgoingAndDeletePokin(ctx, tx, crosscuttingId, ccTo)
		default:
			// Pohon status lain → hapus baris kita saja
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
			); err != nil {
				return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
			}
			return nil
		}
	}
	// ── Cabang B: crosscutting_disetujui ─────────────────────────────────────
	var countRef int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_crosscutting WHERE crosscutting_to = ?`, ccTo,
	).Scan(&countRef); err != nil {
		return fmt.Errorf("gagal hitung referensi crosscutting_to=%d: %w", ccTo, err)
	}
	if countRef > 1 {
		// Pohon masih dipakai OPD lain → hapus baris kita saja
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
		); err != nil {
			return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	// Hanya 1 ref → hapus pohon + child
	return repository.cleanOutgoingAndDeletePokin(ctx, tx, crosscuttingId, ccTo)
}

// ─────────────────────────────────────────────────────────────────────────────
// cleanOutgoingAndDeletePokin
// Bersihkan outgoing crosscuttings dari pohon tujuan, hapus pohon+child,
// lalu reset baris crosscutting utama.
// ─────────────────────────────────────────────────────────────────────────────
func (repository *CrosscuttingOpdRepositoryImpl) cleanOutgoingAndDeletePokin(
	ctx context.Context, tx *sql.Tx, crosscuttingId int, pokinRootId int,
) error {
	// Hapus baris crosscutting yang keluar dari pohon tujuan (rantai outgoing)
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_crosscutting WHERE crosscutting_from = ? AND id != ?`,
		pokinRootId, crosscuttingId,
	); err != nil {
		return fmt.Errorf("gagal hapus outgoing crosscutting pohon id=%d: %w", pokinRootId, err)
	}
	// Kumpulkan subtree
	nodeIds, err := repository.collectSubtreeIdsForCrosscutting(ctx, tx, pokinRootId)
	if err != nil {
		return err
	}
	// Reset semua incoming crosscutting ke node subtree (selain baris utama kita)
	for _, nodeId := range nodeIds {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tb_crosscutting
			SET crosscutting_to = 0, status = 'crosscutting_menunggu'
			WHERE crosscutting_to = ? AND id != ?
		`, nodeId, crosscuttingId); err != nil {
			return fmt.Errorf("gagal reset incoming crosscutting node=%d: %w", nodeId, err)
		}
	}
	// Hapus pohon + dependensi
	for _, nodeId := range nodeIds {
		if err := repository.deletePokinDependenciesOnly(ctx, tx, nodeId); err != nil {
			return err
		}
	}
	// Reset baris crosscutting utama
	if _, err := tx.ExecContext(ctx, `
		UPDATE tb_crosscutting
		SET crosscutting_to = 0, status = 'crosscutting_menunggu'
		WHERE id = ?
	`, crosscuttingId); err != nil {
		return fmt.Errorf("gagal reset crosscutting id=%d: %w", crosscuttingId, err)
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
				) VALUES ('', ?, ?, ?, ?, ?, 'crosscutting_disetujui', ?, '')
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
			UPDATE tb_crosscutting SET status = 'crosscutting_ditolak' WHERE id = ?
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

// DELETE CROSSCUTTING DITERIMA
func (repository *CrosscuttingOpdRepositoryImpl) DeleteCrosscuttingDiterima(
	ctx context.Context, tx *sql.Tx, crosscuttingId int,
) error {
	// 1. Ambil data crosscutting
	var crosscuttingTo int
	var ccStatus string
	err := tx.QueryRowContext(ctx,
		`SELECT COALESCE(crosscutting_to, 0), status FROM tb_crosscutting WHERE id = ?`,
		crosscuttingId,
	).Scan(&crosscuttingTo, &ccStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("crosscutting tidak ditemukan")
		}
		return fmt.Errorf("gagal ambil data crosscutting id=%d: %w", crosscuttingId, err)
	}
	// Hanya boleh diproses jika status disetujui (salah satu dari dua jenis)
	if ccStatus != "crosscutting_disetujui" && ccStatus != "crosscutting_disetujui_existing" {
		return fmt.Errorf("crosscutting id=%d tidak dalam status disetujui (saat ini: %s)", crosscuttingId, ccStatus)
	}
	// crosscutting_to = 0 → tidak ada pohon tujuan, hapus row saja
	if crosscuttingTo == 0 {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
		); err != nil {
			return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	// 2. Hitung berapa banyak baris di tb_crosscutting yang mengarah ke crosscutting_to yang sama
	var countRef int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_crosscutting WHERE crosscutting_to = ?`, crosscuttingTo,
	).Scan(&countRef); err != nil {
		return fmt.Errorf("gagal hitung referensi crosscutting_to=%d: %w", crosscuttingTo, err)
	}
	// ── Cabang: status crosscutting_disetujui ─────────────────────────────────
	if ccStatus == "crosscutting_disetujui" {
		if countRef > 1 {
			// Pohon masih direferensi OPD lain → hapus row ini saja
			if _, err := tx.ExecContext(ctx,
				`UPDATE tb_crosscutting
SET crosscutting_to = 0, status = 'crosscutting_menunggu'
WHERE id = ?`, crosscuttingId,
			); err != nil {
				return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
			}
			return nil
		}
		// Hanya 1 referensi → hapus pohon kinerja + child
		return repository.deleteCrosscuttingPokinAndReset(ctx, tx, crosscuttingId, crosscuttingTo)
	}
	// ── Cabang: status crosscutting_disetujui_existing ────────────────────────
	if countRef > 1 {
		// Ada referensi lain ke pohon yang sama →
		// hanya lepas tautan baris ini (crosscutting_to=0, status=menunggu)
		if _, err := tx.ExecContext(ctx, `
			UPDATE tb_crosscutting
			SET crosscutting_to = 0, status = 'crosscutting_menunggu'
			WHERE id = ?
		`, crosscuttingId); err != nil {
			return fmt.Errorf("gagal reset crosscutting existing id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	// Hanya 1 referensi → cek status pohon kinerja
	var pokinStatus string
	err = tx.QueryRowContext(ctx,
		`SELECT COALESCE(status, '') FROM tb_pohon_kinerja WHERE id = ?`, crosscuttingTo,
	).Scan(&pokinStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			// Pohon sudah tidak ada, reset crosscutting saja
			if _, err := tx.ExecContext(ctx, `
				UPDATE tb_crosscutting
				SET crosscutting_to = 0, status = 'crosscutting_menunggu'
				WHERE id = ?
			`, crosscuttingId); err != nil {
				return fmt.Errorf("gagal reset crosscutting id=%d: %w", crosscuttingId, err)
			}
			return nil
		}
		return fmt.Errorf("gagal cek status pohon id=%d: %w", crosscuttingTo, err)
	}
	switch pokinStatus {
	case "crosscutting_disetujui":
		// Pohon lahir dari crosscutting → hapus pohon + child
		return repository.deleteCrosscuttingPokinAndReset(ctx, tx, crosscuttingId, crosscuttingTo)
	case "crosscutting_disetujui_existing":
		// Pohon existing yang di-link → jangan hapus pohon, hanya reset tautan
		if _, err := tx.ExecContext(ctx, `
			UPDATE tb_crosscutting
			SET crosscutting_to = 0, status = 'crosscutting_menunggu'
			WHERE id = ?
		`, crosscuttingId); err != nil {
			return fmt.Errorf("gagal reset crosscutting existing id=%d: %w", crosscuttingId, err)
		}
		return nil
	default:
		// Pohon sudah ada sebelumnya (status lain) → hanya lepas tautan
		if _, err := tx.ExecContext(ctx, `
			UPDATE tb_crosscutting
			SET crosscutting_to = 0, status = 'crosscutting_menunggu'
			WHERE id = ?
		`, crosscuttingId); err != nil {
			return fmt.Errorf("gagal reset crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Plan B: UnlinkCrosscuttingDiterima
//
// Parameter: crosscuttingId = ID di tb_crosscutting.
//
// Alur:
//   - Sama dengan Plan A untuk kasus > 1 referensi (hapus row crosscutting).
//   - Jika == 1 referensi: TIDAK hapus pohon kinerja.
//     Hanya reset tb_crosscutting: crosscutting_to = 0, status = 'crosscutting_menunggu'.
//     Pohon kinerja tetap ada, status pohon tidak diubah.
//
// ─────────────────────────────────────────────────────────────────────────────
func (repository *CrosscuttingOpdRepositoryImpl) UnlinkCrosscuttingDiterima(
	ctx context.Context, tx *sql.Tx, crosscuttingId int,
) error {
	// 1. Ambil data crosscutting
	var crosscuttingTo int
	var ccStatus string
	err := tx.QueryRowContext(ctx,
		`SELECT crosscutting_to, status FROM tb_crosscutting WHERE id = ?`,
		crosscuttingId,
	).Scan(&crosscuttingTo, &ccStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("crosscutting tidak ditemukan")
		}
		return fmt.Errorf("gagal ambil data crosscutting id=%d: %w", crosscuttingId, err)
	}
	if ccStatus != "crosscutting_disetujui" {
		return fmt.Errorf("crosscutting id=%d bukan status crosscutting_disetujui (saat ini: %s)", crosscuttingId, ccStatus)
	}
	if crosscuttingTo == 0 {
		// Tidak ada pohon tujuan, hapus row saja
		_, err = tx.ExecContext(ctx, `DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId)
		return err
	}
	// 2. Hitung referensi ke crosscutting_to yang sama
	var countRef int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tb_crosscutting WHERE crosscutting_to = ?`, crosscuttingTo,
	).Scan(&countRef); err != nil {
		return fmt.Errorf("gagal hitung referensi crosscutting_to=%d: %w", crosscuttingTo, err)
	}
	if countRef > 1 {
		// Lebih dari 1 referensi → hapus row ini saja, pohon tetap
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM tb_crosscutting WHERE id = ?`, crosscuttingId,
		); err != nil {
			return fmt.Errorf("gagal hapus crosscutting id=%d: %w", crosscuttingId, err)
		}
		return nil
	}
	// count == 1: HANYA lepas tautan, pohon kinerja tidak disentuh
	if _, err := tx.ExecContext(ctx, `
		UPDATE tb_crosscutting
		SET crosscutting_to = 0, status = 'crosscutting_menunggu'
		WHERE id = ?
	`, crosscuttingId); err != nil {
		return fmt.Errorf("gagal unlink crosscutting id=%d: %w", crosscuttingId, err)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Private helper: kumpulkan semua ID subtree dari satu pohon (root + children)
// Dipakai oleh DeleteCrosscuttingDiterima untuk tahu node mana saja yang
// perlu dihapus.
// ─────────────────────────────────────────────────────────────────────────────
func (repository *CrosscuttingOpdRepositoryImpl) collectSubtreeIdsForCrosscutting(
	ctx context.Context, tx *sql.Tx, rootId int,
) ([]int, error) {
	rows, err := tx.QueryContext(ctx, `
		WITH RECURSIVE subtree AS (
			SELECT id FROM tb_pohon_kinerja WHERE id = ?
			UNION ALL
			SELECT p.id FROM tb_pohon_kinerja p
			JOIN subtree s ON p.parent = s.id
		)
		SELECT id FROM subtree
	`, rootId)
	if err != nil {
		return nil, fmt.Errorf("gagal kumpulkan subtree dari id=%d: %w", rootId, err)
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("gagal scan subtree id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// Private helper: hapus semua dependensi + pohon kinerja untuk satu node ID
// (tanpa menyentuh tb_crosscutting — sudah dihandle sebelumnya)
// ─────────────────────────────────────────────────────────────────────────────
func (repository *CrosscuttingOpdRepositoryImpl) deletePokinDependenciesOnly(
	ctx context.Context, tx *sql.Tx, nodeId int,
) error {
	// 1. Hapus target
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM tb_target
		WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)
	`, nodeId); err != nil {
		return fmt.Errorf("gagal hapus target node=%d: %w", nodeId, err)
	}
	// 2. Hapus indikator
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_indikator WHERE pokin_id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus indikator node=%d: %w", nodeId, err)
	}
	// 3. Hapus pelaksana
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus pelaksana node=%d: %w", nodeId, err)
	}
	// 4. Hapus tagging
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_tagging_pokin WHERE id_pokin = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus tagging node=%d: %w", nodeId, err)
	}
	// 5. Hapus pohon kinerja
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_pohon_kinerja WHERE id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus pohon kinerja node=%d: %w", nodeId, err)
	}
	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) deleteCrosscuttingPokinAndReset(
	ctx context.Context, tx *sql.Tx, crosscuttingId int, pokinRootId int,
) error {
	// Kumpulkan seluruh node subtree
	nodeIds, err := repository.collectSubtreeIdsForCrosscutting(ctx, tx, pokinRootId)
	if err != nil {
		return err
	}
	// Reset semua crosscutting yang mengarah ke node subtree
	// (kecuali baris crosscuttingId sendiri — akan direset terpisah di akhir)
	for _, nodeId := range nodeIds {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tb_crosscutting
			SET crosscutting_to = 0, status = 'crosscutting_menunggu'
			WHERE crosscutting_to = ? AND id != ?
		`, nodeId, crosscuttingId); err != nil {
			return fmt.Errorf("gagal reset crosscutting anak node=%d: %w", nodeId, err)
		}
	}
	// Hapus dependensi + pohon kinerja per node (dari child ke atas)
	for _, nodeId := range nodeIds {
		if err := repository.deletePokinDependenciesOnly(ctx, tx, nodeId); err != nil {
			return err
		}
	}
	// Reset crosscutting utama: tautan sudah tidak valid karena pohon dihapus
	if _, err := tx.ExecContext(ctx, `
		UPDATE tb_crosscutting
		SET crosscutting_to = 0, status = 'crosscutting_menunggu'
		WHERE id = ?
	`, crosscuttingId); err != nil {
		return fmt.Errorf("gagal reset crosscutting id=%d: %w", crosscuttingId, err)
	}
	return nil
}
