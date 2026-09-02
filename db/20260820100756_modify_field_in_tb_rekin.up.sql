UPDATE tb_rencana_kinerja 
SET sasaranopd_id = 0 
WHERE sasaranopd_id IS NULL;

ALTER TABLE tb_rencana_kinerja MODIFY COLUMN sasaranopd_id INT NOT NULL DEFAULT 0;