package helper

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func GenerateID(baseCode string) string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5] // Mengambil 5 karakter pertama dari UUID
	return fmt.Sprintf("%s-%s-%s", baseCode, currentYear, uuid)
}
