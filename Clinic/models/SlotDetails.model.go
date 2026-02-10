package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SlotDetails struct {
	ClinicID        primitive.ObjectID `json:"clinicID"`
	DoctorID        primitive.ObjectID `json:"doctorID"`
	MaxAppointments int                `json:"maxAppointments"`
}
