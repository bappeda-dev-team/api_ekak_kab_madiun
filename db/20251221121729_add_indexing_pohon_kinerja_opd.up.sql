-- Tambahkan index untuk query tematik (clone_from)
CREATE INDEX idx_pokin_clone_from ON tb_pohon_kinerja(clone_from, level_pohon);

-- Index untuk parent-child relationship
CREATE INDEX idx_pokin_parent_level ON tb_pohon_kinerja(parent, level_pohon, id);

-- Index untuk pelaksana batch query
CREATE INDEX idx_pelaksana_pokin_id ON tb_pelaksana_pokin(pohon_kinerja_id);

-- Index untuk indikator batch query
CREATE INDEX idx_indikator_pokin_id ON tb_indikator(pokin_id);

-- Index untuk target batch query  
CREATE INDEX idx_target_indikator_id ON tb_target(indikator_id);

-- Index untuk tagging batch query
CREATE INDEX idx_tagging_pokin_id ON tb_tagging_pokin(id_pokin);

-- Index untuk review batch query
CREATE INDEX idx_review_pohon_kinerja ON tb_review(id_pohon_kinerja);

-- Index untuk keterangan tagging
CREATE INDEX idx_keterangan_tagging ON tb_keterangan_tagging_program_unggulan(id_tagging);