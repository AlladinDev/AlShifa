package validators

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/AlladinDev/AlShifa/clinic/models"
	"github.com/AlladinDev/AlShifa/constants"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var mobileRegex = regexp.MustCompile(`^[0-9]{10}$`)

func ValidateDoctor(doc models.Doctor) map[string]string {
	errors := make(map[string]string)

	// Trim all inputs
	name := strings.TrimSpace(doc.Name)
	email := strings.TrimSpace(doc.Email)
	password := strings.TrimSpace(doc.Password)
	mobile := strings.TrimSpace(doc.Mobile)
	address := strings.TrimSpace(doc.Address)
	qualifications := strings.TrimSpace(doc.Qualifications)
	experience := doc.Experience
	post := strings.TrimSpace(doc.Post)

	// 👤 Name
	if name == "" {
		errors["name"] = "Name is required"
	} else if len(name) > constants.MaxNameLength {
		errors["name"] = "Name too long"
	}

	// 📧 Email
	if email == "" {
		errors["email"] = "Email is required"
	} else if len(email) > constants.MaxEmailLength {
		errors["email"] = "Email too long"
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
	}

	// 📱 Mobile
	if mobile == "" {
		errors["mobile"] = "Mobile is required"
	} else if !mobileRegex.MatchString(mobile) {
		errors["mobile"] = "Mobile must be 10 digits"
	}

	// 🔐 Password
	if password == "" {
		errors["password"] = "Password is required"
	} else if len(password) < constants.MinPasswordLength {
		errors["password"] = "Password too short"
	} else if len(password) > constants.MaxPasswordLength {
		errors["password"] = "Password too long"
	} else {
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
			errors["password"] = "Weak password (must include upper, lower, number, special char)"
		}
	}

	// 📍 Address
	if address != "" && len(address) > constants.MaxAddressLength {
		errors["address"] = "Address too long"
	}

	// 🎓 Qualifications
	if qualifications != "" && len(qualifications) > constants.MaxQualificationsLength {
		errors["qualifications"] = "Qualifications too long"
	}

	// 💼 Post
	if post == "" {
		errors["post"] = "Post is required"
	}

	// 🧠 Experience
	if experience == 0 || experience < 0 {
		errors["experience"] = "Experience is required"
	} else if experience > 100 { // e.g. "100 years"
		errors["experience"] = "Experience is too large"
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
