package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return os.Getenv("MONGOURI")
}

func EnvSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return os.Getenv("SECRET_KEY")
}
