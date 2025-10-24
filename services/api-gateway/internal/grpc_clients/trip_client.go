package grpc_clients

import (
	"ride-sharing/shared/env"
	tripv1 "ride-sharing/shared/gen/go/trip/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TripServiceClient struct {
	Client tripv1.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {
	tripServiceUrl := env.GetString("TRIP_SERVICE_URL", "trip-service:9093")
	conn, err := grpc.NewClient(
		tripServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	tripClient := tripv1.NewTripServiceClient(conn)
	return &TripServiceClient{
		Client: tripClient,
		conn:   conn,
	}, nil
}

func (c *TripServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
