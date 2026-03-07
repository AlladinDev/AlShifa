package models

import (
	"time"
)

//AddDoctorOtpPayload this payload is made in service layer by combining clinicdetails sent with otp info
type AddDoctorOtpPayload struct {
	OTP           string       `json:"otp" bson:"otp"`
	Expiry        time.Time    `json:"expiry" bson:"expiry"`
	ClinicDetails ClinicDoctor `json:"clinicDetails"`
}
