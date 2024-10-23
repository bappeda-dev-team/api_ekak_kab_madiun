CREATE TABLE tb_subkegiatan_terpilih (
    id VARCHAR(225) NOT NULL,
    rekin_id VARCHAR(225),
    subkegiatan_id VARCHAR(225),
    pegawai_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE = InnoDB;