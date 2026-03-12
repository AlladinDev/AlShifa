// Package service provides service layer functions for user module
package service

import (
	"context"
	"net/http"

	structs "github.com/AlladinDev/AlShifa/structs"
	interfaces "github.com/AlladinDev/AlShifa/users/interfaces"
	models "github.com/AlladinDev/AlShifa/users/models"

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
