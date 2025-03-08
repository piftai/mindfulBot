package utils

import (
	"log"

	"github.com/joho/godotenv"
)

// Env init environment variables
func Env() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Not found .env")
	}
}
