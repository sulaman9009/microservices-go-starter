package events

import (
	"context"
	"encoding/json"
	"ride-sharing/services/trip-service/internal/domain"
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

func (p *tripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModel) error {
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.rabbitmq.PublishMessage(ctx, contracts.TripEventCreated, &contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    data,
	})
}
