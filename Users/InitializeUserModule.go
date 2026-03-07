// Package users provides various layers for user module
package users

import (
	constants "AlShifa/constants"
	internals "AlShifa/internals"
	"AlShifa/middleware"
	controller "AlShifa/users/controller"
	repository "AlShifa/users/repository"
	service "AlShifa/users/service"

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
