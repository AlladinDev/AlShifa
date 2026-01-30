package main

import (
	clinic "AlShifa/Clinic"
	internals "AlShifa/Internals"
	users "AlShifa/Users"
	utils "AlShifa/Utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/redis/go-redis/v9"
)

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func main() {

	printMemUsage()

	if os.Getenv("APP_ENV") != "production" {
		utils.LoadEnvs()
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	addr := "0.0.0.0:" + port

	//call monogodb connect function
	mongoClient, mongoErr := internals.ConnectMongo(os.Getenv("MONGODB_URL"))
	if mongoErr != nil {
		log.Fatal("Failed to connect to mongodb", mongoErr)
	}

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

	appStore := internals.App{
		DB:     mongoClient.Database("AlShifa"),
		Server: http.NewServeMux(),
		Redis:  &redisInstance,
	}

	//initialise modules
	clinic.InitialiseClinicModule(&appStore)
	users.InitialiseUserModule(&appStore)

	fmt.Print("Server Started")

	if err := http.ListenAndServe(addr, appStore.Server); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
