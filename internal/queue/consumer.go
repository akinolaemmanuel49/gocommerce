package queue

import (
	"encoding/json"
	"log"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderNotification struct {
	OrderID          string `json:"orderID"`
	UserID           string `json:"userID"`
	NotificationType string `json:"type"`
	Message          string `json:"message"`
}

func ConsumeOrderNotifications(config *configs.Config, ch *amqp.Channel, logger, errorLogger *log.Logger) {
	msgs, err := ch.Consume(
		config.OrderQueueName, // Queue name
		"",                    // Consumer tag
		true,                  // Auto-acknowledge
		false,                 // Exclusive
		false,                 // No-local
		false,                 // No-wait
		nil,                   // Arguments
	)
	if err != nil {
		errorLogger.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		var notification OrderNotification
		if err := json.Unmarshal(msg.Body, &notification); err != nil {
			logger.Printf("Failed to decode message: %v", err)
			continue
		}

		// Process the notification
		logger.Printf("Processing notification: %+v", notification)
	}
}
