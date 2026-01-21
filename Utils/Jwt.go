package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struct with userID, role, and expiration
type Claims struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a token with userID, role, and custom expiration
func GenerateJWT(userID string, role string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JwtExpiryTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		panic("JWt SECRET environment variable is not set")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		// Ensure token uses HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
