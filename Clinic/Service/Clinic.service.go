// Package service contains service layer implementation for clinic module
package service

import (
	interfaces "AlShifa/Clinic/Interfaces"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClinicService struct {
	Repo interfaces.IRepository
}

func NewClinicService(repo interfaces.IRepository) *ClinicService {
	return &ClinicService{
		Repo: repo,
	}
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IService = (*ClinicService)(nil)

func (service *ClinicService) RegisterClinic(ctx context.Context, ownerID string, clinicDetails models.Clinic) *structs.IAppError {
	//first convert ownerID to proper mongodb id
	ownerMongoDBID, err := primitive.ObjectIDFromHex(string(ownerID))
	if err != nil {
		return utils.ReturnAppError(err, 500, "Unable TO Register CLinic", "Server Error")
	}

	//now first check if against this ownerId owner exists or not
	_, ownerExistingErr := service.Repo.GetOwnerDetails(ctx, bson.M{"_id": ownerMongoDBID})
	if ownerExistingErr != nil {
		return utils.ReturnAppError(ownerExistingErr, 500, "Failed to register clinic", "Server error")
	}

	//first do validation
	validationErr := validators.ValidateClinicDetails(&clinicDetails)
	if validationErr != nil {
		return utils.ReturnAppError(validationErr, 400, "Registration failed", "Invalid Details")
	}

	// set default values
	clinicDetails.RegistrationDate = time.Now().UTC()
	clinicDetails.Wallet = primitive.NilObjectID
	clinicDetails.Doctors = nil
	registrationErr := service.Repo.RegisterClinic(ctx, ownerMongoDBID, clinicDetails)
	if registrationErr != nil {
		return utils.ReturnAppError(registrationErr, 500, "Registration Failed", "Unknown reason")
	}

	return nil

}

func (service *ClinicService) RegisterClinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError {

	//now check if email or mobile exists
	_, ownerMongoDBErr := service.Repo.GetOwnerDetails(ctx, bson.M{"$or": []bson.M{
		{"email": ownerDetails.Email},
		{"mobile": ownerDetails.Mobile},
	}})

	if ownerMongoDBErr == nil {
		fmt.Print("issue ", ownerMongoDBErr)
		return utils.ReturnAppError(ownerMongoDBErr, 400, "User Already Exists", "Duplicate Email or PhoneNumber")
	}

	//now hash the password
	hashedPassword, hashingErr := utils.HashPasswordArgon2id(ownerDetails.Password)
	if hashingErr != nil {
		return utils.ReturnAppError(hashingErr, 500, "Registration Failed", "Server Issue")
	}

	ownerDetails.Password = hashedPassword

	ownerDetails.RegistrationDate = time.Now().UTC()
	ownerDetails.Clinic = primitive.NilObjectID
	ownerDetails.ID = primitive.NewObjectID()

	//now call the repo method to register owner
	registrationErr := service.Repo.RegisterClinicOwner(ctx, ownerDetails)
	if registrationErr != nil {
		return utils.ReturnAppError(registrationErr, 500, "Failed to register owner", "Server error")

	}

	return nil

	//if error is there it will return it else it will return nil automatically

}

func (service *ClinicService) SearchClinic(ctx context.Context, clinicName string) ([]models.Clinic, *structs.IAppError) {
	clinics, err := service.Repo.SearchClinic(ctx, bson.M{"name": clinicName})
	if err != nil {
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Clinic details", "Server Error")
	}

	return clinics, nil
}
