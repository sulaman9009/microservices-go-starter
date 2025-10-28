package grpc_clients

import (
	"ride-sharing/shared/env"
	driverv1 "ride-sharing/shared/gen/go/driver/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DriverServiceClient struct {
	Client driverv1.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*DriverServiceClient, error) {
	driverServiceUrl := env.GetString("DRIVER_SERVICE_URL", "driver-service:9092")
	conn, err := grpc.NewClient(
		driverServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	driverClient := driverv1.NewDriverServiceClient(conn)
	return &DriverServiceClient{
		Client: driverClient,
		conn:   conn,
	}, nil
}

func (c *DriverServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
