package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewTripService(repo domain.TripRepository) *TripService {
	return &TripService{
		repository: repo,
	}
}

type TripService struct {
	repository domain.TripRepository
}

func (s *TripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	return s.repository.CreateTrip(ctx, &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	})
}

func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*domain.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routeResp domain.OsrmApiResponse

	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &routeResp, nil
}
