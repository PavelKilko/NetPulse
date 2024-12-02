package services

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type MonitoringMessage struct {
	URLID  uint   `json:"url_id"`
	Action string `json:"action"`
	URL    string `json:"url"`
}

func PublishToRabbitMQ(message MonitoringMessage) {
	rabbitMQUrl := os.Getenv("RABBITMQ_URL")
	conn, err := amqp091.Dial(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"monitoring_queue", // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %s", err)
	} else {
		log.Printf("Published message to RabbitMQ: %s", body)
	}
}
