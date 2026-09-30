CREATE TABLE IF NOT EXISTS tb_renaksi_opd (
    id INT AUTO_INCREMENT PRIMARY KEY,
    rekin_id VARCHAR(255) NOT NULL,
    sasaran_id INT NOT NULL,
    tw1 int,
    tw2 int,
    tw3 int,
    tw4 int,
    tahun VARCHAR(255) NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB;

