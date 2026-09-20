// Package doctorcontroller provides controller for alshifa
package doctorcontroller

import (
	"encoding/json"
	"net/http"
	"time"

	doctordtos "github.com/AlladinDev/AlShifa/internal/doctor/dtos"
	doctorinterfaces "github.com/AlladinDev/AlShifa/internal/doctor/interfaces"
	doctormodels "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"github.com/AlladinDev/AlShifa/internal/doctor/validators"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service           doctorinterfaces.IService
	parseUserIDToDBID func(id any) (primitive.ObjectID, error)
}

func NewController(Service doctorinterfaces.IService, userIDParsingFn func(id any) (primitive.ObjectID, error)) *Controller {
	return &Controller{
		service:           Service,
		parseUserIDToDBID: userIDParsingFn,
	}
}

func (c *Controller) RegisterDoctor(res http.ResponseWriter, req *http.Request) {
	var details doctormodels.Doctor

	//load the form data
	if err := req.ParseMultipartForm(20 << 30); err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, &response.IAppError{
			Message:    "Failed to register doctor",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		})
		return
	}

	photoField, photoHeader, err := req.FormFile("photo")
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, &response.IAppError{
			Message:    "Failed to register doctor",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		})
		return
	}

	if photoField == nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Doctor profile photo is required",
			StatusCode: http.StatusBadRequest,
			Reason:     "photo is required",
			ErrorObj:   nil,
		})
		return
	}

	//now get certificates
	certificates := req.MultipartForm.File["certificates"]
	if len(certificates) == 0 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Atleast one specialization certificate is required eg MBBS Certificate",
			StatusCode: http.StatusBadRequest,
			Reason:     "specialization certificate required",
			ErrorObj:   nil,
		})
		return
	}

	//now here get the other details using this util function
	utils.UnmarshalFormValues(req.Form, &details)

	//now here do some validations
	validationErrors := validators.ValidateDoctorDetails(&details)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "validation failed",
			StatusCode: http.StatusBadRequest,
			Reason:     "invalid registration details",
			ErrorObj:   validationErrors,
		})
		return
	}

	successRes, registrationErr := c.service.RegisterDoctor(req.Context(), &details, photoHeader, certificates)
	if registrationErr != nil {
		_ = utils.WriteResponse(res, registrationErr.StatusCode, registrationErr)
		return
	}

	//here set the cookie with otp payload data
	http.SetCookie(res, &http.Cookie{
		Name:    constants.NameOTPToken,
		Value:   successRes,
		Expires: time.Now().Add(constants.OTPExpiry),
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Otp sent to email",
		Data:       nil,
		StatusCode: http.StatusInternalServerError,
	})

}

func (c *Controller) VerifyDoctorRegistrationOtp(res http.ResponseWriter, req *http.Request) {
	authCookie, err := req.Cookie(constants.NameOTPToken)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Cookie not found  Create profile again",
			Reason:     err.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
		})
		return
	}

	//now extract here the otp from req.params
	otpProvided := ""
	if err := json.NewDecoder(req.Body).Decode(&otpProvided); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "verification failed ,failed to parse otp",
			Reason:     err.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
		})
		return
	}

	//now auth cookie is encrypted so decode it first
	decodedCookie, err := utils.Decrypt(authCookie.Value)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Cookie not found  Create profile again",
			Reason:     err.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
		})
		return
	}

	///now typecast decodecookie into otppayload form
	var otpPayload doctordtos.TOTPPayload
	if err := json.Unmarshal(decodedCookie, &otpPayload); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Verification failed",
			Reason:     err.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
		})
		return
	}

	jwtToken, verificationErr := c.service.VerifyDoctorRegistrationOtp(req.Context(), otpProvided, otpPayload)
	if verificationErr != nil {
		_ = utils.WriteResponse(res, verificationErr.StatusCode, verificationErr)
		return
	}

	///now here set jwt in req.Cookie
	http.SetCookie(res, &http.Cookie{
		Name:     "authToken",
		Value:    jwtToken,
		Expires:  time.Now().Add(constants.JwtExpiryTime),
		HttpOnly: true,
		Secure:   true,
	})

	//now send the response
	_ = utils.WriteResponse(res, http.StatusCreated, &response.IAppSuccess{
		Message:    "Registered Successfully",
		Data:       nil,
		StatusCode: http.StatusCreated,
	})
}

func (c *Controller) GetDoctorProfile(res http.ResponseWriter, req *http.Request) {
	//first extract the userid from context
	userIDAny := req.Context().Value(constants.KeyUserID)

	userMongodbID, err := utils.ParseUserID(userIDAny)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Reason:     err.Error(),
			ErrorObj:   err,
		})
		return
	}

	doctor, doctorErr := c.service.GetDoctorDetails(req.Context(), userMongodbID)
	if doctorErr != nil {
		_ = utils.WriteResponse(res, doctorErr.StatusCode, doctorErr)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Doctor fetched successfully",
		Data:       doctor,
		StatusCode: http.StatusOK,
	})

}

func (c *Controller) InitiateLogin(res http.ResponseWriter, req *http.Request) {
	loginDetails := doctordtos.LoginDTO{}

	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid login details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now do some validations here

	otpToken, tokenErr := c.service.InitiateLogin(req.Context(), loginDetails)
	if tokenErr != nil {
		_ = utils.WriteResponse(res, tokenErr.StatusCode, tokenErr)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameOTPToken,
		Value:    otpToken,
		Expires:  time.Now().Add(constants.OTPExpiry),
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Now verify your email",
		Data:       nil,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) VerifyLogin(res http.ResponseWriter, req *http.Request) {
	//get the token from cookie
	otpToken, cookieErr := req.Cookie(constants.NameOTPToken)
	if cookieErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Verification failed",
			Reason:     cookieErr.Error(),
			ErrorObj:   cookieErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now get the otp from payload
	otpShared := ""
	if err := json.NewDecoder(req.Body).Decode(&otpShared); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Otp not provided",
			Reason:     "provide otp to verify email",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//here check cookie should not be empty
	if otpToken.Value == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Verification failed",
			Reason:     "otp cookie cannot be empty",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	jwtToken, err := c.service.VerifyLogin(req.Context(), otpToken.Value, otpShared)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     constants.NameAuthToken,
		Value:    jwtToken,
		Expires:  time.Now().Add(constants.JwtExpiryTime),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Logged in successfully",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func (c *Controller) AllowClinicToOnboard(res http.ResponseWriter, req *http.Request) {
	//from req.context get the doctor id first
	doctorIDAny := req.Context().Value(constants.KeyUserID)

	doctorMongodbID, idErr := utils.ParseUserID(doctorIDAny)
	if idErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid doctor id",
			Reason:     idErr.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now get the clinicid from request body
	clinicIDShared := primitive.ObjectID{}
	if err := json.NewDecoder(req.Body).Decode(&clinicIDShared); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid clinic id",
			Reason:     err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	respErr := c.service.AllowClinicToOnboardDoctor(req.Context(), doctorMongodbID, clinicIDShared)
	if respErr != nil {
		_ = utils.WriteResponse(res, respErr.StatusCode, respErr)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, &response.IAppSuccess{
		Message:    "Clinic allowed to onboard now",
		Data:       nil,
		StatusCode: http.StatusCreated,
	})

}
