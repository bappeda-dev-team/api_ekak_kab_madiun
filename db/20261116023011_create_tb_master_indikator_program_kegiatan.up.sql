CREATE TABLE IF NOT EXISTS tb_master_indikator_program_kegiatan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_indikator VARCHAR(100) NOT NULL,
    kode VARCHAR(100) NOT NULL COMMENT 'kode_program / kode_kegiatan / kode_subkegiatan',
    indikator TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_master_indikator_pk_kode_indikator (kode_indikator),
    KEY idx_master_indikator_pk_kode (kode)
);

