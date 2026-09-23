package controller

import (
	"net/http"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Controller struct {
	sv interfaces.IService
}

func NewController(service interfaces.IService) *Controller {
	return &Controller{
		sv: service,
	}
}

func (c *Controller) FetchClinicsWithDoctors(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	filters := bson.M{}

	_ = utils.TransformParamIDS(params, filters)
	data, err := c.sv.FetchClinicWithDoctors(req.Context(), filters)

	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Successfully Fetched Clinics With Doctors",
		Data:       data,
		StatusCode: http.StatusOK,
	})
}

func (c *Controller) FetchDoctorWithClinics(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	filters := bson.M{}

	_ = utils.TransformParamIDS(params, filters)
	data, err := c.sv.FetchDoctorWithClinics(req.Context(), filters)

	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, &response.IAppSuccess{
		Message:    "Successfully Fetched Doctors With Clinics",
		Data:       data,
		StatusCode: http.StatusOK,
	})
}
