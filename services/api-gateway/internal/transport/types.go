package transport

import "ride-sharing/shared/types"

type previewTripRequest struct {
	UserID      string           `json:"user_id" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}
