package service

import (
	"testing"
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
			expected: "5.01.01.2.01.0001",
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
			expected: "5.99.01.2.01.0001",
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
