// Package controller provides HTTP handlers for managing clinic-related operations.
package controller

import (
	service "AlShifa/Clinic/Service"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	middleware "AlShifa/Middleware"
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Invalid Details Provided", "Json Error"))
		return
	}

	//here validate clinic details
	validationErrors := validators.ValidateClinicDetails(&clinicRegistrationDetails.Clinic)
	if len(validationErrors) != 0 {
		fmt.Print(validationErrors)
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(validationErrors, 400, "Invalid Details", "Validation Failed"))
		return
	}

	registrationErr := controller.Service.RegisterClinic(ctx, clinicRegistrationDetails.OwnerID, clinicRegistrationDetails.Clinic)
	if registrationErr != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, registrationErr)
		return
	}

	response := structs.IAppSuccess{
		Message:    "Clinic Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response)
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
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Invalid Details", "InValid Json"))
		return
	}

	//now validate here
	if err := validators.ValidateOwnerDetails(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Invalid Details", "Validation Failure"))
		return
	}

	if err := controller.Service.RegisterClinicOwner(ctx, ownerDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	response := structs.IAppSuccess{
		Message:    "Owner Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response)
}

func (controller *Controller) SearchClinic(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	// Parse query parameters
	params := req.URL.Query()

	// Initialize empty filter
	filters := bson.M{}

	// Iterate over query params
	for key, values := range params {
		if len(values) == 0 {
			continue
		}
		value := values[0] // take first value for simplicity

		// Special handling for clinicId
		if key == "id" {
			objID, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 400, "Unable To Fetch Clinic Details", "Server Error"))
				return
			}
			filters["_id"] = objID
		} else {
			// Treat all other fields as string match
			filters[key] = value
		}
	}

	// Call your service with filters
	clinics, err := controller.Service.SearchClinic(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(res, http.StatusOK, utils.ReturnAppSuccess(200, "Fetched Successfully", clinics))
}

func (controller *Controller) SearchOwner(res http.ResponseWriter, req *http.Request) {
	userRole := req.Context().Value(middleware.ContextUserRoleKey).(string)
	userID := req.Context().Value(middleware.ContextUserIDKey).(string)

	if userRole == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(nil, 400, "Missing Role", "Missing Role"))
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	// Parse query parameters
	params := req.URL.Query()

	// Initialize empty filter
	filters := bson.M{}

	if userRole == utils.RoleAdmin {
		//only admin can fetch all owners using various filters and clinic owner can see just their details by id in their jwt token
		// Iterate over query params
		for key, values := range params {
			if len(values) == 0 {
				continue
			}
			value := values[0] // take first value for simplicity
			// Treat all other fields as string match
			filters[key] = value
		}
	}

	userMongoDBID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 400, "Invalid UserID", err.Error()))
		return
	}

	filters["_id"] = userMongoDBID

	owner, userSearchErr := controller.Service.SearchOwner(ctx, filters)
	if userSearchErr != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, *userSearchErr)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, utils.ReturnAppSuccess(200, "Fetched Successfully", owner))
}

func (controller *Controller) RegisterDoctor(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()
	var doctor models.Doctor

	if err := json.NewDecoder(req.Body).Decode(&doctor); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Registration Failed", "Invalid Json"))
		return
	}

	//here do validation
	validationErrors := validators.ValidateDoctor(doctor)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(validationErrors, 400, "Registration Failed", "Invalid Details"))
		return
	}

	if err := controller.Service.RegisterDoctor(ctx, doctor); err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	//here it means doctor is successfully registered
	_ = utils.WriteResponse(res, http.StatusCreated, structs.IAppSuccess{
		Message:    "Doctor Registered Successfully",
		Data:       nil,
		StatusCode: 200,
	})

}

func (controller *Controller) SearchDoctor(res http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		_ = utils.InvalidMethodResponse("GET", res)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	// Parse query parameters
	params := req.URL.Query()

	// Initialize empty filter
	filters := bson.M{}

	// Iterate over query params
	for key, values := range params {
		if len(values) == 0 {
			continue
		}
		value := values[0] // take first value for simplicity

		// Special handling for clinicId
		if key == "id" {
			objID, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 400, "Unable To Fetch Doctors Details", "Server Error"))
				return
			}
			filters["_id"] = objID
		} else {
			// Treat all other fields as string match
			filters[key] = value
		}
	}

	doctors, err := controller.Service.SearchDoctor(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Successfully Fetched Details",
		Data:       doctors,
		StatusCode: 200,
	})

}

func (controller *Controller) LoginClinicOwner(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var loginDetails structs.LoginDetails

	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Login Failed", "Invalid Json"))
		return
	}
	jwtToken, err := controller.Service.LoginClinicOwner(ctx, loginDetails.Email, loginDetails.Password)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}
	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Login Successful",
		Data:       utils.JwtPrefix + jwtToken,
		StatusCode: 200,
	})
}

func (controller *Controller) LoginDoctor(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var loginDetails structs.LoginDetails
	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Login Failed", "Invalid Json"))
		return
	}
	jwtToken, err := controller.Service.LoginDoctor(ctx, loginDetails.Email, loginDetails.Password)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}
	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Login Successful",
		Data:       utils.JwtPrefix + jwtToken,
		StatusCode: 200,
	})
}
