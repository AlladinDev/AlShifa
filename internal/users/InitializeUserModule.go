// Package users provides various layers for user module
package users

import (
	app "github.com/AlladinDev/AlShifa/internal/app"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/emailnotifier"
	appInterface "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	controller "github.com/AlladinDev/AlShifa/internal/users/controller"
	repository "github.com/AlladinDev/AlShifa/internal/users/repository"
	service "github.com/AlladinDev/AlShifa/internal/users/service"

	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{constants.NameAppStore}

func InitialiseUserModule(di appInterface.IDependencyInjection) {
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get appstore from di in user module but failed")
	}
	repository := repository.NewRepository(appStore.DB)

	service := service.ReturnNewService(utils.GenerateOTP,
		utils.HashPasswordArgon2id,
		utils.Decrypt, utils.Encrypt,
		emailnotifier.NewEmailNotifier(),
		utils.VerifyPasswordArgon2id,
		utils.GenerateJWT, repository,
	)

	///add this service to di
	di.AddService(constants.NameUserService, service)

	controller := controller.ReturnNewController(service)
	appStore.Route("/user", func(r chi.Router) {
		r.Post("/", controller.InitiateLogin)
		r.Post("/verify", controller.VerifyLogin)
	})

}
