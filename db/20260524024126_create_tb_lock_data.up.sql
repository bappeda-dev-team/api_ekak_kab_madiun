CREATE TABLE tb_lock_data (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    jenis_data VARCHAR(50)  NOT NULL,
    kode_opd   VARCHAR(255) NOT NULL,
    tahun      VARCHAR(4)   NOT NULL,
    locked_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_lock (jenis_data, kode_opd, tahun)
);