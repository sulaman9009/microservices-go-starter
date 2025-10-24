package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Instantiates a new logger
func New() *zerolog.Logger {

	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	return &logger
}
