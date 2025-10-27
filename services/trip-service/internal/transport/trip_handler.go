package transport

import (
	"context"
	"fmt"
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

	userID := req.GetUserID()
	estimatedFares := h.tripService.EstimatePackagesPriceWithRoute(route)
	fares, err := h.tripService.GenerateTripFares(ctx, estimatedFares, userID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares: %v", err)
	}

	return &tripv1.PreviewTripResponse{
		TripID:    "1234",
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}

func (h *tripGrpcHandler) CreateTrip(ctx context.Context, req *tripv1.CreateTripRequest) (*tripv1.CreateTripResponse, error) {
	fmt.Println("in func")
	fareID := req.GetRideFareID()
	userID := req.GetUserID()

	rideFare, err := h.tripService.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate the fare: %v", err)
	}

	trip, err := h.tripService.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the trip: %v", err)
	}

	// Add a comment at the end of the function to publish an event on the Async Comms module.

	return &tripv1.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
