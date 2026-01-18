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

	// Do not load .env file in Railway production
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT_NAME"); !exists {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file:", err)
		}
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

	if err := http.ListenAndServe(":"+port, appStore.Server); err != nil {
		fmt.Print("Failed to start server on error is", err)
	}

}

// package main

// import (
// 	validators "AlShifa/Clinic/Validators"
// 	"AlShifa/Clinic/models"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// // Helper to parse time strings without stopping execution
// func parseTime(s string) time.Time {
// 	t, err := time.Parse("03:04 PM", s)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	fmt.Println(t)
// 	return t
// }

// func main() {
// 	// Generate dummy data for a Clinic instance

// 	dummyClinic := models.Clinic{
// 		RegistrationDate: time.Date(2026, time.January, 17, 10, 0, 0, 0, time.UTC),
// 		Name:             "HealthFirst Wellness Center",
// 		Address:          "Soura Srinagar",
// 		SeasonTimings: []models.SeasonTimingDetails{
// 			{
// 				Name:  "Summer",
// 				Start: parseTime("06:00 AM"), // Matches "03:04 PM" layout
// 				End:   parseTime("08:00 PM"),
// 			},
// 			{
// 				Name:  "Winter",
// 				Start: parseTime("09:00 AM"),
// 				End:   parseTime("05:00 PM"),
// 			},
// 		},
// 		Mobile:  9876543000,
// 		Pincode: 190011,
// 		Doctors: nil,
// 		Wallet:  primitive.NilObjectID,
// 	}

// 	fmt.Print(validators.ValidateClinicDetails(&dummyClinic))
// }
