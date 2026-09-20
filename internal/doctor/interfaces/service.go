// Package doctorinterfaces defines the service contracts for doctor-related operations.
package doctorinterfaces

import (
	"context"
	"mime/multipart"

	doctordtos "github.com/AlladinDev/AlShifa/internal/doctor/dtos"
	models "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	DoctorExists(ctx context.Context, filter bson.M) (bool, error)
	GetDoctorDetails(ctx context.Context, ID primitive.ObjectID) (*models.Doctor, *response.IAppError)
	RegisterDoctor(ctx context.Context, doctorDetails *models.Doctor, profilePhoto *multipart.FileHeader, certificates []*multipart.FileHeader) (token string, error *response.IAppError)
	VerifyDoctorRegistrationOtp(ctx context.Context, OTP string, OTPPayload doctordtos.TOTPPayload) (string, *response.IAppError)
	FetchDoctors(ctx context.Context, filters bson.M) ([]models.Doctor, *response.IAppError)
	IsClinicAllowedToOnboardDoctor(ctx context.Context, clinicID primitive.ObjectID, ownerID primitive.ObjectID) (bool, error)
	GetDoctorEmail(ctx context.Context, doctorID primitive.ObjectID) (string, error)
	InitiateLogin(ctx context.Context, loginPayload doctordtos.LoginDTO) (string, *response.IAppError)
	VerifyLogin(ctx context.Context, otp string, token string) (jwtToken string, err *response.IAppError)
	AllowClinicToOnboardDoctor(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *response.IAppError
}
