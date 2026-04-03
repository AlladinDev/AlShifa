package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/AlladinDev/AlShifa/dtos"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmailPassword(details dtos.LoginEmailPasswordDTO) map[string]string {
	errors := make(map[string]string)

	email := strings.TrimSpace(details.Email)
	password := strings.TrimSpace(details.Password)

	// 📧 Email validation
	if email == "" {
		errors["email"] = "Email is required"
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
	}

	// 🔐 Password validation (manual checks instead of regex)
	if password == "" {
		errors["password"] = "Password is required"
	} else {
		if len(password) < 8 {
			errors["password"] = "Password must be at least 8 characters"
			return errors
		}

		var hasUpper, hasLower, hasNumber, hasSpecial bool

		for _, ch := range password {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsLower(ch):
				hasLower = true
			case unicode.IsDigit(ch):
				hasNumber = true
			case strings.ContainsRune("@$!%*?&", ch):
				hasSpecial = true
			}
		}

		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			errors["password"] = "Password must include uppercase, lowercase, number, and special character"
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
