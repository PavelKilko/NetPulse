package main

import (
	"github.com/PavelKilko/NetPulse/monitor/repository"
	"github.com/PavelKilko/NetPulse/monitor/services"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize MongoDB connection
	mongoUrl := os.Getenv("MONGODB_URL")
	repository.InitMongo(mongoUrl)

	// Start RabbitMQ Consumer
	services.StartRabbitMQConsumer()
}
