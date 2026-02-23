package service

import (
	"AlShifa/clinic/models"
	appInterfaces "AlShifa/interfaces"
	structs "AlShifa/structs"
	"context"
	"errors"
	"net/http"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoDBCacheKeys = map[string]primitive.ObjectID{}

func ReturnNewObjectID(key string) primitive.ObjectID {
	cachedOBjectID, exists := mongoDBCacheKeys[key]
	if !exists {
		mongoDBCacheKeys[key] = primitive.NewObjectID()
		return mongoDBCacheKeys[key]
	}
	return cachedOBjectID
}

func ResetMOngoDBCacheKeys() {
	mongoDBCacheKeys = map[string]primitive.ObjectID{}
}

func TestAddDoctorToclinicOTPGenerationFn(t *testing.T) {

	testCases := []struct {
		Name            string
		OtpGenerator    func(uniquePrefix string) string
		Notifier        appInterfaces.INotifier[string, string]
		mockRepoFn      *MockRepo
		shouldReturnErr bool
		clinicDetails   models.AddDoctorToclinic
		expectedErr     *structs.IAppError
		OTPCache        appInterfaces.ICache[string, models.AddDoctorOtpPayload]
	}{

		//test case when otp should be successfully sent as everything is ok
		{Name: "Successfully sent otp to doctor for onboarding",
			OtpGenerator:    func(uniquePrefix string) string { return "OTP" },
			Notifier:        ReturnNewEmailNotifier(),
			shouldReturnErr: false,
			clinicDetails: models.AddDoctorToclinic{
				DoctorID: ReturnNewObjectID("Dr Saqlain"),
				ClinicID: ReturnNewObjectID("AlMedrid clinic"),
			},
			expectedErr: nil,
			mockRepoFn: &MockRepo{
				SearchclinicMockFn: func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {
					return []models.ClinicDoctor{
						{
							ID: ReturnNewObjectID("AlMedrid clinic"),
						},
					}, nil
				},
				SearchDoctorMockFn: func(ctx context.Context, filter bson.M) (models.Doctor, error) {
					return models.Doctor{
						ID: ReturnNewObjectID("Dr Saqlain"),
					}, nil
				},
			},
			OTPCache: ReturnNewCache(),
		},

		//test case when otp generation fails because of error from searchclinicMockFn
		{Name: "Otp generation fn failed   to search doctor",
			OtpGenerator:    func(uniquePrefix string) string { return "OTP" },
			Notifier:        ReturnNewEmailNotifier(),
			shouldReturnErr: true,
			clinicDetails: models.AddDoctorToclinic{
				DoctorID: ReturnNewObjectID("Dr Rouf"),
				ClinicID: ReturnNewObjectID("AlMedrid clinic"),
			},
			mockRepoFn: &MockRepo{
				SearchclinicMockFn: func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {
					return nil, mongo.ErrNoDocuments
				},
				SearchDoctorMockFn: func(ctx context.Context, filter bson.M) (models.Doctor, error) {
					return models.Doctor{
						ID: ReturnNewObjectID("Dr Rouf"),
					}, nil
				},
			},
			expectedErr: &structs.IAppError{
				Message:    "Failed to Add Doctor To clinic",
				Reason:     mongo.ErrNoDocuments.Error(),
				ErrorObj:   mongo.ErrNoDocuments,
				StatusCode: http.StatusInternalServerError,
			},
			OTPCache: ReturnNewCache(),
		},

		//otp generation failed as doctor is already onboarded to clinic
		{Name: "Otp generation failed  because doctor is already onboarded",
			OtpGenerator:    func(uniquePrefix string) string { return "OTP" },
			Notifier:        ReturnNewEmailNotifier(),
			shouldReturnErr: true,
			clinicDetails: models.AddDoctorToclinic{
				DoctorID: ReturnNewObjectID("Dr Imran_AlreadyOnboarded"),
				ClinicID: ReturnNewObjectID("Alchem clinic"),
			},
			mockRepoFn: &MockRepo{
				SearchclinicMockFn: func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {
					return []models.ClinicDoctor{
						{
							ID: ReturnNewObjectID("Alchem clinic"),
						},
					}, nil
				},
				SearchDoctorMockFn: func(ctx context.Context, filter bson.M) (models.Doctor, error) {
					return models.Doctor{
						ID: ReturnNewObjectID("Dr Imran_AlreadyOnboarded"),
					}, nil
				},
			},
			expectedErr: &structs.IAppError{
				Message:    "Doctor Already Added To clinic",
				Reason:     errors.New("doctor is already added to clinic").Error(),
				ErrorObj:   errors.New("doctor is already added to clinic"),
				StatusCode: http.StatusForbidden,
			},
			OTPCache: ReturnNewCache(),
		},
	}

	for _, tc := range testCases {

		//first clear cache for keys
		ResetMOngoDBCacheKeys()

		service := NewclinicService(tc.mockRepoFn, nil, tc.Notifier, tc.OtpGenerator, tc.OTPCache)
		err := service.AddDoctorToclinic(context.Background(), tc.clinicDetails)
		if !tc.shouldReturnErr && err != nil {
			t.Fatalf("Test name %s Unexpected error %v", tc.Name, err)
		}

		if tc.shouldReturnErr && err == nil {
			t.Fatalf("Test name %s expected err but got no error expected error was %v", tc.Name, tc.expectedErr)
		}

	}
}
