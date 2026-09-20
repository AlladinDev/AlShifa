// Package middleware provides HTTP middleware for handling authentication and other common tasks.
package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
)

func JwtAuthmiddleware(next http.Handler) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, cookieErr := r.Cookie(constants.NameAuthToken)
		if cookieErr != nil {
			_ = utils.WriteResponse(w, http.StatusBadRequest, response.IAppError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid or no cookie found",
				ErrorObj:   cookieErr,
				Reason:     cookieErr.Error(),
			})
			return
		}

		if cookie.Value == "" {
			_ = utils.WriteResponse(w, http.StatusBadRequest, response.IAppError{
				StatusCode: http.StatusUnauthorized,
				Message:    "empty auth cookie",
				ErrorObj:   nil,
				Reason:     "invalid auth cookie it is empty",
			})
			return
		}

		//use utility function to validate token
		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {

			_ = utils.WriteResponse(w, http.StatusBadRequest, response.IAppError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid or expired token",
				ErrorObj:   err,
				Reason:     err.Error(),
			})
			return
		}

		// ---- Inject values into context ----

		ctx := context.WithValue(r.Context(), constants.KeyUserID, claims.UserID)

		ctx = context.WithValue(ctx, constants.KeyEmail, claims.Email)
		ctx = context.WithValue(ctx, constants.KeyMobile, claims.Mobile)
		ctx = context.WithValue(ctx, constants.KeyUserRole, claims.Role)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
