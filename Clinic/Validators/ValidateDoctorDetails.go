package validators

import (
	"AlShifa/Clinic/models"
	"net/mail"
	"strings"
	"unicode/utf8"
)

// --------------------
// Validation constants
// --------------------
const (
	NameMinLen          = 2
	NameMaxLen          = 100
	QualificationMinLen = 2
	QualificationMaxLen = 200
	AddressMinLen       = 5
	AddressMaxLen       = 300
	EmailMaxLen         = 100
	PasswordMinLen      = 6
	PasswordMaxLen      = 100
	WorkingAtMinLen     = 2
	WorkingAtMaxLen     = 200
	MobileMin           = 6000000000 // example: 10-digit starting with 6-9
	MobileMax           = 9999999999
)

// --------------------
// Validation function
// --------------------

// ValidateDoctor validates a Doctor struct and returns a map of field errors
func ValidateDoctor(d models.Doctor) map[string]string {
	errors := make(map[string]string)

	// ---------- Name ----------
	name := strings.TrimSpace(d.Name)
	if name == "" {
		errors["name"] = "Name is required"
	} else {
		length := utf8.RuneCountInString(name)
		if length < NameMinLen {
			errors["name"] = "Name must be at least " + itoa(NameMinLen) + " characters"
		}
		if length > NameMaxLen {
			errors["name"] = "Name cannot exceed " + itoa(NameMaxLen) + " characters"
		}
	}

	// ---------- Qualifications ----------
	qual := strings.TrimSpace(d.Qualifications)
	if qual == "" {
		errors["qualifications"] = "Qualifications are required"
	} else {
		length := utf8.RuneCountInString(qual)
		if length < QualificationMinLen {
			errors["qualifications"] = "Qualifications must be at least " + itoa(QualificationMinLen) + " characters"
		}
		if length > QualificationMaxLen {
			errors["qualifications"] = "Qualifications cannot exceed " + itoa(QualificationMaxLen) + " characters"
		}
	}

	// ---------- Address ----------
	addr := strings.TrimSpace(d.Address)
	if addr == "" {
		errors["address"] = "Address is required"
	} else {
		length := utf8.RuneCountInString(addr)
		if length < AddressMinLen {
			errors["address"] = "Address must be at least " + itoa(AddressMinLen) + " characters"
		}
		if length > AddressMaxLen {
			errors["address"] = "Address cannot exceed " + itoa(AddressMaxLen) + " characters"
		}
	}

	// ---------- Email ----------
	email := strings.TrimSpace(d.Email)
	if email == "" {
		errors["email"] = "Email is required"
	} else {
		if utf8.RuneCountInString(email) > EmailMaxLen {
			errors["email"] = "Email cannot exceed " + itoa(EmailMaxLen) + " characters"
		}
		if _, err := mail.ParseAddress(email); err != nil {
			errors["email"] = "Invalid email format"
		}
	}

	// ---------- Password ----------
	pass := strings.TrimSpace(d.Password)
	if pass == "" {
		errors["password"] = "Password is required"
	} else {
		length := utf8.RuneCountInString(pass)
		if length < PasswordMinLen {
			errors["password"] = "Password must be at least " + itoa(PasswordMinLen) + " characters"
		}
		if length > PasswordMaxLen {
			errors["password"] = "Password cannot exceed " + itoa(PasswordMaxLen) + " characters"
		}
	}

	// ---------- WorkingAt ----------
	work := strings.TrimSpace(d.WorkingAt)
	if work == "" {
		errors["workingAt"] = "WorkingAt is required"
	} else {
		length := utf8.RuneCountInString(work)
		if length < WorkingAtMinLen {
			errors["workingAt"] = "WorkingAt must be at least " + itoa(WorkingAtMinLen) + " characters"
		}
		if length > WorkingAtMaxLen {
			errors["workingAt"] = "WorkingAt cannot exceed " + itoa(WorkingAtMaxLen) + " characters"
		}
	}

	// ---------- Mobile ----------
	mobile := d.Mobile
	if mobile == 0 {
		errors["mobile"] = "Mobile number is required"
	} else {
		if mobile < MobileMin || mobile > MobileMax {
			errors["mobile"] = "Mobile number must be a valid 10-digit number"
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// --------------------
// Helper function
// --------------------
func itoa(i int) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(rune(i))), " "))
}
