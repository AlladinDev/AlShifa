// Package controller implements various handlers for user module
package controller

import (
	"AlShifa/constants"
	structs "AlShifa/structs"
	interfaces "AlShifa/users/interfaces"
	models "AlShifa/users/models"
	validators "AlShifa/users/validators"
	utils "AlShifa/utils"

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

func (controller *UserController) RegisterUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var user models.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Registration Failed",
			StatusCode: 400,
			ErrorObj:   err,
			Reason:     "Invalid Json Details",
		})
		return
	}

	//validate user details also
	validationErrors := validators.ValidateUser(&user)
	if validationErrors != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Registration Failed",
			StatusCode: 400,
			Reason:     "Invalid Details",
			ErrorObj:   validationErrors,
		})
		return
	}

	if err := controller.Service.AddUser(ctx, user); err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, err)
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, structs.IAppSuccess{
		Message:    "User Registered Successfully",
		Data:       nil,
		StatusCode: 201,
	})
}

func (controller *UserController) SearchUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userIDAny := req.Context().Value(constants.KeyUserID)
	userMongoDBIDErr, userMongoDBID := utils.ParseUserID(userIDAny)
	if userMongoDBIDErr != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid UserID",
			Reason:     userMongoDBIDErr.Error(),
			ErrorObj:   userMongoDBIDErr,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	user, searchErr := controller.Service.SearchUserByID(ctx, userMongoDBID)
	if searchErr != nil {
		_ = utils.WriteResponse(res, searchErr.StatusCode, searchErr)
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, user)
}
