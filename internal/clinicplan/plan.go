// Package plan provides clinic plan facility
package plan

import (
	"fmt"
	"reflect"

	"github.com/AlladinDev/AlShifa/internal/app"
	plancontroller "github.com/AlladinDev/AlShifa/internal/clinicplan/controller"
	"github.com/AlladinDev/AlShifa/internal/clinicplan/repository"
	planservice "github.com/AlladinDev/AlShifa/internal/clinicplan/service"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{
	constants.NameAppStore,
}

func InitializePlanModule(di appInterfaces.IDependencyInjection) {
	//get the appStore from di
	appStoreAny, _ := di.GetService(constants.NameAppStore)

	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		fmt.Print(reflect.TypeOf(appStoreAny))
		panic("tried to get appStore from di in planmodule but got error")
	}
	repository := repository.NewRepository(appStore.DB)
	service := planservice.NewService(repository)
	controller := plancontroller.NewController(service)
	di.AddService(constants.NamePlanService, service)

	//now make apis
	appStore.Route("/plan", func(r chi.Router) {
		r.Post("/", controller.AddNewPlan)
		r.Put("/", controller.UpdatePlan)
		r.Get("/", controller.FetchPlans)
	})
}
