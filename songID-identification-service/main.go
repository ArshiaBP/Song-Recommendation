package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"songID-identification-service/configs"
	"songID-identification-service/controllers"
)

func main() {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}
	configs.ConnectToDatabase()
	configs.ConnectToRabbitMQ()
	controllers.Consume()
}
