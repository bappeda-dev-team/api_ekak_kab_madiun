-- Index composite untuk optimasi query FindAll
CREATE INDEX idx_pokin_kode_tahun_status 
ON tb_pohon_kinerja(kode_opd, tahun, status, level_pohon, id);

-- Index untuk ORDER BY jika diperlukan secara terpisah
CREATE INDEX idx_pokin_status_level_id 
ON tb_pohon_kinerja(status, level_pohon, id);