// Package service provides service functions for appointment module
package service

import (
	"AlShifa/appointment/interfaces"
	"AlShifa/appointment/models"
	"AlShifa/constants"
	"net/http"
	"time"

	"AlShifa/structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo          interfaces.IRepository
	clinicService interfaces.IClinicModule
}

func NewService(repo interfaces.IRepository, clinicService interfaces.IClinicModule) *Service {
	return &Service{
		repo:          repo,
		clinicService: clinicService,
	}
}

var _ interfaces.IService = (*Service)(nil)

func (s *Service) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError) {
	//call the clinic module pass the details and it will automatically check whether this clinic exits doctor exists and if appointment date is possible
	err, doctorName, clinicName, clinicAddress, clinicMaxAppointments := s.clinicService.ClinicDoctorDetails(ctx, appointmentDetails.ClinicID, appointmentDetails.UserID, appointmentDetails.AppointmentDate)
	if err != nil {
		return 0, err
	}

	//now override doctorName clinicName and clinicAddress of post method data with details returned from clinicmodule for better safety so that user doesnt send any arbitary clinic doctor name
	appointmentDetails.DoctorName = doctorName
	appointmentDetails.ClinicName = clinicName
	appointmentDetails.ClinicAddress = clinicAddress

	//now add some default things
	appointmentDetails.CreatedAt = time.Now()
	appointmentDetails.ID = primitive.NewObjectID()
	appointmentDetails.Status = constants.StatusAppointmentPending

	//now call the repo to save data
	slotNumber, appointmentSavingErr := s.repo.AddAppointment(ctx, clinicMaxAppointments, appointmentDetails)
	if appointmentSavingErr != nil {
		return 0, &structs.IAppError{
			Message:    "Failed to book appointment",
			Reason:     appointmentSavingErr.Error(),
			ErrorObj:   appointmentSavingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return slotNumber, nil
}

func (s *Service) FetchAppointments(ctx context.Context, filters bson.M) ([]models.Appointment, *structs.IAppError) {
	appointments, err := s.repo.FetchAppointments(ctx, filters)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "failed to fetch appointments",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	return appointments, nil
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) *structs.IAppError {
	if err := s.repo.UpdateAppointmentStatus(ctx, appointmentID, status); err != nil {
		return &structs.IAppError{
			Message:    "Failed to update appointment",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}

	return nil
}

func (s *Service) FetchAppointmentDaysBooked(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) ([]models.Slot, *structs.IAppError) {
	//here first whether this clinic exists or not
	clinicSearchErr := s.clinicService.ClinicExists(ctx, clinicID)
	if clinicSearchErr != nil {
		return nil, clinicSearchErr
	}

	//now check whether this doctor exists or not
	doctorSearchErr := s.clinicService.DoctorExists(ctx, doctorID)
	if doctorSearchErr != nil {
		return nil, doctorSearchErr
	}

	//now check whether this doctor clinic mapping exists or not
	mappingErr := s.clinicService.DoctorClinicMappingExists(ctx, clinicID, doctorID)
	if mappingErr != nil {
		return nil, mappingErr
	}

	//now fetch maxAppointments for this clinic
	maxAppointments, err := s.clinicService.FetchMaxAppointments(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	daysBooked, daysBookedErr := s.repo.FetchAppointmentDaysBooked(ctx, maxAppointments, doctorID, clinicID)
	if daysBookedErr != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch appointments days booked",
			StatusCode: http.StatusInternalServerError,
			Reason:     daysBookedErr.Error(),
			ErrorObj:   daysBookedErr,
		}
	}

	return daysBooked, nil
}
