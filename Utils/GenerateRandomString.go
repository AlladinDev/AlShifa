package utils

import (
	"math/rand"
	"time"
)

func init() {
	// Seed the random number generator ONCE
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)

	for i := range length {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
