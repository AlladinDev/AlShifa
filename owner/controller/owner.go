// Package controller provides handler methods for owner module
package controller

import (
	"encoding/json"
	"net/http"

	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/owner/interfaces"
	"github.com/AlladinDev/AlShifa/owner/models"
	"github.com/AlladinDev/AlShifa/owner/validators"
	"github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/utils"

	"go.mongodb.org/mongo-driver/bson"
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

	var ownerDetails models.Owner
	if err := json.NewDecoder(req.Body).Decode(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := validators.ValidateOwnerDetails(&ownerDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     "Invalid Details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := c.service.RegisterOwner(ctx, ownerDetails); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     err.Reason,
			ErrorObj:   err,
			StatusCode: err.StatusCode,
		})
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, structs.IAppSuccess{
		Message:    "Owner Registered Successfully",
		StatusCode: http.StatusCreated,
		Data:       nil,
	})

}

func (c *Controller) GetOwnerByID(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ownerID := ctx.Value(constants.KeyUserID)
	ownerIDErr, ownerMongoDBID := utils.ParseUserID(ownerID)
	if ownerIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid OwnerID",
			Reason:     ownerIDErr.Error(),
			ErrorObj:   ownerIDErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	owner, err := c.service.GetOwnerByID(ctx, ownerMongoDBID)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, structs.IAppError{
			Message:    "Failed to fetch owner details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: err.StatusCode,
		})
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Failed to fetch owner details",
		Data:       owner,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) GetOwnerDetails(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params := req.URL.Query()
	if len(params) == 0 {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Filters Required For ownerDetails",
			Reason:     "Filters not provided",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	filters := bson.M{}

	utils.TransformParamIDS(params, filters)

	owner, err := c.service.GetOwnerDetails(ctx, filters)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, structs.IAppError{
			Message:    err.Message,
			Reason:     err.Reason,
			ErrorObj:   err,
			StatusCode: err.StatusCode,
		})
		return
	}
	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Succesfully Fetched Owner Details",
		StatusCode: http.StatusOK,
		Data:       owner,
	})

}
