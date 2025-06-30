package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func PublishToTopic[T any](data []T, exchangeName string, routingKey string) bool {
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Println("❌ Falló la conexión a RabbitMQ:", err)
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("❌ Falló al abrir un canal:", err)
		return false
	}
	defer ch.Close()

	// Declarar exchange de tipo topic
	err = ch.ExchangeDeclare(
		exchangeName, // nombre del exchange
		"topic",      // tipo de exchange
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println("❌ Falló al declarar el exchange:", err)
		return false
	}

	body, _ := json.Marshal(data)
	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("❌ Falló al publicar:", err)
		return false
	}
	log.Println("✅ Publicado correctamente en el topic", routingKey)
	return true
}