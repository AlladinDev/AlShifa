package service

import (
	"context"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/dtos"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TService struct {
	repo interfaces.IRepo
}

func NewService(Repo interfaces.IRepo) interfaces.IService {
	return &TService{
		repo: Repo,
	}
}

func (s TService) IsAppointmentDatePossible(ctx context.Context, appointmentDate time.Time, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error) {
	return s.repo.CheckAppointmentDatePossible(ctx, clinicID, doctorID, appointmentDate)
}

func (s TService) FetchClinicWithDoctors(ctx context.Context, filter bson.M) ([]dtos.ClinicDoctorDTO, *response.IAppError) {
	data, err := s.repo.FetchClinicWithDoctors(ctx, filter)
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed to fetch clinic details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return data, nil
}

func (s TService) FetchDoctorWithClinics(ctx context.Context, filter bson.M) ([]dtos.DoctorWithClinic, *response.IAppError) {
	data, err := s.repo.FetchDoctorWithClinics(ctx, filter)
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed to fetch Data",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return data, nil
}
