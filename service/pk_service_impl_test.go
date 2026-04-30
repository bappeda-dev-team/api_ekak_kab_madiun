package service

import (
	"ekak_kabupaten_madiun/model/web/opdmaster"
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
		// [pegawaiIdBawahan] -> []candiateAtasans
		{
			name: "single parent relationship",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{
					IdPohon:     1,
					PegawaiId:   "atasan",
					NamaPegawai: "Atasan A",
					LevelPohon:  4,
					KodeOpd: opdmaster.OpdResponseForAll{
						KodeOpd: "OPD1",
						NamaOpd: "Dinas A",
					},
				},
				{
					IdPohon:       2,
					IdParentPohon: 1,
					LevelPohon:    5,
					PegawaiId:     "bawahan123",
					NamaPegawai:   "Bawahan B",
				},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"bawahan123": {
					{
						IdPegawai:    "atasan",
						NamaPegawai:  "Atasan A",
						LevelPegawai: 4,
						KodeOpd:      "OPD1",
						NamaOpd:      "Dinas A",
					},
				},
			},
		},
		{
			name: "no parent should be skipped",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{
					IdPohon:     1,
					PegawaiId:   "pegawai",
					NamaPegawai: "Tanpa Atasan",
				},
			},
			expected: map[string][]pkopd.AtasanCandidate{},
		},
		{
			name: "duplicate parent should be unique",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{
					IdPohon:            1,
					PegawaiId:          "atasan",
					NamaPegawai:        "Atasan A",
					NamaRencanaKinerja: "rekin-1",
				},
				{
					IdPohon:            3,
					PegawaiId:          "atasan",
					NamaPegawai:        "Atasan A",
					NamaRencanaKinerja: "rekin-2",
				},
				{
					IdPohon:       2,
					IdParentPohon: 1,
					PegawaiId:     "bawahan1",
				},
				{
					IdPohon:       4,
					IdParentPohon: 3,
					PegawaiId:     "bawahan2", // sama pegawai
				},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"bawahan1": {
					{
						IdPegawai:   "atasan",
						NamaPegawai: "Atasan A",
					},
				},
				"bawahan2": {
					{
						IdPegawai:   "atasan",
						NamaPegawai: "Atasan A",
					},
				},
			},
		},
		{
			name: "parent not found should skip",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{
					IdPohon:       2,
					IdParentPohon: 99, // tidak ada
					PegawaiId:     "pegawai",
				},
			},
			expected: map[string][]pkopd.AtasanCandidate{},
		},
		{
			name: "multiple pegawai different parents",
			rekins: []rencanakinerja.RencanaKinerjaResponse{
				{
					IdPohon:     1,
					PegawaiId:   "atasan1",
					NamaPegawai: "Atasan 1",
				},
				{
					IdPohon:     2,
					PegawaiId:   "atasan2",
					NamaPegawai: "Atasan 2",
				},
				{
					IdPohon:     2,
					PegawaiId:   "atasan3",
					NamaPegawai: "Atasan 3",
				},
				{
					IdPohon:     12,
					PegawaiId:   "atasan3",
					NamaPegawai: "Atasan 3",
				},
				{
					IdPohon:       3,
					IdParentPohon: 1,
					PegawaiId:     "pegawai1",
				},
				{
					IdPohon:       4,
					IdParentPohon: 2,
					PegawaiId:     "pegawai2",
				},
			},
			expected: map[string][]pkopd.AtasanCandidate{
				"pegawai1": {
					{IdPegawai: "atasan1", NamaPegawai: "Atasan 1"},
				},
				"pegawai2": {
					{IdPegawai: "atasan2", NamaPegawai: "Atasan 2"},
					{IdPegawai: "atasan3", NamaPegawai: "Atasan 3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildAtasanMap(tt.rekins)

			assert.Equal(t, tt.expected, result)
		})
	}
}
