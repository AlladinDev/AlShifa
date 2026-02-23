// Package users provides various layers for user module
package users

import (
	constants "AlShifa/constants"
	internals "AlShifa/internals"
	middleware "AlShifa/middleware"
	controller "AlShifa/users/controller"
	interfaces "AlShifa/users/interfaces"
	repository "AlShifa/users/repository"
	service "AlShifa/users/service"
	utils "AlShifa/utils"
)

func InitialiseUserModule(app *internals.App) {
	repository := repository.ReturnNewRepository(app.DB)

	//here from dependency injection get the clinicService if not found panic
	clinicServiceAny, _ := app.DI.GetService(constants.NameClinicService)

	if clinicServiceAny == nil {
		panic("ClinicService required in usermodule not found in di")
	}

	clinicService, ok := clinicServiceAny.(interfaces.IUserClinicContract)
	if !ok {
		panic("interface clinicservice conversion failed in user module")
	}

	service := service.ReturnNewService(repository, clinicService)

	router := app.Server

	controller := controller.ReturnNewController(service)
	router.HandleFunc(utils.MakeURL("POST", "/user/register"), controller.RegisterUser)
	router.HandleFunc(utils.MakeURL("POST", "/user/login"), controller.LoginUser)
	router.HandleFunc(utils.MakeURL("GET", "/user/appointments"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.FetchAppointments, utils.RoleUser)))

	app.Server.HandleFunc(utils.MakeURL("GET", "/user/details"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.SearchUser, utils.RoleUser)))
}
