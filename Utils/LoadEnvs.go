package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	_, b, _, _ := runtime.Caller(0)

	basePath := filepath.Dir(b)
	if err := godotenv.Load(filepath.Join(basePath, ".env")); err != nil {
		log.Fatal("error loading .env file:", err)
	}

	fmt.Println("Envs loaded successfully")
}
