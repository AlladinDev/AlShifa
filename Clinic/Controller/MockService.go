package controller

import (
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockService implements IService for testing purposes
type MockService struct {
	RegisterClinicFn          func(ctx context.Context, ownerID string, clinic models.Clinic) *structs.IAppError
	RegisterClinicOwnerFn     func(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	SearchClinicFn            func(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError)
	LoginClinicOwnerFn        func(ctx context.Context, email string, password string) (string, *structs.IAppError)
	LoginDoctorFn             func(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchDoctorFn            func(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError)
	AddDoctorToClinicFn       func(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError
	SearchOwnerFn             func(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError)
	RegisterDoctorFn          func(ctx context.Context, doctor models.Doctor) *structs.IAppError
	VerifyAddDoctorToClinicFn func(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError
	AddAppointmentFn          func(ctx context.Context, appointmentDetails models.Appointment) *structs.IAppError
}

// ---- Interface Methods ----

func (m *MockService) RegisterClinic(ctx context.Context, ownerID string, clinic models.Clinic) *structs.IAppError {
	if m.RegisterClinicFn != nil {
		return m.RegisterClinicFn(ctx, ownerID, clinic)
	}
	return nil
}

func (m *MockService) VerifyAddDoctorToClinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError {
	if m.AddDoctorToClinicFn != nil {
		return m.VerifyAddDoctorToClinicFn(ctx, otp, doctorID, clinicID)
	}
	return nil
}
func (m *MockService) RegisterClinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError {
	if m.RegisterClinicOwnerFn != nil {
		return m.RegisterClinicOwnerFn(ctx, ownerDetails)
	}
	return nil
}

func (m *MockService) SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {
	if m.SearchClinicFn != nil {
		return m.SearchClinicFn(ctx, filter)
	}
	return nil, nil
}

func (m *MockService) LoginClinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	if m.LoginClinicOwnerFn != nil {
		return m.LoginClinicOwnerFn(ctx, email, password)
	}
	return "", nil
}

func (m *MockService) LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	if m.LoginDoctorFn != nil {
		return m.LoginDoctorFn(ctx, email, password)
	}
	return "", nil
}

func (m *MockService) SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError) {
	if m.SearchDoctorFn != nil {
		return m.SearchDoctorFn(ctx, filter)
	}
	return nil, nil
}

func (m *MockService) AddDoctorToClinic(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError {
	if m.AddDoctorToClinicFn != nil {
		return m.AddDoctorToClinicFn(ctx, clinicDetails)
	}
	return nil
}

func (m *MockService) SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError) {
	if m.SearchOwnerFn != nil {
		return m.SearchOwnerFn(ctx, filter)
	}
	return nil, nil
}

func (m *MockService) RegisterDoctor(ctx context.Context, doctor models.Doctor) *structs.IAppError {
	if m.RegisterDoctorFn != nil {
		return m.RegisterDoctorFn(ctx, doctor)
	}
	return nil
}

func (m *MockService) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) *structs.IAppError {
	if m.AddAppointmentFn != nil {
		return m.AddAppointmentFn(ctx, appointmentDetails)
	}
	return nil
}
