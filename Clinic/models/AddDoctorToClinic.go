// Package models provides models for clinic
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ClinicDoctor is the payload which clinic has to send when onboarding a doctor this is a mapper between clinic and doctor
type ClinicDoctor struct {
	ClinicID      primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	DoctorID      primitive.ObjectID `json:"doctorID" bson:"doctorID"`
	WorkingDays   []string           `json:"workingDays" bson:"workingDays"`
	StartTime     time.Time          `json:"startTiming" bson:"startTiming"`
	EndTime       time.Time          `json:"endTime" bson:"endTime"`
	DoctorName    string             `json:"doctorName" bson:"doctorName"`
	Experience    int                `json:"experience"  bson:"experience"`
	ClinicDetails *Clinic            `json:"clinicDetails" bson:"clinicDetails"`
	WorkingAs     string             `json:"workingAs" bson:"workingAs"`
	ClinicName    string             `json:"clinicName" bson:"clinicName"`
	ClinicAddress string             `json:"clinicAddress" bson:"clinicAddress"`
}

/*
   {doctorDetails:{
   name,
   address
   qualifications
   etc
   },

   clinics:[
   {
      workingDays
      start end
      clinic timing
      clinic details
   }
   ]

}

*/
