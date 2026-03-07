// Package clinicbridge serves as the bridge of communication between doctor and appointment module doctor module uses intrface IDoctorClinic which this bridge implements
package clinicbridge

import (
	"AlShifa/appointment/interfaces"
	"AlShifa/structs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClinicBridge struct {
	clinicService interfaces.IClinicModule
}

func AppointmentClinicBridge(clinicService interfaces.IClinicModule) *ClinicBridge {
	return &ClinicBridge{
		clinicService: clinicService,
	}
}

var _ interfaces.IClinicModule = (*ClinicBridge)(nil)

func (b *ClinicBridge) ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string, clinicMaxAppointment int) {
	return b.clinicService.ClinicDoctorDetails(ctx, clinicID, doctorID, appointmentDateRequested)
}

func (b *ClinicBridge) ClinicExists(ctx context.Context, clinicID primitive.ObjectID) *structs.IAppError {
	return b.clinicService.ClinicExists(ctx, clinicID)
}

func (b *ClinicBridge) DoctorExists(ctx context.Context, doctorID primitive.ObjectID) *structs.IAppError {
	return b.clinicService.DoctorExists(ctx, doctorID)
}

func (b *ClinicBridge) DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) *structs.IAppError {
	return b.clinicService.DoctorClinicMappingExists(ctx, clinicID, doctorID)
}

func (b *ClinicBridge) FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, *structs.IAppError) {
	return b.clinicService.FetchMaxAppointments(ctx, clinicID)
}
