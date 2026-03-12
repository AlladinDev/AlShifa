// Package service provides service functions for auth module
package service

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/AlladinDev/AlShifa/auth/interfaces"
	"github.com/AlladinDev/AlShifa/auth/models"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo interfaces.IRepository
}

func NewService(repo interfaces.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, credientials models.Credientials) *structs.IAppError {
	//now here check if role is admin  return error
	if credientials.Role == constants.RoleAdmin {
		return &structs.IAppError{
			Message:    "Invalid Role",
			Reason:     "Invalid Role",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//first check if this email exists if yes return error
	userExists, err := s.repo.SearchCredientials(ctx, bson.M{"email": credientials.Email, "mobile": credientials.Mobile})
	if err != nil {
		return &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if userExists != nil {
		return &structs.IAppError{
			Message:    "This email or mobile already exists",
			Reason:     "Duplicate email or mobile",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	//check role
	if !slices.Contains(constants.RolesAllowed, credientials.Role) {
		return &structs.IAppError{
			Message:    "Invalid Role",
			Reason:     "Invalid Role",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//add default things
	credientials.CreatedAt = time.Now()
	credientials.IsVerified = false
	credientials.ProfileCompleted = false
	credientials.ID = primitive.NewObjectID()

	//hash the password
	hashedPassword, hashingErr := utils.HashPasswordArgon2id(credientials.Password)
	if hashingErr != nil {
		return &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}
	//replace raw password with hashed password
	credientials.Password = hashedPassword

	//now register the user
	if err := s.repo.Register(ctx, credientials); err != nil {
		return &structs.IAppError{
			Message:    "Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *Service) UpdateVerificationStatus(ctx context.Context, userID primitive.ObjectID, status bool) *structs.IAppError {
	//first check whether a credential exists as per this id or not
	_, err := s.repo.SearchCredientials(ctx, bson.M{"_id": userID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Update Status",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	if err := s.repo.UpdateVerificationStatus(ctx, userID, status); err != nil {
		return &structs.IAppError{
			Message:    "Failed to Update Status",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}
	return nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	userExists, err := s.repo.SearchCredientials(ctx, bson.M{"email": email})
	if err != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	if userExists == nil {
		return "", &structs.IAppError{
			Message:    "Invalid Email Or Password",
			StatusCode: http.StatusInternalServerError,
			Reason:     "Invalid Email Or Password",
			ErrorObj:   nil,
		}
	}

	passwordMatches, err := utils.VerifyPasswordArgon2id(password, userExists.Password)
	if err != nil || !passwordMatches {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	token, err := utils.GenerateJWT(&constants.JwtCustomClaims{
		Mobile: userExists.Mobile,
		Email:  userExists.Email,
		Role:   userExists.Email,
		UserID: userExists.ID.Hex(),
	})

	if err != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}
	return token, nil
}

func (s *Service) SearchCredientials(ctx context.Context, filter bson.M) (*models.Credientials, *structs.IAppError) {
	credientials, err := s.repo.SearchCredientials(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to Fetch User",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return credientials, nil
}
