package queue

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectRabbitMQ implements logic for connecting to the RabbitMQ queue
func ConnectRabbitMQ(config *configs.Config, logger, errorLogger *log.Logger, isConsumer bool) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(config.RabbitMQURI)
	if err != nil {
		errorLogger.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		errorLogger.Fatalf("Failed to open a channel: %v", err)
		return conn, nil, err
	}

	err = setupQueue(config, ch, errorLogger)
	if err != nil {
		errorLogger.Fatalf("Failed to setup queue: %v", err)
		return nil, nil, err
	}
	if isConsumer {
		logger.Println("Queue initialized: orderNotifications")
	}
	return conn, ch, nil
}

// setupQueue declares a queue
func setupQueue(config *configs.Config, ch *amqp.Channel, errorLogger *log.Logger) error {
	_, err := ch.QueueDeclare(
		config.OrderQueueName, // Queue name
		true,                  // Durable
		false,                 // Auto-deleted
		false,                 // Exclusive
		false,                 // No-wait
		nil,                   // Arguments
	)
	if err != nil {
		errorLogger.Printf("Failed to declare queue: %v", err)
		return err
	}

	return nil
}
