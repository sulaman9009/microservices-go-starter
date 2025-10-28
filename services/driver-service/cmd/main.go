package main

import (
	"ride-sharing/services/driver-service/internal/transport"
	"ride-sharing/shared/logger"
)

func main() {
	logger := logger.New()
	server := transport.NewGRPCServer(logger)
	if err := server.Start(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start trip service")
	}
}
