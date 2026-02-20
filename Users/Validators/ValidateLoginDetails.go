// Package validators provides validation functions for user module
package validators

import (
	structs "AlShifa/users/structs"
	utils "AlShifa/utils"
	"unicode"
)

func ValidateLoginDetails(l *structs.LoginDetails) map[string]string {
	if l == nil {
		return map[string]string{
			"login": "login details cannot be nil",
		}
	}

	errors := make(map[string]string)

	//email validation
	if l.Email == "" {
		errors["email"] = "Valid email address required"
	}

	// ---------- Password ----------
	password := l.Password
	if password == "" {
		errors["password"] = "password is required"
	} else if len(password) < utils.MinPasswordLength {
		errors["password"] = "password is too short"
	} else if len(password) > utils.MaxPasswordLength {
		errors["password"] = "password is too long"
	} else {
		var hasUpper, hasLower, hasDigit, hasSpecial bool

		for _, r := range password {
			switch {
			case unicode.IsUpper(r):
				hasUpper = true
			case unicode.IsLower(r):
				hasLower = true
			case unicode.IsDigit(r):
				hasDigit = true
			case unicode.IsPunct(r) || unicode.IsSymbol(r):
				hasSpecial = true
			}
		}

		if !(hasUpper && hasLower && hasDigit && hasSpecial) {
			errors["password"] = "password must contain upper, lower, digit and special character"
		}
	}

	// ---------- Final ----------
	if len(errors) == 0 {
		return nil
	}

	return errors
}
