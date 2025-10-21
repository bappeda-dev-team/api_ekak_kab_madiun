package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type RincianBelanjaRepositoryImpl struct {
}

func NewRincianBelanjaRepositoryImpl() *RincianBelanjaRepositoryImpl {
	return &RincianBelanjaRepositoryImpl{}
}

func (repository *RincianBelanjaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, rincianBelanja domain.RincianBelanja) (domain.RincianBelanja, error) {
	script := `
		INSERT INTO tb_rincian_belanja (renaksi_id, anggaran)
		VALUES (?, ?)
	`
	result, err := tx.ExecContext(ctx, script, rincianBelanja.RenaksiId, rincianBelanja.Anggaran)
	if err != nil {
		return domain.RincianBelanja{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.RincianBelanja{}, err
	}
	rincianBelanja.Id = int(id)
	return rincianBelanja, nil
}

func (repository *RincianBelanjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, rincianBelanja domain.RincianBelanja) (domain.RincianBelanja, error) {
	script := `
		UPDATE tb_rincian_belanja SET  anggaran = ? WHERE renaksi_id = ?
	`
	_, err := tx.ExecContext(ctx, script, rincianBelanja.Anggaran, rincianBelanja.RenaksiId)
	if err != nil {
		return domain.RincianBelanja{}, err
	}
	return rincianBelanja, nil
}

func (repository *RincianBelanjaRepositoryImpl) FindByRenaksiId(ctx context.Context, tx *sql.Tx, renaksiId string) (domain.RincianBelanja, error) {
	script := `
		SELECT renaksi_id, anggaran FROM tb_rincian_belanja WHERE renaksi_id = ?
	`
	rows, err := tx.QueryContext(ctx, script, renaksiId)
	if err != nil {
		return domain.RincianBelanja{}, err
	}
	defer rows.Close()

	var rincianBelanja domain.RincianBelanja
	if rows.Next() {
		err = rows.Scan(&rincianBelanja.RenaksiId, &rincianBelanja.Anggaran)
		if err != nil {
			return domain.RincianBelanja{}, err
		}
	}
	return rincianBelanja, nil
}

func (repository *RincianBelanjaRepositoryImpl) FindRincianBelanjaAsn(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.RincianBelanjaAsn, error) {
	query := `
        WITH rencana_kinerja_pegawai AS (
            SELECT 
                rk.id as rekin_id,
                rk.pegawai_id,
                p.nama as nama_pegawai,
                st.kode_subkegiatan,
                sk.nama_subkegiatan,
                rk.nama_rencana_kinerja
            FROM tb_rencana_kinerja rk
            LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
            INNER JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
            LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
            WHERE rk.pegawai_id = ? 
            AND rk.tahun = ?
            AND st.kode_subkegiatan IS NOT NULL
            AND st.kode_subkegiatan != ''
        )
        SELECT 
            rkp.pegawai_id,
            rkp.nama_pegawai,
            rkp.kode_subkegiatan,
            rkp.nama_subkegiatan,
            rkp.rekin_id,
            rkp.nama_rencana_kinerja,
            ra.id as renaksi_id,
            ra.nama_rencana_aksi,
            COALESCE(ra.urutan, 999) as urutan,
            COALESCE(SUM(rb.anggaran), 0) as anggaran
        FROM rencana_kinerja_pegawai rkp
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rkp.rekin_id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        GROUP BY rkp.pegawai_id, rkp.nama_pegawai, rkp.kode_subkegiatan, rkp.nama_subkegiatan,
                 rkp.rekin_id, rkp.nama_rencana_kinerja, ra.id, ra.nama_rencana_aksi, ra.urutan
        ORDER BY rkp.kode_subkegiatan, rkp.rekin_id, ra.urutan ASC, ra.id ASC
    `

	rows, err := tx.QueryContext(ctx, query, pegawaiId, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rincian belanja asn: %v", err)
	}
	defer rows.Close()

	var result []domain.RincianBelanjaAsn
	var currentSubkegiatan *domain.RincianBelanjaAsn
	var currentRencanaKinerja *domain.RencanaKinerjaAsn
	renaksiMap := make(map[string]bool)

	for rows.Next() {
		var (
			pegawaiId, namaPegawai, kodeSubkegiatan, namaSubkegiatan string
			rekinId, namaRencanaKinerja                              string
			renaksiId, namaRenaksi                                   sql.NullString
			urutan                                                   sql.NullInt64
			anggaran                                                 int64
		)

		err := rows.Scan(
			&pegawaiId,
			&namaPegawai,
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
			return nil, fmt.Errorf("error scanning rincian belanja asn: %v", err)
		}

		// Jika subkegiatan baru
		if currentSubkegiatan == nil || currentSubkegiatan.KodeSubkegiatan != kodeSubkegiatan {
			if currentSubkegiatan != nil {
				result = append(result, *currentSubkegiatan)
			}
			currentSubkegiatan = &domain.RincianBelanjaAsn{
				PegawaiId:       pegawaiId,
				NamaPegawai:     namaPegawai,
				KodeSubkegiatan: kodeSubkegiatan,
				NamaSubkegiatan: namaSubkegiatan,
				TotalAnggaran:   0,
				RencanaKinerja:  []domain.RencanaKinerjaAsn{},
			}
			currentRencanaKinerja = nil
			renaksiMap = make(map[string]bool)
		}

		// Jika rencana kinerja baru
		if currentRencanaKinerja == nil || currentRencanaKinerja.RencanaKinerjaId != rekinId {
			currentRencanaKinerja = &domain.RencanaKinerjaAsn{
				RencanaKinerjaId: rekinId,
				RencanaKinerja:   namaRencanaKinerja,
				RencanaAksi:      make([]domain.RincianBelanja, 0),
			}
			currentSubkegiatan.RencanaKinerja = append(currentSubkegiatan.RencanaKinerja, *currentRencanaKinerja)
			renaksiMap = make(map[string]bool) // RESET untuk rencana kinerja baru
		}

		// Tambahkan rencana aksi jika ada dan belum duplikat
		if renaksiId.Valid && namaRenaksi.Valid {
			if !renaksiMap[renaksiId.String] {
				rincianBelanja := domain.RincianBelanja{
					RenaksiId: renaksiId.String,
					Renaksi:   namaRenaksi.String,
					Anggaran:  anggaran,
				}
				lastIdx := len(currentSubkegiatan.RencanaKinerja) - 1
				currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi = append(
					currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi,
					rincianBelanja,
				)
				currentSubkegiatan.TotalAnggaran += int(anggaran)
				renaksiMap[renaksiId.String] = true
			}
		}
	}

	if currentSubkegiatan != nil {
		result = append(result, *currentSubkegiatan)
	}

	return result, nil
}

func (repository *RincianBelanjaRepositoryImpl) FindIndikatorByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.Indikator, error) {
	script := `
        SELECT 
            i.id,
            i.rencana_kinerja_id,
            i.indikator,
            t.id as target_id,
            t.indikator_id,
            COALESCE(t.target, '') as target,
            COALESCE(t.satuan, '') as satuan
        FROM tb_indikator i
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE i.rencana_kinerja_id = ?
        ORDER BY i.id`

	rows, err := tx.QueryContext(ctx, script, rekinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId, rekinId, indikator             string
			targetId, indikatorId, target, satuan string
		)

		err := rows.Scan(
			&indId,
			&rekinId,
			&indikator,
			&targetId,
			&indikatorId,
			&target,
			&satuan,
		)
		if err != nil {
			return nil, err
		}

		// Proses Indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:               indId,
				RencanaKinerjaId: rekinId,
				Indikator:        indikator,
				Target:           []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// Proses Target jika ada
		if targetId != "" && indikatorId != "" {
			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      target,
				Satuan:      satuan,
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

func (repository *RincianBelanjaRepositoryImpl) FindIndikatorSubkegiatanByKodeAndOpd(ctx context.Context, tx *sql.Tx, kodeSubkegiatan string, kodeOpd string, tahun string) ([]domain.Indikator, error) {
	script := `
        SELECT 
            i.id,
            i.kode as kode_subkegiatan,
            i.kode_opd,
            i.indikator,
            i.tahun,
            t.id as target_id,
            t.indikator_id,
            COALESCE(t.target, '') as target,
            COALESCE(t.satuan, '') as satuan
        FROM tb_indikator i
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE i.kode = ?
        AND i.kode_opd = ?
        AND i.tahun = ?
        ORDER BY i.id, i.tahun`

	rows, err := tx.QueryContext(ctx, script, kodeSubkegiatan, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indId, kodeSubkegiatan, kodeOpd, indikator, tahun string
			targetId, indikatorId                             sql.NullString // Menggunakan sql.NullString untuk field yang bisa NULL
			target, satuan                                    sql.NullString
		)

		err := rows.Scan(
			&indId,
			&kodeSubkegiatan,
			&kodeOpd,
			&indikator,
			&tahun,
			&targetId,
			&indikatorId,
			&target,
			&satuan,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning indikator: %v", err)
		}

		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:        indId,
				Kode:      kodeSubkegiatan,
				KodeOpd:   kodeOpd,
				Tahun:     tahun,
				Indikator: indikator,
				Target:    []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// Hanya tambahkan target jika targetId valid
		if targetId.Valid && targetId.String != "" {
			targetObj := domain.Target{
				Id:          targetId.String,
				IndikatorId: helper.GetNullStringValue(indikatorId),
				Target:      helper.GetNullStringValue(target),
				Satuan:      helper.GetNullStringValue(satuan),
			}
			ind.Target = append(ind.Target, targetObj)
		}
	}

	var result []domain.Indikator
	for _, ind := range indikatorMap {
		result = append(result, *ind)
	}

	return result, nil
}

func (repository *RincianBelanjaRepositoryImpl) FindAnggaranByRenaksiId(ctx context.Context, tx *sql.Tx, renaksiId string) (domain.RincianBelanja, error) {
	query := `
        SELECT 
            rb.id,
            rb.renaksi_id,
            rb.anggaran,
            ra.nama_rencana_aksi
        FROM tb_rincian_belanja rb
        LEFT JOIN tb_rencana_aksi ra ON ra.id = rb.renaksi_id
        WHERE rb.renaksi_id = ?
    `

	var result domain.RincianBelanja
	row := tx.QueryRowContext(ctx, query, renaksiId)
	err := row.Scan(
		&result.Id,
		&result.RenaksiId,
		&result.Anggaran,
		&result.Renaksi,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.RincianBelanja{}, nil // Return empty struct if not found
		}
		return domain.RincianBelanja{}, err
	}

	return result, nil
}

// laporan rincian belanja
func (repository *RincianBelanjaRepositoryImpl) LaporanRincianBelanjaOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.RincianBelanjaAsn, error) {
	query := `
         WITH rencana_kinerja_opd AS (
            SELECT 
                rk.id as rekin_id,
                rk.pegawai_id,
                p.nama as nama_pegawai,
                st.kode_subkegiatan,
                sk.nama_subkegiatan,
                rk.nama_rencana_kinerja
            FROM tb_rencana_kinerja rk
            LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
            INNER JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
            LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
            WHERE rk.kode_opd = ? 
            AND rk.tahun = ?
            AND st.kode_subkegiatan IS NOT NULL
            AND st.kode_subkegiatan != ''
        )
        SELECT 
            rkp.pegawai_id,
            rkp.nama_pegawai,
            rkp.kode_subkegiatan,
            rkp.nama_subkegiatan,
            rkp.rekin_id,
            rkp.nama_rencana_kinerja,
            ra.id as renaksi_id,
            ra.nama_rencana_aksi,
            COALESCE(rb.anggaran, 0) as anggaran
        FROM rencana_kinerja_opd rkp
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rkp.rekin_id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        ORDER BY rkp.kode_subkegiatan, rkp.rekin_id, ra.id
`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rincian belanja asn: %v", err)
	}
	defer rows.Close()

	var result []domain.RincianBelanjaAsn
	var currentSubkegiatan *domain.RincianBelanjaAsn
	var currentRencanaKinerja *domain.RencanaKinerjaAsn

	for rows.Next() {
		var (
			pegawaiId, namaPegawai, kodeSubkegiatan, namaSubkegiatan string
			rekinId, namaRencanaKinerja                              string
			renaksiId, namaRenaksi                                   sql.NullString
			anggaran                                                 int64
		)

		err := rows.Scan(
			&pegawaiId,
			&namaPegawai,
			&kodeSubkegiatan,
			&namaSubkegiatan,
			&rekinId,
			&namaRencanaKinerja,
			&renaksiId,
			&namaRenaksi,
			&anggaran,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rincian belanja asn: %v", err)
		}

		// Jika subkegiatan baru
		if currentSubkegiatan == nil || currentSubkegiatan.KodeSubkegiatan != kodeSubkegiatan {
			if currentSubkegiatan != nil {
				result = append(result, *currentSubkegiatan)
			}
			currentSubkegiatan = &domain.RincianBelanjaAsn{
				PegawaiId:       pegawaiId,
				NamaPegawai:     namaPegawai,
				KodeSubkegiatan: kodeSubkegiatan,
				NamaSubkegiatan: namaSubkegiatan,
				TotalAnggaran:   0,
				RencanaKinerja:  []domain.RencanaKinerjaAsn{},
			}
			currentRencanaKinerja = nil
		}

		// Jika rencana kinerja baru dalam subkegiatan yang sama
		if currentRencanaKinerja == nil || currentRencanaKinerja.RencanaKinerja != namaRencanaKinerja {
			currentRencanaKinerja = &domain.RencanaKinerjaAsn{
				RencanaKinerjaId: rekinId,
				RencanaKinerja:   namaRencanaKinerja,
				PegawaiId:        pegawaiId,
				NamaPegawai:      namaPegawai,
				RencanaAksi:      make([]domain.RincianBelanja, 0),
			}
			currentSubkegiatan.RencanaKinerja = append(currentSubkegiatan.RencanaKinerja, *currentRencanaKinerja)
		}

		// Tambahkan rencana aksi jika ada
		if renaksiId.Valid && namaRenaksi.Valid {
			rincianBelanja := domain.RincianBelanja{
				RenaksiId: renaksiId.String,
				Renaksi:   namaRenaksi.String,
				Anggaran:  anggaran,
			}
			lastIdx := len(currentSubkegiatan.RencanaKinerja) - 1
			currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi = append(
				currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi,
				rincianBelanja,
			)
			currentSubkegiatan.TotalAnggaran += int(anggaran)
		}
	}

	// Tambahkan subkegiatan terakhir jika ada
	if currentSubkegiatan != nil {
		result = append(result, *currentSubkegiatan)
	}

	return result, nil
}

func (repository *RincianBelanjaRepositoryImpl) LaporanRincianBelanjaPegawai(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.RincianBelanjaAsn, error) {
	opdQuery := `
    SELECT DISTINCT kode_opd 
    FROM tb_rencana_kinerja 
    WHERE pegawai_id = ? AND tahun = ?
    `
	opdRows, err := tx.QueryContext(ctx, opdQuery, pegawaiId, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying opd: %v", err)
	}
	defer opdRows.Close()

	var kodeOpds []string
	for opdRows.Next() {
		var kodeOpd string
		if err := opdRows.Scan(&kodeOpd); err != nil {
			return nil, fmt.Errorf("error scanning opd: %v", err)
		}
		kodeOpds = append(kodeOpds, kodeOpd)
	}

	fmt.Println(kodeOpds)

	query := `
    WITH pegawai_rencana AS (
        SELECT DISTINCT st.kode_subkegiatan, rk.kode_opd
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
        WHERE rk.pegawai_id = ? AND rk.tahun = ?
    ),
    related_rencana AS (
        SELECT 
            rk.id as rekin_id,
            rk.pegawai_id,
            rk.kode_opd,
            p.nama as nama_pegawai,
            st.kode_subkegiatan,
            sk.nama_subkegiatan,
            rk.nama_rencana_kinerja
        FROM tb_rencana_kinerja rk
        INNER JOIN tb_subkegiatan_terpilih st ON st.rekin_id = rk.id
        INNER JOIN pegawai_rencana pr ON pr.kode_subkegiatan = st.kode_subkegiatan 
            AND pr.kode_opd = rk.kode_opd
        LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
        LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
        WHERE rk.tahun = ?
    )
    SELECT 
        rr.pegawai_id,
        rr.nama_pegawai,
        rr.kode_opd,
        rr.kode_subkegiatan,
        rr.nama_subkegiatan,
        rr.rekin_id,
        rr.nama_rencana_kinerja,
        ra.id as renaksi_id,
        ra.nama_rencana_aksi,
        COALESCE(ra.urutan, 999) as urutan,
        COALESCE(SUM(rb.anggaran), 0) as anggaran
    FROM related_rencana rr
    LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rr.rekin_id
    LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
    GROUP BY rr.pegawai_id, rr.nama_pegawai, rr.kode_opd, rr.kode_subkegiatan, 
             rr.nama_subkegiatan, rr.rekin_id, rr.nama_rencana_kinerja, ra.id, ra.nama_rencana_aksi, ra.urutan
    ORDER BY rr.kode_opd, rr.kode_subkegiatan, rr.pegawai_id, rr.rekin_id, ra.urutan ASC, ra.id ASC
    `

	rows, err := tx.QueryContext(ctx, query, pegawaiId, tahun, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rincian belanja pegawai: %v", err)
	}
	defer rows.Close()

	var result []domain.RincianBelanjaAsn
	var currentSubkegiatan *domain.RincianBelanjaAsn
	var currentRencanaKinerja *domain.RencanaKinerjaAsn
	renaksiMap := make(map[string]bool)

	for rows.Next() {
		var (
			pegawaiId, namaPegawai, kodeOpd, kodeSubkegiatan, namaSubkegiatan string
			rekinId, namaRencanaKinerja                                       string
			renaksiId, namaRenaksi                                            sql.NullString
			urutan                                                            sql.NullInt64
			anggaran                                                          int64
		)

		err := rows.Scan(
			&pegawaiId,
			&namaPegawai,
			&kodeOpd,
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
			return nil, fmt.Errorf("error scanning rincian belanja pegawai: %v", err)
		}

		// Logika pengelompokan subkegiatan
		if currentSubkegiatan == nil ||
			currentSubkegiatan.KodeSubkegiatan != kodeSubkegiatan ||
			currentSubkegiatan.KodeOpd != kodeOpd {
			if currentSubkegiatan != nil {
				result = append(result, *currentSubkegiatan)
			}
			currentSubkegiatan = &domain.RincianBelanjaAsn{
				KodeOpd:         kodeOpd,
				KodeSubkegiatan: kodeSubkegiatan,
				NamaSubkegiatan: namaSubkegiatan,
				TotalAnggaran:   0,
				RencanaKinerja:  []domain.RencanaKinerjaAsn{},
			}
			currentRencanaKinerja = nil
			renaksiMap = make(map[string]bool)
		}

		// Logika rencana kinerja
		if currentRencanaKinerja == nil ||
			currentRencanaKinerja.RencanaKinerjaId != rekinId {
			currentRencanaKinerja = &domain.RencanaKinerjaAsn{
				RencanaKinerjaId: rekinId,
				RencanaKinerja:   namaRencanaKinerja,
				PegawaiId:        pegawaiId,
				NamaPegawai:      namaPegawai,
				RencanaAksi:      make([]domain.RincianBelanja, 0),
			}
			currentSubkegiatan.RencanaKinerja = append(currentSubkegiatan.RencanaKinerja, *currentRencanaKinerja)
			renaksiMap = make(map[string]bool) // RESET untuk rencana kinerja baru
		}

		// Tambahkan rencana aksi jika ada dan belum duplikat
		if renaksiId.Valid && namaRenaksi.Valid {
			if !renaksiMap[renaksiId.String] {
				rincianBelanja := domain.RincianBelanja{
					RenaksiId: renaksiId.String,
					Renaksi:   namaRenaksi.String,
					Anggaran:  anggaran,
				}
				lastIdx := len(currentSubkegiatan.RencanaKinerja) - 1
				currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi = append(
					currentSubkegiatan.RencanaKinerja[lastIdx].RencanaAksi,
					rincianBelanja,
				)
				currentSubkegiatan.TotalAnggaran += int(anggaran)
				renaksiMap[renaksiId.String] = true
			}
		}
	}

	if currentSubkegiatan != nil {
		result = append(result, *currentSubkegiatan)
	}

	return result, nil
}
