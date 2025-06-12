package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func Publish[T any](data []T, queueName string) bool {
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel", err)
		return false
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Println("Queue declare error", err)
		return false
	}

	body, _ := json.Marshal(data)
	err = ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("Publish error", err)
		return false
	}
	return true
}