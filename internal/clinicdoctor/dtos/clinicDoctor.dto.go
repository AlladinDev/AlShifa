package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Doctor struct {
	MappingID        primitive.ObjectID    `json:"mappingID" bson:"mappingID"`
	ID               primitive.ObjectID    `json:"_id" bson:"_id"`
	AvailableOn      []string              `json:"availableOn" bson:"availableOn"`
	ConsultationFees int                   `json:"consultationFees" bson:"consultationFees"`
	Timings          []ConsultationTimings `json:"timings" bson:"timings"`
	Name             string                `json:"name" bson:"name"`
	Experience       int                   `json:"experience" bson:"experience"`
	Speciality       string                `json:"speciality" bson:"speciality"`
	JoinedAt         time.Time             `json:"joinedOn" bson:"joinedOn"`
	PhotoURL         string                `json:"photoUrl" bson:"photoUrl"`
	WorkingAt        string                `json:"workingAt" bson:"workingAt"`
	Post             string                `json:"post" bson:"post"`
	Qualifications   string                `json:"qualifications" bson:"qualifications"`
}
type ClinicDoctorDTO struct {
	ID              primitive.ObjectID    `json:"_id" bson:"_id"`
	ClinicCreatedAt time.Time             `json:"registrationDate" bson:"registrationDate"`
	Name            string                `json:"name" bson:"name"`
	Departments     []string              `json:"departments" bson:"departments"`
	Doctors         []Doctor              `json:"doctors" bson:"doctors"`
	WorkingDays     []string              `json:"workingDays" bson:"workingDays"`
	Address         string                `json:"address" bson:"address"`
	SeasonTimings   []SeasonTimingDetails `json:"seasonTimings" bson:"seasonTimings"`
}
