package transport

import (
	"fmt"
	"net/http"
	"ride-sharing/services/api-gateway/internal/problems"
	util "ride-sharing/services/api-gateway/internal/utils"

	"github.com/labstack/echo/v4"
)

func (s *server) previewTrip(c echo.Context) error {
	var req previewTripRequest
	if err := c.Bind(&req); err != nil {
		return problems.NewBadRequest("invalid request payload", err.Error())
	}
	if err := util.ValidatePayload(&req); err != nil {
		return problems.NewBadRequest("validation failed", err.Error())
	}

	fmt.Printf("preview trip request:", req)
	return c.String(http.StatusOK, "preview trip endpoint")
}

func (s *server) startTrip(c echo.Context) error {
	return nil
}
