package service

import (
	DTO "AlShifa/clinic/dtos"
	interfaces "AlShifa/clinic/interfaces"
	"AlShifa/clinic/models"
	"context"
	"log"

	sharedModels "AlShifa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepo struct {
	RegisterclinicMockFn      func(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	RegisterclinicOwnerMockFn func(ctx context.Context, owner models.Owner) error
	GetOwnerDetailsMockFn     func(ctx context.Context, filter bson.M) ([]models.Owner, error)
	SearchclinicMockFn        func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error)
	RegisterDoctorMockFn      func(ctx context.Context, doctorDetails models.Doctor) error
	SearchDoctorsMockFn       func(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error)
	SearchDoctorMockFn        func(ctx context.Context, filter bson.M) (models.Doctor, error)
	AddDoctorToclinicMockFn   func(ctx context.Context, clinicDetails models.AddDoctorToclinic) error
	AddAppointmentMockFn      func(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (int, error)
	SearchclinicByIDFn        func(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error)
	FetchDoctorAtclinicsFn    func(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, error)
	FetchAppointmentsFn       func(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointments, error)
	AppointmentSlotsBookedFn  func(ctx context.Context, maxAppointments int, clinicID primitive.ObjectID, doctorID primitive.ObjectID) ([]models.Slot, error)
}

var _ interfaces.IRepository = (*MockRepo)(nil)

func (m *MockRepo) Registerclinic(ctx context.Context, clinicID primitive.ObjectID, clinic models.Clinic) error {
	if m.RegisterclinicMockFn == nil {
		log.Fatal("RegisterclinicMockFn not implemented in MockRepository of clinic service")
	}
	return m.RegisterclinicMockFn(ctx, clinicID, clinic)
}

func (m *MockRepo) FetchAppointments(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointments, error) {
	if m.FetchAppointmentsFn == nil {
		log.Fatal("FetchAppointmentsFn not implemented in MockRepository of clinic service")
	}
	return m.FetchAppointmentsFn(ctx, groupBy, filter)
}
func (m *MockRepo) SearchclinicByID(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error) {
	if m.SearchclinicByIDFn == nil {
		log.Fatal("RegisterclinicMockFn not implemented in MockRepository of clinic service")
	}
	return m.SearchclinicByIDFn(ctx, clinicID)
}
func (m *MockRepo) AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) error {
	if m.AddDoctorToclinicMockFn == nil {
		log.Fatal("AddDoctorToclinicMockFn not implemented in MockRepository of clinic service")
	}
	return m.AddDoctorToclinicMockFn(ctx, clinicDetails)
}

func (m *MockRepo) RegisterclinicOwner(ctx context.Context, owner models.Owner) error {
	if m.RegisterclinicOwnerMockFn == nil {
		log.Fatal("RegisterclinicOwner not implemented in MockRepository of clinic service")
	}
	return m.RegisterclinicOwnerMockFn(ctx, owner)
}

func (m *MockRepo) Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {
	if m.SearchclinicMockFn == nil {
		log.Fatal("Searchclinic not implemented in MockRepository of clinic service")
	}
	return m.SearchclinicMockFn(ctx, filter)
}

func (m *MockRepo) GetOwnerDetails(ctx context.Context, filter bson.M) ([]models.Owner, error) {
	if m.GetOwnerDetailsMockFn == nil {
		log.Fatal("GetOwnerDetailsMockFn not implemented in MockRepository of clinic service")
	}
	return m.GetOwnerDetailsMockFn(ctx, filter)
}

func (m *MockRepo) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error {
	if m.RegisterDoctorMockFn == nil {
		log.Fatal("RegisterDoctorMockFn not implemented in MockRepository of clinic service")
	}
	return m.RegisterDoctorMockFn(ctx, doctorDetails)
}

func (m *MockRepo) SearchDoctor(ctx context.Context, filter bson.M) (models.Doctor, error) {
	if m.SearchDoctorMockFn == nil {
		log.Fatal("SearchDoctorMockFn not implemented in MockRepository of clinic service")
	}
	return m.SearchDoctorMockFn(ctx, filter)
}

func (m *MockRepo) SearchDoctors(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error) {
	if m.SearchDoctorsMockFn == nil {
		log.Fatal("SearchDoctorsMockFn not implemented in MockRepository of clinic service")
	}
	return m.SearchDoctorsMockFn(ctx, filter)
}

func (m *MockRepo) AddAppointment(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (int, error) {
	if m.AddAppointmentMockFn == nil {
		log.Fatal("AddAppointmentMockFn not implemented in MockRepository of clinic service")
	}
	return m.AddAppointmentMockFn(ctx, maxAppointments, appointmentDetails)
}

func (m *MockRepo) AppointmentSlotsBooked(ctx context.Context, maxAppointments int, clinicID primitive.ObjectID, doctorID primitive.ObjectID) ([]models.Slot, error) {
	if m.AppointmentSlotsBookedFn == nil {
		log.Fatal("AppointmentSlotsBookedFn not implemented in MockRepository of clinic service")
	}
	return m.AppointmentSlotsBookedFn(ctx, maxAppointments, clinicID, doctorID)
}

func (m *MockRepo) FetchDoctorAtclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, error) {
	if m.FetchDoctorAtclinicsFn == nil {
		log.Fatal("FetchDoctorAtclinicsFn is not implemented in mock repo of clinic module")
	}
	return m.FetchDoctorAtclinics(ctx, filter)
}
