package domain

// MasterIndikatorProgramKegiatan merepresentasikan baris dari tb_master_indikator_program_kegiatan.
// Kode berisi kode_program / kode_kegiatan / kode_subkegiatan.
type MasterIndikatorProgramKegiatan struct {
	Id            int
	KodeIndikator string
	Kode          string // kode_program / kode_kegiatan / kode_subkegiatan
	Indikator     string
}

// TargetProgramKegiatan merepresentasikan baris dari tb_target_program_kegiatan.
// Target ini spesifik per OPD (kode_opd), tahun, dan jenis (renstra/renja/dll).
type TargetProgramKegiatan struct {
	Id            int
	KodeIndikator string
	KodeOpd       string
	Jenis         string // renstra / renja / rankhir / penetapan / dll
	Target        string
	Satuan        string
	Tahun         int
}
