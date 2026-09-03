CREATE TABLE tb_program_opd_terpilih (
    id INT AUTO_INCREMENT PRIMARY KEY,

    pohon_kinerja_id INT NOT NULL,
    program_opd_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uniq_program_terpilih (pohon_kinerja_id, program_opd_id),

    CONSTRAINT fk_pot_pohon
        FOREIGN KEY (pohon_kinerja_id)
        REFERENCES tb_pohon_kinerja(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_pot_program
        FOREIGN KEY (program_opd_id)
        REFERENCES tb_pohon_kinerja(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);