package transport

import (
	tripv1 "ride-sharing/shared/gen/go/trip/v1"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

func (p *previewTripRequest) toProto() *tripv1.PreviewTripRequest {
	return &tripv1.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &tripv1.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		EndLocation: &tripv1.Coordinate{
			Latitude:  p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID" validate:"required"`
	UserID     string `json:"userID" validate:"required"`
}

func (c *startTripRequest) toProto() *tripv1.CreateTripRequest {
	return &tripv1.CreateTripRequest{
		RideFareID: c.RideFareID,
		UserID:     c.UserID,
	}
}
