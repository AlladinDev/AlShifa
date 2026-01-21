// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	controller "AlShifa/Clinic/Controller"
	repository "AlShifa/Clinic/Repository"
	service "AlShifa/Clinic/Service"
	internals "AlShifa/Internals"
	utils "AlShifa/Utils"
	"fmt"
	"net/http"
)

func InitialiseClinicModule(app *internals.App) {
	repository := repository.NewRepository(app.DB)
	service := service.NewClinicService(repository)
	controller := controller.NewController(service)
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/register"), controller.RegisterOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/register"), controller.RegisterClinic)
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/details"), controller.SearchClinic)
	app.Server.HandleFunc(utils.MakeURL("GET", "/owner/details"), controller.SearchOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/register"), controller.RegisterDoctor)
	app.Server.HandleFunc(utils.MakeURL("GET", "/doctor/details"), controller.SearchDoctor)
	app.Server.HandleFunc(utils.MakeURL("GET", "/healthcheck"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hey buddy server is working for client module")
	})
}
