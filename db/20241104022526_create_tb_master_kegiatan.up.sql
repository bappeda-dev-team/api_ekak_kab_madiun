CREATE TABLE tb_master_kegiatan (
	`id` VARCHAR(255) NOT NULL,
	`nama_kegiatan` VARCHAR(255) NOT NULL,
    `kode_kegiatan` VARCHAR(255) NOT NULL, 
	`kode_opd` VARCHAR(255) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;