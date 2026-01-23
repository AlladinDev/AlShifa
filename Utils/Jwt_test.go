package utils

import "testing"

func TestJwtGenerationAndValidation(t *testing.T) {
	t.Setenv("JWT_SECRET", "123")
	token, err := GenerateJWT("98766545", "user")
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

	if claims.UserID != "98766545" {
		t.Errorf("Expected UserID to be '98766545', got '%s'", claims.UserID)
	}

	if claims.Role != "user" {
		t.Errorf("Expected Role to be 'user', got '%s'", claims.Role)
	}
}
