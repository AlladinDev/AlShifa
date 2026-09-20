// Package controller implements various handlers for user module
package controller

import (
	"time"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
	interfaces "github.com/AlladinDev/AlShifa/internal/users/interfaces"

	"encoding/json"
	"net/http"
)

type UserController struct {
	Service interfaces.IService
}

func ReturnNewController(service interfaces.IService) *UserController {
	return &UserController{
		Service: service,
	}
}

func (controller *UserController) InitiateLogin(res http.ResponseWriter, req *http.Request) {
	var mobileNumber string
	ctx := req.Context()
	if err := json.NewDecoder(req.Body).Decode(&mobileNumber); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid json Details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	token, err := controller.Service.InitiateLoginOTP(ctx, mobileNumber)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameOTPToken,
		Value:    token,
		Expires:  time.Now().Add(constants.OTPExpiry),
		Secure:   true,
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusOK, response.IAppSuccess{
		Message:    "Login otp sent now verify your mobile number",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func (controller *UserController) VerifyLogin(res http.ResponseWriter, req *http.Request) {
	//now here try to get the token from cookie
	otpCookie, err := req.Cookie(constants.NameOTPToken)
	if err != nil {
		if err != http.ErrNoCookie {
			_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
				Message:    "No registration token found",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		_ = utils.WriteResponse(res, http.StatusInternalServerError, &response.IAppError{
			Message:    "mobile verification failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	var otp string
	if err := json.NewDecoder(req.Body).Decode(&otp); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "invalid otp provided",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
			Reason:     err.Error(),
		})
	}

	authToken, verificationErr := controller.Service.VerifyLoginOTP(req.Context(), otpCookie.Value, otp)
	if verificationErr != nil {
		_ = utils.WriteResponse(res, verificationErr.StatusCode, verificationErr)
	}

	//now set this authtoken in cookie
	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameAuthToken,
		Value:    authToken,
		Expires:  time.Now().Add(constants.JwtExpiryTime),
		Secure:   true,
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "mobile verified successfully now access our services ",
		Data:       nil,
		StatusCode: http.StatusOK,
	})

}
