package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var MongoClient *mongo.Client

// InitMongo Initialize MongoDB connection
func InitMongo(connectionString string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}
	MongoClient = client
}

// StoreMonitoringResult stores a monitoring result in MongoDB
func StoreMonitoringResult(urlID uint, responseTime int, statusCode int, timestamp time.Time) {
	collection := MongoClient.Database("netpulse").Collection("url_metrics")
	_, err := collection.InsertOne(context.TODO(), bson.M{
		"url_id":        urlID,
		"response_time": responseTime,
		"status_code":   statusCode,
		"timestamp":     timestamp,
	})
	if err != nil {
		log.Printf("Failed to insert monitoring result: %s", err)
	}
}
