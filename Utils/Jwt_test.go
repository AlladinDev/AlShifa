package utils

import (
	"AlShifa/constants"
	"testing"
)

func TestJwtGenerationAndValidation(t *testing.T) {
	t.Setenv("JWT_SECRET", "123")
	token, err := GenerateJWT(&constants.JwtCustomClaims{
		Mobile: "1234567892",
		UserID: "123",
		Role:   "user",
		Email:  "abc@gmail.com",
	})
	if err != nil {
		t.Errorf("Error generating JWT: %v", err)
	}

	if token == "" {
		t.Errorf("Generated token is empty")
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Errorf("Error validating JWT: %v", err)
	}

	if claims.UserID != "123" {
		t.Errorf("Expected UserID to be '123', got '%s'", claims.UserID)
	}

	if claims.Mobile != "1234567892" {
		t.Errorf("Expected mobile  to be '1234567892', got '%s'", claims.Mobile)
	}

	if claims.Email != "abc@gmail.com" {
		t.Errorf("Expected abc@gmail.com  to be 'abc@gmail.com', got '%s'", claims.Mobile)
	}

	if claims.Role != "user" {
		t.Errorf("Expected Role to be 'user', got '%s'", claims.Role)
	}
}
