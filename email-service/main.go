package main

import (
	"email-service/configs"
	"email-service/services"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}
	configs.ConnectToDatabase()
	log.Println("database connected")
	log.Println("service started")
	services.Listen()
}
