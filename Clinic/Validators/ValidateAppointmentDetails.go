package validators

import (
	"AlShifa/clinic/models"
	middleware "AlShifa/middleware"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidateAppointmentDetails validates appointment fields safely.
// Converts string IDs to ObjectID and date string to time.Time if valid.
func ValidateAppointmentDetails(appointmentDetails *models.Appointment, ctx context.Context) map[string]string {
	errors := make(map[string]string)

	// --------------------------
	// 1. Extract user ID from context (string -> ObjectID)
	// --------------------------
	userIDVal := ctx.Value(middleware.ContextUserIDKey)

	if userIDVal == nil {
		errors["userId"] = "User ID is missing from context"
	} else {
		userIDStr, ok := userIDVal.(string)
		if !ok || userIDStr == "" {
			errors["userId"] = "User ID is invalid"
		} else {
			id, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				errors["userId"] = "User ID is not a valid MongoDB ObjectID"
			} else {
				appointmentDetails.User = id
			}
		}
	}

	// --------------------------
	// 2. Validate AppointmentDate (string -> time.Time)
	// --------------------------

	if appointmentDetails.AppointmentDate.Before(time.Now()) {
		errors["appointmentDate"] = "Appointment Date cannot be in past"
	}

	// --------------------------
	// 3. Validate PatientName
	// --------------------------
	if appointmentDetails.PatientName == "" {
		errors["patientName"] = "Patient name is required"
	} else if len(appointmentDetails.PatientName) > 100 {
		errors["patientName"] = "Patient name is too long (max 100 characters)"
	}

	// --------------------------
	// 4. Validate PatientAddress
	// --------------------------
	if appointmentDetails.PatientAddress == "" {
		errors["patientAddress"] = "Patient address is required"
	} else if len(appointmentDetails.PatientAddress) > 200 {
		errors["patientAddress"] = "Patient address is too long (max 200 characters)"
	}

	// --------------------------
	// 5. Validate PatientMobile
	// --------------------------
	if appointmentDetails.PatientMobile <= 0 {
		errors["patientMobile"] = "Patient mobile must be a positive number"
	} else if appointmentDetails.PatientMobile < 1000000000 || appointmentDetails.PatientMobile > 9999999999 {
		errors["patientMobile"] = "Patient mobile must be a valid 10-digit number"
	}

	// --------------------------
	// 6. Validate clinic ID (string -> ObjectID)
	// --------------------------

	id, err := primitive.ObjectIDFromHex(appointmentDetails.Clinic.Hex())
	if err != nil {
		errors["clinic"] = "clinic ID is not a valid MongoDB ObjectID"
	}
	if appointmentDetails.Clinic == primitive.NilObjectID {
		errors["clinic"] = "clinic ID is not a valid MongoDB ObjectID"
	} else {
		appointmentDetails.Clinic = id
	}

	// --------------------------
	// 7. Validate Doctor ID (string -> ObjectID)
	// --------------------------
	doctorID, err := primitive.ObjectIDFromHex(appointmentDetails.Doctor.Hex())
	if err != nil {
		errors["doctor"] = "Doctor ID is not a valid MongoDB ObjectID"
	}
	if appointmentDetails.Doctor == primitive.NilObjectID {
		errors["doctor"] = "Doctor ID is not a valid MongoDB ObjectID"
	} else {
		appointmentDetails.Doctor = doctorID
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
