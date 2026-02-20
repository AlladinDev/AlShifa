package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Appointment struct {
	// 24 bytes each
	AppointmentDate  time.Time `json:"appointmentDate" bson:"appointmentDate"`
	RegistrationDate time.Time `json:"registrationDate" bson:"registrationDate"`

	// 16 bytes each
	PatientName    string `json:"patientName" bson:"patientName"`
	PatientAddress string `json:"patientAddress" bson:"patientAddress"`
	Status         string `json:"status" bson:"status"`
	DoctorName     string `json:"doctorName" bson:"doctorName"`
	ClinicName     string `json:"clinicName" bson:"clinicName"`

	// 12 bytes each
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Clinic primitive.ObjectID `json:"clinic" bson:"clinic"`
	User   primitive.ObjectID `json:"user" bson:"user"`
	Doctor primitive.ObjectID `json:"doctor" bson:"doctor"`

	// 8 bytes
	PatientMobile int `json:"patientMobile" bson:"patientMobile"`

	// 1 byte (placed last to avoid padding waste)
	Slot int `json:"slot" bson:"slot"`
}
