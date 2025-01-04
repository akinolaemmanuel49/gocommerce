package queue

import (
	"errors"
	"fmt"
	"log"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishOrderNotification(config *configs.Config, ch *amqp.Channel, orderID, userID, notificationType, message string, logger, errorLogger *log.Logger) error {
	msg := fmt.Sprintf(`{"orderID": "%s", "userID": "%s", "type": "%s", "message": "%s"}`, orderID, userID, notificationType, message)

	if ch == nil {
		return errors.New("DOOM FLAG")
	}

	err := ch.Publish(
		"",                    // Default exchange
		config.OrderQueueName, // Routing key (queue name)
		false,                 // Mandatory
		false,                 // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		errorLogger.Printf("Failed to publish message: %v", err)
		return err
	}

	logger.Printf("Published notification: %s", msg)
	return nil
}
