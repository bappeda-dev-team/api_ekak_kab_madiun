CREATE TABLE tb_bidangurusan_terpilih (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_bidang_urusan VARCHAR(255) NOT NULL,
    kode_opd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_bidang_urusan_terpilih_bidang_urusan
        FOREIGN KEY (kode_bidang_urusan)
        REFERENCES tb_bidang_urusan(kode_bidang_urusan)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;