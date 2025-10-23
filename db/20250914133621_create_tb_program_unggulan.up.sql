CREATE TABLE tb_program_unggulan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama_tagging VARCHAR(255) NOT NULL,
    kode_program_unggulan VARCHAR(255) NOT NULL,
    keterangan_program_unggulan TEXT,
    keterangan TEXT,
    tahun_awal VARCHAR(255) NOT NULL,
    tahun_akhir VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);