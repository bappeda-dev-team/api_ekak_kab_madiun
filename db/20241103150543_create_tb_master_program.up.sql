CREATE TABLE tb_master_program(
    id VARCHAR(255) PRIMARY KEY,
    nama_program TEXT NOT NULL,
    kode_program VARCHAR(255) NOT NULL,
    kode_opd VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    tahun INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;