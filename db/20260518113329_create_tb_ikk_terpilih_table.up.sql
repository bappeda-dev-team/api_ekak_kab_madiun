CREATE TABLE tb_ikk_terpilih (
    id INT AUTO_INCREMENT PRIMARY KEY,

    pohon_kinerja_id INT NOT NULL,
    ikk_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uniq_ikk_terpilih (pohon_kinerja_id, ikk_id),

    CONSTRAINT fk_tit_pohon
        FOREIGN KEY (pohon_kinerja_id)
        REFERENCES tb_pohon_kinerja(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_tit_program
        FOREIGN KEY (ikk_id)
        REFERENCES tb_ikk(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);