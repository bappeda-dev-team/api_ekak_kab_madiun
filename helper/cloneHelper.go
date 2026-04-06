package helper

import (
	"fmt"
	"time"
)

func GenerateKodeClone(kodeOpd string) string {
	// todays date
	currentTime := time.Now().Format("20060102")
	return fmt.Sprintf("CLONE-REKIN-%s-%s", kodeOpd, currentTime)
}
