package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type CascadingOpdRepositoryImpl struct {
}

func NewCascadingOpdRepositoryImpl(db *sql.DB, rencanaKinerjaRepository RencanaKinerjaRepository) *CascadingOpdRepositoryImpl {
	return &CascadingOpdRepositoryImpl{}
}

func (repository *CascadingOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
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
            COALESCE(is_active) as is_active
        FROM tb_pohon_kinerja 
        WHERE kode_opd = ? 
        AND tahun = ?
        AND status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
        ORDER BY 
            CASE 
                WHEN status = 'pokin dari pemda' THEN 0 
                ELSE 1 
            END,
            level_pohon, 
            id ASC`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
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
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}

	// Inisialisasi slice kosong jika tidak ada data
	if pokins == nil {
		pokins = make([]domain.PohonKinerja, 0)
	}

	return pokins, nil
}

func (repository *CascadingOpdRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error) {
	script := `
        SELECT i.id, i.pokin_id, i.indikator, 
               t.id, t.indikator_id, t.target, t.satuan
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
		var indId, pokinId, indikator string
		var targetId, indikatorId, target, satuan sql.NullString

		err := rows.Scan(
			&indId, &pokinId, &indikator,
			&targetId, &indikatorId, &target, &satuan)
		if err != nil {
			return nil, err
		}

		// Proses Indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:        indId,
				Indikator: indikator,
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

func (repository *CascadingOpdRepositoryImpl) FindByKodeAndOpdAndTahun(ctx context.Context, tx *sql.Tx, kode string, kodeOpd string, tahun string) ([]domain.Indikator, error) {
	query := `
        SELECT 
            id,
            kode,
            indikator
        FROM tb_indikator 
        WHERE kode = ? 
        AND kode_opd = ?
        AND tahun = ?
    `

	rows, err := tx.QueryContext(ctx, query, kode, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var i domain.Indikator
		err := rows.Scan(
			&i.Id,
			&i.Kode,
			&i.Indikator,
		)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, i)
	}

	return indikators, nil
}

func (repository *CascadingOpdRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
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

func (repository *CascadingOpdRepositoryImpl) FindPokinByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (domain.PohonKinerja, error) {
	script := `
        SELECT 
            pk.id,
            COALESCE(pk.nama_pohon, '') as nama_pohon,
            COALESCE(pk.parent, 0) as parent,
            COALESCE(pk.jenis_pohon, '') as jenis_pohon,
            COALESCE(pk.level_pohon, 0) as level_pohon,
            COALESCE(pk.kode_opd, '') as kode_opd,
            COALESCE(pk.keterangan, '') as keterangan,
            COALESCE(pk.keterangan_crosscutting, '') as keterangan_crosscutting,
            COALESCE(pk.tahun, '') as tahun,
            COALESCE(pk.status, '') as status,
            COALESCE(pk.is_active) as is_active
        FROM tb_pohon_kinerja pk
        INNER JOIN tb_rencana_kinerja rk ON rk.id_pohon = pk.id
        WHERE rk.id = ?
    `

	var pokin domain.PohonKinerja
	err := tx.QueryRowContext(ctx, script, rekinId).Scan(
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
	)
	if err != nil {
		return domain.PohonKinerja{}, err
	}

	return pokin, nil
}

func (repository *CascadingOpdRepositoryImpl) FindPokinById(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error) {
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
            COALESCE(is_active) as is_active
        FROM tb_pohon_kinerja 
        WHERE id = ?
    `

	var pokin domain.PohonKinerja
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(
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
	)
	if err != nil {
		return domain.PohonKinerja{}, err
	}

	return pokin, nil
}

func (repository *CascadingOpdRepositoryImpl) FindStrategicByChildPokin(ctx context.Context, tx *sql.Tx, pokinId int) (domain.PohonKinerja, error) {
	// Recursive query untuk trace ke level 4 (Strategic)
	script := `
        WITH RECURSIVE parent_tree AS (
            -- Base case: pohon yang dicari
            SELECT id, parent, level_pohon, nama_pohon, jenis_pohon, kode_opd, 
                   keterangan, keterangan_crosscutting, tahun, status, is_active
            FROM tb_pohon_kinerja
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: parent dari pohon
            SELECT pk.id, pk.parent, pk.level_pohon, pk.nama_pohon, pk.jenis_pohon, 
                   pk.kode_opd, pk.keterangan, pk.keterangan_crosscutting, pk.tahun, 
                   pk.status, pk.is_active
            FROM tb_pohon_kinerja pk
            INNER JOIN parent_tree pt ON pk.id = pt.parent
        )
        SELECT 
            id,
            COALESCE(nama_pohon, '') as nama_pohon,
            COALESCE(parent, 0) as parent,
            COALESCE(jenis_pohon, '') as jenis_pohon,
            level_pohon,
            COALESCE(kode_opd, '') as kode_opd,
            COALESCE(keterangan, '') as keterangan,
            COALESCE(keterangan_crosscutting, '') as keterangan_crosscutting,
            COALESCE(tahun, '') as tahun,
            COALESCE(status, '') as status,
            COALESCE(is_active) as is_active
        FROM parent_tree
        WHERE level_pohon = 4
        LIMIT 1
    `

	var pokin domain.PohonKinerja
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(
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
	)
	if err != nil {
		return domain.PohonKinerja{}, err
	}

	return pokin, nil
}

func (repository *CascadingOpdRepositoryImpl) CalculateTotalAnggaranByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	script := `
        WITH RECURSIVE pohon_tree AS (
            -- Base case: pohon yang dicari
            SELECT id
            FROM tb_pohon_kinerja
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: semua child dari pohon
            SELECT pk.id
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_tree pt ON pk.parent = pt.id
        )
        SELECT COALESCE(SUM(rb.anggaran), 0) as total_anggaran
        FROM pohon_tree
        INNER JOIN tb_rencana_kinerja rk ON rk.id_pohon = pohon_tree.id
        INNER JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
        INNER JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
    `

	var totalAnggaran int64
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(&totalAnggaran)
	if err != nil {
		return 0, err
	}

	return totalAnggaran, nil
}

// Ambil anggaran dari satu rencana kinerja
func (repository *CascadingOpdRepositoryImpl) GetAnggaranByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int64, error) {
	script := `
        SELECT COALESCE(SUM(rb.anggaran), 0) as total_anggaran
        FROM tb_rencana_aksi ra
        INNER JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        WHERE ra.rencana_kinerja_id = ?
    `

	var totalAnggaran int64
	err := tx.QueryRowContext(ctx, script, rekinId).Scan(&totalAnggaran)
	if err != nil {
		return 0, err
	}

	return totalAnggaran, nil
}

// Ambil semua operational children dari tactical
func (repository *CascadingOpdRepositoryImpl) FindOperationalChildrenByTacticalId(ctx context.Context, tx *sql.Tx, tacticalId int) ([]int, error) {
	script := `
		SELECT id 
		FROM tb_pohon_kinerja 
		WHERE parent = ? AND level_pohon = 6
	`

	rows, err := tx.QueryContext(ctx, script, tacticalId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operationalIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			operationalIds = append(operationalIds, id)
		}
	}

	return operationalIds, nil
}

// Ambil semua tactical children dari strategic
func (repository *CascadingOpdRepositoryImpl) FindTacticalChildrenByStrategicId(ctx context.Context, tx *sql.Tx, strategicId int) ([]int, error) {
	script := `
		SELECT id 
		FROM tb_pohon_kinerja 
		WHERE parent = ? AND level_pohon = 5
	`

	rows, err := tx.QueryContext(ctx, script, strategicId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tacticalIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			tacticalIds = append(tacticalIds, id)
		}
	}

	return tacticalIds, nil
}

// Ambil total anggaran dari semua rencana kinerja di pohon ini
func (repository *CascadingOpdRepositoryImpl) GetTotalAnggaranByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	script := `
		SELECT COALESCE(SUM(rb.anggaran), 0) as total_anggaran
		FROM tb_rencana_kinerja rk
		INNER JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
		INNER JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
		WHERE rk.id_pohon = ?
	`

	var totalAnggaran int64
	err := tx.QueryRowContext(ctx, script, pokinId).Scan(&totalAnggaran)
	if err != nil {
		return 0, err
	}

	return totalAnggaran, nil
}

// Ambil kode subkegiatan untuk collect program (level 6)
func (repository *CascadingOpdRepositoryImpl) FindKodeSubkegiatanByPokinId(ctx context.Context, tx *sql.Tx, pokinId int) ([]string, error) {
	script := `
		SELECT DISTINCT kode_subkegiatan
		FROM tb_rencana_kinerja
		WHERE id_pohon = ?
		AND kode_subkegiatan IS NOT NULL AND kode_subkegiatan != ''
	`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kodeList []string
	for rows.Next() {
		var kode string
		if err := rows.Scan(&kode); err == nil {
			kodeList = append(kodeList, kode)
		}
	}

	return kodeList, nil
}

// Ambil kode subkegiatan dari child nodes recursive (level 4 & 5)
func (repository *CascadingOpdRepositoryImpl) FindKodeSubkegiatanFromChildren(ctx context.Context, tx *sql.Tx, pokinId int) ([]string, error) {
	script := `
		WITH RECURSIVE pohon_tree AS (
			SELECT id, level_pohon
			FROM tb_pohon_kinerja
			WHERE parent = ?
			
			UNION ALL
			
			SELECT pk.id, pk.level_pohon
			FROM tb_pohon_kinerja pk
			INNER JOIN pohon_tree pt ON pk.parent = pt.id
		)
		SELECT DISTINCT rk.kode_subkegiatan
		FROM pohon_tree
		INNER JOIN tb_rencana_kinerja rk ON rk.id_pohon = pohon_tree.id
		WHERE rk.kode_subkegiatan IS NOT NULL AND rk.kode_subkegiatan != ''
	`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kodeList []string
	for rows.Next() {
		var kode string
		if err := rows.Scan(&kode); err == nil {
			kodeList = append(kodeList, kode)
		}
	}

	return kodeList, nil
}

func (repository *CascadingOpdRepositoryImpl) FindPokinByNipAndTahun(ctx context.Context, tx *sql.Tx, nip string, tahun string) ([]domain.PohonKinerja, error) {
	script := `
		SELECT DISTINCT
			pk.id,
			COALESCE(pk.nama_pohon, '') as nama_pohon,
			COALESCE(pk.parent, 0) as parent,
			COALESCE(pk.jenis_pohon, '') as jenis_pohon,
			COALESCE(pk.level_pohon, 0) as level_pohon,
			COALESCE(pk.kode_opd, '') as kode_opd,
			COALESCE(pk.keterangan, '') as keterangan,
			COALESCE(pk.keterangan_crosscutting, '') as keterangan_crosscutting,
			COALESCE(pk.tahun, '') as tahun,
			COALESCE(pk.status, '') as status,
			COALESCE(pk.is_active) as is_active
		FROM tb_pohon_kinerja pk
		INNER JOIN tb_rencana_kinerja rk ON rk.id_pohon = pk.id
		WHERE rk.pegawai_id = ?
		AND pk.tahun = ?
		ORDER BY COALESCE(pk.level_pohon, 0), pk.id ASC
	`

	rows, err := tx.QueryContext(ctx, script, nip, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
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
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}

	if pokins == nil {
		pokins = make([]domain.PohonKinerja, 0)
	}

	return pokins, nil
}

func (repository *CascadingOpdRepositoryImpl) GetTotalAnggaranByPokinIdWithPelaksana(ctx context.Context, tx *sql.Tx, pokinId int) (int64, error) {
	script := `
		SELECT COALESCE(SUM(rb.anggaran), 0) as total_anggaran
		FROM tb_rencana_kinerja rk
		INNER JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rk.id
		INNER JOIN tb_rincian_belanja rb ON rb.rencana_aksi_id = ra.id
		INNER JOIN tb_pegawai peg ON peg.nip = rk.pegawai_id
		INNER JOIN tb_pelaksana_pokin pp ON pp.pegawai_id = peg.id AND pp.pokin_id = ?
		WHERE rk.id_pohon = ?
	`

	var totalAnggaran int64
	err := tx.QueryRowContext(ctx, script, pokinId, pokinId).Scan(&totalAnggaran)
	if err != nil {
		return 0, err
	}

	return totalAnggaran, nil
}

func (repository *CascadingOpdRepositoryImpl) FindTargetByIndikatorIdsBatch(ctx context.Context, tx *sql.Tx, indikatorIds []string) (map[string][]domain.Target, error) {
	if len(indikatorIds) == 0 {
		return make(map[string][]domain.Target), nil
	}

	// Build query dengan IN clause
	placeholders := make([]string, len(indikatorIds))
	args := make([]interface{}, len(indikatorIds))
	for i, id := range indikatorIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(
		"SELECT id, indikator_id, target, satuan FROM tb_target WHERE indikator_id IN (%s)",
		strings.Join(placeholders, ","),
	)

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]domain.Target)
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan)
		if err != nil {
			return nil, err
		}
		result[target.IndikatorId] = append(result[target.IndikatorId], target)
	}

	return result, nil
}

func (repository *CascadingOpdRepositoryImpl) FindIndikatorTargetByPokinIds(
	ctx context.Context,
	tx *sql.Tx,
	pokinIds []int,
) (map[int][]domain.Indikator, error) {

	const op = "cascadingopd_repository.FindIndikatorTargetByPokinIds"

	if len(pokinIds) == 0 {
		return map[int][]domain.Indikator{}, nil
	}

	baseQuery := `
		SELECT
			i.id,
			i.pokin_id,
			i.indikator,
			t.id,
			t.target,
			t.satuan
		FROM tb_indikator i
		LEFT JOIN tb_target t ON i.id = t.indikator_id
		WHERE i.pokin_id IN (?)
		ORDER BY i.pokin_id, i.id
	`

	query, args := helper.BuildInQuery(baseQuery, pokinIds)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	// pohonId -> indikatorId -> indikator
	pohonMap := make(map[int]map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indikatorId string
			pohonId     int
			indikator   string
			targetId    sql.NullString
			target      sql.NullString
			satuan      sql.NullString
		)

		if err := rows.Scan(
			&indikatorId,
			&pohonId,
			&indikator,
			&targetId,
			&target,
			&satuan,
		); err != nil {
			return nil, fmt.Errorf("%s: scan failed: %w", op, err)
		}

		if pohonMap[pohonId] == nil {
			pohonMap[pohonId] = make(map[string]*domain.Indikator)
		}

		if pohonMap[pohonId][indikatorId] == nil {
			pohonMap[pohonId][indikatorId] = &domain.Indikator{
				Id:        indikatorId,
				Indikator: indikator,
				Target:    make([]domain.Target, 0),
			}
		}

		// append target jika ada
		if targetId.Valid {
			pohonMap[pohonId][indikatorId].Target = append(
				pohonMap[pohonId][indikatorId].Target,
				domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId,
					Target:      target.String,
					Satuan:      satuan.String,
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	// flatten → map[int][]domain.Indikator
	result := make(map[int][]domain.Indikator)

	for pohonId, indikatorMap := range pohonMap {
		for _, ind := range indikatorMap {
			result[pohonId] = append(result[pohonId], *ind)
		}
	}

	return result, nil
}
