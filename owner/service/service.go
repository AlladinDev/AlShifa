// Package service provides service functions for owner module
package service

import (
	"context"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/dtos"
	"github.com/AlladinDev/AlShifa/owner/interfaces"
	"github.com/AlladinDev/AlShifa/owner/models"
	"github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo interfaces.IRepository
}

func NewService(repo interfaces.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

var _ interfaces.IService = (*Service)(nil)

func (s *Service) RegisterOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError {

	//first check whether this email or mobile already exists if yes throw error
	ownerExists, searchErr := s.repo.GetOwnerDetails(ctx, bson.M{"$or": []bson.M{
		{"email": ownerDetails.Email},
		{"mobile": ownerDetails.Mobile},
	}})

	if searchErr != nil && searchErr != mongo.ErrNoDocuments {
		return &structs.IAppError{
			Message:    "Failed to register owner",
			Reason:     searchErr.Error(),
			ErrorObj:   searchErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if ownerExists != nil {
		return &structs.IAppError{
			Message:    "Owner Already Exists with this email or mobile",
			Reason:     "This Email or mobile already exists",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	//add some default things like createdAt role
	ownerDetails.Role = constants.RoleclinicOwner
	ownerDetails.CreatedAt = time.Now()
	ownerDetails.ID = primitive.NewObjectID()

	//hash the password now
	hashedPassword, hashingErr := utils.HashPasswordArgon2id(ownerDetails.Password)
	if hashingErr != nil {
		return &structs.IAppError{
			Message:    "Failed to register owner",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now update the raw password with this hashedpassword
	ownerDetails.Password = hashedPassword

	if err := s.repo.RegisterOwner(ctx, ownerDetails); err != nil {
		return &structs.IAppError{
			Message:    "Failed to register owner",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *Service) GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, *structs.IAppError) {
	owner, err := s.repo.GetOwnerByID(ctx, ownerID)
	if err != nil {
		errMsg := "Failed to Fetch OwnerDetails"
		errStatusCode := http.StatusInternalServerError
		if err == mongo.ErrNoDocuments {
			errMsg = "Owner Doesnt Exist"
			errStatusCode = http.StatusNotFound
		}
		return nil, &structs.IAppError{
			Message:    errMsg,
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: errStatusCode,
		}
	}

	//hide owner password dont sent it to frontend
	owner.Password = ""
	return owner, nil
}

func (s *Service) GetOwnerDetails(ctx context.Context, filters bson.M) (*models.Owner, *structs.IAppError) {
	owner, err := s.repo.GetOwnerDetails(ctx, filters)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch ownerdetails",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return owner, nil
}

func (s *Service) LoginOwner(ctx context.Context, loginDetails dtos.LoginEmailPasswordDTO) (string, *structs.IAppError) {
	owner, err := s.repo.GetOwnerDetails(ctx, bson.M{"email": loginDetails.Email})
	if err != nil {
		return "", &structs.IAppError{
			Message:    "Invalid Email or password",
			Reason:     "invalid credientials",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now check password matching
	passwordMatches, matchingErr := utils.VerifyPasswordArgon2id(loginDetails.Password, owner.Password)
	if matchingErr != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     matchingErr.Error(),
			ErrorObj:   matchingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !passwordMatches {
		return "", &structs.IAppError{
			Message:    "Invalid Email or password",
			Reason:     "invalid credientials",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now as everything is ok generate jwt token
	token, err := utils.GenerateJWT(&constants.JwtCustomClaims{
		UserID: owner.ID.Hex(),
		Role:   owner.Role,
		Mobile: owner.Mobile,
		Email:  owner.Email,
	})

	if err != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}
