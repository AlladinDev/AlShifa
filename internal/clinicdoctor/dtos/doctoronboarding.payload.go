package dtos

import (
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
)

type DoctorOnboardingOTPPayload struct {
	OTPHash   string    `json:"otpHash" bson:"otpHash"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	Payload   *models.ClinicDoctorMapping
}
