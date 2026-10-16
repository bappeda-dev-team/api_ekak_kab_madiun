ALTER TABLE tb_outcome_matrix
    CHANGE COLUMN kode_subkegiatan kode_opd VARCHAR(255) NOT NULL,
    ADD COLUMN jenis VARCHAR(255) NOT NULL DEFAULT '';

DROP INDEX idx_outcome_matrix_kode_subkegiatan ON tb_outcome_matrix;
DROP INDEX idx_outcome_matrix_kode_pair ON tb_outcome_matrix;

