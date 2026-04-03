// Package service provides service layer functions for user module
package service

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AlladinDev/AlShifa/constants"
	structs "github.com/AlladinDev/AlShifa/structs"
	"github.com/AlladinDev/AlShifa/users/dtos"
	interfaces "github.com/AlladinDev/AlShifa/users/interfaces"
	models "github.com/AlladinDev/AlShifa/users/models"
	"github.com/AlladinDev/AlShifa/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo interfaces.IRepository
}

var _ interfaces.IService = (*Service)(nil)

func ReturnNewService(repo interfaces.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddUser(ctx context.Context, userDetails models.User) *structs.IAppError {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": userDetails.Email},
			bson.M{"mobile": userDetails.Mobile},
		},
	}
	_, userExistsErr := s.repo.SearchUser(ctx, filter)
	//if error is not nil check if it is mongo no document it means user is fresh ignore it but if error is some other one return error
	if userExistsErr != nil {
		if userExistsErr != mongo.ErrNoDocuments {
			return &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: http.StatusInternalServerError,
				Reason:     userExistsErr.Error(),
				ErrorObj:   userExistsErr,
			}
		}
	}

	//if error is nil it definitely means user already exists
	if userExistsErr == nil {
		return &structs.IAppError{
			Message:    "Duplicate email or mobile",
			StatusCode: http.StatusBadRequest,
			Reason:     "duplicate email or mobile",
			ErrorObj:   errors.New("duplicate email or mobile"),
		}
	}

	//hash the password
	if err := s.repo.RegisterUser(ctx, userDetails); err != nil {
		return &structs.IAppError{
			Message:    "Registration Failed",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	return nil
}
func (s *Service) SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {

	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: 500,
				ErrorObj:   err,
				Reason:     "Server Error",
			}
		}
		return nil, &structs.IAppError{
			Message:    "User Not Found",
			StatusCode: 404,
			ErrorObj:   err,
			Reason:     "User Doesnt Exist",
		}

	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, loginDetails dtos.LoginDTO) (string, *structs.IAppError) {
	user, searchErr := s.repo.SearchUser(ctx, bson.M{"email": loginDetails.Email})
	if searchErr != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     searchErr.Error(),
			ErrorObj:   searchErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now match password
	passwordMatches, err := utils.VerifyPasswordArgon2id(loginDetails.Password, user.Password)
	if err != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !passwordMatches {
		return "", &structs.IAppError{
			Message:    "Invalid Credientials",
			Reason:     "invalid password or email",
			ErrorObj:   errors.New("invalid password or email"),
			StatusCode: http.StatusBadRequest,
		}
	}

	token, tokenErr := utils.GenerateJWT(&constants.JwtCustomClaims{
		UserID:     user.ID.Hex(),
		Role:       user.Role,
		Mobile:     strconv.Itoa(user.Mobile),
		IsVerified: true,
		Email:      user.Email,
	})

	if tokenErr != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     tokenErr.Error(),
			ErrorObj:   tokenErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}
