// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	"fmt"

	"github.com/AlladinDev/AlShifa/internal/app"
	controller "github.com/AlladinDev/AlShifa/internal/clinic/controller"
	clinicinterfaces "github.com/AlladinDev/AlShifa/internal/clinic/interfaces"
	repository "github.com/AlladinDev/AlShifa/internal/clinic/repository"
	service "github.com/AlladinDev/AlShifa/internal/clinic/service"
	middlewares "github.com/AlladinDev/AlShifa/internal/shared/middleware"

	constants "github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"

	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{constants.NameAppStore, constants.NamePlanService}

func InitialiseclinicModule(di appInterfaces.IDependencyInjection) {
	//initialize redis module as we need it
	//now get the appStore from di
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get app store in clinic module but failed")
	}

	repository := repository.NewRepository(appStore.DB)

	//get the clinicplan module from di as service needs it
	clinicPlanAny, _ := di.GetService(constants.NamePlanService)
	if clinicPlanAny == nil {
		panic("tried to get clinic plan module from di but not found")
	}

	clinicPlan, ok := clinicPlanAny.(clinicinterfaces.IClinicPlan)
	if !ok {
		panic("tried to get clinic plan module from di but type assertion failed inside clinic module")
	}

	//get service
	service := service.NewclinicService(repository, clinicPlan)

	//get the controller
	controller := controller.NewController(service)

	//now add this clinic service to di container
	di.AddService(constants.NameClinicService, service)

	//-------------------Routes start here----------------------------//

	// Clinic routes
	appStore.Route("/clinic", func(clinic chi.Router) {
		clinic.With(middlewares.JwtAuthmiddleware, middlewares.RoleGuardmiddleware(constants.RoleclinicOwner)).Post("/", controller.Registerclinic)
		clinic.Get("/", controller.Searchclinic)
	})

	fmt.Println("Clinic module initialized")
}
