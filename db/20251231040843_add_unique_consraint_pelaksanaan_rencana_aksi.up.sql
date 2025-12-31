ALTER TABLE tb_pelaksanaan_rencana_aksi 
ADD UNIQUE KEY uk_rencana_aksi_bulan (rencana_aksi_id, bulan);