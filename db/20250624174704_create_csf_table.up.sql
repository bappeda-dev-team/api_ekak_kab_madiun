CREATE TABLE tb_csf (
    id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,
    pohon_id INT NOT NULL UNIQUE,
    pernyataan_kondisi_strategis VARCHAR(255),
    alasan_kondisi_strategis VARCHAR(255),
    data_terukur VARCHAR(255),
    kondisi_terukur VARCHAR(255),
    kondisi_wujud VARCHAR(255),
    tahun INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tb_csf_pohon FOREIGN KEY (pohon_id)
        REFERENCES tb_pohon_kinerja(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE = InnoDB;
