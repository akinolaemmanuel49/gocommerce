package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	l "github.com/akinolaemmanuel49/gocommerce/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderMessage represents the structure of the message to be sent
type OrderMessage struct {
	OrderID          string    `json:"orderId"`
	UserID           string    `json:"userId"`
	EventType        string    `json:"eventType"`                                                  // e.g., "OrderCreated", "OrderConfirmed"
	Status           string    `json:"status" validate:"required,oneof=pending shipped delivered"` // "pending" "shipped" "delivered"
	Message          string    `json:"message"`
	NotificationTime time.Time `json:"notificationTime"`
}

// Publisher handles publishing messages to the queue
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// NewPublisher initializes a new Publisher
func NewPublisher(config *configs.Config) (*Publisher, error) {
	// Setup logger
	logger, err := l.SetupLogger("service.log", "INFO")
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Fatalf("%s", fmt.Sprintf("%-7s: Error setting up error logger: %v", "ERROR", err))
	}
	errorLogger, err := l.SetupLogger("service.log", "ERROR")
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Fatalf("%s", fmt.Sprintf("%-7s: Error setting up error logger: %v", "ERROR", err))
	}

	conn, ch, err := ConnectRabbitMQ(config, logger, errorLogger)

	if err != nil {
		return nil, err
	}

	fmt.Println("CHANNEL CHECK")

	return &Publisher{
		conn:    conn,
		channel: ch,
		queue:   config.OrderQueueName,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, message OrderMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish the message
	return p.channel.Publish(
		"",      // exchange (default)
		p.queue, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close closes the connection and channel
func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
