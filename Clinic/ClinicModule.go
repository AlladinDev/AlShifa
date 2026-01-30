// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	controller "AlShifa/Clinic/Controller"
	repository "AlShifa/Clinic/Repository"
	service "AlShifa/Clinic/Service"
	validators "AlShifa/Clinic/Validators"
	"AlShifa/Clinic/models"
	emailnotifier "AlShifa/EmailNotifier"
	internals "AlShifa/Internals"
	middleware "AlShifa/Middleware"
	redispack "AlShifa/RedisPack"
	utils "AlShifa/Utils"
	"fmt"
	"log"
	"net/http"
)

func InitialiseClinicModule(app *internals.App) {
	//initialize redis module as we need it

	repository := repository.NewRepository(app.DB)

	//now arrange some dependencies for clinic service like cache ,otp generator
	addDoctorToClinicCache := redispack.NewRedisCache[string, models.AddDoctorOtpPayload](app.Redis, "AddDoctorToClinicCache")
	emailNotifier, smtpErr := emailnotifier.NewEmailService()
	if smtpErr != nil {
		log.Fatal("failed to get email service for clinic module", smtpErr)
	}

	service := service.NewClinicService(repository, app.DB.Client(), emailNotifier, utils.GenerateOTP, addDoctorToClinicCache)
	controller := controller.NewController(service, validators.ValidateAddDoctorToClinicDetails)
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/register"), controller.RegisterOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/register"), middleware.JwtAuthMiddleware(middleware.RoleGuardMiddleware(controller.RegisterClinic, utils.RoleClinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("GET", "/clinic/details"), middleware.JwtAuthMiddleware(controller.SearchClinic))
	app.Server.HandleFunc(utils.MakeURL("GET", "/owner/details"), middleware.JwtAuthMiddleware(middleware.RoleGuardMiddleware(controller.SearchOwner, utils.RoleAdmin, utils.RoleClinicOwner)))
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/register"), middleware.JwtAuthMiddleware(controller.RegisterDoctor))
	app.Server.HandleFunc(utils.MakeURL("GET", "/doctor/details"), middleware.JwtAuthMiddleware(controller.SearchDoctor))
	app.Server.HandleFunc(utils.MakeURL("POST", "/owner/login"), controller.LoginClinicOwner)
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/addDoctor"), middleware.JwtAuthMiddleware(controller.AddDoctorToClinic))
	app.Server.HandleFunc(utils.MakeURL("POST", "/clinic/addDoctor/verify"), middleware.JwtAuthMiddleware(controller.VerifyAddDoctorToClinicOtp))
	app.Server.HandleFunc(utils.MakeURL("POST", "/doctor/login"), controller.LoginDoctor)
	app.Server.HandleFunc(utils.MakeURL("GET", "/healthcheck"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hey buddy server is working for client module")
	})
}
