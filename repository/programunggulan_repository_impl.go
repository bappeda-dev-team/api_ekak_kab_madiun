package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"strings"
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

// func (repository *ProgramUnggulanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramUnggulan, error) {
// 	script := "SELECT id, nama_tagging, kode_program_unggulan, keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir FROM tb_program_unggulan WHERE tahun_awal >= ? AND tahun_akhir <= ?"
// 	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir)
// 	if err != nil {
// 		return []domain.ProgramUnggulan{}, err
// 	}
// 	defer rows.Close()
// 	var programUnggulanList []domain.ProgramUnggulan
// 	for rows.Next() {
// 		var programUnggulan domain.ProgramUnggulan
// 		err = rows.Scan(&programUnggulan.Id, &programUnggulan.NamaTagging, &programUnggulan.KodeProgramUnggulan, &programUnggulan.KeteranganProgramUnggulan, &programUnggulan.Keterangan, &programUnggulan.TahunAwal, &programUnggulan.TahunAkhir)
// 		if err != nil {
// 			return []domain.ProgramUnggulan{}, err
// 		}
// 		programUnggulanList = append(programUnggulanList, programUnggulan)
// 	}
// 	return programUnggulanList, nil
// }

func (repository *ProgramUnggulanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string) ([]domain.ProgramUnggulan, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
	SELECT 
		pu.id, pu.nama_tagging, pu.kode_program_unggulan, 
		pu.keterangan_program_unggulan, pu.keterangan, 
		pu.tahun_awal, pu.tahun_akhir,
		CASE 
			WHEN EXISTS (
				SELECT 1 FROM tb_keterangan_tagging_program_unggulan ktpu
				JOIN tb_tagging_pokin tp ON ktpu.id_tagging = tp.id
				WHERE ktpu.kode_program_unggulan = pu.kode_program_unggulan
			) THEN TRUE ELSE FALSE 
		END as is_active
	FROM tb_program_unggulan pu
	WHERE pu.tahun_awal >= ? AND pu.tahun_akhir <= ?`
	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir)
	if err != nil {
		return []domain.ProgramUnggulan{}, err
	}
	defer rows.Close()
	var programUnggulanList []domain.ProgramUnggulan
	var kodes []string // kumpulkan semua kode untuk batch berikutnya
	for rows.Next() {
		var pu domain.ProgramUnggulan
		err = rows.Scan(
			&pu.Id, &pu.NamaTagging, &pu.KodeProgramUnggulan,
			&pu.KeteranganProgramUnggulan, &pu.Keterangan,
			&pu.TahunAwal, &pu.TahunAkhir, &pu.IsActive,
		)
		if err != nil {
			return []domain.ProgramUnggulan{}, err
		}
		programUnggulanList = append(programUnggulanList, pu)
		kodes = append(kodes, pu.KodeProgramUnggulan)
	}
	if len(programUnggulanList) == 0 {
		return programUnggulanList, nil
	}
	// === QUERY 2 & 3: Batch enrichment ===
	opdMap, err := repository.FindOpdBatchByKodes(ctx, tx, kodes)
	if err != nil {
		return []domain.ProgramUnggulan{}, err
	}
	tahunMap, err := repository.FindTahunTerpakaiBatch(ctx, tx, kodes)
	if err != nil {
		return []domain.ProgramUnggulan{}, err
	}
	// === Gabungkan hasil di memory ===
	for i, pu := range programUnggulanList {
		programUnggulanList[i].OpdList = opdMap[pu.KodeProgramUnggulan]
		programUnggulanList[i].TahunTerpakai = tahunMap[pu.KodeProgramUnggulan]
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

func (repository *ProgramUnggulanRepositoryImpl) FindUnusedByTahun(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ProgramUnggulan, error) {
	script := `
        SELECT DISTINCT 
            pu.id, 
            pu.nama_tagging, 
            pu.kode_program_unggulan, 
            pu.keterangan_program_unggulan, 
            pu.keterangan, 
            pu.tahun_awal, 
            pu.tahun_akhir
        FROM 
            tb_program_unggulan pu
        WHERE 
            ? BETWEEN pu.tahun_awal AND pu.tahun_akhir
            AND NOT EXISTS (
                SELECT 1 
                FROM tb_keterangan_tagging_program_unggulan ktpu 
                WHERE ktpu.kode_program_unggulan = pu.kode_program_unggulan 
                AND ktpu.tahun = ?
            )
        ORDER BY pu.tahun_awal ASC`

	rows, err := tx.QueryContext(ctx, script, tahun, tahun)
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

func (repository *ProgramUnggulanRepositoryImpl) FindByIdTerkait(ctx context.Context, tx *sql.Tx, ids []int) ([]domain.ProgramUnggulan, error) {
	if len(ids) == 0 {
		return []domain.ProgramUnggulan{}, errors.New("ids tidak boleh kosong")
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT id, nama_tagging, kode_program_unggulan, 
		       keterangan_program_unggulan, keterangan, tahun_awal, tahun_akhir 
		FROM tb_program_unggulan 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.ProgramUnggulan{}, err
	}
	defer rows.Close()

	var result []domain.ProgramUnggulan
	for rows.Next() {
		var programUnggulan domain.ProgramUnggulan
		err := rows.Scan(
			&programUnggulan.Id,
			&programUnggulan.NamaTagging,
			&programUnggulan.KodeProgramUnggulan,
			&programUnggulan.KeteranganProgramUnggulan,
			&programUnggulan.Keterangan,
			&programUnggulan.TahunAwal,
			&programUnggulan.TahunAkhir,
		)
		if err != nil {
			return []domain.ProgramUnggulan{}, err
		}
		result = append(result, programUnggulan)
	}

	if err = rows.Err(); err != nil {
		return []domain.ProgramUnggulan{}, err
	}

	return result, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindOpdByKodeProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) ([]domain.OpdProgramUnggulan, error) {
	script := `
        SELECT opp.kode_opd, tod.nama_opd
        FROM tb_opd_program_prioritas opp
        JOIN tb_operasional_daerah tod ON opp.kode_opd = tod.kode_opd
        WHERE opp.kode_program_unggulan = ?`

	rows, err := tx.QueryContext(ctx, script, kodeProgramUnggulan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var opdList []domain.OpdProgramUnggulan
	for rows.Next() {
		var opd domain.OpdProgramUnggulan
		err = rows.Scan(&opd.KodeOpd, &opd.NamaOpd)
		if err != nil {
			return nil, err
		}
		opdList = append(opdList, opd)
	}
	return opdList, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindTahunTerpakai(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string) ([]string, error) {
	script := `
        SELECT DISTINCT ktpu.tahun
        FROM tb_keterangan_tagging_program_unggulan ktpu
        WHERE ktpu.kode_program_unggulan = ?
        ORDER BY ktpu.tahun ASC`
	rows, err := tx.QueryContext(ctx, script, kodeProgramUnggulan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tahunList []string
	for rows.Next() {
		var tahun string
		err = rows.Scan(&tahun)
		if err != nil {
			return nil, err
		}
		tahunList = append(tahunList, tahun)
	}
	return tahunList, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindProgramUnggulanByKodesBatch(ctx context.Context, tx *sql.Tx, kodes []string) (map[string]*domain.ProgramUnggulan, error) {
	if len(kodes) == 0 {
		return make(map[string]*domain.ProgramUnggulan), nil
	}

	placeholders := make([]string, len(kodes))
	args := make([]interface{}, len(kodes))
	for i, kode := range kodes {
		placeholders[i] = "?"
		args[i] = kode
	}

	script := fmt.Sprintf(`
		SELECT id, kode_program_unggulan, keterangan_program_unggulan, tahun_awal, tahun_akhir
		FROM tb_program_unggulan
		WHERE kode_program_unggulan IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domain.ProgramUnggulan)
	for rows.Next() {
		var program domain.ProgramUnggulan
		err := rows.Scan(&program.Id, &program.KodeProgramUnggulan, &program.KeteranganProgramUnggulan, &program.TahunAwal, &program.TahunAkhir)
		if err != nil {
			return nil, err
		}
		result[program.KodeProgramUnggulan] = &program
	}

	return result, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindOpdBatchByKodes(ctx context.Context, tx *sql.Tx, kodes []string) (map[string][]domain.OpdProgramUnggulan, error) {
	if len(kodes) == 0 {
		return make(map[string][]domain.OpdProgramUnggulan), nil
	}
	placeholders := make([]string, len(kodes))
	args := make([]interface{}, len(kodes))
	for i, kode := range kodes {
		placeholders[i] = "?"
		args[i] = kode
	}
	script := fmt.Sprintf(`
        SELECT opp.id, opp.kode_program_unggulan, opp.kode_opd, tod.nama_opd
        FROM tb_opd_program_prioritas opp
        JOIN tb_operasional_daerah tod ON opp.kode_opd = tod.kode_opd
        WHERE opp.kode_program_unggulan IN (%s)
        ORDER BY opp.kode_program_unggulan`, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// key: kode_program_unggulan → []OpdProgramUnggulan
	result := make(map[string][]domain.OpdProgramUnggulan)
	for rows.Next() {
		var kodeProg string
		var opd domain.OpdProgramUnggulan
		err = rows.Scan(&opd.Id, &kodeProg, &opd.KodeOpd, &opd.NamaOpd)
		if err != nil {
			return nil, err
		}
		result[kodeProg] = append(result[kodeProg], opd)
	}
	return result, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindTahunTerpakaiBatch(ctx context.Context, tx *sql.Tx, kodes []string) (map[string][]string, error) {
	if len(kodes) == 0 {
		return make(map[string][]string), nil
	}
	placeholders := make([]string, len(kodes))
	args := make([]interface{}, len(kodes))
	for i, kode := range kodes {
		placeholders[i] = "?"
		args[i] = kode
	}
	script := fmt.Sprintf(`
        SELECT kode_program_unggulan, tahun
        FROM tb_keterangan_tagging_program_unggulan
        WHERE kode_program_unggulan IN (%s)
        GROUP BY kode_program_unggulan, tahun
        ORDER BY kode_program_unggulan, tahun ASC`, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// key: kode_program_unggulan → []tahun
	result := make(map[string][]string)
	for rows.Next() {
		var kodeProg, tahun string
		err = rows.Scan(&kodeProg, &tahun)
		if err != nil {
			return nil, err
		}
		result[kodeProg] = append(result[kodeProg], tahun)
	}
	return result, nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindOpdByKodeProgramUnggulanAndKodeOpds(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string, kodeOpds []string) ([]domain.OpdProgramUnggulan, error) {
	if len(kodeOpds) == 0 {
		return []domain.OpdProgramUnggulan{}, nil
	}
	placeholders := make([]string, len(kodeOpds))
	args := make([]interface{}, 0, len(kodeOpds)+1)
	args = append(args, kodeProgramUnggulan)
	for i, kodeOpd := range kodeOpds {
		placeholders[i] = "?"
		args = append(args, kodeOpd)
	}
	script := fmt.Sprintf(`
		SELECT opp.id, opp.kode_opd, tod.nama_opd
		FROM tb_opd_program_prioritas opp
		JOIN tb_operasional_daerah tod ON opp.kode_opd = tod.kode_opd
		WHERE opp.kode_program_unggulan = ?
		  AND opp.kode_opd IN (%s)`, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var opdList []domain.OpdProgramUnggulan
	for rows.Next() {
		var opd domain.OpdProgramUnggulan
		if err = rows.Scan(&opd.Id, &opd.KodeOpd, &opd.NamaOpd); err != nil {
			return nil, err
		}
		opdList = append(opdList, opd)
	}
	return opdList, nil
}

func (repository *ProgramUnggulanRepositoryImpl) CreateOpdProgramUnggulan(ctx context.Context, tx *sql.Tx, kodeProgramUnggulan string, kodeOpd []string) error {
	if len(kodeOpd) == 0 {
		return nil
	}
	query := `INSERT INTO tb_opd_program_prioritas (kode_program_unggulan, kode_opd) VALUES `
	var placeholders []string
	var values []any
	for _, kodeOpd := range kodeOpd {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, kodeProgramUnggulan, kodeOpd)
	}
	query += strings.Join(placeholders, ",")
	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("gagal batch insert opd program unggulan: %w", err)
	}
	return nil
}

func (repository *ProgramUnggulanRepositoryImpl) FindOpdProgramUnggulanById(ctx context.Context, tx *sql.Tx, id int) ([]domain.OpdProgramUnggulan, error) {
	script := `SELECT id, kode_program_unggulan, kode_opd FROM tb_opd_program_prioritas WHERE id = ?`
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var opdList []domain.OpdProgramUnggulan
	for rows.Next() {
		var opd domain.OpdProgramUnggulan
		err = rows.Scan(&opd.Id, &opd.KodeProgramUnggulan, &opd.KodeOpd)
		if err != nil {
			return nil, err
		}
		opdList = append(opdList, opd)
	}
	return opdList, nil
}

func (repository *ProgramUnggulanRepositoryImpl) DeleteOpdProgramUnggulan(ctx context.Context, tx *sql.Tx, id int) error {
	script := `DELETE FROM tb_opd_program_prioritas WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, id)
	return err
}

func (repository *ProgramUnggulanRepositoryImpl) FindByTahunAndKodeOpd(ctx context.Context, tx *sql.Tx, tahun string, kodeOpd string) ([]domain.ProgramUnggulan, error) {
	script := `
        SELECT DISTINCT pu.id, pu.nama_tagging, pu.kode_program_unggulan, 
               pu.keterangan_program_unggulan, pu.keterangan, pu.tahun_awal, pu.tahun_akhir 
        FROM tb_program_unggulan pu
        JOIN tb_opd_program_prioritas opp ON pu.kode_program_unggulan = opp.kode_program_unggulan
        WHERE ? BETWEEN pu.tahun_awal AND pu.tahun_akhir
          AND opp.kode_opd = ?
        ORDER BY pu.tahun_awal ASC`
	rows, err := tx.QueryContext(ctx, script, tahun, kodeOpd)

	var programUnggulanList []domain.ProgramUnggulan
	for rows.Next() {
		var programUnggulan domain.ProgramUnggulan
		err = rows.Scan(&programUnggulan.Id, &programUnggulan.NamaTagging, &programUnggulan.KodeProgramUnggulan, &programUnggulan.KeteranganProgramUnggulan, &programUnggulan.Keterangan, &programUnggulan.TahunAwal, &programUnggulan.TahunAkhir)
		if err != nil {
			return nil, err
		}
		programUnggulanList = append(programUnggulanList, programUnggulan)
	}
	return programUnggulanList, nil
}
