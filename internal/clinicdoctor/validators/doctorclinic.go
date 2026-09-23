package validators

import (
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
)

func ValidateClinicDoctorMapping(details *models.ClinicDoctorMapping) map[string]string {

	errors := make(map[string]string)

	if details.DoctorID.IsZero() {
		errors["doctorID"] = "doctor id cannot be empty"
	}

	if details.ConsultationFees <= 0 {
		errors["consultationFees"] = "consultation fees must be more than 0"
	}

	if details.AvailableOn == nil {
		errors["availableOn"] = "doctor availability days cannot be nil"
	} else if len(details.AvailableOn) == 0 {
		errors["availableOn"] = "doctor must be available for atleast one day in a week"
	}

	if details.Timings == nil {
		errors["timings"] = "doctor availability timings cannot be nil"
	} else if len(details.Timings) == 0 {
		errors["timings"] = "doctor availability timings cannot be empty"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
