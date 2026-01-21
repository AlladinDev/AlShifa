package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// VerifyPasswordArgon2id verifies a password against an Argon2id hash
func VerifyPasswordArgon2id(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")

	// Remove optional leading $
	if parts[0] == "" {
		parts = parts[1:]
	}

	if len(parts) != 5 {
		return false, errors.New("invalid hash format")
	}

	if parts[0] != "argon2id" {
		return false, errors.New("not an argon2id hash")
	}

	if parts[1] != "v=19" {
		return false, errors.New("unsupported argon2 version")
	}

	var memory uint32
	var time uint32
	var threads uint8

	// ✅ THIS IS THE FIX
	_, err := fmt.Sscanf(
		parts[2],
		"m=%d,t=%d,p=%d",
		&memory,
		&time,
		&threads,
	)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(hash)),
	)

	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}

	return false, nil
}
