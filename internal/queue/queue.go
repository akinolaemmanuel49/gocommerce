package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupQueue(ch *amqp.Channel, logger, errorLogger *log.Logger) error {
	_, err := ch.QueueDeclare(
		"orderNotifications", // Queue name
		true,                 // Durable
		false,                // Auto-deleted
		false,                // Exclusive
		false,                // No-wait
		nil,                  // Arguments
	)
	if err != nil {
		errorLogger.Printf("Failed to declare queue: %v", err)
		return err
	}

	logger.Println("Queue initialized: orderNotifications")
	return nil
}
