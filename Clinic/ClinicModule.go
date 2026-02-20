// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	emailnotifier "AlShifa/EmailNotifier"
	redispack "AlShifa/RedisPack"
	controller "AlShifa/clinic/Controller"
	repository "AlShifa/clinic/Repository"
	service "AlShifa/clinic/Service"
	validators "AlShifa/clinic/Validators"
	"AlShifa/clinic/models"
	constants "AlShifa/constants"
	internals "AlShifa/internals"
	middleware "AlShifa/middleware"
	utils "AlShifa/utils"
	"fmt"
	"log"
	"net/http"
)

func InitialiseclinicModule(app *internals.App) {
	//initialize redis module as we need it

	repository := repository.NewRepository(app.DB)

	//now arrange some dependencies for clinic service like cache ,otp generator
	addDoctorToclinicCache := redispack.NewRedisCache[string, models.AddDoctorOtpPayload](app.Redis, "AddDoctorToclinicCache")
	emailNotifier, smtpErr := emailnotifier.NewEmailService()
	if smtpErr != nil {
		log.Fatal("failed to get email service for clinic module", smtpErr)
	}

	service := service.NewclinicService(repository, app.DB.Client(), emailNotifier, utils.GenerateOTP, addDoctorToclinicCache)

	//add this service to di container so that other modules can use it
	app.DI.AddService(constants.NameClinicService, service)

	controller := controller.NewController(service, validators.ValidateAddDoctorToclinicDetails)
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/register"), controller.RegisterOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/register"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.Registerclinic, utils.RoleclinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/details"), controller.Searchclinic)
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/doctorAtclinics"), middleware.JwtAuthmiddleware(controller.FetchDoctorWithItsclinics))
	app.Server.HandleFunc(utils.MakeURL("GET", "/owner/details"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.SearchOwner, utils.RoleAdmin, utils.RoleclinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/register"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.RegisterDoctor, utils.RoleclinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("GET", "/doctor/details"), controller.SearchDoctor)
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/login"), controller.LoginclinicOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/addDoctor"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.AddDoctorToclinic, utils.RoleclinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/appointments/book"), middleware.JwtAuthmiddleware(controller.AddAppointment))
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/addDoctor/verify"), middleware.JwtAuthmiddleware(middleware.RoleGuardmiddleware(controller.VerifyAddDoctorToclinicOtp, utils.RoleclinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/login"), controller.LoginDoctor)
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/appointments/full"), middleware.JwtAuthmiddleware(controller.AppointmentSlotsBooked))
	app.Server.HandleFunc(utils.MakeURL("GET", "/healthcheck"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hey buddy server is working for client module")
	})

	fmt.Println("Clinic module initialized")
}
