package main

import (
	"github.com/PavelKilko/NetPulse/database"
	"github.com/PavelKilko/NetPulse/routes"
	"github.com/PavelKilko/NetPulse/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Fiber app
	app := fiber.New()

	// Connect to PostgreSQL
	database.ConnectDB()

	// Connect to Redis
	database.ConnectRedis()

	// Connect to MongoDB
	database.ConnectMongoDB()

	// Publish monitoring tasks to RabbitMQ
	services.PublishInitialMonitoringTasks()

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
