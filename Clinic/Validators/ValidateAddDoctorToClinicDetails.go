package validators

import (
	"github.com/AlladinDev/AlShifa/clinic/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateAddDoctorToclinicDetails(details *models.ClinicDoctor) map[string]string {
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

	//validate clinicid whether this is a valid clinicid or not
	clinicMongodbID, clinicIDErr := primitive.ObjectIDFromHex(details.ClinicID.Hex())
	if clinicIDErr != nil {
		validationErrors["clinicID"] = "Invalid clinicID"
	} else {
		details.ClinicID = clinicMongodbID
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
