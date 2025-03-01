package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"time"

	"github.com/redis/go-redis/v9"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ctx = context.Background()

var rdb *redis.Client

func init() {
	// Initialize Redis Client once
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

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
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
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

		log.Printf("📤 Task sent to queue %s", queueName)
	}

	responseWithJSON(w, http.StatusOK, "Task created successfully")

}

// ✅ Move Redis Subscription to Server Startup
func StartRedisListener() {
	log.Println("🚀 Redis listener started")
	pubsub := rdb.Subscribe(ctx, "task-status")
	defer pubsub.Close()

	log.Println("📡 Listening for task status updates from Redis...")
	for msg := range pubsub.Channel() {
		var data models.Task
		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			log.Printf("❌ Failed to decode Redis message: %s", err)
			continue
		}

		log.Printf("🔄 Task Update from Redis: %v", data)
		jsonData, _ := json.Marshal(data)
		broadcastToClients(jsonData)
	}
}
