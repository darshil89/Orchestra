package main

import (
	"encoding/json"
	"log"
	"time"

	"api/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"APITask", // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag (empty for auto-generated)
		false,  // auto-ack (set to false for manual acknowledgment)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to register a consumer")

	// Create a channel to keep the program running
	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			log.Printf("Received Task")
			var d models.Task
			err := json.Unmarshal(msg.Body, &d)
			if err != nil {
				log.Printf("Failed to unmarshal task: %s", err)
			}
			time.Sleep(5 * time.Second) // Simulate processing time
			log.Printf("%s", "API Task is done "+d.Title)

			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
