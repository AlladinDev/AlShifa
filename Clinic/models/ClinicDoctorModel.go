package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorDetails struct {
	RegistrationDate time.Time          `json:"registrationDate,omitempty" bson:"registrationDate"`
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name"`
	Qualifications   string             `json:"qualifications,omitempty" bson:"qualifications"`
	Address          string             `json:"address,omitempty" bson:"address"`
	Email            string             `json:"email,omitempty" bson:"email"`
	Password         string             `json:"password,omitempty" bson:"password"`
	WorkingAt        string             `json:"workingAt,omitempty" bson:"workingAt"`
	Mobile           int64              `json:"mobile,omitempty" bson:"mobile"`
	Role             string             `json:"role,omitempty" bson:"role"`
	StartTime        time.Time          `json:"startTime" bson:"startTime"`
	EndTime          time.Time          `json:"endTime" bson:"endTime"`
	WorkingDays      []string           `json:"workingDays" bson:"workingDays"`
}

// clinicDoctor is a combining model for clinic and doctor relationship it will allow flexibility to add more fields without need to modify clinic or doctor documents
//using aggregation pipeline we can inject clinic details and doctor details into this collection and it will be easy to query for doctors of a clinic or clinics of a doctor without need to do multiple queries
//repository will just send clinicDetails and doctorDetails as per this model
type ClinicDoctor struct {
	ClinicDetails Clinic          `json:"clinicDetails" bson:"clinicDetails"`
	Doctors       []DoctorDetails `json:"doctors"  bson:"doctors"`
}
