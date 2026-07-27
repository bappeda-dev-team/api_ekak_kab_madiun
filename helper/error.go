package helper

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	MsgTargetFormatInvalid    = "target format tidak valid"
	MsgTargetFormatKoma       = "target format tidak valid: gunakan titik (.) sebagai desimal, bukan koma (,)"
	MsgTargetFormatHuruf      = "target format tidak valid: tidak boleh mengandung huruf"
	MsgTargetFormatBukanAngka = "target format tidak valid: harus berupa angka"
)

func ValidateTargetRawString(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return nil
	}
	if strings.Contains(raw, ",") {
		return fmt.Errorf(MsgTargetFormatKoma)
	}
	for _, r := range raw {
		if unicode.IsLetter(r) {
			return fmt.Errorf(MsgTargetFormatHuruf)
		}
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf(MsgTargetFormatBukanAngka)
	}
	return nil
}
func ValidateTargetFloat(v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf(MsgTargetFormatInvalid)
	}
	return nil
}
func TargetToDBString(v float64) (string, error) {
	if err := ValidateTargetFloat(v); err != nil {
		return "", err
	}
	if v == math.Trunc(v) {
		return strconv.FormatInt(int64(v), 10), nil
	}
	return strconv.FormatFloat(v, 'f', -1, 64), nil
}

// Wrapper pesan siap tampil ke user
func ErrTargetIndikator(indikator, tahun, pesan string) error {
	return fmt.Errorf("indikator '%s' tahun %s: %s", indikator, tahun, pesan)
}
func ErrTargetLayer(kodeIndikator, tahun, pesan string) error {
	return fmt.Errorf("nilai target tidak valid untuk kode_indikator '%s' tahun %s: %s", kodeIndikator, tahun, pesan)
}

type TargetValidationError struct {
	Message string
}

func (e *TargetValidationError) Error() string {
	return e.Message
}
func NewTargetValidationError(message string) error {
	return &TargetValidationError{Message: message}
}

// IsTargetValidationError — dipakai controller untuk tentukan 400 vs 500
func IsTargetValidationError(err error) bool {
	if err == nil {
		return false
	}
	var tv *TargetValidationError
	if errors.As(err, &tv) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, MsgTargetFormatInvalid) ||
		strings.Contains(msg, "nilai target tidak valid")
}
