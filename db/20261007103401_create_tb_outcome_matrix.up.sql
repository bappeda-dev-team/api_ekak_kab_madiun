CREATE TABLE IF NOT EXISTS tb_outcome_matrix (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_subkegiatan VARCHAR(255) NOT NULL,
    kode VARCHAR(255) NOT NULL,
    outcome TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_outcome_matrix_kode (kode),
    INDEX idx_outcome_matrix_kode_subkegiatan (kode_subkegiatan),
    INDEX idx_outcome_matrix_kode_pair (kode, kode_subkegiatan)
) ENGINE=InnoDB;
