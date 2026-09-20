// Package service provides service layer functions for user module
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/users/dtos"
	"github.com/AlladinDev/AlShifa/internal/users/interfaces"
	"github.com/AlladinDev/AlShifa/internal/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	otpGenerator func() string
	hashingFn    func(data string) (string, error)
	encryptionFn func(data any) (string, error)
	notifier     appInterfaces.INotifier
	decryptionFn func(token string) ([]byte, error)
	verifyHash   func(input string, hash string) (bool, error)
	generateJWT  func(token *constants.JwtCustomClaims) (string, error)
	repo         interfaces.IRepository
}

var _ interfaces.IService = (*Service)(nil)

func ReturnNewService(
	OTPGenerator func() string,
	HashingFn func(data string) (string, error),
	DecryptionFn func(token string) ([]byte, error),
	EncryptionFn func(data any) (string, error),
	Notifier appInterfaces.INotifier,
	VerifyHash func(input string, hash string) (bool, error),
	GenerateJWT func(claims *constants.JwtCustomClaims) (string, error),
	Repository interfaces.IRepository,
) *Service {
	return &Service{
		otpGenerator: OTPGenerator,
		hashingFn:    HashingFn,
		decryptionFn: DecryptionFn,
		encryptionFn: EncryptionFn,
		notifier:     Notifier,
		verifyHash:   VerifyHash,
		generateJWT:  GenerateJWT,
		repo:         Repository,
	}
}

// InitiateLoginOTP generates and sends a one-time password (OTP) to the
// provided mobile number and returns an encrypted payload that can later be
// used to verify the OTP.
//
// The complete flow is:
//
//	1. Generate a new OTP.
//	2. Hash the generated OTP.
//	3. Create an OTP payload containing:
//	   - the hashed OTP,
//	   - the expiration time,
//	   - the user's mobile number.
//	4. Encrypt the OTP payload.
//	5. Send the plain OTP to the user's mobile number through the notifier.
//	6. Return the encrypted payload.
//
// The plain OTP itself is never returned to the caller. It is only included
// in the notification sent to the user's mobile number.
//
// The encrypted payload contains the information required by the subsequent
// OTP verification process. Since the OTP is stored only in hashed form
// inside the payload, the verification process can later compare the OTP
// supplied by the user against the stored hash.
//
// The method follows a fail-fast approach. If any operation fails, the
// remaining operations are not executed and an IAppError is returned.
//
// Error flow:
//
//	OTP generation
//	    │
//	    ▼
//	OTP hashing ────────── failure ──► return error
//	    │
//	    ▼
//	Create OTP payload
//	    │
//	    ▼
//	Encrypt payload ────── failure ──► return error
//	    │
//	    ▼
//	Send OTP notification ─ failure ─► return error
//	    │
//	    ▼
//	Return encrypted payload
//
// The context is accepted so that the method follows the service-layer
// convention for request-scoped operations. Currently, the context is not
// directly used by the operations inside this method because the configured
// OTP generator, hashing function, encryption function, and notifier do not
// require it.
//
// Parameters:
//
//	ctx
//	    Context associated with the current request or operation.
//
//	mobileNumber
//	    Mobile number to which the OTP should be sent. This value is also
//	    stored inside the encrypted OTP payload so that the subsequent
//	    verification process can associate the OTP with the same mobile
//	    number.
//
// Returns:
//
//	string
//	    The encrypted OTP payload when all operations succeed. An empty
//	    string is returned when any operation fails.
//
//	*response.IAppError
//	    nil when the operation succeeds. Otherwise, an IAppError containing:
//	    - a generic failure message,
//	    - the underlying reason,
//	    - the original error object,
//	    - HTTP 500 Internal Server Error status.
//
// Security considerations:
//
//	- The plain OTP is never returned from the service.
//	- Only the hash of the OTP is placed inside the OTP payload.
//	- The OTP payload is encrypted before being returned.
//	- The plain OTP is sent only through the configured notification
//	  mechanism.
//	- If encryption fails, the OTP is not sent.
//	- If notification fails, the encrypted payload is not returned.
//
// The ordering of encryption before notification is intentional. This
// ensures that the caller receives a usable verification payload only when
// both payload creation and OTP delivery have succeeded.
func (s *Service) InitiateLoginOTP(
	ctx context.Context,
	mobileNumber string,
) (string, *response.IAppError) {

	// Generate a new OTP.
	//
	// The plain OTP is kept only in memory and is later sent to the user
	// through the notification service.
	otp := s.otpGenerator()

	// Hash the OTP before storing it inside the payload.
	//
	// The plain OTP should not be stored because anyone obtaining the
	// payload should not be able to directly retrieve the OTP.
	otpHash, hashingErr := s.hashingFn(otp)

	if hashingErr != nil {
		return "", &response.IAppError{
			Message:    "failed to login user",
			Reason:     hashingErr.Error(),
			ErrorObj:   hashingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Create the OTP payload.
	//
	// The payload contains the hashed OTP, the time at which the OTP
	// expires, and the mobile number associated with the OTP.
	otpPayload := dtos.OTPPayload{
		OTPHash:      otpHash,
		ExpiresAt:    time.Now().Add(constants.OTPExpiry),
		MobileNumber: mobileNumber,
	}

	// Encrypt the OTP payload before returning it to the caller.
	//
	// The encrypted payload will later be provided during OTP verification
	// so that the server can recover the hashed OTP, expiration time, and
	// mobile number.
	encryptedPayload, encryptionErr := s.encryptionFn(otpPayload)

	if encryptionErr != nil {
		return "", &response.IAppError{
			Message:    "failed to login user",
			Reason:     encryptionErr.Error(),
			ErrorObj:   encryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Construct the notification message containing the plain OTP.
	//
	// The plain OTP is intentionally used only here. It is not included in
	// the encrypted payload in plain form and is not returned to the caller.
	message := fmt.Sprintf(
		"otp to verify your mobile number is %s valid for %s minutes",
		otp,
		constants.OTPExpiry,
	)

	title := "Mobile number verification"
	info := "dont share this otp with anyone "

	// Send the OTP to the user's mobile number.
	//
	// If notification fails, do not return the encrypted payload because
	// the user has not successfully received the OTP required to verify
	// the login request.
	if err := s.notifier.SendNotification(
		mobileNumber,
		message,
		title,
		info,
	); err != nil {
		return "", &response.IAppError{
			Message:    "failed to login user",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// All operations completed successfully:
	//
	//	OTP generated
	//	↓
	//	OTP hashed
	//	↓
	//	Payload created
	//	↓
	//	Payload encrypted
	//	↓
	//	OTP sent
	//
	// Return the encrypted payload to the caller.
	return encryptedPayload, nil
}

// VerifyLoginOTP verifies the OTP provided by the user and completes the
// login process.
//
// The function performs two main responsibilities:
//
//	1. Verify that the OTP supplied by the user is valid.
//	2. Authenticate the user by finding or creating their user account and
//	   generating a JWT login token.
//
// The OTP verification flow is:
//
//	1. Decrypt the login cookie.
//	2. Deserialize the decrypted data into an OTPPayload.
//	3. Check whether the OTP has expired.
//	4. Compare the supplied OTP with the stored OTP hash.
//	5. Create the initial JWT claims using the verified mobile number.
//	6. Find an existing user using the mobile number.
//	7. If the user exists, use the existing user's ID and email.
//	8. If the user does not exist, create a new user containing the
//	   verified mobile number and a newly generated ID.
//	9. Put the appropriate user ID and email into the JWT claims.
//	10. Generate the signed JWT.
//	11. Return the JWT login token.
//
// Existing user flow:
//
//	The user has already registered previously.
//
//	OTP verification
//	     │
//	     ▼
//	Find user by mobile number
//	     │
//	     ▼
//	User exists
//	     │
//	     ├── Get existing user ID
//	     ├── Get existing user email
//	     │
//	     ▼
//	Build JWT claims
//	     │
//	     ▼
//	Generate JWT
//	     │
//	     ▼
//	Return login token
//
// New user flow:
//
//	The mobile number has been successfully verified but no account exists.
//
//	OTP verification
//	     │
//	     ▼
//	Find user by mobile number
//	     │
//	     ▼
//	mongo.ErrNoDocuments
//	     │
//	     ▼
//	Create new user
//	     │
//	     ├── Generate new user ID
//	     ├── Store mobile number
//	     ├── Set role to RoleUser
//	     ├── Leave optional fields empty
//	     │
//	     ▼
//	Store new user ID in JWT claims
//	     │
//	     ▼
//	Generate JWT
//	     │
//	     ▼
//	Return login token
//
// OTP verification:
//
//	The login cookie contains an encrypted OTPPayload. The payload contains
//	the OTP hash, expiration time, and mobile number.
//
//	The supplied OTP is never compared directly with the stored hash.
//	Instead, verifyHash is used to determine whether the supplied OTP
//	corresponds to the stored hash.
//
// The OTP is considered valid only when:
//
//	- the encrypted login cookie can be decrypted,
//	- the decrypted data is a valid OTPPayload,
//	- the OTP has not expired,
//	- the supplied OTP matches the stored OTP hash.
//
// User handling:
//
//	After successful OTP verification, the mobile number is considered
//	verified. The service then searches for an existing user using that
//	mobile number.
//
//	If the user exists:
//
//	- their existing ID is placed into the JWT,
//	- their existing email is placed into the JWT.
//
//	If the user does not exist:
//
//	- a new user is created,
//	- a new ObjectID is generated,
//	- the verified mobile number is stored,
//	- the role is set to RoleUser,
//	- optional fields such as name, email, and password are left empty,
//	- the newly generated user ID is placed into the JWT.
//
// JWT claims:
//
//	The JWT payload always contains:
//
//	- Role: RoleUser
//	- Mobile: verified mobile number
//	- IsVerified: true
//	- ID: existing or newly created user ID
//
//	For an existing user, Email contains the user's existing email address.
//	For a newly created user, Email remains empty until the user completes
//	their profile.
//
// Error handling:
//
//	The function uses fail-fast error handling. If any required operation
//	fails, the function immediately returns an appropriate IAppError.
//
//	Possible failures include:
//
//	- login cookie decryption failure,
//	- invalid OTP payload,
//	- expired OTP,
//	- OTP hash comparison failure,
//	- invalid OTP,
//	- user lookup failure,
//	- new user creation failure,
//	- JWT generation failure.
//
// Security considerations:
//
//	- The plain OTP is not stored in the login cookie.
//	- Only the OTP hash is stored in the OTP payload.
//	- The OTP must be verified before accessing or creating a user account.
//	- The OTP expiration time is checked before hash verification.
//	- The mobile number used for account lookup comes from the decrypted
//	  OTP payload rather than being independently supplied by the client.
//	- A JWT is generated only after successful OTP verification and
//	  successful user lookup or creation.
//
// Parameters:
//
//	ctx
//	    Context associated with the current request. It is passed to
//	    repository operations so database operations remain request-scoped.
//
//	loginCookie
//	    Encrypted OTP payload generated by InitiateLoginOTP.
//
//	OTP
//	    Plain OTP entered by the user.
//
// Returns:
//
//	string
//	    A signed JWT login token when OTP verification and user
//	    authentication succeed.
//
//	*response.IAppError
//	    nil when the login succeeds. Otherwise, an IAppError describing
//	    the operation that failed.
func (s *Service) VerifyLoginOTP(
	ctx context.Context,
	loginCookie string,
	OTP string,
) (string, *response.IAppError) {

	// First decrypt the login cookie.
	//
	// The cookie contains the encrypted OTP payload generated during
	// the OTP initiation process.
	decryptedToken, decryptionErr := s.decryptionFn(loginCookie)

	if decryptionErr != nil {
		return "", &response.IAppError{
			Message:    "failed to verify otp",
			Reason:     decryptionErr.Error(),
			ErrorObj:   decryptionErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Convert the decrypted data into its OTPPayload representation.
	//
	// OTPPayload contains:
	//
	//	OTPHash
	//	ExpiresAt
	//	MobileNumber
	var otpPayload dtos.OTPPayload

	if err := json.Unmarshal(decryptedToken, &otpPayload); err != nil {
		return "", &response.IAppError{
			Message:    "failed to verify otp",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Check whether the OTP has expired.
	//
	// An expired OTP cannot be used even if the supplied OTP is correct.
	if otpPayload.ExpiresAt.Before(time.Now()) {
		return "", &response.IAppError{
			Message:    "otp is expired request a new otp",
			Reason:     "otp is expired",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	// Compare the supplied OTP with the stored OTP hash.
	//
	// The plain OTP is not stored or compared directly.
	hashMatches, comparisonErr := s.verifyHash(
		OTP,
		otpPayload.OTPHash,
	)

	if comparisonErr != nil {
		return "", &response.IAppError{
			Message:    "failed to verify otp",
			Reason:     comparisonErr.Error(),
			ErrorObj:   comparisonErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Reject the request when the supplied OTP does not match.
	if !hashMatches {
		return "", &response.IAppError{
			Message:    "invalid otp provided",
			Reason:     "otp is invalid enter correct otp",
			ErrorObj:   nil,
			StatusCode: http.StatusBadRequest,
		}
	}

	// At this point the mobile number has been successfully verified.
	//
	// Create the initial JWT claims. The ID and email will be populated
	// after checking whether the user already exists.
	jwtPayload := &constants.JwtCustomClaims{
		Role:       constants.RoleUser,
		Mobile:     otpPayload.MobileNumber,
		Email:      "",
		IsVerified: true,
	}

	// Search for an existing user using the verified mobile number.
	existingUser, userFetchingErr := s.repo.GetUserDetails(
		ctx,
		bson.M{
			"mobile": otpPayload.MobileNumber,
		},
	)

	if userFetchingErr != nil {

		// If no user exists, create a new user.
		if userFetchingErr == mongo.ErrNoDocuments {

			// Create a minimal user account.
			//
			// The mobile number has already been verified through the OTP
			// process, so it can safely be associated with the new account.
			newUser := models.User{
				Name:      "",
				Role:      constants.RoleUser,
				Email:     "",
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				Mobile:    otpPayload.MobileNumber,
				Password:  "",
			}

			// Persist the newly created user.
			if err := s.repo.CreateUser(ctx, &newUser); err != nil {
				return "", &response.IAppError{
					Message:    "Failed to verify user",
					Reason:     err.Error(),
					ErrorObj:   err,
					StatusCode: http.StatusInternalServerError,
				}
			}

			// Store the newly created user's ID in the JWT claims.
			//
			// Email remains empty because the new user has not completed
			// their profile yet.
			jwtPayload.ID = newUser.ID.String()

		} else {

			// Any repository error other than "user not found" means
			// that the user could not be safely retrieved.
			return "", &response.IAppError{
				Message:    "Failed to verify user",
				Reason:     userFetchingErr.Error(),
				ErrorObj:   userFetchingErr,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	// If the JWT ID is still empty, the user already existed.
	//
	// Use the existing user's ID and email in the JWT payload.
	//
	// For a newly created user, jwtPayload.ID was already populated above,
	// so this block is skipped.
	if jwtPayload.ID == "" {
		jwtPayload.ID = existingUser.ID.Hex()
		jwtPayload.Email = existingUser.Email
	}

	// Generate the actual signed JWT from the completed claims.
	token, err := s.generateJWT(jwtPayload)

	if err != nil {
		return "", &response.IAppError{
			Message:    "failed to verify otp",
			Reason:     err.Error(),
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	// OTP verification and user authentication are complete.
	//
	// Return the JWT login token to the caller.
	return token, nil
}
