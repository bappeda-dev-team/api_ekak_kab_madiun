// repository/rencanaaksiopd_repository_impl.go
package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
)

type RencanaAksiOpdRepositoryImpl struct {
}

func NewRencanaAksiOpdRepositoryImpl() *RencanaAksiOpdRepositoryImpl {
	return &RencanaAksiOpdRepositoryImpl{}
}

func (repository *RencanaAksiOpdRepositoryImpl) FindBySasaranOpdAndTahun(ctx context.Context, tx *sql.Tx, sasaranOpdId int, tahun string) ([]domain.RencanaAksiOpd, error) {
	//INI SEBELUM DIEDIT
	query := `
	WITH renaksi_opd_data AS (
		SELECT 
			ro.id,
			ro.rekin_id,
			ro.sasaran_id,
			ro.tahun,
			ro.tw1,
			ro.tw2,
			ro.tw3,
			ro.tw4,
			ro.keterangan,
			so.nama_sasaran_opd as nama_sasaran_opd
		FROM tb_renaksi_opd ro
		LEFT JOIN tb_sasaran_opd so ON so.id = ro.sasaran_id
		WHERE ro.sasaran_id = ? 
		AND ro.tahun = ?
	),
	rencana_kinerja_data AS (
		SELECT 
			rk.id as rekin_id,
			rk.nama_rencana_kinerja,
			rk.pegawai_id as nip,
			rk.kode_opd,
			p.nama as nama_pegawai,
			rod.tw1,
			rod.tw2,
			rod.tw3,
			rod.tw4,
			rod.tahun,
			rod.keterangan,
			rod.nama_sasaran_opd
		FROM tb_rencana_kinerja rk
		JOIN renaksi_opd_data rod ON rod.rekin_id = rk.id
		LEFT JOIN tb_pegawai p ON p.nip = rk.pegawai_id
	),
	subkegiatan_data AS (
		SELECT 
			st.rekin_id,
			st.kode_subkegiatan,
			sk.nama_subkegiatan
		FROM tb_subkegiatan_terpilih st
		JOIN rencana_kinerja_data rkd ON rkd.rekin_id = st.rekin_id
		LEFT JOIN tb_subkegiatan sk ON sk.kode_subkegiatan = st.kode_subkegiatan
	),
	indikator_data AS (
		SELECT 
			i.id,
			i.kode,
			i.indikator,
			i.pagu_anggaran,
			t.target,
			t.satuan,
			rkd.rekin_id,
			sd.kode_subkegiatan
		FROM tb_indikator i
		JOIN subkegiatan_data sd ON sd.kode_subkegiatan = i.kode
		JOIN rencana_kinerja_data rkd ON rkd.rekin_id = sd.rekin_id
		LEFT JOIN tb_target t ON t.indikator_id = i.id AND t.tahun = rkd.tahun
		WHERE i.kode_opd = rkd.kode_opd
		AND i.tahun = rkd.tahun
	),
	renaksi_anggaran AS (
		SELECT 
			ra.rencana_kinerja_id as rekin_id,
			COALESCE(SUM(rb.anggaran), 0) as total_anggaran
		FROM tb_rencana_aksi ra
		LEFT JOIN tb_rincian_belanja rb ON rb.renaksi_id = ra.id
		WHERE ra.rencana_kinerja_id IN (SELECT rekin_id FROM rencana_kinerja_data)
		GROUP BY ra.rencana_kinerja_id
	)
	SELECT 
		rod.id,
		rod.sasaran_id,
		rod.nama_sasaran_opd,
		rkd.rekin_id,
		rkd.nama_rencana_kinerja,
		rkd.nip,
		rkd.nama_pegawai,
		rkd.kode_opd,
		sd.kode_subkegiatan,
		sd.nama_subkegiatan,
		id.id as indikator_id,
		id.indikator,
		id.target,
		id.satuan,
		id.pagu_anggaran,
		COALESCE(ra.total_anggaran, 0) as total_anggaran,
		rkd.tahun,
		rkd.tw1,
		rkd.tw2,
		rkd.tw3,
		rkd.tw4,
		rod.keterangan
	FROM renaksi_opd_data rod
	JOIN rencana_kinerja_data rkd ON rkd.rekin_id = rod.rekin_id
	LEFT JOIN subkegiatan_data sd ON sd.rekin_id = rkd.rekin_id
	LEFT JOIN indikator_data id ON id.rekin_id = rkd.rekin_id 
		AND id.kode_subkegiatan = sd.kode_subkegiatan
	LEFT JOIN renaksi_anggaran ra ON ra.rekin_id = rkd.rekin_id
	ORDER BY 
		rkd.rekin_id,
		sd.kode_subkegiatan,
		id.id
`

	// Tambahkan namaSasaranOpd dalam scanning
	rows, err := tx.QueryContext(ctx, query, sasaranOpdId, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying rencana aksi opd: %v", err)
	}
	defer rows.Close()

	var result []domain.RencanaAksiOpd
	var currentRencanaAksi *domain.RencanaAksiOpd
	var currentRencanaKinerja *domain.RencanaKinerjaOpd
	var currentSubKegiatan *domain.SubKegiatanOpdRenaksi

	for rows.Next() {
		var (
			id, sasaranId                    int
			namaSasaranOpd                   string
			rekinId, namaRencanaKinerja      string
			nip, namaPegawai, kodeOpd        string
			kodeSubKegiatan, namaSubKegiatan sql.NullString
			indikatorId, indikator           sql.NullString
			target, satuan                   sql.NullString
			paguAnggaran                     sql.NullInt64
			totalAnggaran                    int64
			tahun                            string
			tw1, tw2, tw3, tw4               sql.NullInt32
			keterangan                       sql.NullString
		)

		err := rows.Scan(
			&id,
			&sasaranId,
			&namaSasaranOpd,
			&rekinId,
			&namaRencanaKinerja,
			&nip,
			&namaPegawai,
			&kodeOpd,
			&kodeSubKegiatan,
			&namaSubKegiatan,
			&indikatorId,
			&indikator,
			&target,
			&satuan,
			&paguAnggaran,
			&totalAnggaran,
			&tahun,
			&tw1,
			&tw2,
			&tw3,
			&tw4,
			&keterangan,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rencana aksi opd: %v", err)
		}

		// Inisialisasi RencanaAksiOpd baru
		if currentRencanaAksi == nil || currentRencanaAksi.Id != id {
			var ket *string
			if keterangan.Valid {
				ket = &keterangan.String
			}
			currentRencanaAksi = &domain.RencanaAksiOpd{
				Id:             id,
				RekinId:        rekinId,
				SasaranOpdId:   sasaranId,
				NamaSasaranOpd: namaSasaranOpd,
				TahunRenaksi:   tahun,
				Tw1:            int(tw1.Int32),
				Tw2:            int(tw2.Int32),
				Tw3:            int(tw3.Int32),
				Tw4:            int(tw4.Int32),
				Keterangan:     ket,
				RencanaKinerja: []domain.RencanaKinerjaOpd{},
			}
			result = append(result, *currentRencanaAksi)
			currentRencanaKinerja = nil
		}

		// Inisialisasi RencanaKinerja baru
		if currentRencanaKinerja == nil || currentRencanaKinerja.RekinId != rekinId {
			var ket *string
			if keterangan.Valid {
				ket = &keterangan.String
			}
			currentRencanaKinerja = &domain.RencanaKinerjaOpd{
				Id:                 id,
				RekinId:            rekinId,
				NamaRencanaKinerja: namaRencanaKinerja,
				NipPegawai:         nip,
				NamaPegawai:        namaPegawai,
				KodeOpd:            kodeOpd,
				TotalAnggaran:      totalAnggaran,
				Tw1:                int(tw1.Int32),
				Tw2:                int(tw2.Int32),
				Tw3:                int(tw3.Int32),
				Tw4:                int(tw4.Int32),
				Keterangan:         ket,
				SubKegiatan:        []domain.SubKegiatanOpdRenaksi{},
			}
			lastIdx := len(result) - 1
			result[lastIdx].RencanaKinerja = append(result[lastIdx].RencanaKinerja, *currentRencanaKinerja)
			currentSubKegiatan = nil
		}

		// Inisialisasi SubKegiatan baru
		if kodeSubKegiatan.Valid && (currentSubKegiatan == nil || currentSubKegiatan.KodeSubKegiatan != kodeSubKegiatan.String) {
			currentSubKegiatan = &domain.SubKegiatanOpdRenaksi{
				KodeSubKegiatan: kodeSubKegiatan.String,
				NamaSubKegiatan: namaSubKegiatan.String,
				Indikator:       []domain.IndikatorSubKegiatanOpd{},
			}
			lastRKIdx := len(result[len(result)-1].RencanaKinerja) - 1
			result[len(result)-1].RencanaKinerja[lastRKIdx].SubKegiatan = append(
				result[len(result)-1].RencanaKinerja[lastRKIdx].SubKegiatan,
				*currentSubKegiatan,
			)
		}

		// Tambahkan Indikator jika ada
		if indikatorId.Valid && currentSubKegiatan != nil {
			indikatorData := domain.IndikatorSubKegiatanOpd{
				Id:           indikatorId.String,
				Indikator:    indikator.String,
				Target:       target.String,
				Satuan:       satuan.String,
				PaguAnggaran: paguAnggaran.Int64,
			}
			lastRKIdx := len(result[len(result)-1].RencanaKinerja) - 1
			lastSKIdx := len(result[len(result)-1].RencanaKinerja[lastRKIdx].SubKegiatan) - 1
			result[len(result)-1].RencanaKinerja[lastRKIdx].SubKegiatan[lastSKIdx].Indikator = append(
				result[len(result)-1].RencanaKinerja[lastRKIdx].SubKegiatan[lastSKIdx].Indikator,
				indikatorData,
			)
		}
	}

	return result, nil
}

func (repository *RencanaAksiOpdRepositoryImpl) SyncJadwalPelaksanaan(ctx context.Context, tx *sql.Tx, rekinId string) error {
	// Query untuk mendapatkan total bobot per triwindu
	query := `
        WITH pelaksanaan_data AS (
            -- Ambil semua pelaksanaan rencana aksi untuk rencana kinerja tertentu
            SELECT 
                pra.bulan,
                pra.bobot
            FROM tb_rencana_aksi ra
            JOIN tb_pelaksanaan_rencana_aksi pra ON ra.id = pra.rencana_aksi_id
            WHERE ra.rencana_kinerja_id = ?
        ),
        triwindu_totals AS (
            -- Hitung total bobot per triwindu
            SELECT
                CASE 
                    WHEN bulan BETWEEN 1 AND 3 THEN 1
                    WHEN bulan BETWEEN 4 AND 6 THEN 2
                    WHEN bulan BETWEEN 7 AND 9 THEN 3
                    WHEN bulan BETWEEN 10 AND 12 THEN 4
                END as triwindu,
                SUM(bobot) as total_bobot
            FROM pelaksanaan_data
            GROUP BY 
                CASE 
                    WHEN bulan BETWEEN 1 AND 3 THEN 1
                    WHEN bulan BETWEEN 4 AND 6 THEN 2
                    WHEN bulan BETWEEN 7 AND 9 THEN 3
                    WHEN bulan BETWEEN 10 AND 12 THEN 4
                END
        )
        SELECT 
            COALESCE(SUM(CASE WHEN triwindu = 1 THEN total_bobot ELSE 0 END), 0) as tw1,
            COALESCE(SUM(CASE WHEN triwindu = 2 THEN total_bobot ELSE 0 END), 0) as tw2,
            COALESCE(SUM(CASE WHEN triwindu = 3 THEN total_bobot ELSE 0 END), 0) as tw3,
            COALESCE(SUM(CASE WHEN triwindu = 4 THEN total_bobot ELSE 0 END), 0) as tw4
        FROM triwindu_totals
    `

	var tw1, tw2, tw3, tw4 int
	err := tx.QueryRowContext(ctx, query, rekinId).Scan(&tw1, &tw2, &tw3, &tw4)
	if err != nil {
		return fmt.Errorf("gagal mengambil data triwindu: %v", err)
	}

	// Update tb_renaksi_opd dengan nilai triwindu yang baru
	updateQuery := `
        UPDATE tb_renaksi_opd 
        SET tw1 = ?, tw2 = ?, tw3 = ?, tw4 = ?
        WHERE rekin_id = ?
    `

	_, err = tx.ExecContext(ctx, updateQuery, tw1, tw2, tw3, tw4, rekinId)
	if err != nil {
		return fmt.Errorf("gagal mengupdate triwindu di renaksi opd: %v", err)
	}

	return nil
}

func (repository *RencanaAksiOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, rencanaAksiOpd domain.RencanaAksiOpd) (domain.RencanaAksiOpd, error) {
	script := `
	INSERT INTO tb_renaksi_opd (id, rekin_id, sasaran_id, tahun, keterangan)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, script, rencanaAksiOpd.Id, rencanaAksiOpd.RekinId, rencanaAksiOpd.SasaranOpdId, rencanaAksiOpd.TahunRenaksi, rencanaAksiOpd.Keterangan)
	if err != nil {
		return domain.RencanaAksiOpd{}, err
	}

	return rencanaAksiOpd, nil
}

func (repository *RencanaAksiOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, rencanaAksiOpd domain.RencanaAksiOpd) domain.RencanaAksiOpd {
	script := `
	UPDATE tb_renaksi_opd 
	SET rekin_id = ?, keterangan = ?
	WHERE id = ?


	`
	_, err := tx.ExecContext(ctx, script, rencanaAksiOpd.RekinId, rencanaAksiOpd.Keterangan, rencanaAksiOpd.Id)
	if err != nil {
		return domain.RencanaAksiOpd{}
	}

	return rencanaAksiOpd
}

func (repository *RencanaAksiOpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, Id int) error {
	script := `	DELETE FROM tb_renaksi_opd WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, Id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *RencanaAksiOpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, Id int) (domain.RencanaAksiOpd, error) {
	query := `
    WITH target_data AS (
        SELECT 
            id,
            indikator_id,
            tahun,
            target,
            satuan
        FROM tb_target 
        WHERE tahun = (SELECT tahun FROM tb_renaksi_opd WHERE id = ?)
    )
    SELECT 
        ro.id,
        ro.rekin_id,
        ro.sasaran_id,
        ro.tahun,
        ro.keterangan,
        rk.nama_rencana_kinerja,
        so.nama_sasaran_opd,
        so.tahun_awal,
        so.tahun_akhir,
        so.jenis_periode,
        i.id as indikator_id,
        i.indikator,
        i.rumus_perhitungan,
        i.sumber_data,
        t.id as target_id,
        t.tahun as target_tahun,
        t.target,
        t.satuan
    FROM tb_renaksi_opd ro
    LEFT JOIN tb_rencana_kinerja rk ON ro.rekin_id = rk.id
    LEFT JOIN tb_sasaran_opd so ON ro.sasaran_id = so.id
    LEFT JOIN tb_indikator i ON so.id = i.sasaran_opd_id
    LEFT JOIN target_data t ON i.id = t.indikator_id
    WHERE ro.id = ?
    ORDER BY i.id ASC`

	rows, err := tx.QueryContext(ctx, query, Id, Id)
	if err != nil {
		return domain.RencanaAksiOpd{}, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var rencanaAksiOpd domain.RencanaAksiOpd
	var indikatorMap = make(map[string]*domain.IndikatorSasaranOpd)
	firstRow := true

	for rows.Next() {
		var (
			id, sasaranId                int
			rekinId                      string
			tahun                        string
			keterangan                   sql.NullString
			namaRencanaKinerja           string
			namaSasaranOpd               string
			tahunAwal, tahunAkhir        string
			jenisPeriode                 string
			indikatorId, indikatorNama   sql.NullString
			rumusPerhitungan, sumberData sql.NullString
			targetId, targetTahun        sql.NullString
			targetValue, targetSatuan    sql.NullString
		)

		err := rows.Scan(
			&id,
			&rekinId,
			&sasaranId,
			&tahun,
			&keterangan,
			&namaRencanaKinerja,
			&namaSasaranOpd,
			&tahunAwal,
			&tahunAkhir,
			&jenisPeriode,
			&indikatorId,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetTahun,
			&targetValue,
			&targetSatuan,
		)
		if err != nil {
			return domain.RencanaAksiOpd{}, fmt.Errorf("error scanning row: %v", err)
		}

		if firstRow {
			var keteranganString *string
			if keterangan.Valid {
				keteranganString = &keterangan.String
			}
			rencanaAksiOpd = domain.RencanaAksiOpd{
				Id:                 id,
				RekinId:            rekinId,
				SasaranOpdId:       sasaranId,
				TahunRenaksi:       tahun,
				Keterangan:         keteranganString,
				NamaRencanaKinerja: namaRencanaKinerja,
				SasaranOpd: domain.SasaranOpdDetailRenaksi{
					Id:             sasaranId,
					NamaSasaranOpd: namaSasaranOpd,
					TahunAwal:      tahunAwal,
					TahunAkhir:     tahunAkhir,
					JenisPeriode:   jenisPeriode,
					Indikator:      []domain.IndikatorSasaranOpd{},
				},
			}
			firstRow = false
		}

		// Process Indikator
		if indikatorId.Valid {
			indikator, exists := indikatorMap[indikatorId.String]
			if !exists {
				// Buat indikator baru
				indikator = &domain.IndikatorSasaranOpd{
					Id:               indikatorId.String,
					SasaranOpdId:     sasaranId,
					Indikator:        indikatorNama.String,
					RumusPerhitungan: rumusPerhitungan.String,
					SumberData:       sumberData.String,
					Target:           domain.TargetRenaksiOpd{},
				}

				// Set target default (kosong) untuk tahun yang sesuai
				indikator.Target = domain.TargetRenaksiOpd{
					Id:          "",
					IndikatorId: indikatorId.String,
					Tahun:       tahun,
					Target:      "",
					Satuan:      "",
				}

				// Update target jika ada
				if targetId.Valid && targetTahun.String == tahun {
					indikator.Target = domain.TargetRenaksiOpd{
						Id:          targetId.String,
						IndikatorId: indikatorId.String,
						Tahun:       targetTahun.String,
						Target:      targetValue.String,
						Satuan:      targetSatuan.String,
					}
				}

				indikatorMap[indikatorId.String] = indikator
				rencanaAksiOpd.SasaranOpd.Indikator = append(rencanaAksiOpd.SasaranOpd.Indikator, *indikator)
			}
		}
	}

	return rencanaAksiOpd, nil
}

// repository/rencanaaksiopd_repository_impl.go

func (repository *RencanaAksiOpdRepositoryImpl) FindAllSasaranByTahun(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.SasaranOpdDetailRenaksi, error) {
	query := `
    WITH target_data AS (
        SELECT 
            id,
            indikator_id,
            tahun,
            target,
            satuan
        FROM tb_target 
        WHERE tahun = ?
    ),
    valid_periode AS (
        -- Ambil periode yang valid dari tb_periode
        SELECT 
            tahun_awal,
            tahun_akhir,
            jenis_periode
        FROM tb_periode

    )
    SELECT DISTINCT
        so.id as sasaran_id,
        so.nama_sasaran_opd,
        so.tahun_awal,
        so.tahun_akhir,
        so.jenis_periode,
        i.id as indikator_id,
        i.indikator,
        i.rumus_perhitungan,
        i.sumber_data,
        t.id as target_id,
        t.tahun as target_tahun,
        t.target,
        t.satuan
    FROM tb_sasaran_opd so
    JOIN tb_pohon_kinerja pk ON so.pokin_id = pk.id
    LEFT JOIN tb_sasaran_opd_view sov ON sov.id_pokin = pk.id
    -- Join dengan valid_periode untuk memastikan periode yang valid
    INNER JOIN valid_periode vp ON (
        so.tahun_awal = vp.tahun_awal 
        AND so.tahun_akhir = vp.tahun_akhir 
        AND so.jenis_periode = vp.jenis_periode
    )
    LEFT JOIN tb_indikator i ON so.id = i.sasaran_opd_id
    LEFT JOIN target_data t ON i.id = t.indikator_id
    WHERE pk.kode_opd = ?
    AND ? BETWEEN so.tahun_awal AND so.tahun_akhir
    AND COALESCE(sov.is_hide, 0) = 0
    ORDER BY so.nama_sasaran_opd ASC, i.indikator ASC`

	rows, err := tx.QueryContext(ctx, query, tahun, kodeOpd, tahun)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	sasaranMap := make(map[int]*domain.SasaranOpdDetailRenaksi)
	indikatorMap := make(map[string]*domain.IndikatorSasaranOpd)

	for rows.Next() {
		var (
			sasaranId                                        int
			namaSasaran, tahunAwal, tahunAkhir, jenisPeriode string
			indikatorId, indikatorNama                       sql.NullString
			rumusPerhitungan, sumberData                     sql.NullString
			targetId, targetTahun                            sql.NullString
			targetValue, targetSatuan                        sql.NullString
		)

		err := rows.Scan(
			&sasaranId,
			&namaSasaran,
			&tahunAwal,
			&tahunAkhir,
			&jenisPeriode,
			&indikatorId,
			&indikatorNama,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetTahun,
			&targetValue,
			&targetSatuan,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Process Sasaran OPD
		sasaran, exists := sasaranMap[sasaranId]
		if !exists {
			sasaran = &domain.SasaranOpdDetailRenaksi{
				Id:             sasaranId,
				NamaSasaranOpd: namaSasaran,
				TahunAwal:      tahunAwal,
				TahunAkhir:     tahunAkhir,
				JenisPeriode:   jenisPeriode,
				Indikator:      []domain.IndikatorSasaranOpd{},
			}
			sasaranMap[sasaranId] = sasaran
		}

		// Process Indikator
		if indikatorId.Valid {
			indikator, exists := indikatorMap[indikatorId.String]
			if !exists {
				// Buat indikator baru
				indikator = &domain.IndikatorSasaranOpd{
					Id:               indikatorId.String,
					SasaranOpdId:     sasaranId,
					Indikator:        indikatorNama.String,
					RumusPerhitungan: rumusPerhitungan.String,
					SumberData:       sumberData.String,
					Target:           domain.TargetRenaksiOpd{},
				}

				// Set target default (kosong) untuk tahun yang sesuai
				indikator.Target = domain.TargetRenaksiOpd{
					Id:          "",
					IndikatorId: indikatorId.String,
					Tahun:       tahun,
					Target:      "",
					Satuan:      "",
				}

				// Update target jika ada
				if targetId.Valid && targetTahun.String == tahun {
					indikator.Target = domain.TargetRenaksiOpd{
						Id:          targetId.String,
						IndikatorId: indikatorId.String,
						Tahun:       targetTahun.String,
						Target:      targetValue.String,
						Satuan:      targetSatuan.String,
					}
				}

				indikatorMap[indikatorId.String] = indikator
				sasaran.Indikator = append(sasaran.Indikator, *indikator)
			}
		}
	}

	// Convert map to slice dan sort berdasarkan nama sasaran
	var result []domain.SasaranOpdDetailRenaksi
	for _, sasaran := range sasaranMap {
		// Sort indikator berdasarkan nama indikator
		sort.Slice(sasaran.Indikator, func(i, j int) bool {
			return sasaran.Indikator[i].Indikator < sasaran.Indikator[j].Indikator
		})
		result = append(result, *sasaran)
	}

	// Sort sasaran berdasarkan nama
	sort.Slice(result, func(i, j int) bool {
		return result[i].NamaSasaranOpd < result[j].NamaSasaranOpd
	})

	return result, nil
}
