CREATE TABLE datamaster_rb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    jenis_rb VARCHAR(10) NOT NULL,
    kegiatan_utama TEXT NOT NULL,
    keterangan TEXT,
    tahun_baseline INT,
    tahun_next INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    current_version INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    last_updated_by INT,
    CONSTRAINT fk_last_updated_by FOREIGN KEY (last_updated_by)
        REFERENCES tb_users(id)
) ENGINE=InnoDB;
