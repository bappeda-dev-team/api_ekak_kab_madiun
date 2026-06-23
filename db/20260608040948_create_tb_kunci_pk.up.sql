CREATE TABLE tb_kunci_pk (
    id BIGINT NOT NULL AUTO_INCREMENT,
    id_pegawai VARCHAR(50) NOT NULL,
    kode_opd VARCHAR(50) NOT NULL,
    tahun INT NOT NULL,
    dikunci_oleh VARCHAR(50) NOT NULL,
    dikunci_pada DATETIME NOT NULL,
    status_pk VARCHAR(20) NOT NULL,
    pk_terkunci BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_kunci_pk_pegawai_tahun (
        id_pegawai,
        tahun
    ),
    INDEX idx_kunci_pk_kode_opd (kode_opd),
    INDEX idx_kunci_pk_tahun (tahun)
);
