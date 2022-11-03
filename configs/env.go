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

	if os.Getenv("MONGO_URI") == "" {
		log.Fatal("MONGO_URI is not set")
	}

	return os.Getenv("MONGOURI")
}

func EnvSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	if os.Getenv("SECRET_KEY") == "" {
		log.Fatal("SECRET is not set")
	}

	return os.Getenv("SECRET_KEY")
}
