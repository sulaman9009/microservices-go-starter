package main

import (
	"ride-sharing/services/api-gateway/internal/grpc_clients"
	"ride-sharing/services/api-gateway/internal/logger"
	"ride-sharing/services/api-gateway/internal/transport"

	"github.com/rs/zerolog"
)

func main() {
	log := logger.New()
	if err := run(log); err != nil {
		log.Fatal().Err(err).Msg("application error")
	}
}

func run(logger *zerolog.Logger) error {
	tripClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		return err
	}
	srv := transport.NewHTTPServer(logger, tripClient)
	srv.Start()
	return nil
}
