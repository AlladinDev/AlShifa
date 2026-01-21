// Package models provides models for user module
package models

import (
	"AlShifa/Clinic/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	RegistrationDate time.Time            `json:"registrationDate" bson:"registrationDate"`
	Age              int                  `json:"age" bson:"age"`
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	Name             string               `json:"name" bson:"name"`
	Email            string               `json:"email" bson:"email"`
	Role             string               `json:"role" bson:"role"`
	Address          string               `json:"address" bson:"address"`
	Password         string               `json:"password" bson:"password"`
	AppointmentIDS   []primitive.ObjectID `json:"appointmentIDS" bson:"appointmentIDS"`
	Appointments     []models.Appointment `json:"appointments" bson:"appointments"`
	Mobile           int                  `json:"mobile" bson:"mobile"`
	Pincode          int                  `json:"pincode" bson:"pincode"`
}
