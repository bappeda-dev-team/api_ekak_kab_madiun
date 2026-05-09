package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkdRepositoryImpl struct {
}

func NewIkdRepositoryImpl() *IkdRepositoryImpl {
	return &IkdRepositoryImpl{}
}

func (repository *IkdRepositoryImpl) FindAll(
	ctx context.Context,
	tx *sql.Tx,
	kodeOpd string,
	tahun string,
	jenisPeriode string,
) ([]domain.IkdDetail, error) {

	script := `
	SELECT
		pk.id,
		COALESCE(pk.nama_pohon, '') as nama_pohon,
		COALESCE(pk.parent, 0) as parent,
		COALESCE(pk.jenis_pohon, '') as jenis_pohon,
		COALESCE(pk.level_pohon, 0) as level_pohon,
		COALESCE(pk.kode_opd, '') as kode_opd,
		COALESCE(pk.keterangan, '') as keterangan,
		COALESCE(pk.keterangan_crosscutting, '') as keterangan_crosscutting,
		COALESCE(pk.tahun, '') as tahun,
		COALESCE(pk.status, '') as status,
		COALESCE(pk.is_active, 0) as is_active,

		-- PELAKSANA
		COALESCE(pp.id, '') as pelaksana_id,
		COALESCE(pp.pegawai_id, '') as pegawai_id,
		COALESCE(pg.nip, '') as nip,
		COALESCE(pg.nama, '') as nama_pegawai,

		-- SASARAN OPD
		COALESCE(so.id, 0) as sasaran_id,
		COALESCE(so.nama_sasaran_opd, '') as nama_sasaran_opd,
		COALESCE(so.id_tujuan_opd, 0) as id_tujuan_opd,
		COALESCE(so.tahun_awal, '') as tahun_awal,
		COALESCE(so.tahun_akhir, '') as tahun_akhir,
		COALESCE(so.jenis_periode, '') as jenis_periode,

		-- TUJUAN OPD
		COALESCE(to2.tujuan, '') as tujuan_opd,

		-- INDIKATOR
		COALESCE(i.id, '') as indikator_id,
		COALESCE(i.indikator, '') as indikator,
		COALESCE(i.rumus_perhitungan, '') as rumus_perhitungan,
		COALESCE(i.sumber_data, '') as sumber_data,

		-- TARGET
		COALESCE(t.id, '') as target_id,
		COALESCE(t.tahun, '') as target_tahun,
		COALESCE(t.target, '') as target_value,
		COALESCE(t.satuan, '') as target_satuan,

		-- PROGRAM OPD (TACTICAL)
		COALESCE(tp.id, 0) as program_id,
		COALESCE(tp.parent, 0) as program_parent,
		COALESCE(tp.nama_pohon, '') as nama_program

	FROM tb_pohon_kinerja pk

	LEFT JOIN tb_pelaksana_pokin pp
		ON pk.id = pp.pohon_kinerja_id

	LEFT JOIN tb_pegawai pg
		ON pp.pegawai_id = pg.id

	LEFT JOIN tb_sasaran_opd so
		ON pk.id = so.pokin_id
		AND pk.level_pohon = 4
		AND CAST(? AS SIGNED)
			BETWEEN CAST(so.tahun_awal AS SIGNED)
			AND CAST(so.tahun_akhir AS SIGNED)
		AND so.jenis_periode = ?

	LEFT JOIN tb_tujuan_opd to2
		ON so.id_tujuan_opd = to2.id

	LEFT JOIN tb_indikator i
		ON so.id = i.sasaran_opd_id

	LEFT JOIN tb_target t
		ON i.id = t.indikator_id
		AND t.tahun = ?

	-- PROGRAM OPD DARI TACTICAL
	LEFT JOIN tb_pohon_kinerja tp
		ON tp.parent = pk.id
		AND LOWER(tp.jenis_pohon) = 'tactical'

	WHERE pk.kode_opd = ?
	AND pk.tahun = ?
	AND LOWER(pk.jenis_pohon) = 'strategic'

	AND pk.status NOT IN (
		'menunggu_disetujui',
		'tarik pokin opd',
		'disetujui',
		'ditolak',
		'crosscutting_menunggu',
		'crosscutting_ditolak'
	)

	ORDER BY
		pk.level_pohon,
		pk.id ASC
	`

	rows, err := tx.QueryContext(
		ctx,
		script,
		tahun,
		jenisPeriode,
		tahun,
		kodeOpd,
		tahun,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pohonMap := make(map[int]*domain.IkdDetail)

	for rows.Next() {

		var (
			pohon domain.IkdDetail

			pelaksanaId string
			pegawaiId   string
			nip         string
			namaPegawai string

			sasaranId         int
			namaSasaranOpd    string
			idTujuanOpd       int
			tahunAwal         string
			tahunAkhir        string
			jenisPeriodeValue string
			tujuanOpd         string

			indikatorId      string
			indikator        string
			rumusPerhitungan string
			sumberData       string

			targetId     string
			targetTahun  string
			targetValue  string
			targetSatuan string

			programId     int
			programParent int
			namaProgram   string
		)

		err := rows.Scan(
			&pohon.Id,
			&pohon.NamaPohon,
			&pohon.Parent,
			&pohon.JenisPohon,
			&pohon.LevelPohon,
			&pohon.KodeOpd,
			&pohon.Keterangan,
			&pohon.KeteranganCrosscutting,
			&pohon.Tahun,
			&pohon.Status,
			&pohon.IsActive,

			&pelaksanaId,
			&pegawaiId,
			&nip,
			&namaPegawai,

			&sasaranId,
			&namaSasaranOpd,
			&idTujuanOpd,
			&tahunAwal,
			&tahunAkhir,
			&jenisPeriodeValue,

			&tujuanOpd,

			&indikatorId,
			&indikator,
			&rumusPerhitungan,
			&sumberData,

			&targetId,
			&targetTahun,
			&targetValue,
			&targetSatuan,

			&programId,
			&programParent,
			&namaProgram,
		)

		if err != nil {
			return nil, err
		}

		existingPohon, exists := pohonMap[pohon.Id]

		if !exists {

			pohon.Pelaksana = make([]domain.PelaksanaDetail, 0)
			pohon.SasaranOpd = make([]domain.SasaranOpdDetail, 0)
			pohon.ProgramOpd = make([]domain.ProgramOpdDetail, 0)

			pohonMap[pohon.Id] = &pohon
			existingPohon = &pohon
		}

		// PELAKSANA
		if pelaksanaId != "" {

			existsPelaksana := false

			for _, p := range existingPohon.Pelaksana {
				if p.Id == pelaksanaId {
					existsPelaksana = true
					break
				}
			}

			if !existsPelaksana {

				existingPohon.Pelaksana = append(
					existingPohon.Pelaksana,
					domain.PelaksanaDetail{
						Id:          pelaksanaId,
						PegawaiId:   pegawaiId,
						Nip:         nip,
						NamaPegawai: namaPegawai,
					},
				)
			}
		}

		// SASARAN
		var existingSasaran *domain.SasaranOpdDetail

		if sasaranId != 0 {

			sasaranExists := false

			for i := range existingPohon.SasaranOpd {

				if existingPohon.SasaranOpd[i].Id == sasaranId {

					sasaranExists = true
					existingSasaran = &existingPohon.SasaranOpd[i]
					break
				}
			}

			if !sasaranExists {

				newSasaran := domain.SasaranOpdDetail{
					Id:             sasaranId,
					IdPohon:        pohon.Id,
					NamaSasaranOpd: namaSasaranOpd,
					IdTujuanOpd:    idTujuanOpd,
					NamaTujuanOpd:  tujuanOpd,
					TahunAwal:      tahunAwal,
					TahunAkhir:     tahunAkhir,
					JenisPeriode:   jenisPeriodeValue,
					Indikator:      make([]domain.Indikator, 0),
				}

				existingPohon.SasaranOpd = append(
					existingPohon.SasaranOpd,
					newSasaran,
				)

				existingSasaran = &existingPohon.SasaranOpd[len(existingPohon.SasaranOpd)-1]
			}
		}

		// INDIKATOR
		var existingIndikator *domain.Indikator

		if existingSasaran != nil && indikatorId != "" {

			indikatorExists := false

			for i := range existingSasaran.Indikator {

				if existingSasaran.Indikator[i].Id == indikatorId {

					indikatorExists = true
					existingIndikator = &existingSasaran.Indikator[i]
					break
				}
			}

			if !indikatorExists {

				newIndikator := domain.Indikator{
					Id:        indikatorId,
					Indikator: indikator,
					RumusPerhitungan: sql.NullString{
						String: rumusPerhitungan,
						Valid:  rumusPerhitungan != "",
					},
					SumberData: sql.NullString{
						String: sumberData,
						Valid:  sumberData != "",
					},
					Target: make([]domain.Target, 0),
				}

				existingSasaran.Indikator = append(
					existingSasaran.Indikator,
					newIndikator,
				)

				existingIndikator = &existingSasaran.Indikator[len(existingSasaran.Indikator)-1]
			}
		}

		// TARGET
		if existingIndikator != nil && targetId != "" {

			targetExists := false

			for _, t := range existingIndikator.Target {
				if t.Id == targetId {
					targetExists = true
					break
				}
			}

			if !targetExists {

				existingIndikator.Target = append(
					existingIndikator.Target,
					domain.Target{
						Id:          targetId,
						IndikatorId: indikatorId,
						Tahun:       targetTahun,
						Target:      targetValue,
						Satuan:      targetSatuan,
					},
				)
			}
		}

		// PROGRAM OPD
		if programId != 0 {

			programExists := false

			for _, p := range existingPohon.ProgramOpd {
				if p.Id == programId {
					programExists = true
					break
				}
			}

			if !programExists {

				existingPohon.ProgramOpd = append(
					existingPohon.ProgramOpd,
					domain.ProgramOpdDetail{
						Id:          programId,
						Parent:      programParent,
						NamaProgram: namaProgram,
					},
				)
			}
		}
	}

	var result []domain.IkdDetail

	for _, value := range pohonMap {
		result = append(result, *value)
	}

	if result == nil {
		result = make([]domain.IkdDetail, 0)
	}

	return result, nil
}