// Package validators contains validation functions for owner module
package validators

import (
	"strings"

	"github.com/AlladinDev/AlShifa/internal/owner/models"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
)

// MobileLength is 10 for standard 10-digit numbers
const MobileLength = 10

var (
	// Regex to match exactly 10 digits
	InvalidMobileNumberMsg  = "Invalid mobile number. Must be 10 digits."
	InvalidAadhaarNumberMsg = "Invalid aadhaar number. Must be 12 digits."
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
	if len(owner.Mobile) != 10 {
		errors["mobile"] = InvalidMobileNumberMsg
	}

	//now validate aadhaar number
	if len(owner.AadhaarNumber) != 12 {
		errors["aadhaarNumber"] = InvalidAadhaarNumberMsg
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
