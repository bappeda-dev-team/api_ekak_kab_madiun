    CREATE TABLE tb_keterangan_tagging_program_prioritas_pusat (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_tagging INT NOT NULL,
        kode_program_prioritas_pusat VARCHAR(255) NOT NULL,
        tahun VARCHAR(10),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tb_tagging_keterangan_program_prioritas_pusat FOREIGN KEY (id_tagging)
        REFERENCES tb_tagging_pokin(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE

    ) ENGINE = InnoDB;