package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/validators"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
)

type onboardingController struct {
	service interfaces.IDoctorOnboarding
}

func NewOnboardingController(Service interfaces.IDoctorOnboarding) *onboardingController {
	return &onboardingController{
		service: Service,
	}
}

func (c *onboardingController) InitiateDoctorClinicOnboarding(res http.ResponseWriter, req *http.Request) {
	///first get the clinic owner id from req.context
	ownerIDAny := req.Context().Value(constants.KeyUserID)
	ownerMongodbID, IDErr := utils.ParseUserID(ownerIDAny)
	if IDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid owner id provided",
			Reason:     IDErr.Error(),
			ErrorObj:   IDErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//get the data first
	var data models.ClinicDoctorMapping
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid json Details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now here do some validations
	validationErrors := validators.ValidateClinicDoctorMapping(&data)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid details validation failed",
			ErrorObj:   validationErrors,
			StatusCode: http.StatusBadRequest,
			Reason:     "invalid details",
		})
		return
	}

	token, err := c.service.InitiateDoctorClinicOnboarding(req.Context(), &data, ownerMongodbID)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	//here set the token in cookies and set its name and expiry from constants package
	http.SetCookie(res, &http.Cookie{
		Name:     "onboarding_token",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(constants.OTPExpiry),
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "onboarding initiated ,otp has been sent to doctor email",
		Data:       nil,
		StatusCode: http.StatusOK,
	})
}
