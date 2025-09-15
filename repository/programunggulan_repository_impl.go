package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type ProgramUnggulanRepositoryImpl struct {
}

func NewProgramUnggulanRepositoryImpl() *ProgramUnggulanRepositoryImpl {
	return &ProgramUnggulanRepositoryImpl{}
}

func (repository *ProgramUnggulanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, programUnggulan domain.ProgramUnggulan) (domain.ProgramUnggulan, error) {
	script := "INSERT INTO tb_program_unggulan (nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		programUnggulan.NamaTagging,
		programUnggulan.KodeProgramUnggulan,
		programUnggulan.KeteranganProgramUnggulan,
		programUnggulan.Keterangan,
		programUnggulan.TahunAwal,
		programUnggulan.TahunAkhir)
	if err != nil {
		return domain.ProgramUnggulan{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.ProgramUnggulan{}, err
	}
	programUnggulan.Id = int(id)

	return programUnggulan, nil
}

func (repository *ProgramUnggulanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, programUnggulan domain.ProgramUnggulan) (domain.ProgramUnggulan, error) {
	script := "UPDATE tb_program_unggulan SET nama_tagging = ?, keterangan_program_unggulan = ?, keterangan = ?, tahun_awal = ?, tahun_akhir = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, programUnggulan.NamaTagging, programUnggulan.KeteranganProgramUnggulan, programUnggulan.Keterangan, programUnggulan.TahunAwal, programUnggulan.TahunAkhir, programUnggulan.Id)
	if err != nil {
		return domain.ProgramUnggulan{}, err
	}
	return programUnggulan, nil
}

func (repository *ProgramUnggulanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_program_unggulan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}
func (repository *ProgramUnggulanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ProgramUnggulan, error) {
	script := "SELECT id, nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir FROM tb_program_unggulan WHERE id = ?"
	var programUnggulan domain.ProgramUnggulan
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&programUnggulan.Id,
		&programUnggulan.NamaTagging,
		&programUnggulan.KodeProgramUnggulan,
		&programUnggulan.KeteranganProgramUnggulan,
		&programUnggulan.Keterangan,
		&programUnggulan.TahunAwal,
		&programUnggulan.TahunAkhir,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ProgramUnggulan{}, errors.New("program unggulan tidak ditemukan")
		}
		return domain.ProgramUnggulan{}, err
	}
	return programUnggulan, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramUnggulan, error) {
	script := "SELECT id, nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir FROM tb_program_unggulan WHERE tahun_awal >= ? AND tahun_akhir <= ?"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domain.ProgramUnggulan{}, err
	}
	defer rows.Close()
	var programUnggulanList []domain.ProgramUnggulan
	for rows.Next() {
		var programUnggulan domain.ProgramUnggulan
		err = rows.Scan(&programUnggulan.Id, &programUnggulan.NamaTagging, &programUnggulan.KodeProgramUnggulan, &programUnggulan.KeteranganProgramUnggulan, &programUnggulan.Keterangan, &programUnggulan.TahunAwal, &programUnggulan.TahunAkhir)
		if err != nil {
			return []domain.ProgramUnggulan{}, err
		}
		programUnggulanList = append(programUnggulanList, programUnggulan)
	}
	return programUnggulanList, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindByKodeProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) (domain.ProgramUnggulan, error) {
	script := "SELECT id, nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir FROM tb_program_unggulan WHERE kode_program_unggulan = ?"
	var programUnggulan domain.ProgramUnggulan
	err := tx.QueryRowContext(ctx, script, kodeProgramUnggulan).Scan(
		&programUnggulan.Id,
		&programUnggulan.NamaTagging,
		&programUnggulan.KodeProgramUnggulan,
		&programUnggulan.KeteranganProgramUnggulan,
		&programUnggulan.Keterangan,
		&programUnggulan.TahunAwal,
		&programUnggulan.TahunAkhir,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ProgramUnggulan{}, errors.New("program unggulan tidak ditemukan")
		}
		return domain.ProgramUnggulan{}, err
	}
	return programUnggulan, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramUnggulan, error) {
	script := `
        SELECT id, nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir 
        FROM tb_program_unggulan 
        WHERE ? BETWEEN tahun_awal AND tahun_akhir
        ORDER BY tahun_awal ASC`

	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programUnggulanList []domain.ProgramUnggulan
	for rows.Next() {
		var programUnggulan domain.ProgramUnggulan
		err = rows.Scan(
			&programUnggulan.Id,
			&programUnggulan.NamaTagging,
			&programUnggulan.KodeProgramUnggulan,
			&programUnggulan.KeteranganProgramUnggulan,
			&programUnggulan.Keterangan,
			&programUnggulan.TahunAwal,
			&programUnggulan.TahunAkhir,
		)
		if err != nil {
			return nil, err
		}
		programUnggulanList = append(programUnggulanList, programUnggulan)
	}
	return programUnggulanList, nil
}
