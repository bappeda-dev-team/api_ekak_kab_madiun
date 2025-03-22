ALTER TABLE tb_urusan ADD CONSTRAINT unique_kode UNIQUE (kode_urusan);

ALTER TABLE tb_bidang_urusan ADD CONSTRAINT unique_kode UNIQUE (kode_bidang_urusan);

ALTER TABLE tb_master_kegiatan ADD CONSTRAINT unique_kode UNIQUE (kode_kegiatan);

ALTER TABLE tb_master_program ADD CONSTRAINT unique_kode UNIQUE (kode_program);

ALTER TABLE tb_subkegiatan ADD CONSTRAINT unique_kode UNIQUE (kode_subkegiatan);

