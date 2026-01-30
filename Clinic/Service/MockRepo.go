package service

import (
	interfaces "AlShifa/Clinic/Interfaces"
	"AlShifa/Clinic/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockRepo struct {
	RegisterClinicMockFn      func(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	RegisterClinicOwnerMockFn func(ctx context.Context, owner models.Owner) error
	GetOwnerDetailsMockFn     func(ctx context.Context, filter bson.M) ([]models.Owner, error)
	SearchClinicMockFn        func(ctx context.Context, filter bson.M) ([]models.Clinic, error)
	RegisterDoctorMockFn      func(ctx context.Context, doctorDetails models.Doctor) error
	SearchDoctorsMockFn       func(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error)
	SearchDoctorMockFn        func(ctx context.Context, filter bson.M) (models.Doctor, error)
	AddDoctorToClinicMockFn   func(ctx context.Context, clinicDetails models.AddDoctorToClinic) error
}

var _ interfaces.IRepository = (*MockRepo)(nil)

func (m *MockRepo) RegisterClinic(ctx context.Context, clinicID primitive.ObjectID, clinic models.Clinic) error {
	if m.RegisterClinicMockFn == nil {
		log.Fatal("RegisterClinicMockFn not implemented in MockRepository of clinic service")
	}
	return m.RegisterClinicMockFn(ctx, clinicID, clinic)
}

func (m *MockRepo) AddDoctorToClinic(ctx mongo.SessionContext, clinicDetails models.AddDoctorToClinic) error {
	if m.AddDoctorToClinicMockFn == nil {
		log.Fatal("AddDoctorToClinicMockFn not implemented in MockRepository of clinic service")
	}
	return m.AddDoctorToClinicMockFn(ctx, clinicDetails)
}

func (m *MockRepo) RegisterClinicOwner(ctx context.Context, owner models.Owner) error {
	if m.RegisterClinicOwnerMockFn == nil {
		log.Fatal("RegisterClinicOwner not implemented in MockRepository of clinic service")
	}
	return m.RegisterClinicOwnerMockFn(ctx, owner)
}

func (m *MockRepo) SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, error) {
	if m.SearchClinicMockFn == nil {
		log.Fatal("SearchClinic not implemented in MockRepository of clinic service")
	}
	return m.SearchClinicMockFn(ctx, filter)
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
