package validators

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	doctormodels "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func ValidateDoctorDetails(details *doctormodels.Doctor) map[string]string {
	validationErrors := map[string]string{}

	if details == nil {
		validationErrors["doctorDetails"] = "details cannot be empty"
		return validationErrors
	}

	// Name validation
	nameLength := utf8.RuneCountInString(details.Name)

	if nameLength > constants.MaxNameLength {
		validationErrors["name"] = fmt.Sprintf(
			"name cannot be more than %d characters",
			constants.MaxNameLength,
		)
	} else if nameLength < constants.MinNameLength {
		validationErrors["name"] = fmt.Sprintf(
			"name cannot be less than %d characters",
			constants.MinNameLength,
		)
	}

	// Address validation
	addressLength := utf8.RuneCountInString(details.Address)

	if addressLength > constants.MaxAddressLength {
		validationErrors["address"] = fmt.Sprintf(
			"address cannot be more than %d characters",
			constants.MaxAddressLength,
		)
	} else if addressLength < constants.MinAddressLength {
		validationErrors["address"] = fmt.Sprintf(
			"address cannot be less than %d characters",
			constants.MinAddressLength,
		)
	}

	// Qualifications validation
	qualificationsLength := utf8.RuneCountInString(details.Qualifications)

	if qualificationsLength > constants.MaxQualificationsLength {
		validationErrors["qualifications"] = fmt.Sprintf(
			"qualifications cannot be more than %d characters",
			constants.MaxQualificationsLength,
		)
	} else if qualificationsLength < constants.MinQualificationsLength {
		validationErrors["qualifications"] = fmt.Sprintf(
			"qualifications cannot be less than %d characters",
			constants.MinQualificationsLength,
		)
	}

	// Experience validation
	if details.Experience < 0 {
		validationErrors["experience"] = "experience cannot be less than 0"
	}

	// Mobile validation
	if details.Mobile < 5000000000 || details.Mobile > 9999999999 {
		validationErrors["mobile"] = "valid 10-digit mobile number is required"
	}

	// Email validation
	if !emailRegex.MatchString(details.Email) {
		validationErrors["email"] = "valid email id is required"
	}

	//field validation
	if len(details.Field) < constants.MinSpecializationFieldLength {
		validationErrors["field"] = fmt.Sprintf(
			"field of specialization cannot be less than %d characters",
			constants.MinSpecializationFieldLength,
		)
	} else if len(details.Field) > constants.MaxSpecializationFieldLength {
		validationErrors["field"] = fmt.Sprintf(
			"field of specialization cannot be more than %d characters",
			constants.MaxSpecializationFieldLength,
		)
	}

	// Password validation
	passwordLength := utf8.RuneCountInString(details.Password)

	if passwordLength < 8 {
		validationErrors["password"] = "password should be at least 8 characters long"
	} else if passwordLength > 20 {
		validationErrors["password"] = "password cannot be more than 20 characters"
	} else {
		containsNumeric := false
		containsUpperCase := false
		containsLowerCase := false
		containsSpecialCase := false

		for _, char := range details.Password {
			switch {
			case unicode.IsDigit(char):
				containsNumeric = true

			case unicode.IsUpper(char):
				containsUpperCase = true

			case unicode.IsLower(char):
				containsLowerCase = true

			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				containsSpecialCase = true
			}
		}

		if !containsNumeric {
			validationErrors["password"] = "password must contain at least one number"
		} else if !containsUpperCase {
			validationErrors["password"] = "password must contain at least one uppercase letter"
		} else if !containsLowerCase {
			validationErrors["password"] = "password must contain at least one lowercase letter"
		} else if !containsSpecialCase {
			validationErrors["password"] = "password must contain at least one special character"
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
