package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"strings"
)

type ReviewRepositoryImpl struct {
}

func NewReviewRepositoryImpl() *ReviewRepositoryImpl {
	return &ReviewRepositoryImpl{}
}

func (repository *ReviewRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, review domain.Review) (domain.Review, error) {
	script := "INSERT INTO tb_review (id, id_pohon_kinerja, review, keterangan, jenis_pokin, created_by) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, review.Id, review.IdPohonKinerja, review.Review, review.Keterangan, review.Jenis_pokin, review.CreatedBy)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (repository *ReviewRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, review domain.Review) (domain.Review, error) {
	script := "UPDATE tb_review SET review = ?, keterangan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, review.Review, review.Keterangan, review.Id)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (repository *ReviewRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_review WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ReviewRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Review, error) {
	script := "SELECT id, id_pohon_kinerja, review, keterangan, jenis_pokin, created_by, created_at, updated_at FROM tb_review WHERE id = ?"
	row := tx.QueryRowContext(ctx, script, id)
	var review domain.Review
	err := row.Scan(&review.Id, &review.IdPohonKinerja, &review.Review, &review.Keterangan, &review.Jenis_pokin, &review.CreatedBy, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (repository *ReviewRepositoryImpl) FindByPohonKinerja(ctx context.Context, tx *sql.Tx, idPohonKinerja int) ([]domain.Review, error) {
	script := "SELECT id, id_pohon_kinerja, review, keterangan, jenis_pokin, created_by, created_at, updated_at FROM tb_review WHERE id_pohon_kinerja = ?"
	rows, err := tx.QueryContext(ctx, script, idPohonKinerja)
	if err != nil {
		return []domain.Review{}, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var review domain.Review
		err := rows.Scan(&review.Id, &review.IdPohonKinerja, &review.Review, &review.Keterangan, &review.Jenis_pokin, &review.CreatedBy, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return []domain.Review{}, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (repository *ReviewRepositoryImpl) CountReviewByPohonKinerja(ctx context.Context, tx *sql.Tx, idPohonKinerja int) (int, error) {
	script := "SELECT COUNT(*) FROM tb_review WHERE id_pohon_kinerja = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, idPohonKinerja).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repository *ReviewRepositoryImpl) FindAllReviewByTematik(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.ReviewTematik, error) {
	query := `
        WITH RECURSIVE pohon_hierarchy AS (
            -- Base case: ambil semua tematik (level 0)
            SELECT 
                id, nama_pohon, parent, level_pohon, jenis_pohon, created_at, updated_at
            FROM tb_pohon_kinerja
            WHERE level_pohon = 0
            AND tahun = ?

            UNION ALL

            -- Recursive case: ambil semua turunan
            SELECT 
                c.id, c.nama_pohon, c.parent, c.level_pohon, c.jenis_pohon, c.created_at, c.updated_at
            FROM tb_pohon_kinerja c
            INNER JOIN pohon_hierarchy p ON c.parent = p.id
            WHERE c.tahun = ?
        )
        SELECT 
            t.id as id_tematik,
            t.nama_pohon as nama_tematik,
            t.level_pohon as level_tematik,
            ph.id as pohon_id,
            ph.parent,
            ph.nama_pohon,
            ph.level_pohon,
			ph.jenis_pohon,
			ph.created_at, 
			ph.updated_at,
            r.review,
            r.keterangan,
            r.created_by,
            r.jenis_pokin
        FROM tb_pohon_kinerja t
        -- Mulai dari tematik level 0
        LEFT JOIN pohon_hierarchy ph ON 
            ph.id = t.id OR 
            EXISTS (
                WITH RECURSIVE tree AS (
                    SELECT id, parent FROM tb_pohon_kinerja WHERE id = ph.id AND tahun = ?
                    UNION ALL
                    SELECT p.id, p.parent FROM tb_pohon_kinerja p
                    INNER JOIN tree tr ON p.id = tr.parent
                    WHERE p.tahun = ?
                )
                SELECT 1 FROM tree WHERE parent = t.id
            )
        LEFT JOIN tb_review r ON r.id_pohon_kinerja = ph.id
        WHERE t.level_pohon = 0 
        AND t.tahun = ?
        ORDER BY t.id, COALESCE(ph.level_pohon, -1), COALESCE(ph.id, 0)`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun, tahun, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ReviewTematik
	var currentTematik *domain.ReviewTematik

	for rows.Next() {
		var (
			idTematik    int
			namaTematik  string
			levelTematik int
			pohonId      sql.NullInt64
			parent       sql.NullInt64
			namaPohon    sql.NullString
			levelPohon   sql.NullInt64
			jenispohon   sql.NullString
			review       sql.NullString
			keterangan   sql.NullString
			createdBy    sql.NullString
			jenisPokin   sql.NullString
			created_at   sql.NullString
			updated_at   sql.NullString
		)

		err := rows.Scan(
			&idTematik,
			&namaTematik,
			&levelTematik,
			&pohonId,
			&parent,
			&namaPohon,
			&levelPohon,
			&jenispohon,
			&created_at,
			&updated_at,
			&review,
			&keterangan,
			&createdBy,
			&jenisPokin,
		)
		if err != nil {
			return nil, err
		}

		// Jika tematik baru atau pertama kali
		if currentTematik == nil || currentTematik.IdTematik != idTematik {
			if currentTematik != nil {
				result = append(result, *currentTematik)
			}
			currentTematik = &domain.ReviewTematik{
				IdTematik:  idTematik,
				NamaPohon:  namaTematik,
				LevelPohon: levelTematik,
				Review:     []domain.ReviewDetail{},
			}
		}

		// Hanya tambahkan review detail jika ada data review
		if pohonId.Valid && review.Valid {
			reviewDetail := domain.ReviewDetail{
				IdPohon:    int(pohonId.Int64),
				Parent:     int(parent.Int64),
				NamaPohon:  namaPohon.String,
				LevelPohon: int(levelPohon.Int64),
				JenisPohon: jenispohon.String,
				Review:     review.String,
				Keterangan: keterangan.String,
				CreatedBy:  createdBy.String,
				JenisPokin: jenisPokin.String,
				CreatedAt:  created_at.String,
				UpdatedAt:  updated_at.String,
			}
			currentTematik.Review = append(currentTematik.Review, reviewDetail)
		}
	}

	// Tambahkan tematik terakhir jika ada
	if currentTematik != nil {
		result = append(result, *currentTematik)
	}

	return result, nil
}

func (repository *ReviewRepositoryImpl) FindAllReviewOpd(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.ReviewOpd, error) {
	query := `
        SELECT 
            pk.id as id_pohon,
            COALESCE(pk.parent, 0) as parent,
            COALESCE(pk.nama_pohon, '') as nama_pohon,
            COALESCE(pk.level_pohon, 0) as level_pohon,
            COALESCE(pk.jenis_pohon, '') as jenis_pohon,
            COALESCE(r.review, '') as review,
            COALESCE(r.keterangan, '') as keterangan,
            COALESCE(r.created_by, '') as created_by,
            COALESCE(r.created_at, CURRENT_TIMESTAMP) as created_at,
            COALESCE(r.updated_at, CURRENT_TIMESTAMP) as updated_at
        FROM tb_pohon_kinerja pk
        INNER JOIN tb_review r ON r.id_pohon_kinerja = pk.id  -- Ganti LEFT JOIN menjadi INNER JOIN
        WHERE pk.kode_opd = ?
        AND pk.tahun = ?
        AND pk.level_pohon >= 4
        AND pk.status NOT IN (
            'menunggu_disetujui', 
            'tarik pokin opd', 
            'disetujui', 
            'ditolak', 
            'crosscutting_menunggu', 
            'crosscutting_ditolak'
        )
        AND r.review IS NOT NULL  -- Tambahan filter untuk memastikan ada review
        AND r.review != ''        -- Tambahan filter untuk memastikan review tidak kosong
        ORDER BY pk.level_pohon, pk.id ASC`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.ReviewOpd
	for rows.Next() {
		var review domain.ReviewOpd
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&review.IdPohon,
			&review.Parent,
			&review.NamaPohon,
			&review.LevelPohon,
			&review.JenisPohon,
			&review.Review,
			&review.Keterangan,
			&review.CreatedBy,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			review.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			review.UpdatedAt = updatedAt.String
		}

		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if reviews == nil {
		reviews = make([]domain.ReviewOpd, 0)
	}

	return reviews, nil
}

func (repository *ReviewRepositoryImpl) FindByPokinIdBatch(ctx context.Context, tx *sql.Tx, pokinIds []int) ([]domain.ReviewWithNama, error) {
	if len(pokinIds) == 0 {
		return make([]domain.ReviewWithNama, 0), nil
	}

	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}

	script := fmt.Sprintf(`
		SELECT 
			r.id, 
			r.id_pohon_kinerja, 
			r.review, 
			r.keterangan, 
			r.created_by, 
			COALESCE(p.nama, '') as nama_reviewer,
			r.jenis_pokin
		FROM tb_review r
		LEFT JOIN tb_pegawai p ON p.nip = r.created_by
		WHERE r.id_pohon_kinerja IN (%s)
		ORDER BY r.id_pohon_kinerja, r.id
	`, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.ReviewWithNama
	for rows.Next() {
		var review domain.ReviewWithNama
		err := rows.Scan(
			&review.Id,
			&review.IdPohonKinerja,
			&review.Review,
			&review.Keterangan,
			&review.CreatedBy,
			&review.NamaReviewer,
			&review.Jenis_pokin,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if reviews == nil {
		reviews = make([]domain.ReviewWithNama, 0)
	}

	return reviews, nil
}
