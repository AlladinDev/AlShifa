// Package service contains service layer implementation for clinic module
package service

import (
	"context"

	"fmt"

	"net/http"

	"time"

	interfaces "github.com/AlladinDev/AlShifa/internal/clinic/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinic/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"

	customerrors "github.com/AlladinDev/AlShifa/internal/shared/customerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clinicService struct {
	Repo       interfaces.IRepository
	clinicPlan interfaces.IClinicPlan
}

func NewclinicService(repo interfaces.IRepository, clinicPlanModule interfaces.IClinicPlan) *clinicService {
	return &clinicService{
		Repo:       repo,
		clinicPlan: clinicPlanModule,
	}
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IService = (*clinicService)(nil)

//FetchMaxAppointments function needed by appointment module for booking appointments
func (service *clinicService) FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error) {
	return service.Repo.FetchMaxAppointments(ctx, clinicID)
}

//DeductClinicMoneyForAppointment function needed by appointment module for booking appointments
func (service *clinicService) DeductClinicMoneyForAppointment(ctx context.Context, clinicID primitive.ObjectID) error {
	clinics, clinicSearchErr := service.Repo.Searchclinic(ctx, bson.M{"_id": clinicID})
	if clinicSearchErr != nil {
		return clinicSearchErr
	}

	if len(clinics) == 0 {
		return customerrors.ErrClinicNotFound
	}
	clinic := clinics[0]
	//now check whether it has wallet or not
	if clinic.Wallet == nil {
		return customerrors.ErrClinicWalletNotFound
	}

	//here check if clinic has a plan attached or not
	if clinic.PlanID == primitive.NilObjectID {
		return customerrors.ErrClinicPlanNotFound
	}

	//now get the plan details to check how much to deduct and whether this clinic has that amount or not
	_, amountToDeduct, planErr := service.clinicPlan.FetchAmountToDeduct(ctx, bson.M{"_id": clinic.PlanID})

	if planErr != nil {
		return planErr
	}

	//now check its balance
	if clinic.Wallet.AvailableBalance < int64(amountToDeduct) {
		return customerrors.ErrClinicInsufficientWalletBalance
	}

	return service.Repo.DeductClinicWallet(ctx, amountToDeduct, clinicID)
}

//ClinicExists function needed by appointment module for booking appointments
func (service *clinicService) ClinicExists(ctx context.Context, filter bson.M) (bool, error) {
	return service.Repo.ClinicExists(ctx, filter)
}

func (service *clinicService) Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinicDetails models.Clinic) *response.IAppError {
	// set default values
	clinicDetails.CreatedAt = time.Now().UTC()

	clinicDetails.ID = primitive.NewObjectID()

	//here get basic plan and add its id to clinnic
	planID, planErr := service.clinicPlan.FetchPlanID(ctx, constants.ClinicPlanBasic)
	if planErr != nil {
		return &response.IAppError{
			Message:    "Failed to register clinic",
			Reason:     planErr.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   planErr,
		}
	}

	//add planid to clinic
	clinicDetails.PlanID = planID

	//add free wallet balance to use our service first
	clinicDetails.Wallet = &models.WalletDetails{
		AvailableBalance: constants.FreeWalletBalance,
	}

	clinicDetails.OwnerID = ownerID

	registrationErr := service.Repo.Registerclinic(ctx, ownerID, clinicDetails)
	if registrationErr != nil {
		fmt.Print(registrationErr)
		return utils.ReturnAppError(registrationErr, 500, "Registration Failed", "Unknown reason")
	}

	return nil

}

func (service *clinicService) Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, *response.IAppError) {

	clinics, err := service.Repo.Searchclinic(ctx, filter)
	if err != nil {
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch clinic details", err.Error())
	}

	return clinics, nil
}

func (service *clinicService) GetClinicIDByOwnerID(ctx context.Context, ownerID primitive.ObjectID) (clinicID primitive.ObjectID, err error) {
	clinic, err := service.Repo.FetchSingleClinic(ctx, bson.M{"ownerID": ownerID})
	if err != nil {
		return primitive.NilObjectID, err
	}

	return clinic.ID, nil
}

func (service *clinicService) GetClinicName(ctx context.Context, clinicID primitive.ObjectID) (clinicName string, err error) {
	clinic, err := service.Repo.FetchSingleClinic(ctx, bson.M{"_id": clinicID})
	if err != nil {
		return "", err
	}

	return clinic.Name, nil
}
