package transport

import (
	"net/http"
	"ride-sharing/services/api-gateway/internal/problems"
	util "ride-sharing/services/api-gateway/internal/utils"
	"ride-sharing/shared/contracts"

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

	resp := contracts.APIResponse{
		Data: "ok",
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *server) startTrip(c echo.Context) error {
	return nil
}
