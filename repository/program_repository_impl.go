package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
)

type ProgramRepositoryImpl struct {
}

func NewProgramRepositoryImpl() *ProgramRepositoryImpl {
	return &ProgramRepositoryImpl{}
}

func (repository *ProgramRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error) {
	scriptProgram := "INSERT INTO tb_master_program (id, kode_program, nama_program, tahun, is_active) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, scriptProgram, program.Id, program.KodeProgram, program.NamaProgram, program.Tahun, program.IsActive)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}
	scriptIndikator := "INSERT INTO tb_indikator (id, program_id, indikator, tahun) VALUES (?, ?, ?, ?)"
	for _, indikator := range program.Indikator {
		_, err = tx.ExecContext(ctx, scriptIndikator, indikator.Id, indikator.ProgramId, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
		for _, target := range indikator.Target {
			_, err = tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domainmaster.ProgramKegiatan{}, err
			}
		}
	}
	return program, nil
}

func (repository *ProgramRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, program domainmaster.ProgramKegiatan) (domainmaster.ProgramKegiatan, error) {
	scriptProgram := "UPDATE tb_master_program SET kode_program = ?, nama_program = ?, tahun = ?, is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptProgram, program.KodeProgram, program.NamaProgram, program.Tahun, program.IsActive, program.Id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE program_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, program.Id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}

	scriptIndikator := "INSERT INTO tb_indikator (id, program_id, indikator, tahun) VALUES (?, ?, ?, ?)"
	for _, indikator := range program.Indikator {
		_, err = tx.ExecContext(ctx, scriptIndikator, indikator.Id, indikator.ProgramId, indikator.Indikator, indikator.Tahun)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTarget, indikator.Id)
		if err != nil {
			return domainmaster.ProgramKegiatan{}, err
		}

		scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
		for _, target := range indikator.Target {
			_, err = tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return domainmaster.ProgramKegiatan{}, err
			}
		}
	}

	return program, nil
}

func (repository *ProgramRepositoryImpl) FindIndikatorByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]domain.Indikator, error) {
	script := "SELECT id, program_id, indikator, tahun FROM tb_indikator WHERE program_id = ?"
	params := []interface{}{programId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.ProgramId, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *ProgramRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
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

func (repository *ProgramRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.ProgramKegiatan, error) {
	// Query untuk mendapatkan data program
	scriptProgram := "SELECT id, kode_program, nama_program, tahun, is_active FROM tb_master_program WHERE id = ?"
	rows, err := tx.QueryContext(ctx, scriptProgram, id)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("gagal mengambil data program: %v", err)
	}
	defer rows.Close()

	var program domainmaster.ProgramKegiatan
	if !rows.Next() {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("program dengan id %s tidak ditemukan", id)
	}

	// Scan data program
	err = rows.Scan(
		&program.Id,
		&program.KodeProgram,
		&program.NamaProgram,
		&program.Tahun,
		&program.IsActive,
	)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, fmt.Errorf("gagal scan data program: %v", err)
	}

	return program, nil
}

func (repository *ProgramRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	scriptIndikator := "DELETE FROM tb_indikator WHERE program_id = ?"
	_, err := tx.ExecContext(ctx, scriptIndikator, id)
	if err != nil {
		return err
	}

	scriptProgram := "DELETE FROM tb_master_program WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptProgram, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *ProgramRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.ProgramKegiatan, error) {
	script := "SELECT id, kode_program, nama_program, tahun, is_active FROM tb_master_program ORDER BY id"

	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []domainmaster.ProgramKegiatan
	for rows.Next() {
		var program domainmaster.ProgramKegiatan
		err := rows.Scan(
			&program.Id,
			&program.KodeProgram,
			&program.NamaProgram,
			&program.Tahun,
			&program.IsActive,
		)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	return programs, nil
}

func (repository *ProgramRepositoryImpl) FindByKodeProgram(ctx context.Context, tx *sql.Tx, kodeProgram string) (domainmaster.ProgramKegiatan, error) {
	SQL := `
    SELECT 
        kode_program,
        nama_program
    FROM tb_master_program
    WHERE kode_program = ?
    `

	var program domainmaster.ProgramKegiatan
	err := tx.QueryRowContext(ctx, SQL, kodeProgram).Scan(
		&program.KodeProgram,
		&program.NamaProgram,
	)
	if err != nil {
		return domainmaster.ProgramKegiatan{}, err
	}

	return program, nil
}

// cukup ini saja untuk tarik indikator
// program, kegiatan, subkegiatan
func (repository *ProgramRepositoryImpl) FindIndikatorTargetByKodeRelasiBatch(
	ctx context.Context,
	tx *sql.Tx,
	kodeRelasis []string,
) (map[string][]domain.Indikator, error) {

	const op = "program_repository.FindIndikatorTargetByKodeRelasiBatch"

	if len(kodeRelasis) == 0 {
		return map[string][]domain.Indikator{}, nil
	}

	baseQuery := `
		SELECT
			i.id,
			i.kode,
			i.indikator,
			t.id,
			t.target,
			t.satuan
		FROM tb_indikator i
		LEFT JOIN tb_target t ON i.id = t.indikator_id
		WHERE i.kode IN (?)
		ORDER BY i.kode, i.id
	`

	query, args := helper.BuildInQueryString(baseQuery, kodeRelasis)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	// kodeRelasi -> indikatorId -> indikator
	indikatorMap := make(map[string]map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indikatorId string
			kodeRelasi  string
			indikator   string
			targetId    sql.NullString
			target      sql.NullString
			satuan      sql.NullString
		)

		if err := rows.Scan(
			&indikatorId,
			&kodeRelasi,
			&indikator,
			&targetId,
			&target,
			&satuan,
		); err != nil {
			return nil, fmt.Errorf("%s: scan failed: %w", op, err)
		}

		if indikatorMap[kodeRelasi] == nil {
			indikatorMap[kodeRelasi] = make(map[string]*domain.Indikator)
		}

		if indikatorMap[kodeRelasi][indikatorId] == nil {
			indikatorMap[kodeRelasi][indikatorId] = &domain.Indikator{
				Id:        indikatorId,
				Indikator: indikator,
				Target:    make([]domain.Target, 0),
			}
		}

		if targetId.Valid {
			indikatorMap[kodeRelasi][indikatorId].Target = append(
				indikatorMap[kodeRelasi][indikatorId].Target,
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

	// flatten → map[string][]domain.Indikator
	result := make(map[string][]domain.Indikator)
	for kode, indMap := range indikatorMap {
		for _, ind := range indMap {
			result[kode] = append(result[kode], *ind)
		}
	}

	return result, nil
}

// queryMasterIndikatorProgramKegiatanByKodes — dipakai program/kegiatan/subkegiatan (match indikator.kode = kode master).
func queryMasterIndikatorProgramKegiatanByKodes(
	ctx context.Context,
	tx *sql.Tx,
	kodes []string,
) (map[string][]domainmaster.MasterIndikatorProgramKegiatan, error) {
	if len(kodes) == 0 {
		return map[string][]domainmaster.MasterIndikatorProgramKegiatan{}, nil
	}
	baseQuery := `
		SELECT id, kode_indikator, kode, indikator
		FROM tb_master_indikator_program_kegiatan
		WHERE kode IN (?)
		ORDER BY kode ASC, id ASC
	`
	query, args := helper.BuildInQueryString(baseQuery, kodes)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]domainmaster.MasterIndikatorProgramKegiatan)
	for rows.Next() {
		var row domainmaster.MasterIndikatorProgramKegiatan
		if err := rows.Scan(&row.Id, &row.KodeIndikator, &row.Kode, &row.Indikator); err != nil {
			return nil, err
		}
		result[row.Kode] = append(result[row.Kode], row)
	}
	return result, rows.Err()
}

func (repository *ProgramRepositoryImpl) FindMasterIndikatorProgramKegiatanByKodes(
	ctx context.Context, tx *sql.Tx, kodes []string,
) (map[string][]domainmaster.MasterIndikatorProgramKegiatan, error) {
	return queryMasterIndikatorProgramKegiatanByKodes(ctx, tx, kodes)
}

// ─────────────────────────────────────────────────────────────────────────────
// Upsert indikator program
// ─────────────────────────────────────────────────────────────────────────────

// FindIndikatorIdsByProgramId — mengambil semua id indikator untuk 1 program.
// Dipakai untuk diff: indikator existing yang tidak ada di request → dihapus.
func (r *ProgramRepositoryImpl) FindIndikatorIdsByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM tb_indikator WHERE program_id = ?`, programId)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil id indikator program %s: %w", programId, err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// UpsertIndikator — INSERT baru atau UPDATE jika id sudah ada (ON DUPLICATE KEY).
func (r *ProgramRepositoryImpl) UpsertIndikator(ctx context.Context, tx *sql.Tx, ind domain.Indikator) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO tb_indikator (id, program_id, indikator, tahun)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			indikator = VALUES(indikator),
			tahun     = VALUES(tahun)`,
		ind.Id, ind.ProgramId, ind.Indikator, ind.Tahun,
	)
	if err != nil {
		return fmt.Errorf("gagal upsert indikator %s: %w", ind.Id, err)
	}
	return nil
}

// DeleteIndikatorWithTargets — hapus target lalu indikator (cascade manual).
func (r *ProgramRepositoryImpl) DeleteIndikatorWithTargets(ctx context.Context, tx *sql.Tx, indikatorId string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_target WHERE indikator_id = ?`, indikatorId); err != nil {
		return fmt.Errorf("gagal hapus target indikator %s: %w", indikatorId, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_indikator WHERE id = ?`, indikatorId); err != nil {
		return fmt.Errorf("gagal hapus indikator %s: %w", indikatorId, err)
	}
	return nil
}

// ReplaceTargets — hapus semua target lama lalu insert target baru.
func (r *ProgramRepositoryImpl) ReplaceTargets(ctx context.Context, tx *sql.Tx, indikatorId string, targets []domain.Target) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM tb_target WHERE indikator_id = ?`, indikatorId); err != nil {
		return fmt.Errorf("gagal hapus target lama indikator %s: %w", indikatorId, err)
	}
	for _, t := range targets {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO tb_target (id, indikator_id, target, satuan, tahun)
			VALUES (?, ?, ?, ?, ?)`,
			t.Id, indikatorId, t.Target, t.Satuan, t.Tahun,
		)
		if err != nil {
			return fmt.Errorf("gagal insert target baru %s: %w", t.Id, err)
		}
	}
	return nil
}

// FindIndikatorWithTargetsByProgramId — mengambil semua indikator + target untuk 1 program.
// Dipakai untuk membangun response setelah upsert.
func (r *ProgramRepositoryImpl) FindIndikatorWithTargetsByProgramId(ctx context.Context, tx *sql.Tx, programId string) ([]domain.Indikator, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT
			i.id, i.program_id, i.indikator, i.tahun,
			COALESCE(t.id, '')     AS target_id,
			COALESCE(t.target, '') AS target,
			COALESCE(t.satuan, '') AS satuan,
			COALESCE(t.tahun, '')  AS target_tahun
		FROM tb_indikator i
		LEFT JOIN tb_target t ON t.indikator_id = i.id
		WHERE i.program_id = ?
		ORDER BY i.id, t.id`,
		programId,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal query indikator program %s: %w", programId, err)
	}
	defer rows.Close()

	// Map indikatorId → *domain.Indikator untuk dedup
	indMap := make(map[string]*domain.Indikator)
	var order []string // pertahankan urutan insert

	for rows.Next() {
		var (
			indId, progId, indikator, tahunInd string
			targetId, target, satuan, tahunTgt string
		)
		if err := rows.Scan(&indId, &progId, &indikator, &tahunInd,
			&targetId, &target, &satuan, &tahunTgt); err != nil {
			return nil, err
		}

		if indMap[indId] == nil {
			indMap[indId] = &domain.Indikator{
				Id:        indId,
				ProgramId: progId,
				Indikator: indikator,
				Tahun:     tahunInd,
				Target:    []domain.Target{},
			}
			order = append(order, indId)
		}

		if targetId != "" {
			indMap[indId].Target = append(indMap[indId].Target, domain.Target{
				Id:          targetId,
				IndikatorId: indId,
				Target:      target,
				Satuan:      satuan,
				Tahun:       tahunTgt,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]domain.Indikator, 0, len(order))
	for _, id := range order {
		result = append(result, *indMap[id])
	}
	return result, nil
}
