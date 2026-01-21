package middleware

import (
	utils "AlShifa/Utils"
	"net/http"
	"slices"
)

func RoleGuardMiddleware(handler http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	if len(allowedRoles) == 0 {
		panic("At least one role must be specified for RoleGuardMiddleware")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get user role from context
		userRole := r.Context().Value(ContextUserRoleKey)

		if userRole == nil {
			_ = utils.WriteResponse(w, http.StatusBadRequest, utils.ReturnAppError(nil, 400, "Missing Role", "Missing Role"))
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok || userRoleStr == "" {
			_ = utils.WriteResponse(w, http.StatusBadRequest, utils.ReturnAppError(nil, 400, "Invalid Role Type", "Invalid Role Type"))
			return
		}

		// Check if user role is allowed
		if slices.Contains(allowedRoles, userRoleStr) {
			handler(w, r)
			return
		}

		_ = utils.WriteResponse(w, http.StatusForbidden, utils.ReturnAppError(nil, 403, "Forbidden To Access This Api", "Forbidden"))
	}
}
