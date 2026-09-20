package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlladinDev/AlShifa/internal/clinic"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor"
	plan "github.com/AlladinDev/AlShifa/internal/clinicplan"
	"github.com/AlladinDev/AlShifa/internal/database"
	"github.com/AlladinDev/AlShifa/internal/dependencyinjection"
	doctor "github.com/AlladinDev/AlShifa/internal/doctor"
	"github.com/AlladinDev/AlShifa/internal/owner"
	user "github.com/AlladinDev/AlShifa/internal/users"

	app "github.com/AlladinDev/AlShifa/internal/app"
	appointmentModule "github.com/AlladinDev/AlShifa/internal/appointment"
	"github.com/AlladinDev/AlShifa/internal/shared/cloudinary"
	constants "github.com/AlladinDev/AlShifa/internal/shared/constants"
	alShifaMiddlewares "github.com/AlladinDev/AlShifa/internal/shared/middleware"
	natsbroker "github.com/AlladinDev/AlShifa/internal/shared/nats"
	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"

	chiMiddlewares "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	//here use recover function inside defer
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panicked but recovered inside main function error is ", err)
		}

	}()

	utils.LoadEnvs(".env")

	//get the port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	addr := ":" + port

	//create chi router parent router
	chiRouter := chi.NewRouter()

	//--------------------apply middlewares----------------
	//timeout middewares
	//chiRouter.Use(chiMiddlewares.Recoverer)
	chiRouter.Use(chiMiddlewares.Timeout(constants.RequestTimeout))
	//chiRouter.Use(chiMiddlewares.Logger)
	chiRouter.Use(alShifaMiddlewares.Cors)

	//call monogodb connect function
	mongoClient, mongoErr := database.ConnectMongo(os.Getenv("MONGODB_URL"))
	if mongoErr != nil {
		log.Fatal("Failed to connect to mongodb ", mongoErr)
	}

	defer database.Disconnect(mongoClient)

	//here get the nats broker
	broker, err := natsbroker.New()
	if err != nil {
		log.Fatalf("tried to get setup nats in main.go but failed err is : %v", err)
	}

	// version router this will be v1 router for api versioning
	v1 := chi.NewRouter()
	chiRouter.Mount("/v1", v1)

	//call initialization function in cloudinary so that modules can use it to upload files
	cloudinaryInstance := cloudinary.NewCloudinaryInstance()
	di := dependencyinjection.NewDI()

	appStore := app.NewApp().
		WithDB(mongoClient.Database("AlShifa")).
		WithServer(v1).
		WithEmailNotifier().
		WithBroker(broker).
		WithImageUploader(cloudinaryInstance)

	//add the appstore to di because modules need it
	di.AddService(constants.NameAppStore, appStore)

	di.RegisterModule(constants.NameClinicModule, clinic.InitialiseclinicModule, clinic.Dependencies)
	di.RegisterModule(constants.NamePlanModule, plan.InitializePlanModule, plan.Dependencies)
	di.RegisterModule(constants.NameDoctorModule, doctor.InitializeDoctorModule, doctor.Dependencies)
	di.RegisterModule(constants.NameDoctorClinicMappingModule, clinicdoctor.InitClinicDoctorMapping, clinicdoctor.Dependencies)
	di.RegisterModule(constants.NameUserModule, user.InitialiseUserModule, user.Dependencies)
	di.RegisterModule(constants.NameAppointmentModule, appointmentModule.InitAppointmentModule, appointmentModule.Dependendies)
	di.RegisterModule(constants.NameOwnerModule, owner.InitOwner, owner.Dependencies)
	di.Instantiatemodules()

	if err := http.ListenAndServe(addr, chiRouter); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
