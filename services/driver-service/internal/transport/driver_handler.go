package transport

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
	driverv1 "ride-sharing/shared/gen/go/driver/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type driverGrpcHandler struct {
	service domain.DriverService

	// This is embedded for forward compatibility.
	// It embeds all required methods to satisfy the grpc contract, stubs each method by returning nil, (unimplemented) err.
	// for each stubbed method, we will provide our own actual implementation below.
	driverv1.UnimplementedDriverServiceServer
}

func NewDriverGrpcHandler(service domain.DriverService) *driverGrpcHandler {
	return &driverGrpcHandler{
		service: service,
	}
}

// Implement gRPC methods here
func (h *driverGrpcHandler) RegisterDriver(ctx context.Context, req *driverv1.RegisterDriverRequest) (*driverv1.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(req.GetDriverID(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver: %v", err)
	}

	return &driverv1.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *driverGrpcHandler) UnregisterDriver(ctx context.Context, req *driverv1.RegisterDriverRequest) (*driverv1.RegisterDriverResponse, error) {
	h.service.UnregisterDriver(req.GetDriverID())
	return &driverv1.RegisterDriverResponse{
		Driver: &driverv1.Driver{
			Id: req.GetDriverID(),
		},
	}, nil
}
