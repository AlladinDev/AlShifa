package main

import (
	clinic "AlShifa/Clinic"
	internals "AlShifa/Internals"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	//call monogodb connect function
	mongoClient, mongoErr := internals.ConnectMongo(os.Getenv("MONGODB_URL"))
	if mongoErr != nil {
		log.Fatal("Failed to connect to mongodb", mongoErr)
	}

	appStore := internals.App{
		DB:     mongoClient.Database("AlShifa"),
		Server: http.NewServeMux(),
	}

	// appStore.Server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	json.NewEncoder(w).Encode("Server working")
	// })
	//initialise modules
	clinic.InitialiseClinicModule(&appStore)
	fmt.Print("Server Started")

	if err := http.ListenAndServe(port, appStore.Server); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}
