package main

import (
	clinic "AlShifa/Clinic"
	internals "AlShifa/Internals"
	users "AlShifa/Users"
	utils "AlShifa/Utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/redis/go-redis/v9"
)

func printMemUsage() (string, error) {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)

	byteData, err := json.Marshal(&m)

	if err != nil {
		return "", err
	}

	return string(byteData), nil

}

func main() {

	printMemUsage()

	log.Printf("App running mode:  APP_ENV is = '%s'\n", os.Getenv("APP_ENV"))

	utils.LoadEnvs(".env")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	addr := ":" + port

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

	appStore := internals.App{
		DB:     mongoClient.Database("AlShifa"),
		Server: http.NewServeMux(),
		Redis:  &redisInstance,
	}

	//initialise modules
	clinic.InitialiseClinicModule(&appStore)
	users.InitialiseUserModule(&appStore)

	fmt.Print("Server Started")

	appStore.Server.HandleFunc("GET /cpu/usage", func(w http.ResponseWriter, r *http.Request) {
		cpuStatistics, err := printMemUsage()
		if err != nil {
			_ = utils.WriteResponse(w, http.StatusInternalServerError, "Failed to get cpu usage")
			return
		}

		_ = utils.WriteResponse(w, http.StatusOK, cpuStatistics)
	})

	if err := http.ListenAndServe(addr, appStore.Server); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
