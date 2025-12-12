package domain

import (
	"context"
)

type TripEventPublisher interface {
	PublishTripCreated(ctx context.Context, trip *TripModel) error
}
