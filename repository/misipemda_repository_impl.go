package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"math/rand"
	"strings"
	"time"
)

type MisiPemdaRepositoryImpl struct{}

func NewMisiPemdaRepositoryImpl() *MisiPemdaRepositoryImpl {
	return &MisiPemdaRepositoryImpl{}
}

func (repository *MisiPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	visiPemda.Id = r.Intn(9000) + 1000

	script := "INSERT INTO tb_misi_pemda (id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, visiPemda.Id, visiPemda.IdVisi, visiPemda.Misi, visiPemda.Urutan, visiPemda.TahunAwalPeriode, visiPemda.TahunAkhirPeriode, visiPemda.JenisPeriode, visiPemda.Keterangan)
	if err != nil {
		return domain.MisiPemda{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.MisiPemda{}, err
	}
	visiPemda.Id = int(id)

	return visiPemda, nil
}

func (repository *MisiPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error) {
	script := "UPDATE tb_misi_pemda SET id_visi = ?, misi = ?, urutan = ?, tahun_awal_periode = ?, tahun_akhir_periode = ?, jenis_periode = ?, keterangan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, visiPemda.IdVisi, visiPemda.Misi, visiPemda.Urutan, visiPemda.TahunAwalPeriode, visiPemda.TahunAkhirPeriode, visiPemda.JenisPeriode, visiPemda.Keterangan, visiPemda.Id)
	if err != nil {
		return domain.MisiPemda{}, err
	}

	return visiPemda, nil
}

func (repository *MisiPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, visiPemdaId int) error {
	script := "DELETE FROM tb_misi_pemda WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, visiPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *MisiPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error) {
	script := "SELECT id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan FROM tb_misi_pemda WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, visiPemdaId)
	if err != nil {
		return domain.MisiPemda{}, err
	}
	defer rows.Close()

	var visiPemda domain.MisiPemda
	if rows.Next() {
		err := rows.Scan(
			&visiPemda.Id,
			&visiPemda.IdVisi,
			&visiPemda.Misi,
			&visiPemda.Urutan,
			&visiPemda.TahunAwalPeriode,
			&visiPemda.TahunAkhirPeriode,
			&visiPemda.JenisPeriode,
			&visiPemda.Keterangan,
		)
		helper.PanicIfError(err)

	}

	return visiPemda, nil
}

func (repository *MisiPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.MisiPemda, error) {
	var conditions []string
	var params []interface{}

	baseQuery := `SELECT id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, 
                  jenis_periode, keterangan FROM tb_misi_pemda`

	// Membangun query dinamis berdasarkan filter yang ada
	if tahunAwal != "" && tahunAkhir != "" {
		conditions = append(conditions, "CAST(? AS SIGNED) BETWEEN CAST(tahun_awal_periode AS SIGNED) AND CAST(tahun_akhir_periode AS SIGNED)")
		params = append(params, tahunAwal)
	}

	if jenisPeriode != "" {
		conditions = append(conditions, "jenis_periode = ?")
		params = append(params, jenisPeriode)
	}

	// Menggabungkan kondisi WHERE jika ada
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Tambahkan ORDER BY untuk mengurutkan berdasarkan id_visi dan urutan
	baseQuery += " ORDER BY id_visi, urutan"

	rows, err := tx.QueryContext(ctx, baseQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var misiPemdaList []domain.MisiPemda
	for rows.Next() {
		var misiPemda domain.MisiPemda
		err := rows.Scan(
			&misiPemda.Id,
			&misiPemda.IdVisi,
			&misiPemda.Misi,
			&misiPemda.Urutan,
			&misiPemda.TahunAwalPeriode,
			&misiPemda.TahunAkhirPeriode,
			&misiPemda.JenisPeriode,
			&misiPemda.Keterangan,
		)
		if err != nil {
			return nil, err
		}
		misiPemdaList = append(misiPemdaList, misiPemda)
	}

	return misiPemdaList, nil
}

func (repository *MisiPemdaRepositoryImpl) FindByIdWithDefault(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error) {
	if visiPemdaId == 0 {
		return domain.MisiPemda{
			Id:                0,
			IdVisi:            0,
			Misi:              "",
			Urutan:            0,
			TahunAwalPeriode:  "",
			TahunAkhirPeriode: "",
			JenisPeriode:      "",
			Keterangan:        "",
		}, nil
	}

	return repository.FindById(ctx, tx, visiPemdaId)
}

func (repository *MisiPemdaRepositoryImpl) CheckUrutanExists(ctx context.Context, tx *sql.Tx, idVisi int, urutan int) (bool, error) {
	// Ubah query untuk hanya mengecek urutan yang sama pada id_visi yang sama
	script := "SELECT EXISTS(SELECT 1 FROM tb_misi_pemda WHERE id_visi = ? AND urutan = ?)"
	var exists bool
	err := tx.QueryRowContext(ctx, script, idVisi, urutan).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repository *MisiPemdaRepositoryImpl) CheckUrutanExistsExcept(ctx context.Context, tx *sql.Tx, idVisi int, urutan int, id int) (bool, error) {
	// Ubah query untuk hanya mengecek urutan yang sama pada id_visi yang sama, kecuali untuk id yang sedang diupdate
	script := "SELECT EXISTS(SELECT 1 FROM tb_misi_pemda WHERE id_visi = ? AND urutan = ? AND id != ?)"
	var exists bool
	err := tx.QueryRowContext(ctx, script, idVisi, urutan, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repository *MisiPemdaRepositoryImpl) FindByIdVisi(ctx context.Context, tx *sql.Tx, idVisi int) ([]domain.MisiPemda, error) {
	script := `SELECT id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, 
               jenis_periode, keterangan 
               FROM tb_misi_pemda 
               WHERE id_visi = ? 
               ORDER BY urutan`

	rows, err := tx.QueryContext(ctx, script, idVisi)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var misiPemdaList []domain.MisiPemda
	for rows.Next() {
		var misiPemda domain.MisiPemda
		err := rows.Scan(
			&misiPemda.Id,
			&misiPemda.IdVisi,
			&misiPemda.Misi,
			&misiPemda.Urutan,
			&misiPemda.TahunAwalPeriode,
			&misiPemda.TahunAkhirPeriode,
			&misiPemda.JenisPeriode,
			&misiPemda.Keterangan,
		)
		if err != nil {
			return nil, err
		}
		misiPemdaList = append(misiPemdaList, misiPemda)
	}

	return misiPemdaList, nil
}

func (repository *MisiPemdaRepositoryImpl) FindByTematikId(ctx context.Context, tx *sql.Tx, tematikId int) ([]domain.MisiPemda, error) {
	script := `
		SELECT DISTINCT
			mp.id,
			mp.id_visi,
			COALESCE(vp.visi, '') AS visi,
			mp.misi,
			mp.urutan,
			mp.tahun_awal_periode,
			mp.tahun_akhir_periode,
			mp.jenis_periode,
			COALESCE(mp.keterangan, '') AS keterangan
		FROM tb_tujuan_pemda tp
		INNER JOIN tb_misi_pemda mp ON tp.id_misi = mp.id
		LEFT JOIN tb_visi_pemda vp ON mp.id_visi = vp.id
		WHERE tp.tematik_id = ?
		  AND tp.id_misi > 0
		ORDER BY mp.urutan ASC, mp.id ASC`

	rows, err := tx.QueryContext(ctx, script, tematikId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	misiList := make([]domain.MisiPemda, 0)
	for rows.Next() {
		var misi domain.MisiPemda
		if err := rows.Scan(
			&misi.Id,
			&misi.IdVisi,
			&misi.Visi,
			&misi.Misi,
			&misi.Urutan,
			&misi.TahunAwalPeriode,
			&misi.TahunAkhirPeriode,
			&misi.JenisPeriode,
			&misi.Keterangan,
		); err != nil {
			return nil, err
		}
		misiList = append(misiList, misi)
	}

	return misiList, rows.Err()
}
