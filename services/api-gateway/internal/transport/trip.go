package transport

import (
	"net/http"
	"ride-sharing/services/api-gateway/internal/problems"
	"ride-sharing/shared/contracts"

	"github.com/labstack/echo/v4"
)

func (s *server) previewTrip(c echo.Context) error {
	var req previewTripRequest
	if err := c.Bind(&req); err != nil {
		return problems.NewBadRequest("invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	tripPreview, err := s.tripClient.Client.PreviewTrip(c.Request().Context(), req.toProto())
	if err != nil {
		return problems.NewInternal("failed to preview trip", err.Error())
	}

	resp := contracts.APIResponse{
		Data: tripPreview,
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *server) startTrip(c echo.Context) error {
	var req startTripRequest
	if err := c.Bind(&req); err != nil {
		return problems.NewBadRequest("invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	trip, err := s.tripClient.Client.CreateTrip(c.Request().Context(), req.toProto())
	if err != nil {
		return problems.NewInternal("failed to start trip", err.Error())
	}
	resp := contracts.APIResponse{
		Data: trip,
	}
	return c.JSON(http.StatusOK, resp)
}
