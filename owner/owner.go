// Package owner provides functionality for managing owners.
package owner

import (
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/internals"
	"github.com/AlladinDev/AlShifa/middleware"
	"github.com/AlladinDev/AlShifa/owner/controller"
	"github.com/AlladinDev/AlShifa/owner/repository"
	"github.com/AlladinDev/AlShifa/owner/service"
	"github.com/go-chi/chi/v5"
)

func InitOwner(app *internals.App) {
	Repository := repository.NewRepository(app.DB)
	service := service.NewService(Repository)
	controller := controller.NewController(service)

	app.Route("/owner", func(r chi.Router) {
		r.Post("/", controller.RegisterOwner)
		r.Post("/login", controller.LoginOwner)
		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleclinicOwner)).Get("/", controller.GetOwnerByID)
	})

}
