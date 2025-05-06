CREATE TABLE tb_sasaran_opd (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pokin_id INT NOT NULL,
    nama_sasaran_opd TEXT,
    tahun_awal VARCHAR(30),
    tahun_akhir VARCHAR(30),
    jenis_periode VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;