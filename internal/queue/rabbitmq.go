package queue

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(config *configs.Config, logger, errorLogger *log.Logger) (*amqp.Connection, *amqp.Channel, error) {
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

	err = SetupQueue(ch, logger, errorLogger)
	if err != nil {
		errorLogger.Fatalf("Failed to setup queue: %v", err)
		return nil, nil, err
	}
	return conn, ch, nil
}
