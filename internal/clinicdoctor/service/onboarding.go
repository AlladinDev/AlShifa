package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/dtos"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterface "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo             interfaces.IRepo
	doctorModule     interfaces.IDoctorModule
	clinicModule     interfaces.IClinicModule
	OTPGenerator     func() string
	notifier         appInterface.INotifier
	hashingFn        func(data string) (string, error)
	hashComparisonFn func(firstHash string, secondHash string) (bool, error)
	encrypt          func(data any) (string, error)
	decrypt          func(data string) ([]byte, error)
}

func NewOnboardingService(Repository interfaces.IRepo,
	Notifier appInterface.INotifier,
	EncryptionFn func(data any) (string, error),
	DecryptionFn func(data string) ([]byte, error),
	HashingFn func(data string) (string, error),
	otpGenerator func() string,
	ClinicModule interfaces.IClinicModule,
	DoctorModule interfaces.IDoctorModule,
	HashingComparisonFn func(firstHash string,
		secondHash string) (bool, error)) interfaces.IDoctorOnboarding {
	return &Service{
		repo:             Repository,
		doctorModule:     DoctorModule,
		clinicModule:     ClinicModule,
		OTPGenerator:     otpGenerator,
		hashingFn:        HashingFn,
		encrypt:          EncryptionFn,
		decrypt:          DecryptionFn,
		notifier:         Notifier,
		hashComparisonFn: HashingComparisonFn,
	}
}

// InitiateDoctorClinicOnboarding function initiates the process for onboarding
// ============================================================================
// INITIATE DOCTOR CLINIC ONBOARDING
// ============================================================================
//
// This method STARTS the onboarding process.
//
// IMPORTANT:
// It does NOT create the clinic-doctor mapping.
//
// Instead:
//
//	Clinic requests onboarding
//	        ↓
//	Validate doctor + clinic
//	        ↓
//	Check doctor's permission
//	        ↓
//	Generate OTP
//	        ↓
//	Hash OTP
//	        ↓
//	Create encrypted pending payload
//	        ↓
//	Send OTP to doctor
//	        ↓
//	Return encrypted onboarding token
//
// The actual mapping is created only by VerifyDoctorClinicOnboarding after
// the doctor successfully proves ownership of the OTP.
//
// ============================================================================
func (s *Service) InitiateDoctorClinicOnboarding(
	ctx context.Context,
	details *models.ClinicDoctorMapping,
	userID primitive.ObjectID,
) (string, *response.IAppError) {

	// -------------------------------------------------------------------------
	// 1. Verify that the doctor exists.
	// -------------------------------------------------------------------------
	//
	// Exists check has three possible outcomes:
	//
	//	true, nil   -> doctor exists
	//	false, nil  -> doctor does not exist
	//	false, err  -> database/service failure
	//
	// We keep "doesn't exist" separate from "check failed".

	doctorExists, err := s.doctorModule.DoctorExists(
		ctx,
		bson.M{"_id": details.DoctorID},
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !doctorExists {
		return "", &response.IAppError{
			Message:    "This doctor doesn't exist",
			Reason:     "doctor was not found",
			ErrorObj:   nil,
			StatusCode: http.StatusNotFound,
		}
	}

	// -------------------------------------------------------------------------
	// 2. Resolve the clinic from the authenticated user.
	// -------------------------------------------------------------------------
	//
	// The caller does NOT provide the clinic ID.
	//
	// We derive it from userID so that a clinic owner/receptionist can only
	// initiate onboarding for their own clinic.
	//
	// This prevents:
	//
	//	user A → authenticated to clinic A
	//	user A → sends clinic B ID
	//
	// from onboarding a doctor into clinic B.

	clinicID, err := s.clinicModule.GetClinicIDByOwnerID(
		ctx,
		userID,
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to resolve clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// The clinic ID is now server-trusted.
	details.ClinicID = clinicID

	// -------------------------------------------------------------------------
	// 3. Verify that the resolved clinic exists.
	// -------------------------------------------------------------------------

	clinicExists, err := s.clinicModule.ClinicExists(
		ctx,
		bson.M{"_id": details.ClinicID},
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to verify clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !clinicExists {
		return "", &response.IAppError{
			Message:    "This clinic doesn't exist",
			Reason:     "clinic was not found",
			ErrorObj:   nil,
			StatusCode: http.StatusNotFound,
		}
	}

	// -------------------------------------------------------------------------
	// 4. Make sure the doctor isn't already associated with this clinic.
	// -------------------------------------------------------------------------
	//
	// We don't want to start an onboarding flow for an already-existing
	// relationship.
	//
	//	true, nil  -> mapping already exists
	//	false, nil -> mapping doesn't exist
	//	false, err -> check failed

	mappingExists, err := s.repo.CheckClinicDoctorMappingExists(
		ctx,
		details.DoctorID,
		details.ClinicID,
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to check clinic doctor mapping",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if mappingExists {
		return "", &response.IAppError{
			Message:    "Clinic Doctor Mapping already exists",
			Reason:     "doctor is already associated with this clinic",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	// -------------------------------------------------------------------------
	// 5. Verify that the doctor has allowed this clinic to onboard them.
	// -------------------------------------------------------------------------
	//
	// This is the doctor's consent check.
	//
	// We perform it BEFORE generating or sending an OTP because there is no
	// reason to start an onboarding flow that the doctor has not permitted.

	fmt.Printf("%+v", details)
	allowed, err := s.doctorModule.IsClinicAllowedToOnboardDoctor(
		ctx,
		details.ClinicID,
		details.DoctorID,
	)
	if err != nil {
		// Any other error means the permission state could not be determined.
		return "", &response.IAppError{
			Message:    "Failed to verify onboarding permission",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if !allowed {
		return "", &response.IAppError{
			Message:    "Failed to onboard tell doctor to allow clinic to onboard this doctor",
			Reason:     "permission is needed from doctor for onboarding",
			ErrorObj:   nil,
			StatusCode: http.StatusForbidden,
		}
	}

	// -------------------------------------------------------------------------
	// 6. Get the doctor's email.
	// -------------------------------------------------------------------------
	//
	// The OTP must be sent directly to the doctor.

	doctorEmail, err := s.doctorModule.GetDoctorEmail(
		ctx,
		details.DoctorID,
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to get doctor email",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 7. Get the clinic name.
	// -------------------------------------------------------------------------
	//
	// The clinic name is included in the OTP notification so the doctor knows
	// which clinic is requesting onboarding.

	clinicName, err := s.clinicModule.GetClinicName(
		ctx,
		details.ClinicID,
	)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to get clinic name",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 8. Generate the OTP.
	// -------------------------------------------------------------------------
	//
	// The plaintext OTP exists only in memory.
	// It will be sent to the doctor but never stored in the payload.

	otp := s.OTPGenerator()

	// -------------------------------------------------------------------------
	// 9. Hash the OTP.
	// -------------------------------------------------------------------------
	//
	// We store the hash rather than the plaintext OTP.
	//
	// During verification:
	//
	//	submitted OTP
	//	    ↓
	//	hash
	//	    ↓
	//	compare with OTPHash
	//
	// This means the encrypted payload never contains the actual OTP.

	hashedOTP, err := s.hashingFn(otp)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to secure OTP",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 10. Create the pending onboarding payload.
	// -------------------------------------------------------------------------
	//
	// This payload represents:
	//
	//	"Someone has requested onboarding, but the doctor has not verified it."
	//
	// Nothing has been persisted to the clinic-doctor mapping collection yet.

	onboardingPayload := dtos.DoctorOnboardingOTPPayload{
		OTPHash:   hashedOTP,
		ExpiresAt: time.Now().Add(constants.OTPExpiry),
		Payload:   details,
	}

	// -------------------------------------------------------------------------
	// 11. Encrypt the pending onboarding payload.
	// -------------------------------------------------------------------------
	//
	// The payload contains sensitive information:
	//
	//	- Doctor ID
	//	- Clinic ID
	//	- OTP hash
	//	- Expiration time
	//	- Mapping information
	//
	// Encrypting it prevents the client from reading or modifying these
	// values before verification.

	encryptedPayload, err := s.encrypt(onboardingPayload)

	if err != nil {
		return "", &response.IAppError{
			Message:    "Failed to create onboarding token",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 12. Send the plaintext OTP to the doctor.
	// -------------------------------------------------------------------------
	//
	// IMPORTANT:
	//
	// The OTP is sent through the doctor's registered email.
	// The encrypted payload contains only the OTP hash.

	message := fmt.Sprintf(
		"OTP for onboarding to %s clinic is %s. If you did not expect this request, you can safely ignore it.",
		clinicName,
		otp,
	)

	title := "Doctor Clinic Onboarding"

	info := "Do not share this OTP with anyone. Enter it only in the AlShifa platform."

	if err := s.notifier.SendNotification(
		doctorEmail,
		message,
		title,
		info,
	); err != nil {
		return "", &response.IAppError{
			Message:    "Failed to send onboarding OTP",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 13. Return the encrypted pending onboarding token.
	// -------------------------------------------------------------------------
	//
	// At this point:
	//
	//	Doctor validated       ✓
	//	Clinic validated       ✓
	//	Mapping doesn't exist  ✓
	//	Doctor allowed it      ✓
	//	OTP generated          ✓
	//	OTP hashed             ✓
	//	Payload encrypted      ✓
	//	OTP sent               ✓
	//
	// But:
	//
	//	Clinic-doctor mapping  ✗ NOT CREATED
	//
	// The verification method must perform the final database write.

	return encryptedPayload, nil
}

// VerifyDoctorClinicOnboarding function verifies the otp and creates the doctor clinic mapping in database
// ============================================================================
// VERIFY DOCTOR CLINIC ONBOARDING
// ============================================================================
//
// This method FINISHES the onboarding process.
//
// It receives:
//
//  1. OTP entered by the doctor
//  2. Encrypted onboarding token
//  3. Authenticated user's ID
//
// Flow:
//
//	Encrypted token + OTP
//	        ↓
//	Decrypt token
//	        ↓
//	Check expiry
//	        ↓
//	Compare otp
//	        ↓
//	Resolve authenticated clinic
//	        ↓
//	Verify clinic matches token
//	        ↓
//	Verify doctor exists
//	        ↓
//	Verify doctor's permission still exists
//	        ↓
//	Verify mapping doesn't exist
//	        ↓
//	Create clinic-doctor mapping
//
// Only after every check succeeds is the mapping persisted.
//
// ============================================================================
func (s *Service) VerifyDoctorClinicOnboarding(
	ctx context.Context,
	OTP string,
	onboardingToken string,
	userID primitive.ObjectID,
) *response.IAppError {

	// -------------------------------------------------------------------------
	// 1. Decrypt the onboarding token.
	// -------------------------------------------------------------------------
	//
	// The token contains the temporary onboarding state created by
	// InitiateDoctorClinicOnboarding.

	decryptedToken, err := s.decrypt(onboardingToken)

	if err != nil {
		return &response.IAppError{
			Message:    "Failed to decrypt onboarding token",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 2. Convert decrypted data into the onboarding payload.
	// -------------------------------------------------------------------------

	otpPayload := dtos.DoctorOnboardingOTPPayload{}

	if err := json.Unmarshal(decryptedToken, &otpPayload); err != nil {
		return &response.IAppError{
			Message:    "Invalid onboarding token",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusBadRequest,
		}
	}

	// -------------------------------------------------------------------------
	// 3. Check whether the OTP has expired.
	// -------------------------------------------------------------------------
	//
	// Even a correct OTP must be rejected after its expiration time.

	if time.Now().After(otpPayload.ExpiresAt) {
		return &response.IAppError{
			Message:    "OTP expired, request a new OTP",
			Reason:     "OTP expiration time has passed",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	// -------------------------------------------------------------------------
	// 5. Compare the submitted OTP with the stored OTP hash.
	// -------------------------------------------------------------------------

	otpMatches, err := s.hashComparisonFn(
		OTP,
		otpPayload.OTPHash,
	)

	if err != nil {
		return &response.IAppError{
			Message:    "Failed to verify OTP",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Wrong OTP means the onboarding stops immediately.
	if !otpMatches {
		return &response.IAppError{
			Message:    "Incorrect OTP",
			Reason:     "submitted OTP does not match the onboarding OTP",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	// -------------------------------------------------------------------------
	// 6. Resolve the clinic belonging to the authenticated user.
	// -------------------------------------------------------------------------
	//
	// We NEVER trust the clinic ID supplied by the client.
	//
	// Instead:
	//
	//	authenticated user
	//	       ↓
	//	GetClinicIDByUserID
	//	       ↓
	//	actual clinic
	//
	// This clinic must match the clinic stored in the encrypted token.

	clinicID, err := s.clinicModule.GetClinicIDByOwnerID(
		ctx,
		userID,
	)

	if err != nil {
		return &response.IAppError{
			Message:    "Failed to resolve clinic",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 7. Make sure the authenticated user belongs to the target clinic.
	// -------------------------------------------------------------------------
	//
	// This prevents a valid onboarding token from being used by another clinic.

	if clinicID != otpPayload.Payload.ClinicID {
		return &response.IAppError{
			Message:  "Failed to onboard doctor",
			Reason:   "authenticated clinic does not match onboarding clinic",
			ErrorObj: nil,

			// This is an authorization failure.
			StatusCode: http.StatusForbidden,
		}
	}

	// -------------------------------------------------------------------------
	// 8. Verify that the doctor still exists.
	// -------------------------------------------------------------------------
	//
	// The doctor may have been deleted after initiation but before OTP
	// verification, so we check again before creating the mapping.

	doctorExists, err := s.doctorModule.DoctorExists(
		ctx,
		bson.M{"_id": otpPayload.Payload.DoctorID},
	)

	if err != nil {
		return &response.IAppError{
			Message:    "Failed to verify doctor",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !doctorExists {
		return &response.IAppError{
			Message:    "This doctor doesn't exist",
			Reason:     "doctor no longer exists",
			ErrorObj:   nil,
			StatusCode: http.StatusNotFound,
		}
	}

	// -------------------------------------------------------------------------
	// 9. Verify that the doctor still allows this clinic to onboard them.
	// -------------------------------------------------------------------------
	//
	// This is checked again because permission could have been revoked after
	// the OTP was generated.

	allowed, err := s.doctorModule.IsClinicAllowedToOnboardDoctor(
		ctx,
		otpPayload.Payload.ClinicID,
		otpPayload.Payload.DoctorID,
	)
	if err != nil {
		return &response.IAppError{
			Message:    "Failed to verify onboarding permission",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if !allowed {
		return &response.IAppError{
			Message:    "Onboarding not allowed; ask the doctor to enable it",
			Reason:     "doctor has not allowed the clinic to onboard him",
			ErrorObj:   nil,
			StatusCode: http.StatusForbidden,
		}
	}

	// -------------------------------------------------------------------------
	// 10. Check whether the mapping already exists.
	// -------------------------------------------------------------------------
	//
	// We check again because the state may have changed since initiation.
	//
	// Example:
	//
	//	Initiation
	//	    ↓
	//	Mapping doesn't exist
	//	    ↓
	//	Doctor enters OTP
	//	    ↓
	//	Another request creates mapping
	//	    ↓
	//	This verification request arrives
	//
	// Therefore the mapping must be checked immediately before insertion.

	mappingExists, err := s.repo.CheckClinicDoctorMappingExists(
		ctx,
		otpPayload.Payload.DoctorID,
		otpPayload.Payload.ClinicID,
	)

	if err != nil {
		return &response.IAppError{
			Message:    "Failed to check clinic doctor mapping",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if mappingExists {
		return &response.IAppError{
			Message:    "Onboarding failed; mapping already exists",
			Reason:     "doctor is already associated with this clinic",
			ErrorObj:   nil,
			StatusCode: http.StatusConflict,
		}
	}

	// -------------------------------------------------------------------------
	// 11. Finally create the clinic-doctor mapping.
	// -------------------------------------------------------------------------
	//
	// We have now passed every validation:
	//
	//	✓ Token valid
	//	✓ Payload valid
	//	✓ OTP not expired
	//	✓ OTP correct
	//	✓ Authenticated user belongs to target clinic
	//	✓ Doctor exists
	//	✓ Doctor still allows onboarding
	//	✓ Mapping does not already exist
	//
	// ONLY NOW do we create the relationship in the database.

	//here set some default things like objectid createdat etc
	otpPayload.Payload.ID = primitive.NewObjectID()
	otpPayload.Payload.CreatedAt = time.Now()

	if err := s.repo.RegisterDoctorClinicMapping(
		ctx,
		otpPayload.Payload,
	); err != nil {
		return &response.IAppError{
			Message:    "Failed to create clinic doctor mapping",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// -------------------------------------------------------------------------
	// 12. Onboarding completed successfully.
	// -------------------------------------------------------------------------
	//
	// nil means the doctor has successfully been onboarded to the clinic.

	return nil
}
