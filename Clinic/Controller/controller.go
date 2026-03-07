// Package controller provides HTTP handlers for managing clinic-related operations.
package controller

import (
	interfaces "AlShifa/clinic/interfaces"
	"AlShifa/clinic/models"
	validators "AlShifa/clinic/validators"
	"AlShifa/constants"
	structs "AlShifa/structs"
	utils "AlShifa/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service                            interfaces.IService
	ValidateAddDoctorToclinicDetailsFn func(details *models.ClinicDoctor) map[string]string
}

func NewController(svr interfaces.IService, validateAddDoctorToclinicFn func(details *models.ClinicDoctor) map[string]string) *Controller {
	return &Controller{
		Service:                            svr,
		ValidateAddDoctorToclinicDetailsFn: validateAddDoctorToclinicFn,
	}
}

func (controller *Controller) Registerclinic(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var clinicRegistrationDetails models.Clinic
	if err := json.NewDecoder(req.Body).Decode(&clinicRegistrationDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid Details Provided", "Json Error"))
		return
	}

	//here validate clinic details
	validationErrors := validators.ValidateClinicDetails(&clinicRegistrationDetails)
	if len(validationErrors) != 0 {
		fmt.Print(validationErrors)
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(validationErrors, 400, "Invalid Details", "Validation Failed"))
		return
	}

	//here extract the ownerId from req.context fed by jwt middleware
	ownerID := req.Context().Value(constants.KeyUserID)

	ownerMongodbIDErr, ownerMongoDBID := utils.ParseUserID(ownerID)
	if ownerMongodbIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid OwnerID",
			Reason:     "Owner id is invalid",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   errors.New("ownerid is not a valid mongodb ID"),
		})
		return
	}

	//here add this id to clinicregistration details so that user can send any other owners id
	clinicRegistrationDetails.OwnerID = ownerMongoDBID

	registrationErr := controller.Service.Registerclinic(ctx, clinicRegistrationDetails.OwnerID, clinicRegistrationDetails)
	if registrationErr != nil {
		_ = utils.WriteResponse(res, registrationErr.StatusCode, registrationErr)
		return
	}

	response := structs.IAppSuccess{
		Message:    "clinic Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response)
}

func (controller *Controller) Searchclinic(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	// Parse query parameters
	params := req.URL.Query()

	// Initialize empty filter
	filters := bson.M{}

	_ = utils.TransformParamIDS(params, filters)

	// Call your service with filters
	clinics, err := controller.Service.Searchclinic(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, utils.ReturnAppSuccess(200, "Fetched Successfully", clinics))
}

func (controller *Controller) SearchDoctor(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Parse query parameters
	params := req.URL.Query()

	// Initialize empty filter
	filters := bson.M{}

	_ = utils.TransformParamIDS(params, filters)

	doctors, err := controller.Service.FetchDoctors(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Successfully Fetched Details",
		Data:       doctors,
		StatusCode: 200,
	})

}

func (controller *Controller) AddDoctorToclinic(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	//first try to see if client id is present in req.context
	ownerID := req.Context().Value(constants.KeyUserID)
	if ownerID == "" {
		_ = utils.WriteResponse(res, 400, utils.ReturnAppError(errors.New("ownerid is missing"), 400, "OwnerID missing", "OwnerId is missing for authentication"))
		return
	}

	//now try to convert ownerID into string
	ownerIDStr, ok := ownerID.(string)
	if !ok {
		_ = utils.WriteResponse(res, 400, utils.ReturnAppError(errors.New("ownerid is invalid"), 400, "OwnerID invalid", "OwnerId is invalid for authentication"))
		return
	}

	ownerMongoDBID, err := primitive.ObjectIDFromHex(ownerIDStr)
	if err != nil {
		_ = utils.WriteResponse(res, 400, utils.ReturnAppError(errors.New("ownerid is invalid"), 400, "OwnerID invalid", "OwnerId is invalid for authentication"))
		return
	}

	//now extract clinicdetails from req.payload
	var clinicDetails models.ClinicDoctor
	if err := json.NewDecoder(req.Body).Decode(&clinicDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "InValid JsonDetails", "Invalid Json Details"))
		return
	}

	//now validate details first
	validationErrors := controller.ValidateAddDoctorToclinicDetailsFn(&clinicDetails)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(validationErrors, http.StatusBadRequest, "Invalid Details", "validation Failed"))
		return
	}

	//now as we have clinic also now fit details such as clinicid in clinic details and pass this info to service layer
	if err := controller.Service.AddDoctorToclinic(ctx, ownerMongoDBID, clinicDetails); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "OTP sent to doctor email valid for " + utils.OTPExpiry.String() + " seconds",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func (controller *Controller) VerifyAddDoctorToclinicOtp(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var otpPayload struct {
		OTP string `json:"otp"`
	}

	if err := json.NewDecoder(req.Body).Decode(&otpPayload); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid Json Details", "Invalid details"))
		return
	}

	//extract the owner id or receptionist id whoever is authorized for adding doctor ,from context
	ownerIDAny := req.Context().Value(constants.KeyUserID)

	mongodbIDError, ownerMongoDBID := utils.ParseUserID(ownerIDAny)
	if mongodbIDError != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid UserId",
			Reason:     mongodbIDError.Error(),
			ErrorObj:   mongodbIDError,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now do some validation
	if len(otpPayload.OTP) != 6 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("6 digit  otp required"))
		return
	}

	//now pass these details to service layer
	if err := controller.Service.VerifyAddDoctorToclinicOTP(ctx, otpPayload.OTP, ownerMongoDBID); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &structs.IAppSuccess{
		Message:    "Doctor Onboarded Successfully",
		Data:       nil,
		StatusCode: http.StatusOK,
	})

}

func (controller *Controller) FetchDoctorWithItsclinics(res http.ResponseWriter, req *http.Request) {
	//try to parse query params
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	params := req.URL.Query()

	filter := bson.M{}

	_ = utils.TransformParamIDS(params, filter)

	fmt.Print("filter is", filter)

	doctorWithItsclinics, err := controller.Service.FetchDoctorClinicMappings(ctx, filter)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Successfully fetched Doctor Details",
		Data:       doctorWithItsclinics,
		StatusCode: http.StatusOK,
	})

}

func (controller *Controller) GetDoctors(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params := req.URL.Query()

	filters := bson.M{}

	//transform param filters in mongodb format
	_ = utils.TransformParamIDS(params, filters)

	doctors, err := controller.Service.FetchDoctors(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &structs.IAppSuccess{
		Message:    "Successfully Fetched Doctors",
		Data:       doctors,
		StatusCode: http.StatusOK,
	})
}

func (controller *Controller) RegisterDoctor(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var doctorDetails models.Doctor
	if jsonErr := json.NewDecoder(req.Body).Decode(&doctorDetails); jsonErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Json Details",
			Reason:     jsonErr.Error(),
			ErrorObj:   jsonErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now do some validations here

	//extract doctorID here from req.context
	doctorIDAny := req.Context().Value(constants.KeyUserID)
	doctorMongoDBErr, doctorMongoDBID := utils.ParseUserID(doctorIDAny)
	if doctorMongoDBErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid doctorID",
			Reason:     doctorMongoDBErr.Error(),
			ErrorObj:   doctorMongoDBErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	//now add this doctorMongoDBID to doctor details
	doctorDetails.ID = doctorMongoDBID

	if registrationErr := controller.Service.RegisterDoctor(ctx, doctorDetails); registrationErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, registrationErr)
		return
	}

}
