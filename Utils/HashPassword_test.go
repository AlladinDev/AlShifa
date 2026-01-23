package utils

import "testing"

func TestHashPasswordAndVerifyPasword(t *testing.T) {
	password := "MySecureP@ssw0rd!98273493akdnfhkdhf"
	hashedPassword, err := HashPasswordArgon2id(password)
	if err != nil {
		t.Errorf("Error While Hashing password '%s'", err.Error())
	}

	isValid, err := VerifyPasswordArgon2id(password, hashedPassword)
	if err != nil {
		t.Errorf("Error while verifying password '%s'", err.Error())
	}

	if !isValid {
		t.Errorf("Expected password to be valid, but got invalid")
	}
}
