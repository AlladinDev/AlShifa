package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	// 24 bytes each
	AppointmentDate time.Time `json:"appointmentDate" bson:"appointmentDate"`
	CreatedAt       time.Time `json:"registrationDate" bson:"registrationDate"`

	// 16 bytes each
	PatientName    string `json:"patientName" bson:"patientName"`
	PatientAddress string `json:"patientAddress" bson:"patientAddress"`
	Status         string `json:"status" bson:"status"`
	DoctorName     string `json:"doctorName" bson:"doctorName"`
	ClinicName     string `json:"clinicName" bson:"clinicName"`
	ClinicAddress  string `json:"clinicAddress" bson:"clinicAddress"`

	// 12 bytes each
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	ClinicID primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	DoctorID primitive.ObjectID `json:"doctorID" bson:"doctorID"`

	// 8 bytes
	PatientMobile string `json:"patientMobile" bson:"patientMobile"`

	// 1 byte (placed last to avoid padding waste)
	Slot int `json:"slot" bson:"slot"`
}
