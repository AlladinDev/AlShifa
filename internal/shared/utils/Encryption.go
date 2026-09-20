package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

var encryptionKey = []byte(
	"12345678901234567890123456789012", // exactly 32 bytes
)

// Encrypt accepts any data type
func Encrypt(data any) (string, error) {
	if data == nil {
		return "", errors.New("data to encrypt cannot be nil")
	}

	// Convert anything to JSON bytes
	plainBytes, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	// Create random nonce
	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(
		rand.Reader,
		nonce,
	); err != nil {

		return "", err
	}

	// Encrypt
	cipherText := gcm.Seal(
		nonce,
		nonce,
		plainBytes,
		nil,
	)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt returns JSON bytes
func Decrypt(cipherText string) ([]byte, error) {

	data, err := base64.StdEncoding.DecodeString(cipherText)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(encryptionKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {

		return nil, errors.New(
			"invalid encrypted data",
		)
	}

	nonce := data[:nonceSize]

	encrypted := data[nonceSize:]

	plain, err := gcm.Open(
		nil,
		nonce,
		encrypted,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return plain, nil
}
