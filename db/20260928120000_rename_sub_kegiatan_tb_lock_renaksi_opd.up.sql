ALTER TABLE tb_lock_renaksi_opd CHANGE COLUMN sub_kegiatan nama_subkegiatan TEXT;
ALTER TABLE tb_lock_renaksi_opd ADD COLUMN kode_subkegiatan TEXT NULL AFTER nama_subkegiatan;
