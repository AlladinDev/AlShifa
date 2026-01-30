package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AddDoctorToClinic is the payload which clinic has to send when onboarding a doctor
type AddDoctorToClinic struct {
	ClinicID    primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	DoctorID    primitive.ObjectID `json:"doctorID" bson:"doctorID"`
	WorkingDays []string           `json:"workingDays" bson:"workingDays"`
	StartTime   time.Time          `json:"startTiming" bson:"startTiming"`
	EndTime     time.Time          `json:"endTime" bson:"endTime"`
}
