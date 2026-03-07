// Package controller provides controllers for appointment module
package controller

import (
	"AlShifa/appointment/interfaces"
	"AlShifa/appointment/models"
	"AlShifa/structs"
	"AlShifa/utils"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Appointment Details",
			Reason:     "json details are invalid",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	ctx := req.Context()

	//now do some validations

	//now call service layer
	slot, err := c.service.AddAppointment(ctx, appointmentDetails)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, structs.IAppSuccess{
		Message:    "Appointment Booked Successfully",
		Data:       slot,
		StatusCode: http.StatusCreated,
	})

}

func (c *Controller) FetchAppointments(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	paramsPassed := req.URL.Query()

	filters := bson.M{}

	_ = utils.TransformParamIDS(paramsPassed, filters)
	appointments, err := c.service.FetchAppointments(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Successfully Fetched Appointments",
		StatusCode: http.StatusOK,
		Data:       appointments,
	})

}

func (c *Controller) UpdateAppointmentStatus(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var appointmentUpdationDTO struct {
		AppointmentID primitive.ObjectID `json:"_id" bson:"_id"`
		Status        bool               `json:"status" bson:"status"`
	}
	if err := json.NewDecoder(req.Body).Decode(&appointmentUpdationDTO); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
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

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Updated Appointment Successfully",
		StatusCode: http.StatusOK,
		Data:       nil,
	})

}
