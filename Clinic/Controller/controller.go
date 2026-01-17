// Package controller provides HTTP handlers for managing clinic-related operations.
package controller

import (
	service "AlShifa/Clinic/Service"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Controller struct {
	Service *service.ClinicService
}

type ClinicRegistration struct {
	OwnerID string        `json:"ownerId"`
	Clinic  models.Clinic `json:"clinicDetails"`
}

func NewController(svr *service.ClinicService) *Controller {
	return &Controller{
		Service: svr,
	}
}

func (controller *Controller) RegisterClinic(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		_ = utils.InvalidMethodResponse("POST", res)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var clinicRegistrationDetails ClinicRegistration
	if err := json.NewDecoder(req.Body).Decode(&clinicRegistrationDetails); err != nil {
		_ = utils.WriteResponse(res, utils.ReturnAppError(err, 500, "Invalid Details Provided", "Json Error"))
		return
	}

	//here validate clinic details
	validationErrors := validators.ValidateClinicDetails(&clinicRegistrationDetails.Clinic)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, utils.ReturnAppError(validationErrors, 400, "Invalid Details", "Invalid Details"))
		return
	}

	response := controller.Service.RegisterClinic(ctx, clinicRegistrationDetails.OwnerID, clinicRegistrationDetails.Clinic)
	_ = utils.WriteResponse(res, response)
}

func (controller *Controller) RegisterOwner(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		_ = utils.InvalidMethodResponse("POST", res)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var ownerDetails models.Owner
	if err := json.NewDecoder(req.Body).Decode(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, utils.ReturnAppError(err, 500, "Invalid Details", "InValid Json"))
		return
	}

	//now validate here
	if err := validators.ValidateOwnerDetails(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, utils.ReturnAppError(err, 500, "Invalid Details", "Validation Failure"))
		return
	}

	fmt.Print("reached here")
	if err := controller.Service.RegisterClinicOwner(ctx, ownerDetails); err != nil {
		_ = utils.WriteResponse(res, err)
		return
	}

	response := structs.IAppSuccess{
		Message:    "Owner Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	}

	_ = utils.WriteResponse(res, response)
}
