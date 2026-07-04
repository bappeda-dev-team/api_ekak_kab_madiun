package tujuanpemda

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type TargetInput float64

func (t *TargetInput) UnmarshalJSON(data []byte) error {
	// Path 1: JSON number → 85 atau 85.5
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		f, err := num.Float64()
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return fmt.Errorf("target format tidak valid: harus berupa angka")
		}
		*t = TargetInput(f)
		return nil
	}
	// Path 2: JSON string → "85" atau "85.5"
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("target format tidak valid")
	}
	if err := validateTargetRawString(str); err != nil {
		return err
	}
	str = strings.TrimSpace(str)
	if str == "" || str == "-" {
		*t = 0
		return nil
	}
	f, _ := strconv.ParseFloat(str, 64)
	*t = TargetInput(f)
	return nil
}
func (t TargetInput) Float64() float64 {
	return float64(t)
}

const TargetInvalidFormat = "format tidak valid"

// TargetDisplay — di JSON: number jika valid, string "format tidak valid" jika tidak, 0 jika kosong.
type TargetDisplay struct {
	empty   bool
	invalid bool
	value   float64
}

func NewTargetDisplayFromString(raw string) TargetDisplay {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return TargetDisplay{empty: true}
	}
	// Tolak koma sebagai desimal: "89,4"
	if strings.Contains(raw, ",") {
		return TargetDisplay{invalid: true}
	}
	// Tolak huruf/alphabet
	for _, r := range raw {
		if unicode.IsLetter(r) {
			return TargetDisplay{invalid: true}
		}
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return TargetDisplay{invalid: true}
	}
	return TargetDisplay{value: f}
}
func (t TargetDisplay) MarshalJSON() ([]byte, error) {
	if t.invalid {
		return json.Marshal(TargetInvalidFormat)
	}
	if t.empty {
		return json.Marshal(float64(0))
	}
	return json.Marshal(t.value)
}

type LayerTargetItemRequest struct {
	KodeIndikator string      `json:"kode_indikator"`
	Tahun         string      `json:"tahun"`
	Target        TargetInput `json:"target"`
	Satuan        string      `json:"satuan"`
}
type LayerTargetBatchRequest struct {
	Targets []LayerTargetItemRequest `json:"targets"`
}

type LayerTargetUpdateItemRequest struct {
	Id     int         `json:"id"`
	Target TargetInput `json:"target"`
	Satuan string      `json:"satuan"`
}
type LayerTargetUpdateBatchRequest struct {
	Targets []LayerTargetUpdateItemRequest `json:"targets"`
}

func validateTargetRawString(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return nil // placeholder kosong → 0
	}
	// Tolak koma desimal: "89,4"
	if strings.Contains(raw, ",") {
		return fmt.Errorf("target format tidak valid: gunakan titik (.) sebagai desimal, bukan koma (,)")
	}
	// Tolak huruf: "abc", "85persen"
	for _, r := range raw {
		if unicode.IsLetter(r) {
			return fmt.Errorf("target format tidak valid: tidak boleh mengandung huruf")
		}
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf("target format tidak valid: harus berupa angka")
	}
	return nil
}
