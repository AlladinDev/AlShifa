// Package validators contains validation functions for clinic module
package validators

import (
	"AlShifa/Clinic/models"
	"fmt"
	"strings"
	"time"
)

const (
	MaxClinicNameLength = 150
	MaxSeasonNameLength = 50
	MaxPincodeDigits    = 6
	MaxOwnerNameLength  = 50
	MaxEmailLength      = 50
	MaxPasswordLength   = 50
	MaxAddressLength    = 200
	MaxGenderLength     = 10
	MaxMobileDigits     = 10
)

type ClinicRegistrationValidationErrors struct {
	ClinicDetailsError map[string]string `json:"clinicDetailsError"`
	OwnerDetailsError  map[string]string `json:"ownerDetailsError"`
}

// func ValidateOwnerDetails(owner *models.Owner, errors *clinicRegistrationValidationErrors) {

// 	if owner == nil {
// 		errors.OwnerDetailsError["owner"] = "owner details are required"
// 		return
// 	}

// 	// Name
// 	name := strings.TrimSpace(owner.Name)
// 	if name == "" {
// 		errors.OwnerDetailsError["name"] = "owner name is required"
// 	} else if len(name) > MaxOwnerNameLength {
// 		errors.OwnerDetailsError["name"] = "owner name is too long"
// 	}

// 	// Email
// 	email := strings.TrimSpace(owner.Email)
// 	if email == "" {
// 		errors.OwnerDetailsError["email"] = "email is required"
// 	} else if len(email) > MaxEmailLength {
// 		errors.OwnerDetailsError["email"] = "email is too long"
// 	} else if !strings.Contains(email, "@") {
// 		errors.OwnerDetailsError["email"] = "invalid email format"
// 	}

// 	// Password
// 	pass := strings.TrimSpace(owner.Password)
// 	if pass == "" {
// 		errors.OwnerDetailsError["password"] = "password is required"
// 	} else if len(pass) < 8 {
// 		errors.OwnerDetailsError["password"] = "password must be at least 8 characters"
// 	} else if len(pass) > MaxPasswordLength {
// 		errors.OwnerDetailsError["password"] = "password is too long"
// 	}

// 	// Address
// 	addr := strings.TrimSpace(owner.Address)
// 	if addr == "" {
// 		errors.OwnerDetailsError["address"] = "owner address is required"
// 	} else if len(addr) > MaxAddressLength {
// 		errors.OwnerDetailsError["address"] = "address is too long"
// 	}

// 	// Gender
// 	gender := strings.TrimSpace(owner.Gender)
// 	if gender == "" {
// 		errors.OwnerDetailsError["gender"] = "gender is required"
// 	} else if len(gender) > MaxGenderLength {
// 		errors.OwnerDetailsError["gender"] = "gender is too long"
// 	} else {
// 		g := strings.ToLower(gender)
// 		if g != "male" && g != "female" && g != "other" {
// 			errors.OwnerDetailsError["gender"] = "gender must be male, female, or other"
// 		}
// 	}

// 	// Mobile
// 	if owner.Mobile <= 0 || lenInt64Digits(owner.Mobile) > MaxMobileDigits {
// 		errors.OwnerDetailsError["mobile"] = "invalid mobile number"
// 	}

// 	// Clinic reference
// 	if owner.Clinic == primitive.NilObjectID {
// 		errors.OwnerDetailsError["clinic"] = "clinic reference is required"
// 	}

// 	if len(errors.OwnerDetailsError) == 0 {
// 		return
// 	}

// }
func ValidateOwnerDetails(owner *models.Owner) map[string]string {
	errors := make(map[string]string)

	if owner == nil {
		errors["owner"] = "owner details are required"
		return errors
	}

	// Name
	name := strings.TrimSpace(owner.Name)
	if name == "" {
		errors["name"] = "owner name is required"
	} else if len(name) > MaxOwnerNameLength {
		errors["name"] = "owner name is too long"
	}

	// Email
	email := strings.TrimSpace(owner.Email)
	if email == "" {
		errors["email"] = "email is required"
	} else if len(email) > MaxEmailLength {
		errors["email"] = "email is too long"
	} else if !strings.Contains(email, "@") {
		errors["email"] = "invalid email format"
	}

	// Password
	pass := strings.TrimSpace(owner.Password)
	if pass == "" {
		errors["password"] = "password is required"
	} else if len(pass) < 8 {
		errors["password"] = "password must be at least 8 characters"
	} else if len(pass) > MaxPasswordLength {
		errors["password"] = "password is too long"
	}

	// Address
	addr := strings.TrimSpace(owner.Address)
	if addr == "" {
		errors["address"] = "owner address is required"
	} else if len(addr) > MaxAddressLength {
		errors["address"] = "address is too long"
	}

	// Gender
	gender := strings.TrimSpace(owner.Gender)
	if gender == "" {
		errors["gender"] = "gender is required"
	} else if len(gender) > MaxGenderLength {
		errors["gender"] = "gender is too long"
	} else {
		g := strings.ToLower(gender)
		if g != "male" && g != "female" && g != "other" {
			errors["gender"] = "gender must be male, female, or other"
		}
	}

	// Mobile
	if owner.Mobile <= 0 || lenInt64Digits(owner.Mobile) > MaxMobileDigits {
		errors["mobile"] = "invalid mobile number"
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

// func ValidateClinicDetails(clinic *models.Clinic, errors *ClinicRegistrationValidationErrors) {

// 	if clinic == nil {
// 		errors.ClinicDetailsError["clinic"] = "clinic details are required"
// 		return
// 	}

// 	// Name
// 	name := strings.TrimSpace(clinic.Name)
// 	if name == "" {
// 		errors.ClinicDetailsError["name"] = "clinic name is required"
// 	} else if len(name) > MaxClinicNameLength {
// 		errors.ClinicDetailsError["name"] = "clinic name is too long"
// 	}

// 	// Address
// 	address := strings.TrimSpace(clinic.Address)
// 	if address == "" {
// 		errors.ClinicDetailsError["address"] = "clinic address is required"
// 	} else if len(address) > MaxAddressLength {
// 		errors.ClinicDetailsError["address"] = "clinic address is too long"
// 	}

// 	// Mobile
// 	if clinic.Mobile <= 0 || lenInt64Digits(clinic.Mobile) > MaxMobileDigits {
// 		errors.ClinicDetailsError["mobile"] = "invalid clinic mobile number"
// 	}

// 	// Pincode
// 	if clinic.Pincode <= 0 || clinic.Pincode < 100000 || clinic.Pincode > 999999 {
// 		errors.ClinicDetailsError["pincode"] = "invalid pincode"
// 	}

// 	// SeasonTiming
// 	if clinic.SeasonTimings == nil {
// 		errors.ClinicDetailsError["seasonTimings"] = "season timing details required"
// 	}

// 	if clinic.SeasonTimings != nil {
// 		for _, seasonDetails := range clinic.SeasonTimings {

// 			if strings.TrimSpace(seasonDetails.Name) == "" {
// 				errors.ClinicDetailsError["seasonTiming.name"] = "season name is required"
// 			} else if len(seasonDetails.Name) > MaxSeasonNameLength {
// 				errors.ClinicDetailsError["seasonTiming.name"] = "season name is too long"
// 			}

// 			if seasonDetails.Start.IsZero() {
// 				errors.ClinicDetailsError["seasonTiming.start"] = "season start time is required"
// 			}
// 			if seasonDetails.End.IsZero() {
// 				errors.ClinicDetailsError["seasonTiming.end"] = "season end time is required"
// 			}
// 			if !seasonDetails.Start.IsZero() && !seasonDetails.End.IsZero() && seasonDetails.End.Before(seasonDetails.Start) {
// 				errors.ClinicDetailsError["seasonTiming.range"] = "season end must be after start"
// 			}
// 			if !seasonDetails.Start.IsZero() && seasonDetails.Start.Before(time.Now().AddDate(-10, 0, 0)) {
// 				errors.ClinicDetailsError["seasonTiming.start"] = "season start date is too old"
// 			}
// 		}
// 	}

// 	if len(errors.ClinicDetailsError) == 0 {
// 		return
// 	}

// }

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
	} else if len(name) > MaxClinicNameLength {
		errors["name"] = "clinic name is too long"
	}

	// Address
	address := strings.TrimSpace(clinic.Address)
	if address == "" {
		errors["address"] = "clinic address is required"
	} else if len(address) > MaxAddressLength {
		errors["address"] = "clinic address is too long"
	}

	// Mobile
	if clinic.Mobile <= 0 || lenInt64Digits(clinic.Mobile) > MaxMobileDigits {
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

	for i, season := range clinic.SeasonTimings {
		prefix := fmt.Sprintf("seasonTimings[%d]", i)

		if strings.TrimSpace(season.Name) == "" {
			errors[prefix+".name"] = "season name is required"
		} else if len(season.Name) > MaxSeasonNameLength {
			errors[prefix+".name"] = "season name is too long"
		}

		if season.Start.IsZero() {
			errors[prefix+".start"] = "season start time is required"
		}

		if season.End.IsZero() {
			errors[prefix+".end"] = "season end time is required"
		}

		if !season.Start.IsZero() && !season.End.IsZero() &&
			season.End.Before(season.Start) {
			errors[prefix+".range"] = "season end must be after start"
		}

		if !season.Start.IsZero() &&
			season.Start.Before(time.Now().AddDate(-10, 0, 0)) {
			errors[prefix+".start"] = "season start date is too old"
		}
	}

	return errors
}

// //now here write single function to handle both validation
// func ValidateClinicRegistration(clinicDetails *models.Clinic, ownerDetails *models.Owner) *clinicRegistrationValidationErrors {
// 	errors := ClinicRegistrationValidationErrors{}
// 	ValidateOwnerDetails(ownerDetails, &errors)
// 	ValidateClinicDetails(clinicDetails, &errors)
// 	return &errors
// }
