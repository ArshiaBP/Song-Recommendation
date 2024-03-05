package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"request-registeration-service/configs"
)

func main() {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}
	configs.ConnectToDatabase()
	log.Println("connected to database")
}
