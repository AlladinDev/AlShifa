// Package auth handles auth related tasks registration login jwt sign verification etc but just related to bare auth
package auth

import (
	"github.com/AlladinDev/AlShifa/auth/controller"
	"github.com/AlladinDev/AlShifa/auth/repository"
	"github.com/AlladinDev/AlShifa/auth/service"
	"github.com/AlladinDev/AlShifa/internals"

	"github.com/go-chi/chi/v5"
)

func InitAuth(app *internals.App) {
	repository := repository.NewRepository(app.DB)
	service := service.NewService(repository)
	controller := controller.NewController(service)

	app.Route("/auth", func(r chi.Router) {
		r.Post("/register", controller.Register)
		r.Post("/login", controller.Login)
	})
}
