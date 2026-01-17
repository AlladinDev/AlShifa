package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argon2idTime    uint32 = 2
	argon2idMemory  uint32 = 64 * 1024 // 64 MB
	argon2idThreads uint8  = 4
	argon2idKeyLen  uint32 = 32
	argon2idSaltLen        = 16
)

// HashPasswordArgon2id hashes a password using Argon2id
// Format:
// $argon2id$v=19$m=65536,t=2,p=4$<salt>$<hash>
func HashPasswordArgon2id(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, argon2idSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Derive key (Argon2id)
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2idTime,
		argon2idMemory,
		argon2idThreads,
		argon2idKeyLen,
	)

	// Encode salt and hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2idMemory,
		argon2idTime,
		argon2idThreads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}
