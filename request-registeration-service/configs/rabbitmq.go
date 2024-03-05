package configs

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
	"log"
	"os"
	"sync"
	"time"
)

var (
	Ch     *amqp.Channel
	Ctx    context.Context
	Queue  amqp.Queue
	onceMQ sync.Once
)

func ConnectToRabbitMQ() {
	onceMQ.Do(func() {
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		Ch = ch
		Queue = q
		Ctx = ctx
	})
}
