// Package cliniccontroller provides HTTP handlers for managing clinic-related operations.
package cliniccontroller

import (
	"encoding/json"
	"fmt"

	"net/http"

	clinicinterfaces "github.com/AlladinDev/AlShifa/internal/clinic/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinic/models"
	validators "github.com/AlladinDev/AlShifa/internal/clinic/validators"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type Controller struct {
	Service clinicinterfaces.IService
}

func NewController(svr clinicinterfaces.IService) *Controller {
	return &Controller{
		Service: svr,
	}
}

func (controller *Controller) Registerclinic(res http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

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

	ownerMongoDBID, ownerMongodbIDErr := utils.ParseUserID(ownerID)
	if ownerMongodbIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid OwnerID",
			Reason:     ownerMongodbIDErr.Error(),
			StatusCode: http.StatusBadRequest,
			ErrorObj:   ownerMongodbIDErr,
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

	response := response.IAppSuccess{
		Message:    "clinic Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response)
}

func (controller *Controller) Searchclinic(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

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
