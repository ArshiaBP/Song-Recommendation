package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"request-registeration-service/configs"
	"request-registeration-service/controllers"
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
	server := echo.New()
	server.POST("/send-song/:email", controllers.SaveRequestHandler)
	log.Fatal(server.Start("localhost:8080"))
}
