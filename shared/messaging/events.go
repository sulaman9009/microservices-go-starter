package messaging

import tripv1 "ride-sharing/shared/gen/go/trip/v1"

const (
	FindAvailableDriversQueue = "find_available_drivers"
)

type TripEventData struct {
	Trip *tripv1.Trip `json:"trip"`
}
