package transport

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	tripv1 "ride-sharing/shared/gen/go/trip/v1"
)

type tripGrpcHandler struct {
	tripService domain.TripService

	// This is embedded for forward compatibility.
	// It embeds all required methods to satisfy the grpc contract, stubs each method by returning nil, (unimplemented) err.
	// for each stubbed method, we will provide our own actual implementation below.
	tripv1.UnimplementedTripServiceServer
}

func NewTripGrpcHandler(tripService domain.TripService) *tripGrpcHandler {
	return &tripGrpcHandler{
		tripService: tripService,
	}
}

func (h *tripGrpcHandler) PreviewTrip(
	context.Context,
	*tripv1.PreviewTripRequest,
) (*tripv1.PreviewTripResponse, error) {
	return nil, nil
}
