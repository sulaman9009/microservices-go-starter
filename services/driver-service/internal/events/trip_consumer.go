package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *tripConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var tripEvents contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvents); err != nil {
			return fmt.Errorf("error unmarshaling trip event: %v", err)
		}
		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvents.Data, &payload); err != nil {
			return fmt.Errorf("error unmarshaling trip event data: %v", err)
		}
		log.Printf("driver received message: %+v", payload)
		return nil
	})
}
