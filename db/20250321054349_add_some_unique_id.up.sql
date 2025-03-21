ALTER TABLE tb_urusan ADD CONSTRAINT unique_id UNIQUE (id);

ALTER TABLE tb_bidang_urusan ADD CONSTRAINT unique_id UNIQUE (id);

ALTER TABLE tb_master_kegiatan ADD CONSTRAINT unique_id UNIQUE (id);

ALTER TABLE tb_master_program ADD CONSTRAINT unique_id UNIQUE (id);

ALTER TABLE tb_subkegiatan ADD CONSTRAINT unique_id UNIQUE (id);

