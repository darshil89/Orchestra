package main

import (
	"context"
	"db/models"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

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
		"DBTask", // queue name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
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

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create a channel to keep the program running
	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			var d models.Task
			err := json.Unmarshal(msg.Body, &d)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("🚀 Processing Task: %s (ID: %d)", d.Title, d.ID)

			// Publish "processing" status to Redis
			statusUpdate := models.Task{
				ID:          d.ID,
				Title:       d.Title,
				Status:      models.InProgress,
				Description: d.Description,
				Function:    d.Function,
			}
			publishToRedis(rdb, "task-status", statusUpdate)

			// Simulate task processing
			time.Sleep(2 * time.Second)

			// Publish "completed" status to Redis
			statusUpdate.Status = models.Done
			publishToRedis(rdb, "task-status", statusUpdate)

			log.Printf("✅ Task %s completed", d.Title)

			// Acknowledge message after processing
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Function to properly format and publish messages to Redis
func publishToRedis(client *redis.Client, channel string, message models.Task) {
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Error encoding JSON: %s", err)
		return
	}
	client.Publish(ctx, channel, jsonMsg)
}
