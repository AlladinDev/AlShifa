// Package controller provides controllers for appointment module
package controller

import (
	"encoding/json"
	"net/http"

	"github.com/AlladinDev/AlShifa/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/appointment/models"
	"github.com/AlladinDev/AlShifa/appointment/validators"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/utils"

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
	validationErrors := validators.ValidateAppointment(appointmentDetails)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Appointment Details",
			Reason:     "Invalid appointment details",
			ErrorObj:   validationErrors,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now from ctx get the userid and add it to appointment payload
	userIDAny := ctx.Value(constants.KeyUserID)
	userIDErr, userMongoID := utils.ParseUserID(userIDAny)
	if userIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid UserID",
			Reason:     userIDErr.Error(),
			ErrorObj:   userIDErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	appointmentDetails.UserID = userMongoID

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

	//get the userid and user type from ctx passed by jwt middleware
	userTypeAny := req.Context().Value(constants.KeyUserRole)

	userType, ok := userTypeAny.(string)
	if !ok {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid User Type",
			Reason:     "Invalid UserType",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   nil,
		})
		return
	}

	//now get the userid also
	userIDAny := req.Context().Value(constants.KeyUserID)
	userIDErr, userID := utils.ParseUserID(userIDAny)
	if userIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid UserID",
			Reason:     userIDErr.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   nil,
		})
		return
	}

	paramsPassed := req.URL.Query()

	filters := bson.M{}

	//add some default filters like userid and usertype so that only associated user or clinic and see its own related appointments

	filters["userID"] = userID
	filters["userType"] = userType

	//get the userID because a user is allowed to see only his appointments

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
