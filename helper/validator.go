package helper

import (
	"errors"
	"strconv"
)

func ValidateTahunRange(awalStr, akhirStr string) (int, int, error) {
	if awalStr == "" || akhirStr == "" {
		return 0, 0, errors.New("tahun_awal dan tahun_akhir wajib diisi")
	}

	awal, err1 := strconv.Atoi(awalStr)
	akhir, err2 := strconv.Atoi(akhirStr)

	if err1 != nil || err2 != nil {
		return 0, 0, errors.New("tahun harus angka")
	}

	if awal > akhir {
		return 0, 0, errors.New("range tidak valid")
	}

	return awal, akhir, nil
}
