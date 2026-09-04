package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
)

type PptkRepositoryImpl struct {
}

func NewPptkRepositoryImpl() *PptkRepositoryImpl {
	return &PptkRepositoryImpl{}
}

func (repository *PptkRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error) {
	script := `INSERT INTO tb_pptk 
	(nip, kode_opd, tahun, kode_sub_kegiatan, nip_atasan, nonaktif_at) 
	VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pptk.Nip,
		pptk.KodeOpd,
		pptk.Tahun,
		pptk.KodeSubKegiatan,
		pptk.NipAtasan,
		pptk.NonAktifAt)
	if err != nil {
		return domain.Pptk{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Pptk{}, err
	}
	pptk.Id = int(id)

	return pptk, nil
}

func (repository *PptkRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pptk domain.Pptk) (domain.Pptk, error) {
	script := `UPDATE tb_pptk SET 
			   nip = ?, 
			   kode_opd = ?, 
			   tahun = ?, 
			   kode_sub_kegiatan = ?, 
			   nip_atasan = ?,
			   nonaktif_at = ?
			   WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, 
		pptk.Nip, 
		pptk.KodeOpd, 
		pptk.Tahun, 
		pptk.KodeSubKegiatan, 
		pptk.NipAtasan, 
		pptk.NonAktifAt, 
		pptk.Id)
	if err != nil {
		return domain.Pptk{}, err
	}
	return pptk, nil
}

func (repository *PptkRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_pptk WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PptkRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pptk, error) {
	script := `SELECT 
            tp.id, 
			tp.nip, 
			peg.nama, 
			tp.kode_opd, 
			tp.tahun, 
			tp.kode_sub_kegiatan, 
			tp.nip_atasan, 
			tpa.nama, 
			tp.aktif_at, 
			tp.nonaktif_at
        FROM tb_pptk tp
		LEFT JOIN tb_pegawai peg ON tp.nip = peg.nip
		LEFT JOIN tb_pegawai tpa ON tp.nip_atasan = tpa.nip
        WHERE tp.id = ?`
	var Pptk domain.Pptk
	err := tx.QueryRowContext(ctx, script, id).Scan(
		&Pptk.Id,
		&Pptk.Nip,
		&Pptk.NamaPegawai,
		&Pptk.KodeOpd,
		&Pptk.Tahun,
		&Pptk.KodeSubKegiatan,
		&Pptk.NipAtasan,
		&Pptk.NamaAtasan,
		&Pptk.AktifAt,
		&Pptk.NonAktifAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Pptk{}, errors.New("pptk tidak ditemukan")
		}
		return domain.Pptk{}, err
	}
	return Pptk, nil
}

func (repository *PptkRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx,kodeSubkegiatan string, kodeOpd string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            tp.id, 
			tp.nip, 
			peg.nama, 
			tp.kode_opd, 
			tp.tahun, 
			tp.kode_sub_kegiatan, 
			tp.nip_atasan, 
			tpa.nama, 
			tp.aktif_at, 
			tp.nonaktif_at
        FROM tb_pptk tp
		LEFT JOIN tb_pegawai peg ON tp.nip = peg.nip
		LEFT JOIN tb_pegawai tpa ON tp.nip_atasan = tpa.nip
        WHERE tp.kode_sub_kegiatan = ? AND tp.kode_opd = ? AND tp.tahun = ?`

	rows, err := tx.QueryContext(ctx, script, kodeSubkegiatan, kodeOpd, tahun)
	if err != nil {
		return []domain.Pptk{}, err
	}
	defer rows.Close()
	// log.Printf("TahunAWal: %s TAhun Akhir: %s", tahunAwal, tahunAkhir)

	var pptkList []domain.Pptk
	for rows.Next() {
		var pptk domain.Pptk
		err = rows.Scan(
			&pptk.Id,
			&pptk.Nip,
			&pptk.NamaPegawai,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.NamaAtasan,
			&pptk.AktifAt,
			&pptk.NonAktifAt,
		)
		if err != nil {
			return []domain.Pptk{}, err
		}
		pptkList = append(pptkList, pptk)
	}
	return pptkList, nil
}
func (repository *PptkRepositoryImpl) FindAllByNip(ctx context.Context, tx *sql.Tx, kodeSubkegiatan string, pegawaiId string, tahun string) ([]domain.Pptk, error) {
	// Query untuk mengambil program unggulan beserta status aktifnya
	script := `
        SELECT 
            tp.id, 
			tp.nip, 
			peg.nama, 
			tp.kode_opd, 
			tp.tahun, 
			tp.kode_sub_kegiatan, 
			tp.nip_atasan, 
			tpa.nama, 
			tp.aktif_at, 
			tp.nonaktif_at
        FROM tb_pptk tp
		LEFT JOIN tb_pegawai peg ON tp.nip = peg.nip
		LEFT JOIN tb_pegawai tpa ON tp.nip_atasan = tpa.nip
        WHERE tp.kode_sub_kegiatan = ? AND tp.nip = ? AND tp.tahun = ?`

	rows, err := tx.QueryContext(ctx, script, kodeSubkegiatan, pegawaiId, tahun)
	if err != nil {
		return []domain.Pptk{}, err
	}
	defer rows.Close()
	// log.Printf("TahunAWal: %s TAhun Akhir: %s", tahunAwal, tahunAkhir)

	var pptkList []domain.Pptk
	for rows.Next() {
		var pptk domain.Pptk
		err = rows.Scan(
			&pptk.Id,
			&pptk.Nip,
			&pptk.NamaPegawai,
			&pptk.KodeOpd,
			&pptk.Tahun,
			&pptk.KodeSubKegiatan,
			&pptk.NipAtasan,
			&pptk.NamaAtasan,
			&pptk.AktifAt,
			&pptk.NonAktifAt,
		)
		if err != nil {
			return []domain.Pptk{}, err
		}
		pptkList = append(pptkList, pptk)
	}
	return pptkList, nil
}

func (repository *PptkRepositoryImpl) KandidatPptkOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.KandidatPptk, error) {

	query := `
	WITH rencana_kinerja_opd AS (
	   SELECT DISTINCT
		   rk.id as rekin_id,
		   rk.pegawai_id,
		   p.nama as nama_pegawai,

		   (SELECT st2.kode_subkegiatan 
			FROM tb_subkegiatan_terpilih st2 
			WHERE st2.rekin_id = rk.id 
			LIMIT 1) as kode_subkegiatan,

		   (SELECT sk2.nama_subkegiatan 
			FROM tb_subkegiatan_terpilih st2 
			LEFT JOIN tb_subkegiatan sk2 
				ON sk2.kode_subkegiatan = st2.kode_subkegiatan 
			WHERE st2.rekin_id = rk.id 
			LIMIT 1) as nama_subkegiatan,

		   rk.nama_rencana_kinerja
	   FROM tb_rencana_kinerja rk
	   LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
	   WHERE rk.kode_opd = ? 
	   AND rk.tahun = ?
	   AND EXISTS (
		   SELECT 1 FROM tb_subkegiatan_terpilih st 
		   WHERE st.rekin_id = rk.id 
		   AND st.kode_subkegiatan IS NOT NULL 
		   AND st.kode_subkegiatan != ''
	   )
   )
   SELECT 
	   rkp.pegawai_id,
	   rkp.nama_pegawai,
	   (
		   SELECT r.role
		   FROM tb_users u
		   LEFT JOIN tb_user_role ur ON ur.user_id = u.id
		   LEFT JOIN tb_role r ON r.id = ur.role_id
		   WHERE u.nip = rkp.pegawai_id
		   LIMIT 1
	   ) as level,
	   rkp.kode_subkegiatan,
	   rkp.nama_subkegiatan,
	   rkp.rekin_id,
	   rkp.nama_rencana_kinerja,
	   ra.id as renaksi_id,
	   ra.nama_rencana_aksi,
	   COALESCE(ra.urutan, 999) as urutan,
	   COALESCE(SUM(rb.anggaran), 0) as anggaran
   FROM rencana_kinerja_opd rkp
   LEFT JOIN tb_rencana_aksi ra 
	   ON ra.rencana_kinerja_id = rkp.rekin_id
   LEFT JOIN tb_rincian_belanja rb 
	   ON rb.renaksi_id = ra.id

   GROUP BY rkp.pegawai_id, rkp.nama_pegawai, rkp.kode_subkegiatan, rkp.nama_subkegiatan,
            rkp.rekin_id, rkp.nama_rencana_kinerja, ra.id, ra.nama_rencana_aksi, ra.urutan

   ORDER BY rkp.kode_subkegiatan, rkp.rekin_id, ra.urutan ASC, ra.id
	`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying kandidat pptk opd: %v", err)
	}
	defer rows.Close()

	var result []domain.KandidatPptk

	// Tracking pegawai agar tidak duplicate
	pegawaiTracker := make(map[string]bool)

	for rows.Next() {
		var (
			pegawaiId, namaPegawai string
			level                 sql.NullString
			kodeSubkegiatan       string
			namaSubkegiatan       string
			rekinId               string
			namaRencanaKinerja    string
			renaksiId             sql.NullString
			namaRenaksi           sql.NullString
			urutan                sql.NullInt64
			anggaran              int64
		)

		err := rows.Scan(
			&pegawaiId,
			&namaPegawai,
			&level,
			&kodeSubkegiatan,
			&namaSubkegiatan,
			&rekinId,
			&namaRencanaKinerja,
			&renaksiId,
			&namaRenaksi,
			&urutan,
			&anggaran,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning kandidat pptk opd: %v", err)
		}

		// Hanya satu kandidat untuk setiap pegawai
		if !pegawaiTracker[pegawaiId] {
			result = append(result, domain.KandidatPptk{
				PegawaiId:   pegawaiId,
				NamaPegawai: namaPegawai,
				Level:       level.String,
			})

			pegawaiTracker[pegawaiId] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kandidat pptk opd: %v", err)
	}

	return result, nil
}

func (repository *PptkRepositoryImpl) KandidatPptkPegawai(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.KandidatPptk, error) {

	query := `
    WITH pegawai_rencana AS (
        SELECT DISTINCT st.kode_subkegiatan, rk.kode_opd
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_subkegiatan_terpilih st 
            ON st.rekin_id = rk.id
        WHERE rk.pegawai_id = ? 
          AND rk.tahun = ?
    ),
    related_rencana AS (
        SELECT DISTINCT
            rk.id as rekin_id,
            rk.pegawai_id,
            rk.kode_opd,
            p.nama as nama_pegawai,
            st.kode_subkegiatan,
            sk.nama_subkegiatan,
            rk.nama_rencana_kinerja
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_subkegiatan_terpilih st 
            ON st.rekin_id = rk.id
        INNER JOIN pegawai_rencana pr 
            ON pr.kode_subkegiatan = st.kode_subkegiatan
            AND pr.kode_opd = rk.kode_opd
        LEFT JOIN tb_pegawai p 
            ON p.nip = rk.pegawai_id
        LEFT JOIN tb_subkegiatan sk 
            ON sk.kode_subkegiatan = st.kode_subkegiatan
        WHERE rk.tahun = ?
    )
    SELECT
        rr.pegawai_id,
        rr.nama_pegawai,
        (
            SELECT r.role
            FROM tb_users u
            LEFT JOIN tb_user_role ur 
                ON ur.user_id = u.id
            LEFT JOIN tb_role r 
                ON r.id = ur.role_id
            WHERE u.nip = rr.pegawai_id
            LIMIT 1
        ) as level,
        rr.kode_opd,
        rr.kode_subkegiatan,
        rr.rekin_id,
        rr.nama_rencana_kinerja
    FROM related_rencana rr
    GROUP BY
        rr.pegawai_id,
        rr.nama_pegawai,
        rr.kode_opd,
        rr.kode_subkegiatan,
        rr.rekin_id,
        rr.nama_rencana_kinerja
    ORDER BY
        rr.kode_opd,
        rr.kode_subkegiatan,
        rr.pegawai_id,
        rr.rekin_id
    `

	rows, err := tx.QueryContext(
		ctx,
		query,
		pegawaiId,
		tahun,
		tahun,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying kandidat pptk pegawai: %v", err)
	}
	defer rows.Close()

	var result []domain.KandidatPptk
	pegawaiTracker := make(map[string]bool)

	for rows.Next() {
		var (
			pegawaiIdResult    string
			namaPegawai        string
			level              sql.NullString
			kodeOpd            string
			kodeSubkegiatan    string
			rekinId            string
			namaRencanaKinerja string
		)

		err := rows.Scan(
			&pegawaiIdResult,
			&namaPegawai,
			&level,
			&kodeOpd,
			&kodeSubkegiatan,
			&rekinId,
			&namaRencanaKinerja,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"error scanning kandidat pptk pegawai: %v",
				err,
			)
		}

		if !pegawaiTracker[pegawaiIdResult] {
			result = append(result, domain.KandidatPptk{
				PegawaiId:   pegawaiIdResult,
				NamaPegawai: namaPegawai,
				Level:       level.String,
			})

			pegawaiTracker[pegawaiIdResult] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error iterating kandidat pptk pegawai: %v",
			err,
		)
	}

	return result, nil
}