// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	controller "AlShifa/Clinic/Controller"
	repository "AlShifa/Clinic/Repository"
	service "AlShifa/Clinic/Service"
	internals "AlShifa/Internals"
	middleware "AlShifa/Middleware"
	utils "AlShifa/Utils"
	"fmt"
	"net/http"
)

func InitialiseClinicModule(app *internals.App) {
	repository := repository.NewRepository(app.DB)
	service := service.NewClinicService(repository)
	controller := controller.NewController(service)
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/register"), controller.RegisterOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/register"), middleware.JwtAuthMiddleware(middleware.RoleGuardMiddleware(controller.RegisterClinic, utils.RoleClinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/details"), middleware.JwtAuthMiddleware(controller.SearchClinic))
	app.Server.HandleFunc(utils.MakeURL("GET", "/owner/details"), middleware.JwtAuthMiddleware(middleware.RoleGuardMiddleware(controller.SearchOwner, utils.RoleAdmin, utils.RoleClinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/register"), middleware.JwtAuthMiddleware(controller.RegisterDoctor))
	app.Server.HandleFunc(utils.MakeURL("GET", "/doctor/details"), middleware.JwtAuthMiddleware(controller.SearchDoctor))
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/login"), controller.LoginClinicOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/login"), controller.LoginDoctor)
	app.Server.HandleFunc(utils.MakeURL("GET", "/healthcheck"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hey buddy server is working for client module")
	})
}
