CREATE TABLE tb_inovasi_rekin (
    id VARCHAR(255) NOT NULL,
    rekin_id VARCHAR(255) NOT NULL,
    nama_inovasi TEXT NOT NULL,
    jenis_inovasi_id INT NOT NULL,
    waktu_implementasi DATE NOT NULL,
    instansi VARCHAR(255) DEFAULT NULL,
    inovator VARCHAR(255) NOT NULL,
    kode_opd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;