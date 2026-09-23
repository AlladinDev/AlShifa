package clinicdoctor

import (
	"github.com/AlladinDev/AlShifa/internal/app"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/controller"
	interfaces "github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/repository"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/service"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterface "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/middleware"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{constants.NameAppStore, constants.NameDoctorService, constants.NameClinicService}

func InitClinicDoctorMapping(di appInterface.IDependencyInjection) {

	//get the appstore
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get appstore in clinic doctor module but failed")
	}

	//here get the doctor module and clinic module dependency from dependecy injection
	doctorModuleAny, _ := di.GetService(constants.NameDoctorService)
	doctorModule, ok := doctorModuleAny.(interfaces.IDoctorModule)
	if !ok {
		panic("tried to get doctor module dependency in ClinicDoctor module but failed")
	}

	//get the clinic module dependency also
	clinicModuleAny, _ := di.GetService(constants.NameClinicService)
	clinicModule, ok := clinicModuleAny.(interfaces.IClinicModule)
	if !ok {
		panic("tried to get clinic module dependency in ClinicDoctor module but failed")
	}

	repository := repository.NewRepo(appStore.DB)
	clinicDoctorService := service.NewService(repository)
	onboardingService := service.NewOnboardingService(repository, appStore.EmailNotifier, utils.Encrypt, utils.Decrypt, utils.HashPasswordArgon2id, utils.GenerateOTP, clinicModule, doctorModule, utils.VerifyPasswordArgon2id)
	onboardingController := controller.NewOnboardingController(onboardingService)
	controller := controller.NewController(clinicDoctorService)
	//add this clinic doctor service to di
	di.AddService(constants.NameDoctorClinicMappingModule, clinicDoctorService)

	appStore.Route("/clinicdoctor", func(r chi.Router) {
		r.Get("/clinics", controller.FetchClinicsWithDoctors)
		r.Get("/doctors", controller.FetchDoctorWithClinics)
		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleclinicOwner)).Post("/", onboardingController.InitiateDoctorClinicOnboarding)
		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleclinicOwner)).Post("/onboarding/verify", onboardingController.VerifyDoctorClinicOnboarding)
	})

}
