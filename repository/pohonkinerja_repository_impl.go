package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PohonKinerjaRepositoryImpl struct {
}

func NewPohonKinerjaRepositoryImpl() *PohonKinerjaRepositoryImpl {
	return &PohonKinerjaRepositoryImpl{}
}
func (repository *PohonKinerjaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pohonKinerja.NamaPohon,
		pohonKinerja.Parent,
		pohonKinerja.JenisPohon,
		pohonKinerja.LevelPohon,
		pohonKinerja.KodeOpd,
		pohonKinerja.Keterangan,
		pohonKinerja.Tahun,
		pohonKinerja.Status)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}
	pohonKinerja.Id = int(lastInsertId)

	// Insert pelaksana
	scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
	for _, pelaksana := range pohonKinerja.Pelaksana {
		_, err = tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pohonKinerja.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pohonKinerja, err
		}
	}

	// Insert indikator
	for _, indikator := range pohonKinerja.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pohonKinerja.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pohonKinerja, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pohonKinerja, err
			}
		}
	}

	// Insert tagging
	scriptTagging := "INSERT INTO tb_tagging_pokin (id_pokin, nama_tagging) VALUES (?, ?)"
	for _, tagging := range pohonKinerja.TaggingPokin {
		// Insert tagging pokin
		result, err := tx.ExecContext(ctx, scriptTagging,
			pohonKinerja.Id,
			tagging.NamaTagging)
		if err != nil {
			return pohonKinerja, err
		}

		// Dapatkan ID tagging yang baru dibuat
		lastTaggingId, err := result.LastInsertId()
		if err != nil {
			return pohonKinerja, err
		}
		tagging.Id = int(lastTaggingId)
		tagging.IdPokin = pohonKinerja.Id

		// Insert keterangan tagging program untuk setiap tagging
		scriptKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			_, err := tx.ExecContext(ctx, scriptKeterangan,
				tagging.Id,
				keterangan.KodeProgramUnggulan,
				keterangan.Tahun) // Gunakan tahun yang sudah diset di service
			if err != nil {
				return pohonKinerja, err
			}
		}
	}

	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update tb_pohon_kinerja
	scriptPokin := `
        UPDATE tb_pohon_kinerja 
        SET nama_pohon = ?, 
            parent = CASE 
                WHEN clone_from = 0 THEN ? 
                ELSE parent 
            END,
            jenis_pohon = ?, 
            level_pohon = ?, 
            kode_opd = ?, 
            keterangan = ?, 
            tahun = ?,
			updated_by = ?
        WHERE id = ?`

	_, err := tx.ExecContext(ctx, scriptPokin,
		pohonKinerja.NamaPohon,
		pohonKinerja.Parent,
		pohonKinerja.JenisPohon,
		pohonKinerja.LevelPohon,
		pohonKinerja.KodeOpd,
		pohonKinerja.Keterangan,
		pohonKinerja.Tahun,
		pohonKinerja.UpdatedBy,
		pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	// Update pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pohonKinerja.Id))
	if err != nil {
		return pohonKinerja, err
	}

	for _, pelaksana := range pohonKinerja.Pelaksana {
		scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pohonKinerja.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pohonKinerja, err
		}
	}

	// Proses indikator
	for _, indikator := range pohonKinerja.Indikator {
		// Update atau insert indikator dengan clone_from
		scriptUpdateIndikator := `
			INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) 
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				indikator = VALUES(indikator),
				tahun = VALUES(tahun),
				clone_from = VALUES(clone_from)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			pohonKinerja.Id,
			indikator.Indikator,
			indikator.Tahun,
			indikator.CloneFrom)
		if err != nil {
			return pohonKinerja, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
		if err != nil {
			return pohonKinerja, err
		}

		// Insert target baru dengan clone_from
		for _, target := range indikator.Target {
			// Log untuk debugging
			fmt.Printf("Inserting target: ID=%s, IndikatorID=%s, CloneFrom=%s\n",
				target.Id, target.IndikatorId, target.CloneFrom)

			scriptInsertTarget := `
				INSERT INTO tb_target 
					(id, indikator_id, target, satuan, tahun, clone_from)
				VALUES 
					(?, ?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				target.IndikatorId,
				target.Target,
				target.Satuan,
				target.Tahun,
				target.CloneFrom) // Pastikan clone_from dimasukkan
			if err != nil {
				return pohonKinerja, fmt.Errorf("error inserting target: %v", err)
			}
		}
	}

	// Hapus indikator yang tidak ada dalam request
	var existingIndikatorIds []string
	scriptGetExisting := "SELECT id FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, scriptGetExisting, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return pohonKinerja, err
		}
		existingIndikatorIds = append(existingIndikatorIds, id)
	}

	// Buat map untuk indikator baru
	newIndikatorIds := make(map[string]bool)
	for _, ind := range pohonKinerja.Indikator {
		newIndikatorIds[ind.Id] = true
	}

	// Hapus indikator yang tidak ada dalam request
	for _, existingId := range existingIndikatorIds {
		if !newIndikatorIds[existingId] {
			// Hapus target terlebih dahulu
			scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteTargets, existingId)
			if err != nil {
				return pohonKinerja, err
			}

			// Kemudian hapus indikator
			scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteIndikator, existingId)
			if err != nil {
				return pohonKinerja, err
			}
		}
	}

	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	scriptPokin := `
        SELECT 
            pk.id, 
            COALESCE(pk.parent, 0) as parent, 
            COALESCE(pk.nama_pohon, '') as nama_pohon, 
            COALESCE(pk.jenis_pohon, '') as jenis_pohon, 
            COALESCE(pk.level_pohon, 0) as level_pohon, 
            COALESCE(pk.kode_opd, '') as kode_opd, 
            COALESCE(pk.keterangan, '') as keterangan, 
			COALESCE(pk.keterangan_crosscutting, '') as keterangan_crosscutting,
            COALESCE(pk.tahun, '') as tahun,
            COALESCE(pk.status, '') as status
        FROM 
            tb_pohon_kinerja pk 
        WHERE 
            pk.id = ?`

	rows, err := tx.QueryContext(ctx, scriptPokin, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pohonKinerja := domain.PohonKinerja{}
	if rows.Next() {
		err := rows.Scan(
			&pohonKinerja.Id,
			&pohonKinerja.Parent,
			&pohonKinerja.NamaPohon,
			&pohonKinerja.JenisPohon,
			&pohonKinerja.LevelPohon,
			&pohonKinerja.KodeOpd,
			&pohonKinerja.Keterangan,
			&pohonKinerja.KeteranganCrosscutting,
			&pohonKinerja.Tahun,
			&pohonKinerja.Status,
		)
		if err != nil {
			return domain.PohonKinerja{}, err
		}
	}

	return pohonKinerja, nil
}

//pokin lama
// func (repository *PohonKinerjaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
// 	script := `
//         SELECT
//             id,
//             COALESCE(nama_pohon, '') as nama_pohon,
//             COALESCE(parent, 0) as parent,
//             COALESCE(jenis_pohon, '') as jenis_pohon,
//             COALESCE(level_pohon, 0) as level_pohon,
//             COALESCE(kode_opd, '') as kode_opd,
//             COALESCE(keterangan, '') as keterangan,
//             COALESCE(keterangan_crosscutting, '') as keterangan_crosscutting,
//             COALESCE(tahun, '') as tahun,
//             COALESCE(status, '') as status,
// 			COALESCE(is_active) as is_active
//         FROM tb_pohon_kinerja
//         WHERE kode_opd = ?
// 		AND tahun = ?
// 		AND status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
//         ORDER BY level_pohon, id ASC`

// 	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var pokins []domain.PohonKinerja
// 	for rows.Next() {
// 		var pokin domain.PohonKinerja
// 		err := rows.Scan(
// 			&pokin.Id,
// 			&pokin.NamaPohon,
// 			&pokin.Parent,
// 			&pokin.JenisPohon,
// 			&pokin.LevelPohon,
// 			&pokin.KodeOpd,
// 			&pokin.Keterangan,
// 			&pokin.KeteranganCrosscutting,
// 			&pokin.Tahun,
// 			&pokin.Status,
// 			&pokin.IsActive,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		pokins = append(pokins, pokin)
// 	}

// 	// Inisialisasi slice kosong jika tidak ada data
// 	if pokins == nil {
// 		pokins = make([]domain.PohonKinerja, 0)
// 	}

// 	return pokins, nil
// }

// pokin baru
func (repository *PohonKinerjaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := `
    SELECT 
            id,
            COALESCE(nama_pohon, '') as nama_pohon,
            COALESCE(parent, 0) as parent,
            COALESCE(jenis_pohon, '') as jenis_pohon,
            COALESCE(level_pohon, 0) as level_pohon,
            COALESCE(kode_opd, '') as kode_opd,
            COALESCE(keterangan, '') as keterangan,
            COALESCE(keterangan_crosscutting, '') as keterangan_crosscutting,
            COALESCE(tahun, '') as tahun,
            COALESCE(status, '') as status,
            COALESCE(is_active, 0) as is_active,
            COALESCE(clone_from, 0) as clone_from
        FROM tb_pohon_kinerja 
        WHERE kode_opd = ? 
        AND tahun = ?
        AND level_pohon >= 4
        AND status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
        ORDER BY 
            level_pohon ASC, 
            id ASC
        LIMIT 10000`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Pre-allocate dengan estimasi capacity
	pokins := make([]domain.PohonKinerja, 0, 1000)
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.KeteranganCrosscutting,
			&pokin.Tahun,
			&pokin.Status,
			&pokin.IsActive,
			&pokin.CloneFrom,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}

	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindStrategicNoParent(ctx context.Context, tx *sql.Tx, levelPohon, parent int, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE level_pohon = ? AND parent = ? AND kode_opd = ? AND tahun = ?"
	rows, err := tx.QueryContext(ctx, script, levelPohon, parent, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.Parent, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Keterangan, &pokin.Tahun)
		if err != nil {
			return nil, err
		}
		result = append(result, pokin)
	}
	return result, nil
}

//POKIN LAMA
// func (repository *PohonKinerjaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
// 	// Temukan semua ID anak pohon secara rekursif
// 	scriptFindChildren := `
//         WITH RECURSIVE child_tree AS (
//             SELECT id, parent, clone_from, level_pohon
//             FROM tb_pohon_kinerja
//             WHERE id = ?

//             UNION ALL

//             SELECT t.id, t.parent, t.clone_from, t.level_pohon
//             FROM tb_pohon_kinerja t
//             JOIN child_tree ct ON t.parent = ct.id
//         )
//         SELECT id, clone_from FROM child_tree
//     `

// 	rows, err := tx.QueryContext(ctx, scriptFindChildren, id)
// 	if err != nil {
// 		return fmt.Errorf("gagal mencari turunan pohon: %v", err)
// 	}
// 	defer rows.Close()

// 	// Kumpulkan semua ID anak pohon dan clone_from
// 	var nodeIds []string
// 	cloneFromMap := make(map[string]int) // map[nodeId]cloneFromId
// 	for rows.Next() {
// 		var nodeId string
// 		var cloneFrom sql.NullInt64
// 		if err := rows.Scan(&nodeId, &cloneFrom); err != nil {
// 			return fmt.Errorf("gagal membaca ID turunan pohon: %v", err)
// 		}
// 		nodeIds = append(nodeIds, nodeId)
// 		if cloneFrom.Valid {
// 			cloneFromMap[nodeId] = int(cloneFrom.Int64)
// 		}
// 	}

// 	// Update status untuk semua node asli yang di-clone
// 	for _, cloneFromId := range cloneFromMap {
// 		scriptUpdateStatus := `
//             UPDATE tb_pohon_kinerja
//             SET status = 'menunggu_disetujui'
//             WHERE id = ?
//         `
// 		if _, err := tx.ExecContext(ctx, scriptUpdateStatus, cloneFromId); err != nil {
// 			return fmt.Errorf("gagal mengupdate status node asli ID %d: %v", cloneFromId, err)
// 		}
// 	}

// 	// Proses penghapusan untuk setiap node
// 	for _, nodeId := range nodeIds {
// 		// 1. Hapus target
// 		scriptDeleteTarget := `
//             DELETE FROM tb_target
//             WHERE indikator_id IN (
//                 SELECT id FROM tb_indikator
//                 WHERE pokin_id = ?
//             )`
// 		if _, err := tx.ExecContext(ctx, scriptDeleteTarget, nodeId); err != nil {
// 			return fmt.Errorf("gagal menghapus target: %v", err)
// 		}

// 		// 2. Hapus indikator
// 		scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
// 		if _, err := tx.ExecContext(ctx, scriptDeleteIndikator, nodeId); err != nil {
// 			return fmt.Errorf("gagal menghapus indikator: %v", err)
// 		}

// 		// 3. Hapus pelaksana
// 		scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
// 		if _, err := tx.ExecContext(ctx, scriptDeletePelaksana, nodeId); err != nil {
// 			return fmt.Errorf("gagal menghapus pelaksana: %v", err)
// 		}

// 		// 4. Tangani crosscutting
// 		// Hapus crosscutting yang menunggu/ditolak
// 		scriptDeletePendingCrosscutting := `
//             DELETE FROM tb_crosscutting
//             WHERE crosscutting_from = ?
//             AND status IN ('menunggu_disetujui', 'ditolak')`
// 		if _, err := tx.ExecContext(ctx, scriptDeletePendingCrosscutting, nodeId); err != nil {
// 			return fmt.Errorf("gagal menghapus crosscutting pending: %v", err)
// 		}

// 		// Update status crosscutting yang disetujui
// 		scriptUpdateCrosscuttingStatus := `
//             UPDATE tb_crosscutting
//             SET status = 'crosscutting_ditolak',
//                 crosscutting_to = 0
//             WHERE crosscutting_to = ?
//             AND status = 'crosscutting_disetujui'`
// 		if _, err := tx.ExecContext(ctx, scriptUpdateCrosscuttingStatus, nodeId); err != nil {
// 			return fmt.Errorf("gagal mengupdate status crosscutting: %v", err)
// 		}

// 		// 5. Hapus pohon kinerja
// 		scriptDeletePokin := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
// 		if _, err := tx.ExecContext(ctx, scriptDeletePokin, nodeId); err != nil {
// 			return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
// 		}

// 		// Hapus tagging sebelum menghapus pohon kinerja
// 		scriptDeleteTagging := fmt.Sprintf("DELETE FROM tb_tagging_pokin WHERE id_pokin IN (%s)", placeholders(len(nodeIds)))
// 		_, err = tx.ExecContext(ctx, scriptDeleteTagging, convertToInterface(nodeIds)...)
// 		if err != nil {
// 			return fmt.Errorf("gagal menghapus tagging: %v", err)
// 		}
// 	}

// 	return nil
// }
//END POKIN LAMA DELETE

// DELETE POKIN TRIAL
func (repository *PohonKinerjaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	// 1. Kumpulkan semua ID subtree (akar + seluruh turunan via parent)
	nodeIds, cloneFromMap, err := repository.collectSubtreeIds(ctx, tx, id)
	if err != nil {
		return err
	}
	// 2. Update status node asli yang di-clone → kembali ke 'menunggu_disetujui'
	for _, cloneFromId := range cloneFromMap {
		if _, err := tx.ExecContext(ctx,
			`UPDATE tb_pohon_kinerja SET status = 'menunggu_disetujui' WHERE id = ?`,
			cloneFromId,
		); err != nil {
			return fmt.Errorf("gagal update status clone_from id=%d: %w", cloneFromId, err)
		}
	}
	// 3. Handle crosscutting per node, kumpulkan pohon yg lahir dari crosscutting_disetujui
	//    dan hanya punya 1 referensi → perlu ikut dihapus
	var crosscuttingDerivedIds []int
	for _, nodeIdStr := range nodeIds {
		nodeId, _ := strconv.Atoi(nodeIdStr)
		derived, err := repository.processCrosscuttingForDelete(ctx, tx, nodeId)
		if err != nil {
			return fmt.Errorf("gagal proses crosscutting node id=%d: %w", nodeId, err)
		}
		crosscuttingDerivedIds = append(crosscuttingDerivedIds, derived...)
	}
	// 4. Hapus semua data pendukung + pohon kinerja utama (subtree)
	for _, nodeIdStr := range nodeIds {
		if err := repository.deletePokinAndDependencies(ctx, tx, nodeIdStr); err != nil {
			return err
		}
	}
	// 5. Rekursif hapus pohon yg lahir dari crosscutting_disetujui (single ref)
	//    Buat set nodeIds yang sudah dihapus supaya tidak double-process
	deletedSet := make(map[int]bool, len(nodeIds))
	for _, s := range nodeIds {
		if idInt, err := strconv.Atoi(s); err == nil {
			deletedSet[idInt] = true
		}
	}
	for _, derivedId := range crosscuttingDerivedIds {
		if deletedSet[derivedId] {
			continue
		}
		deletedSet[derivedId] = true
		if err := repository.Delete(ctx, tx, derivedId); err != nil {
			return fmt.Errorf("gagal hapus pohon crosscutting derived id=%d: %w", derivedId, err)
		}
	}
	return nil
}

//ENDING

// ─────────────────────────────────────────────────────────
// collectSubtreeIds — CTE rekursif ambil seluruh turunan
// ─────────────────────────────────────────────────────────
func (repository *PohonKinerjaRepositoryImpl) collectSubtreeIds(
	ctx context.Context, tx *sql.Tx, id int,
) (nodeIds []string, cloneFromMap map[string]int, err error) {
	rows, err := tx.QueryContext(ctx, `
		WITH RECURSIVE child_tree AS (
			SELECT id, clone_from
			FROM tb_pohon_kinerja
			WHERE id = ?
			UNION ALL
			SELECT t.id, t.clone_from
			FROM tb_pohon_kinerja t
			JOIN child_tree ct ON t.parent = ct.id
		)
		SELECT id, clone_from FROM child_tree
	`, id)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mencari turunan pohon: %w", err)
	}
	defer rows.Close()
	cloneFromMap = make(map[string]int)
	for rows.Next() {
		var nodeId string
		var cloneFrom sql.NullInt64
		if err := rows.Scan(&nodeId, &cloneFrom); err != nil {
			return nil, nil, fmt.Errorf("gagal scan turunan pohon: %w", err)
		}
		nodeIds = append(nodeIds, nodeId)
		if cloneFrom.Valid && cloneFrom.Int64 > 0 {
			cloneFromMap[nodeId] = int(cloneFrom.Int64)
		}
	}
	return nodeIds, cloneFromMap, rows.Err()
}

// ─────────────────────────────────────────────────────────
// processCrosscuttingForDelete — handle semua relasi crosscutting
// untuk satu nodeId sebelum pohon dihapus.
// Returns: slice crosscutting_to yang pohonnya perlu ikut dihapus.
// ─────────────────────────────────────────────────────────
func (repository *PohonKinerjaRepositoryImpl) processCrosscuttingForDelete(
	ctx context.Context, tx *sql.Tx, nodeId int,
) ([]int, error) {
	// ── A. Outgoing crosscutting: node ini sebagai crosscutting_from ──
	rows, err := tx.QueryContext(ctx, `
		SELECT id, COALESCE(crosscutting_to, 0), status
		FROM tb_crosscutting
		WHERE crosscutting_from = ?
	`, nodeId)
	if err != nil {
		return nil, fmt.Errorf("gagal query outgoing crosscutting nodeId=%d: %w", nodeId, err)
	}
	type ccRow struct {
		Id             int
		CrosscuttingTo int
		Status         string
	}
	var outgoing []ccRow
	for rows.Next() {
		var r ccRow
		if err := rows.Scan(&r.Id, &r.CrosscuttingTo, &r.Status); err != nil {
			rows.Close()
			return nil, fmt.Errorf("gagal scan crosscutting: %w", err)
		}
		outgoing = append(outgoing, r)
	}
	rows.Close()
	var toDeletePokinIds []int
	for _, cc := range outgoing {
		switch cc.Status {
		case "crosscutting_menunggu", "crosscutting_ditolak":
			// Hapus crosscutting row saja; pohon tujuan belum terbentuk
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM tb_crosscutting WHERE id = ?`, cc.Id,
			); err != nil {
				return nil, fmt.Errorf("gagal hapus crosscutting id=%d: %w", cc.Id, err)
			}
		case "crosscutting_disetujui":
			// Hitung berapa banyak crosscutting yang menuju pohon yang sama
			var countRef int
			if err := tx.QueryRowContext(ctx,
				`SELECT COUNT(*) FROM tb_crosscutting
				 WHERE crosscutting_to = ? AND status = 'crosscutting_disetujui'`,
				cc.CrosscuttingTo,
			).Scan(&countRef); err != nil {
				return nil, fmt.Errorf("gagal count ref crosscutting_to=%d: %w", cc.CrosscuttingTo, err)
			}
			// Hapus crosscutting row terlebih dahulu
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM tb_crosscutting WHERE id = ?`, cc.Id,
			); err != nil {
				return nil, fmt.Errorf("gagal hapus crosscutting id=%d: %w", cc.Id, err)
			}
			if cc.CrosscuttingTo > 0 && countRef == 1 {
				// Satu-satunya crosscutting ke pohon ini → pohon ikut dihapus
				toDeletePokinIds = append(toDeletePokinIds, cc.CrosscuttingTo)
			}
			// countRef > 1 → hanya crosscutting row dihapus; pohon tetap (masih dipakai lain)
		case "crosscutting_disetujui_existing":
			// Pohon existing TIDAK boleh dihapus; hanya putus tautan crosscutting
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM tb_crosscutting WHERE id = ?`, cc.Id,
			); err != nil {
				return nil, fmt.Errorf("gagal hapus crosscutting existing id=%d: %w", cc.Id, err)
			}
			// Fix status pohon existing (legacy repair — repo terpisah di crosscuttingopd_repository_impl)
			if cc.CrosscuttingTo > 0 {
				if err := repository.fixPokinStatusAfterExistingUnlink(ctx, tx, cc.CrosscuttingTo); err != nil {
					return nil, err
				}
			}
		}
	}
	// ── B. Incoming crosscutting: node ini sebagai crosscutting_to ──
	// (OPD lain yang crosscutting KE node ini)
	if _, err := tx.ExecContext(ctx, `
    UPDATE tb_crosscutting
    SET crosscutting_to = 0, status = 'crosscutting_menunggu'
    WHERE crosscutting_to = ?
	`, nodeId); err != nil {
		return nil, fmt.Errorf("gagal reset incoming crosscutting ke node=%d: %w", nodeId, err)
	}
	return toDeletePokinIds, nil
}

// ─────────────────────────────────────────────────────────
// deletePokinAndDependencies — hapus satu node + data pendukungnya
// (indikator, target, pelaksana, tagging, pohon kinerja)
// CATATAN: crosscutting sudah dihandle di processCrosscuttingForDelete
// ─────────────────────────────────────────────────────────
func (repository *PohonKinerjaRepositoryImpl) deletePokinAndDependencies(
	ctx context.Context, tx *sql.Tx, nodeId string,
) error {
	// 1. Hapus target
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM tb_target
		WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)
	`, nodeId); err != nil {
		return fmt.Errorf("gagal hapus target node=%s: %w", nodeId, err)
	}
	// 2. Hapus indikator
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_indikator WHERE pokin_id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus indikator node=%s: %w", nodeId, err)
	}
	// 3. Hapus pelaksana
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus pelaksana node=%s: %w", nodeId, err)
	}
	// 4. Hapus tagging
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_tagging_pokin WHERE id_pokin = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus tagging node=%s: %w", nodeId, err)
	}
	// 5. Hapus pohon kinerja
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM tb_pohon_kinerja WHERE id = ?`, nodeId,
	); err != nil {
		return fmt.Errorf("gagal hapus pohon kinerja node=%s: %w", nodeId, err)
	}
	return nil
}

// ─────────────────────────────────────────────────────────
// fixPokinStatusAfterExistingUnlink (private)
// Dipanggil saat crosscutting_disetujui_existing diputus;
// kembalikan status pohon existing sesuai jenis_pohon.
// (Logic ini juga tersedia di CrosscuttingOpdRepositoryImpl sebagai legacy repair)
// ─────────────────────────────────────────────────────────
func (repository *PohonKinerjaRepositoryImpl) fixPokinStatusAfterExistingUnlink(
	ctx context.Context, tx *sql.Tx, pokinId int,
) error {
	var jenisPohon string
	err := tx.QueryRowContext(ctx,
		`SELECT jenis_pohon FROM tb_pohon_kinerja WHERE id = ?`, pokinId,
	).Scan(&jenisPohon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // pohon tidak ada, skip
		}
		return fmt.Errorf("fixPokinStatus gagal ambil jenis_pohon id=%d: %w", pokinId, err)
	}
	var newStatus string
	switch jenisPohon {
	case "Strategic Pemda", "Tactical Pemda", "Operational Pemda":
		newStatus = "pokin dari pemda"
	default:
		newStatus = "" // string kosong untuk non-Pemda
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?`, newStatus, pokinId,
	); err != nil {
		return fmt.Errorf("fixPokinStatus gagal update id=%d: %w", pokinId, err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) recursiveDelete(ctx context.Context, tx *sql.Tx, id string) error {
	// Temukan semua ID anak pohon secara rekursif
	scriptFindRelatedIds := `
        WITH RECURSIVE related_ids AS (
            SELECT id FROM tb_pohon_kinerja WHERE id = ?

            UNION ALL

            SELECT t.id
            FROM tb_pohon_kinerja t
            JOIN related_ids r ON t.parent = r.id
        )
        SELECT id FROM related_ids
    `

	rows, err := tx.QueryContext(ctx, scriptFindRelatedIds, id)
	if err != nil {
		return fmt.Errorf("gagal mencari ID terkait: %v", err)
	}
	defer rows.Close()

	// Kumpulkan semua ID terkait
	var relatedIds []string
	for rows.Next() {
		var relatedId string
		if err := rows.Scan(&relatedId); err != nil {
			return fmt.Errorf("gagal membaca ID terkait: %v", err)
		}
		relatedIds = append(relatedIds, relatedId)
	}

	if len(relatedIds) == 0 {
		// Jika tidak ada ID terkait, tambahkan ID asli
		relatedIds = append(relatedIds, id)
	}

	// Hapus data yang terkait
	return repository.deleteRelatedData(ctx, tx, relatedIds)
}

func (repository *PohonKinerjaRepositoryImpl) deleteRelatedData(ctx context.Context, tx *sql.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil // Return early jika tidak ada data yang perlu dihapus
	}

	// Hapus crosscutting yang terkait dengan pohon kinerja yang akan dihapus
	scriptDeleteCrosscutting := `
        DELETE FROM tb_crosscutting 
        WHERE crosscutting_from IN (` + placeholders(len(ids)) + `)
        OR crosscutting_to IN (` + placeholders(len(ids)) + `)`
	if _, err := tx.ExecContext(ctx, scriptDeleteCrosscutting,
		append(convertToInterface(ids), convertToInterface(ids)...)...); err != nil {
		return fmt.Errorf("gagal menghapus crosscutting: %v", err)
	}

	// Hapus target
	scriptDeleteTarget := `
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id IN (` + placeholders(len(ids)) + `)
        )`
	if _, err := tx.ExecContext(ctx, scriptDeleteTarget, convertToInterface(ids)...); err != nil {
		return fmt.Errorf("gagal menghapus target: %v", err)
	}

	// Hapus indikator
	scriptDeleteIndikator := `
        DELETE FROM tb_indikator 
        WHERE pokin_id IN (` + placeholders(len(ids)) + `)`
	if _, err := tx.ExecContext(ctx, scriptDeleteIndikator, convertToInterface(ids)...); err != nil {
		return fmt.Errorf("gagal menghapus indikator: %v", err)
	}

	// Hapus pelaksana
	scriptDeletePelaksana := `
        DELETE FROM tb_pelaksana_pokin 
        WHERE pohon_kinerja_id IN (` + placeholders(len(ids)) + `)`
	if _, err := tx.ExecContext(ctx, scriptDeletePelaksana, convertToInterface(ids)...); err != nil {
		return fmt.Errorf("gagal menghapus pelaksana: %v", err)
	}

	// Hapus pohon kinerja22
	scriptDeletePokin := `
        DELETE FROM tb_pohon_kinerja 
        WHERE id IN (` + placeholders(len(ids)) + `)`
	if _, err := tx.ExecContext(ctx, scriptDeletePokin, convertToInterface(ids)...); err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

// Fungsi untuk membuat placeholder dinamis
func placeholders(n int) string {
	ph := make([]string, n)
	for i := range ph {
		ph[i] = "?"
	}
	return strings.Join(ph, ", ")
}

// Fungsi untuk mengonversi slice string ke slice interface{}
func convertToInterface(ids []string) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}

func (repository *PohonKinerjaRepositoryImpl) FindPelaksanaPokin(ctx context.Context, tx *sql.Tx, pohonKinerjaId string) ([]domain.PelaksanaPokin, error) {
	script := "SELECT id, pohon_kinerja_id, pegawai_id FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	rows, err := tx.QueryContext(ctx, script, pohonKinerjaId)
	helper.PanicIfError(err)
	defer rows.Close()

	var result []domain.PelaksanaPokin
	for rows.Next() {
		var pelaksana domain.PelaksanaPokin
		err := rows.Scan(&pelaksana.Id, &pelaksana.PohonKinerjaId, &pelaksana.PegawaiId)
		helper.PanicIfError(err)
		result = append(result, pelaksana)
	}
	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePelaksanaPokin(ctx context.Context, tx *sql.Tx, pelaksanaId string) error {
	script := "DELETE FROM tb_pelaksana_pokin WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pelaksanaId)
	return err
}

// admin pokin
func (repository *PohonKinerjaRepositoryImpl) CreatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Insert pohon kinerja tanpa ID
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon, pokinAdmin.Parent, pokinAdmin.JenisPohon, pokinAdmin.LevelPohon, pokinAdmin.KodeOpd, pokinAdmin.Keterangan, pokinAdmin.Tahun, pokinAdmin.Status)
	if err != nil {
		return pokinAdmin, err
	}

	// Dapatkan ID yang baru dibuat
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pokinAdmin, err
	}
	pokinAdmin.Id = int(lastInsertId)

	// // Tambahkan insert pelaksana
	// scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
	// for _, pelaksana := range pokinAdmin.Pelaksana {
	// 	_, err = tx.ExecContext(ctx, scriptPelaksana,
	// 		pelaksana.Id,
	// 		fmt.Sprint(pokinAdmin.Id),
	// 		pelaksana.PegawaiId)
	// 	if err != nil {
	// 		return pokinAdmin, err
	// 	}
	// }

	// Insert indikator
	for _, indikator := range pokinAdmin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id, pokinAdmin.Id, indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	// Insert tagging
	for _, tagging := range pokinAdmin.TaggingPokin {
		// Insert tagging
		scriptTagging := "INSERT INTO tb_tagging_pokin (id_pokin, nama_tagging) VALUES (?, ?)"
		resultTagging, err := tx.ExecContext(ctx, scriptTagging,
			pokinAdmin.Id,
			tagging.NamaTagging)
		if err != nil {
			return pokinAdmin, err
		}

		// Dapatkan ID tagging yang baru dibuat
		lastTaggingId, err := resultTagging.LastInsertId()
		if err != nil {
			return pokinAdmin, err
		}
		tagging.Id = int(lastTaggingId)
		tagging.IdPokin = pokinAdmin.Id

		// Insert keterangan program unggulan
		for _, keterangan := range tagging.KeteranganTaggingProgram {
			scriptKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
			_, err = tx.ExecContext(ctx, scriptKeterangan,
				lastTaggingId,
				keterangan.KodeProgramUnggulan,
				keterangan.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update tb_pohon_kinerja dengan mempertahankan status
	scriptPokin := `
        UPDATE tb_pohon_kinerja 
        SET nama_pohon = ?, 
            parent = CASE 
                WHEN clone_from = 0 THEN ? 
                ELSE parent 
            END,
            jenis_pohon = ?, 
            level_pohon = ?, 
            kode_opd = ?, 
            keterangan = ?, 
            tahun = ?,
			updated_by = ?
        WHERE id = ?`

	_, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon,
		pokinAdmin.Parent,
		pokinAdmin.JenisPohon,
		pokinAdmin.LevelPohon,
		pokinAdmin.KodeOpd,
		pokinAdmin.Keterangan,
		pokinAdmin.Tahun,
		pokinAdmin.UpdatedBy,
		pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Update pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pokinAdmin.Id))
	if err != nil {
		return pokinAdmin, err
	}

	for _, pelaksana := range pokinAdmin.Pelaksana {
		scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pokinAdmin.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pokinAdmin, err
		}
	}

	// Proses indikator
	for _, indikator := range pokinAdmin.Indikator {
		// Update atau insert indikator
		scriptUpdateIndikator := `
			INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) 
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				indikator = VALUES(indikator),
				tahun = VALUES(tahun)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			pokinAdmin.Id,
			indikator.Indikator,
			indikator.Tahun,
			indikator.CloneFrom)
		if err != nil {
			return pokinAdmin, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target baru
		for _, target := range indikator.Target {
			scriptInsertTarget := `
				INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, clone_from)
				VALUES (?, ?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun,
				target.CloneFrom)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	// Hapus indikator yang tidak ada dalam request
	var existingIndikatorIds []string
	scriptGetExisting := "SELECT id FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, scriptGetExisting, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return pokinAdmin, err
		}
		existingIndikatorIds = append(existingIndikatorIds, id)
	}

	// Buat map untuk indikator baru
	newIndikatorIds := make(map[string]bool)
	for _, ind := range pokinAdmin.Indikator {
		newIndikatorIds[ind.Id] = true
	}

	// Hapus indikator yang tidak ada dalam request
	for _, existingId := range existingIndikatorIds {
		if !newIndikatorIds[existingId] {
			// Hapus target terlebih dahulu
			scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteTargets, existingId)
			if err != nil {
				return pokinAdmin, err
			}

			// Kemudian hapus indikator
			scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteIndikator, existingId)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	// Update tagging
	// Hapus tagging yang tidak ada di request baru
	existingTaggings, err := repository.FindTaggingByPokinId(ctx, tx, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Buat map untuk tracking ID yang masih digunakan
	existingIds := make(map[int]bool)
	for _, tagging := range pokinAdmin.TaggingPokin {
		if tagging.Id != 0 {
			existingIds[tagging.Id] = true
		}
	}

	// Hapus tagging yang tidak ada di request baru
	for _, existing := range existingTaggings {
		if !existingIds[existing.Id] {
			// Hapus keterangan program unggulan terlebih dahulu
			scriptDeleteKeterangan := "DELETE FROM tb_keterangan_tagging_program_unggulan WHERE id_tagging = ?"
			_, err := tx.ExecContext(ctx, scriptDeleteKeterangan, existing.Id)
			if err != nil {
				return pokinAdmin, err
			}

			// Kemudian hapus tagging
			scriptDelete := "DELETE FROM tb_tagging_pokin WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDelete, existing.Id)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	// Update atau insert tagging
	for _, tagging := range pokinAdmin.TaggingPokin {
		if tagging.Id != 0 {
			// Update existing tagging
			script := "UPDATE tb_tagging_pokin SET nama_tagging = ? WHERE id = ? AND id_pokin = ?"
			_, err := tx.ExecContext(ctx, script,
				tagging.NamaTagging,
				tagging.Id,
				pokinAdmin.Id)
			if err != nil {
				return pokinAdmin, err
			}

			// Hapus keterangan lama
			scriptDeleteKeterangan := "DELETE FROM tb_keterangan_tagging_program_unggulan WHERE id_tagging = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteKeterangan, tagging.Id)
			if err != nil {
				return pokinAdmin, err
			}

			// Insert keterangan baru
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				scriptInsertKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
				_, err = tx.ExecContext(ctx, scriptInsertKeterangan,
					tagging.Id,
					keterangan.KodeProgramUnggulan,
					keterangan.Tahun)
				if err != nil {
					return pokinAdmin, err
				}
			}
		} else {
			// Insert new tagging
			scriptTagging := "INSERT INTO tb_tagging_pokin (id_pokin, nama_tagging) VALUES (?, ?)"
			result, err := tx.ExecContext(ctx, scriptTagging,
				pokinAdmin.Id,
				tagging.NamaTagging)
			if err != nil {
				return pokinAdmin, err
			}

			newId, err := result.LastInsertId()
			if err != nil {
				return pokinAdmin, err
			}
			tagging.Id = int(newId)

			// Insert keterangan program unggulan
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				scriptKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
				_, err = tx.ExecContext(ctx, scriptKeterangan,
					newId,
					keterangan.KodeProgramUnggulan,
					keterangan.Tahun)
				if err != nil {
					return pokinAdmin, err
				}
			}
		}
		tagging.IdPokin = pokinAdmin.Id
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePelaksanaOnly(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err := tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pokin.Id))
	if err != nil {
		return pokin, err
	}

	for _, pelaksana := range pokin.Pelaksana {
		scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pokin.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pokin, err
		}
	}

	return pokin, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePokinAdmin(ctx context.Context, tx *sql.Tx, id int) error {
	// Query untuk mendapatkan semua ID yang akan dihapus
	findIdsScript := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: node yang akan dihapus
            SELECT id, parent, level_pohon, clone_from 
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: child nodes dan data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON 
                -- Ambil child nodes langsung
                pk.parent = ph.id OR 
                -- Jika data asli, ambil yang mengclone-nya
                (ph.clone_from = 0 AND pk.clone_from = ph.id)
        ),
        clone_hierarki AS (
            -- Base case: data yang mengclone dan data yang parent-nya terhubung dengan id yang dihapus
            SELECT id, parent, level_pohon, clone_from
            FROM tb_pohon_kinerja
            WHERE clone_from IN (SELECT id FROM pohon_hierarki)
            OR parent IN (SELECT id FROM pohon_hierarki)
            
            UNION ALL
            
            -- Recursive case: child nodes dari data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN clone_hierarki ch ON 
                pk.parent = ch.id
        ),
        parent_hierarki AS (
            -- Ambil data yang parent-nya adalah id yang akan dihapus
            SELECT id, parent, level_pohon, clone_from
            FROM tb_pohon_kinerja
            WHERE parent = ?
        )
        SELECT id FROM pohon_hierarki
        UNION
        SELECT id FROM clone_hierarki
        UNION
        SELECT id FROM parent_hierarki;
    `

	rows, err := tx.QueryContext(ctx, findIdsScript, id, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil hierarki ID: %v", err)
	}
	defer rows.Close()

	// Kumpulkan semua ID yang akan dihapus
	var idsToDelete []interface{}
	for rows.Next() {
		var idToDelete int
		if err := rows.Scan(&idToDelete); err != nil {
			return fmt.Errorf("gagal scan ID: %v", err)
		}
		idsToDelete = append(idsToDelete, idToDelete)
	}

	if len(idsToDelete) == 0 {
		return fmt.Errorf("tidak ada data yang akan dihapus")
	}

	// Buat placeholder untuk query IN clause
	placeholders := make([]string, len(idsToDelete))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	// Hapus target terlebih dahulu
	scriptDeleteTarget := fmt.Sprintf(`
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id IN (%s)
        )`, inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus target: %v", err)
	}

	// Hapus indikator
	scriptDeleteIndikator := fmt.Sprintf("DELETE FROM tb_indikator WHERE pokin_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus indikator: %v", err)
	}

	// Hapus pelaksana
	scriptDeletePelaksana := fmt.Sprintf("DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pelaksana: %v", err)
	}

	// Hapus pohon kinerja
	scriptDeletePokin := fmt.Sprintf("DELETE FROM tb_pohon_kinerja WHERE id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePokin, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := `
        SELECT 
            pk.id, 
            pk.parent, 
            pk.nama_pohon, 
            pk.jenis_pohon, 
            pk.level_pohon, 
            pk.kode_opd, 
            pk.keterangan, 
            pk.tahun,
            pk.status,
			pk.is_active,
			pk.clone_from
        FROM 
            tb_pohon_kinerja pk 
        WHERE 
            pk.id = ?`

	var pokin domain.PohonKinerja
	var cloneFromNullInt sql.NullInt64
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&pokin.Id,
		&pokin.Parent,
		&pokin.NamaPohon,
		&pokin.JenisPohon,
		&pokin.LevelPohon,
		&pokin.KodeOpd,
		&pokin.Keterangan,
		&pokin.Tahun,
		&pokin.Status,
		&pokin.IsActive,
		&cloneFromNullInt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PohonKinerja{}, fmt.Errorf("pohon kinerja tidak ditemukan")
		}
		return domain.PohonKinerja{}, err
	}
	if cloneFromNullInt.Valid {
		pokin.CloneFrom = int(cloneFromNullInt.Int64)
	} else {

		pokin.CloneFrom = 0
	}

	// Ambil data tagging
	scriptTagging := `
    SELECT 
        t.id,
        t.id_pokin,
        t.nama_tagging,
        k.id as keterangan_id,
        k.kode_program_unggulan,
        k.tahun
    FROM tb_tagging_pokin t
    LEFT JOIN tb_keterangan_tagging_program_unggulan k ON t.id = k.id_tagging
    WHERE t.id_pokin = ?`

	taggingRows, err := tx.QueryContext(ctx, scriptTagging, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer taggingRows.Close()

	taggingMap := make(map[int]*domain.TaggingPokin)
	for taggingRows.Next() {
		var (
			taggingId, idPokin  int
			namaTagging         string
			keteranganId        sql.NullInt64
			kodeProgramUnggulan sql.NullString
			tahun               sql.NullString
		)

		err := taggingRows.Scan(
			&taggingId,
			&idPokin,
			&namaTagging,
			&keteranganId,
			&kodeProgramUnggulan,
			&tahun,
		)
		if err != nil {
			return domain.PohonKinerja{}, err
		}

		// Ambil atau buat tagging baru
		tagging, exists := taggingMap[taggingId]
		if !exists {
			tagging = &domain.TaggingPokin{
				Id:          taggingId,
				IdPokin:     idPokin,
				NamaTagging: namaTagging,
			}
			taggingMap[taggingId] = tagging
		}

		// Tambahkan keterangan jika ada
		if keteranganId.Valid && kodeProgramUnggulan.Valid {
			keterangan := domain.KeteranganTagging{
				Id:                  int(keteranganId.Int64),
				IdTagging:           taggingId,
				KodeProgramUnggulan: kodeProgramUnggulan.String,
				Tahun:               tahun.String,
			}
			tagging.KeteranganTaggingProgram = append(tagging.KeteranganTaggingProgram, keterangan)
		}
	}

	// Konversi map ke slice
	var taggings []domain.TaggingPokin
	for _, tagging := range taggingMap {
		taggings = append(taggings, *tagging)
	}
	pokin.TaggingPokin = taggings

	return pokin, nil
}
func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            i.id as indikator_id,
            i.indikator as nama_indikator,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan
        FROM 
            tb_pohon_kinerja pk
        LEFT JOIN 
            tb_indikator i ON pk.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        WHERE 
            pk.tahun = ?
        ORDER BY 
            pk.level_pohon, pk.id, i.id, t.id
    `

	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan pohon kinerja yang sudah diproses
	pokinMap := make(map[int]domain.PohonKinerja)
	indikatorMap := make(map[string]domain.Indikator)

	for rows.Next() {
		var (
			pokinId, parent, levelPohon                            int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin string
			indikatorId, namaIndikator                             sql.NullString
			targetId, targetValue, targetSatuan                    sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon, &levelPohon,
			&kodeOpd, &keterangan, &tahunPokin,
			&indikatorId, &namaIndikator,
			&targetId, &targetValue, &targetSatuan,
		)
		if err != nil {
			return nil, err
		}

		// Proses Pohon Kinerja
		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = domain.PohonKinerja{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahunPokin,
			}
			pokinMap[pokinId] = pokin
		}

		// Proses Indikator jika ada
		if indikatorId.Valid && namaIndikator.Valid {
			indikator, exists := indikatorMap[indikatorId.String]
			if !exists {
				indikator = domain.Indikator{
					Id:        indikatorId.String,
					PokinId:   fmt.Sprint(pokinId),
					Indikator: namaIndikator.String,
					Tahun:     tahunPokin,
				}
			}

			// Proses Target jika ada
			if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      targetValue.String,
					Satuan:      targetSatuan.String,
					Tahun:       tahunPokin,
				}
				indikator.Target = append(indikator.Target, target)
			}

			indikatorMap[indikatorId.String] = indikator

			// Update indikator di pokin
			pokin.Indikator = append(pokin.Indikator, indikator)
			pokinMap[pokinId] = pokin
		}
	}

	// Konversi map ke slice
	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		result = append(result, pokin)
	}

	// Urutkan berdasarkan level dan ID
	sort.Slice(result, func(i, j int) bool {
		if result[i].LevelPohon == result[j].LevelPohon {
			return result[i].Id < result[j].Id
		}
		return result[i].LevelPohon < result[j].LevelPohon
	})

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminByIdHierarki(ctx context.Context, tx *sql.Tx, idPokin int) ([]domain.PohonKinerja, error) {
	script := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: pilih node yang diminta
            SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, is_active
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: ambil semua child nodes
            SELECT pk.id, pk.nama_pohon, pk.parent, pk.jenis_pohon, pk.level_pohon, pk.kode_opd, pk.keterangan, pk.tahun, pk.status, pk.is_active
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON pk.parent = ph.id
        )
        SELECT 
            ph.id,
            ph.nama_pohon,
            ph.parent,
            ph.jenis_pohon,
            ph.level_pohon,
            ph.kode_opd,
            ph.keterangan,
            ph.tahun,
            ph.status,
            ph.is_active,
            i.id as indikator_id,
            i.indikator as nama_indikator,
            COALESCE(i.created_at, '') as indikator_created_at,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            pp.id as pelaksana_id,
            pp.pegawai_id
        FROM 
            pohon_hierarki ph
        LEFT JOIN 
            tb_indikator i ON ph.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        LEFT JOIN 
            tb_pelaksana_pokin pp ON ph.id = pp.pohon_kinerja_id
        ORDER BY 
            ph.level_pohon, ph.id, i.created_at, i.id, t.id, pp.id
    `

	rows, err := tx.QueryContext(ctx, script, idPokin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]domain.PohonKinerja)
	processedIndikators := make(map[string]bool)

	for rows.Next() {
		var (
			pokinId, parent, levelPohon                                    int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin, status string
			is_active                                                      bool
			indikatorId, namaIndikator                                     sql.NullString
			indikatorCreatedAt                                             sql.NullString // Ubah ke NullString
			targetId, targetValue, targetSatuan                            sql.NullString
			pelaksanaId, pegawaiId                                         sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon, &levelPohon,
			&kodeOpd, &keterangan, &tahunPokin, &status, &is_active,
			&indikatorId, &namaIndikator, &indikatorCreatedAt,
			&targetId, &targetValue, &targetSatuan,
			&pelaksanaId, &pegawaiId,
		)
		if err != nil {
			return nil, err
		}

		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = domain.PohonKinerja{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahunPokin,
				Status:     status,
				IsActive:   is_active,
			}
		}

		if pelaksanaId.Valid && pegawaiId.Valid {
			pelaksana := domain.PelaksanaPokin{
				Id:        pelaksanaId.String,
				PegawaiId: pegawaiId.String,
			}
			isDuplicate := false
			for _, p := range pokin.Pelaksana {
				if p.Id == pelaksana.Id {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				pokin.Pelaksana = append(pokin.Pelaksana, pelaksana)
			}
		}

		if indikatorId.Valid && namaIndikator.Valid {
			if !processedIndikators[indikatorId.String] {
				processedIndikators[indikatorId.String] = true

				// Parse created_at string ke time.Time jika ada
				var createdAt time.Time
				if indikatorCreatedAt.Valid && indikatorCreatedAt.String != "" {
					parsedTime, err := time.Parse("2006-01-02 15:04:05", indikatorCreatedAt.String)
					if err == nil {
						createdAt = parsedTime
					}
				}

				indikator := domain.Indikator{
					Id:        indikatorId.String,
					PokinId:   fmt.Sprint(pokinId),
					Indikator: namaIndikator.String,
					Tahun:     tahunPokin,
					CreatedAt: createdAt,
				}

				processedTargets := make(map[string]bool)

				if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
					if !processedTargets[targetId.String] {
						processedTargets[targetId.String] = true
						target := domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       tahunPokin,
						}
						indikator.Target = append(indikator.Target, target)
					}
				}

				pokin.Indikator = append(pokin.Indikator, indikator)
			}
		}

		pokinMap[pokinId] = pokin
	}

	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		// Urutkan indikator berdasarkan created_at sebelum menambahkan ke result
		sort.Slice(pokin.Indikator, func(i, j int) bool {
			return pokin.Indikator[i].CreatedAt.Before(pokin.Indikator[j].CreatedAt)
		})
		result = append(result, pokin)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].LevelPohon == result[j].LevelPohon {
			return result[i].Id < result[j].Id
		}
		return result[i].LevelPohon < result[j].LevelPohon
	})

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error) {
	script := `
        SELECT i.id, i.pokin_id, i.indikator, i.tahun, i.clone_from,
               t.id, t.indikator_id, t.target, t.satuan, t.tahun, t.clone_from
        FROM tb_indikator i
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE i.pokin_id = ?`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var indId, pokinId, indikator, indTahun, indCloneFrom string
		var targetId, indikatorId, target, satuan, targetTahun sql.NullString
		var targetCloneFrom sql.NullString

		err := rows.Scan(
			&indId, &pokinId, &indikator, &indTahun, &indCloneFrom,
			&targetId, &indikatorId, &target, &satuan, &targetTahun, &targetCloneFrom)
		if err != nil {
			return nil, err
		}

		// Proses Indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:        indId,
				Indikator: indikator,
				Tahun:     indTahun,
				CloneFrom: indCloneFrom,
				Target:    []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// Proses Target jika ada
		if targetId.Valid && indikatorId.Valid {
			target := domain.Target{
				Id:          targetId.String,
				IndikatorId: indikatorId.String,
				Target:      target.String,
				Satuan:      satuan.String,
				Tahun:       targetTahun.String,
			}
			if targetCloneFrom.Valid {
				target.CloneFrom = targetCloneFrom.String
			}
			ind.Target = append(ind.Target, target)
		}
	}

	// Convert map to slice
	var result []domain.Indikator
	for _, ind := range indikatorMap {
		result = append(result, *ind)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinToClone(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, is_active FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, fmt.Errorf("gagal memeriksa data yang akan di-clone: %v", err)
	}
	defer rows.Close()

	var existingPokin domain.PohonKinerja
	if rows.Next() {
		err := rows.Scan(
			&existingPokin.Id,
			&existingPokin.NamaPohon,
			&existingPokin.Parent,
			&existingPokin.JenisPohon,
			&existingPokin.LevelPohon,
			&existingPokin.KodeOpd,
			&existingPokin.Keterangan,
			&existingPokin.Tahun,
			&existingPokin.Status,
			&existingPokin.IsActive,
		)
		if err != nil {
			return domain.PohonKinerja{}, fmt.Errorf("gagal membaca data yang akan di-clone: %v", err)
		}
		return existingPokin, nil
	}
	return domain.PohonKinerja{}, fmt.Errorf("data dengan ID %d tidak ditemukan", id)
}

func (repository *PohonKinerjaRepositoryImpl) ValidateParentLevel(ctx context.Context, tx *sql.Tx, parentId int, levelPohon int) error {
	// Validasi dasar: level tidak boleh kurang dari 4
	if levelPohon < 4 {
		return fmt.Errorf("level pohon tidak boleh kurang dari 4")
	}

	// Untuk level 4, parent bisa memiliki level 0 hingga 3
	if levelPohon == 4 {
		if parentId == 0 {
			return nil
		}
		// Cek level parentnya
		script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
		var parentLevel int
		err := tx.QueryRowContext(ctx, script, parentId).Scan(&parentLevel)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("parent dengan ID %d tidak ditemukan", parentId)
			}
			return fmt.Errorf("gagal memeriksa level parent: %v", err)
		}

		// Validasi level parent untuk level 4
		if parentLevel < 0 || parentLevel > 3 {
			return fmt.Errorf("level pohon 4 harus memiliki parent dengan level 0 hingga 3, bukan level %d", parentLevel)
		}
		return nil
	}

	// Untuk level > 4, parent tidak boleh 0
	if parentId == 0 {
		return fmt.Errorf("level pohon %d harus memiliki parent", levelPohon)
	}

	// Cek level parent untuk level > 4
	script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
	var parentLevel int
	err := tx.QueryRowContext(ctx, script, parentId).Scan(&parentLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parent dengan ID %d tidak ditemukan", parentId)
		}
		return fmt.Errorf("gagal memeriksa level parent: %v", err)
	}

	// Validasi: level parent harus tepat 1 tingkat di atas level saat ini
	expectedParentLevel := levelPohon - 1
	if parentLevel != expectedParentLevel {
		return fmt.Errorf("level pohon %d harus memiliki parent dengan level %d, bukan level %d",
			levelPohon, expectedParentLevel, parentLevel)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorToClone(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.Indikator, error) {
	script := "SELECT id, indikator, tahun FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data indikator: %v", err)
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data indikator: %v", err)
		}
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetToClone(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data target: %v", err)
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data target: %v", err)
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPokin(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error) {
	script := `INSERT INTO tb_pohon_kinerja 
        (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from, is_active) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pokin.NamaPohon,
		pokin.Parent,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status,
		pokin.CloneFrom,
		pokin.IsActive,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal menyimpan data pohon kinerja yang di-clone: %v", err)
	}
	return result.LastInsertId()
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedIndikator(ctx context.Context, tx *sql.Tx, indikatorId string, pokinId int64, indikator domain.Indikator) error {
	script := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script,
		indikatorId,
		pokinId,
		indikator.Indikator,
		indikator.Tahun,
		indikator.Id, // Id indikator asli sebagai clone_from
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan indikator baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedTarget(ctx context.Context, tx *sql.Tx, targetId string, indikatorId string, target domain.Target) error {
	fmt.Printf("Inserting target with clone_from: %s\n", target.Id) // Log sementara
	script := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, clone_from) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script,
		targetId,
		indikatorId,
		target.Target,
		target.Satuan,
		target.Tahun,
		target.Id, // Id target asli sebagai clone_from
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan target baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByJenisPohon(ctx context.Context, tx *sql.Tx, jenisPohon string, levelPohon int, tahun string, kodeOpd string, status string) ([]domain.PohonKinerja, error) {
	script := `
SELECT id, nama_pohon, jenis_pohon, level_pohon, kode_opd, tahun, keterangan, status, is_active
FROM tb_pohon_kinerja
WHERE 1=1`
	parameters := []interface{}{}
	if jenisPohon != "" {
		script += " AND jenis_pohon = ?"
		parameters = append(parameters, jenisPohon)
	}
	if levelPohon != 0 {
		script += " AND level_pohon = ?"
		parameters = append(parameters, levelPohon)
	}
	if kodeOpd != "" {
		script += " AND kode_opd = ?"
		parameters = append(parameters, kodeOpd)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		parameters = append(parameters, tahun)
	}
	if status != "" {
		script += " AND status = ?"
		parameters = append(parameters, status)
	}
	// Hierarki yang akarnya langsung bertaut ke parent -100 tidak ditampilkan
	// (akar + seluruh turunan).
	script += `
	AND id NOT IN (
		WITH RECURSIVE excluded_tree AS (
			SELECT id FROM tb_pohon_kinerja WHERE parent = -100
			UNION ALL
			SELECT p.id FROM tb_pohon_kinerja p
			INNER JOIN excluded_tree et ON p.parent = et.id
		)
		SELECT id FROM excluded_tree
	)`
	script += " ORDER BY nama_pohon ASC"
	rows, err := tx.QueryContext(ctx, script, parameters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Tahun, &pokin.Keterangan, &pokin.Status, &pokin.IsActive)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, rows.Err()
}

// func (repository *PohonKinerjaRepositoryImpl) FindPokinByJenisPohon(ctx context.Context, tx *sql.Tx, jenisPohon string, levelPohon int, tahun string, kodeOpd string, status string) ([]domain.PohonKinerja, error) {
// 	script := "SELECT id, nama_pohon, jenis_pohon, level_pohon, kode_opd, tahun, keterangan, status, is_active FROM tb_pohon_kinerja WHERE 1=1"
// 	parameters := []interface{}{}
// 	if jenisPohon != "" {
// 		script += " AND jenis_pohon = ?"
// 		parameters = append(parameters, jenisPohon)
// 	}
// 	if levelPohon != 0 {
// 		script += " AND level_pohon = ?"
// 		parameters = append(parameters, levelPohon)
// 	}
// 	if kodeOpd != "" {
// 		script += " AND kode_opd = ?"
// 		parameters = append(parameters, kodeOpd)
// 	}
// 	if tahun != "" {
// 		script += " AND tahun = ?"
// 		parameters = append(parameters, tahun)
// 	}
// 	if status != "" {
// 		script += " AND status = ?"
// 		parameters = append(parameters, status)
// 	}
// 	script += " ORDER BY nama_pohon asc"

// 	rows, err := tx.QueryContext(ctx, script, parameters...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var pokins []domain.PohonKinerja
// 	for rows.Next() {
// 		var pokin domain.PohonKinerja
// 		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Tahun, &pokin.Keterangan, &pokin.Status, &pokin.IsActive)
// 		if err != nil {
// 			return nil, err
// 		}
// 		pokins = append(pokins, pokin)
// 	}
// 	return pokins, nil
// }

func (repository *PohonKinerjaRepositoryImpl) FindPokinByPelaksana(ctx context.Context, tx *sql.Tx, nip string, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT DISTINCT
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            pk.created_at,
            pp.id as pelaksana_id,
            pp.pegawai_id,
            p.nip,
            p.nama as nama_pegawai
        FROM 
            tb_pohon_kinerja pk
        INNER JOIN 
            tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
        INNER JOIN 
            tb_pegawai p ON pp.pegawai_id = p.id
        WHERE 
            p.nip = ?  -- ✅ FILTER BERDASARKAN NIP
            AND pk.tahun = ?
        ORDER BY 
            pk.level_pohon, pk.id, pk.created_at ASC
    `

	rows, err := tx.QueryContext(ctx, script, nip, tahun) // ✅ PARAMETER NIP
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}
	defer rows.Close()

	pokinMap := make(map[int]domain.PohonKinerja)

	for rows.Next() {
		var pokin domain.PohonKinerja
		var pelaksana domain.PelaksanaPokin

		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.CreatedAt,
			&pelaksana.Id,
			&pelaksana.PegawaiId,
			&pelaksana.Nip,         // ✅ TAMBAHKAN NIP
			&pelaksana.NamaPegawai, // ✅ TAMBAHKAN NAMA PEGAWAI
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data pohon kinerja: %v", err)
		}

		// Cek apakah pohon kinerja sudah ada di map
		existingPokin, exists := pokinMap[pokin.Id]
		if exists {
			// Jika sudah ada, tambahkan pelaksana baru ke slice pelaksana yang ada
			existingPokin.Pelaksana = append(existingPokin.Pelaksana, pelaksana)
			pokinMap[pokin.Id] = existingPokin
		} else {
			// Jika belum ada, buat entry baru dengan pelaksana pertama
			pokin.Pelaksana = []domain.PelaksanaPokin{pelaksana}
			pokinMap[pokin.Id] = pokin
		}
	}

	// Konversi map ke slice untuk hasil akhir
	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		result = append(result, pokin)
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, status string) ([]domain.PohonKinerja, error) {
	SQL := `SELECT id, nama_pohon, kode_opd, tahun, jenis_pohon, level_pohon, parent, status 
            FROM tb_pohon_kinerja 
            WHERE kode_opd = ? AND tahun = ? AND status = ?`

	rows, err := tx.QueryContext(ctx, SQL, kodeOpd, tahun, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		pokin := domain.PohonKinerja{}
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.KodeOpd,
			&pokin.Tahun,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.Parent,
			&pokin.Status,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	script := "UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, status, id)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckPokinStatus(ctx context.Context, tx *sql.Tx, id int) (string, error) {
	script := "SELECT status FROM tb_pohon_kinerja WHERE id = ?"
	var status string
	err := tx.QueryRowContext(ctx, script, id).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("gagal mengecek status: %v", err)
	}
	return status, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPokinWithStatus(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error) {
	script := `INSERT INTO tb_pohon_kinerja 
        (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pokin.NamaPohon,
		pokin.Parent,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status,
		pokin.CloneFrom,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal menyimpan data pohon kinerja yang di-clone: %v", err)
	}
	return result.LastInsertId()
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinStatusTolak(ctx context.Context, tx *sql.Tx, id int, status string) error {
	script := "UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, status, id)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status dan alasan: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckCloneFrom(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	script := "SELECT COALESCE(clone_from, 0) FROM tb_pohon_kinerja WHERE id = ?"
	var cloneFrom int
	err := tx.QueryRowContext(ctx, script, id).Scan(&cloneFrom)
	if err != nil {
		return 0, fmt.Errorf("gagal mengecek clone_from: %v", err)
	}
	return cloneFrom, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByCloneFrom(ctx context.Context, tx *sql.Tx, cloneFromId int) ([]domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from FROM tb_pohon_kinerja WHERE clone_from = ?"
	rows, err := tx.QueryContext(ctx, script, cloneFromId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.Parent,
			&pokin.NamaPohon,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&pokin.CloneFrom,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorByCloneFrom(ctx context.Context, tx *sql.Tx, pokinId int, cloneFromId string) (domain.Indikator, error) {
	script := "SELECT id, indikator, tahun FROM tb_indikator WHERE pokin_id = ? AND clone_from = ?"
	var indikator domain.Indikator
	err := tx.QueryRowContext(ctx, script, pokinId, cloneFromId).Scan(
		&indikator.Id,
		&indikator.Indikator,
		&indikator.Tahun,
	)
	if err != nil {
		return domain.Indikator{}, err
	}
	return indikator, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetByCloneFrom(ctx context.Context, tx *sql.Tx, indikatorId string, cloneFromId string) (domain.Target, error) {
	script := "SELECT id, target, satuan, tahun FROM tb_target WHERE indikator_id = ? AND clone_from = ?"
	var target domain.Target
	err := tx.QueryRowContext(ctx, script, indikatorId, cloneFromId).Scan(
		&target.Id,
		&target.Target,
		&target.Satuan,
		&target.Tahun,
	)
	if err != nil {
		return domain.Target{}, err
	}
	return target, nil
}

// Tambahkan method baru untuk FindPokinByCrosscuttingStatus
func (repository *PohonKinerjaRepositoryImpl) FindPokinByCrosscuttingStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PohonKinerja, error) {
	script := `SELECT 
        id, nama_pohon, parent, jenis_pohon, level_pohon, 
        kode_opd, keterangan, tahun, status 
        FROM tb_pohon_kinerja 
        WHERE kode_opd = ? 
        AND tahun = ? 
        AND status IN ('crosscutting_menunggu','crosscutting_ditolak')
        ORDER BY level_pohon, id ASC`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id, &pokin.NamaPohon, &pokin.Parent, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Keterangan, &pokin.Tahun, &pokin.Status,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeleteClonedPokinHierarchy(ctx context.Context, tx *sql.Tx, id int) error {
	// Query untuk mendapatkan hierarki dari data clone
	findIdsScript := `
        WITH RECURSIVE clone_hierarki AS (
            -- Base case: node clone yang akan dihapus
            SELECT id, parent, level_pohon, clone_from 
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: child nodes dari data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN clone_hierarki ch ON 
                pk.parent = ch.id
        )
        SELECT id FROM clone_hierarki;
    `

	rows, err := tx.QueryContext(ctx, findIdsScript, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil hierarki clone: %v", err)
	}
	defer rows.Close()

	var idsToDelete []interface{}
	for rows.Next() {
		var idToDelete int
		if err := rows.Scan(&idToDelete); err != nil {
			return fmt.Errorf("gagal scan ID: %v", err)
		}
		idsToDelete = append(idsToDelete, idToDelete)
	}

	if len(idsToDelete) == 0 {
		return fmt.Errorf("tidak ada data yang akan dihapus")
	}

	// Buat placeholder untuk query IN clause
	placeholders := make([]string, len(idsToDelete))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	// Hapus target terlebih dahulu
	scriptDeleteTarget := fmt.Sprintf(`
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id IN (%s)
        )`, inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus target: %v", err)
	}

	// Hapus indikator
	scriptDeleteIndikator := fmt.Sprintf("DELETE FROM tb_indikator WHERE pokin_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus indikator: %v", err)
	}

	// Hapus pelaksana
	scriptDeletePelaksana := fmt.Sprintf("DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pelaksana: %v", err)
	}

	// Hapus pohon kinerja
	scriptDeletePokin := fmt.Sprintf("DELETE FROM tb_pohon_kinerja WHERE id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePokin, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

func (r *PohonKinerjaRepositoryImpl) FindChildPokins(ctx context.Context, tx *sql.Tx, parentId int64) ([]domain.PohonKinerja, error) {
	SQL := `SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from, is_active
            FROM tb_pohon_kinerja 
            WHERE parent = ?`

	rows, err := tx.QueryContext(ctx, SQL, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.Parent,
			&pokin.NamaPohon,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&pokin.CloneFrom,
			&pokin.IsActive,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPelaksana(ctx context.Context, tx *sql.Tx, newId string, pokinId int64, pelaksana domain.PelaksanaPokin) error {
	SQL := `INSERT INTO tb_pelaksana_pokin (id, pokin_id, pegawai_id)
            VALUES (?, ?, ?)`

	_, err := tx.ExecContext(ctx, SQL, newId, pokinId, pelaksana.PegawaiId)
	return err
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinStatusFromApproved(ctx context.Context, tx *sql.Tx, id int) error {
	SQL := `
        UPDATE tb_pohon_kinerja 
        SET status = 'menunggu_disetujui' 
        WHERE id = ? 
        AND status = 'disetujui'
    `

	result, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("tidak ada data yang diupdate untuk ID %d", id)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdateParent(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := `UPDATE tb_pohon_kinerja SET parent = ? WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, pohonKinerja.Parent, pohonKinerja.Id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindidPokinWithAllTema(ctx context.Context, tx *sql.Tx, id int) ([]domain.PohonKinerja, error) {
	script := `
                 WITH RECURSIVE ancestor_tree AS (
            -- Base case: node yang dicari
            SELECT 
                pk.id, pk.nama_pohon, pk.parent, pk.jenis_pohon, 
                pk.level_pohon, pk.kode_opd, pk.keterangan, 
                pk.tahun, pk.status,
                i.id as indikator_id, i.indikator as nama_indikator,
                t.id as target_id, t.target as target_value, 
                t.satuan as target_satuan,
                pp.id as pelaksana_id, pp.pegawai_id
            FROM tb_pohon_kinerja pk
            LEFT JOIN tb_indikator i ON pk.id = i.pokin_id
            LEFT JOIN tb_target t ON i.id = t.indikator_id
            LEFT JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
            WHERE pk.id = ?
            
            UNION ALL
            
            -- Recursive case: ambil parent nodes
            SELECT 
                pk.id, pk.nama_pohon, pk.parent, pk.jenis_pohon, 
                pk.level_pohon, pk.kode_opd, pk.keterangan, 
                pk.tahun, pk.status,
                i.id as indikator_id, i.indikator as nama_indikator,
                t.id as target_id, t.target as target_value, 
                t.satuan as target_satuan,
                pp.id as pelaksana_id, pp.pegawai_id
            FROM tb_pohon_kinerja pk
            LEFT JOIN tb_indikator i ON pk.id = i.pokin_id
            LEFT JOIN tb_target t ON i.id = t.indikator_id
            LEFT JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
            INNER JOIN ancestor_tree at ON pk.id = at.parent
        )
        SELECT * FROM ancestor_tree
        ORDER BY level_pohon ASC`

	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]*domain.PohonKinerja)
	processedIndikators := make(map[string]bool)

	for rows.Next() {
		var (
			pokinId                                    int
			namaPohon, jenisPohon, kodeOpd, keterangan string
			parent, levelPohon                         int
			tahun, status                              string
			indikatorId, namaIndikator                 sql.NullString
			targetId, targetValue, targetSatuan        sql.NullString
			pelaksanaId, pegawaiId                     sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon,
			&levelPohon, &kodeOpd, &keterangan,
			&tahun, &status,
			&indikatorId, &namaIndikator,
			&targetId, &targetValue, &targetSatuan,
			&pelaksanaId, &pegawaiId,
		)
		if err != nil {
			return nil, err
		}

		// Inisialisasi atau ambil pohon kinerja dari map
		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = &domain.PohonKinerja{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahun,
				Status:     status,
			}
			pokinMap[pokinId] = pokin
		}

		// Proses Indikator dan Target
		if indikatorId.Valid && namaIndikator.Valid {
			if !processedIndikators[indikatorId.String] {
				processedIndikators[indikatorId.String] = true
				indikator := domain.Indikator{
					Id:        indikatorId.String,
					PokinId:   fmt.Sprint(pokinId),
					Indikator: namaIndikator.String,
					Tahun:     tahun,
				}

				if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
					target := domain.Target{
						Id:          targetId.String,
						IndikatorId: indikatorId.String,
						Target:      targetValue.String,
						Satuan:      targetSatuan.String,
						Tahun:       tahun,
					}
					indikator.Target = append(indikator.Target, target)
				}

				pokin.Indikator = append(pokin.Indikator, indikator)
			}
		}

		// Proses Pelaksana
		if pelaksanaId.Valid && pegawaiId.Valid {
			pelaksana := domain.PelaksanaPokin{
				Id:        pelaksanaId.String,
				PegawaiId: pegawaiId.String,
			}
			// Cek duplikasi
			isDuplicate := false
			for _, p := range pokin.Pelaksana {
				if p.Id == pelaksana.Id {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				pokin.Pelaksana = append(pokin.Pelaksana, pelaksana)
			}
		}
	}

	// Convert map to sorted slice
	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		result = append(result, *pokin)
	}

	// Sort berdasarkan level_pohon
	sort.Slice(result, func(i, j int) bool {
		return result[i].LevelPohon < result[j].LevelPohon
	})

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) ValidatePokinId(ctx context.Context, tx *sql.Tx, pokinId int) error {
	script := "SELECT COUNT(*) FROM tb_pohon_kinerja WHERE id = ?"

	var count int
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(&count)
	if err != nil {
		return fmt.Errorf("gagal melakukan validasi pohon kinerja: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("pohon kinerja dengan ID %d tidak ditemukan", pokinId)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) ValidatePokinLevel(ctx context.Context, tx *sql.Tx, pokinId int, expectedLevel int, purpose string) error {
	script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"

	var levelPohon int
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(&levelPohon)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("pohon kinerja dengan ID %d tidak ditemukan", pokinId)
		}
		return fmt.Errorf("gagal melakukan validasi pohon kinerja: %v", err)
	}

	if levelPohon != expectedLevel {
		return fmt.Errorf("%s hanya bisa dibuat dari pohon kinerja level %d, bukan level %d",
			purpose, expectedLevel, levelPohon)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinWithPeriode(ctx context.Context, tx *sql.Tx, pokinId int, jenisPeriode string) (domain.PohonKinerja, domain.Periode, error) {
	query := `
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            pk.status,
            COALESCE(p.id, 0) as periode_id,
            COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            i.rumus_perhitungan,
            i.sumber_data,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan
        FROM 
            tb_pohon_kinerja pk
        LEFT JOIN 
            tb_periode p ON CAST(pk.tahun AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
            AND p.jenis_periode = ?
        LEFT JOIN 
            tb_indikator i ON pk.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        WHERE 
            pk.id = ?
        ORDER BY 
            i.id`

	rows, err := tx.QueryContext(ctx, query, jenisPeriode, pokinId)
	if err != nil {
		return domain.PohonKinerja{}, domain.Periode{}, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()

	var pokin domain.PohonKinerja
	var periode domain.Periode
	indikatorMap := make(map[string]*domain.Indikator)
	firstRow := true

	for rows.Next() {
		var (
			periodeId                           int
			tahunAwal, tahunAkhir               string
			indikatorId, indikatorText          sql.NullString
			rumusPerhitungan, sumberData        sql.NullString
			targetId, targetValue, targetSatuan sql.NullString
		)

		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetValue,
			&targetSatuan,
		)
		if err != nil {
			return domain.PohonKinerja{}, domain.Periode{}, fmt.Errorf("error scanning row: %v", err)
		}

		if firstRow {
			periode = domain.Periode{
				Id:         periodeId,
				TahunAwal:  tahunAwal,
				TahunAkhir: tahunAkhir,
			}
			firstRow = false
		}

		// Proses indikator jika ada
		if indikatorId.Valid && indikatorText.Valid {
			indikator, exists := indikatorMap[indikatorId.String]
			if !exists {
				indikator = &domain.Indikator{
					Id:               indikatorId.String,
					PokinId:          fmt.Sprint(pokin.Id),
					Indikator:        indikatorText.String,
					RumusPerhitungan: rumusPerhitungan,
					SumberData:       sumberData,
					Target:           []domain.Target{},
				}
				indikatorMap[indikatorId.String] = indikator

				// Buat target untuk setiap tahun dalam periode
				if periode.Id != 0 && periode.TahunAwal != "" && periode.TahunAkhir != "" {
					tahunAwalInt, _ := strconv.Atoi(periode.TahunAwal)
					tahunAkhirInt, _ := strconv.Atoi(periode.TahunAkhir)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						target := domain.Target{
							Id:          "-",
							IndikatorId: indikatorId.String,
							Target:      "",
							Satuan:      "",
							Tahun:       tahunStr,
						}
						indikator.Target = append(indikator.Target, target)
					}
				}
			}

			// Update target jika ada data
			if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
				for i := range indikator.Target {
					if indikator.Target[i].Tahun == pokin.Tahun {
						indikator.Target[i] = domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       pokin.Tahun,
						}
						break
					}
				}
			}
		}
	}

	// Konversi map indikator ke slice
	for _, indikator := range indikatorMap {
		// Sort target berdasarkan tahun
		sort.Slice(indikator.Target, func(i, j int) bool {
			tahunI, _ := strconv.Atoi(indikator.Target[i].Tahun)
			tahunJ, _ := strconv.Atoi(indikator.Target[j].Tahun)
			return tahunI < tahunJ
		})
		pokin.Indikator = append(pokin.Indikator, *indikator)
	}

	// Sort indikator berdasarkan ID
	sort.Slice(pokin.Indikator, func(i, j int) bool {
		return pokin.Indikator[i].Id < pokin.Indikator[j].Id
	})

	if pokin.Id == 0 {
		return pokin, periode, fmt.Errorf("pohon kinerja tidak ditemukan")
	}

	return pokin, periode, nil
}

// aktif / nonaktif tematik
func (repository *PohonKinerjaRepositoryImpl) UpdateTematikStatus(ctx context.Context, tx *sql.Tx, id int, isActive bool) error {
	query := "UPDATE tb_pohon_kinerja SET is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, isActive, id)
	return err
}

func (repository *PohonKinerjaRepositoryImpl) GetChildrenAndClones(ctx context.Context, tx *sql.Tx, parentId int, isActivating bool) ([]int, error) {
	var query string

	if isActivating {
		// Query untuk mengaktifkan: ambil semua yang terhubung tanpa memandang status is_active
		query = `
            WITH RECURSIVE tree AS (
                -- Base case: direct children and clones yang nonaktif
                SELECT id, parent, clone_from, level_pohon
                FROM tb_pohon_kinerja
                WHERE (parent = ? OR clone_from = ?) 
                AND is_active = false
                
                UNION ALL
                
                -- Recursive case: children dan clone yang nonaktif
                SELECT pk.id, pk.parent, pk.clone_from, pk.level_pohon
                FROM tb_pohon_kinerja pk
                INNER JOIN tree t ON (pk.parent = t.id OR pk.clone_from = t.id)
                WHERE pk.is_active = false
            )
            SELECT DISTINCT id FROM tree`
	} else {
		// Query untuk menonaktifkan: ambil semua yang terhubung dan masih aktif
		query = `
            WITH RECURSIVE tree AS (
                -- Base case: direct children and clones yang aktif
                SELECT id, parent, clone_from, level_pohon
                FROM tb_pohon_kinerja
                WHERE (parent = ? OR clone_from = ?)
                AND is_active = true
                
                UNION ALL
                
                -- Recursive case: children dan clone yang aktif
                SELECT pk.id, pk.parent, pk.clone_from, pk.level_pohon
                FROM tb_pohon_kinerja pk
                INNER JOIN tree t ON (pk.parent = t.id OR pk.clone_from = t.id)
                WHERE pk.is_active = true
            )
            SELECT DISTINCT id FROM tree`
	}

	rows, err := tx.QueryContext(ctx, query, parentId, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

// clone pokin opd
func (repository *PohonKinerjaRepositoryImpl) IsExistsByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) bool {
	script := "SELECT COUNT(*) FROM tb_pohon_kinerja WHERE kode_opd = ? AND tahun = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, kodeOpd, tahun).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

//clone pokin lama
// func (repository *PohonKinerjaRepositoryImpl) ClonePokinOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, sourceTahun string, targetTahun string) error {
// 	// 1. Dapatkan daftar ID yang valid (status kosong dan parent dengan status kosong)
// 	scriptValidIds := `
// 	  SELECT p1.id
// 	  FROM tb_pohon_kinerja p1
// 	  LEFT JOIN tb_pohon_kinerja p2 ON p1.parent = p2.id
// 	  WHERE p1.kode_opd = ?
// 	  AND p1.tahun = ?
// 	  AND p1.status = ''
// 	  AND (p1.parent = 0 OR (p2.status = '' OR p2.status IS NULL))
//   `
// 	validRows, err := tx.QueryContext(ctx, scriptValidIds, kodeOpd, sourceTahun)
// 	if err != nil {
// 		return err
// 	}
// 	defer validRows.Close()

// 	var validIds []int
// 	for validRows.Next() {
// 		var id int
// 		if err := validRows.Scan(&id); err != nil {
// 			return err
// 		}
// 		validIds = append(validIds, id)
// 	}

// 	// 2. Clone Pohon Kinerja yang valid
// 	for _, validId := range validIds {
// 		scriptPokin := `
// 		  INSERT INTO tb_pohon_kinerja (
// 			  nama_pohon, parent, jenis_pohon, level_pohon,
// 			  kode_opd, keterangan, keterangan_crosscutting,
// 			  tahun, status, is_active
// 		  )
// 		  SELECT
// 			  nama_pohon, parent, jenis_pohon, level_pohon,
// 			  kode_opd, keterangan, keterangan_crosscutting,
// 			  ?, '', is_active
// 		  FROM tb_pohon_kinerja
// 		  WHERE id = ?
// 	  `
// 		_, err := tx.ExecContext(ctx, scriptPokin, targetTahun, validId)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// 3. Dapatkan mapping ID lama ke ID baru
// 	scriptMapping := `
// 	  SELECT
// 		  src.id as old_id,
// 		  dst.id as new_id
// 	  FROM tb_pohon_kinerja src
// 	  JOIN tb_pohon_kinerja dst ON
// 		  src.nama_pohon = dst.nama_pohon AND
// 		  src.kode_opd = dst.kode_opd AND
// 		  src.level_pohon = dst.level_pohon
// 	  WHERE src.tahun = ? AND dst.tahun = ?
// 	  AND src.kode_opd = ?
//   `
// 	rows, err := tx.QueryContext(ctx, scriptMapping, sourceTahun, targetTahun, kodeOpd)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	idMapping := make(map[int]int)
// 	for rows.Next() {
// 		var oldId, newId int
// 		if err := rows.Scan(&oldId, &newId); err != nil {
// 			return err
// 		}
// 		idMapping[oldId] = newId
// 	}

// 	// 4. Update parent IDs menggunakan mapping
// 	for oldId, newId := range idMapping {
// 		scriptUpdateParent := `
// 		  UPDATE tb_pohon_kinerja
// 		  SET parent = ?
// 		  WHERE id = ? AND tahun = ?
// 	  `
// 		oldPokin := `SELECT parent FROM tb_pohon_kinerja WHERE id = ? AND tahun = ?`
// 		var oldParent int
// 		err := tx.QueryRowContext(ctx, oldPokin, oldId, sourceTahun).Scan(&oldParent)
// 		if err != nil {
// 			continue
// 		}

// 		if newParent, exists := idMapping[oldParent]; exists {
// 			_, err = tx.ExecContext(ctx, scriptUpdateParent, newParent, newId, targetTahun)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	// 5. Clone Indikator untuk pohon kinerja yang valid
// 	for oldId, newId := range idMapping {
// 		// Clone indikator - Tambahkan kolom tahun
// 		scriptIndikator := `
// 			INSERT INTO tb_indikator (
// 				id, pokin_id, indikator, tahun
// 			)
// 			SELECT
// 				CONCAT('IND-', UUID_SHORT()),
// 				?,
// 				indikator,
// 				''
// 			FROM tb_indikator
// 			WHERE pokin_id = ?
// 		`
// 		_, err := tx.ExecContext(ctx, scriptIndikator, newId, oldId)
// 		if err != nil {
// 			return err
// 		}

// 		// Dapatkan mapping ID indikator (menggunakan string untuk ID)
// 		scriptIndikatorMapping := `
//             SELECT
//                 src.id as old_indikator_id,
//                 dst.id as new_indikator_id
//             FROM tb_indikator src
//             JOIN tb_indikator dst ON
//                 src.indikator = dst.indikator AND
//                 dst.pokin_id = ?
//             WHERE src.pokin_id = ?
//         `
// 		indikatorRows, err := tx.QueryContext(ctx, scriptIndikatorMapping, newId, oldId)
// 		if err != nil {
// 			return err
// 		}
// 		defer indikatorRows.Close()

// 		// Ubah tipe mapping menjadi string
// 		indikatorMapping := make(map[string]string)
// 		for indikatorRows.Next() {
// 			var oldIndikatorId, newIndikatorId string
// 			if err := indikatorRows.Scan(&oldIndikatorId, &newIndikatorId); err != nil {
// 				return err
// 			}
// 			indikatorMapping[oldIndikatorId] = newIndikatorId
// 		}

// 		// Clone target untuk setiap indikator
// 		for oldIndikatorId, newIndikatorId := range indikatorMapping {
// 			scriptTarget := `
// 				INSERT INTO tb_target (
// 					id, indikator_id, target, satuan, tahun
// 				)
// 				SELECT
// 					CONCAT('TRG-', UUID_SHORT()),
// 					?,
// 					target,
// 					satuan,
// 					?
// 				FROM tb_target
// 				WHERE indikator_id = ?
// 			`
// 			_, err := tx.ExecContext(ctx, scriptTarget, newIndikatorId, targetTahun, oldIndikatorId)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	// Clone tagging untuk setiap pohon kinerja
// 	for oldId, newId := range idMapping {
// 		// Clone tagging dengan menambahkan clone_from
// 		scriptTagging := `
// 		INSERT INTO tb_tagging_pokin (
// 			id_pokin,
// 			nama_tagging,
// 			keterangan_tagging,
// 			clone_from
// 		)
// 		SELECT
// 			?,  -- new pokin_id
// 			nama_tagging,
// 			keterangan_tagging,
// 			id  -- menyimpan id lama sebagai clone_from
// 		FROM tb_tagging_pokin
// 		WHERE id_pokin = ?  -- old pokin_id
// 	`
// 		_, err := tx.ExecContext(ctx, scriptTagging, newId, oldId)
// 		if err != nil {
// 			return fmt.Errorf("gagal mengkloning tagging: %v", err)
// 		}
// 	}

// 	return nil
// }

// clone pokin opd baru
func (repository *PohonKinerjaRepositoryImpl) ClonePokinOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, sourceTahun string, targetTahun string) error {
	const parentPlaceholderPemda = -100
	const statusPokinPemda = "pokin dari pemda"
	// 1. ID valid: status kosong, parent = 0 atau parent status kosong/NULL, atau parent "pokin dari pemda"
	scriptValidIds := `
	  SELECT p1.id 
	  FROM tb_pohon_kinerja p1
	  LEFT JOIN tb_pohon_kinerja p2 ON p1.parent = p2.id 
	  WHERE p1.kode_opd = ? 
	  AND p1.tahun = ?
	  AND p1.status = ''
	  AND (
	    p1.parent = 0 
	    OR (p2.status = '' OR p2.status IS NULL)
	    OR p2.status = 'pokin dari pemda'
	  )
	`
	validRows, err := tx.QueryContext(ctx, scriptValidIds, kodeOpd, sourceTahun)
	if err != nil {
		return err
	}
	defer validRows.Close()
	var validIds []int
	for validRows.Next() {
		var id int
		if err := validRows.Scan(&id); err != nil {
			return err
		}
		validIds = append(validIds, id)
	}
	// 2. Clone Pohon Kinerja yang valid
	for _, validId := range validIds {
		scriptPokin := `
	  INSERT INTO tb_pohon_kinerja (
		  nama_pohon, parent, jenis_pohon, level_pohon, 
		  kode_opd, keterangan, keterangan_crosscutting, 
		  tahun, status, is_active, keterangan_tahun_clone
	  )
	  SELECT 
		  nama_pohon, parent, jenis_pohon, level_pohon,
		  kode_opd, keterangan, keterangan_crosscutting,
		  ?, '', is_active, ?
	  FROM tb_pohon_kinerja
	  WHERE id = ?
	`
		_, err := tx.ExecContext(ctx, scriptPokin, targetTahun, sourceTahun, validId)
		if err != nil {
			return err
		}
	}
	// 3. Mapping ID lama -> baru
	scriptMapping := `
	  SELECT 
		  src.id as old_id,
		  dst.id as new_id
	  FROM tb_pohon_kinerja src
	  JOIN tb_pohon_kinerja dst ON 
		  src.nama_pohon = dst.nama_pohon AND
		  src.kode_opd = dst.kode_opd AND
		  src.level_pohon = dst.level_pohon
	  WHERE src.tahun = ? AND dst.tahun = ? 
	  AND src.kode_opd = ?
	`
	rows, err := tx.QueryContext(ctx, scriptMapping, sourceTahun, targetTahun, kodeOpd)
	if err != nil {
		return err
	}
	defer rows.Close()
	idMapping := make(map[int]int)
	for rows.Next() {
		var oldId, newId int
		if err := rows.Scan(&oldId, &newId); err != nil {
			return err
		}
		idMapping[oldId] = newId
	}
	// 4. Update parent: mapping biasa, atau -100 jika parent sumber "pokin dari pemda" (tidak diklon)
	scriptUpdateParent := `
		  UPDATE tb_pohon_kinerja 
		  SET parent = ?
		  WHERE id = ? AND tahun = ?
	  `
	oldPokin := `SELECT parent FROM tb_pohon_kinerja WHERE id = ? AND tahun = ?`
	parentStatusQuery := `SELECT COALESCE(status, '') FROM tb_pohon_kinerja WHERE id = ? AND tahun = ?`
	for oldId, newId := range idMapping {
		var oldParent int
		err := tx.QueryRowContext(ctx, oldPokin, oldId, sourceTahun).Scan(&oldParent)
		if err != nil {
			continue
		}
		if oldParent == 0 {
			continue
		}
		if newParent, exists := idMapping[oldParent]; exists {
			_, err = tx.ExecContext(ctx, scriptUpdateParent, newParent, newId, targetTahun)
			if err != nil {
				return err
			}
			continue
		}
		var parentStatus string
		err = tx.QueryRowContext(ctx, parentStatusQuery, oldParent, sourceTahun).Scan(&parentStatus)
		if err != nil {
			continue
		}
		if parentStatus == statusPokinPemda {
			_, err = tx.ExecContext(ctx, scriptUpdateParent, parentPlaceholderPemda, newId, targetTahun)
			if err != nil {
				return err
			}
		}
	}
	// 5. Clone indikator
	for oldId, newId := range idMapping {
		scriptIndikator := `
			INSERT INTO tb_indikator (
				id, pokin_id, indikator, tahun
			)
			SELECT 
				CONCAT('IND-', UUID_SHORT()), 
				?, 
				indikator,
				''
			FROM tb_indikator
			WHERE pokin_id = ?
		`
		_, err := tx.ExecContext(ctx, scriptIndikator, newId, oldId)
		if err != nil {
			return err
		}
		scriptIndikatorMapping := `
            SELECT 
                src.id as old_indikator_id,
                dst.id as new_indikator_id
            FROM tb_indikator src
            JOIN tb_indikator dst ON 
                src.indikator = dst.indikator AND
                dst.pokin_id = ?
            WHERE src.pokin_id = ?
        `
		indikatorRows, err := tx.QueryContext(ctx, scriptIndikatorMapping, newId, oldId)
		if err != nil {
			return err
		}
		indikatorMapping := make(map[string]string)
		for indikatorRows.Next() {
			var oldIndikatorId, newIndikatorId string
			if err := indikatorRows.Scan(&oldIndikatorId, &newIndikatorId); err != nil {
				indikatorRows.Close()
				return err
			}
			indikatorMapping[oldIndikatorId] = newIndikatorId
		}
		indikatorRows.Close()
		for oldIndikatorId, newIndikatorId := range indikatorMapping {
			scriptTarget := `
				INSERT INTO tb_target (
					id, indikator_id, target, satuan, tahun
				)
				SELECT 
					CONCAT('TRG-', UUID_SHORT()), 
					?, 
					target, 
					satuan,
					? 
				FROM tb_target
				WHERE indikator_id = ?
			`
			_, err := tx.ExecContext(ctx, scriptTarget, newIndikatorId, targetTahun, oldIndikatorId)
			if err != nil {
				return err
			}
		}
	}
	// Clone tagging
	for oldId, newId := range idMapping {
		scriptTagging := `
		INSERT INTO tb_tagging_pokin (
			id_pokin,
			nama_tagging,
			keterangan_tagging,
			clone_from
		)
		SELECT 
			?,
			nama_tagging,
			keterangan_tagging,
			id
		FROM tb_tagging_pokin
		WHERE id_pokin = ?
	`
		_, err := tx.ExecContext(ctx, scriptTagging, newId, oldId)
		if err != nil {
			return fmt.Errorf("gagal mengkloning tagging: %v", err)
		}
	}
	return nil
}

// count pokin pemda in opd
func (repository *PohonKinerjaRepositoryImpl) CountPokinPemdaByLevel(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) (map[int]int, error) {
	script := `
  WITH RECURSIVE pohon_all AS (
    SELECT 
        id,
        parent,
        level_pohon,
        status,
        jenis_pohon,
        CASE 
            WHEN level_pohon = 4 AND parent = 0 AND (
                status = 'pokin dari pemda'
                OR (
                    status = 'crosscutting_disetujui_existing'
                    AND jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                )
            ) THEN TRUE
            ELSE FALSE
        END as is_counted
    FROM tb_pohon_kinerja
    WHERE kode_opd = ? AND tahun = ?
),
valid_level_4 AS (
    SELECT id 
    FROM pohon_all 
    WHERE level_pohon = 4 
    AND parent = 0
    AND (
        status = 'pokin dari pemda'
        OR (
            status = 'crosscutting_disetujui_existing'
            AND jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
        )
    )
),
pohon_hierarchy AS (
    SELECT 
        p.*,
        p.is_counted as should_count
    FROM pohon_all p
    WHERE p.level_pohon = 4
    UNION ALL
    SELECT 
        child.*,
        CASE
            WHEN child.level_pohon = 5 AND (
                child.status = 'pokin dari pemda'
                OR (
                    child.status = 'crosscutting_disetujui_existing'
                    AND child.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                )
            ) THEN
                CASE
                    WHEN (
                        parent.status = 'pokin dari pemda'
                        OR (
                            parent.status = 'crosscutting_disetujui_existing'
                            AND parent.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                        )
                    ) THEN
                        CASE WHEN (SELECT p2.parent FROM pohon_all p2 WHERE p2.id = parent.id) = 0 THEN TRUE
                        ELSE FALSE END
                    WHEN parent.status = '' THEN TRUE
                    WHEN EXISTS (
                        SELECT 1 FROM valid_level_4 
                        WHERE id = child.parent
                    ) THEN TRUE
                    ELSE FALSE
                END
            WHEN child.level_pohon >= 6 AND (
                child.status = 'pokin dari pemda'
                OR (
                    child.status = 'crosscutting_disetujui_existing'
                    AND child.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                )
            ) THEN
                CASE
                    WHEN (
                        parent.status = 'pokin dari pemda'
                        OR (
                            parent.status = 'crosscutting_disetujui_existing'
                            AND parent.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                        )
                    ) THEN
                        CASE
                            WHEN EXISTS (
                                SELECT 1 FROM pohon_all p2
                                WHERE p2.id = parent.parent
                                AND (
                                    p2.status = 'pokin dari pemda'
                                    OR (
                                        p2.status = 'crosscutting_disetujui_existing'
                                        AND p2.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
                                    )
                                )
                            ) THEN
                                CASE WHEN (SELECT p3.parent FROM pohon_all p3 WHERE p3.id = parent.parent) = 0 THEN TRUE
                                ELSE FALSE END
                            WHEN (SELECT p2.status FROM pohon_all p2 WHERE p2.id = parent.parent) = '' THEN TRUE
                            ELSE FALSE
                        END
                    WHEN parent.status = '' THEN TRUE
                    WHEN EXISTS (
                        WITH RECURSIVE ancestors AS (
                            SELECT p2.id, p2.parent, p2.level_pohon
                            FROM pohon_all p2
                            WHERE p2.id = child.parent
                            
                            UNION ALL
                            
                            SELECT p3.id, p3.parent, p3.level_pohon
                            FROM pohon_all p3
                            JOIN ancestors a ON p3.id = a.parent
                            WHERE p3.level_pohon >= 4
                        )
                        SELECT 1 
                        FROM ancestors a
                        JOIN valid_level_4 v ON v.id = a.id
                        WHERE a.level_pohon = 4
                    ) THEN TRUE
                    ELSE FALSE
                END
            ELSE FALSE
        END as should_count
    FROM pohon_all child
    JOIN pohon_hierarchy parent ON child.parent = parent.id
    WHERE child.level_pohon > 4
)
SELECT 
    level_pohon,
    COUNT(*) as jumlah
FROM pohon_hierarchy
WHERE should_count = TRUE
GROUP BY level_pohon
ORDER BY level_pohon;`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var level, count int
		if err := rows.Scan(&level, &count); err != nil {
			return nil, err
		}
		result[level] = count
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckAsalPokin(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	SQL := "SELECT clone_from FROM pohon_kinerja WHERE id = ?"
	var cloneFrom sql.NullInt64

	err := tx.QueryRowContext(ctx, SQL, id).Scan(&cloneFrom)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}

	if cloneFrom.Valid {
		return int(cloneFrom.Int64), nil
	}
	return 0, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePokinWithIndikatorAndTarget(ctx context.Context, tx *sql.Tx, id int) error {
	// Hapus target terlebih dahulu
	scriptDeleteTarget := `
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id = ?
        )`
	_, err := tx.ExecContext(ctx, scriptDeleteTarget, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus target: %v", err)
	}

	// Hapus indikator
	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus indikator: %v", err)
	}

	// Hapus pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pelaksana: %v", err)
	}

	// Hapus pohon kinerja
	scriptDeletePokin := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePokin, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindListOpdAllTematik(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: pilih semua node level 0
            SELECT 
                id, 
                nama_pohon, 
                parent, 
                jenis_pohon, 
                level_pohon, 
                kode_opd, 
                keterangan, 
                tahun, 
                status, 
                is_active,
                id as root_id
            FROM tb_pohon_kinerja 
            WHERE level_pohon = 0 AND tahun = ?
            
            UNION ALL
            
            -- Recursive case: ambil semua child nodes
            SELECT 
                pk.id, 
                pk.nama_pohon, 
                pk.parent, 
                pk.jenis_pohon, 
                pk.level_pohon, 
                pk.kode_opd, 
                pk.keterangan, 
                pk.tahun, 
                pk.status, 
                pk.is_active,
                ph.root_id
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON pk.parent = ph.id
        )
        SELECT 
            ph.id,
            ph.nama_pohon,
            ph.parent,
            ph.jenis_pohon,
            ph.level_pohon,
            ph.kode_opd,
            ph.keterangan,
            ph.tahun,
            ph.status,
            ph.is_active,
            o.nama_opd,
            ph.root_id
        FROM 
            pohon_hierarki ph
        LEFT JOIN 
            tb_operasional_daerah o ON ph.kode_opd = o.kode_opd
        ORDER BY 
            ph.level_pohon, ph.id
    `

	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan pohon kinerja level 0
	pokinMap := make(map[int]domain.PohonKinerja)
	// Map untuk menyimpan OPD yang sudah diproses per pohon kinerja
	opdMap := make(map[int]map[string]string)

	for rows.Next() {
		var (
			pokinId, parent, levelPohon, rootId                            int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin, status string
			is_active                                                      bool
			namaOpd                                                        sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon, &levelPohon,
			&kodeOpd, &keterangan, &tahunPokin, &status, &is_active,
			&namaOpd, &rootId,
		)
		if err != nil {
			return nil, err
		}

		// Jika ini adalah level 0, simpan ke pokinMap
		if levelPohon == 0 {
			pokinMap[pokinId] = domain.PohonKinerja{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahunPokin,
				Status:     status,
				IsActive:   is_active,
				ListOpd:    []domain.OpdList{},
			}
			// Inisialisasi map OPD untuk level 0
			opdMap[pokinId] = make(map[string]string)
		}

		// Jika ini adalah child node (level > 0) dan memiliki kode_opd
		if levelPohon > 0 && kodeOpd != "" {
			// Gunakan rootId untuk menambahkan OPD ke parent level 0
			if _, exists := opdMap[rootId]; !exists {
				opdMap[rootId] = make(map[string]string)
			}
			// Tambahkan OPD ke map jika belum ada
			if _, exists := opdMap[rootId][kodeOpd]; !exists {
				opdMap[rootId][kodeOpd] = namaOpd.String
			}
		}
	}

	// Konversi map ke slice
	var result []domain.PohonKinerja
	for id, pokin := range pokinMap {
		// Tambahkan list OPD ke pohon kinerja
		for kodeOpd, namaOpd := range opdMap[id] {
			pokin.ListOpd = append(pokin.ListOpd, domain.OpdList{
				KodeOpd:         kodeOpd,
				PerangkatDaerah: namaOpd,
			})
		}
		result = append(result, pokin)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) ValidateParentLevelTarikStrategiOpd(ctx context.Context, tx *sql.Tx, parentId int, childLevel int) error {
	// Jika tidak ada parent (parent = 0), tidak perlu validasi
	if parentId == 0 {
		return nil
	}

	// Ambil data parent
	query := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
	var parentLevel int
	err := tx.QueryRowContext(ctx, query, parentId).Scan(&parentLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parent dengan id %d tidak ditemukan", parentId)
		}
		return err
	}

	// Validasi khusus untuk level 4 (Strategic)
	if childLevel == 4 {
		// Level 4 bisa memiliki parent dari level 0,1,2,3
		if parentLevel < 0 || parentLevel > 3 {
			return fmt.Errorf("untuk level Strategic (4), parent harus memiliki level 0-3, parent level saat ini: %d", parentLevel)
		}
		return nil
	}

	// Untuk level lainnya, parent harus memiliki level tepat 1 di atasnya
	expectedParentLevel := childLevel - 1
	if parentLevel != expectedParentLevel {
		return fmt.Errorf("level parent (%d) tidak sesuai dengan yang diharapkan (%d) untuk child level %d",
			parentLevel, expectedParentLevel, childLevel)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAtasan(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, []domain.PelaksanaPokin, error) {
	scriptPokin := `
        SELECT parent 
        FROM tb_pohon_kinerja 
        WHERE id = ?`

	var parentId int
	err := tx.QueryRowContext(ctx, scriptPokin, id).Scan(&parentId)
	if err != nil {
		return domain.PohonKinerja{}, nil, err
	}

	// Ambil data pokin parent
	scriptParentPokin := `
        SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, 
               kode_opd, keterangan, tahun, status, is_active
        FROM tb_pohon_kinerja 
        WHERE id = ?`

	var pokinAtasan domain.PohonKinerja
	err = tx.QueryRowContext(ctx, scriptParentPokin, parentId).Scan(
		&pokinAtasan.Id,
		&pokinAtasan.NamaPohon,
		&pokinAtasan.Parent,
		&pokinAtasan.JenisPohon,
		&pokinAtasan.LevelPohon,
		&pokinAtasan.KodeOpd,
		&pokinAtasan.Keterangan,
		&pokinAtasan.Tahun,
		&pokinAtasan.Status,
		&pokinAtasan.IsActive,
	)
	if err != nil {
		return domain.PohonKinerja{}, nil, err
	}

	// Ambil data pegawai dari pelaksana pokin
	scriptPegawai := `
        SELECT DISTINCT 
            p.id as pegawai_id,
            p.nip as nip_pegawai,
            p.nama
        FROM tb_pelaksana_pokin pp
        JOIN tb_pegawai p ON pp.pegawai_id = p.id
        WHERE pp.pohon_kinerja_id = ?`

	rows, err := tx.QueryContext(ctx, scriptPegawai, parentId)
	if err != nil {
		return domain.PohonKinerja{}, nil, err
	}
	defer rows.Close()

	var pegawaiList []domain.PelaksanaPokin
	for rows.Next() {
		var pegawai domain.PelaksanaPokin
		err := rows.Scan(
			&pegawai.Id,
			&pegawai.Nip,
			&pegawai.NamaPegawai,
		)
		if err != nil {
			return domain.PohonKinerja{}, nil, err
		}
		pegawaiList = append(pegawaiList, pegawai)
	}

	return pokinAtasan, pegawaiList, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdateTagging(ctx context.Context, tx *sql.Tx, pokinId int, taggings []domain.TaggingPokin) ([]domain.TaggingPokin, error) {
	// Hapus tagging yang tidak ada di request baru
	existingTaggings, err := repository.FindTaggingByPokinId(ctx, tx, pokinId)
	if err != nil {
		return nil, err
	}

	// Buat map untuk tracking ID yang masih digunakan
	newTaggingIds := make(map[int]bool)
	for _, tagging := range taggings {
		if tagging.Id != 0 {
			newTaggingIds[tagging.Id] = true
		}
	}

	// Hapus tagging yang tidak ada dalam request
	for _, existingTagging := range existingTaggings {
		if !newTaggingIds[existingTagging.Id] {
			// Hapus keterangan tagging program terlebih dahulu
			scriptDeleteKeterangan := "DELETE FROM tb_keterangan_tagging_program_unggulan WHERE id_tagging = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteKeterangan, existingTagging.Id)
			if err != nil {
				return nil, err
			}

			// Kemudian hapus tagging
			scriptDeleteTagging := "DELETE FROM tb_tagging_pokin WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteTagging, existingTagging.Id)
			if err != nil {
				return nil, err
			}
		}
	}

	var results []domain.TaggingPokin

	// Update atau insert tagging
	for _, tagging := range taggings {
		if tagging.Id != 0 {
			// Update existing tagging
			scriptUpdateTagging := "UPDATE tb_tagging_pokin SET nama_tagging = ? WHERE id = ? AND id_pokin = ?"
			_, err := tx.ExecContext(ctx, scriptUpdateTagging,
				tagging.NamaTagging,
				tagging.Id,
				pokinId)
			if err != nil {
				return nil, err
			}

			// Hapus keterangan lama
			scriptDeleteKeterangan := "DELETE FROM tb_keterangan_tagging_program_unggulan WHERE id_tagging = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteKeterangan, tagging.Id)
			if err != nil {
				return nil, err
			}

			// Insert keterangan baru dengan tahun
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				scriptInsertKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
				_, err = tx.ExecContext(ctx, scriptInsertKeterangan,
					tagging.Id,
					keterangan.KodeProgramUnggulan,
					keterangan.Tahun) // Pastikan tahun dimasukkan
				if err != nil {
					return nil, err
				}
			}

			results = append(results, tagging)
		} else {
			// Insert new tagging
			scriptInsertTagging := "INSERT INTO tb_tagging_pokin (id_pokin, nama_tagging) VALUES (?, ?)"
			result, err := tx.ExecContext(ctx, scriptInsertTagging,
				pokinId,
				tagging.NamaTagging)
			if err != nil {
				return nil, err
			}

			newId, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}

			// Insert keterangan untuk tagging baru dengan tahun
			for _, keterangan := range tagging.KeteranganTaggingProgram {
				scriptInsertKeterangan := "INSERT INTO tb_keterangan_tagging_program_unggulan (id_tagging, kode_program_unggulan, tahun) VALUES (?, ?, ?)"
				_, err = tx.ExecContext(ctx, scriptInsertKeterangan,
					newId,
					keterangan.KodeProgramUnggulan,
					keterangan.Tahun) // Pastikan tahun dimasukkan
				if err != nil {
					return nil, err
				}
			}

			tagging.Id = int(newId)
			tagging.IdPokin = pokinId
			results = append(results, tagging)
		}
	}

	return results, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTaggingByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.TaggingPokin, error) {
	// Query untuk mengambil tagging dan keterangan program
	script := `
        SELECT 
            t.id,
            t.id_pokin,
            t.nama_tagging,
            t.clone_from,
            k.id as keterangan_id,
            k.kode_program_unggulan,
            k.tahun
        FROM tb_tagging_pokin t
        LEFT JOIN tb_keterangan_tagging_program_unggulan k ON t.id = k.id_tagging
        WHERE t.id_pokin = ?
        ORDER BY t.id, k.id`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan tagging
	taggingMap := make(map[int]*domain.TaggingPokin)

	for rows.Next() {
		var (
			taggingId, idPokin  int
			namaTagging         string
			cloneFrom           sql.NullInt64
			keteranganId        sql.NullInt64
			kodeProgramUnggulan sql.NullString
			tahun               sql.NullString
		)

		err := rows.Scan(
			&taggingId,
			&idPokin,
			&namaTagging,
			&cloneFrom,
			&keteranganId,
			&kodeProgramUnggulan,
			&tahun,
		)
		if err != nil {
			return nil, err
		}

		// Ambil atau buat tagging baru
		tagging, exists := taggingMap[taggingId]
		if !exists {
			tagging = &domain.TaggingPokin{
				Id:                       taggingId,
				IdPokin:                  idPokin,
				NamaTagging:              namaTagging,
				KeteranganTaggingProgram: []domain.KeteranganTagging{},
			}
			if cloneFrom.Valid {
				tagging.CloneFrom = int(cloneFrom.Int64)
			}
			taggingMap[taggingId] = tagging
		}

		// Tambahkan keterangan program jika ada
		if keteranganId.Valid && kodeProgramUnggulan.Valid {
			keterangan := domain.KeteranganTagging{
				Id:                  int(keteranganId.Int64),
				IdTagging:           taggingId,
				KodeProgramUnggulan: kodeProgramUnggulan.String,
				Tahun:               tahun.String,
			}
			tagging.KeteranganTaggingProgram = append(tagging.KeteranganTaggingProgram, keterangan)
		}
	}

	// Konversi map ke slice
	var result []domain.TaggingPokin
	for _, tagging := range taggingMap {
		result = append(result, *tagging)
	}

	// Sort berdasarkan ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTematikByCloneFrom(ctx context.Context, tx *sql.Tx, cloneFromId int) (*domain.PohonKinerja, error) {
	script := `
        WITH RECURSIVE parent_tree AS (
            -- Base case: start from the cloned node
            SELECT id, parent, nama_pohon, level_pohon
            FROM tb_pohon_kinerja
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: get parent nodes
            SELECT pk.id, pk.parent, pk.nama_pohon, pk.level_pohon
            FROM tb_pohon_kinerja pk
            INNER JOIN parent_tree pt ON pk.id = pt.parent
            WHERE pk.level_pohon >= 0
        )
        SELECT id, nama_pohon
        FROM parent_tree
        WHERE level_pohon = 0
        LIMIT 1`

	var tematik struct {
		Id        int
		NamaPohon string
	}

	err := tx.QueryRowContext(ctx, script, cloneFromId).Scan(&tematik.Id, &tematik.NamaPohon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.PohonKinerja{
		Id:        tematik.Id,
		NamaPohon: tematik.NamaPohon,
	}, nil
}

func (repository *PohonKinerjaRepositoryImpl) ClonePokinPemda(ctx context.Context, tx *sql.Tx, sourceId int, targetTahun string) (int64, error) {
	// 1. Ambil data pohon kinerja source
	scriptGetSource := `
		SELECT 
			nama_pohon, parent, jenis_pohon, level_pohon, 
			kode_opd, keterangan, status, is_active
		FROM tb_pohon_kinerja
		WHERE id = ? AND status != 'tarik pokin opd'
	`

	var source struct {
		NamaPohon  string
		Parent     int
		JenisPohon string
		LevelPohon int
		KodeOpd    string
		Keterangan string
		Status     string
		IsActive   bool
	}

	err := tx.QueryRowContext(ctx, scriptGetSource, sourceId).Scan(
		&source.NamaPohon, &source.Parent, &source.JenisPohon, &source.LevelPohon,
		&source.KodeOpd, &source.Keterangan, &source.Status, &source.IsActive,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data source: %w", err)
	}

	// 2. Tentukan status untuk clone
	var newStatus string
	if source.LevelPohon <= 3 {
		newStatus = source.Status
	} else {
		newStatus = "menunggu_disetujui"
	}

	// 3. Insert pohon kinerja baru
	scriptInsert := `
		INSERT INTO tb_pohon_kinerja 
		(nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, is_active, keterangan_tahun_clone, keterangan_clone_dari)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, scriptInsert,
		source.NamaPohon,
		0, // Parent akan diupdate dalam recursive
		source.JenisPohon,
		source.LevelPohon,
		source.KodeOpd,
		source.Keterangan,
		targetTahun,
		newStatus,
		source.IsActive,
		targetTahun,
		sourceId, // ✅ clone_from = 0 (default)
	)
	if err != nil {
		return 0, fmt.Errorf("gagal insert pohon kinerja: %w", err)
	}

	newPokinId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan ID baru: %w", err)
	}

	// 4. Clone indikator dan target
	err = repository.cloneIndikatorAndTarget(ctx, tx, sourceId, newPokinId)
	if err != nil {
		fmt.Printf("Warning: Gagal clone indikator: %v\n", err)
	}

	// 5. Clone pelaksana
	// err = repository.clonePelaksana(ctx, tx, sourceId, newPokinId)
	// if err != nil {
	// 	fmt.Printf("Warning: Gagal clone pelaksana: %v\n", err)
	// }

	return newPokinId, nil
}

func (repository *PohonKinerjaRepositoryImpl) cloneIndikatorAndTarget(ctx context.Context, tx *sql.Tx, sourceId int, newPokinId int64) error {
	// Gunakan CTE untuk 1 query saja - TIDAK ADA NESTED LOOP
	query := `
		WITH source_data AS (
			SELECT 
				i.id as old_indikator_id,
				i.indikator,
				i.tahun,
				t.id as old_target_id,
				t.target,
				t.satuan,
				t.tahun as target_tahun
			FROM tb_indikator i
			LEFT JOIN tb_target t ON t.indikator_id = i.id
			WHERE i.pokin_id = ?
		)
		SELECT * FROM source_data
	`

	rows, err := tx.QueryContext(ctx, query, sourceId)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Collect all data first
	type dataRow struct {
		OldIndikatorId sql.NullString
		Indikator      sql.NullString
		Tahun          sql.NullString
		OldTargetId    sql.NullString
		Target         sql.NullString
		Satuan         sql.NullString
		TargetTahun    sql.NullString
	}

	var allData []dataRow
	for rows.Next() {
		var d dataRow
		if err := rows.Scan(&d.OldIndikatorId, &d.Indikator, &d.Tahun,
			&d.OldTargetId, &d.Target, &d.Satuan, &d.TargetTahun); err == nil {
			allData = append(allData, d)
		}
	}

	// Process all data
	indikatorMap := make(map[string]string)
	for _, d := range allData {
		if !d.OldIndikatorId.Valid {
			continue
		}

		// Insert indikator jika belum ada di map
		if _, exists := indikatorMap[d.OldIndikatorId.String]; !exists {
			newIndikatorId := fmt.Sprintf("IND-POKIN-%s", uuid.New().String()[:8])

			_, err = tx.ExecContext(ctx, `
				INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from)
				VALUES (?, ?, ?, ?, ?)
			`, newIndikatorId, newPokinId, d.Indikator.String, d.Tahun.String, d.OldIndikatorId.String)

			if err == nil {
				indikatorMap[d.OldIndikatorId.String] = newIndikatorId
			}
		}

		// Insert target
		if d.OldTargetId.Valid && d.Target.Valid {
			newTargetId := fmt.Sprintf("TRGT-IND-%s", uuid.New().String()[:8])

			_, _ = tx.ExecContext(ctx, `
				INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, clone_from)
				VALUES (?, ?, ?, ?, ?, ?)
			`, newTargetId, indikatorMap[d.OldIndikatorId.String], d.Target.String,
				d.Satuan.String, d.TargetTahun.String, d.OldTargetId.String)
		}
	}

	return nil
}

// REPLACEMENT UNTUK clonePelaksana - Lebih sederhana
func (repository *PohonKinerjaRepositoryImpl) clonePelaksana(ctx context.Context, tx *sql.Tx, sourceId int, newPokinId int64) error {
	// Gunakan cara sederhana dengan SELECT dan INSERT
	query := `
		INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id)
		SELECT CONCAT('PLKS-', UUID()), ?, pegawai_id
		FROM tb_pelaksana_pokin
		WHERE pohon_kinerja_id = ?
	`

	_, err := tx.ExecContext(ctx, query, newPokinId, sourceId)
	if err != nil {
		return fmt.Errorf("gagal clone pelaksana: %w", err)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) CloneHierarchyRecursive(ctx context.Context, tx *sql.Tx, sourceId int, newParentId int64, targetTahun string) (int64, error) {
	// Ambil data source
	scriptGetSource := `
		SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, 
		       kode_opd, keterangan, tahun, status, is_active
		FROM tb_pohon_kinerja
		WHERE id = ? AND status != 'tarik pokin opd'
	`

	var source struct {
		Id         int
		NamaPohon  string
		Parent     int
		JenisPohon string
		LevelPohon int
		KodeOpd    string
		Keterangan string
		Tahun      string
		Status     string
		IsActive   bool
	}

	err := tx.QueryRowContext(ctx, scriptGetSource, sourceId).Scan(
		&source.Id, &source.NamaPohon, &source.Parent, &source.JenisPohon, &source.LevelPohon,
		&source.KodeOpd, &source.Keterangan, &source.Tahun, &source.Status, &source.IsActive,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data source: %w", err)
	}

	// Clone pohon kinerja ini
	newId, err := repository.ClonePokinPemda(ctx, tx, sourceId, targetTahun)
	if err != nil {
		return 0, fmt.Errorf("gagal clone pohon kinerja: %w", err)
	}

	// Update parent ID hanya jika newParentId > 0
	if newParentId > 0 {
		_, err = tx.ExecContext(ctx, `
			UPDATE tb_pohon_kinerja SET parent = ? WHERE id = ?
		`, newParentId, newId)
		if err != nil {
			return 0, fmt.Errorf("gagal update parent: %w", err)
		}
	}

	// Ambil semua child dari source
	scriptGetChildren := `
		SELECT id FROM tb_pohon_kinerja 
		WHERE parent = ? AND status != 'tarik pokin opd'
	`
	rows, err := tx.QueryContext(ctx, scriptGetChildren, sourceId)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil child: %w", err)
	}
	defer rows.Close()

	var childIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		childIds = append(childIds, id)
	}

	// ✅ CHECK ERROR DARI ROWS
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating child rows: %w", err)
	}

	// Clone setiap child secara rekursif
	for _, childId := range childIds {
		_, err := repository.CloneHierarchyRecursive(ctx, tx, childId, newId, targetTahun)
		if err != nil {
			return 0, err
		}
	}

	return newId, nil
}

type ControlPokinLevel struct {
	LevelPohon                int
	JumlahPokin               int
	JumlahPelaksana           int
	JumlahPokinAdaPelaksana   int
	JumlahPokinTanpaPelaksana int
	JumlahRencanaKinerja      int
	JumlahPokinAdaRekin       int
	JumlahPokinTanpaRekin     int
}

func (repository *PohonKinerjaRepositoryImpl) ControlPokinOpdByLevel(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) (map[int]ControlPokinLevel, error) {
	query := `
		WITH RECURSIVE valid_pokin AS (
			-- ✅ BASE CASE: Strategic (level 4) dengan parent = 0 atau parent level 0-3
			SELECT
				pk.id,
				pk.level_pohon,
				pk.parent,
				pk.tahun
			FROM tb_pohon_kinerja pk
			WHERE pk.kode_opd = ?
			AND pk.tahun = ?
			AND pk.level_pohon = 4
			AND pk.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			AND (
				pk.parent = 0
				OR EXISTS (
					-- Parent adalah level 0-3 (tematik dari tahun berbeda)
					SELECT 1 FROM tb_pohon_kinerja p2
					WHERE p2.id = pk.parent
					AND p2.level_pohon BETWEEN 0 AND 3
				)
			)

			UNION ALL

			-- ✅ RECURSIVE: Level 5+ harus punya parent valid dengan tahun yang sama
			SELECT
				child.id,
				child.level_pohon,
				child.parent,
				child.tahun
			FROM tb_pohon_kinerja child
			INNER JOIN valid_pokin vp ON child.parent = vp.id
			WHERE child.kode_opd = ?
			AND child.tahun = ?
			AND child.level_pohon > 4
			AND child.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			-- ✅ PENTING: Parent harus tahun yang sama
			AND child.tahun = vp.tahun
		),
		pokin_pelaksana AS (
			SELECT
				vp.level_pohon,
				COUNT(DISTINCT vp.id) as total_pokin,
				COUNT(DISTINCT pp.pegawai_id) as total_pelaksana,
				COUNT(DISTINCT CASE WHEN pp.pegawai_id IS NOT NULL THEN vp.id END) as pokin_ada_pelaksana
			FROM valid_pokin vp
			LEFT JOIN tb_pelaksana_pokin pp ON vp.id = pp.pohon_kinerja_id
			GROUP BY vp.level_pohon
		),
		pokin_rekin AS (
			SELECT
				vp.level_pohon,
				-- Total rencana kinerja yang pegawai-nya adalah pelaksana pohon kinerja
				COUNT(DISTINCT CASE
					WHEN rk.id IS NOT NULL
					AND EXISTS (
						SELECT 1
						FROM tb_pelaksana_pokin pp2
						INNER JOIN tb_pegawai pg ON pp2.pegawai_id = pg.id
						WHERE pp2.pohon_kinerja_id = vp.id
						AND pg.nip = rk.pegawai_id
					)
					THEN rk.id
				END) as total_rencana_kinerja,
				-- Pokin yang punya minimal 1 rencana kinerja (dari pelaksananya)
				COUNT(DISTINCT CASE
					WHEN EXISTS (
						SELECT 1
						FROM tb_rencana_kinerja rk2
						INNER JOIN tb_pelaksana_pokin pp3 ON pp3.pohon_kinerja_id = vp.id
						INNER JOIN tb_pegawai pg2 ON pp3.pegawai_id = pg2.id
						WHERE rk2.id_pohon = vp.id
						AND pg2.nip = rk2.pegawai_id
					)
					THEN vp.id
				END) as pokin_ada_rekin_pelaksana
			FROM valid_pokin vp
			LEFT JOIN tb_rencana_kinerja rk ON vp.id = rk.id_pohon
			GROUP BY vp.level_pohon
		)
		SELECT
			pp.level_pohon,
			pp.total_pokin as jumlah_pokin,
			pp.total_pelaksana as jumlah_pelaksana,
			pp.pokin_ada_pelaksana as jumlah_pokin_ada_pelaksana,
			(pp.total_pokin - pp.pokin_ada_pelaksana) as jumlah_pokin_tanpa_pelaksana,
			COALESCE(pr.total_rencana_kinerja, 0) as jumlah_rencana_kinerja,
			COALESCE(pr.pokin_ada_rekin_pelaksana, 0) as jumlah_pokin_ada_rekin,
			(pp.total_pokin - COALESCE(pr.pokin_ada_rekin_pelaksana, 0)) as jumlah_pokin_tanpa_rekin
		FROM pokin_pelaksana pp
		LEFT JOIN pokin_rekin pr ON pp.level_pohon = pr.level_pohon
		ORDER BY pp.level_pohon
	`

	// ✅ PARAMETER: kodeOpd, tahun (untuk base case), kodeOpd, tahun (untuk recursive)
	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data control pokin: %w", err)
	}
	defer rows.Close()

	result := make(map[int]ControlPokinLevel)
	for rows.Next() {
		var data ControlPokinLevel
		err := rows.Scan(
			&data.LevelPohon,
			&data.JumlahPokin,
			&data.JumlahPelaksana,
			&data.JumlahPokinAdaPelaksana,
			&data.JumlahPokinTanpaPelaksana,
			&data.JumlahRencanaKinerja,
			&data.JumlahPokinAdaRekin,
			&data.JumlahPokinTanpaRekin,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data: %w", err)
		}
		result[data.LevelPohon] = data
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindControlPokinTematikNodes(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]LeaderboardTematikNode, error) {
	query := `
		WITH RECURSIVE
		valid_pokin AS (
			SELECT
				pk.id,
				pk.level_pohon,
				pk.kode_opd,
				pk.clone_from,
				pk.status,
				pk.tahun,
				pk.parent,
				pk.jenis_pohon
			FROM tb_pohon_kinerja pk
			WHERE pk.tahun = ?
			AND pk.kode_opd = ?
			AND pk.level_pohon = 4
			AND pk.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			AND (
				pk.parent = 0
				OR pk.parent IN (
					SELECT id FROM tb_pohon_kinerja
					WHERE level_pohon BETWEEN 0 AND 3
				)
			)
			UNION ALL
			SELECT
				child.id,
				child.level_pohon,
				child.kode_opd,
				child.clone_from,
				child.status,
				child.tahun,
				child.parent,
				child.jenis_pohon
			FROM tb_pohon_kinerja child
			INNER JOIN valid_pokin vp ON child.parent = vp.id
			WHERE child.tahun = ?
			AND child.level_pohon > 4
			AND child.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			AND child.kode_opd = vp.kode_opd
			AND child.tahun = vp.tahun
		),
		tematik_trace AS (
			SELECT
				vp.kode_opd,
				pk_src.id,
				pk_src.parent,
				pk_src.level_pohon,
				pk_src.nama_pohon,
				1 AS depth
			FROM valid_pokin vp
			INNER JOIN tb_pohon_kinerja pk_src ON vp.clone_from = pk_src.id
			WHERE (
				vp.status = 'pokin dari pemda'
				OR (
					vp.status = 'crosscutting_disetujui_existing'
					AND vp.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
				)
			)
			AND vp.clone_from > 0
			AND vp.level_pohon = 4
			UNION ALL
			SELECT
				tt.kode_opd,
				pk_parent.id,
				pk_parent.parent,
				pk_parent.level_pohon,
				pk_parent.nama_pohon,
				tt.depth + 1
			FROM tematik_trace tt
			INNER JOIN tb_pohon_kinerja pk_parent ON tt.parent = pk_parent.id
			WHERE tt.level_pohon > 0
			AND tt.depth < 5
		)
		SELECT DISTINCT
			tt.kode_opd,
			tt.id,
			tt.parent,
			tt.level_pohon,
			tt.nama_pohon
		FROM tematik_trace tt
		WHERE tt.level_pohon BETWEEN 0 AND 3
		ORDER BY tt.level_pohon, tt.id
	`
	rows, err := tx.QueryContext(ctx, query, tahun, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil node tematik control pokin: %w", err)
	}
	defer rows.Close()
	var out []LeaderboardTematikNode
	for rows.Next() {
		var n LeaderboardTematikNode
		if err := rows.Scan(&n.KodeOpd, &n.Id, &n.Parent, &n.LevelPohon, &n.NamaPohon); err != nil {
			return nil, fmt.Errorf("gagal scan node tematik control: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

type LeaderboardOpdData struct {
	KodeOpd             string
	NamaOpd             string
	TotalPokin          int
	TotalPokinAdaRekin  int
	PersentaseCascading float64
	TematikL0           []string
	TematikL1           []string
	TematikL2           []string
	TematikL3           []string
	IsHidden            bool
}

type LeaderboardTematikNode struct {
	KodeOpd    string
	Id         int
	Parent     int
	LevelPohon int
	NamaPohon  string
}

func (repository *PohonKinerjaRepositoryImpl) LeaderboardPokinOpd(ctx context.Context, tx *sql.Tx, tahun string) ([]LeaderboardOpdData, error) {
	query := `
	WITH RECURSIVE 
	valid_pokin AS (
		SELECT 
			pk.id,
			pk.level_pohon,
			pk.kode_opd,
			pk.clone_from,
			pk.status,
			pk.tahun,
			pk.parent,
			pk.jenis_pohon
		FROM tb_pohon_kinerja pk
		WHERE pk.tahun = ?
		AND pk.level_pohon = 4
		AND pk.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
		AND (
			pk.parent = 0 
			OR pk.parent IN (
				SELECT id FROM tb_pohon_kinerja 
				WHERE level_pohon BETWEEN 0 AND 3
			)
		)

		UNION ALL
		
		SELECT 
			child.id,
			child.level_pohon,
			child.kode_opd,
			child.clone_from,
			child.status,
			child.tahun,
			child.parent,
			child.jenis_pohon
		FROM tb_pohon_kinerja child
		INNER JOIN valid_pokin vp ON child.parent = vp.id
		WHERE child.tahun = ?
		AND child.level_pohon > 4
		AND child.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
		AND child.kode_opd = vp.kode_opd
		AND child.tahun = vp.tahun
	),

	pokin_pelaksana_valid AS (
		SELECT DISTINCT
			pp.pohon_kinerja_id,
			pg.nip
		FROM tb_pelaksana_pokin pp
		INNER JOIN tb_pegawai pg ON pp.pegawai_id = pg.id
		WHERE pp.pohon_kinerja_id IN (SELECT id FROM valid_pokin)
	),

	-- ✅ OPTIMASI: ganti EXISTS jadi JOIN
	pokin_with_rekin AS (
		SELECT DISTINCT 
			vp.id,
			vp.kode_opd
		FROM valid_pokin vp
		JOIN pokin_pelaksana_valid ppv 
			ON ppv.pohon_kinerja_id = vp.id
		JOIN tb_rencana_kinerja rk 
			ON rk.id_pohon = vp.id 
			AND rk.pegawai_id = ppv.nip
	),

	opd_cascading AS (
		SELECT 
			vp.kode_opd,
			COUNT(DISTINCT vp.id) as total_pokin,
			COUNT(DISTINCT pr.id) as total_pokin_ada_rekin
		FROM valid_pokin vp
		LEFT JOIN pokin_with_rekin pr ON pr.id = vp.id
		GROUP BY vp.kode_opd
	),

	tematik_trace AS (
		SELECT 
			vp.kode_opd,
			pk_src.id,
			pk_src.parent,
			pk_src.level_pohon,
			pk_src.nama_pohon,
			1 as depth
		FROM valid_pokin vp
		INNER JOIN tb_pohon_kinerja pk_src ON vp.clone_from = pk_src.id
		WHERE (
			vp.status = 'pokin dari pemda'
			OR (
				vp.status = 'crosscutting_disetujui_existing' 
				AND vp.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
			)
		)
		AND vp.clone_from > 0
		AND vp.level_pohon = 4
		
		UNION ALL
		
		SELECT 
			tt.kode_opd,
			pk_parent.id,
			pk_parent.parent,
			pk_parent.level_pohon,
			pk_parent.nama_pohon,
			tt.depth + 1
		FROM tematik_trace tt
		INNER JOIN tb_pohon_kinerja pk_parent ON tt.parent = pk_parent.id
		WHERE tt.level_pohon > 0
		AND tt.depth < 5
	),

	tematik_agregat AS (
		SELECT 
			kode_opd,
			GROUP_CONCAT(DISTINCT CASE WHEN level_pohon = 0 THEN nama_pohon END SEPARATOR '|||') AS tematik_l0,
			GROUP_CONCAT(DISTINCT CASE WHEN level_pohon = 1 THEN nama_pohon END SEPARATOR '|||') AS tematik_l1,
			GROUP_CONCAT(DISTINCT CASE WHEN level_pohon = 2 THEN nama_pohon END SEPARATOR '|||') AS tematik_l2,
			GROUP_CONCAT(DISTINCT CASE WHEN level_pohon = 3 THEN nama_pohon END SEPARATOR '|||') AS tematik_l3,
			COUNT(DISTINCT CASE WHEN level_pohon = 0 AND parent = 0 THEN nama_pohon END) AS jumlah_tematik
		FROM tematik_trace
		WHERE level_pohon BETWEEN 0 AND 3
		GROUP BY kode_opd
	)

		SELECT 
		opd.kode_opd,
		opd.nama_opd,
		COALESCE(oc.total_pokin, 0) AS total_pokin,
		COALESCE(oc.total_pokin_ada_rekin, 0) AS total_pokin_ada_rekin,
		CASE 
			WHEN COALESCE(oc.total_pokin, 0) > 0 
			THEN (COALESCE(oc.total_pokin_ada_rekin, 0) * 100.0 / oc.total_pokin)
			ELSE 0 
		END AS persentase_cascading,
		COALESCE(ta.tematik_l0, '') AS tematik_l0,
		COALESCE(ta.tematik_l1, '') AS tematik_l1,
		COALESCE(ta.tematik_l2, '') AS tematik_l2,
		COALESCE(ta.tematik_l3, '') AS tematik_l3,
		CASE WHEN COALESCE(ta.jumlah_tematik, 0) > 0 THEN 1 ELSE 0 END AS has_tematik,
		COALESCE(lh.is_hidden, 0) AS is_hidden
	FROM tb_operasional_daerah opd
	LEFT JOIN opd_cascading oc ON opd.kode_opd = oc.kode_opd
	LEFT JOIN tematik_agregat ta ON opd.kode_opd = ta.kode_opd
	LEFT JOIN tb_leaderboard_hidden lh 
		ON opd.kode_opd = lh.kode_opd 
		AND lh.tahun = ?

	ORDER BY persentase_cascading DESC, opd.nama_opd ASC
	`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data leaderboard: %w", err)
	}
	defer rows.Close()

	var result []LeaderboardOpdData
	for rows.Next() {
		var data LeaderboardOpdData
		var tL0, tL1, tL2, tL3 string
		var hasTematik int

		err := rows.Scan(
			&data.KodeOpd,
			&data.NamaOpd,
			&data.TotalPokin,
			&data.TotalPokinAdaRekin,
			&data.PersentaseCascading,
			&tL0, &tL1, &tL2, &tL3,
			&hasTematik,
			&data.IsHidden,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data: %w", err)
		}

		splitNames := func(s string) []string {
			if s == "" {
				return []string{}
			}
			return strings.Split(s, "|||")
		}

		data.TematikL0 = splitNames(tL0)
		data.TematikL1 = splitNames(tL1)
		data.TematikL2 = splitNames(tL2)
		data.TematikL3 = splitNames(tL3)

		result = append(result, data)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindLeaderboardTematikNodes(ctx context.Context, tx *sql.Tx, tahun string) ([]LeaderboardTematikNode, error) {
	query := `
		WITH RECURSIVE
		valid_pokin AS (
			SELECT
				pk.id,
				pk.level_pohon,
				pk.kode_opd,
				pk.clone_from,
				pk.status,
				pk.tahun,
				pk.parent,
				pk.jenis_pohon
			FROM tb_pohon_kinerja pk
			WHERE pk.tahun = ?
			AND pk.level_pohon = 4
			AND pk.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			AND (
				pk.parent = 0
				OR pk.parent IN (
					SELECT id FROM tb_pohon_kinerja
					WHERE level_pohon BETWEEN 0 AND 3
				)
			)
			UNION ALL
			SELECT
				child.id,
				child.level_pohon,
				child.kode_opd,
				child.clone_from,
				child.status,
				child.tahun,
				child.parent,
				child.jenis_pohon
			FROM tb_pohon_kinerja child
			INNER JOIN valid_pokin vp ON child.parent = vp.id
			WHERE child.tahun = ?
			AND child.level_pohon > 4
			AND child.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
			AND child.kode_opd = vp.kode_opd
			AND child.tahun = vp.tahun
		),
		tematik_trace AS (
			SELECT
				vp.kode_opd,
				pk_src.id,
				pk_src.parent,
				pk_src.level_pohon,
				pk_src.nama_pohon,
				1 AS depth
			FROM valid_pokin vp
			INNER JOIN tb_pohon_kinerja pk_src ON vp.clone_from = pk_src.id
			WHERE (
				vp.status = 'pokin dari pemda'
				OR (
					vp.status = 'crosscutting_disetujui_existing'
					AND vp.jenis_pohon IN ('Strategic Pemda', 'Tactical Pemda', 'Operasional Pemda')
				)
			)
			AND vp.clone_from > 0
			AND vp.level_pohon = 4
			UNION ALL
			SELECT
				tt.kode_opd,
				pk_parent.id,
				pk_parent.parent,
				pk_parent.level_pohon,
				pk_parent.nama_pohon,
				tt.depth + 1
			FROM tematik_trace tt
			INNER JOIN tb_pohon_kinerja pk_parent ON tt.parent = pk_parent.id
			WHERE tt.level_pohon > 0
			AND tt.depth < 5
		)
		SELECT DISTINCT
			tt.kode_opd,
			tt.id,
			tt.parent,
			tt.level_pohon,
			tt.nama_pohon
		FROM tematik_trace tt
		WHERE tt.level_pohon BETWEEN 0 AND 3
		ORDER BY tt.kode_opd, tt.level_pohon, tt.id
	`
	rows, err := tx.QueryContext(ctx, query, tahun, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil node tematik leaderboard: %w", err)
	}
	defer rows.Close()
	var out []LeaderboardTematikNode
	for rows.Next() {
		var n LeaderboardTematikNode
		if err := rows.Scan(&n.KodeOpd, &n.Id, &n.Parent, &n.LevelPohon, &n.NamaPohon); err != nil {
			return nil, fmt.Errorf("gagal scan node tematik: %w", err)
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpsertLeaderboardHidden(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, hidden bool) error {
	// Satu pasangan (kode_opd, tahun) hanya boleh 0 atau 1 baris; DELETE membersihkan duplikat lama.
	const del = `DELETE FROM tb_leaderboard_hidden WHERE kode_opd = ? AND tahun = ?`
	if _, err := tx.ExecContext(ctx, del, kodeOpd, tahun); err != nil {
		return fmt.Errorf("upsert leaderboard hidden: %w", err)
	}
	if !hidden {
		return nil
	}
	const ins = `INSERT INTO tb_leaderboard_hidden (kode_opd, tahun, is_hidden) VALUES (?, ?, TRUE)`
	if _, err := tx.ExecContext(ctx, ins, kodeOpd, tahun); err != nil {
		return fmt.Errorf("upsert leaderboard hidden: %w", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindLeaderboardHiddenKodeOpdsByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]string, error) {
	const q = `
SELECT kode_opd
FROM tb_leaderboard_hidden
WHERE tahun = ? AND is_hidden = TRUE
ORDER BY kode_opd
`
	rows, err := tx.QueryContext(ctx, q, tahun)
	if err != nil {
		return nil, fmt.Errorf("daftar leaderboard hidden: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, fmt.Errorf("scan kode_opd hidden: %w", err)
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func (repository *PohonKinerjaRepositoryImpl) FindPelaksanaPokinBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.PelaksanaPokin, error) {
	if len(pokinIds) == 0 {
		return make(map[int][]domain.PelaksanaPokin), nil
	}

	// Build IN clause dengan placeholders
	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT 
			pp.id, 
			pp.pohon_kinerja_id, 
			pp.pegawai_id,
			p.nip,
			p.nama as nama_pegawai
		FROM tb_pelaksana_pokin pp
		INNER JOIN tb_pegawai p ON pp.pegawai_id = p.id
		WHERE pp.pohon_kinerja_id IN (%s)
		ORDER BY pp.pohon_kinerja_id, pp.id
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]domain.PelaksanaPokin)
	for rows.Next() {
		var pelaksana domain.PelaksanaPokin
		var pokinId int
		err := rows.Scan(
			&pelaksana.Id,
			&pokinId,
			&pelaksana.PegawaiId,
			&pelaksana.Nip,
			&pelaksana.NamaPegawai,
		)
		if err != nil {
			return nil, err
		}
		result[pokinId] = append(result[pokinId], pelaksana)
	}

	return result, nil
}

// FindIndikatorByPokinIdsBatch mengambil semua indikator untuk multiple pokin dalam 1 query
func (repository *PohonKinerjaRepositoryImpl) FindIndikatorByPokinIdsBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.Indikator, error) {
	if len(pokinIds) == 0 {
		return make(map[int][]domain.Indikator), nil
	}

	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT 
			i.id, 
			i.pokin_id, 
			i.indikator, 
			i.tahun, 
			i.clone_from
		FROM tb_indikator i
		WHERE i.pokin_id IN (%s)
		ORDER BY i.pokin_id, i.id
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]domain.Indikator)
	for rows.Next() {
		var indikator domain.Indikator
		var pokinId int
		err := rows.Scan(
			&indikator.Id,
			&pokinId,
			&indikator.Indikator,
			&indikator.Tahun,
			&indikator.CloneFrom,
		)
		if err != nil {
			return nil, err
		}
		result[pokinId] = append(result[pokinId], indikator)
	}

	return result, nil
}

// FindTargetByIndikatorIdsBatch mengambil semua target untuk multiple indikator dalam 1 query
func (repository *PohonKinerjaRepositoryImpl) FindTargetByIndikatorIdsBatch(ctx context.Context, tx *sql.Tx, indikatorIds []string) (map[string][]domain.Target, error) {
	if len(indikatorIds) == 0 {
		return make(map[string][]domain.Target), nil
	}

	placeholders := make([]string, len(indikatorIds))
	args := make([]interface{}, len(indikatorIds))
	for i, id := range indikatorIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT 
			id, 
			indikator_id, 
			target, 
			satuan, 
			tahun
		FROM tb_target
		WHERE indikator_id IN (%s)
		ORDER BY indikator_id, id
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]domain.Target)
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(
			&target.Id,
			&target.IndikatorId,
			&target.Target,
			&target.Satuan,
			&target.Tahun,
		)
		if err != nil {
			return nil, err
		}
		result[target.IndikatorId] = append(result[target.IndikatorId], target)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTaggingByPokinIdsBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) (map[int][]domain.TaggingPokin, error) {
	if len(pokinIds) == 0 {
		return make(map[int][]domain.TaggingPokin), nil
	}

	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	// OPTIMASI: JOIN langsung dengan program unggulan untuk menghindari query terpisah
	script := fmt.Sprintf(`
		SELECT 
			t.id,
			t.id_pokin,
			t.nama_tagging,
			t.clone_from,
			k.id as keterangan_id,
			k.kode_program_unggulan,
			k.tahun,
			pu.nama_tagging AS nama_program_prioritas,
			pu.keterangan_program_unggulan
		FROM tb_tagging_pokin t
		LEFT JOIN tb_keterangan_tagging_program_unggulan k ON t.id = k.id_tagging
		LEFT JOIN tb_program_unggulan pu ON k.kode_program_unggulan = pu.kode_program_unggulan
		WHERE t.id_pokin IN (%s)
		ORDER BY t.id_pokin, t.id, k.id
		LIMIT 10000
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taggingMap := make(map[int]*domain.TaggingPokin)
	result := make(map[int][]domain.TaggingPokin, len(pokinIds))

	for rows.Next() {
		var (
			taggingId, idPokin   int
			namaTagging          string
			cloneFrom            sql.NullInt64
			keteranganId         sql.NullInt64
			kodeProgramUnggulan  sql.NullString
			tahun                sql.NullString
			namaProgramPrioritas sql.NullString
			rencanaImplementasi  sql.NullString
		)

		err := rows.Scan(
			&taggingId,
			&idPokin,
			&namaTagging,
			&cloneFrom,
			&keteranganId,
			&kodeProgramUnggulan,
			&tahun,
			&namaProgramPrioritas,
			&rencanaImplementasi,
		)
		if err != nil {
			return nil, err
		}

		key := taggingId
		tagging, exists := taggingMap[key]
		if !exists {
			tagging = &domain.TaggingPokin{
				Id:                       taggingId,
				IdPokin:                  idPokin,
				NamaTagging:              namaTagging,
				KeteranganTaggingProgram: []domain.KeteranganTagging{},
			}
			if cloneFrom.Valid {
				tagging.CloneFrom = int(cloneFrom.Int64)
			}
			taggingMap[key] = tagging
		}

		if keteranganId.Valid && kodeProgramUnggulan.Valid {
			var rencanaImpl *string
			if rencanaImplementasi.Valid && rencanaImplementasi.String != "" {
				rencanaImpl = &rencanaImplementasi.String
			}
			var namaPrgPrio *string
			if namaProgramPrioritas.Valid && namaProgramPrioritas.String != "" {
				namaPrgPrio = &namaProgramPrioritas.String
			}

			keterangan := domain.KeteranganTagging{
				Id:                   int(keteranganId.Int64),
				IdTagging:            taggingId,
				KodeProgramUnggulan:  kodeProgramUnggulan.String,
				NamaProgramPrioritas: namaPrgPrio,
				RencanaImplementasi:  rencanaImpl,
				Tahun:                tahun.String,
			}
			tagging.KeteranganTaggingProgram = append(tagging.KeteranganTaggingProgram, keterangan)
		}
	}

	// Group by pokinId
	for _, tagging := range taggingMap {
		pokinId := tagging.IdPokin
		if result[pokinId] == nil {
			result[pokinId] = make([]domain.TaggingPokin, 0, 2)
		}
		result[pokinId] = append(result[pokinId], *tagging)
	}

	return result, nil
}

// func (repository *PohonKinerjaRepositoryImpl) FindTematikByCloneFromBatch(ctx context.Context, tx *sql.Tx, cloneFromIds []int) (map[int]*domain.PohonKinerja, error) {
// 	if len(cloneFromIds) == 0 {
// 		return make(map[int]*domain.PohonKinerja), nil
// 	}

// 	uniqueIds := make(map[int]bool)
// 	var uniqueCloneFromIds []int
// 	for _, id := range cloneFromIds {
// 		if id > 0 && !uniqueIds[id] {
// 			uniqueIds[id] = true
// 			uniqueCloneFromIds = append(uniqueCloneFromIds, id)
// 		}
// 	}

// 	if len(uniqueCloneFromIds) == 0 {
// 		return make(map[int]*domain.PohonKinerja), nil
// 	}

// 	result := make(map[int]*domain.PohonKinerja)

// 	// OPTIMASI: Ambil semua data yang relevan dalam query sederhana
// 	// Kumpulkan semua IDs yang mungkin relevan (clone_from_ids + semua parent-nya)
// 	nodeMap := make(map[int]*struct {
// 		Id         int
// 		Parent     sql.NullInt64
// 		NamaPohon  string
// 		LevelPohon int
// 	})

// 	// Step 1: Ambil clone_from_ids dan parent-nya secara iterative
// 	currentIds := uniqueCloneFromIds
// 	maxIterations := 10

// 	for iteration := 0; iteration < maxIterations && len(currentIds) > 0; iteration++ {
// 		placeholders := make([]string, len(currentIds))
// 		args := make([]interface{}, len(currentIds))
// 		for i, id := range currentIds {
// 			placeholders[i] = "?"
// 			args[i] = id
// 		}

// 		script := fmt.Sprintf(`
// 			SELECT id, parent, nama_pohon, level_pohon
// 			FROM tb_pohon_kinerja
// 			WHERE id IN (%s)
// 		`, strings.Join(placeholders, ","))

// 		rows, err := tx.QueryContext(ctx, script, args...)
// 		if err != nil {
// 			log.Printf("[ERROR] FindTematikByCloneFromBatch query error: %v", err)
// 			break
// 		}

// 		var nextIds []int
// 		nextIdSet := make(map[int]bool)

// 		for rows.Next() {
// 			var id, levelPohon int
// 			var parent sql.NullInt64
// 			var namaPohon string

// 			if err := rows.Scan(&id, &parent, &namaPohon, &levelPohon); err != nil {
// 				continue
// 			}

// 			nodeMap[id] = &struct {
// 				Id         int
// 				Parent     sql.NullInt64
// 				NamaPohon  string
// 				LevelPohon int
// 			}{
// 				Id:         id,
// 				Parent:     parent,
// 				NamaPohon:  namaPohon,
// 				LevelPohon: levelPohon,
// 			}

// 			// PERBAIKAN: Convert int64 ke int
// 			if parent.Valid && parent.Int64 > 0 {
// 				parentId := int(parent.Int64)
// 				if !nextIdSet[parentId] {
// 					nextIds = append(nextIds, parentId)
// 					nextIdSet[parentId] = true
// 				}
// 			}
// 		}
// 		rows.Close()

// 		currentIds = nextIds
// 	}

// 	// Step 2: Process di Go untuk mencari tematik
// 	for _, cloneFromId := range uniqueCloneFromIds {
// 		tematik := repository.findTematikInMemorySimple(nodeMap, cloneFromId)
// 		if tematik != nil {
// 			result[cloneFromId] = tematik
// 		}
// 	}

// 	log.Printf("[DEBUG] FindTematikByCloneFromBatch: Total results=%d out of %d cloneFromIds",
// 		len(result), len(uniqueCloneFromIds))

// 	return result, nil
// }

// func (repository *PohonKinerjaRepositoryImpl) findTematikInMemorySimple(
// 	nodeMap map[int]*struct {
// 		Id         int
// 		Parent     sql.NullInt64
// 		NamaPohon  string
// 		LevelPohon int
// 	},
// 	startId int,
// ) *domain.PohonKinerja {
// 	currentId := startId
// 	visited := make(map[int]bool)
// 	maxDepth := 15

// 	for depth := 0; depth < maxDepth; depth++ {
// 		if visited[currentId] {
// 			break
// 		}
// 		visited[currentId] = true

// 		node, exists := nodeMap[currentId]
// 		if !exists {
// 			break
// 		}

// 		// Jika ketemu level 0 dengan parent NULL atau 0, return
// 		if node.LevelPohon == 0 && (!node.Parent.Valid || node.Parent.Int64 == 0) {
// 			return &domain.PohonKinerja{
// 				Id:        node.Id,
// 				NamaPohon: node.NamaPohon,
// 			}
// 		}

// 		// Naik ke parent dengan convert int64 ke int
// 		if node.Parent.Valid && node.Parent.Int64 > 0 {
// 			currentId = int(node.Parent.Int64)
// 		} else {
// 			break
// 		}
// 	}

// 	return nil
// }

func (repository *PohonKinerjaRepositoryImpl) FindTematikByCloneFromBatch(ctx context.Context, tx *sql.Tx, cloneFromIds []int) (map[int]*domain.PohonKinerja, error) {
	if len(cloneFromIds) == 0 {
		return make(map[int]*domain.PohonKinerja), nil
	}

	// Batasi maksimal 50 IDs untuk performa
	if len(cloneFromIds) > 50 {
		cloneFromIds = cloneFromIds[:50]
	}

	placeholders := make([]string, len(cloneFromIds))
	args := make([]interface{}, len(cloneFromIds))
	for i, id := range cloneFromIds {
		placeholders[i] = "?"
		args[i] = id
	}

	// Query optimized dengan CTE yang lebih efisien
	script := fmt.Sprintf(`
		WITH RECURSIVE tematik_tree AS (
			SELECT 
				id,
				parent,
				nama_pohon,
				level_pohon,
				id as original_id,
				0 as depth
			FROM tb_pohon_kinerja
			WHERE id IN (%s)
			
			UNION ALL
			
			SELECT 
				p.id,
				p.parent,
				p.nama_pohon,
				p.level_pohon,
				t.original_id,
				t.depth + 1
			FROM tb_pohon_kinerja p
			INNER JOIN tematik_tree t ON p.id = t.parent
			WHERE t.depth < 8
				AND p.level_pohon >= 0
		)
		SELECT 
			original_id,
			id as tematik_id,
			nama_pohon
		FROM (
			SELECT 
				original_id,
				id,
				nama_pohon,
				level_pohon,
				ROW_NUMBER() OVER (PARTITION BY original_id ORDER BY depth ASC) as rn
			FROM tematik_tree
			WHERE level_pohon = 0
		) ranked
		WHERE rn = 1
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return make(map[int]*domain.PohonKinerja), nil // Return empty map pada error
	}
	defer rows.Close()

	result := make(map[int]*domain.PohonKinerja, len(cloneFromIds))
	for rows.Next() {
		var originalId, tematikId int
		var namaPohon string
		if err := rows.Scan(&originalId, &tematikId, &namaPohon); err != nil {
			continue
		}
		result[originalId] = &domain.PohonKinerja{
			Id:        tematikId,
			NamaPohon: namaPohon,
		}
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindByIds(ctx context.Context, tx *sql.Tx, ids []int) (map[int]domain.PohonKinerja, error) {
	if len(ids) == 0 {
		return make(map[int]domain.PohonKinerja), nil
	}

	// Buat placeholders untuk IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT id, nama_pohon, tahun, level_pohon 
		FROM tb_pohon_kinerja 
		WHERE id IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pohonMap := make(map[int]domain.PohonKinerja)
	for rows.Next() {
		var pohon domain.PohonKinerja
		err := rows.Scan(&pohon.Id, &pohon.NamaPohon, &pohon.Tahun, &pohon.LevelPohon)
		if err != nil {
			return nil, err
		}
		pohonMap[pohon.Id] = pohon
	}

	return pohonMap, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPelaksanaPokinBatchForCascading(
	ctx context.Context,
	tx *sql.Tx,
	pohonKinerjaIds []int,
) ([]domain.PelaksanaPokin, error) {
	const op = "pohonkinerja_repository.FindPelaksanaPokinBatch"

	if len(pohonKinerjaIds) == 0 {
		return []domain.PelaksanaPokin{}, nil
	}

	baseQuery := `
		SELECT tpokin.id, tpokin.pohon_kinerja_id, tpokin.pegawai_id, pg.nama, pg.nip
		FROM tb_pelaksana_pokin tpokin
        JOIN tb_pegawai pg ON tpokin.pegawai_id = pg.id
		WHERE tpokin.pohon_kinerja_id IN (?)
	`

	query, args := helper.BuildInQuery(baseQuery, pohonKinerjaIds)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var result []domain.PelaksanaPokin
	for rows.Next() {
		var pelaksana domain.PelaksanaPokin
		if err := rows.Scan(
			&pelaksana.Id,
			&pelaksana.PohonKinerjaId,
			&pelaksana.PegawaiId,
			&pelaksana.NamaPegawai,
			&pelaksana.Nip,
		); err != nil {
			return nil, fmt.Errorf("%s: query failed: %w", op, err)
		}
		result = append(result, pelaksana)
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckIfSourceAlreadyCloned(
	ctx context.Context,
	tx *sql.Tx,
	sourceId int,
	tahunTarget string,
) (bool, error) {

	script := `
        SELECT 1
        FROM tb_pohon_kinerja pk
        WHERE pk.keterangan_clone_dari = ? AND pk.keterangan_tahun_clone = ?
        LIMIT 1
    `

	var exists int
	err := tx.QueryRowContext(ctx, script, sourceId, tahunTarget).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByParentClonePokinOpd(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string, levelPohon *int) ([]domain.PohonKinerja, error) {
	const parentPlaceholder = -100
	script := `
		SELECT 
			id,
			COALESCE(nama_pohon, '') AS nama_pohon,
			COALESCE(parent, 0) AS parent,
			COALESCE(jenis_pohon, '') AS jenis_pohon,
			COALESCE(level_pohon, 0) AS level_pohon,
			COALESCE(kode_opd, '') AS kode_opd,
			COALESCE(keterangan, '') AS keterangan,
			COALESCE(keterangan_crosscutting, '') AS keterangan_crosscutting,
			COALESCE(tahun, '') AS tahun,
			COALESCE(status, '') AS status,
			COALESCE(is_active, 0) AS is_active,
			COALESCE(clone_from, 0) AS clone_from,
			COALESCE(keterangan_tahun_clone, '') AS keterangan_tahun_clone
		FROM tb_pohon_kinerja
		WHERE kode_opd = ?
		  AND tahun = ?
		  AND parent = ?
		  AND status = ''
		  AND jenis_pohon IN ('Strategic', 'Tactical', 'Operational', 'Operational N')
	`
	args := []interface{}{kodeOpd, tahun, parentPlaceholder}
	if levelPohon != nil {
		script += ` AND level_pohon = ?`
		args = append(args, *levelPohon)
	}
	script += ` ORDER BY level_pohon ASC, id ASC`
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.PohonKinerja
	for rows.Next() {
		var p domain.PohonKinerja
		err := rows.Scan(
			&p.Id,
			&p.NamaPohon,
			&p.Parent,
			&p.JenisPohon,
			&p.LevelPohon,
			&p.KodeOpd,
			&p.Keterangan,
			&p.KeteranganCrosscutting,
			&p.Tahun,
			&p.Status,
			&p.IsActive,
			&p.CloneFrom,
			&p.KeteranganTahunClone,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
