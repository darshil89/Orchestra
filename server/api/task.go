package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	models "server/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var data models.Tasks
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Invalid request payload")
		return
	}

	// connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queues := map[string]string{
		"DB Push":        "DBTask",
		"Dockerhub Push": "DockerhubTask",
		"API Call":       "APITask",
	}

	// Declare the response queue
	responseQueue, err := ch.QueueDeclare(
		"ResponseQueue", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare response queue")
	for _, task := range data {
		queueName, exists := queues[task.Title]
		if !exists {
			log.Printf("Unknown task type: %s", task.Title)
			continue
		}

		_, err := ch.QueueDeclare(
			queueName, // Queue name
			false,     // Durable
			false,     // Delete when unused
			false,     // Exclusive
			false,     // No-wait
			nil,       // Arguments
		)
		failOnError(err, "Failed to declare a queue")

		// Publish task to the correct queue
		taskBody, _ := json.Marshal(task)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = ch.PublishWithContext(ctx,
			"",        // Exchange
			queueName, // Routing key (queue name)
			false,     // Mandatory
			false,     // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        taskBody,
				ReplyTo:     "ResponseQueue",
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent task to %s ", queueName)
	}

	// Listen for responses
	msgs, err := ch.Consume(
		responseQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to register consumer for response")

	go func() {
		for d := range msgs {
			log.Printf("Received response: %s", d.Body)
		}
	}()

	log.Println(" [*] Waiting for responses. To exit, press CTRL+C")
	select {} // Keep running

	// responseWithJSON(w, http.StatusOK, "Task created successfully")

}
