// Package appointment contains code for appoinment module
package appointment

import (
	"github.com/AlladinDev/AlShifa/appointment/controller"
	"github.com/AlladinDev/AlShifa/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/appointment/repository"
	"github.com/AlladinDev/AlShifa/appointment/service"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/internals"
	"github.com/AlladinDev/AlShifa/middleware"

	"github.com/go-chi/chi/v5"
)

func InitAppointmentModule(app *internals.App) {
	repository := repository.NewRepository(app.DB)

	//get the clinic module from di
	clinicModuleAny, _ := app.DI.GetService(constants.NameClinicService)
	clinicModule := clinicModuleAny.(interfaces.IClinicModule)
	// if !ok {
	// 	panic("Interface conversion  failed in appointment module when getting clinic module from di")
	// }

	service := service.NewService(repository, clinicModule)

	//add this service to di
	app.DI.AddService(constants.NameAppointmentService, service)

	controller := controller.NewController(service)

	app.Route("/appointments", func(r chi.Router) {
		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleUser)).Post("/", controller.AddAppointment)
		r.With(middleware.JwtAuthmiddleware).Put("/", controller.UpdateAppointmentStatus)
		r.With(middleware.JwtAuthmiddleware).Get("/", controller.FetchAppointments)
	})
}
