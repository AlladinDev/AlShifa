// Package controller provides controller functions for auth module
package controller

import (
	"AlShifa/auth/dtos"
	"AlShifa/auth/interfaces"
	"AlShifa/auth/models"
	validation "AlShifa/auth/validators"
	"AlShifa/constants"
	"AlShifa/structs"
	"AlShifa/utils"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type Controller struct {
	sv interfaces.IService
}

func NewController(sv interfaces.IService) *Controller {
	return &Controller{
		sv: sv,
	}
}

func (c *Controller) Register(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var credientials models.Credientials
	if err := json.NewDecoder(req.Body).Decode(&credientials); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "InValid JsonDetails",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now do some validation
	if err := validation.ValidateCredentials(&credientials); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Details",
			Reason:     "Invalid details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := c.sv.Register(ctx, credientials); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	_ = utils.WriteResponse(res, http.StatusCreated, structs.IAppSuccess{
		Message:    "Credientials Saved Successfully",
		Data:       nil,
		StatusCode: http.StatusCreated,
	})
}

func (c *Controller) Login(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var loginDetails dtos.LoginDetails
	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid Details",
			Reason:     "Invalid details",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now check whether this credientials exists or not
	credientials, err := c.sv.SearchCredientials(ctx, bson.M{"email": loginDetails.Email})
	if err != nil {
		_ = utils.WriteResponse(res, err.StatusCode, structs.IAppError{
			Message:    "Login Failed",
			Reason:     err.Reason,
			ErrorObj:   err,
			StatusCode: err.StatusCode,
		})
		return
	}

	//now match password
	passwordMatches, passwordMatchingErr := utils.VerifyPasswordArgon2id(loginDetails.Password, credientials.Password)
	if passwordMatchingErr != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, structs.IAppError{
			Message:    "Login Failed",
			Reason:     passwordMatchingErr.Error(),
			ErrorObj:   passwordMatchingErr,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if !passwordMatches {
		_ = utils.WriteResponse(res, http.StatusBadRequest, structs.IAppError{
			Message:    "Invalid email or password",
			Reason:     "Invalid email or password",
			ErrorObj:   "Invalid email or password",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	//now generate jwt
	token, tokenErr := utils.GenerateJWT(&constants.JwtCustomClaims{
		UserID:     credientials.ID.Hex(),
		Role:       credientials.Role,
		IsVerified: false,
	})

	if tokenErr != nil {
		_ = utils.WriteResponse(res, http.StatusInternalServerError, structs.IAppError{
			Message:    "Login Failed",
			Reason:     tokenErr.Error(),
			ErrorObj:   tokenErr,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	_ = utils.WriteResponse(res, http.StatusOK, structs.IAppSuccess{
		Message:    "Login Successfull",
		Data:       token,
		StatusCode: http.StatusOK,
	})
}
