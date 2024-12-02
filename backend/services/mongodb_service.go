package services

import (
	"context"
	"log"
	"time"

	"github.com/PavelKilko/NetPulse/database"
	"github.com/PavelKilko/NetPulse/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMetricsForURL(urlID uint, duration time.Duration) ([]models.Metrics, error) {
	// Reference the "url_metrics" collection in the "netpulse" database
	collection := database.MongoClient.Database("netpulse").Collection("url_metrics")

	// Calculate the start time for the query
	startTime := time.Now().Add(-duration)

	// Build the query filter
	filter := bson.M{
		"url_id":    urlID,
		"timestamp": bson.M{"$gte": startTime},
	}

	// Execute the query
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to query metrics for URL ID %d: %v", urlID, err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Parse the results
	var metrics []models.Metrics
	if err := cursor.All(context.TODO(), &metrics); err != nil {
		log.Printf("Failed to decode metrics for URL ID %d: %v", urlID, err)
		return nil, err
	}

	return metrics, nil
}
