package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func HashJson(item any) (string, error) {
	bytes, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:]), nil
}
