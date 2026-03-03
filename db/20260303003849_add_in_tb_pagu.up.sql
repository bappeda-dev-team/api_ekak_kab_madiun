ALTER TABLE tb_pagu 
    MODIFY COLUMN pagu BIGINT UNSIGNED,
    ADD COLUMN kode_subkegiatan VARCHAR(50),
    ADD COLUMN kode_opd VARCHAR(50),
    ADD COLUMN jenis VARCHAR(100),
    MODIFY COLUMN id INT AUTO_INCREMENT PRIMARY KEY,
    ADD UNIQUE KEY unique_pagu_subkegiatan (kode_opd, kode_subkegiatan, tahun, jenis);
