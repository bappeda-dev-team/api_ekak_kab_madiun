CREATE TABLE tb_indikator_ikk (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_ikk INT(50),
    kode_opd VARCHAR(255),
    kode_bidang_urusan VARCHAR(255),
    indikator VARCHAR(255),
    tahun INT(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);