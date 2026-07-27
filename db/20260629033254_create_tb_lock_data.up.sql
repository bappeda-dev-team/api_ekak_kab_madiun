CREATE TABLE tb_lock_data_pemda (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    jenis      VARCHAR(100) NOT NULL COMMENT 'contoh: tujuan_pemda',
    tahun      VARCHAR(4)   NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_lock_pemda_jenis_tahun (jenis, tahun)
) ENGINE=InnoDB;