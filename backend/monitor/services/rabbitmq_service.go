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

func StartRabbitMQConsumer() {
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var msg MonitoringMessage
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}

			if msg.Action == "enable" {
				StartMonitoring(msg.URLID, msg.URL)
			} else if msg.Action == "disable" {
				StopMonitoring(msg.URLID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}
