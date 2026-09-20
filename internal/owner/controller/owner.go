// Package controller provides handler methods for owner module
package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/owner/dtos"
	"github.com/AlladinDev/AlShifa/internal/owner/interfaces"
	"github.com/AlladinDev/AlShifa/internal/owner/models"
	"github.com/AlladinDev/AlShifa/internal/owner/validators"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
)

type Controller struct {
	service interfaces.IService
}

func NewController(service interfaces.IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) RegisterOwner(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	parsingErr := req.ParseMultipartForm(32 << 20) // 32 MB
	if parsingErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "File Size Uploaded Exceeds Max Limit",
			Reason:     "File Limit exceeds",
			ErrorObj:   parsingErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// 2. Extract the file using the HTML input field name (e.g., "photo")
	// "file" is the data stream, "header" contains the filename and size.
	profilePhotoFile, _, err := req.FormFile("photo")
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Profile Photo Missing ",
			Reason:     "Invalid Details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	defer profilePhotoFile.Close()

	var ownerDetails models.Owner
	//it will extract the form values and copy them into ownerDetails model form
	utils.UnmarshalFormValues(req.Form, &ownerDetails)

	if err := validators.ValidateOwnerDetails(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Registration Failed",
			Reason:     "Invalid Details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	authToken, registrationErr := c.service.RegisterOwner(ctx, ownerDetails, profilePhotoFile)
	if registrationErr != nil {
		_ = utils.WriteResponse(res, registrationErr.StatusCode, registrationErr)
		return
	}

	//here set the token in cookie
	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameOTPToken,
		Value:    authToken,
		Expires:  time.Now().Add(constants.OTPExpiry),
		Secure:   true,
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusCreated, response.IAppSuccess{
		Message:    "Now Verify Your Email",
		StatusCode: http.StatusCreated,
		Data:       nil,
	})

}

func (c *Controller) VerifyOwnerRegistrationOTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var otp string
	if err := json.NewDecoder(req.Body).Decode(&otp); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid Json Details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now check for cookie it should have the cookie
	cookie, err := req.Cookie(constants.NameOTPToken)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
				Message:    "Verification failed try again submitting profile details",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Verification Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now decode the value in cookie
	cookieValueDecoded, decodingErr := utils.Decrypt(cookie.Value)
	if decodingErr != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, response.IAppError{
			Message:    "Registration Failed",
			Reason:     decodingErr.Error(),
			ErrorObj:   decodingErr,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	//now parse cookie decoded into suitable form otpPayload form
	var otpPayload dtos.OTPPayload
	if err := json.Unmarshal(cookieValueDecoded, &otpPayload); err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, response.IAppError{
			Message:    constants.NameOTPToken,
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	///here it otpPayload will have expiry time check it
	if otpPayload.ExpiresAt.Before(time.Now()) {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, response.IAppError{
			Message:    "Profile Creation timed out plz create your profile again",
			Reason:     "Profile Creation Timed out",
			ErrorObj:   nil,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	//now call service function to verify it
	jwtToken, verificationErr := c.service.VerifyOwnerRegisterOtp(ctx, otp, otpPayload)
	if verificationErr != nil {
		_ = utils.WriteResponse(res, verificationErr.StatusCode, verificationErr)
		return
	}

	//set the cookie here
	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameAuthToken,
		Value:    fmt.Sprintf("BEARER %s", jwtToken),
		Expires:  time.Now().Add(constants.JwtExpiryTime),
		Secure:   true,
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusCreated, response.IAppSuccess{
		Message:    "Owner Registered Successfully",
		Data:       nil,
		StatusCode: http.StatusCreated,
	})

}
func (c *Controller) GetOwnerByID(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ownerID := ctx.Value(constants.KeyUserID)
	ownerMongoDBID, ownerIDErr := utils.ParseUserID(ownerID)
	if ownerIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid OwnerID",
			Reason:     ownerIDErr.Error(),
			ErrorObj:   ownerIDErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	owner, err := c.service.GetOwnerByID(ctx, ownerMongoDBID)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, response.IAppError{
			Message:    "Failed to fetch owner details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: err.StatusCode,
		})
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, response.IAppSuccess{
		Message:    "Owner Details Fetched Successfully",
		Data:       owner,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) InitiateLogin(res http.ResponseWriter, req *http.Request) {
	//we have to get email and password
	loginDto := dtos.LoginDTO{}
	if err := json.NewDecoder(req.Body).Decode(&loginDto); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid json details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now do some verification
	validationErrors := map[string]string{}
	if len(loginDto.Email) > constants.MaxEmailLength {
		validationErrors["email"] = "email too long"
	}
	if len(loginDto.Password) < 8 {
		validationErrors["password"] = "password cannot be less than 8"
	}

	if len(validationErrors) > 0 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "validation errors",
			Reason:     "validation erros",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   validationErrors,
		})
		return
	}

	temporaryToken, loginErr := c.service.InitiateLogin(req.Context(), loginDto.Email, loginDto.Password)
	if loginErr != nil {
		_ = utils.WriteResponse(res, loginErr.StatusCode, loginErr)
		return
	}

	//now set this temporary token into cookie
	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameOTPToken,
		Value:    temporaryToken,
		Expires:  time.Now().Add(constants.OTPExpiry),
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Verify your email now",
		Data:       nil,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) VerifyLogin(res http.ResponseWriter, req *http.Request) {
	//here get the token from cookie
	authToken, cookieErr := req.Cookie(constants.NameOTPToken)
	if cookieErr != nil {
		if errors.Is(cookieErr, http.ErrNoCookie) {
			_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
				Message:    "Auth token is missing",
				Reason:     "cookie is missing",
				ErrorObj:   nil,
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     cookieErr.Error(),
			ErrorObj:   cookieErr,
			StatusCode: http.StatusInternalServerError,
		})
	}

	//now get the otp value here
	var otpProvided string
	if err := json.NewDecoder(req.Body).Decode(&otpProvided); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid otp failed to parse it",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	jwtToken, err := c.service.VerifyLogin(req.Context(), authToken.Value, otpProvided)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	//now here set this jwttoken into cookie
	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameAuthToken,
		Value:    jwtToken,
		Expires:  time.Now().Add(constants.JwtExpiryTime),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "logged in successfully",
		Data:       nil,
		StatusCode: http.StatusOK,
	})
}
