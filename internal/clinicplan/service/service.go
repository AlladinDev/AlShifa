// Package planservice provides service function for plan module
package planservice

import (
	"context"
	"net/http"
	"time"

	planinterface "github.com/AlladinDev/AlShifa/internal/clinicplan/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"github.com/AlladinDev/AlShifa/internal/clinicplan/validators"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo planinterface.IRepository
}

func NewService(Repo planinterface.IRepository) planinterface.IService {
	return &Service{
		repo: Repo,
	}
}

func (s *Service) FetchAmountToDeduct(ctx context.Context, filter bson.M) (primitive.ObjectID, int, error) {
	return s.repo.FetchAmountToDeduct(ctx, filter)
}

func (s *Service) RegisterPlan(ctx context.Context, details *models.Plan) *response.IAppError {
	//here add plan createdAt and updatedAt date and also generate unique id
	details.CreatedAt = time.Now()
	details.UpdatedAt = time.Now()
	details.ID = primitive.NewObjectID()

	//now here do a validation
	if err := validators.ValidatePlanDetails(details); err != nil {
		return &response.IAppError{
			Message:    "Invalid plan details",
			Reason:     "validation failed",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := s.repo.AddPlan(ctx, details); err != nil {
		return &response.IAppError{
			Message:    "Failed to add plan",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *Service) FetchPlanID(ctx context.Context, planName string) (primitive.ObjectID, error) {
	return s.repo.FetchPlanID(ctx, planName)
}

func (s *Service) FetchPlans(ctx context.Context, filter bson.M) ([]models.Plan, *response.IAppError) {
	plans, err := s.repo.FetchPlans(ctx, filter)
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed to Fetch Plans",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return plans, nil
}

func (s *Service) UpdatePlan(ctx context.Context, ID primitive.ObjectID, data *models.Plan) *response.IAppError {
	if err := validators.ValidatePlanDetails(data); err != nil {
		return &response.IAppError{
			Message:    "Invalid plan details",
			Reason:     "validation failed",
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if err := s.repo.UpdatePlan(ctx, ID, data); err != nil {
		return &response.IAppError{
			Message:    "Failed to update plan",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   err,
			Reason:     err.Error(),
		}
	}

	return nil
}
