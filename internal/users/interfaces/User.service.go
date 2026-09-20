package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/shared/response"
)

type IService interface {
	InitiateLoginOTP(ctx context.Context, MobileNumber string) (string, *response.IAppError)
	VerifyLoginOTP(ctx context.Context, loginCookie string, OTP string) (string, *response.IAppError)
}
