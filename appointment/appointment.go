// Package appointment contains code for appoinment module
package appointment

import (
	"AlShifa/appointment/controller"
	"AlShifa/appointment/interfaces"
	"AlShifa/appointment/repository"
	"AlShifa/appointment/service"
	"AlShifa/constants"
	"AlShifa/internals"
	"AlShifa/middleware"

	"github.com/go-chi/chi/v5"
)

func InitAppointmentModule(app *internals.App) {
	repository := repository.NewRepository(app.DB)

	//get the clinic module from di
	clinicModuleAny, _ := app.DI.GetService(constants.NameClinicService)
	clinicModule, ok := clinicModuleAny.(interfaces.IClinicModule)
	if !ok {
		panic("Interface conversion  failed in appointment module when getting clinic module from di")
	}

	service := service.NewService(repository, clinicModule)
	controller := controller.NewController(service)

	app.Route("/appointment", func(r chi.Router) {
		r.With(middleware.JwtAuthmiddleware).Post("/", controller.AddAppointment)
		r.With(middleware.JwtAuthmiddleware).Put("/", controller.UpdateAppointmentStatus)
		r.With(middleware.JwtAuthmiddleware).Get("/", controller.FetchAppointments)
	})
}
