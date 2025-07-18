package consumers

import (
	"encoding/json"
	"log"
	"product-service/config"
	"product-service/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartProductConsumer(ch *amqp.Channel, cfg *config.Config) {
	msgs, err := ch.Consume(
		cfg.ProductQueue,
		"product-service", // consumers tag
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to register consumers: %v", err)
		return
	}

	go func() {
		for msg := range msgs {
			processProductMessage(msg)
			msg.Ack(false) // 手动确认消息
		}
	}()
}

func processProductMessage(msg amqp.Delivery) {
	var event models.ProductEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return
	}

	switch event.EventType {
	case models.EventProductCreated:
		log.Printf("New product created: %d - %s", event.ProductID, event.ProductData.Name)
	case models.EventProductUpdated:
		log.Printf("Product updated: %d", event.ProductID)
	case models.EventProductDeleted:
		log.Printf("Product deleted: %d", event.ProductID)
	case models.EventCategoryCreated:
		log.Printf("New category created: %d", event.CategoryID)
	case models.EventImageAdded:
		log.Printf("Image added to product %d: %s", event.ProductID, event.ImageData.ImageURL)
	case models.EventAttributeAdded:
		log.Printf("Attribute added to product %d: %s=%s",
			event.ProductID, event.Attribute.Name, event.Attribute.Value)
	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}
}
