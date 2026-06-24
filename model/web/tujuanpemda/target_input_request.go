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
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		f, err := num.Float64()
		if err != nil {
			return fmt.Errorf("target format tidak valid")
		}
		*t = TargetInput(f)
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("target format tidak valid")
	}
	str = strings.TrimSpace(str)
	if str == "" || str == "-" {
		*t = 0
		return nil
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("target format tidak valid")
	}
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
