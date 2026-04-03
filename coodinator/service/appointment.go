// Package service provides service functions for coordinator service
package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/coodinator/interfaces"
	sharedModels "github.com/AlladinDev/AlShifa/models"
	"github.com/AlladinDev/AlShifa/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentService struct {
	ClinicService      interfaces.IClinicService
	AppointmentService interfaces.IAppointmentService
	mongoClient        *mongo.Client
}

func NewAppointmentCoordinator(clinicService interfaces.IClinicService, mongoclient *mongo.Client, appointmentService interfaces.IAppointmentService) *AppointmentService {
	return &AppointmentService{
		ClinicService:      clinicService,
		mongoClient:        mongoclient,
		AppointmentService: appointmentService,
	}
}

func (a *AppointmentService) BookAppointment(ctx context.Context, appointmentDetails sharedModels.Appointment) (int, *structs.IAppError) {
	//call the clinic module pass the details and it will automatically check whether this clinic exits doctor exists and if appointment date is possible
	err, doctorName, clinicName, clinicAddress, clinicMaxAppointments := a.ClinicService.ClinicDoctorDetails(ctx, appointmentDetails.ClinicID, appointmentDetails.DoctorID, appointmentDetails.AppointmentDate)
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

	//start mongodb session
	session, sessionErr := a.mongoClient.StartSession()
	if sessionErr != nil {
		return 0, &structs.IAppError{
			Message:    "Failed to book appointment",
			Reason:     sessionErr.Error(),
			ErrorObj:   sessionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	txnFn := func(sessCtx mongo.SessionContext) (any, error) {
		//now call the repo to save data
		slotNumber, appointmentSavingErr := a.AppointmentService.AddAppointment(sessCtx, clinicMaxAppointments, appointmentDetails)
		if appointmentSavingErr != nil {
			return 0, appointmentSavingErr
		}

		if err := a.ClinicService.DeductClinicMoneyForAppointment(sessCtx, appointmentDetails.ClinicID); err != nil {
			return 0, err
		}

		return slotNumber, nil
	}

	transactionRes, transactionErr := session.WithTransaction(ctx, txnFn)
	if transactionErr != nil {
		return 0, &structs.IAppError{
			Message:    "Failed to book appointment",
			Reason:     transactionErr.Error(),
			ErrorObj:   transactionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	slotNumber, ok := transactionRes.(int)
	if !ok {
		return 0, &structs.IAppError{
			Message:    "Failed to Book Appointment",
			Reason:     "failed to convert slotNumber any to int",
			ErrorObj:   errors.New("failed to convert slotNumber any to int"),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return slotNumber, nil
}
