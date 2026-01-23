// Package service contains service layer implementation for clinic module
package service

import (
	interfaces "AlShifa/Clinic/Interfaces"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		fmt.Print(err)
		return utils.ReturnAppError(err, 500, "Unable TO Register CLinic", "Server Error")
	}

	//now first check if against this ownerId owner exists or not
	_, ownerExistingErr := service.Repo.GetOwnerDetails(ctx, bson.M{"_id": ownerMongoDBID})
	if ownerExistingErr != nil {
		return utils.ReturnAppError(ownerExistingErr, 500, "Failed to register clinic", "Server error")
	}

	//first do validation
	validationErr := validators.ValidateClinicDetails(&clinicDetails)
	if len(validationErr) != 0 {
		return utils.ReturnAppError(validationErr, 400, "Registration failed", "Invalid Details")
	}

	// set default values
	clinicDetails.RegistrationDate = time.Now().UTC()
	clinicDetails.Wallet = primitive.NilObjectID
	clinicDetails.ID = primitive.NewObjectID()
	clinicDetails.Doctors = nil
	registrationErr := service.Repo.RegisterClinic(ctx, ownerMongoDBID, clinicDetails)
	if registrationErr != nil {
		fmt.Print(registrationErr)
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
	ownerDetails.Role = utils.RoleClinicOwner

	//now call the repo method to register owner
	registrationErr := service.Repo.RegisterClinicOwner(ctx, ownerDetails)
	if registrationErr != nil {
		return utils.ReturnAppError(registrationErr, 500, "Failed to register owner", "Server error")

	}

	return nil

	//if error is there it will return it else it will return nil automatically

}

func (service *ClinicService) SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {
	clinics, err := service.Repo.SearchClinic(ctx, filter)
	if err != nil {
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Clinic details", "Server Error")
	}

	return clinics, nil
}

func (service *ClinicService) SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError) {
	owner, err := service.Repo.GetOwnerDetails(ctx, filter)
	if err != nil {
		fmt.Print(err, owner)
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Owner details", "Server Error")
	}

	return owner, nil
}

func (service *ClinicService) RegisterDoctor(ctx context.Context, doctor models.Doctor) *structs.IAppError {

	//here check if doctor exists using mobile and phoneNumber
	existingDoctors, err := service.Repo.SearchDoctors(ctx, bson.M{
		"$or": []bson.M{
			{"email": doctor.Email},
			{"mobile": doctor.Mobile},
		},
	})

	///if error is nill check if it is of other type  and return error
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return utils.ReturnAppError(err, 500, "Registration Failed", "Server Error")
		}
	}

	if len(existingDoctors) > 0 {
		return utils.ReturnAppError(errors.New("doctor already exists"), 400, "Email or Mobile Already Exists", "Duplicate Email or Mobile")
	}

	//here set the default values
	doctor.Clinics = nil
	doctor.Appointments = nil
	doctor.RegistrationDate = time.Now()
	doctor.ID = primitive.NewObjectID()
	doctor.Role = utils.RoleDoctor

	hashedPassword, err := utils.HashPasswordArgon2id(doctor.Password)
	if err != nil {
		return utils.ReturnAppError(err, 500, "Registration Failed", "Server Error")
	}
	doctor.Password = hashedPassword

	if err := service.Repo.RegisterDoctor(ctx, doctor); err != nil {
		return utils.ReturnAppError(err, 500, "Unable To Add Doctor", "Server Error")
	}

	return nil
}

func (service *ClinicService) SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error) {
	// //here validate filters
	// allowedFilters := []string{"_id", "name", "mobile", "email"}
	// for keys := range filter {

	// }
	return service.Repo.SearchDoctors(ctx, filter)
}

func (service *ClinicService) LoginClinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	owners, err := service.Repo.GetOwnerDetails(ctx, bson.M{"email": email})
	if err != nil {
		return "", utils.ReturnAppError(err, 404, "Owner Not Found", "Invalid Email or Password")
	}

	if len(owners) == 0 {
		return "", utils.ReturnAppError(errors.New("owner not found"), 404, "Owner Not Found", "Invalid Email or Password")
	}

	owner := owners[0]

	passwordMatches, err := utils.VerifyPasswordArgon2id(password, owner.Password)
	if err != nil || !passwordMatches {
		return "", utils.ReturnAppError(err, 401, "Unauthorized", "Invalid Email or Password")
	}

	token, err := utils.GenerateJWT(owner.ID.Hex(), owner.Role)
	if err != nil {
		return "", utils.ReturnAppError(err, 500, "Login Failed", "Server Error")
	}

	return token, nil
}

func (service *ClinicService) LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	doctor, err := service.Repo.SearchDoctor(ctx, bson.M{"email": email})
	if err != nil {
		return "", utils.ReturnAppError(err, 404, "Doctor Not Found", "Invalid Email or Password")
	}

	passwordMatches, err := utils.VerifyPasswordArgon2id(password, doctor.Password)
	if err != nil || !passwordMatches {
		return "", utils.ReturnAppError(err, 401, "Unauthorized", "Invalid Email or Password")
	}
	token, err := utils.GenerateJWT(doctor.ID.Hex(), doctor.Role)
	if err != nil {
		return "", utils.ReturnAppError(err, 500, "Login Failed", "Server Error")
	}
	return token, nil
}
