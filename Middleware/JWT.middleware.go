// Package middleware provides HTTP middleware for handling authentication and other common tasks.
package middleware

import (
	"AlShifa/constants"
	structs "AlShifa/structs"
	utils "AlShifa/utils"
	"context"
	"net/http"
	"os"
	"strings"
)

func JwtAuthmiddleware(next http.Handler) http.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			_ = utils.WriteResponse(w, http.StatusBadRequest, structs.IAppError{
				StatusCode: http.StatusUnauthorized,
				Message:    "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			_ = utils.WriteResponse(w, http.StatusBadRequest, structs.IAppError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid Authentication Header",
			})
			return
		}

		tokenStr := parts[1]

		//use utility function to validate token
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {

			_ = utils.WriteResponse(w, http.StatusBadRequest, structs.IAppError{
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
