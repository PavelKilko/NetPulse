package services

import (
	"log"

	"github.com/PavelKilko/NetPulse/database"
	"github.com/PavelKilko/NetPulse/models"
)

func PublishInitialMonitoringTasks() {
	var urls []models.URL

	// Fetch URLs with monitoring enabled
	if err := database.DB.Where("monitoring = ?", true).Find(&urls).Error; err != nil {
		log.Printf("Failed to fetch URLs with monitoring enabled: %s", err)
		return
	}

	// Publish each URL to RabbitMQ
	for _, url := range urls {
		message := MonitoringMessage{
			URLID:  url.ID,
			Action: "enable",
			URL:    url.Address,
		}
		PublishToRabbitMQ(message)
	}
}
