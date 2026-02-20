// Package controller provides HTTP handlers for managing clinic-related operations.
package controller

import (
	interfaces "AlShifa/clinic/Interfaces"
	validators "AlShifa/clinic/Validators"
	"AlShifa/clinic/models"
	middleware "AlShifa/middleware"
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
	ValidateAddDoctorToclinicDetailsFn func(details *models.AddDoctorToclinic) map[string]string
}

type clinicRegistration struct {
	OwnerID primitive.ObjectID `json:"ownerId"`
	clinic  models.Clinic      `json:"clinicDetails"`
}

func NewController(svr interfaces.IService, validateAddDoctorToclinicFn func(details *models.AddDoctorToclinic) map[string]string) *Controller {
	return &Controller{
		Service:                            svr,
		ValidateAddDoctorToclinicDetailsFn: validateAddDoctorToclinicFn,
	}
}

func (controller *Controller) Registerclinic(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var clinicRegistrationDetails clinicRegistration
	if err := json.NewDecoder(req.Body).Decode(&clinicRegistrationDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid Details Provided", "Json Error"))
		return
	}

	//here validate clinic details
	validationErrors := validators.ValidateclinicDetails(&clinicRegistrationDetails.clinic)
	if len(validationErrors) != 0 {
		fmt.Print(validationErrors)
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(validationErrors, 400, "Invalid Details", "Validation Failed"))
		return
	}

	//here extract the ownerId from req.context fed by jwt middleware
	ownerID := req.Context().Value(middleware.ContextUserIDKey)

	//try to convert it into string
	ownerIDStr, ok := ownerID.(string)
	if !ok {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(errors.New("invalid ownerid"), 400, "Invalid OwnerId String Failed to register clinic", "Invalid ownerid"))
		return
	}

	if len(ownerIDStr) == 0 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(errors.New("owner Id missing"), 400, "Ownerid is missing", "Missing OwnerId"))
		return
	}

	//now try to convert this ownerIDStr into mongodbID
	ownerMongoDBID, err := primitive.ObjectIDFromHex(ownerIDStr)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(errors.New("invalid ownerid"), 400, "Invalid OwnerId Failed to register clinic", "Invalid ownerid"))
		return
	}

	//here add this id to clinicregistration details so that user can send any other owners id
	clinicRegistrationDetails.OwnerID = ownerMongoDBID

	registrationErr := controller.Service.Registerclinic(ctx, clinicRegistrationDetails.OwnerID, clinicRegistrationDetails.clinic)
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

	if err := controller.Service.RegisterclinicOwner(ctx, ownerDetails); err != nil {
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

	_ = utils.TransformParamIDS(params, filters)

	doctors, err := controller.Service.SearchDoctor(ctx, filters)
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

func (controller *Controller) LoginclinicOwner(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var loginDetails structs.LoginDetails

	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 500, "Login Failed", "Invalid Json"))
		return
	}
	jwtToken, err := controller.Service.LoginclinicOwner(ctx, loginDetails.Email, loginDetails.Password)
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

func (controller *Controller) AddDoctorToclinic(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	//first try to see if client id is present in req.context
	ownerID := req.Context().Value(middleware.ContextUserIDKey)
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
	fmt.Print("reached here")

	//now extract clinicdetails from req.payload
	var clinicDetails models.AddDoctorToclinic
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

	//now using ownerMongoDbID find its associated clinic
	owners, searchclinicErr := controller.Service.SearchOwner(ctx, bson.M{"_id": ownerMongoDBID})
	if searchclinicErr != nil {
		_ = utils.WriteResponse(res, searchclinicErr.StatusCode, searchclinicErr)
		return
	}

	if len(owners) == 0 {
		_ = utils.WriteResponse(res, http.StatusNotFound, &structs.IAppError{
			Message:    "No Owner Found",
			Reason:     "No Owner Found For this OwnerID",
			StatusCode: http.StatusNotFound,
			ErrorObj:   nil,
		})
		return
	}

	//add id to clinicDetails and dont trust frontend for sending it
	clinicDetails.ClinicID = owners[0].Clinic

	//now as we have clinic also now fit details such as clinicid in clinic details and pass this info to service layer
	if err := controller.Service.AddDoctorToclinic(ctx, clinicDetails); err != nil {
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
		OTP      string `json:"otp"`
		DoctorID string `json:"doctorID"`
	}

	if err := json.NewDecoder(req.Body).Decode(&otpPayload); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid Json Details", "Invalid details"))
		return
	}

	//now try to extract ownerid from context
	ownerID := req.Context().Value(middleware.ContextUserIDKey)
	if ownerID == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("missing clinicid"))
		return
	}

	//now try to convert clinicid to string
	ownerIDString, ok := ownerID.(string)
	if !ok {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("invalid clinicId"))
		return
	}

	//now convert clinicIDString into mongodb format
	ownerIDMongoDBID, err := primitive.ObjectIDFromHex(ownerIDString)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("invalid clinicId"))
		return
	}

	//now do some validation
	if len(otpPayload.OTP) != 6 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("6 digit  otp required"))
		return
	}

	if otpPayload.DoctorID == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("DoctorID cannot be empty"))
		return
	}

	//now convert doctorID to  mongoDB ID
	doctorMongoDBID, err := primitive.ObjectIDFromHex(otpPayload.DoctorID)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, errors.New("invalid doctorID"))
		return
	}

	//now fetch the clinic using ownerID
	owners, ownerSearchErr := controller.Service.SearchOwner(ctx, bson.M{"_id": ownerIDMongoDBID})
	if ownerSearchErr != nil {
		_ = utils.WriteResponse(res, ownerSearchErr.StatusCode, err)
		return
	}

	if len(owners) == 0 {
		_ = utils.WriteResponse(res, http.StatusNotFound, &structs.IAppError{
			Message:    "No Owner Found",
			ErrorObj:   errors.New("no Owner found"),
			Reason:     errors.New("no Owner found").Error(),
			StatusCode: http.StatusNotFound,
		})
		return
	}
	clinicID := owners[0].Clinic

	//now pass these details to service layer
	if err := controller.Service.VerifyAddDoctorToclinicOTP(ctx, otpPayload.OTP, doctorMongoDBID, clinicID); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &structs.IAppSuccess{
		Message:    "Doctor Onboarded Successfully",
		Data:       nil,
		StatusCode: http.StatusOK,
	})

}

func (controller *Controller) AddAppointment(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()
	//first parse the details
	var appointmentDetails models.Appointment
	if err := json.NewDecoder(req.Body).Decode(&appointmentDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid json", err.Error()))
		return
	}

	//now do the validations
	if err := validators.ValidateAppointmentDetails(&appointmentDetails, req.Context()); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, http.StatusBadRequest, "Invalid json", "Invalid details"))
		return
	}

	//now call the service layer method
	slotBooked, err := controller.Service.AddAppointment(ctx, appointmentDetails)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, utils.ReturnAppSuccess(http.StatusCreated, "Appointment Booked", slotBooked))

}

func (controller *Controller) AppointmentSlotsBooked(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()
	var slotDetails models.SlotDetails
	if err := json.NewDecoder(req.Body).Decode(&slotDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(errors.New("invalid Json Details"), http.StatusBadRequest, "Invalid Json Details", "Invalid Json Details"))
		return
	}

	validationErrors := validators.ValidateSlotDetails(&slotDetails)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, validationErrors)
		return
	}

	slots, err := controller.Service.AppointmentSlotsBooked(ctx, slotDetails)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, utils.ReturnAppSuccess(http.StatusOK, "Successfully Fetched Slots booked", slots))

}

func (controller *Controller) FetchDoctorWithItsclinics(res http.ResponseWriter, req *http.Request) {
	//try to parse query params
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	params := req.URL.Query()

	filter := bson.M{}

	_ = utils.TransformParamIDS(params, filter)

	fmt.Print("filter is", filter)

	doctorWithItsclinics, err := controller.Service.DoctorWithItsclinics(ctx, filter)
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

func (controller *Controller) FetchClinicAppointments(res http.ResponseWriter, req *http.Request) {

}
