package validators

import (
	models "AlShifa/Users/Models"
	utils "AlShifa/Utils"
	"fmt"
	"regexp"
	"strings"
	"unicode"
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

	// ---------- Email ----------
	email := strings.TrimSpace(strings.ToLower(u.Email))
	if email == "" {
		errors["email"] = "email is required"
	} else if !regexp.MustCompile(utils.EmailRegex).MatchString(email) {
		errors["email"] = "invalid email format"
	}

	// ---------- Password ----------
	password := u.Password
	if password == "" {
		errors["password"] = "password is required"
	} else if len(password) > utils.MaxPasswordLength {
		errors["password"] = "password is too long"
	} else if len(password) < utils.MinPasswordLength {
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

	// ---------- Age ----------
	if u.Age < utils.MinAge || u.Age > utils.MaxAge {
		errors["age"] = fmt.Sprintf("age must be between  %d and %d", utils.MinAge, utils.MaxAge)
	}

	// ---------- Address ----------
	address := strings.TrimSpace(u.Address)
	if address == "" {
		errors["address"] = "address is required"
	} else if len(address) < utils.MinAddressLength {
		errors["address"] = "address is too short"
	}

	// ---------- Mobile ----------
	mobile := fmt.Sprintf("%d", u.Mobile)
	if u.Mobile < 0 {
		errors["mobile"] = "mobile number must be positive"
	} else if u.Mobile == 0 {
		errors["mobile"] = "mobile number is required"
	} else if len(mobile) != utils.MobileLength {
		errors["mobile"] = "mobile number must be 10 digits"
	}

	// ---------- Pincode ----------
	pincode := fmt.Sprintf("%d", u.Pincode)
	if u.Pincode == 0 {
		errors["pincode"] = "pincode is required"
	} else if len(pincode) != utils.PincodeLength {
		errors["pincode"] = "pincode must be 6 digits"
	}

	// ---------- Final ----------
	if len(errors) == 0 {
		return nil
	}

	return errors
}
