package service

import (
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pkopd"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestReplaceKode(t *testing.T) {
	tests := []struct {
		name     string
		kode     string
		kodeOpd  string
		expected string
	}{
		{
			name:     "normal case",
			kode:     "X.XX.01.2.01.0001",
			kodeOpd:  "5.01.5.05.0.00.01.0000",
			expected: "X.XX.01.2.01.0001",
		},
		{
			name:     "tidak terdeteksi X.XX",
			kode:     "5.01.01.2.01.0001",
			kodeOpd:  "5.01.5.05.0.00.01.0000",
			expected: "5.01.01.2.01.0001",
		},
		{
			name:     "tidak sama dengan kode opd",
			kode:     "5.99.01.2.01.0001",
			kodeOpd:  "5.01.5.05.0.00.01.0000",
			expected: "5.01.01.2.01.0001",
		},
		{
			name:     "invalid kode",
			kode:     "-",
			kodeOpd:  "5.01.5.05.0.00.01.0000",
			expected: "-",
		},
		{
			name:     "invalid kode opd",
			kode:     "5.21.01.2.01.0001",
			kodeOpd:  "--",
			expected: "5.21.01.2.01.0001",
		},
		{
			name:     "invalid kode opd dan butuh",
			kode:     "X.XX.01.2.01.0001",
			kodeOpd:  "--",
			expected: "X.XX.01.2.01.0001",
		},
		{
			name:     "empty",
			kode:     "",
			kodeOpd:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceKode(tt.kode, tt.kodeOpd)
			if result != tt.expected {
				t.Errorf("replaceKode(%q, %q) = %q; want %q",
					tt.kode, tt.kodeOpd, result, tt.expected)
			}
		})
	}
}

func TestBuildAtasanMap(t *testing.T) {
	tests := []struct {
		name     string
		rekins   []rencanakinerja.RencanaKinerjaResponse
		expected map[string][]pkopd.AtasanCandidate
	}{
		{
			name: "level 5 -> only level 4",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 1, PegawaiId: "kadis", NamaPegawai: "Kadis", LevelPohon: 4},
				{IdPohon: 2, IdParentPohon: 1, PegawaiId: "kabid", NamaPegawai: "Kabid", LevelPohon: 5},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"kabid": {
					{IdPegawai: "kadis", NamaPegawai: "Kadis", LevelPegawai: 4, IdPohonAtasan: 1, IdParentPohonAtasan: 0},
				},
			},
		},
		{
			name: "level 6 -> level 4 and 5",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 1, PegawaiId: "kadis", LevelPohon: 4},
				{IdPohon: 2, IdParentPohon: 1, PegawaiId: "kabid", LevelPohon: 5},
				{IdPohon: 3, IdParentPohon: 2, PegawaiId: "subkor", LevelPohon: 6},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"kabid": {
					{IdPegawai: "kadis", LevelPegawai: 4, IdPohonAtasan: 1, IdParentPohonAtasan: 0},
				},
				"subkor": {
					{IdPegawai: "kadis", LevelPegawai: 4, IdPohonAtasan: 1, IdParentPohonAtasan: 0},
					{IdPegawai: "kabid", LevelPegawai: 5, IdPohonAtasan: 2, IdParentPohonAtasan: 1},
				},
			},
		},
		{
			name: "level 7 -> only level 6 segaris",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 3, PegawaiId: "subkor-1", LevelPohon: 6},
				{IdPohon: 3, PegawaiId: "subkor-2", LevelPohon: 6},
				{IdPohon: 1, PegawaiId: "not-included", LevelPohon: 6},
				{IdPohon: 4, IdParentPohon: 3, PegawaiId: "staff", LevelPohon: 7},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"staff": {
					{IdPegawai: "subkor-1", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
					{IdPegawai: "subkor-2", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
				},
			},
		},
		{
			name: "level 8 -> level 6 segaris",
			// staff-8 -> staff -> subkor
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 3, PegawaiId: "subkor-1", LevelPohon: 6},
				{IdPohon: 3, PegawaiId: "subkor-1", LevelPohon: 6},
				{IdPohon: 3, PegawaiId: "subkor-2", LevelPohon: 6},
				{IdPohon: 8, PegawaiId: "subkor-no", LevelPohon: 6},
				{IdPohon: 4, IdParentPohon: 3, PegawaiId: "staff", LevelPohon: 7},
				{IdPohon: 5, IdParentPohon: 4, PegawaiId: "staff-8", LevelPohon: 8},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"staff": {
					{IdPegawai: "subkor-1", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
					{IdPegawai: "subkor-2", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
				},
				"staff-8": {
					{IdPegawai: "subkor-1", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
					{IdPegawai: "subkor-2", LevelPegawai: 6, IdPohonAtasan: 3, IdParentPohonAtasan: 0},
				},
			},
		},
		{
			name: "level 4 -> no atasan",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 1, PegawaiId: "kadis", LevelPohon: 4},
			},
			expected: map[string][]pkopd.AtasanCandidate{},
		},
		{
			name: "multiple candidates same parent (level 5 -> 4)",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 1, PegawaiId: "kadis1", LevelPohon: 4},
				{IdPohon: 1, PegawaiId: "kadis2", LevelPohon: 4},
				{IdPohon: 2, IdParentPohon: 1, PegawaiId: "kabid", LevelPohon: 5},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"kabid": {
					{IdPegawai: "kadis1", LevelPegawai: 4, IdPohonAtasan: 1, IdParentPohonAtasan: 0},
					{IdPegawai: "kadis2", LevelPegawai: 4, IdPohonAtasan: 1, IdParentPohonAtasan: 0},
				},
			},
		},
		{
			name: "parent not found",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{IdPohon: 2, IdParentPohon: 99, PegawaiId: "staff", LevelPohon: 7},
			},
			expected: map[string][]pkopd.AtasanCandidate{},
		},
		{
			name:     "empty rekins, no candidates",
			rekins:   []rencanakinerja.RencanaKinerjaResponse{},
			expected: map[string][]pkopd.AtasanCandidate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildAtasanMap(tt.rekins)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildLevel4Candidates(t *testing.T) {
	tests := []struct {
		name          string
		sasaranPemdas []domain.AllSasaranPemdaPk
		expected      []pkopd.AtasanCandidate
	}{
		{
			name:          "sasaranPemda empty, no candidates",
			sasaranPemdas: []domain.AllSasaranPemdaPk{},
			expected:      []pkopd.AtasanCandidate{},
		},
		{
			name: "expected flow",
			sasaranPemdas: []domain.AllSasaranPemdaPk{
				{JabatanKepalaPemda: "KEPALA DAERAH XX",
					NamaKepalaPemda: "namakepaladaerah",
					NipKepalaPemda:  "---",
					SasaranPemdaId:  123,
					SasaranPemda:    "sasaranpemda",
				},
			},
			expected: []pkopd.AtasanCandidate{
				{
					IdPegawai:           "---",
					NamaPegawai:         "namakepaladaerah",
					LevelPegawai:        3,
					IdPohonAtasan:       0,
					IdParentPohonAtasan: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildLevel4Candidates(tt.sasaranPemdas)
			assert.Equal(t, tt.expected, result)
		})
	}
}
