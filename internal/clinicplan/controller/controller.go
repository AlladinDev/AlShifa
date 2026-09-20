package plancontroller

import (
	"encoding/json"
	"net/http"

	planInterface "github.com/AlladinDev/AlShifa/internal/clinicplan/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service planInterface.IService
}

func NewController(Service planInterface.IService) *Controller {
	return &Controller{
		service: Service,
	}
}

func (c *Controller) AddNewPlan(res http.ResponseWriter, req *http.Request) {
	var planDetails models.Plan

	if err := json.NewDecoder(req.Body).Decode(&planDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, response.IAppError{
			Message:    "Invalid Json Details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now here do some validation
	if err := c.service.RegisterPlan(req.Context(), &planDetails); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response.IAppSuccess{
		Message:    "Plan Added Successfully",
		Data:       nil,
		StatusCode: http.StatusCreated,
	})

}

func (c *Controller) FetchPlans(res http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()

	filters := bson.M{}

	_ = utils.TransformParamIDS(urlParams, filters)

	plans, err := c.service.FetchPlans(req.Context(), filters)
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, response.IAppSuccess{
		Message:    "Successfully Fetched Plans",
		Data:       plans,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) UpdatePlan(res http.ResponseWriter, req *http.Request) {
	var dataToUpdate models.Plan
	if err := json.NewDecoder(req.Body).Decode(&dataToUpdate); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid Json details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
			Reason:     err.Error(),
		})
		return
	}

	//now extract the id from req.url param
	planIDString := req.URL.Query().Get("id")
	if planIDString == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "plan id is required",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
			Reason:     "plan id is not present",
		})
		return
	}

	//now convert this id into mongodb format
	planID, err := primitive.ObjectIDFromHex(planIDString)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, &response.IAppError{
			Message:    "Invalid Plan ID ",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
			Reason:     err.Error(),
		})
		return
	}

	if err := c.service.UpdatePlan(req.Context(), planID, &dataToUpdate); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, response.IAppSuccess{
		Message:    "Successfully updated plan",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}
