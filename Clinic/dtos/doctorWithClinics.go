package dtos

import (
	"time"

	"github.com/AlladinDev/AlShifa/clinic/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clinicDTO struct {
	ID              primitive.ObjectID           `json:"_id" bson:"id"`
	Name            string                       `json:"name,omitempty"  bson:"name"`                  // 16 bytes
	Address         string                       `json:"address,omitempty" bson:"address"`             // 16 bytes
	SeasonTimings   []models.SeasonTimingDetails `json:"seasonTimings,omitempty" bson:"seasonTimings"` // 8 bytes (pointer)
	Mobile          string                       `json:"mobile,omitempty" bson:"mobile"`               // 8 bytes (int64 for phone numbers)
	Pincode         int32                        `json:"pincode,omitempty" bson:"pincode"`
	MaxAppointments int                          `json:"maxAppointments" bson:"maxAppointments"`
	WorkingDays     []string                     `json:"workingDays" bson:"workingDays"`
	StartTime       time.Time                    `json:"startTiming" bson:"startTiming"`
	EndTime         time.Time                    `json:"endTime" bson:"endTime"`
}
type DoctorWithClinics struct {
	ID            primitive.ObjectID `json:"_id" bson:"id"`
	DoctorDetails *models.Doctor     `json:"doctorDetails" bson:"doctorDetails"`
	Clinics       []clinicDTO        `json:"clinics" bson:"clinics"`
}
