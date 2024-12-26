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
}

func ConsumeOrderNotifications(config *configs.Config, ch *amqp.Channel) {
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
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		var notification OrderNotification
		if err := json.Unmarshal(msg.Body, &notification); err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}

		// Process the notification
		log.Printf("Processing notification: %+v", notification)
		sendNotification(notification)
	}
}

func sendNotification(notification OrderNotification) {
	switch notification.NotificationType {
	case "email":
		log.Printf("Sending email for order %s to user %s", notification.OrderID, notification.UserID)
	case "sms":
		log.Printf("Sending SMS for order %s to user %s", notification.OrderID, notification.UserID)
	case "push":
		log.Printf("Sending push notification for order %s to user %s", notification.OrderID, notification.UserID)
	default:
		log.Printf("Unknown notification type for order %s", notification.OrderID)
	}
}
