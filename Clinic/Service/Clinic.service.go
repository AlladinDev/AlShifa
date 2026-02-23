// Package service contains service layer implementation for clinic module
package service

import (
	appInterfaces "AlShifa/Interfaces"
	interfaces "AlShifa/clinic/Interfaces"
	DTO "AlShifa/clinic/dtos"
	"AlShifa/clinic/models"
	sharedModels "AlShifa/models"
	structs "AlShifa/structs"
	utils "AlShifa/utils"
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

type clinicService struct {
	OTPNotifier  appInterfaces.INotifier[string, string]
	OTPGenerator func(uniquePrefix string) string
	Repo         interfaces.IRepository
	MongoClient  *mongo.Client
	//this is the cache which will be used to store otp when a clinic wants to add a doctor
	OtpCache appInterfaces.ICache[string, models.AddDoctorOtpPayload]
}

func NewclinicService(repo interfaces.IRepository, mongoClient *mongo.Client, OTPNotifier appInterfaces.INotifier[string, string],
	OTPGenerator func(uniquePrefix string) string,
	authCache appInterfaces.ICache[string, models.AddDoctorOtpPayload]) *clinicService {
	return &clinicService{
		Repo:         repo,
		OtpCache:     authCache,
		OTPNotifier:  OTPNotifier,
		OTPGenerator: OTPGenerator,
		MongoClient:  mongoClient,
	}
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IService = (*clinicService)(nil)

func (service *clinicService) FetchAppointments(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointments, *structs.IAppError) {
	//do little bit validation here filter can be empty but groupBy must be one of these user,doctor,clinic
	groupingTagsAllowed := []string{"user", "doctor", "clinic"}
	if !slices.Contains(groupingTagsAllowed, groupBy) {
		return nil, &structs.IAppError{
			Message:    "Invalid GroupBy Tags",
			Reason:     "Invalid Grouping tag",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	appointments, err := service.Repo.FetchAppointments(ctx, groupBy, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch appointments",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return appointments, nil
}

//AddDoctorToclinic function generate otp for adding doctor to clinic this process needs another function or handler to verify otp and then the process is completed
func (service *clinicService) AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError {
	//now check if clinic exists with this ClinicID
	fmt.Println("clinicid is", clinicDetails.ClinicID)
	clinic, err := service.Repo.SearchclinicByID(ctx, clinicDetails.ClinicID)
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinics",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now check whether this doctorID exists or not
	doctor, err := service.Repo.SearchDoctor(ctx, bson.M{"_id": clinicDetails.DoctorID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//here check clinicDoctor model to check if this doctorid and ClinicID exists if yes it means doctor is already onboarded there
	clinicDoctors, err := service.Repo.Searchclinic(ctx, bson.M{"ClinicID": clinic.ID, "doctorID": doctor.ID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//if len of clinicdoctors is more than 0 with this clinicid and doctorid it means doctor is already onboarded to this clinic
	if len(clinicDoctors) > 0 {
		fmt.Println(clinicDoctors)
		return &structs.IAppError{
			Message:    "Doctor is already onboarded to this clinic",
			Reason:     errors.New("duplicate onboarding").Error(),
			ErrorObj:   errors.New("duplicate onboarding"),
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
		ClinicDetails: &models.AddDoctorToclinic{
			DoctorName:    doctor.Name,
			ClinicName:    clinic.Name,
			ClinicAddress: clinic.Address,
			WorkingDays:   clinicDetails.WorkingDays,
			StartTime:     clinicDetails.EndTime,
			EndTime:       clinicDetails.EndTime,
			ClinicID:      clinicDetails.ClinicID,
			DoctorID:      clinicDetails.DoctorID,
		},
	}

	//now save this otp payload in cache
	if err := service.OtpCache.Set(ctx, otp, otpPayload, utils.CacheTTL); err != nil {
		return &structs.IAppError{
			Message:    "Failed To Add Doctor To clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusForbidden,
		}
	}

	//now send the otp to doctor email through notification module
	notifierErr := service.OTPNotifier.SendNotification(doctor.Email, otp)
	if notifierErr != nil {
		return &structs.IAppError{
			Message:    "Failed To Add Doctor To clinic OTP Error",
			Reason:     notifierErr.Error(),
			ErrorObj:   notifierErr,
			StatusCode: http.StatusForbidden,
		}
	}

	//now send nil as error
	return nil

}

func (service *clinicService) VerifyAddDoctorToclinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, ClinicID primitive.ObjectID) *structs.IAppError {
	//check whether in cache this otp exists or not
	otpToFind := fmt.Sprintf("%s:%s:%s", ClinicID.Hex(), doctorID.Hex(), otp)
	fmt.Print(otpToFind)
	otpPayload, otpExists, err := service.OtpCache.Get(ctx, otpToFind)
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
	if otpPayload.ClinicDetails.ClinicID != ClinicID {
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

	if err := service.Repo.AddDoctorToclinic(ctx, *otpPayload.ClinicDetails); err != nil {
		return &structs.IAppError{
			Message:    "Failed to add doctor to clinic",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: 500,
		}
	}

	//here after verifying remove otp from redis also
	if err := service.OtpCache.Delete(ctx, otpToFind); err != nil {
		fmt.Print("Failed to remove otp aftr verification")
	}

	return nil

}

func (service *clinicService) Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinicDetails models.Clinic) *structs.IAppError {

	//now first check if against this ownerId owner exists or not
	owners, ownerExistingErr := service.Repo.GetOwnerDetails(ctx, bson.M{"_id": ownerID})
	if ownerExistingErr != nil {
		return utils.ReturnAppError(ownerExistingErr, 500, "Failed to register clinic", "Server error")
	}

	//now here check if owner already has a clinic dont allow another one
	owner := owners[0]
	if owner.Clinic != primitive.NilObjectID {
		return utils.ReturnAppError(errors.New("clinic Already Exists For This Owner"), http.StatusUnauthorized, "clinic Already exists", "clinic Already exists")
	}

	// set default values
	clinicDetails.RegistrationDate = time.Now().UTC()
	clinicDetails.Wallet = nil
	clinicDetails.ID = primitive.NewObjectID()
	clinicDetails.PlanType = utils.PlanPaid
	registrationErr := service.Repo.Registerclinic(ctx, ownerID, clinicDetails)
	if registrationErr != nil {
		fmt.Print(registrationErr)
		return utils.ReturnAppError(registrationErr, 500, "Registration Failed", "Unknown reason")
	}

	return nil

}

func (service *clinicService) RegisterclinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError {

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
	ownerDetails.Role = utils.RoleclinicOwner

	//now call the repo method to register owner
	registrationErr := service.Repo.RegisterclinicOwner(ctx, ownerDetails)
	if registrationErr != nil {
		return utils.ReturnAppError(registrationErr, 500, "Failed to register owner", "Server error")

	}

	return nil

	//if error is there it will return it else it will return nil automatically

}

func (service *clinicService) Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError) {

	clinics, err := service.Repo.Searchclinic(ctx, filter)
	if err != nil {
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch clinic details", err.Error())
	}

	return clinics, nil
}

func (service *clinicService) SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError) {
	owners, err := service.Repo.GetOwnerDetails(ctx, filter)
	if err != nil {
		fmt.Print(err, owners)
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch Owner details", "Server Error")
	}

	return owners, nil
}

func (service *clinicService) RegisterDoctor(ctx context.Context, doctor models.Doctor) *structs.IAppError {

	//here check if doctor exists using mobile and phoneNumber
	existingDoctors, err := service.Repo.SearchDoctors(ctx, bson.M{
		"$or": []bson.M{
			{"email": doctor.Email},
			{"mobile": doctor.Mobile},
		},
	})

	fmt.Print(existingDoctors)

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

func (service *clinicService) SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError) {

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

func (service *clinicService) LoginclinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError) {
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

func (service *clinicService) LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError) {
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

func (service *clinicService) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError) {
	//here check whether this doctorid and clinicid exists or not  and userid will be injected by jwt middleware

	//now do the searches
	doctor, err := service.Repo.SearchDoctor(ctx, bson.M{"_id": appointmentDetails.Doctor})
	if err != nil {
		return 0, utils.ReturnAppError(err, http.StatusInternalServerError, "Failed to add appoinment", err.Error())
	}

	//now search for this clinic
	clinics, clinicSearchErr := service.Searchclinic(ctx, bson.M{"clinicID": appointmentDetails.Clinic})
	if clinicSearchErr != nil {
		return 0, utils.ReturnAppError(clinicSearchErr, http.StatusInternalServerError, "Failed to add appoinment", clinicSearchErr.Error())
	}

	if len(clinics) == 0 {
		return 0, utils.ReturnAppError(errors.New("no clinic exists with this id"), http.StatusInternalServerError, "Failed to add appoinment", "no clinic exists with this id")
	}

	//get the first clinic
	clinic := clinics[0]

	//here check using clinicDoctor whether this doctor is onboarded to this clinic or not
	_, clinicsDetailsErr := service.Repo.Searchclinic(ctx, bson.M{"ClinicID": appointmentDetails.Clinic, "doctorID": doctor.ID})
	if clinicsDetailsErr != nil {
		err := errors.New("this doctor is not onboarded to this clinic")
		return 0, utils.ReturnAppError(err, http.StatusNotFound, "Doctor Not Onboarded", err.Error())
	}

	//now add some defaults like status date
	appointmentDetails.RegistrationDate = time.Now()
	appointmentDetails.Status = "pending"
	appointmentDetails.DoctorName = doctor.Name

	appointmentDetails.ClinicName = clinic.ClinicDetails.Name
	appointmentDetails.ID = primitive.NewObjectID()

	slotNumber, appointmentErr := service.Repo.AddAppointment(ctx, clinic.ClinicDetails.MaxAppointments, appointmentDetails)
	if appointmentErr != nil {
		//here check for duplicate key error if yes then send error as MaxSlots booked
		if mongo.IsDuplicateKeyError(appointmentErr) {
			return 0, utils.ReturnAppError(errors.New("max Appointments Reached For Today"), http.StatusExpectationFailed, "Max Appointments Reached For today", "Max Appointments Reached For today")
		}
		return 0, utils.ReturnAppError(appointmentErr, http.StatusInternalServerError, "Failed to add appointment", appointmentErr.Error())
	}
	return slotNumber, nil
}

//AppointmentSlotsBooked function returns slots for whom maxAppointment has reached so all those slot documents where maxAppointment has reached are returned
func (service *clinicService) AppointmentSlotsBooked(ctx context.Context, slotDetais models.SlotDetails) ([]models.Slot, *structs.IAppError) {
	slots, err := service.Repo.AppointmentSlotsBooked(ctx, slotDetais.MaxAppointments, slotDetais.ClinicID, slotDetais.DoctorID)
	if err != nil {
		return nil, utils.ReturnAppError(err, http.StatusInternalServerError, "Failed to Fetch Booked Slots", err.Error())
	}

	return slots, nil
}

func (service *clinicService) DoctorWithItsclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, *structs.IAppError) {
	doctorsDetails, err := service.Repo.FetchDoctorAtclinics(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch  details",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	return doctorsDetails, nil
}
