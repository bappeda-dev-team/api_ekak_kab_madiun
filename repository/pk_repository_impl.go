package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type PkRepositoryImpl struct{}

func NewPkRepositoryImpl() *PkRepositoryImpl {
	return &PkRepositoryImpl{}
}

func (repository *PkRepositoryImpl) FindByKodeOpdTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun int) (map[int][]domain.PkOpd, error) {
	query := `
    SELECT pk.id,
           pk.kode_opd,
           pk.nama_opd,
           pk.level_pk,
           pk.nama_atasan,
           pk.nip_atasan,
           pk.id_rekin_atasan,
           pk.rekin_atasan,
           pk.nip_pemilik_pk,
           pk.nama_pemilik_pk,
           pk.id_rekin_pemilik_pk,
           pk.rekin_pemilik_pk,
           pk.tahun,
           pk.keterangan
    FROM pk_opd pk
    WHERE pk.kode_opd = ? AND pk.tahun = ?
    ORDER BY pk.level_pk
    `

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// group by level
	results := make(map[int][]domain.PkOpd)

	for rows.Next() {
		var pkOpd domain.PkOpd

		err := rows.Scan(
			&pkOpd.Id,
			&pkOpd.KodeOpd,
			&pkOpd.NamaOpd,
			&pkOpd.LevelPk,
			&pkOpd.NamaAtasan,
			&pkOpd.NipAtasan,
			&pkOpd.IdRekinAtasan,
			&pkOpd.RekinAtasan,
			&pkOpd.NipPemilikPk,
			&pkOpd.NamaPemilikPk,
			&pkOpd.IdRekinPemilikPk,
			&pkOpd.RekinPemilikPk,
			&pkOpd.Tahun,
			&pkOpd.Keterangan,
		)
		if err != nil {
			return nil, fmt.Errorf("scan pk_opd failed: %w", err)
		}
		results[pkOpd.LevelPk] = append(results[pkOpd.LevelPk], pkOpd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repository *PkRepositoryImpl) HubungkanRekin(
	ctx context.Context,
	tx *sql.Tx,
	pk domain.PkOpd,
) error {

	query := `
	INSERT INTO pk_opd (
		id,
		kode_opd,
		nama_opd,
		level_pk,
		nip_atasan,
		nama_atasan,
		id_rekin_atasan,
		rekin_atasan,
		nip_pemilik_pk,
		nama_pemilik_pk,
		id_rekin_pemilik_pk,
		rekin_pemilik_pk,
		tahun,
		keterangan
	) VALUES (
		UUID(),
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
	ON DUPLICATE KEY UPDATE
		level_pk           = VALUES(level_pk),
		nip_atasan         = VALUES(nip_atasan),
		nama_atasan        = VALUES(nama_atasan),
		id_rekin_atasan    = VALUES(id_rekin_atasan),
		rekin_atasan       = VALUES(rekin_atasan),
		keterangan         = VALUES(keterangan),
		updated_at         = CURRENT_TIMESTAMP
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		pk.KodeOpd,
		pk.NamaOpd,
		pk.LevelPk,
		pk.NipAtasan,
		pk.NamaAtasan,
		pk.IdRekinAtasan,
		pk.RekinAtasan,
		pk.NipPemilikPk,
		pk.NamaPemilikPk,
		pk.IdRekinPemilikPk,
		pk.RekinPemilikPk,
		pk.Tahun,
		pk.Keterangan,
	)

	if err != nil {
		return fmt.Errorf("hubungkan rekin failed: %w", err)
	}

	return nil
}

func (repository *PkRepositoryImpl) FindSubkegiatanByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) (map[string]domain.AllItemPk, error) {
	if len(rekinIds) == 0 {
		return make(map[string]domain.AllItemPk), nil
	}

	placeholders := make([]string, len(rekinIds))
	args := make([]any, len(rekinIds))
	for i, id := range rekinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT st.rekin_id,
          p.kode_program,
          p.nama_program,
          k.kode_kegiatan,
          k.nama_kegiatan,
          sub.kode_subkegiatan,
          sub.nama_subkegiatan
		FROM tb_subkegiatan_terpilih st
        JOIN tb_subkegiatan sub ON sub.id = st.subkegiatan_id
        LEFT JOIN tb_master_kegiatan k ON k.kode_kegiatan = SUBSTRING_INDEX(st.kode_subkegiatan, '.', 5)
        LEFT JOIN tb_master_program p ON p.kode_program = SUBSTRING_INDEX(st.kode_subkegiatan, '.', 3)
		WHERE st.rekin_id IN (%s)`,
		strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subMap := make(map[string]domain.AllItemPk)
	for rows.Next() {
		var itemPk domain.AllItemPk
		var kodeProgram, namaProgram sql.NullString
		var kodeKegiatan, namaKegiatan sql.NullString

		err := rows.Scan(&itemPk.RekinId,
			&kodeProgram,
			&namaProgram,
			&kodeKegiatan,
			&namaKegiatan,
			&itemPk.KodeSubkegiatan,
			&itemPk.NamaSubkegiatan)
		if err != nil {
			return nil, err
		}
		if kodeProgram.Valid {
			itemPk.KodeProgram = kodeProgram.String
		}
		if namaProgram.Valid {
			itemPk.NamaProgram = namaProgram.String
		}
		if kodeKegiatan.Valid {
			itemPk.KodeKegiatan = kodeKegiatan.String
		}
		if namaKegiatan.Valid {
			itemPk.NamaKegiatan = namaKegiatan.String
		}
		subMap[itemPk.RekinId] = itemPk
	}

	return subMap, nil
}

func (repository *PkRepositoryImpl) FindTotalPaguAnggaranByRekinIds(ctx context.Context, tx *sql.Tx, rekinIds []string) (map[string]int, error) {
	if len(rekinIds) == 0 {
		return make(map[string]int), nil
	}

	placeholders := make([]string, len(rekinIds))
	args := make([]any, len(rekinIds))
	for i, id := range rekinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
         SELECT ra.rencana_kinerja_id, SUM(rb.anggaran)
         FROM tb_rincian_belanja rb
         JOIN tb_rencana_aksi ra ON ra.id = rb.renaksi_id
         WHERE ra.rencana_kinerja_id IN (%s)
         GROUP BY ra.rencana_kinerja_id
    `, strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	totalPaguMap := make(map[string]int)
	for rows.Next() {
		var rekinIdNs sql.NullString
		var rekinId string
		var totalAnggaranNs sql.NullInt64
		var totalAnggaran int
		err := rows.Scan(
			&rekinIdNs,
			&totalAnggaranNs,
		)
		if err != nil {
			return nil, err
		}
		if rekinIdNs.Valid {
			rekinId = rekinIdNs.String
		}
		if totalAnggaranNs.Valid {
			totalAnggaran = int(totalAnggaranNs.Int64)
		}
		totalPaguMap[rekinId] = totalAnggaran
	}
	return totalPaguMap, nil
}
