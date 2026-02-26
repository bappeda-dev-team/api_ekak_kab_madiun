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
