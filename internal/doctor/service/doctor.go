// Package doctorservice provides the core business logic for doctor registration,
// verification, and profile-related operations.
package doctorservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	doctorDtos "github.com/AlladinDev/AlShifa/internal/doctor/dtos"
	doctorinterfaces "github.com/AlladinDev/AlShifa/internal/doctor/interfaces"
	doctormodels "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"github.com/AlladinDev/AlShifa/internal/owner/dtos"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/customerrors"
	appInterface "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IService struct {
	repo          doctorinterfaces.IRepository
	OTPGenerator  func() string
	imageUploader appInterface.IImageUploader
	notifier      appInterface.INotifier
	decryptionFn  func(token string) ([]byte, error)
	clinicModule  doctorinterfaces.IClinic
	encryptionFn  func(data any) (string, error)
	jwtGenerator  func(claims *constants.JwtCustomClaims) (string, error)
	calculateHash func(data string) (string, error)
	hashVerifier  func(rawString string, hash string) (bool, error)
}

func NewService(Repo doctorinterfaces.IRepository,
	OTPGenerator func() string,
	ImageUploader appInterface.IImageUploader,
	Notifier appInterface.INotifier,
	EncryptionFn func(data any) (string, error),
	DecryptionFn func(token string) ([]byte, error),
	HashingFn func(data string) (string, error),
	ClinicModule doctorinterfaces.IClinic,
	GenerateJWT func(claims *constants.JwtCustomClaims) (string, error),
	HashVerifier func(rawString string, hash string) (bool, error),
) doctorinterfaces.IService {
	return &IService{
		repo:          Repo,
		OTPGenerator:  OTPGenerator,
		imageUploader: ImageUploader,
		notifier:      Notifier,
		clinicModule:  ClinicModule,
		hashVerifier:  HashVerifier,
		encryptionFn:  EncryptionFn,
		jwtGenerator:  GenerateJWT,
		calculateHash: HashingFn,
		decryptionFn:  DecryptionFn,
	}
}

func (s *IService) GetDoctorDetails(ctx context.Context, ID primitive.ObjectID) (*doctormodels.Doctor, *response.IAppError) {
	doctor, err := s.repo.FetchDoctor(ctx, bson.M{"_id": ID})
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed To fetch Doctor details",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return doctor, nil
}

func (s *IService) RegisterDoctor(ctx context.Context, doctorDetails *doctormodels.Doctor, profilePhotoFile *multipart.FileHeader, certificates []*multipart.FileHeader) (string, *response.IAppError) {
	//check if this doctor already exists
	existingDoctor, err := s.repo.FetchDoctor(ctx, bson.M{"$or": []bson.M{
		{"email": doctorDetails.Email},
		{"aadhaarNumber": doctorDetails.AadhaarNumber},
		{"mobile": doctorDetails.Mobile},
	}})
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", &response.IAppError{
				Message:    "Doctor Registration Failed",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	if existingDoctor != nil {
		return "", &response.IAppError{
			Message:    "Email or mobile already exists",
			Reason:     "Email or mobile Already Exists",
			ErrorObj:   errors.New("doctor already exists"),
			StatusCode: http.StatusConflict,
		}
	}

	// Upload profile photo
	if profilePhotoFile != nil {
		// Open the file to get io.Reader
		file, openErr := profilePhotoFile.Open()
		if openErr != nil {
			return "", &response.IAppError{
				Message:    "Doctor Registration Failed",
				Reason:     "Failed to open profile photo",
				ErrorObj:   openErr,
				StatusCode: http.StatusInternalServerError,
			}
		}
		defer file.Close()

		profilePhotoURL, profilePhotoID, profileUploadErr := s.imageUploader.Upload(file, "DoctorProfiles")
		if profileUploadErr != nil {
			return "", &response.IAppError{
				Message:    "Doctor Registration Failed",
				Reason:     profileUploadErr.Error(),
				ErrorObj:   profileUploadErr,
				StatusCode: http.StatusInternalServerError,
			}
		}

		doctorDetails.ProfilePhoto = profilePhotoURL
		doctorDetails.ProfilePhotoID = profilePhotoID
	}

	//upload certificates to cloudinary to get url and public id
	for _, certificate := range certificates {
		if certificate == nil {
			continue
		}

		// Open the certificate file to get io.Reader
		file, openErr := certificate.Open()
		if openErr != nil {
			return "", &response.IAppError{
				Message:    "Doctor Registration Failed",
				Reason:     "Failed to open certificate file",
				ErrorObj:   openErr,
				StatusCode: http.StatusInternalServerError,
			}
		}
		defer file.Close()

		certificateURL, certificateID, certificateUploadErr := s.imageUploader.Upload(file, "SpecializationCertificates")
		if certificateUploadErr != nil {
			return "", &response.IAppError{
				Message:    "Doctor Registration Failed",
				Reason:     certificateUploadErr.Error(),
				ErrorObj:   certificateUploadErr,
				StatusCode: http.StatusInternalServerError,
			}
		}

		certificatesStruct := doctormodels.TSpecializationCertificates{
			ID:  certificateID,
			URL: certificateURL,
		}
		doctorDetails.SpecializationCertificates = append(doctorDetails.SpecializationCertificates, certificatesStruct)
	}

	//hash the password first
	hashedPassword, hashingErr := utils.HashPasswordArgon2id(doctorDetails.Password)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//set password to hashed one
	doctorDetails.Password = hashedPassword

	//here add some defaults like id createdAt role
	doctorDetails.CreatedAt = time.Now()
	doctorDetails.Role = constants.RoleDoctor
	doctorDetails.ID = primitive.NewObjectID()

	//now add some boolean fields
	doctorDetails.AadhaarVerified = false
	doctorDetails.EmailVerified = false
	doctorDetails.ProfessionalDegreesVerified = false
	doctorDetails.MobileVerified = false
	doctorDetails.AccountOpeningAmountPaid = false
	doctorDetails.ClinicsAllowedToOnboard = []primitive.ObjectID{}

	//now at this point send otp to user email and then return a encrypted payload containing hashedOTP userID userType
	otp := s.OTPGenerator()
	otpHash, hashingErr := utils.HashPasswordArgon2id(otp)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//first send otp
	otpMessage := fmt.Sprintf("Your AlShifa verification code is %s.This OTP will expire in %d minutes.", otp, 10)
	notifierErr := s.notifier.SendNotification(doctorDetails.Email, otpMessage, "Email Verification", "Enter this OTP in the AlShifa application to complete verification.")
	if notifierErr != nil {
		return "", &response.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     notifierErr.Error(),
			ErrorObj:   notifierErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now encrypt full user model
	doctorDetailedEncrypted, encryptionErr := utils.Encrypt(doctorDetails)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now prepare encryption payload to be send with otp
	authPayload := doctorDtos.TOTPPayload{
		OTPHash:   otpHash,
		Payload:   doctorDetailedEncrypted,
		ExpiresAt: time.Now().Add(constants.OTPExpiry),
	}

	//now encrypt this authPayload
	authPayloadEncrypted, encryptionErr := utils.Encrypt(authPayload)

	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Doctor Registration Failed",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return authPayloadEncrypted, nil
}

func (s *IService) VerifyDoctorRegistrationOtp(ctx context.Context, otp string, otpPayload doctorDtos.TOTPPayload) (string, *response.IAppError) {
	//first check for expiry
	if time.Now().After(otpPayload.ExpiresAt) {
		return "", &response.IAppError{
			Message:    "Otp expired ",
			Reason:     "otp expired",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now check for hash
	hashMatches, err := utils.VerifyPasswordArgon2id(otp, otpPayload.OTPHash)
	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor ",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !hashMatches {
		return "", &response.IAppError{
			Message:    "Invalid Otp",
			Reason:     "invalid otp",
			ErrorObj:   hashMatches,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now as otp payload is  not expired decrypt the payload
	decryptedPayload, decryptionErr := utils.Decrypt(otpPayload.Payload)
	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     decryptionErr.Error(),
			ErrorObj:   decryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now try to convert decrypted payload into doctor model type
	var doctor doctormodels.Doctor
	if err := json.Unmarshal(decryptedPayload, &doctor); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now here again verify if this doctor is already registered because if doctor hits same api again it will be again registered
	doctorExists, existenceCheckErr := s.repo.DoctorExists(ctx, bson.M{"$or": []bson.M{
		{"email": doctor.Email},
		{"aadhaarNumber": doctor.AadhaarNumber},
		{"mobile": doctor.Mobile},
	}})

	if existenceCheckErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     existenceCheckErr.Error(),
			ErrorObj:   existenceCheckErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if doctorExists {
		return "", &response.IAppError{
			Message:    "Doctor is already registered with this email or mobile or aadhaar",
			Reason:     "duplicate entry is not possible",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	//as email is verified now so update the bool value
	doctor.EmailVerified = true

	if err := s.repo.RegisterDoctor(ctx, doctor); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now prepare jwt payload
	jwtPayload, err := utils.GenerateJWT(&constants.JwtCustomClaims{
		UserID:     doctor.ID.Hex(),
		Email:      doctor.Email,
		Mobile:     string(doctor.Mobile),
		Role:       doctor.Role,
		IsVerified: true,
	})
	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwtPayload, nil
}

func (s *IService) DoctorExists(ctx context.Context, filter bson.M) (bool, error) {
	return s.repo.DoctorExists(ctx, filter)
}

func (s *IService) FetchDoctors(ctx context.Context, filters bson.M) ([]doctormodels.Doctor, *response.IAppError) {
	doctors, err := s.repo.FetchDoctors(ctx, filters)
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed to fetch doctors",
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
			Reason:     err.Error(),
		}
	}

	return doctors, nil
}

func (s *IService) IsClinicAllowedToOnboardDoctor(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) (bool, error) {
	return s.repo.CheckClinicAllowedToOnboard(ctx, clinicID, doctorID)
}

func (s *IService) GetDoctorEmail(ctx context.Context, doctorID primitive.ObjectID) (string, error) {
	doctor, err := s.repo.FetchDoctor(ctx, bson.M{"_id": doctorID})
	if err != nil {
		return "", err
	}
	return doctor.Email, nil
}

func (s *IService) AllowClinicToOnboardDoctor(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *response.IAppError {
	//first thing here we have to check whether this clinic exists or not for that clinic module interface to communicate with clinic module
	clinicExists, clinicCheckingErr := s.clinicModule.ClinicExists(ctx, bson.M{"_id": clinicID})
	if clinicCheckingErr != nil {
		return &response.IAppError{
			Message:    "Failed to allow clinic to onboard doctor",
			Reason:     clinicCheckingErr.Error(),
			ErrorObj:   clinicCheckingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !clinicExists {
		return &response.IAppError{
			Message:    "No clinic exits with this clinic id",
			Reason:     "clinic not found",
			ErrorObj:   nil,
			StatusCode: http.StatusNotFound,
		}
	}

	if err := s.repo.AllowClinicToOnboard(ctx, clinicID, doctorID); err != nil {
		//now errors can be of three types it can be either no document found,no document modified due to some issues,or some other error
		if err == mongo.ErrNoDocuments {
			return &response.IAppError{
				Message:    "No record found with this doctor id",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		//document was found but failed to update
		if err == customerrors.ErrNoDocumentModified {
			return &response.IAppError{
				Message:    "Failed to update doctor document",
				Reason:     err.Error(),
				ErrorObj:   err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		//some other error
		return &response.IAppError{
			Message:    "Failed to update doctor document",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}

	}

	return nil
}

func (s *IService) InitiateLogin(ctx context.Context, loginPayload doctorDtos.LoginDTO) (string, *response.IAppError) {
	//here first check whether this email exists or not
	doctor, existenceCheckErr := s.repo.FetchDoctor(ctx, bson.M{"email": loginPayload.Email})
	if existenceCheckErr != nil {
		if errors.Is(existenceCheckErr, mongo.ErrNoDocuments) {
			return "", &response.IAppError{
				Message:    "This email doesnt exist",
				Reason:     existenceCheckErr.Error(),
				ErrorObj:   existenceCheckErr,
				StatusCode: http.StatusBadRequest,
			}
		}
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     existenceCheckErr.Error(),
			ErrorObj:   existenceCheckErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now check for whether the password matches or not
	passwordMatches, verificationErr := s.hashVerifier(loginPayload.Password, doctor.Password)
	if verificationErr != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     verificationErr.Error(),
			ErrorObj:   verificationErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !passwordMatches {
		return "", &response.IAppError{
			Message:    "Invalid email or password",
			Reason:     "credientials are not matching",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now as everything is ok
	//generate payload for otp
	otp := s.OTPGenerator()

	otphash, hashingErr := s.calculateHash(otp)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	otpPayload := dtos.OTPPayload{
		OTPHash:   otphash,
		ExpiresAt: time.Now().Add(constants.OTPExpiry),
		Payload:   doctor.ID.Hex(),
	}

	//now encrypt this otppayload
	otpPayloadEncrypted, encryptionErr := s.encryptionFn(otpPayload)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now at this point send the otp to doctor email
	msg := fmt.Sprintf("OTP for logging into your account is %s ,valid for %s ", otp, constants.OTPExpiry)
	if err := s.notifier.SendNotification(doctor.Email, msg, "OTP For Authentication", "Do not share this otp with anyone except on authentication forms of AlShifa Platform"); err != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return otpPayloadEncrypted, nil
}

func (s *IService) VerifyLogin(ctx context.Context, token string, otp string) (jwtToken string, err *response.IAppError) {
	//first decrypt the token
	decryptedToken, decryptionErr := s.decryptionFn(token)
	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     decryptionErr.Error(),
			ErrorObj:   decryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now try to convert this decryptedToken into otpPayload
	otpPayload := dtos.OTPPayload{}
	if err := json.Unmarshal(decryptedToken, &otpPayload); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now here check the expiry time
	if time.Now().After(otpPayload.ExpiresAt) {
		return "", &response.IAppError{
			Message:    "OTP Expired",
			Reason:     "otp is expired request new otp to continue",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now compare otp hash
	otpMatches, otpMatchingErr := s.hashVerifier(otp, otpPayload.OTPHash)
	if otpMatchingErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     otpMatchingErr.Error(),
			ErrorObj:   otpMatchingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !otpMatches {
		return "", &response.IAppError{
			Message:    "Incorrect otp provided",
			Reason:     "otp is incorrect",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now here as everything is correct prepare jwt payload for that get the doctor using doctorid in payload field of otpPayload dto
	doctorMongoID, idErr := primitive.ObjectIDFromHex(otpPayload.Payload)
	if idErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     idErr.Error(),
			ErrorObj:   idErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now using this doctormongodb id get the doctor details
	doctor, doctorFetchingErr := s.repo.FetchDoctor(ctx, bson.M{"_id": doctorMongoID})
	if doctorFetchingErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     doctorFetchingErr.Error(),
			ErrorObj:   doctorFetchingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	token, tokenErr := s.jwtGenerator(&constants.JwtCustomClaims{
		UserID:     doctor.ID.Hex(),
		Email:      doctor.Email,
		Role:       doctor.Role,
		Mobile:     string(doctor.Mobile),
		IsVerified: true,
	})
	if tokenErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify otp",
			Reason:     tokenErr.Error(),
			ErrorObj:   tokenErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}
