package utils

import (
	"net/http"
)

func MethodCheckMiddleware(methodAllowed string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != methodAllowed {
			_ = InvalidMethodResponse(methodAllowed, w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
