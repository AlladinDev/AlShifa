// Package controller implements various handlers for user module
package controller

import (
	middleware "AlShifa/Middleware"
	structs "AlShifa/Structs"
	validators "AlShifa/Users/Validators"

	interfaces "AlShifa/Users/Interfaces"
	models "AlShifa/Users/Models"
	userModuleStructs "AlShifa/Users/Structs"
	utils "AlShifa/Utils"

	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()
	var user models.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Registration Failed",
			StatusCode: 400,

			Reason: "Invalid Json Details",
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
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	userID := req.Context().Value(middleware.ContextUserIDKey)
	if userID == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(nil, 400, "Missing ID", "Missing ID"))
		return
	}
	// //now safely try to parse userID which is of type any to string and if it fails throw error
	// userIDString,
	userIDString, ok := userID.(string)
	if !ok || userIDString == "" {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(nil, 400, "Invalid or missing Userid ", "Invalid or missing Userid"))
		return
	}

	objectUserID, err := primitive.ObjectIDFromHex(userIDString)
	if err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 400, "Invalid UserID", err.Error()))
		return
	}

	//if userId is nill or zero typed return error
	if objectUserID == primitive.NilObjectID {
		_ = utils.WriteResponse(res, http.StatusBadRequest, utils.ReturnAppError(err, 400, "Zero valued ObjectId", "Zero value userid passed"))
		return
	}

	user, searchErr := controller.Service.SearchUserByID(ctx, objectUserID)
	if searchErr != nil {
		_ = utils.WriteResponse(res, searchErr.StatusCode, utils.ReturnAppError(err, searchErr.StatusCode, "Failed To Get User Details", "Server Error"))
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, user)
}

func (controller *UserController) LoginUser(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), utils.RequestTimeout)
	defer cancel()

	var loginDetails userModuleStructs.LoginDetails
	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Login Failed",
			StatusCode: http.StatusBadRequest,
			Reason:     "Invalid Json",
			ErrorObj:   err,
		})
		return
	}

	//now do some validation
	validationErrs := validators.ValidateLoginDetails(&loginDetails)
	if validationErrs != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Details",
			StatusCode: 400,
			Reason:     "Invalid Details",
			ErrorObj:   validationErrs,
		})
		return
	}

	jwtToken, err := controller.Service.LoginUser(ctx, loginDetails.Email, loginDetails.Password)
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
