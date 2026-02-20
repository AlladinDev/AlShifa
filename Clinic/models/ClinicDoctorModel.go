package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorDetails struct {
	RegistrationDate time.Time          `json:"registrationDate,omitempty" bson:"registrationDate"`
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name"`
	Qualifications   string             `json:"qualifications,omitempty" bson:"qualifications"`
	Address          string             `json:"address,omitempty" bson:"address"`
	Email            string             `json:"email,omitempty" bson:"email"`
	Password         string             `json:"password,omitempty" bson:"password"`
	WorkingAt        string             `json:"workingAt,omitempty" bson:"workingAt"`
	Mobile           int64              `json:"mobile,omitempty" bson:"mobile"`
	Role             string             `json:"role,omitempty" bson:"role"`
	StartTime        time.Time          `json:"startTime" bson:"startTime"`
	EndTime          time.Time          `json:"endTime" bson:"endTime"`
	WorkingDays      []string           `json:"workingDays" bson:"workingDays"`
}

// clinicDoctor is a combining model for clinic and doctor relationship it will allow flexibility to add more fields without need to modify clinic or doctor documents
type ClinicDoctor struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	ClinicDetails Clinic             `json:"clinicDetails" bson:"clinicDetails"`
	ClinicID      primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	DoctorID      primitive.ObjectID `json:"doctorID" bson:"doctorID"`
	Doctors       []DoctorDetails    `json:"doctors"  bson:"doctors"`
	StartTime     time.Time          `json:"startTime"`
	EndTime       time.Time          `json:"endTime"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	DoctorName    string             `json:"doctorName" bson:"doctorName"`
	ClinicName    string             `json:"clinicName" bson:"clinicName"`
	ClinicAddress string             `json:"clinicAddress" bson:"clinicAddress"`
}
