package main

import (
	"AlShifa/appointment"
	clinic "AlShifa/clinic"
	constants "AlShifa/constants"
	internals "AlShifa/internals"
	alShifaMiddlewares "AlShifa/middleware"
	"AlShifa/owner"
	users "AlShifa/users"
	utils "AlShifa/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	chiMiddlewares "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/redis/go-redis/v9"
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

	//create chi router
	chiRouter := chi.NewRouter()

	//--------------------apply middlewares----------------
	//timeout middewares
	chiRouter.Use(chiMiddlewares.Timeout(constants.RequestTimeout))
	chiRouter.Use(alShifaMiddlewares.Cors)
	chiRouter.Use(chiMiddlewares.Recoverer)

	//call monogodb connect function
	mongoClient, mongoErr := internals.ConnectMongo(os.Getenv("MONGODB_URL"))
	if mongoErr != nil {
		log.Fatal("Failed to connect to mongodb", mongoErr)
	}

	defer internals.Disconnect(mongoClient)

	//redis initialization
	//get redis url
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS url is missing in env file")
	}

	//get redis username
	redisUsername := os.Getenv("REDIS_USERNAME")
	if redisUsername == "" {
		log.Fatal("redisUsername is missing in env file")
	}

	//get redis passowrd
	redisUserPassword := os.Getenv("REDIS_PASSWORD")
	if redisUserPassword == "" {
		log.Fatal("REDIS userPassword is missing in env file")
	}

	//create redis client with all credientials
	redisInstance := *redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisUserPassword,
		Username: redisUsername,
	})

	// Ping the Redis server to verify the connection
	_, err := redisInstance.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("REDIS Connected successfully")
	///----------------------redis initialization end

	appStore := internals.NewApp().
		WithDB(mongoClient.Database("AlShifa")).
		WithDI().
		WithRedis(&redisInstance).
		WithServer(chiRouter)

	//add central middlewares

	//initialise modules
	clinic.InitialiseclinicModule(appStore)

	users.InitialiseUserModule(appStore)
	owner.InitOwner(appStore)

	appointment.InitAppointmentModule(appStore)

	fmt.Println("Server Started with dependency injection system initialized")

	//print services registered in di for debugging
	appStore.DI.PrintServicesInDI()

	if err := http.ListenAndServe(addr, chiRouter); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
