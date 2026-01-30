package validators

import (
	"AlShifa/Clinic/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateAddDoctorToClinicDetails(details *models.AddDoctorToClinic) map[string]string {
	validationErrors := make(map[string]string)
	if details == nil {
		validationErrors["details"] = "missing clinic details"
		return validationErrors
	}

	//check whether doctorID is a valid mongodb id or not
	doctorMongoDBID, err := primitive.ObjectIDFromHex(details.DoctorID.Hex())
	if err != nil {
		validationErrors["doctorID"] = "Invalid doctor id"
	} else {
		details.DoctorID = doctorMongoDBID
	}

	if len(details.WorkingDays) == 0 {
		validationErrors["workingDays"] = "WorkingDays cannot be zero"
	}

	if details.EndTime.Before(details.StartTime) {
		validationErrors["startTime"] = "EndTime cannot be before startTime"
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}
