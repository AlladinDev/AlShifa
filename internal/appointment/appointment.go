// Package appointment contains code for appoinment module
package appointment

import (
	"github.com/AlladinDev/AlShifa/internal/app"
	"github.com/AlladinDev/AlShifa/internal/appointment/controller"
	"github.com/AlladinDev/AlShifa/internal/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/internal/appointment/repository"
	"github.com/AlladinDev/AlShifa/internal/appointment/service"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterface "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/middleware"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

var Dependendies = []string{constants.NameAppStore, constants.NameClinicService, constants.NameDoctorClinicMappingModule, constants.NameDoctorService}

func InitAppointmentModule(di appInterface.IDependencyInjection) {
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get appstore from appointment module but failed")
	}
	repository := repository.NewRepository(appStore.DB)

	//get the clinic module from di
	clinicServiceAny, _ := di.GetService(constants.NameClinicService)
	clinicService, ok := clinicServiceAny.(interfaces.IClinicModule)
	if !ok {
		panic("Interface conversion  failed in appointment module when getting clinic module from di")
	}

	//get the doctor service
	doctorServiceAny, _ := di.GetService(constants.NameDoctorService)
	if doctorServiceAny == nil {
		panic("doctorModule needed in appointment module but not found in di")
	}
	doctorService, ok := doctorServiceAny.(interfaces.IDoctorModule)
	if !ok {
		panic("Interface conversion  failed in appointment module when getting doctor  module from di")
	}

	//get the clinicdoctor service from di
	clinicdoctorServiceAny, _ := di.GetService(constants.NameDoctorClinicMappingModule)
	clinicDoctorService, ok := clinicdoctorServiceAny.(interfaces.IClinicDoctor)
	if !ok {
		panic("tried to get clinic doctor mapping module from di in appointment module but failed")
	}

	//get the planservice from di
	planServiceAny, _ := di.GetService(constants.NamePlanService)
	planService, ok := planServiceAny.(interfaces.IPlanService)
	if !ok {
		panic("tried to get plan service from di in appointment module but failed")
	}
	service := service.NewService(repository, clinicService, appStore.DB.Client(), planService, doctorService, clinicDoctorService, utils.Encrypt, appStore.Broker, appStore.KVStore)

	//add this service to di
	di.AddService(constants.NameAppointmentService, service)

	controller := controller.NewController(service)

	//routes
	appStore.Route("/appointments", func(r chi.Router) {
		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleUser)).Post("/", controller.AddAppointment)
		r.With(middleware.JwtAuthmiddleware).Put("/", controller.UpdateAppointmentStatus)
	})
}
