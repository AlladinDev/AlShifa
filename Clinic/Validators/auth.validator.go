// Package validators contains validation functions for clinic module
package validators

import (
	"AlShifa/Clinic/models"
	utils "AlShifa/Utils"
	"regexp"
	"strings"
	"unicode"
)

func ValidateOwnerDetails(owner *models.Owner) map[string]string {
	errors := make(map[string]string)

	if owner == nil {
		errors["owner"] = "owner details are required"
		return errors
	}

	// Name
	name := strings.TrimSpace(owner.Name)
	if name == "" {
		errors["name"] = utils.NameMissingErrMsg
	} else if len(name) > utils.MaxNameLength {
		errors["name"] = utils.LongNameErrMsg
	} else if len(name) < utils.MinNameLength {
		errors["name"] = utils.ShortNameErrMsg
	}

	// Email
	email := strings.TrimSpace(owner.Email)
	if email == "" {
		errors["email"] = utils.EmailMissingErrMsg
	} else if len(email) > utils.MaxEmailLength {
		errors["email"] = utils.LongEmailErrMsg
	} else if len(email) < utils.MinEmailLength {
		errors["email"] = utils.ShortEmailErrMsg
	} else if !regexp.MustCompile(utils.EmailRegex).MatchString(email) {
		errors["email"] = utils.InvalidEmailFormatMsg
	}

	// Password
	pass := strings.TrimSpace(owner.Password)
	if pass == "" {
		errors["password"] = utils.PasswordMissingErrMsg
	} else if len(pass) < 8 {
		errors["password"] = utils.ShortPasswordErrMsg
	} else if len(pass) > utils.MaxPasswordLength {
		errors["password"] = utils.LongPasswordErrMsg
	} else {
		var hasUpper bool
		var hasLower bool
		var hasDigit bool
		var hasSpecial bool
		for _, r := range pass {
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
			errors["password"] = utils.PasswordWeakErrMsg
		}

	}

	// Address
	addr := strings.TrimSpace(owner.Address)
	if addr == "" {
		errors["address"] = utils.AddressMissingErrMsg
	} else if len(addr) > utils.MaxAddressLength {
		errors["address"] = utils.LongAddressErrMsg
	} else if len(addr) < utils.MinAddressLength {
		errors["address"] = utils.ShortAddressErrMsg
	}

	// Gender
	gender := strings.TrimSpace(owner.Gender)
	if gender == "" {
		errors["gender"] = utils.GenderMissingErrMsg
	} else {
		g := strings.ToLower(gender)
		if g != "male" && g != "female" && g != "other" {
			errors["gender"] = utils.InvalidGenderErrMsg
		}
	}

	// Mobile

	if lenInt64Digits(owner.Mobile) != utils.MobileLength {
		errors["mobile"] = utils.InvalidMobileNumberMsg
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// helper function reused from clinic validator
func lenInt64Digits(n int64) int {
	count := 0
	if n == 0 {
		return 1
	}
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func ValidateClinicDetails(clinic *models.Clinic) map[string]string {
	errors := make(map[string]string)

	if clinic == nil {
		errors["clinic"] = "clinic details are required"
		return errors
	}

	// Name
	name := strings.TrimSpace(clinic.Name)
	if name == "" {
		errors["name"] = "clinic name is required"
	} else if len(name) > utils.MaxNameLength {
		errors["name"] = "clinic name is too long"
	}

	// Address
	address := strings.TrimSpace(clinic.Address)
	if address == "" {
		errors["address"] = "clinic address is required"
	} else if len(address) > utils.MaxAddressLength {
		errors["address"] = "clinic address is too long"
	}

	// Mobile
	if clinic.Mobile <= 0 || lenInt64Digits(clinic.Mobile) != utils.MobileLength {
		errors["mobile"] = "invalid clinic mobile number"
	}

	// Pincode
	if clinic.Pincode < 100000 || clinic.Pincode > 999999 {
		errors["pincode"] = "invalid pincode"
	}

	// Season timings
	if len(clinic.SeasonTimings) == 0 {
		errors["seasonTimings"] = "season timing details required"
		return errors
	}

	//fmt.Print("season details passed are", clinic.SeasonTimings)
	for _, season := range clinic.SeasonTimings {
		prefix := season.Name

		if strings.TrimSpace(season.Name) == "" {
			errors[prefix+" name"] = " season name is required"
		} else if len(season.Name) > utils.MaxNameLength {
			errors[prefix+" name"] = " season name is too long"
		}

		if season.Start.IsZero() {
			errors[prefix+" start"] = " season start time is required"
		}

		if season.End.IsZero() {
			errors[prefix+" end"] = " season end time is required"
		}

		if !season.Start.IsZero() && !season.End.IsZero() &&
			season.End.Before(season.Start) {
			errors[prefix+" range"] = " season end must be after start"
		}

		// if !season.Start.IsZero() &&
		// 	season.Start.Before(time.Now().AddDate(-10, 0, 0)) {
		// 	errors[prefix+" start "] = " season start date is too old"
		// }
	}

	return errors
}
