CREATE TABLE tb_leaderboard_hidden (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kode_opd VARCHAR(50),
    tahun VARCHAR(4),
    is_hidden BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);