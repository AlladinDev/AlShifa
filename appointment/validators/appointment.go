// Package validators provides validation functions for appointment module
package validators

import (
	"regexp"
	"strings"
	"time"

	"github.com/AlladinDev/AlShifa/appointment/models"
	"github.com/AlladinDev/AlShifa/constants"
)

// Regex for 10-digit numeric mobile
var mobileRegex = regexp.MustCompile(`^[0-9]{10}$`)

func ValidateAppointment(app models.Appointment) map[string]string {
	errors := make(map[string]string)

	// Trim string fields
	patientName := strings.TrimSpace(app.PatientName)
	patientAddress := strings.TrimSpace(app.PatientAddress)
	patientMobile := strings.TrimSpace(app.PatientMobile)

	// 👤 Patient Name
	if patientName == "" {
		errors["patientName"] = "Patient name is required"
	} else if len(patientName) > constants.MaxNameLength {
		errors["patientName"] = "Patient name too long"
	}

	// 🏠 Patient Address
	if patientAddress == "" {
		errors["patientAddress"] = "Patient address is required"
	} else if patientAddress != "" && len(patientAddress) > constants.MaxAddressLength {
		errors["patientAddress"] = "Patient address too long"
	}

	// 📱 Patient Mobile (regex)
	if patientMobile == "" {
		errors["patientMobile"] = "Patient mobile is required"
	} else if !mobileRegex.MatchString(patientMobile) {
		errors["patientMobile"] = "Patient mobile must be exactly 10 digits and numeric"
	}

	// 🆔 ObjectIDs
	if app.ClinicID.IsZero() {
		errors["clinicID"] = "Clinic ID is required"
	}
	if app.DoctorID.IsZero() {
		errors["doctorID"] = "Doctor ID is required"
	}

	// 📅 AppointmentDate
	if app.AppointmentDate.IsZero() {
		errors["appointmentDate"] = "Appointment date is required"
	} else if app.AppointmentDate.Before(time.Now()) {
		errors["appointmentDate"] = "Appointment date cannot be in past"
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
