CREATE TABLE tb_target_pemda (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_indikator VARCHAR(255) NOT NULL,
    target VARCHAR(255),
    satuan VARCHAR(255),
    tahun VARCHAR(20),
    jenis VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_target_pemda_kode_indikator (kode_indikator),
    INDEX idx_target_pemda_jenis (jenis),
    INDEX idx_target_pemda_kode_tahun_jenis (kode_indikator, tahun, jenis)
) ENGINE=InnoDB;