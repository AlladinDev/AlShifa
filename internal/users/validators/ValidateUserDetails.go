// Package validators provides validation functions for user package
package validators

import (
	"strings"
	"unicode"

	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
	models "github.com/AlladinDev/AlShifa/internal/users/models"
)

func ValidateUser(u *models.User) map[string]string {
	if u == nil {
		return map[string]string{
			"user": "user cannot be nil",
		}
	}

	errors := make(map[string]string)

	// ---------- Name ----------
	name := strings.TrimSpace(u.Name)
	if name == "" {
		errors["name"] = "name is required"
	} else if len(name) < utils.MinNameLength || len(name) > utils.MaxNameLength {
		errors["name"] = "name length is invalid"
	} else {
		for _, r := range name {
			if !unicode.IsLetter(r) && r != ' ' {
				errors["name"] = "name must contain only letters and spaces"
				break
			}
		}
	}

	// ---------- Final ----------
	if len(errors) == 0 {
		return nil
	}

	return errors
}
