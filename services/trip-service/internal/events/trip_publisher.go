package events

import (
	"context"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type tripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *tripEventPublisher {
	return &tripEventPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *tripEventPublisher) PublishTripCreated(ctx context.Context) error {
	return p.rabbitmq.PublishMessage(ctx, contracts.TripEventCreated, "trip created event")
}
