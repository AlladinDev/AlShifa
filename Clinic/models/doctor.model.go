package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clinicDetails struct {
	StartTime   time.Time          `json:"startTiming,omitempty" bson:"startTiming"`
	EndTime     time.Time          `json:"endTime,omitempty" bson:"endTime"`
	Clinic      primitive.ObjectID `json:"clinic,omitempty" bson:"clinic"`
	Information *Clinic            `json:"information,omitempty" bson:"information"`
	WorkingDays []string           `json:"workingDays,omitempty" bson:"workingDays"`
}

type Doctor struct {
	RegistrationDate time.Time          `json:"registrationDate,omitempty" bson:"registrationDate"`
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name"`
	Qualifications   string             `json:"qualifications,omitempty" bson:"qualifications"`
	Address          string             `json:"address,omitempty" bson:"address"`
	Email            string             `json:"email,omitempty" bson:"email"`
	Password         string             `json:"password,omitempty" bson:"password"`
	WorkingAt        string             `json:"workingAt,omitempty" bson:"workingAt"`
	Mobile           int64              `json:"mobile,omitempty" bson:"mobile"`
	//clinics is for injecting clinicdetails into it during mongodb pipeline
	Clinics []clinicDetails `json:"clinics,omitempty" bson:"clinics"`
	Role    string          `json:"role,omitempty" bson:"role"`
}

//DoctorPublicDetails details which will be sent to public
type DoctorPublicDetails struct {
	ID             primitive.ObjectID `json:"id,omitzero" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty" bson:"name"`
	Qualifications string             `json:"qualifications,omitempty" bson:"qualifications"`
	WorkingAt      string             `json:"workingAt,omitempty" bson:"workingAt"`
	Clinics        []clinicDetails    `json:"clinics,omitempty" bson:"clinics"`
}
