-- 1. Tambah kolom nullable terlebih dahulu agar aman untuk data lama
ALTER TABLE tb_renaksi_opd ADD COLUMN urutan INT NULL;

-- 2. Backfill urutan 1..N per (sasaran_id, tahun)
UPDATE tb_renaksi_opd ro
JOIN (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY sasaran_id, tahun
            ORDER BY rekin_id
        ) AS rn
    FROM tb_renaksi_opd
) t ON t.id = ro.id
SET ro.urutan = t.rn;

ALTER TABLE tb_renaksi_opd MODIFY COLUMN urutan INT NOT NULL;

ALTER TABLE tb_renaksi_opd ADD UNIQUE KEY uq_renaksi_opd_urutan (sasaran_id, tahun, urutan);
