CREATE TABLE struktur_organisasi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nip_bawahan VARCHAR(20) NOT NULL,
    nip_atasan VARCHAR(20) NOT NULL,
    kode_opd VARCHAR(255) NOT NULL,
    tahun INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_bawahan_opd_tahun
        UNIQUE (nip_bawahan, kode_opd, tahun),

    CONSTRAINT chk_not_self
        CHECK (nip_bawahan <> nip_atasan),

    INDEX idx_opd_tahun (kode_opd, tahun)
);
