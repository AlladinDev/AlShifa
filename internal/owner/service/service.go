// Package service provides service functions for owner module
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/AlladinDev/AlShifa/internal/owner/dtos"
	"github.com/AlladinDev/AlShifa/internal/owner/interfaces"
	"github.com/AlladinDev/AlShifa/internal/owner/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	alshifaInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"

	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//
type userProfile struct {
	Email  string             `json:"email" bson:"email"`
	UserID primitive.ObjectID `json:"userID" bson:"userID"`
	Mobile string             `json:"mobile" bson:"mobile"`
}
type Service struct {
	repo            interfaces.IRepository
	notifier        alshifaInterfaces.INotifier
	passwordMatcher func(password string, hash string) (bool, error)
	fileUploader    alshifaInterfaces.IImageUploader
	calculateHash   func(data string) (hash string, err error)
	kvStore         alshifaInterfaces.Cache[string, []byte]
	otpGenerator    func() string
	verifyHash      func(rawToken string, hash string) (bool, error)
	decryptionFn    func(token string) ([]byte, error)
	encryptFn       func(data any) (string, error)
	generateJWT     func(claims *constants.JwtCustomClaims) (string, error)
}

func NewService(repo interfaces.IRepository, GenerateJWT func(claims *constants.JwtCustomClaims) (string, error), HashVerificationFn func(rawToken string, hash string) (bool, error), DecryptionFn func(token string) ([]byte, error), EncryptionFn func(data any) (string, error), KVStore alshifaInterfaces.Cache[string, []byte], HashingFn func(data string) (string, error), OTPGenerator func() string, FileUploader alshifaInterfaces.IImageUploader, Notifier alshifaInterfaces.INotifier, PasswordMatcher func(password string, hash string) (bool, error)) *Service {
	return &Service{
		repo:            repo,
		notifier:        Notifier,
		passwordMatcher: PasswordMatcher,
		fileUploader:    FileUploader,
		otpGenerator:    OTPGenerator,
		calculateHash:   HashingFn,
		decryptionFn:    DecryptionFn,
		kvStore:         KVStore,
		encryptFn:       EncryptionFn,
		generateJWT:     GenerateJWT,
		verifyHash:      HashVerificationFn,
	}
}

var _ interfaces.IService = (*Service)(nil)

func (s *Service) RegisterOwner(ctx context.Context, ownerDetails models.Owner, photoFile io.Reader) (string, *response.IAppError) {

	//first check whether this email or mobile already exists if yes throw error
	ownerExists, searchErr := s.repo.GetOwnerDetails(ctx, bson.M{"$or": []bson.M{
		{"email": ownerDetails.Email},
		{"mobile": ownerDetails.Mobile},
		{"aadhaarNumber": ownerDetails.AadhaarNumber},
	}})

	if searchErr != nil && searchErr != mongo.ErrNoDocuments {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     searchErr.Error(),
			ErrorObj:   searchErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if ownerExists != nil {
		return "", &response.IAppError{
			Message:    "Owner Already Exists with this email or mobile",
			Reason:     "This Email or mobile already exists",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	//add some default things like createdAt role
	ownerDetails.Role = constants.RoleclinicOwner
	ownerDetails.CreatedAt = time.Now()
	ownerDetails.EmailVerified = false
	ownerDetails.MobileVerified = false
	ownerDetails.AadhaarVerified = false
	ownerDetails.AccountOpeningAmountPaid = false
	ownerDetails.ID = primitive.NewObjectID()

	//hash the password now
	hashedPassword, hashingErr := utils.HashPasswordArgon2id(ownerDetails.Password)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now update the raw password with this hashedpassword
	ownerDetails.Password = hashedPassword

	//at this point save the user photo
	url, photoID, err := s.fileUploader.Upload(photoFile, "Owner Photos")
	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//update the ownerDetails with profile url and unique id of profile photo form uploader
	ownerDetails.Photo = url
	ownerDetails.PhotoID = photoID

	//now at this point create otp payload to sent to user
	encryptedOwnerDetails, encryptionErr := utils.Encrypt(ownerDetails)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//generate otp
	otp := utils.GenerateOTP()

	//hash the otp
	hashedOtp, hashingErr := utils.HashPasswordArgon2id(otp)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     hashingErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now create payload
	payload := dtos.OTPPayload{
		ExpiresAt: time.Now().Add(constants.OTPExpiry),
		Payload:   encryptedOwnerDetails,
		OTPHash:   hashedOtp,
	}

	//message to send to user
	otpMessage := fmt.Sprintf("Your AlShifa verification code is %s.This OTP will expire in %d minutes.", otp, 10)
	//now send otp to user
	if err := s.notifier.SendNotification(ownerDetails.Email, otpMessage, "Authentication", "Enter this OTP in AlShifa platform to create your account"); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now encrypt the payload
	encryptedPayload, encryptionErr := utils.Encrypt(payload)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to register owner",
			Reason:     encryptionErr.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return encryptedPayload, nil
}

func (s *Service) GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, *response.IAppError) {
	owner, err := s.repo.GetOwnerByID(ctx, ownerID)
	if err != nil {
		errMsg := "Failed to Fetch OwnerDetails"
		errStatusCode := http.StatusInternalServerError
		if err == mongo.ErrNoDocuments {
			errMsg = "Owner Doesnt Exist"
			errStatusCode = http.StatusNotFound
		}
		return nil, &response.IAppError{
			Message:    errMsg,
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: errStatusCode,
		}
	}

	//hide owner password dont sent it to frontend
	owner.Password = ""
	return owner, nil
}

func (s *Service) GetOwnerDetails(ctx context.Context, filters bson.M) (*models.Owner, *response.IAppError) {
	owner, err := s.repo.GetOwnerDetails(ctx, filters)
	if err != nil {
		return nil, &response.IAppError{
			Message:    "Failed to fetch ownerdetails",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return owner, nil
}

func (s *Service) VerifyOwnerRegisterOtp(ctx context.Context, otp string, otpPayload dtos.OTPPayload) (string, *response.IAppError) {
	//first check for expired otp
	if otpPayload.ExpiresAt.Before(time.Now()) {
		return "", &response.IAppError{
			Message:    "Otp expired ",
			ErrorObj:   nil,
			Reason:     "Otp is expired",
			StatusCode: http.StatusBadRequest,
		}
	}

	//now check for otp hash if it doesnt match return error
	hashMatches, hashingErr := utils.VerifyPasswordArgon2id(otp, otpPayload.OTPHash)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Failed To Verify User",
			ErrorObj:   hashingErr,
			Reason:     hashingErr.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	if !hashMatches {
		return "", &response.IAppError{
			Message:    "InCorrect otp",
			ErrorObj:   nil,
			Reason:     "incorrect otp sent",
			StatusCode: http.StatusBadRequest,
		}
	}

	//now as otp is correct decrypt the authpayload in it which contains user profile data
	decryiptedOwnerProfile, decryptionErr := utils.Decrypt(otpPayload.Payload)
	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Failed to authenticate User",
			ErrorObj:   decryptionErr,
			Reason:     decryptionErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now convert decrypted owner data into model form as it will be byte form
	var ownerData models.Owner

	if err := json.Unmarshal(decryiptedOwnerProfile, &ownerData); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to authenticate User",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now check here if email aadhaar mobile number already exists for another safety check
	findQuery := bson.M{"$or": []bson.M{
		{"email": ownerData.Email},
		{"mobile": ownerData.Mobile},
		{"aadhaarNumber": ownerData.AadhaarNumber},
	}}

	existingOwner, err := s.repo.GetOwnerDetails(ctx, findQuery)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", &response.IAppError{
				Message:    "Failed to authenticate User",
				ErrorObj:   err,
				Reason:     err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	if existingOwner != nil {
		return "", &response.IAppError{
			Message:    "Email or mobile or aadhaar number already exists ,login instead",
			ErrorObj:   nil,
			Reason:     "Duplicate unique fields login instead",
			StatusCode: http.StatusConflict,
		}
	}

	//now here set email verified as true
	ownerData.EmailVerified = true

	//now as here as everything is correct save the user
	if err := s.repo.RegisterOwner(ctx, ownerData); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to authenticate User",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now as everything is ok user has also been created so generate jwt token
	jwtToken, err := utils.GenerateJWT(&constants.JwtCustomClaims{
		UserID:     ownerData.ID.Hex(),
		Mobile:     ownerData.Mobile,
		Email:      ownerData.Email,
		IsVerified: true,
		Role:       ownerData.Role,
	})

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to authenticate User",
			ErrorObj:   err,
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwtToken, nil

}

func (s *Service) InitiateLogin(ctx context.Context, email string, password string) (temporaryLoginToken string, err *response.IAppError) {
	//first check whether this email exists or not
	owner, ownerSearchErr := s.repo.GetOwnerDetails(ctx, bson.M{"email": email})
	if ownerSearchErr != nil {
		if errors.Is(ownerSearchErr, mongo.ErrNoDocuments) {
			return "", &response.IAppError{
				Message:    "this email doesnt exist",
				Reason:     ownerSearchErr.Error(),
				ErrorObj:   ownerSearchErr,
				StatusCode: http.StatusNotFound,
			}
		} else {
			return "", &response.IAppError{
				Message:    "Login Failed",
				Reason:     ownerSearchErr.Error(),
				ErrorObj:   ownerSearchErr,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	//now check whether the password matches or not
	passwordMatches, matchingErr := s.passwordMatcher(password, owner.Password)
	if matchingErr != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     matchingErr.Error(),
			ErrorObj:   matchingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !passwordMatches {
		return "", &response.IAppError{
			Message:    "invalid email or password",
			Reason:     "password or email is incorrect",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//generate otp now
	otp := s.otpGenerator()

	//now hash this otp
	otpHashed, hashingErr := s.calculateHash(otp)
	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now make a payload for otp payload using user id useremail and mobile
	userProfiePayload := userProfile{
		UserID: owner.ID,
		Email:  owner.Email,
		Mobile: owner.Mobile,
	}

	userProfileEncrypted, encryptionErr := s.encryptFn(userProfiePayload)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now here as everything is correct make a payload containing this ownerid and also an otp hashed with expiry time
	otpPayload := dtos.OTPPayload{
		OTPHash:   otpHashed,
		ExpiresAt: time.Now().Add(constants.OTPExpiry),
		Payload:   userProfileEncrypted,
	}

	encryptedToken, encryptionErr := s.encryptFn(otpPayload)
	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now here send the otp to owner through messanger
	msg := fmt.Sprintf("Your otp for logging into your account is %s it will expire in %s minutes ", otp, constants.OTPExpiry)
	if err := s.notifier.SendNotification(owner.Email, msg, "OTP for login into AlShifa platform", "Do not share your otp with anyone except alshifa platform"); err != nil {
		return "", &response.IAppError{
			Message:    "Login Failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now send this otpPayload back
	return encryptedToken, nil
}

func (s *Service) VerifyLogin(ctx context.Context, temporaryLoginToken string, otp string) (jwtToken string, err *response.IAppError) {
	//first this temporary logintoken is encrypted so decrypt it back
	decryptedLoginToken, decryptionErr := s.decryptionFn(temporaryLoginToken)
	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Verification failed",
			Reason:     decryptionErr.Error(),
			ErrorObj:   decryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now parse this token into otp payload form
	otpPayload := dtos.OTPPayload{}
	if err := json.Unmarshal(decryptedLoginToken, &otpPayload); err != nil {
		return "", &response.IAppError{
			Message:    "Invalid login token format ,login failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now check for expiry of otp
	if time.Now().After(otpPayload.ExpiresAt) {
		return "", &response.IAppError{
			Message:    "Otp expired fetch a new otp and try again",
			Reason:     "otp expired",
			ErrorObj:   "otp expired",
			StatusCode: http.StatusBadRequest,
		}
	}

	//now check for otp hash comparison
	otpCorrect, verificationErr := s.verifyHash(otp, otpPayload.OTPHash)
	if verificationErr != nil {
		return "", &response.IAppError{
			Message:    "login failed",
			Reason:     verificationErr.Error(),
			ErrorObj:   verificationErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !otpCorrect {
		return "", &response.IAppError{
			Message:    "invalid otp",
			Reason:     "otp is incorrect",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	//now as everything is correct  get the user profile payload from otp payload
	decryptedData, decryptionErr := s.decryptionFn(otpPayload.Payload)
	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     decryptionErr.Error(),
			ErrorObj:   decryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now cast this decryptedData into user profile format
	ownerDetails := userProfile{}
	if err := json.Unmarshal(decryptedData, &ownerDetails); err != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	jwtPayload := constants.JwtCustomClaims{
		UserID: ownerDetails.UserID.Hex(),
		Email:  ownerDetails.Email,
		Mobile: ownerDetails.Mobile,
		Role:   constants.RoleclinicOwner,
	}

	jwt, jwtErr := s.generateJWT(&jwtPayload)
	if jwtErr != nil {
		return "", &response.IAppError{
			Message:    "Login failed",
			Reason:     jwtErr.Error(),
			ErrorObj:   jwtErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwt, nil
}
