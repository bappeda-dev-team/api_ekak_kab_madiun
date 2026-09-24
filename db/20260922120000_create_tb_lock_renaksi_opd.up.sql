CREATE TABLE tb_lock_renaksi_opd (
    id            INT AUTO_INCREMENT PRIMARY KEY,
    kode_opd      VARCHAR(255) NOT NULL,
    tahun         VARCHAR(4)   NOT NULL,
    sasaran_id    INT          NOT NULL,
    rekin_id      VARCHAR(255) NOT NULL,
    aksi_kegiatan TEXT,
    sub_kegiatan  TEXT,
    anggaran      BIGINT DEFAULT 0,
    nama_pemilik  VARCHAR(255),
    tw1           INT DEFAULT 0,
    tw2           INT DEFAULT 0,
    tw3           INT DEFAULT 0,
    tw4           INT DEFAULT 0,
    UNIQUE KEY uq_lock_renaksi_opd (kode_opd, tahun, sasaran_id, rekin_id)
);
