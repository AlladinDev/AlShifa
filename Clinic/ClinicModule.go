// Package clinic provides functionalities related to clinic management,
// including owner registration and health checks.
package clinic

import (
	controller "github.com/AlladinDev/AlShifa/clinic/controller"
	"github.com/AlladinDev/AlShifa/clinic/models"
	repository "github.com/AlladinDev/AlShifa/clinic/repository"
	service "github.com/AlladinDev/AlShifa/clinic/service"
	validators "github.com/AlladinDev/AlShifa/clinic/validators"
	constants "github.com/AlladinDev/AlShifa/constants"
	emailnotifier "github.com/AlladinDev/AlShifa/emailnotifier"
	internals "github.com/AlladinDev/AlShifa/internals"
	middleware "github.com/AlladinDev/AlShifa/middleware"
	redispack "github.com/AlladinDev/AlShifa/redispack"

	"fmt"
	"log"

	utils "github.com/AlladinDev/AlShifa/utils"

	"github.com/go-chi/chi/v5"
)

func InitialiseclinicModule(app *internals.App) {
	//initialize redis module as we need it

	repository := repository.NewRepository(app.DB)

	//now arrange some dependencies for clinic service like cache ,otp generator
	otpCache := redispack.NewRedisCache[string, models.AddDoctorOtpPayload](app.Redis, "clinic")
	emailNotifier, smtpErr := emailnotifier.NewEmailService()
	if smtpErr != nil {
		log.Fatal("failed to get email service for clinic module", smtpErr)
	}

	//get service
	service := service.NewclinicService(repository, app.DB.Client(), emailNotifier, utils.GenerateOTP, otpCache)

	//add this service to di container so that other modules can use it
	app.DI.AddService(constants.NameClinicService, service)

	//get the controller
	controller := controller.NewController(service, validators.ValidateAddDoctorToclinicDetails)

	//-------------------Routes start here----------------------------//

	//doctor routes
	//healthcheck for doctor routes
	utils.HealthCheck("doctor", app.Mux)

	//doctor routes without middlewares
	app.Get("/doctors", controller.SearchDoctor)
	app.Post("/doctor", controller.RegisterDoctor)
	app.Post("/doctor/login", controller.LoginDoctor)
	app.With(middleware.JwtAuthmiddleware).Get("/doctor/clinics", controller.FetchDoctorWithClinics)
	app.With(middleware.JwtAuthmiddleware).Get("/clinic/doctors", controller.FetchClinicWithDoctors)
	// Clinic routes
	app.Route("/clinic", func(clinic chi.Router) {
		app.Post("/", controller.Registerclinic)

		//healthcheck for clinic routes
		utils.HealthCheck("clinic", app.Mux)

		clinic.Get("/", controller.Searchclinic)
		clinic.Group(func(c chi.Router) {

			c.Use(middleware.JwtAuthmiddleware)

			// c.Use(middleware.RoleGuardmiddleware(controller.Registerclinic, utils.RoleclinicOwner))

			//c.With(middleware.RoleGuardmiddleware(constants.RoleclinicOwner)).Post("/register", controller.Registerclinic)
			c.Post("/addDoctor", controller.AddDoctorToclinic)
			c.Post("/addDoctor/verify", controller.VerifyAddDoctorToclinicOtp)
		})

	})

	fmt.Println("Clinic module initialized")
}
