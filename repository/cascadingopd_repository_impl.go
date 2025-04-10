package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type CascadingOpdRepositoryImpl struct {
}

func NewCascadingOpdRepositoryImpl(db *sql.DB, rencanaKinerjaRepository RencanaKinerjaRepository) *CascadingOpdRepositoryImpl {
	return &CascadingOpdRepositoryImpl{}
}

func (repository *CascadingOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            id,
            COALESCE(nama_pohon, '') as nama_pohon,
            COALESCE(parent, 0) as parent,
            COALESCE(jenis_pohon, '') as jenis_pohon,
            COALESCE(level_pohon, 0) as level_pohon,
            COALESCE(kode_opd, '') as kode_opd,
            COALESCE(keterangan, '') as keterangan,
            COALESCE(keterangan_crosscutting, '') as keterangan_crosscutting,
            COALESCE(tahun, '') as tahun,
            COALESCE(status, '') as status,
			COALESCE(is_active) as is_active
        FROM tb_pohon_kinerja 
        WHERE kode_opd = ? 
		AND tahun = ?
		AND status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak', 'crosscutting_menunggu', 'crosscutting_ditolak')
        ORDER BY level_pohon, id ASC`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.KeteranganCrosscutting,
			&pokin.Tahun,
			&pokin.Status,
			&pokin.IsActive,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}

	// Inisialisasi slice kosong jika tidak ada data
	if pokins == nil {
		pokins = make([]domain.PohonKinerja, 0)
	}

	return pokins, nil
}

func (repository *CascadingOpdRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error) {
	script := `
        SELECT i.id, i.pokin_id, i.indikator, 
               t.id, t.indikator_id, t.target, t.satuan
        FROM tb_indikator i
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE i.pokin_id = ?`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var indId, pokinId, indikator string
		var targetId, indikatorId, target, satuan sql.NullString

		err := rows.Scan(
			&indId, &pokinId, &indikator,
			&targetId, &indikatorId, &target, &satuan)
		if err != nil {
			return nil, err
		}

		// Proses Indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:        indId,
				Indikator: indikator,
				Target:    []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// Proses Target jika ada
		if targetId.Valid && indikatorId.Valid {
			target := domain.Target{
				Id:          targetId.String,
				IndikatorId: indikatorId.String,
				Target:      target.String,
				Satuan:      satuan.String,
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
