package transport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/events"
	"ride-sharing/services/trip-service/internal/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	tripv1 "ride-sharing/shared/gen/go/trip/v1"
	"ride-sharing/shared/messaging"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var (
	grpcPort    = 9093
	rabbitmqURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
)

type gRPCServer struct {
	logger *zerolog.Logger
}

func NewGRPCServer(logger *zerolog.Logger) *gRPCServer {
	return &gRPCServer{
		logger: logger,
	}
}

func (s *gRPCServer) Start() error {
	// create tcp connection
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		return err
	}

	// create server
	server := grpc.NewServer()

	// create messaging client
	msgClient, err := messaging.NewRabbitMQClient(rabbitmqURI)
	if err != nil {
		return err
	}
	defer msgClient.Close()

	// create event publisher
	publisher := events.NewTripEventPublisher(msgClient)

	// register trip service
	tripRepo := repository.NewInMemRepository()
	tripService := service.NewTripService(tripRepo)
	tripv1.RegisterTripServiceServer(server, NewTripGrpcHandler(tripService, publisher))

	// listen and serve
	s.logger.Info().Msgf("server listening on port: %d", grpcPort)
	go func() {
		if err := server.Serve(listen); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("failed to start grpc server")
		}
	}()

	return s.waitForShutdown(server)
}

func (s *gRPCServer) waitForShutdown(grpcServer *grpc.Server) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sig
	s.logger.Info().Msg("shutting down gRPC server...")

	// Use a goroutine and context timeout to enforce graceful shutdown limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.logger.Warn().Msg("graceful shutdown timed out; forcing stop")
		grpcServer.Stop()
	case <-stopped:
		s.logger.Info().Msg("gRPC server shut down gracefully")
	}

	return nil
}
