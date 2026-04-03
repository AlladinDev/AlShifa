package middleware

import (
	"net/http"
	"slices"

	"github.com/AlladinDev/AlShifa/constants"
	utils "github.com/AlladinDev/AlShifa/utils"
)

func RoleGuardmiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	if len(allowedRoles) == 0 {
		panic("At least one role must be specified for RoleGuardmiddleware")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			userRole := r.Context().Value(constants.KeyUserRole)

			if userRole == nil {
				_ = utils.WriteResponse(w, http.StatusBadRequest,
					utils.ReturnAppError(nil, 400, "Missing Role", "Missing Role"))
				return
			}

			userRoleStr, ok := userRole.(string)
			if !ok || userRoleStr == "" {
				_ = utils.WriteResponse(w, http.StatusBadRequest,
					utils.ReturnAppError(nil, 400, "Invalid Role Type", "Invalid Role Type"))
				return
			}

			// Check if user role is allowed
			if slices.Contains(allowedRoles, userRoleStr) {
				next.ServeHTTP(w, r) // ✅ IMPORTANT: call next
				return
			}

			_ = utils.WriteResponse(w, http.StatusForbidden,
				utils.ReturnAppError(nil, 403, "Forbidden To Access This Api only  can access it", "Forbidden"))
		})
	}
}
