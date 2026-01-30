package utils

import (
	"fmt"
	"math/rand"
)

func GenerateOTP(uniquePrefix string) string {

	code := rand.Intn(900000) + 100000 // ensures a number between 100000-999999
	return fmt.Sprintf("%s:%d", uniquePrefix, code)
}
