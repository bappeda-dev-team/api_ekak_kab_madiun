CREATE TABLE tb_indikator_matrix_pemda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_indikator VARCHAR(255) NOT NUll UNIQUE,
    tujuan_pemda_id INT DEFAULT 0,
    sasaran_pemda_id INT DEFAULT 0,
    indikator TEXT,
    rumus_perhitungan TEXT,
    sumber_data TEXT,
    definisi_operasional TEXT,
    jenis VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_indikator_matrix_pemda_tujuan (tujuan_pemda_id),
    INDEX idx_indikator_matrix_pemda_jenis (jenis)
) ENGINE=InnoDB;