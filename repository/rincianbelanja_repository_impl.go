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
            COALESCE(rb.anggaran, 0) as anggaran
        FROM rencana_kinerja_pegawai rkp
        LEFT JOIN tb_rencana_aksi ra ON ra.rencana_kinerja_id = rkp.rekin_id
        LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
        ORDER BY rkp.kode_subkegiatan, rkp.rekin_id, ra.id
    `

	rows, err := tx.QueryContext(ctx, query, pegawaiId, tahun)
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
    WITH valid_users AS (
            SELECT DISTINCT p.nip
            FROM tb_pegawai p
            JOIN tb_users u ON u.nip = p.nip
            JOIN tb_user_role ur ON ur.user_id = u.id
            JOIN tb_role r ON r.id = ur.role_id
            WHERE r.role = 'level_3'
        ),
        rencana_kinerja_opd AS (
            SELECT 
                rk.id as rekin_id,
                rk.pegawai_id,
                p.nama as nama_pegawai,
                st.kode_subkegiatan,
                sk.nama_subkegiatan,
                rk.nama_rencana_kinerja
            FROM tb_rencana_kinerja rk
            INNER JOIN valid_users vu ON vu.nip = rk.pegawai_id
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
