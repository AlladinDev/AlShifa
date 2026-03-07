// Package service contains service layer implementation for clinic module
package service

import (
	interfaces "AlShifa/clinic/interfaces"
	"AlShifa/clinic/models"
	"AlShifa/constants"
	appInterfaces "AlShifa/interfaces"
	structs "AlShifa/structs"
	utils "AlShifa/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type clinicService struct {
	OTPNotifier  appInterfaces.INotifier[string, string]
	OTPGenerator func() string
	Repo         interfaces.IRepository
	MongoClient  *mongo.Client

	//this is the cache which will be used to store otp when a clinic wants to add a doctor
	OtpCache appInterfaces.ICache[string, models.AddDoctorOtpPayload]
}

func NewclinicService(repo interfaces.IRepository, mongoClient *mongo.Client, OTPNotifier appInterfaces.INotifier[string, string],
	OTPGenerator func() string,
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

func (service *clinicService) DoctorExists(ctx context.Context, doctorID primitive.ObjectID) *structs.IAppError {
	if err := service.Repo.DoctorExists(ctx, doctorID); err != nil {
		return &structs.IAppError{
			Message:    "Failed to fetch doctor details",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}
	return nil
}

func (service *clinicService) ClinicExists(ctx context.Context, clinicID primitive.ObjectID) *structs.IAppError {
	if err := service.Repo.ClinicExists(ctx, clinicID); err != nil {
		return &structs.IAppError{
			Message:    "Failed to fetch doctor details",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}
	return nil

}

func (service *clinicService) FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, *structs.IAppError) {
	maxAppointments, err := service.Repo.FetchMaxAppointments(ctx, clinicID)
	if err != nil {
		return 0, &structs.IAppError{
			Message:    "Failed to fetch maxAppointments",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}

	return maxAppointments, nil
}

func (service *clinicService) DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) *structs.IAppError {
	if err := service.Repo.DoctorClinicMappingExists(ctx, clinicID, doctorID); err != nil {
		return &structs.IAppError{
			Message:    "Failed to check doctor clinic mapping",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}

	return nil
}
func (service *clinicService) FetchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, *structs.IAppError) {
	doctors, err := service.Repo.FetchDoctors(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch doctors",
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
			Reason:     err.Error(),
		}
	}
	return doctors, nil
}

func (service *clinicService) SearchDoctor(ctx context.Context, filter bson.M) (*models.Doctor, *structs.IAppError) {
	doctor, err := service.Repo.FetchDoctors(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch doctors",
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
			Reason:     err.Error(),
		}
	}
	if len(doctor) == 0 {
		return nil, &structs.IAppError{
			Message:    "No Doctor Found",
			StatusCode: http.StatusNotFound,
			ErrorObj:   errors.New("no doctor found with this filter"),
			Reason:     "no doctor found with this filter",
		}
	}
	return &doctor[0], nil
}

func (service *clinicService) ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string, maxAppointmentsAllowed int) {
	//checkout whether this clinic exists or not
	clinic, clinicSearchErr := service.Repo.SearchclinicByID(ctx, clinicID)
	if clinicSearchErr != nil {
		return &structs.IAppError{
			Message:    "Appointment  Booking Failed",
			Reason:     clinicSearchErr.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   clinicSearchErr,
		}, "", "", "", 0
	}

	clinicDoctors, err := service.Repo.FetchDoctorClinicMappings(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}, "", "", "", 0
	}

	if len(clinicDoctors) == 0 {
		return &structs.IAppError{
			Message:    "Appointment Booking Failed",
			Reason:     "This Doctor and clinic mapping doesnt exist",
			StatusCode: http.StatusNotFound,
			ErrorObj:   errors.New("this mapping doesnt exist"),
		}, "", "", "", 0
	}

	clinicDoctor := clinicDoctors[0]

	//now here check whether this doctor is available on requested appointment date
	requestedAppointmentDay := appointmentDateRequested.Weekday()

	appointmentDayPossible := false
	for _, day := range clinicDoctor.WorkingDays {
		if strings.EqualFold(day, requestedAppointmentDay.String()) {
			appointmentDayPossible = true
			break
		}
	}
	if !appointmentDayPossible {
		return &structs.IAppError{
			Message:    "Appointment Day Requested Not Possible",
			Reason:     "This doctor is not available on this day",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   errors.New("this doctor is not available on this day"),
		}, "", "", "", 0
	}

	return nil, clinicDoctor.DoctorName, clinicDoctor.ClinicName, clinicDoctor.ClinicAddress, clinic.MaxAppointments

}

//AddDoctorToclinic function generate otp for adding doctor to clinic this process needs another function or handler to verify otp and then the process is completed
func (service *clinicService) AddDoctorToclinic(ctx context.Context, userID primitive.ObjectID, clinicDetails models.ClinicDoctor) *structs.IAppError {
	//now check if clinic exists with this ClinicID

	clinic, err := service.Repo.SearchclinicByID(ctx, clinicDetails.ClinicID)
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinics",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//safety net here
	//here now check whether this owner owns this clinic by checking clinic detais
	if clinic.OwnerID != userID {
		return &structs.IAppError{
			Message:    "Owner doesnt owns this clinic",
			Reason:     "this clinic belongs to someone else",
			ErrorObj:   errors.New("owner doesnt own this clinic"),
			StatusCode: http.StatusForbidden,
		}
	}

	//now check whether this doctorID exists or not
	doctors, err := service.Repo.FetchDoctors(ctx, bson.M{"_id": clinicDetails.DoctorID})
	if err != nil {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//if length of doctors is 0 it means no doctor exists with this doctorid
	if len(doctors) == 0 {
		return &structs.IAppError{
			Message:    "Failed to Add Doctor To clinic",
			Reason:     "this doctor doesnt exist",
			ErrorObj:   errors.New("this doctor doesnt exist"),
			StatusCode: http.StatusNotFound,
		}
	}
	//get the first doctor as per this id
	doctor := doctors[0]

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
		return &structs.IAppError{
			Message:    "Doctor is already onboarded to this clinic",
			Reason:     errors.New("duplicate onboarding").Error(),
			ErrorObj:   errors.New("duplicate onboarding"),
			StatusCode: http.StatusForbidden,
		}
	}

	//now here extract the email from ctx
	doctorEmailAny := ctx.Value(constants.KeyEmail)
	doctorEmail, validDoctorEmail := doctorEmailAny.(string)
	if !validDoctorEmail {
		return &structs.IAppError{
			Message:    "Invalid Doctor Email",
			Reason:     "doctor Email is not valid",
			ErrorObj:   errors.New("doctor Email is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	//now send otp to doctor email through notification module and rest logic will be handled by otp verifier function
	//prepare otp payload
	uniquePrefix := userID.Hex()
	otpStr := service.OTPGenerator()
	otp := fmt.Sprintf("%s:%s", uniquePrefix, otpStr)
	otpPayload := models.AddDoctorOtpPayload{
		OTP:    otp,
		Expiry: time.Now().Add(utils.OTPExpiry),
		ClinicDetails: models.ClinicDoctor{
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
	if err := service.OtpCache.Set(ctx, otp, otpPayload, constants.CacheTTL); err != nil {
		return &structs.IAppError{
			Message:    "Failed To Add Doctor To clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusForbidden,
		}
	}

	//now send the otp to doctor email through notification module
	notifierErr := service.OTPNotifier.SendNotification(doctorEmail, otp)
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

func (service *clinicService) VerifyAddDoctorToclinicOTP(ctx context.Context, otp string, userID primitive.ObjectID) *structs.IAppError {
	//check whether in cache this otp exists or not
	otpToFind := fmt.Sprintf("%service:%service", userID.Hex(), otp)
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

	//now check whether otp expiry hasnt reached if yes return error
	if time.Now().After(otpPayload.Expiry) {
		return &structs.IAppError{
			Message:    "OTP Expired",
			ErrorObj:   errors.New("OTP Expired"),
			Reason:     errors.New("OTP Expired").Error(),
			StatusCode: 400,
		}
	}

	if err := service.Repo.AddDoctorToclinic(ctx, otpPayload.ClinicDetails); err != nil {
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
	// set default values
	clinicDetails.CreatedAt = time.Now().UTC()
	clinicDetails.Wallet = nil
	clinicDetails.ID = primitive.NewObjectID()
	clinicDetails.PlanType = utils.PlanPaid
	clinicDetails.OwnerID = ownerID
	registrationErr := service.Repo.Registerclinic(ctx, ownerID, clinicDetails)
	if registrationErr != nil {
		fmt.Print(registrationErr)
		return utils.ReturnAppError(registrationErr, 500, "Registration Failed", "Unknown reason")
	}

	return nil

}

func (service *clinicService) Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {

	clinics, err := service.Repo.Searchclinic(ctx, filter)
	if err != nil {
		return nil, utils.ReturnAppError(err, 500, "Unable To Fetch clinic details", err.Error())
	}

	return clinics, nil
}

func (service *clinicService) FetchDoctorClinicMappings(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError) {
	doctorsDetails, err := service.Repo.FetchDoctorClinicMappings(ctx, filter)
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

func (service *clinicService) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) *structs.IAppError {
	if err := service.Repo.RegisterDoctor(ctx, doctorDetails); err != nil {
		return &structs.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
