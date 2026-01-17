package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	AppointmentDate  time.Time          `json:"appointmentDate" bson:"appointmentDate"`
	RegistrationDate time.Time          `json:"registrationDate" bson:"registrationDate"`
	Status           string             `json:"status" bson:"status"`
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Clinic           primitive.ObjectID `json:"clinic" bson:"clinic"`
	User             primitive.ObjectID `json:"user" bson:"user"`
	Doctor           primitive.ObjectID `json:"doctor" bson:"doctor"`
	Slot             int8               `json:"slot" bson:"slot"`
}
