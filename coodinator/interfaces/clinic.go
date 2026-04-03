// Package interfaces provides interfaces for coordinator service
package interfaces

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicService interface {
	ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string, clinicMaxAppointments int)
	DeductClinicMoneyForAppointment(ctx context.Context, clinicID primitive.ObjectID) error
}
