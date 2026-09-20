// Package service contains unit tests for the user authentication service.
//
// This file tests:
//
//   - InitiateLoginOTP
//   - VerifyLoginOTP
//
// The tests use:
//
//   - table-driven tests
//   - reusable configurable mocks
//   - functional options
//   - isolated dependencies
//   - side-effect assertions
//   - JWT claim assertions
//   - repository assertions
//   - fail-fast behaviour verification
//   - error/status verification
//   - context propagation checks
//
// No real MongoDB, SMS provider, encryption service, or JWT service is
// required to execute these tests.
package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"github.com/AlladinDev/AlShifa/internal/users/dtos"
	"github.com/AlladinDev/AlShifa/internal/users/interfaces"
	"github.com/AlladinDev/AlShifa/internal/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	testMobile = "9797798243"
	testOTP    = "123456"
	testJWT    = "test-jwt"
)

// =============================================================================
// MOCK NOTIFIER
// =============================================================================

// MockMessenger is a configurable mock implementation of INotifier.
//
// SendNotificationFunc allows individual tests to control notification
// behaviour and inspect the arguments passed by the service.
type MockMessenger struct {
	SendNotificationFunc func(
		channel string,
		message string,
		title string,
		info string,
	) error
}

// SendNotification implements INotifier.
//
// If SendNotificationFunc is configured, it delegates to that function.
// Otherwise the notification succeeds automatically.
func (m *MockMessenger) SendNotification(
	channel string,
	message string,
	title string,
	info string,
) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(
			channel,
			message,
			title,
			info,
		)
	}

	return nil
}

// =============================================================================
// MOCK REPOSITORY
// =============================================================================

// MockRepository is a configurable mock implementation of IRepository.
//
// Each repository operation can be independently controlled by assigning
// a function field.
type MockRepository struct {
	CreateUserFunc func(
		ctx context.Context,
		user *models.User,
	) error

	GetUserDetailsFunc func(
		ctx context.Context,
		filter bson.M,
	) (*models.User, error)
}

// CreateUser implements IRepository.CreateUser.
func (r *MockRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	if r.CreateUserFunc != nil {
		return r.CreateUserFunc(ctx, user)
	}

	return nil
}

// GetUserDetails implements IRepository.GetUserDetails.
func (r *MockRepository) GetUserDetails(
	ctx context.Context,
	filter bson.M,
) (*models.User, error) {
	if r.GetUserDetailsFunc != nil {
		return r.GetUserDetailsFunc(ctx, filter)
	}

	// Default behaviour is an existing user.
	return &models.User{
		ID:     primitive.NewObjectID(),
		Name:   "Test User",
		Email:  "test@example.com",
		Mobile: testMobile,
		Role:   constants.RoleUser,
	}, nil
}

// =============================================================================
// MOCK SERVICE OPTIONS
// =============================================================================

// MockServiceOptions contains every replaceable dependency used by Service.
//
// Functional options allow tests to replace only the dependency relevant
// to a particular test.
type MockServiceOptions struct {
	GenerateOTP func() string

	Hash    func(string) (string, error)
	Decrypt func(string) ([]byte, error)
	Encrypt func(any) (string, error)

	VerifyHash  func(string, string) (bool, error)
	GenerateJWT func(*constants.JwtCustomClaims) (string, error)

	Messenger  appInterfaces.INotifier
	Repository interfaces.IRepository
}

// ReturnNewMockService creates a Service with safe default dependencies.
//
// Individual tests can override dependencies using functional options.
func ReturnNewMockService(
	opts ...func(*MockServiceOptions),
) *Service {

	options := MockServiceOptions{
		GenerateOTP: func() string {
			return testOTP
		},

		Hash: func(data string) (string, error) {
			return "hashed-otp", nil
		},

		Decrypt: func(token string) ([]byte, error) {
			return utils.Decrypt(token)
		},

		Encrypt: func(data any) (string, error) {
			return utils.Encrypt(data)
		},

		VerifyHash: func(
			input string,
			hash string,
		) (bool, error) {
			return true, nil
		},

		GenerateJWT: func(
			claims *constants.JwtCustomClaims,
		) (string, error) {
			return testJWT, nil
		},

		Messenger: &MockMessenger{},

		Repository: &MockRepository{},
	}

	for _, opt := range opts {
		opt(&options)
	}

	return ReturnNewService(
		options.GenerateOTP,
		options.Hash,
		options.Decrypt,
		options.Encrypt,
		options.Messenger,
		options.VerifyHash,
		options.GenerateJWT,
		options.Repository,
	)
}

// =============================================================================
// FUNCTIONAL OPTIONS
// =============================================================================

// WithGenerateOTP replaces the OTP generator.
func WithGenerateOTP(
	fn func() string,
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.GenerateOTP = fn
	}
}

// WithHashing replaces the hashing implementation.
func WithHashing(
	fn func(string) (string, error),
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.Hash = fn
	}
}

// WithDecryption replaces the decryption implementation.
func WithDecryption(
	fn func(string) ([]byte, error),
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.Decrypt = fn
	}
}

// WithEncryption replaces the encryption implementation.
func WithEncryption(
	fn func(any) (string, error),
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.Encrypt = fn
	}
}

// WithVerifyHash replaces OTP hash verification.
func WithVerifyHash(
	fn func(string, string) (bool, error),
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.VerifyHash = fn
	}
}

// WithGenerateJWT replaces JWT generation.
func WithGenerateJWT(
	fn func(*constants.JwtCustomClaims) (string, error),
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.GenerateJWT = fn
	}
}

// WithMessenger replaces the notification service.
func WithMessenger(
	messenger appInterfaces.INotifier,
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.Messenger = messenger
	}
}

// WithRepository replaces the repository.
func WithRepository(
	repository interfaces.IRepository,
) func(*MockServiceOptions) {
	return func(options *MockServiceOptions) {
		options.Repository = repository
	}
}

// =============================================================================
// TEST HELPERS
// =============================================================================

// assertAppError verifies the common IAppError contract.
//
// This helper keeps individual tests concise and ensures that every error
// test checks the same important fields.
func assertAppError(
	t *testing.T,
	err *response.IAppError,
	expectedReason string,
	expectedStatus int,
	expectedErrorObj error,
) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected application error, got nil")
	}

	if err.Reason != expectedReason {
		t.Errorf(
			"expected reason %q, got %q",
			expectedReason,
			err.Reason,
		)
	}

	if err.StatusCode != expectedStatus {
		t.Errorf(
			"expected status %d, got %d",
			expectedStatus,
			err.StatusCode,
		)
	}

	if expectedErrorObj != nil && err.ErrorObj != expectedErrorObj {
		t.Errorf(
			"expected ErrorObj %v, got %v",
			expectedErrorObj,
			err.ErrorObj,
		)
	}
}

// createEncryptedOTPPayload creates a realistic encrypted OTP payload.
//
// This allows VerifyLoginOTP to be tested independently without having to
// call InitiateLoginOTP first.
func createEncryptedOTPPayload(
	t *testing.T,
	mobile string,
	otp string,
	expiresAt time.Time,
) string {
	t.Helper()

	hash, err := utils.HashPasswordArgon2id(otp)
	if err != nil {
		t.Fatalf(
			"failed to hash test OTP: %v",
			err,
		)
	}

	payload := dtos.OTPPayload{
		MobileNumber: mobile,
		OTPHash:      hash,
		ExpiresAt:    expiresAt,
	}

	encrypted, err := utils.Encrypt(payload)
	if err != nil {
		t.Fatalf(
			"failed to encrypt test OTP payload: %v",
			err,
		)
	}

	return encrypted
}

// =============================================================================
// INITIATE LOGIN OTP
// =============================================================================

// Test_InitiateLoginOTP verifies the complete OTP initiation flow.
//
// Expected flow:
//
//	Generate OTP
//	    ↓
//	Hash OTP
//	    ↓
//	Create OTP payload
//	    ↓
//	Encrypt payload
//	    ↓
//	Send notification
//	    ↓
//	Return encrypted payload
//
// The test also verifies fail-fast behaviour.
func Test_InitiateLoginOTP(t *testing.T) {
	hashErr := errors.New("hashing error")
	encryptErr := errors.New("encryption error")
	notifyErr := errors.New("notification error")

	tests := []struct {
		name string

		hashFn     func(string) (string, error)
		encryptFn  func(any) (string, error)
		notifyFunc func(string, string, string, string) error

		expectedOutput string
		expectedError  error
	}{
		{
			name: "success",

			expectedOutput: "encrypted",
			expectedError:  nil,

			hashFn: func(data string) (string, error) {
				if data != testOTP {
					t.Errorf(
						"expected OTP %q, got %q",
						testOTP,
						data,
					)
				}

				return "hashed", nil
			},

			encryptFn: func(data any) (string, error) {
				payload, ok := data.(dtos.OTPPayload)
				if !ok {
					t.Fatalf(
						"expected OTPPayload, got %T",
						data,
					)
				}

				if payload.MobileNumber != testMobile {
					t.Errorf(
						"expected mobile %q, got %q",
						testMobile,
						payload.MobileNumber,
					)
				}

				if payload.OTPHash != "hashed" {
					t.Errorf(
						"expected hash %q, got %q",
						"hashed",
						payload.OTPHash,
					)
				}

				if !payload.ExpiresAt.After(time.Now()) {
					t.Error(
						"expected OTP expiration to be in the future",
					)
				}

				return "encrypted", nil
			},

			notifyFunc: func(
				channel string,
				message string,
				title string,
				info string,
			) error {
				if channel != testMobile {
					t.Errorf(
						"expected channel %q, got %q",
						testMobile,
						channel,
					)
				}

				if title != "Mobile number verification" {
					t.Errorf(
						"unexpected title %q",
						title,
					)
				}

				if message == "" {
					t.Error("expected notification message")
				}

				if info != "dont share this otp with anyone " {
					t.Errorf(
						"unexpected info %q",
						info,
					)
				}

				return nil
			},
		},

		{
			name: "hashing fails",

			hashFn: func(data string) (string, error) {
				return "", hashErr
			},

			encryptFn: func(data any) (string, error) {
				t.Fatal(
					"encryption must not be called after hashing failure",
				)

				return "", nil
			},

			notifyFunc: func(
				channel string,
				message string,
				title string,
				info string,
			) error {
				t.Fatal(
					"notification must not be called after hashing failure",
				)

				return nil
			},

			expectedError: hashErr,
		},

		{
			name: "encryption fails",

			hashFn: func(data string) (string, error) {
				return "hashed", nil
			},

			encryptFn: func(data any) (string, error) {
				return "", encryptErr
			},

			notifyFunc: func(
				channel string,
				message string,
				title string,
				info string,
			) error {
				t.Fatal(
					"notification must not be called after encryption failure",
				)

				return nil
			},

			expectedError: encryptErr,
		},

		{
			name: "notification fails",

			hashFn: func(data string) (string, error) {
				return "hashed", nil
			},

			encryptFn: func(data any) (string, error) {
				return "encrypted", nil
			},

			notifyFunc: func(
				channel string,
				message string,
				title string,
				info string,
			) error {
				return notifyErr
			},

			expectedError: notifyErr,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messenger := &MockMessenger{
				SendNotificationFunc: tt.notifyFunc,
			}

			service := ReturnNewMockService(
				WithHashing(tt.hashFn),
				WithEncryption(tt.encryptFn),
				WithMessenger(messenger),
			)

			output, appErr := service.InitiateLoginOTP(
				t.Context(),
				testMobile,
			)

			if tt.expectedError == nil {
				if appErr != nil {
					t.Fatalf(
						"expected success, got %+v",
						appErr,
					)
				}

				if output != tt.expectedOutput {
					t.Fatalf(
						"expected output %q, got %q",
						tt.expectedOutput,
						output,
					)
				}

				return
			}

			if output != "" {
				t.Errorf(
					"expected empty output on failure, got %q",
					output,
				)
			}

			assertAppError(
				t,
				appErr,
				tt.expectedError.Error(),
				http.StatusInternalServerError,
				tt.expectedError,
			)
		})
	}
}

// =============================================================================
// VERIFY LOGIN OTP
// =============================================================================

// Test_VerifyLoginOTP tests OTP verification, user lookup/creation and
// JWT generation.
//
// Cases covered:
//
//   - successful existing user login
//   - decryption failure
//   - malformed JSON
//   - expired OTP
//   - hash verification failure
//   - invalid OTP
//   - repository failure
//   - new user creation
//   - new user creation failure
//   - JWT generation failure
func Test_VerifyLoginOTP(t *testing.T) {
	decryptionErr := errors.New("decryption error")
	verifyHashErr := errors.New("hash comparison error")
	repositoryErr := errors.New("database error")
	createUserErr := errors.New("create user error")
	jwtErr := errors.New("jwt generation error")

	tests := []struct {
		name string

		cookieOTP string
		inputOTP  string

		serviceOptions []func(*MockServiceOptions)

		expectedOutput string
		expectedReason string
		expectedStatus int
		expectedObj    error
	}{
		{
			name: "success existing user",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP:       testOTP,
			expectedOutput: testJWT,
		},

		{
			name: "decryption fails",

			cookieOTP: "invalid-cookie",
			inputOTP:  testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithDecryption(
					func(token string) ([]byte, error) {
						return nil, decryptionErr
					},
				),
			},

			expectedReason: decryptionErr.Error(),
			expectedStatus: http.StatusInternalServerError,
			expectedObj:    decryptionErr,
		},

		{
			name: "decrypted payload contains invalid JSON",

			cookieOTP: "encrypted-cookie",
			inputOTP:  testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithDecryption(
					func(token string) ([]byte, error) {
						return []byte(`invalid-json`), nil
					},
				),
			},

			expectedReason: "invalid character 'i' looking for beginning of value",
			expectedStatus: http.StatusInternalServerError,
		},

		{
			name: "OTP is expired",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(-5*time.Minute),
			),

			inputOTP:       testOTP,
			expectedReason: "otp is expired",
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "OTP hash verification fails",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithVerifyHash(
					func(
						input string,
						hash string,
					) (bool, error) {
						return false, verifyHashErr
					},
				),
			},

			expectedReason: verifyHashErr.Error(),
			expectedStatus: http.StatusInternalServerError,
			expectedObj:    verifyHashErr,
		},

		{
			name: "OTP is invalid",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: "999999",

			serviceOptions: []func(*MockServiceOptions){
				WithVerifyHash(
					func(
						input string,
						hash string,
					) (bool, error) {
						return false, nil
					},
				),
			},

			expectedReason: "otp is invalid enter correct otp",
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "repository lookup fails",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithRepository(
					&MockRepository{
						GetUserDetailsFunc: func(
							ctx context.Context,
							filter bson.M,
						) (*models.User, error) {
							return nil, repositoryErr
						},
					},
				),
			},

			expectedReason: repositoryErr.Error(),
			expectedStatus: http.StatusInternalServerError,
			expectedObj:    repositoryErr,
		},

		{
			name: "creates new user",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithRepository(
					&MockRepository{
						GetUserDetailsFunc: func(
							ctx context.Context,
							filter bson.M,
						) (*models.User, error) {
							return nil, mongo.ErrNoDocuments
						},

						CreateUserFunc: func(
							ctx context.Context,
							user *models.User,
						) error {
							if user.Mobile != testMobile {
								t.Errorf(
									"expected mobile %q, got %q",
									testMobile,
									user.Mobile,
								)
							}

							if user.Role != constants.RoleUser {
								t.Errorf(
									"expected role %v, got %v",
									constants.RoleUser,
									user.Role,
								)
							}

							if user.ID.IsZero() {
								t.Error(
									"expected generated user ID",
								)
							}

							return nil
						},
					},
				),
			},

			expectedOutput: testJWT,
		},

		{
			name: "new user creation fails",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithRepository(
					&MockRepository{
						GetUserDetailsFunc: func(
							ctx context.Context,
							filter bson.M,
						) (*models.User, error) {
							return nil, mongo.ErrNoDocuments
						},

						CreateUserFunc: func(
							ctx context.Context,
							user *models.User,
						) error {
							return createUserErr
						},
					},
				),
			},

			expectedReason: createUserErr.Error(),
			expectedStatus: http.StatusInternalServerError,
			expectedObj:    createUserErr,
		},

		{
			name: "JWT generation fails",

			cookieOTP: createEncryptedOTPPayload(
				t,
				testMobile,
				testOTP,
				time.Now().Add(5*time.Minute),
			),

			inputOTP: testOTP,

			serviceOptions: []func(*MockServiceOptions){
				WithGenerateJWT(
					func(
						claims *constants.JwtCustomClaims,
					) (string, error) {
						return "", jwtErr
					},
				),
			},

			expectedReason: jwtErr.Error(),
			expectedStatus: http.StatusInternalServerError,
			expectedObj:    jwtErr,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			options := []func(*MockServiceOptions){
				WithVerifyHash(
					func(
						input string,
						hash string,
					) (bool, error) {
						return true, nil
					},
				),

				WithRepository(
					&MockRepository{},
				),

				WithGenerateJWT(
					func(
						claims *constants.JwtCustomClaims,
					) (string, error) {
						return testJWT, nil
					},
				),
			}

			options = append(
				options,
				tt.serviceOptions...,
			)

			service := ReturnNewMockService(options...)

			output, appErr := service.VerifyLoginOTP(
				t.Context(),
				tt.cookieOTP,
				tt.inputOTP,
			)

			if tt.expectedReason == "" {
				if appErr != nil {
					t.Fatalf(
						"expected success, got %+v",
						appErr,
					)
				}

				if output != tt.expectedOutput {
					t.Fatalf(
						"expected output %q, got %q",
						tt.expectedOutput,
						output,
					)
				}

				return
			}

			if output != "" {
				t.Errorf(
					"expected empty output on failure, got %q",
					output,
				)
			}

			assertAppError(
				t,
				appErr,
				tt.expectedReason,
				tt.expectedStatus,
				tt.expectedObj,
			)
		})
	}
}

// =============================================================================
// VERIFY LOGIN OTP - VERIFIED MOBILE
// =============================================================================

// Test_VerifyLoginOTP_UsesVerifiedMobileNumber verifies that the mobile
// number used for repository lookup comes from the decrypted OTP payload.
//
// This prevents the authentication flow from trusting a different mobile
// number supplied by the client.
func Test_VerifyLoginOTP_UsesVerifiedMobileNumber(t *testing.T) {
	verifiedMobile := "9000000000"

	cookie := createEncryptedOTPPayload(
		t,
		verifiedMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	repositoryCalled := false

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {
			repositoryCalled = true

			if filter["mobile"] != verifiedMobile {
				t.Fatalf(
					"expected repository mobile %q, got %v",
					verifiedMobile,
					filter["mobile"],
				)
			}

			return &models.User{
				ID:     primitive.NewObjectID(),
				Email:  "test@example.com",
				Mobile: verifiedMobile,
				Role:   constants.RoleUser,
			}, nil
		},
	}

	service := ReturnNewMockService(
		WithRepository(repository),
	)

	_, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}

	if !repositoryCalled {
		t.Error(
			"expected repository to be called",
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - EXISTING USER
// =============================================================================

// Test_VerifyLoginOTP_DoesNotCreateExistingUser verifies that an existing
// user is not recreated.
func Test_VerifyLoginOTP_DoesNotCreateExistingUser(t *testing.T) {
	userID := primitive.NewObjectID()

	createCalled := false

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {
			return &models.User{
				ID:     userID,
				Email:  "existing@example.com",
				Mobile: testMobile,
				Role:   constants.RoleUser,
			}, nil
		},

		CreateUserFunc: func(
			ctx context.Context,
			user *models.User,
		) error {
			createCalled = true
			return nil
		},
	}

	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	service := ReturnNewMockService(
		WithRepository(repository),
	)

	_, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}

	if createCalled {
		t.Error(
			"CreateUser should not be called for an existing user",
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - NEW USER JWT CLAIMS
// =============================================================================

// Test_VerifyLoginOTP_NewUserJWTClaims verifies that a newly created user's
// ID is correctly placed into the JWT claims.
func Test_VerifyLoginOTP_NewUserJWTClaims(t *testing.T) {
	var createdUser *models.User

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {
			return nil, mongo.ErrNoDocuments
		},

		CreateUserFunc: func(
			ctx context.Context,
			user *models.User,
		) error {
			copy := *user
			createdUser = &copy
			return nil
		},
	}

	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	service := ReturnNewMockService(
		WithRepository(repository),
		WithGenerateJWT(
			func(
				claims *constants.JwtCustomClaims,
			) (string, error) {

				if createdUser == nil {
					t.Fatal(
						"expected user to be created before JWT generation",
					)
				}

				if claims.ID != createdUser.ID.Hex() {
					t.Errorf(
						"expected JWT ID %q, got %q",
						createdUser.ID.Hex(),
						claims.ID,
					)
				}

				if claims.Mobile != testMobile {
					t.Errorf(
						"expected mobile %q, got %q",
						testMobile,
						claims.Mobile,
					)
				}

				if claims.Role != constants.RoleUser {
					t.Errorf(
						"expected role %v, got %v",
						constants.RoleUser,
						claims.Role,
					)
				}

				if !claims.IsVerified {
					t.Error(
						"expected IsVerified=true",
					)
				}

				if claims.Email != "" {
					t.Errorf(
						"expected empty email for new user, got %q",
						claims.Email,
					)
				}

				return testJWT, nil
			},
		),
	)

	output, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}

	if output != testJWT {
		t.Errorf(
			"expected JWT %q, got %q",
			testJWT,
			output,
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - EXISTING USER JWT CLAIMS
// =============================================================================

// Test_VerifyLoginOTP_ExistingUserJWTClaims verifies that an existing user's
// ID and email are correctly included in the JWT claims.
func Test_VerifyLoginOTP_ExistingUserJWTClaims(t *testing.T) {
	userID := primitive.NewObjectID()

	const email = "existing@example.com"

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {
			return &models.User{
				ID:     userID,
				Email:  email,
				Mobile: testMobile,
				Role:   constants.RoleUser,
			}, nil
		},
	}

	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	service := ReturnNewMockService(
		WithRepository(repository),

		WithGenerateJWT(
			func(
				claims *constants.JwtCustomClaims,
			) (string, error) {

				if claims.ID != userID.Hex() {
					t.Errorf(
						"expected ID %q, got %q",
						userID.Hex(),
						claims.ID,
					)
				}

				if claims.Email != email {
					t.Errorf(
						"expected email %q, got %q",
						email,
						claims.Email,
					)
				}

				if claims.Mobile != testMobile {
					t.Errorf(
						"expected mobile %q, got %q",
						testMobile,
						claims.Mobile,
					)
				}

				if claims.Role != constants.RoleUser {
					t.Errorf(
						"expected role %v, got %v",
						constants.RoleUser,
						claims.Role,
					)
				}

				if !claims.IsVerified {
					t.Error(
						"expected IsVerified=true",
					)
				}

				return testJWT, nil
			},
		),
	)

	output, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}

	if output != testJWT {
		t.Errorf(
			"expected JWT %q, got %q",
			testJWT,
			output,
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - HASH INPUT
// =============================================================================

// Test_VerifyLoginOTP_PassesCorrectOTPToHashVerifier verifies that the OTP
// entered by the user and the OTP hash from the encrypted payload are passed
// correctly to verifyHash.
func Test_VerifyLoginOTP_PassesCorrectOTPToHashVerifier(t *testing.T) {
	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	hashCalled := false

	service := ReturnNewMockService(
		WithVerifyHash(
			func(
				input string,
				hash string,
			) (bool, error) {

				hashCalled = true

				if input != testOTP {
					t.Errorf(
						"expected OTP %q, got %q",
						testOTP,
						input,
					)
				}

				if hash == "" {
					t.Error(
						"expected stored OTP hash to be non-empty",
					)
				}

				return true, nil
			},
		),
	)

	_, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}

	if !hashCalled {
		t.Error(
			"expected verifyHash to be called",
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - CONTEXT PROPAGATION
// =============================================================================

// Test_VerifyLoginOTP_PropagatesContext verifies that the context received
// by VerifyLoginOTP is passed unchanged to repository operations.
func Test_VerifyLoginOTP_PropagatesContext(t *testing.T) {
	type contextKey string

	const key contextKey = "test-key"

	expectedValue := "test-value"

	ctx := context.WithValue(
		context.Background(),
		key,
		expectedValue,
	)

	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {

			if ctx.Value(key) != expectedValue {
				t.Errorf(
					"expected context value %q, got %v",
					expectedValue,
					ctx.Value(key),
				)
			}

			return &models.User{
				ID:     primitive.NewObjectID(),
				Email:  "test@example.com",
				Mobile: testMobile,
				Role:   constants.RoleUser,
			}, nil
		},
	}

	service := ReturnNewMockService(
		WithRepository(repository),
	)

	_, err := service.VerifyLoginOTP(
		ctx,
		cookie,
		testOTP,
	)

	if err != nil {
		t.Fatalf(
			"expected success, got %v",
			err,
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - FAIL FAST
// =============================================================================

// Test_VerifyLoginOTP_DoesNotCallRepositoryForInvalidOTP verifies that the
// service stops immediately when OTP verification fails.
func Test_VerifyLoginOTP_DoesNotCallRepositoryForInvalidOTP(t *testing.T) {
	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(5*time.Minute),
	)

	repositoryCalled := false

	repository := &MockRepository{
		GetUserDetailsFunc: func(
			ctx context.Context,
			filter bson.M,
		) (*models.User, error) {

			repositoryCalled = true

			return &models.User{
				ID: primitive.NewObjectID(),
			}, nil
		},
	}

	service := ReturnNewMockService(
		WithRepository(repository),

		WithVerifyHash(
			func(
				input string,
				hash string,
			) (bool, error) {
				return false, nil
			},
		),
	)

	_, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		"999999",
	)

	if err == nil {
		t.Fatal(
			"expected invalid OTP error",
		)
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			err.StatusCode,
		)
	}

	if repositoryCalled {
		t.Error(
			"repository should not be called when OTP is invalid",
		)
	}
}

// =============================================================================
// VERIFY LOGIN OTP - EXPIRED OTP
// =============================================================================

// Test_VerifyLoginOTP_DoesNotVerifyHashForExpiredOTP verifies that the
// expensive hash verification operation is skipped when the OTP has expired.
func Test_VerifyLoginOTP_DoesNotVerifyHashForExpiredOTP(t *testing.T) {
	cookie := createEncryptedOTPPayload(
		t,
		testMobile,
		testOTP,
		time.Now().Add(-5*time.Minute),
	)

	hashCalled := false

	service := ReturnNewMockService(
		WithVerifyHash(
			func(
				input string,
				hash string,
			) (bool, error) {

				hashCalled = true

				return true, nil
			},
		),
	)

	_, err := service.VerifyLoginOTP(
		t.Context(),
		cookie,
		testOTP,
	)

	if err == nil {
		t.Fatal(
			"expected expired OTP error",
		)
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			err.StatusCode,
		)
	}

	if hashCalled {
		t.Error(
			"hash verification should not be called for expired OTP",
		)
	}
}
