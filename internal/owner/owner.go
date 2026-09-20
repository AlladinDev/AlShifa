// Package owner provides functionality for managing owners.
package owner

import (
	"github.com/AlladinDev/AlShifa/internal/app"
	"github.com/AlladinDev/AlShifa/internal/owner/controller"
	"github.com/AlladinDev/AlShifa/internal/owner/repository"
	"github.com/AlladinDev/AlShifa/internal/owner/service"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/AlladinDev/AlShifa/internal/shared/middleware"
	"github.com/AlladinDev/AlShifa/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

var Dependencies = []string{constants.NameAppStore}

func InitOwner(di appInterfaces.IDependencyInjection) {
	appStoreAny, _ := di.GetService(constants.NameAppStore)
	appStore, ok := appStoreAny.(*app.App)
	if !ok {
		panic("tried to get appstore in owner module but failed")
	}
	Repository := repository.NewRepository(appStore.DB)
	service := service.NewService(Repository,
		utils.GenerateJWT,
		utils.VerifyPasswordArgon2id,
		utils.Decrypt,
		utils.Encrypt,
		appStore.KVStore,
		utils.HashPasswordArgon2id,
		utils.GenerateOTP,
		appStore.ImageUploader,
		appStore.EmailNotifier,
		utils.VerifyPasswordArgon2id)

	controller := controller.NewController(service)

	appStore.Route("/owner", func(r chi.Router) {
		r.Get("/", controller.GetOwnerByID)
		r.Post("/otp", controller.RegisterOwner)
		r.Post("/otp/verify", controller.VerifyOwnerRegistrationOTP)

		r.Post("/login", controller.InitiateLogin)
		r.Post("/login/verify", controller.VerifyLogin)

		r.With(middleware.JwtAuthmiddleware, middleware.RoleGuardmiddleware(constants.RoleclinicOwner)).Get("/", controller.GetOwnerByID)
	})

}
