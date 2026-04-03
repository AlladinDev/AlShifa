package utils

import (
	"fmt"
	"net/http"

	structs "github.com/AlladinDev/AlShifa/structs"
	"github.com/go-chi/chi/v5"
)

func HealthCheck(routerName string, chiRouter *chi.Mux) {
	fullPath := fmt.Sprintf("/%s/health", routerName)

	chiRouter.Get(fullPath, func(w http.ResponseWriter, r *http.Request) {
		_ = WriteResponse(w, http.StatusOK, structs.IAppSuccess{
			Message:    fmt.Sprintf("%s working", routerName),
			Data:       nil,
			StatusCode: http.StatusOK,
		})
	})
}
