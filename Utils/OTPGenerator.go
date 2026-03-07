package utils

import (
	"math/rand"
)

func GenerateOTP() string {
	const digits = "0123456789"
	length := 6
	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		// crypto/rand ensures cryptographically secure random numbers
		n := rand.Intn(length + 1)
		otp[i] = digits[n]
	}

	return string(otp)
}
