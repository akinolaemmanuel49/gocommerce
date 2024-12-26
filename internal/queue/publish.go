package queue

import (
	"fmt"
	"log"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishOrderNotification(config *configs.Config, ch *amqp.Channel, orderID, userID, notificationType string) error {
	message := fmt.Sprintf(`{"orderID": "%s", "userID": "%s", "type": "%s"}`, orderID, userID, notificationType)

	err := ch.Publish(
		"",                    // Default exchange
		config.OrderQueueName, // Routing key (queue name)
		false,                 // Mandatory
		false,                 // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Published notification: %s", message)
	return nil
}
