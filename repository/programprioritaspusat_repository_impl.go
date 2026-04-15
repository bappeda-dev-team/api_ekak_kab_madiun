package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"log"
	"strings"
)

type ProgramPrioritasPusatRepositoryImpl struct {
}

func NewProgramPrioritasPusatRepositoryImpl() *ProgramPrioritasPusatRepositoryImpl {
	return &ProgramPrioritasPusatRepositoryImpl{}
}

func (repository *ProgramPrioritasPusatRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, programPrioritasPusat domain.ProgramPrioritasPusat) (domain.ProgramPrioritasPusat, error) {
	script := "INSERT INTO tb_program_prioritas_pusat (nama_tagging, kode_program_prioritas_pusat, keterangan_program_prioritas_pusat, keterangan, tahun_awal, tahun_akhir) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		programPrioritasPusat.NamaTagging,
		programPrioritasPusat.KodeProgramPrioritasPusat,
		programPrioritasPusat.KeteranganProgramPrioritasPusat,
		programPrioritasPusat.Keterangan,
		programPrioritasPusat.TahunAwal,
		programPrioritasPusat.TahunAkhir)
	if err != nil {
		return domain.ProgramPrioritasPusat{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.ProgramPrioritasPusat{}, err
	}
	programPrioritasPusat.Id = int(id)

	return programPrioritasPusat, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, programPrioritasPusat domain.ProgramPrioritasPusat) (domain.ProgramPrioritasPusat, error) {
	script := "UPDATE tb_program_prioritas_pusat SET nama_tagging = ?, keterangan_program_prioritas_pusat = ?, keterangan = ?, tahun_awal = ?, tahun_akhir = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, programPrioritasPusat.NamaTagging, programPrioritasPusat.KeteranganProgramPrioritasPusat, programPrioritasPusat.Keterangan, programPrioritasPusat.TahunAwal, programPrioritasPusat.TahunAkhir, programPrioritasPusat.Id)
	if err != nil {
		return domain.ProgramPrioritasPusat{}, err
	}
	return programPrioritasPusat, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_program_prioritas_pusat WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.ProgramPrioritasPusat, error) {
	script := "SELECT id, nama_tagging, kode_program_prioritas_pusat, keterangan_program_prioritas_pusat, keterangan, tahun_awal, tahun_akhir FROM tb_program_prioritas_pusat WHERE id = ?"
	var programPrioritasPusat domain.ProgramPrioritasPusat
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&programPrioritasPusat.Id,
		&programPrioritasPusat.NamaTagging,
		&programPrioritasPusat.KodeProgramPrioritasPusat,
		&programPrioritasPusat.KeteranganProgramPrioritasPusat,
		&programPrioritasPusat.Keterangan,
		&programPrioritasPusat.TahunAwal,
		&programPrioritasPusat.TahunAkhir,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ProgramPrioritasPusat{}, errors.New("program prioritas pusat tidak ditemukan")
		}
		return domain.ProgramPrioritasPusat{}, err
	}
	return programPrioritasPusat, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramPrioritasPusat, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            pu.id, 
            pu.nama_tagging, 
            pu.kode_program_prioritas_pusat, 
            pu.keterangan_program_prioritas_pusat, 
            pu.keterangan, 
            pu.tahun_awal, 
            pu.tahun_akhir,
            CASE 
                WHEN EXISTS (
                    SELECT 1 
                    FROM tb_keterangan_tagging_program_prioritas_pusat ktpu
                    JOIN tb_tagging_pokin tp ON ktpu.id_tagging = tp.id
                    WHERE ktpu.kode_program_prioritas_pusat = pu.kode_program_prioritas_pusat
                ) THEN TRUE 
                ELSE FALSE 
            END as is_active
        FROM tb_program_prioritas_pusat pu
        WHERE pu.tahun_awal >= ? AND pu.tahun_akhir <= ?`

	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir)
	if err != nil {
		return []domain.ProgramPrioritasPusat{}, err
	}
	defer rows.Close()
	log.Printf("TahunAWal: %s TAhun Akhir: %s", tahunAwal, tahunAkhir)

	var programPrioritasPusatList []domain.ProgramPrioritasPusat
	for rows.Next() {
		var programPrioritasPusat domain.ProgramPrioritasPusat
		err = rows.Scan(
			&programPrioritasPusat.Id,
			&programPrioritasPusat.NamaTagging,
			&programPrioritasPusat.KodeProgramPrioritasPusat,
			&programPrioritasPusat.KeteranganProgramPrioritasPusat,
			&programPrioritasPusat.Keterangan,
			&programPrioritasPusat.TahunAwal,
			&programPrioritasPusat.TahunAkhir,
			&programPrioritasPusat.IsActive,
		)
		if err != nil {
			return []domain.ProgramPrioritasPusat{}, err
		}
		programPrioritasPusatList = append(programPrioritasPusatList, programPrioritasPusat)
	}
	return programPrioritasPusatList, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindByKodeProgramPrioritasPusat(ctx context.Context, tx *sql.Tx, kodeProgramPrioritasPusat string) (domain.ProgramPrioritasPusat, error) {
	script := "SELECT id, nama_tagging, kode_program_prioritas_pusat, keterangan_program_prioritas_pusat, keterangan, tahun_awal, tahun_akhir FROM tb_program_prioritas_pusat WHERE kode_program_prioritas_pusat = ?"
	var programPrioritasPusat domain.ProgramPrioritasPusat
	err := tx.QueryRowContext(ctx, script, kodeProgramPrioritasPusat).Scan(
		&programPrioritasPusat.Id,
		&programPrioritasPusat.NamaTagging,
		&programPrioritasPusat.KodeProgramPrioritasPusat,
		&programPrioritasPusat.KeteranganProgramPrioritasPusat,
		&programPrioritasPusat.Keterangan,
		&programPrioritasPusat.TahunAwal,
		&programPrioritasPusat.TahunAkhir,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ProgramPrioritasPusat{}, errors.New("program prioritas pusat tidak ditemukan")
		}
		return domain.ProgramPrioritasPusat{}, err
	}
	return programPrioritasPusat, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramPrioritasPusat, error) {
	script := `
        SELECT id, nama_tagging, kode_program_prioritas_pusat, keterangan_program_prioritas_pusat, keterangan, tahun_awal, tahun_akhir 
        FROM tb_program_prioritas_pusat 
        WHERE ? BETWEEN tahun_awal AND tahun_akhir
        ORDER BY tahun_awal ASC`

	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programPrioritasPusatList []domain.ProgramPrioritasPusat
	for rows.Next() {
		var programPrioritasPusat domain.ProgramPrioritasPusat
		err = rows.Scan(
			&programPrioritasPusat.Id,
			&programPrioritasPusat.NamaTagging,
			&programPrioritasPusat.KodeProgramPrioritasPusat,
			&programPrioritasPusat.KeteranganProgramPrioritasPusat,
			&programPrioritasPusat.Keterangan,
			&programPrioritasPusat.TahunAwal,
			&programPrioritasPusat.TahunAkhir,
		)
		if err != nil {
			return nil, err
		}
		programPrioritasPusatList = append(programPrioritasPusatList, programPrioritasPusat)
	}
	return programPrioritasPusatList, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindUnusedByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramPrioritasPusat, error) {
	script := `
        SELECT DISTINCT 
            pu.id, 
            pu.nama_tagging, 
            pu.kode_program_prioritas_pusat, 
            pu.keterangan_program_prioritas_pusat, 
            pu.keterangan, 
            pu.tahun_awal, 
            pu.tahun_akhir
        FROM 
            tb_program_prioritas_pusat pu
        WHERE 
            ? BETWEEN pu.tahun_awal AND pu.tahun_akhir
            AND NOT EXISTS (
                SELECT 1 
                FROM tb_keterangan_tagging_program_prioritas_pusat ktpu 
                WHERE ktpu.kode_program_prioritas_pusat = pu.kode_program_prioritas_pusat 
                AND ktpu.tahun = ?
            )
        ORDER BY pu.tahun_awal ASC`

	rows, err := tx.QueryContext(ctx, script, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programPrioritasPusatList []domain.ProgramPrioritasPusat
	for rows.Next() {
		var programPrioritasPusat domain.ProgramPrioritasPusat
		err = rows.Scan(
			&programPrioritasPusat.Id,
			&programPrioritasPusat.NamaTagging,
			&programPrioritasPusat.KodeProgramPrioritasPusat,
			&programPrioritasPusat.KeteranganProgramPrioritasPusat,
			&programPrioritasPusat.Keterangan,
			&programPrioritasPusat.TahunAwal,
			&programPrioritasPusat.TahunAkhir,
		)
		if err != nil {
			return nil, err
		}
		programPrioritasPusatList = append(programPrioritasPusatList, programPrioritasPusat)
	}
	return programPrioritasPusatList, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindByIdTerkait(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.ProgramPrioritasPusat, error) {
	if len(ids) == 0 {
		return []domain.ProgramPrioritasPusat{}, errors.New("ids tidak boleh kosong")
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT id, nama_tagging, kode_program_prioritas_pusat, 
		       keterangan_program_prioritas_pusat, keterangan, tahun_awal, tahun_akhir 
		FROM tb_program_prioritas_pusat 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.ProgramPrioritasPusat{}, err
	}
	defer rows.Close()

	var result []domain.ProgramPrioritasPusat
	for rows.Next() {
		var programPrioritasPusat domain.ProgramPrioritasPusat
		err := rows.Scan(
			&programPrioritasPusat.Id,
			&programPrioritasPusat.NamaTagging,
			&programPrioritasPusat.KodeProgramPrioritasPusat,
			&programPrioritasPusat.KeteranganProgramPrioritasPusat,
			&programPrioritasPusat.Keterangan,
			&programPrioritasPusat.TahunAwal,
			&programPrioritasPusat.TahunAkhir,
		)
		if err != nil {
			return []domain.ProgramPrioritasPusat{}, err
		}
		result = append(result, programPrioritasPusat)
	}

	if err = rows.Err(); err != nil {
		return []domain.ProgramPrioritasPusat{}, err
	}

	return result, nil
}

func (repository *ProgramPrioritasPusatRepositoryImpl) FindProgramPrioritasPusatByKodesBatch(ctx context.Context, tx *sql.Tx, kodes []string) (map[string]*domain.ProgramPrioritasPusat, error) {
	if len(kodes) == 0 {
		return make(map[string]*domain.ProgramPrioritasPusat), nil
	}

	placeholders := make([]string, len(kodes))
	args := make([]interface{}, len(kodes))
	for i, kode := range kodes {
		placeholders[i] = "?"
		args[i] = kode
	}

	script := fmt.Sprintf(`
		SELECT id, kode_program_prioritas_pusat, keterangan_program_prioritas_pusat, tahun_awal, tahun_akhir
		FROM tb_program_prioritas_pusat
		WHERE kode_program_prioritas_pusat IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domain.ProgramPrioritasPusat)
	for rows.Next() {
		var program domain.ProgramPrioritasPusat
		err := rows.Scan(&program.Id, &program.KodeProgramPrioritasPusat, &program.KeteranganProgramPrioritasPusat, &program.TahunAwal, &program.TahunAkhir)
		if err != nil {
			return nil, err
		}
		result[program.KodeProgramPrioritasPusat] = &program
	}

	return result, nil
}
