// Package service contains service layer implementation for clinic module
package service

import (
	interfaces "AlShifa/Clinic/Interfaces"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	appInterfaces "AlShifa/Interfaces"
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClinicService struct {
	OTPNotifier  appInterfaces.INotifier[string, string]
	OTPGenerator func(uniquePrefix string) string
	Repo         interfaces.IRepository
	MongoClient  *mongo.Client
	//this is the cache which will be used to store otp when a clinic wants to add a doctor
	AddDoctorToClinicAuthCache appInterfaces.ICache[string, models.AddDoctorOtpPayload]
}

func NewClinicService(repo interfaces.IRepository, mongoClient *mongo.Client, OTPNotifier appInterfaces.INotifier[string, string],
	OTPGenerator func(uniquePrefix string) string,
	authCache appInterfaces.ICache[string, models.AddDoctorOtpPayload]) *ClinicService {
	return &ClinicService{
		Repo:                       repo,
		AddDoctorToClinicAuthCache: authCache,
		OTPNotifier:                OTPNotifier,
		OTPGenerator:               OTPGenerator,
		MongoClient:                mongoClient,
	}
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IService = (*ClinicService)(nil)

//AddDoctorToClinic function generate otp for adding doctor to clinic this process needs another function or handler to verify otp and then the process is completed
func (service *ClinicService) AddDoctorToClinic(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError {

	//now check if clinic exists with this clinicID
	clinics, err := service.Repo.SearchClinic(ctx, bson.M{"_id": clinicDetails.ClinicID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To Clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if len(clinics) == 0 {
		return &structs.IAppError{
			Message:    "No Clinic Found With This Clinic",
			Reason:     errors.New("no clinic found with clinic id").Error(),
			ErrorObj:   errors.New("no clinic found with clinic id"),
			StatusCode: http.StatusNotFound,
		}
	}

	fmt.Print("reached here")
	clinic := clinics[0]

	//now check whether this doctorID exists or not
	doctor, err := service.Repo.SearchDoctor(ctx, bson.M{"_id": clinicDetails.DoctorID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To Clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now check whether doctor already has clinic or clinic has this doctor
	doctorAlreadyAdded := slices.Contains(clinic.Doctors, doctor.ID)
	if doctorAlreadyAdded {
		return &structs.IAppError{
			Message:    "Doctor Already Added To Clinic",
			Reason:     errors.New("doctor is already added to clinic").Error(),
			ErrorObj:   errors.New("doctor is already added to clinic"),
			StatusCode: http.StatusForbidden,
		}
	}

	//now send otp to doctor email through notification module and rest logic will be handled by otp verifier function
	//prepare otp payload
	uniquePrefix := fmt.Sprintf("%s:%s", clinicDetails.ClinicID.Hex(), clinicDetails.DoctorID.Hex())
	otp := service.OTPGenerator(uniquePrefix)
	otpPayload := models.AddDoctorOtpPayload{
		OTP:    otp,
		Expiry: time.Now().Add(utils.OTPExpiry),
		ClinicDetails: &models.AddDoctorToClinic{
			WorkingDays: clinicDetails.WorkingDays,
			StartTime:   clinicDetails.EndTime,
			EndTime:     clinicDetails.EndTime,
			ClinicID:    clinicDetails.ClinicID,
			DoctorID:    clinicDetails.DoctorID,
		},
	}

	//now save this otp payload in cache
	if err := service.AddDoctorToClinicAuthCache.Set(ctx, otp, otpPayload, utils.CacheTTL); err != nil {
		return &structs.IAppError{
			Message:    "Failed To Add Doctor To Clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusForbidden,
		}
	}

	//now send the otp to doctor email through notification module
	notifierErr := service.OTPNotifier.SendNotification(doctor.Email, otp)
	if notifierErr != nil {
		return &structs.IAppError{
			Message:    "Failed To Add Doctor To Clinic OTP Error",
			Reason:     notifierErr.Error(),
			ErrorObj:   notifierErr,
			StatusCode: http.StatusForbidden,
		}
	}

	//now send nil as error
	return nil

}

func (service *ClinicService) VerifyAddDoctorToClinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError {
	//check whether in cache this otp exists or not
	otpToFind := fmt.Sprintf("%s:%s:%s", clinicID.Hex(), doctorID.Hex(), otp)
	fmt.Print(otpToFind)
	otpPayload, otpExists, err := service.AddDoctorToClinicAuthCache.Get(ctx, otpToFind)
	if !otpExists {
		return &structs.IAppError{
			Message:    "OTP Doesnt Exist",
			ErrorObj:   errors.New("OTP doesnt exist"),
			Reason:     "OTP doesnt exist",
			StatusCode: 404,
		}
	}
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to verify otp",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: 500,
		}
	}

	//now in otp payload check whether this clinicid is present in otp payload or not
	if otpPayload.ClinicDetails.ClinicID != clinicID {
		return &structs.IAppError{
			Message:    "OTP Verification failed",
			ErrorObj:   errors.New("clinic Doesnt Match"),
			Reason:     errors.New("clinic Doesnt Match").Error(),
			StatusCode: 400,
		}
	}

	//now check whether otp expiry hasnt reached if yes return error
	if time.Now().After(otpPayload.Expiry) {
		return &structs.IAppError{
			Message:    "OTP Expired",
			ErrorObj:   errors.New("OTP Expired"),
			Reason:     errors.New("OTP Expired").Error(),
			StatusCode: 400,
		}
	}

	//now as everything is ok now using transaction add doctorID to clinic and also add clinicID to doctor data in database
	session, err := service.MongoClient.StartSession()
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Verify Otp",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: 500,
		}
	}
	defer session.EndSession(ctx)

	//actual transaction function in which logic to add doctor to clinic resides
	transactionFn := func(ctx mongo.SessionContext) (any, error) {
		if err := service.Repo.AddDoctorToClinic(ctx, *otpPayload.ClinicDetails); err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, transactionErr := session.WithTransaction(ctx, transactionFn)
	if transactionErr != nil {
		return &structs.IAppError{
			Message:    "OTP verification Failed",
			ErrorObj:   transactionErr,
			Reason:     transactionErr.Error(),
			StatusCode: 500,
		}
	}

	//here after verifying remove otp from redis also
	if err := service.AddDoctorToClinicAuthCache.Delete(ctx, otp); err != nil {
		fmt.Print("Failed to remove otp aftr verification")
	}

	return nil

}

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
	owners, ownerMongoDBErr := service.Repo.GetOwnerDetails(ctx, bson.M{"$or": []bson.M{
		{"email": ownerDetails.Email},
		{"mobile": ownerDetails.Mobile},
	}})

	if ownerMongoDBErr != nil {
		return utils.ReturnAppError(ownerMongoDBErr, http.StatusInternalServerError, "Failed to Register Owner", ownerMongoDBErr.Error())
	}

	if len(owners) > 0 {
		return utils.ReturnAppError(ownerMongoDBErr, http.StatusForbidden, "This email or number already exists", "This email or mobile already exists")
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
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Clinic details", err.Error())
	}

	return clinics, nil
}

func (service *ClinicService) SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError) {
	owners, err := service.Repo.GetOwnerDetails(ctx, filter)
	if err != nil {
		fmt.Print(err, owners)
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Owner details", "Server Error")
	}

	return owners, nil
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

func (service *ClinicService) SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError) {
	// //here validate filters
	// allowedFilters := []string{"_id", "name", "mobile", "email"}
	// for keys := range filter {

	// }
	doctors, err := service.Repo.SearchDoctors(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to Fetch Doctors",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return doctors, nil
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

func (service *ClinicService) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (*models.Appointment, *structs.IAppError) {
	//here check whether this doctorid and clinicid exists or not  and userid will be injected by jwt middleware

	//now do the searches
	doctors, err := service.SearchDoctor(ctx, bson.M{"_id": appointmentDetails.Doctor})
	if err != nil {
		return nil, utils.ReturnAppError(err, http.StatusInternalServerError, "Failed to add appoinment", err.Error())
	}

	if len(doctors) == 0 {
		return nil, utils.ReturnAppError(errors.New("no doctor exists with this id"), http.StatusInternalServerError, "Failed to add appoinment", "no doctor exists with this id")
	}

	//now search for this clinic
	clinics, err := service.SearchClinic(ctx, bson.M{"_id": appointmentDetails.Clinic})
	if err != nil {
		return nil, utils.ReturnAppError(err, http.StatusInternalServerError, "Failed to add appoinment", err.Error())
	}

	if len(clinics) == 0 {
		return nil, utils.ReturnAppError(errors.New("no clinic exists with this id"), http.StatusInternalServerError, "Failed to add appoinment", "no clinic exists with this id")
	}

	//get the first clinic
	clinic := clinics[0]

	//here check this doctor id must be present in the clinic.doctors array means it must be onboarded in this clinic
	if !slices.Contains(clinic.Doctors, appointmentDetails.Doctor) {
		return nil, utils.ReturnAppError(errors.New("this doctor is not onboarded in this clinic"), http.StatusInternalServerError, "Failed to add appoinment", "this doctor is not onboarded in this clinic")
	}

	//now add some defaults like status date
	appointmentDetails.RegistrationDate = time.Now()
	appointmentDetails.Status = "pending"
	appointmentDetails.ID = primitive.NewObjectID()

	appointmentCreated, appointmentErr := service.Repo.AddAppointment(ctx, clinic.MaxAppointments, appointmentDetails)
	if appointmentErr != nil {
		//here check for duplicate key error if yes then send error as MaxSlots booked
		if mongo.IsDuplicateKeyError(appointmentErr) {
			return nil, utils.ReturnAppError(errors.New("max Appointments Reached For Today"), http.StatusExpectationFailed, "Max Appointments Reached For today", "Max Appointments Reached For today")
		}
		return nil, utils.ReturnAppError(appointmentErr, http.StatusInternalServerError, "Failed to add appointment", appointmentErr.Error())
	}
	return appointmentCreated, nil
}

//AppointmentSlotsBooked function returns slots for whom maxAppointment has reached so all those slot documents where maxAppointment has reached are returned
func (service *ClinicService) AppointmentSlotsBooked(ctx context.Context, slotDetais models.SlotDetails) ([]models.Slot, *structs.IAppError) {
	slots, err := service.Repo.AppointmentSlotsBooked(ctx, slotDetais.MaxAppointments, slotDetais.ClinicID, slotDetais.DoctorID)
	if err != nil {
		return nil, utils.ReturnAppError(err, http.StatusInternalServerError, "Failed to Fetch Booked Slots", err.Error())
	}

	return slots, nil
}
