// Package middleware provides HTTP middleware for handling authentication and other common tasks.
package middleware

import (
	structs "AlShifa/Structs"
	utils "AlShifa/Utils"
	"context"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	//context related values
	ContextUserIDKey   contextKey = "userID"
	ContextUserRoleKey contextKey = "role"
)

func JwtAuthMiddleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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
		ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, ContextUserRoleKey, claims.Role)

		handler(w, r.WithContext(ctx))
	})
}
