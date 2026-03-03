ALTER TABLE tb_pagu 
    DROP PRIMARY KEY,
    MODIFY COLUMN id VARCHAR(255) NOT NULL,
    MODIFY COLUMN pagu INT,
    DROP COLUMN kode_subkegiatan,
    DROP COLUMN kode_opd,
    DROP COLUMN jenis,
    DROP INDEX unique_pagu_subkegiatan;