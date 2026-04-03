// Package validators provides validators for coordinator service
package validators

import (
	"strings"
	"time"

	sharedModels "github.com/AlladinDev/AlShifa/models"
)

// ValidateAppointmentDetails validates critical fields of Appointment and returns map[field]error
func ValidateAppointmentDetails(details sharedModels.Appointment) map[string]string {
	errs := make(map[string]string)

	// Validate AppointmentDate
	if details.AppointmentDate.IsZero() {
		errs["AppointmentDate"] = "appointment date is required"
	} else if details.AppointmentDate.Before(time.Now()) {
		errs["AppointmentDate"] = "appointment date must be in the future"
	}

	// Validate PatientName
	if strings.TrimSpace(details.PatientName) == "" {
		errs["PatientName"] = "patient name is required"
	}
	if len(details.PatientName) > 50 {
		errs["PatientName"] = "patient name cannot be more than 50"
	}
	if len(details.PatientName) < 3 {
		errs["PatientName"] = "patient name must be more than 3 chars"
	}

	// Validate PatientAddress
	if strings.TrimSpace(details.PatientAddress) == "" {
		errs["PatientAddress"] = "patient address is required"
	}

	if len(details.PatientAddress) > 50 {
		errs["address"] = "patient address cannot be more than 50"
	}
	if len(details.PatientAddress) < 3 {
		errs["PatientName"] = "patient name must be more than 3 chars"
	}

	// Validate ClinicID
	if details.ClinicID.IsZero() {
		errs["ClinicID"] = "invalid clinic ID"
	}

	// Validate DoctorID
	if details.DoctorID.IsZero() {
		errs["DoctorID"] = "invalid doctor ID"
	}

	// Validate PatientMobile
	if details.PatientMobile <= 0 {
		errs["PatientMobile"] = "patient mobile must be details positive number"
	} else if lenInt(details.PatientMobile) != 10 {
		errs["PatientMobile"] = "patient mobile must be 10 digits"
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Helper function to get the number of digits in an int
func lenInt(n int) int {
	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}
