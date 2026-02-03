package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvs(path string) {

	environment := os.Getenv("APP_ENV")

	if environment == "github_actions" {
		return
	}

	if err := godotenv.Load(path); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	fmt.Println("Envs loaded successfully")
}
