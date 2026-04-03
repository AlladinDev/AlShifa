// Package coordinator handles coodination between different modules
package coordinator

import (
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/coodinator/controller"
	"github.com/AlladinDev/AlShifa/coodinator/interfaces"
	"github.com/AlladinDev/AlShifa/coodinator/service"
	"github.com/AlladinDev/AlShifa/internals"
	"github.com/AlladinDev/AlShifa/middleware"
)

func InitCoordinator(app *internals.App) {
	clinicServiceAny, _ := app.DI.GetService(constants.NameClinicService)

	clinicService, ok := clinicServiceAny.(interfaces.IClinicService)
	if !ok {
		panic("failed to convert clinic service from di in coordinator it failes to implement interface")
	}

	appointmentServiceAny, _ := app.DI.GetService(constants.NameAppointmentService)
	appointmentService, ok := appointmentServiceAny.(interfaces.IAppointmentService)
	if !ok {
		panic("failed to convert clinic service from di in coordinator it failes to implement interface")
	}

	service := service.NewAppointmentCoordinator(clinicService, app.DB.Client(), appointmentService)
	controller := controller.NewController(service)

	app.With(middleware.JwtAuthmiddleware).Post("/appointment", controller.BookAppointment)
}
