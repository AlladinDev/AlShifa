// Package service provides service functions for appointment module
package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/internal/appointment/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"

	"context"

	"github.com/AlladinDev/AlShifa/internal/shared/response"

	appDTOS "github.com/AlladinDev/AlShifa/internal/shared/dtos"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo          interfaces.IRepository
	clinicService interfaces.IClinicModule
	doctorService interfaces.IDoctorModule
	mongoClient   *mongo.Client
	planService   interfaces.IPlanService
	doctorClinic  interfaces.IClinicDoctor
	encrypt       func(data any) (string, error)
	broker        appInterfaces.IMessageBroker
	kvStore       appInterfaces.Cache[string, []byte]
}

func NewService(repo interfaces.IRepository,
	clinicService interfaces.IClinicModule,
	mongoClient *mongo.Client,
	PlanService interfaces.IPlanService,
	DoctorService interfaces.IDoctorModule,
	DoctorClinic interfaces.IClinicDoctor,
	EncryptionFn func(data any) (string, error),
	Broker appInterfaces.IMessageBroker,
	KVStore appInterfaces.Cache[string, []byte],

) *Service {
	return &Service{
		repo:          repo,
		clinicService: clinicService,
		mongoClient:   mongoClient,
		doctorService: DoctorService,
		doctorClinic:  DoctorClinic,
		planService:   PlanService,
		encrypt:       EncryptionFn,
		broker:        Broker,
	}
}

var _ interfaces.IService = (*Service)(nil)

func (s *Service) AddAppointment(ctx context.Context, userID primitive.ObjectID, appointmentDetails models.Appointment) (slot int, appointmentToken string, respErr *response.IAppError) {
	//now add userid to appointment doc
	appointmentDetails.UserID = userID

	//call doctorservice to check whether this doctor exists or not
	doctorExists, err := s.doctorService.DoctorExists(ctx, bson.M{"_id": appointmentDetails.DoctorID, "department": appointmentDetails.Department, "doctorName": appointmentDetails.DoctorName})
	if err != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if !doctorExists {
		return 0, "", &response.IAppError{
			Message:    "Doctor doesnt exist",
			Reason:     "doctor with these details doesnt exist",
			ErrorObj:   nil,
			StatusCode: http.StatusNotFound,
		}
	}

	//now call clinic service to check whether this clinic exists or not
	clinicExists, clinicSearchErr := s.clinicService.ClinicExists(ctx, bson.M{"_id": appointmentDetails.ClinicID, "name": appointmentDetails.ClinicName, "address": appointmentDetails.ClinicAddress, "maxAppointments": appointmentDetails.ClinicMaxAppointments})
	if clinicSearchErr != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     clinicSearchErr.Error(),
			ErrorObj:   clinicSearchErr,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if !clinicExists {
		return 0, "", &response.IAppError{
			Message:    "Clinic not present with these details",
			Reason:     "no clinic found",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now here call the clinic doctor mapping module to check whether on this requested appointment date ,doctor  is available or not or in other words is appointment possible or not
	appointmentPossible, err := s.doctorClinic.IsAppointmentDatePossible(ctx, appointmentDetails.AppointmentDate, appointmentDetails.DoctorID, appointmentDetails.ClinicID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, "", &response.IAppError{
				Message:    "This doctor is not onboarded to this clinic yet",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusNotFound,
			}
		}
	}

	if !appointmentPossible {
		return 0, "", &response.IAppError{
			Message:    "Appointment not possible on this date",
			Reason:     "doctor is not available on this date",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now here call the plan service to get user plan details it will return planid amount to deduct from patint for appointment
	_, amountToDeduct, planErr := s.planService.FetchAmountToDeduct(ctx, bson.M{"name": constants.PlanAppointment})
	if planErr != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     planErr.Error(),
			ErrorObj:   planErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now add some default things like status createdAt
	appointmentDetails.AppointmentBookingFeeStatus = constants.AppointmentBookingFeesPending
	appointmentDetails.AppointmentStatus = constants.AppointmentPending
	appointmentDetails.AppointmentFeesDetails = models.IAppointmentFees{
		Amount:    0,
		Paid:      false,
		CreatedAt: time.Now(),
	}

	appointmentDetails.CreatedAt = time.Now()
	appointmentDetails.ID = primitive.NewObjectID()

	//now create the token for this appointment so that payment module when it receives it through cookie which this function will set,
	//through this token payment module can get the plan and decide how much amount to ask from customer and also emit event for success or failure of appointment to appointment module
	//we will then use that appointment id and then get the corresponding appointment doc and update it
	token := appDTOS.AppointmentPaymentToken{
		AppointmentID:  appointmentDetails.ID,
		AmountToDeduct: amountToDeduct,
		UserID:         userID,
		ExpiresAt:      time.Now().Add(constants.AppointmentBookingPaymentExpiryTime),
	}

	///now encrypt this token into json format
	byteToken, jsonErr := json.Marshal(token)
	if jsonErr != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     jsonErr.Error(),
			ErrorObj:   jsonErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	mongoSession, sessionErr := s.mongoClient.StartSession()
	if sessionErr != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     sessionErr.Error(),
			ErrorObj:   sessionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//here save the message in kv store
	if err := s.kvStore.Set(ctx, constants.AppointmentInitiated, byteToken, constants.AppointmentBookingPaymentExpiryTime*3); err != nil {
		return 0, "", &response.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	type TransactionRes struct {
		SlotDocumentID          primitive.ObjectID `json:"slotDocumentID"`
		SlotNumber              int                `json:"slotNumber"`
		AppointmentPaymentToken string             `json:"appointmentPaymentToken"`
	}

	transactionFn := func(mongoCtx mongo.SessionContext) (any, error) {
		respObj := TransactionRes{}
		//now call the repo to save data
		slotDocumentID, slotNumber, appointmentSavingErr := s.repo.AddAppointment(ctx, appointmentDetails.ClinicMaxAppointments, appointmentDetails)
		if appointmentSavingErr != nil {
			return respObj, appointmentSavingErr
		}

		//now encrypt this token
		encryptedToken, encryptionErr := s.encrypt(token)
		if encryptionErr != nil {
			return respObj, encryptionErr
		}

		respObj.SlotDocumentID = slotDocumentID
		respObj.SlotNumber = slotNumber
		respObj.AppointmentPaymentToken = encryptedToken
		return respObj, nil
	}

	dbResponseAny, err := mongoSession.WithTransaction(ctx, transactionFn)

	dbResponse, ok := dbResponseAny.(TransactionRes)
	if !ok {
		return 0, "", &response.IAppError{
			Message:    "Appointment booking failed try again",
			Reason:     "failed to parse transaction response for appointment document",
			ErrorObj:   nil,
			StatusCode: http.StatusInsufficientStorage,
		}
	}
	return dbResponse.SlotNumber, dbResponse.AppointmentPaymentToken, nil

}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) *response.IAppError {
	if err := s.repo.UpdateAppointmentStatus(ctx, appointmentID, status); err != nil {
		return &response.IAppError{
			Message:    "Failed to update appointment",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}

	return nil
}
