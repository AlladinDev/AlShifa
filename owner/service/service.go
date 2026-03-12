// Package service provides service functions for owner module
package service

import (
	"context"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/owner/interfaces"
	"github.com/AlladinDev/AlShifa/owner/models"
	"github.com/AlladinDev/AlShifa/structs"

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

	//add some default things like createdAt role
	ownerDetails.Role = constants.RoleclinicOwner
	ownerDetails.CreatedAt = time.Now()

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
