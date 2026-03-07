// Package models provides models for appointment module
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Slot struct {
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	BookingDate time.Time          `json:"bookingDate" bson:"bookingDate"`
	DoctorID    primitive.ObjectID `json:"doctorID" bson:"doctorID"`
	ClinicID    primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	SlotsBooked int                `json:"slotsBooked" bson:"slotsBooked"`
}
