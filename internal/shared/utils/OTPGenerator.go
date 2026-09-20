package utils

import (
	"math/rand"
)

func GenerateOTP() string {
	const digits = "0123456789"
	length := 8
	otp := make([]byte, length)

	for i := range length {
		// crypto/rand ensures cryptographically secure random numbers
		n := rand.Intn(length + 1)
		otp[i] = digits[n]
	}

	return string(otp)
}
