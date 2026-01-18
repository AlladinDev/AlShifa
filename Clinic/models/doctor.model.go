package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClinicDetails struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	StartTime   time.Time          `json:"startTiming" bson:"startTiming"`
	EndTime     time.Time          `json:"endTime" bson:"endTime"`
	Clinic      primitive.ObjectID `json:"clinic" bson:"clinic"`
	Information []Clinic           `json:"information" bson:"-"`
	WorkingDays []string           `json:"workingDays" bson:"workingDays"`
}

type Doctor struct {
	RegistrationDate time.Time            `json:"registrationDate" bson:"registrationDate"`
	ID               primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name             string               `json:"name" bson:"name"`
	Qualifications   string               `json:"qualifications" bson:"qualifications"`
	Address          string               `json:"address" bson:"address"`
	Email            string               `json:"email" bson:"email"`
	Password         string               `json:"password" bson:"password"`
	WorkingAt        string               `json:"workingAt" bson:"workingAt"`
	Mobile           int64                `json:"mobile" bson:"mobile"`
	Appointments     []primitive.ObjectID `json:"appointments" bson:"appointments"`
	Clinics          []ClinicDetails      `json:"clinics" bson:"clinics"`
}
