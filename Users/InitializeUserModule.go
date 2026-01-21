// Package users provides various layers for user module
package users

import (
	internals "AlShifa/Internals"
	middleware "AlShifa/Middleware"
	controller "AlShifa/Users/Controller"
	repository "AlShifa/Users/Repository"
	service "AlShifa/Users/Service"
	utils "AlShifa/Utils"
)

func InitialiseUserModule(app *internals.App) {
	repository := repository.ReturnNewRepository(app.DB)
	service := service.ReturnNewService(repository)
	controller := controller.ReturnNewController(service)
	app.Server.HandleFunc(utils.MakeURL("POST", "/user/register"), controller.RegisterUser)
	app.Server.HandleFunc(utils.MakeURL("POST", "/user/login"), controller.LoginUser)
	app.Server.HandleFunc(utils.MakeURL("GET", "/user/details"), middleware.JwtAuthMiddleware(controller.SearchUser))
}
