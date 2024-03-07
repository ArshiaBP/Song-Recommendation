package configs

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

var (
	Ch    *amqp.Channel
	Queue amqp.Queue
)

func ConnectToRabbitMQ() {
	conn, err := amqp.Dial(os.Getenv("RabbitMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	q, err := ch.QueueDeclare("requestIDs", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}
	Ch = ch
	Queue = q
}
