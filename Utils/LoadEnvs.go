package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnvs() {

	environment := os.Getenv("APP_ENV")

	if environment == "github_actions" {
		return
	}

	_, b, _, _ := runtime.Caller(0)

	basePath := filepath.Dir(b)
	if err := godotenv.Load(filepath.Join(basePath, ".env")); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	fmt.Println("Envs loaded successfully")
}
