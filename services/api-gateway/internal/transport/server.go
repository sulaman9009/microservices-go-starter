package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/services/api-gateway/internal/problems"
	"ride-sharing/shared/env"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

type server struct {
	mux    *echo.Echo
	logger *zerolog.Logger
}

func NewHTTPServer(logger *zerolog.Logger) *server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Timeout())
	e.Use(middleware.RequestID())

	srv := &server{
		mux:    e,
		logger: logger,
	}
	e.HTTPErrorHandler = srv.ErrorHandler

	return srv
}

func (s *server) Start() {
	s.mountHandlers()

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: s.mux,
	}

	s.logger.Info().Msgf("Starting HTTP server on %s", httpAddr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	if err := s.waitForShutdown(srv); err != nil {
		s.logger.Fatal().Err(err).Msgf("failed to shutdown server correctly")
	}
}

func (s *server) ErrorHandler(err error, c echo.Context) {
	var prob *problems.Problem
	if errors.As(err, &prob) {
		// Log internal message if present
		if prob.InternalDetail != "" {
			s.logger.
				Error().
				Err(err).
				Msgf("[ERROR] %s (status=%d): %s", prob.InternalDetail, prob.Status, prob.Detail)
		} else {
			s.logger.Error().Err(err).Msgf("[ERROR] %s (status=%d)", prob.Detail, prob.Status)
		}

		// Send RFC 9457-safe response
		_ = c.JSON(prob.Status, prob)
		return
	}
	// Default case for unknown errors
	s.logger.Error().Err(err).Msgf("[ERROR] %v", err)
	_ = c.JSON(http.StatusInternalServerError, &problems.Problem{
		Type:   "about:blank",
		Title:  http.StatusText(http.StatusInternalServerError),
		Status: http.StatusInternalServerError,
		Detail: "unexpected server error",
	})
}

func (s *server) waitForShutdown(server *http.Server) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// block here
	<-sig
	s.logger.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
