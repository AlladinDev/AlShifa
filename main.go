package main

import (
	clinic "AlShifa/Clinic"
	internals "AlShifa/Internals"
	users "AlShifa/Users"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
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

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file:", err)
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

	appStore := internals.App{
		DB:     mongoClient.Database("AlShifa"),
		Server: http.NewServeMux(),
	}

	//initialise modules
	clinic.InitialiseClinicModule(&appStore)
	users.InitialiseUserModule(&appStore)

	fmt.Print("Server Started")

	if err := http.ListenAndServe(addr, appStore.Server); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
