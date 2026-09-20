package interfaces

import (
	"context"
	"io"

	"github.com/AlladinDev/AlShifa/internal/owner/dtos"
	"github.com/AlladinDev/AlShifa/internal/owner/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	RegisterOwner(ctx context.Context, ownerDetails models.Owner, photoFile io.Reader) (string, *response.IAppError)
	VerifyOwnerRegisterOtp(ctx context.Context, otp string, otpPayload dtos.OTPPayload) (string, *response.IAppError)
	GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, *response.IAppError)
	GetOwnerDetails(ctx context.Context, filter bson.M) (*models.Owner, *response.IAppError)
	InitiateLogin(ctx context.Context, email string, password string) (temporaryLoginToken string, err *response.IAppError)
	VerifyLogin(ctx context.Context, temporaryLoginToken string, otp string) (jwtToken string, err *response.IAppError)
}
