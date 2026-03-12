// Package validation provides validation functions for auth module
package validation

import (
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/AlladinDev/AlShifa/auth/models"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/utils"
)

func ValidateCredentials(c *models.Credientials) map[string]string {
	errors := make(map[string]string)

	// -------------------
	// Email Validation
	// -------------------
	if strings.TrimSpace(c.Email) == "" {
		errors["email"] = "Email is required"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(c.Email) {
			errors["email"] = "Invalid email format"
		}
	}

	// -------------------
	// Password Validation
	// -------------------
	// ---------- Password ----------
	password := c.Password
	if password == "" {
		errors["password"] = "password is required"
	} else if len(password) > 20 {
		errors["password"] = "password is too long"
	} else if len(password) < 8 {
		errors["password"] = "password is too short"
	} else {
		var hasUpper, hasLower, hasDigit, hasSpecial bool

		if len(password) < utils.MinPasswordLength {
			errors["password"] = "password is too short"
		}

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

		if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
			errors["password"] = "password must contain upper, lower, digit and special character"
		}
	}

	// -------------------
	// Role Validation
	// -------------------
	if strings.TrimSpace(c.Role) == "" {
		errors["role"] = "Role is required"
	} else {
		if !slices.Contains(constants.RolesAllowed, c.Role) {
			errors["role"] = "Invalid Role"
		}
	}

	// -------------------
	// Mobile Validation
	// -------------------
	if strings.TrimSpace(c.Mobile) == "" {
		errors["mobile"] = "Mobile number is required"
	} else {
		mobileRegex := regexp.MustCompile(`^[0-9]{10,15}$`)
		if !mobileRegex.MatchString(c.Mobile) {
			errors["mobile"] = "Mobile number must contain 10 to 15 digits"
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
