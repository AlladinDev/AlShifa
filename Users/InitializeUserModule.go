// Package users provides various layers for user module
package users

import (
	constants "github.com/AlladinDev/AlShifa/constants"
	internals "github.com/AlladinDev/AlShifa/internals"
	"github.com/AlladinDev/AlShifa/middleware"
	controller "github.com/AlladinDev/AlShifa/users/controller"
	repository "github.com/AlladinDev/AlShifa/users/repository"
	service "github.com/AlladinDev/AlShifa/users/service"

	"github.com/go-chi/chi/v5"
)

func InitialiseUserModule(app *internals.App) {
	repository := repository.ReturnNewRepository(app.DB)

	//here from dependency injection get the clinicService if not found panic
	clinicServiceAny, _ := app.DI.GetService(constants.NameClinicService)

	if clinicServiceAny == nil {
		panic("ClinicService required in usermodule not found in di")
	}

	service := service.ReturnNewService(repository)

	controller := controller.ReturnNewController(service)
	app.Route("/user", func(r chi.Router) {
		r.With(middleware.JwtAuthmiddleware).Post("/", controller.RegisterUser)
		r.With(middleware.JwtAuthmiddleware).Get("/details", controller.SearchUser)
	})

}
