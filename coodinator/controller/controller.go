// Package controller provides handlers for coordinator service
package controller

import (
	"encoding/json"
	"net/http"

	"github.com/AlladinDev/AlShifa/coodinator/interfaces"
	"github.com/AlladinDev/AlShifa/coodinator/validators"
	sharedModels "github.com/AlladinDev/AlShifa/models"
	"github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/utils"
)

type Controller struct {
	sv interfaces.ICoodinatorService
}

func NewController(service interfaces.ICoodinatorService) *Controller {
	return &Controller{
		sv: service,
	}
}

func (c *Controller) BookAppointment(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var appointmentDetails sharedModels.Appointment
	if err := json.NewDecoder(req.Body).Decode(&appointmentDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid json Details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//do some validations here
	validationErrs := validators.ValidateAppointmentDetails(appointmentDetails)
	if validationErrs != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Details",
			Reason:     "Invalid Details",
			ErrorObj:   validationErrs,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	slotNumber, err := c.sv.BookAppointment(ctx, appointmentDetails)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppSuccess{
		Message:    "Appointment Booked Successfully",
		Data:       slotNumber,
		StatusCode: http.StatusCreated,
	})

}
