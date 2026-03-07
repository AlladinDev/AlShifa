package utils

import (
	"AlShifa/constants"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a token with userID, role, and custom expiration
func GenerateJWT(claimOptions *constants.JwtCustomClaims) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	claims := &constants.JwtCustomClaims{
		UserID:     claimOptions.UserID,
		Role:       claimOptions.Role,
		Email:      claimOptions.Email,
		Mobile:     claimOptions.Mobile,
		IsVerified: claimOptions.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JwtExpiryTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*constants.JwtCustomClaims, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		panic("JWt SECRET environment variable is not set")
	}

	claims := &constants.JwtCustomClaims{}

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
