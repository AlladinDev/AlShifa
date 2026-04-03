package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlladinDev/AlShifa/appointment"
	"github.com/AlladinDev/AlShifa/auth"
	clinic "github.com/AlladinDev/AlShifa/clinic"
	constants "github.com/AlladinDev/AlShifa/constants"
	coordinator "github.com/AlladinDev/AlShifa/coodinator"
	internals "github.com/AlladinDev/AlShifa/internals"
	alShifaMiddlewares "github.com/AlladinDev/AlShifa/middleware"
	"github.com/AlladinDev/AlShifa/owner"
	users "github.com/AlladinDev/AlShifa/users"
	utils "github.com/AlladinDev/AlShifa/utils"

	chiMiddlewares "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	log.Printf("App running mode:  APP_ENV is = '%s'\n", os.Getenv("APP_ENV"))

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
	chiRouter.Use(chiMiddlewares.Logger)
	chiRouter.Use(chiMiddlewares.Recoverer)
	chiRouter.Use(alShifaMiddlewares.Cors)
	chiRouter.Use(chiMiddlewares.Timeout(constants.RequestTimeout))

	//all route handler to see if request is coming or  not

	//call monogodb connect function
	mongoClient, mongoErr := internals.ConnectMongo(os.Getenv("MONGODB_URL"))
	if mongoErr != nil {
		log.Fatal("Failed to connect to mongodb", mongoErr)
	}

	defer internals.Disconnect(mongoClient)

	//create redis client with all credientials
	redisInstance, redisConnErr := internals.ConnectToRedis()
	if redisConnErr != nil {
		log.Fatalf("Could not connect to Redis: %v", redisConnErr)
	}

	// version router this will be v1 router for api versioning
	v1 := chi.NewRouter()
	chiRouter.Mount("/v1", v1)

	appStore := internals.NewApp().
		WithDB(mongoClient.Database("AlShifa")).
		WithDI().
		WithRedis(redisInstance).
		WithServer(v1)

	//initialise modules
	clinic.InitialiseclinicModule(appStore)
	users.InitialiseUserModule(appStore)
	owner.InitOwner(appStore)
	auth.InitAuth(appStore)
	appointment.InitAppointmentModule(appStore)

	//at end initialise coordinator module
	coordinator.InitCoordinator(appStore)

	//print services registered in di for debugging
	appStore.DI.PrintServicesInDI()

	//now version apis

	utils.HealthCheck("alshifa", v1)

	//append prefix to this router such as v1

	if err := http.ListenAndServe(addr, chiRouter); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
