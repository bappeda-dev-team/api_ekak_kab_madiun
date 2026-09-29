package web

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// StringID adalah ID numerik yang pada JSON boleh dikirim sebagai angka maupun string,
// tetapi selalu disimpan sebagai string agar presisi ID yang besar tidak hilang.
type StringID string

// normalizeID memangkas spasi dan menolak nilai yang bukan bilangan bulat.
// String kosong tetap diterima karena `null` maupun `""` berarti ID tidak diisi;
// emptiness tetap ditangani oleh validator `required` dan oleh Int.
func normalizeID(raw string) (StringID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	if _, err := strconv.ParseInt(trimmed, 10, 64); err != nil {
		return "", fmt.Errorf("id harus berupa angka bulat")
	}
	return StringID(trimmed), nil
}

// UnmarshalJSON menerima `12`, `"12"`, dan `null`, serta menolak desimal,
// huruf, boolean, dan objek.
func (s *StringID) UnmarshalJSON(data []byte) error {
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		id, err := normalizeID(string(num))
		if err != nil {
			return err
		}
		*s = id
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("id harus berupa angka atau string angka")
	}

	id, err := normalizeID(str)
	if err != nil {
		return err
	}
	*s = id
	return nil
}

// Int mengubah ID menjadi int untuk dipakai domain layer.
func (s StringID) Int() (int, error) {
	if s == "" {
		return 0, fmt.Errorf("id kosong")
	}
	value, err := strconv.Atoi(string(s))
	if err != nil {
		return 0, fmt.Errorf("id harus berupa angka bulat")
	}
	return value, nil
}

// MarshalJSON selalu menghasilkan string agar tipe keluaran konsisten.
func (s StringID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}
