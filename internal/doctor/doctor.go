// Package doctor provides doctor functionality for doctor module
package doctor

import (
	"github.com/AlladinDev/AlShifa/internal/app"
	doctorcontroller "github.com/AlladinDev/AlShifa/internal/doctor/controller"
	doctorinterfaces "github.com/AlladinDev/AlShifa/internal/doctor/interfaces"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/middleware"

	doctorrepository "github.com/AlladinDev/AlShifa/internal/doctor/repository"
	doctorservice "github.com/AlladinDev/AlShifa/internal/doctor/service"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{constants.NameAppStore, constants.NameClinicService}

func InitializeDoctorModule(di appInterfaces.IDependencyInjection) {
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get appstore from doctor module but unable to do it")
	}

	//here get the clinic dependency
	clinicServiceAny, _ := di.GetService(constants.NameClinicService)
	clinicModule, ok := clinicServiceAny.(doctorinterfaces.IClinic)
	if !ok {
		panic("tried to get clinic module dependency in doctor module but failed to get it")
	}

	repository := doctorrepository.NewRepository(appStore.DB)
	service := doctorservice.NewService(repository,
		utils.GenerateOTP,
		appStore.ImageUploader,
		appStore.EmailNotifier,
		utils.Encrypt,
		utils.Decrypt,
		utils.HashPasswordArgon2id,
		clinicModule, utils.GenerateJWT,
		utils.VerifyPasswordArgon2id)
	controller := doctorcontroller.NewController(service, utils.ParseUserID)

	//now add this service to di container
	di.AddService(constants.NameDoctorService, service)

	appStore.Route("/doctor", func(r chi.Router) {
		r.Post("/", controller.RegisterDoctor)
		r.Post("/verify", controller.VerifyDoctorRegistrationOtp)

		//login endpoints
		r.Post("/login", controller.InitiateLogin)
		r.Post("/login/verify", controller.VerifyLogin)

		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleDoctor)).Post("/onboard/allow", controller.AllowClinicToOnboard)
		r.Get("/", controller.GetDoctorProfile)
	})

}
