// Package controller provides controllers for appointment module
package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/internal/appointment/models"
	"github.com/AlladinDev/AlShifa/internal/appointment/validators"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service interfaces.IService
}

func NewController(service interfaces.IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) AddAppointment(res http.ResponseWriter, req *http.Request) {

	var appointmentDetails models.Appointment
	if err := json.NewDecoder(req.Body).Decode(&appointmentDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid Appointment Details",
			Reason:     "json details are invalid",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	ctx := req.Context()

	//now do some validations
	validationErrors := validators.ValidateAppointment(appointmentDetails)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid Appointment Details",
			Reason:     "Invalid appointment details",
			ErrorObj:   validationErrors,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now get the userid from req.context
	userIDAny := req.Context().Value(constants.KeyUserID)
	userMongoDBID, idErr := utils.ParseUserID(userIDAny)
	if idErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Userid invalid",
			Reason:     idErr.Error(),
			ErrorObj:   idErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now call service layer
	slot, appointmentPaymentToken, err := c.service.AddAppointment(ctx, userMongoDBID, appointmentDetails)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameAppointmentBookingPaymentToken,
		Value:    appointmentPaymentToken,
		Expires:  time.Now().Add(constants.AppointmentBookingPaymentExpiryTime),
		Secure:   true,
		HttpOnly: true,
	})
	_ = utils.WriteResponse(res, http.StatusCreated, response.IAppSuccess{
		Message:    "Appointment Booked Successfully",
		Data:       slot,
		StatusCode: http.StatusCreated,
	})

}

func (c *Controller) UpdateAppointmentStatus(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var appointmentUpdationDTO struct {
		AppointmentID primitive.ObjectID `json:"_id" bson:"_id"`
		Status        bool               `json:"status" bson:"status"`
	}
	if err := json.NewDecoder(req.Body).Decode(&appointmentUpdationDTO); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid Json Details",
			Reason:     "Invalid json details provided",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err := c.service.UpdateAppointmentStatus(ctx, appointmentUpdationDTO.AppointmentID, appointmentUpdationDTO.Status)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, response.IAppSuccess{
		Message:    "Updated Appointment Successfully",
		StatusCode: http.StatusOK,
		Data:       nil,
	})

}
