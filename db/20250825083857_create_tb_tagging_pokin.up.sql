CREATE TABLE tb_tagging_pokin (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_pokin INT NOT NULL,
    nama_tagging VARCHAR(255) NOT NULL,
    keterangan_tagging TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_tb_tagging_pokin_pohon FOREIGN KEY (id_pokin)
        REFERENCES tb_pohon_kinerja(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)ENGINE=InnoDB;