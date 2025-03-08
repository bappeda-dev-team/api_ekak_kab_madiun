package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
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
