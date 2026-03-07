// Package validators provides validators for clinic module
package validators

import (
	"AlShifa/clinic/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateSlotDetails(input *models.SlotDetails) map[string]string {
	errors := make(map[string]string)

	// Validate clinicID
	clinicID, err := primitive.ObjectIDFromHex(input.ClinicID.Hex())
	if err != nil {
		errors["clinicID"] = "invalid clinic id"
	}
	input.ClinicID = clinicID

	// Validate DoctorID
	doctorID, err := primitive.ObjectIDFromHex(input.DoctorID.Hex())
	if err != nil {
		errors["doctorID"] = "invalid doctor id"
	}

	input.DoctorID = doctorID

	// Validate MaxAppointments
	if input.MaxAppointments <= 0 {
		errors["maxAppointments"] = "must be greater than 0"
	}

	// If validation failed, return errors
	if len(errors) > 0 {
		return errors
	}

	return nil
}
