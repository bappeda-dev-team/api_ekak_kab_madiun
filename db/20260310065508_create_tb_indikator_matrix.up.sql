CREATE TABLE tb_indikator_matrix (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_indikator VARCHAR(255) NOT NULL DEFAULT '' UNIQUE,
    kode VARCHAR(255) NOT NULL DEFAULT '',
    kode_opd VARCHAR(255) NOT NULL DEFAULT '',
    tujuan_opd_id INT,
    sasaran_opd_id INT,
    jenis VARCHAR(255) NOT NULL DEFAULT '',
    rumus_perhitungan TEXT,
    sumber_data TEXT,
    definisi_operasional TEXT,
    indikator VARCHAR(255) NOT NULL DEFAULT '',
    tahun VARCHAR(20) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;