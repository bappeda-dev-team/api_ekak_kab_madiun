 ALTER TABLE tb_subkegiatan
    ADD COLUMN tahun VARCHAR(20),
    ADD COLUMN kode_opd VARCHAR(255),
    ADD COLUMN status VARCHAR(255) NOT NULL DEFAULT 'belum_diambil';