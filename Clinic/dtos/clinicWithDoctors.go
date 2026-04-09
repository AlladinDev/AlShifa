// Package dtos provides dtos for clinic module
package dtos

import (
	"time"

	"github.com/AlladinDev/AlShifa/clinic/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type doctorDTO struct {
	ID             primitive.ObjectID `json:"_id" bson:"id"`
	Qualifications string             `json:"qualifications" bson:"qualifications"`
	Name           string             `json:"name" bson:"name"`
	Experience     int8               `json:"experience" bson:"experience"`
	Post           string             `json:"post" bson:"post"`
	WorkingAt      string             `json:"workingAt,omitempty" bson:"workingAt"`
	Mobile         string             `json:"mobile" bson:"mobile"`
	WorkingDays    []string           `json:"workingDays" bson:"workingDays"`
	StartTime      time.Time          `json:"startTiming" bson:"startTiming"`
	EndTime        time.Time          `json:"endTime" bson:"endTime"`
}

type ClinicWithDoctors struct {
	ID            primitive.ObjectID `json:"_id" bson:"id"`
	ClinicDetails *models.Clinic     `json:"clinicDetails" bson:"clinicDetails"`
	Doctors       []doctorDTO        `json:"doctors" bson:"doctors"`
}
