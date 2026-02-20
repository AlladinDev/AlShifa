package controller

import (
	DTO "AlShifa/clinic/dtos"
	"AlShifa/clinic/models"
	structs "AlShifa/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockService implements IService for testing purposes
type MockService struct {
	RegisterclinicFn          func(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) *structs.IAppError
	RegisterclinicOwnerFn     func(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	SearchclinicFn            func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError)
	LoginclinicOwnerFn        func(ctx context.Context, email string, password string) (string, *structs.IAppError)
	LoginDoctorFn             func(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchDoctorFn            func(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError)
	AddDoctorToclinicFn       func(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError
	SearchOwnerFn             func(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError)
	RegisterDoctorFn          func(ctx context.Context, doctor models.Doctor) *structs.IAppError
	VerifyAddDoctorToclinicFn func(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError
	AddAppointmentFn          func(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (int, *structs.IAppError)
	AppointmentSlotsBookedFn  func(ctx context.Context, slotDetails models.SlotDetails) ([]models.Slot, *structs.IAppError)
	DoctorWithItsclinicsFn    func(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, *structs.IAppError)
}

// ---- Interface Methods ----

func (m *MockService) Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) *structs.IAppError {
	if m.RegisterclinicFn != nil {
		return m.RegisterclinicFn(ctx, ownerID, clinic)
	}
	return nil
}

func (m *MockService) VerifyAddDoctorToclinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError {
	if m.AddDoctorToclinicFn != nil {
		return m.VerifyAddDoctorToclinicFn(ctx, otp, doctorID, clinicID)
	}
	return nil
}
func (m *MockService) RegisterclinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError {
	if m.RegisterclinicOwnerFn != nil {
		return m.RegisterclinicOwnerFn(ctx, ownerDetails)
	}
	return nil
}

func (m *MockService) Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError) {
	if m.SearchclinicFn != nil {
		return m.SearchclinicFn(ctx, filter)
	}
	return nil, nil
}

func (m *MockService) LoginclinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	if m.LoginclinicOwnerFn != nil {
		return m.LoginclinicOwnerFn(ctx, email, password)
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

func (m *MockService) AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError {
	if m.AddDoctorToclinicFn != nil {
		return m.AddDoctorToclinicFn(ctx, clinicDetails)
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

func (m *MockService) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError) {
	if m.AddAppointmentFn != nil {
		return m.AddAppointmentFn(ctx, 1, appointmentDetails)
	}
	return 0, nil
}

func (m *MockService) AppointmentSlotsBooked(ctx context.Context, slotDetails models.SlotDetails) ([]models.Slot, *structs.IAppError) {
	if m.AppointmentSlotsBookedFn != nil {
		return m.AppointmentSlotsBookedFn(ctx, slotDetails)
	}
	return nil, nil
}

func (m *MockService) DoctorWithItsclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, *structs.IAppError) {
	if m.DoctorWithItsclinicsFn != nil {
		return m.DoctorWithItsclinicsFn(ctx, filter)
	}
	return nil, nil
}
