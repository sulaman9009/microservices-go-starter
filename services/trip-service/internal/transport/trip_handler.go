package transport

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	tripv1 "ride-sharing/shared/gen/go/trip/v1"
	"ride-sharing/shared/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *tripGrpcHandler) PreviewTrip(ctx context.Context, req *tripv1.PreviewTripRequest) (*tripv1.PreviewTripResponse, error) {
	route, err := h.tripService.GetRoute(
		ctx,
		&types.Coordinate{
			Latitude:  req.StartLocation.GetLatitude(),
			Longitude: req.StartLocation.GetLongitude(),
		},
		&types.Coordinate{
			Latitude:  req.EndLocation.GetLatitude(),
			Longitude: req.EndLocation.GetLongitude(),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}
	return &tripv1.PreviewTripResponse{
		TripID:    "1234",
		Route:     route.ToProto(),
		RideFares: []*tripv1.RideFare{},
	}, nil
}
