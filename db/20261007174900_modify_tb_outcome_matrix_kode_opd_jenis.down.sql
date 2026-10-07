ALTER TABLE tb_outcome_matrix
    DROP COLUMN jenis,
    CHANGE COLUMN kode_opd kode_subkegiatan VARCHAR(255) NOT NULL;

CREATE INDEX idx_outcome_matrix_kode_subkegiatan ON tb_outcome_matrix (kode_subkegiatan);
CREATE INDEX idx_outcome_matrix_kode_pair ON tb_outcome_matrix (kode, kode_subkegiatan);
