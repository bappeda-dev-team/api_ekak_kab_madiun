package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type TaggingPokinRepositoryImpl struct {
}

func NewTaggingPokinRepositoryImpl() *TaggingPokinRepositoryImpl {
	return &TaggingPokinRepositoryImpl{}
}

func (repository *TaggingPokinRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tagging domain.TaggingPokin) (domain.TaggingPokin, error) {
	script := "INSERT INTO tb_master_tagging (nama_tagging, keterangan_tagging) VALUES (?,?)"
	result, err := tx.ExecContext(ctx, script, tagging.NamaTagging, tagging.KeteranganTagging)
	if err != nil {
		return domain.TaggingPokin{}, fmt.Errorf("error saat menyimpan tagging: %v", err)
	}

	// Mengambil ID yang baru dibuat
	id, err := result.LastInsertId()
	if err != nil {
		return domain.TaggingPokin{}, fmt.Errorf("error saat mengambil id: %v", err)
	}

	// Set ID ke struct tagging
	tagging.Id = int(id)
	return tagging, nil
}

func (repository *TaggingPokinRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tagging domain.TaggingPokin) (domain.TaggingPokin, error) {
	script := "UPDATE tb_master_tagging SET nama_tagging = ?, keterangan_tagging = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, tagging.NamaTagging, tagging.KeteranganTagging, tagging.Id)
	if err != nil {
		return domain.TaggingPokin{}, fmt.Errorf("error saat mengupdate tagging: %v", err)
	}
	return tagging, nil
}

func (repository *TaggingPokinRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_master_tagging WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return fmt.Errorf("error saat menghapus tagging: %v", err)
	}
	return nil
}

func (repository *TaggingPokinRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.TaggingPokin, error) {
	script := "SELECT id, nama_tagging, keterangan_tagging FROM tb_master_tagging"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domain.TaggingPokin{}, fmt.Errorf("error saat mengambil semua tagging: %v", err)
	}
	defer rows.Close()

	var taggings []domain.TaggingPokin
	for rows.Next() {
		var tagging domain.TaggingPokin
		err := rows.Scan(&tagging.Id, &tagging.NamaTagging, &tagging.KeteranganTagging)
		if err != nil {
			return []domain.TaggingPokin{}, fmt.Errorf("error saat mengambil semua tagging: %v", err)
		}
		taggings = append(taggings, tagging)
	}
	return taggings, nil
}

func (repository *TaggingPokinRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.TaggingPokin, error) {
	script := "SELECT id, nama_tagging, keterangan_tagging FROM tb_master_tagging WHERE id = ?"
	var tagging domain.TaggingPokin
	err := tx.QueryRowContext(ctx, script, id).Scan(&tagging.Id, &tagging.NamaTagging, &tagging.KeteranganTagging)
	if err != nil {
		return domain.TaggingPokin{}, fmt.Errorf("error saat mengambil tagging berdasarkan id: %v", err)
	}
	return tagging, nil
}
