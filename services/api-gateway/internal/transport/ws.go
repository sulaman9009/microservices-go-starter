package transport

import (
	"context"
	"ride-sharing/services/api-gateway/internal/problems"
	"ride-sharing/shared/contracts"
	driverv1 "ride-sharing/shared/gen/go/driver/v1"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/labstack/echo/v4"
)

type driverWSQueryParams struct {
	UserID      string `query:"userID" validate:"required"`
	PackageSlug string `query:"packageSlug" validate:"required"`
}

func (s *server) handleRiderWS(c echo.Context) error {
	ws, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{
		OriginPatterns:     []string{"*"},
		InsecureSkipVerify: true, // Only in development
	})
	if err != nil {
		s.logger.Err(err).Msg("websocket accept error")
		return err
	}
	defer ws.CloseNow()

	userID := c.QueryParam("userID")
	if userID == "" {
		ws.Close(websocket.StatusUnsupportedData, "userID is required")
		return problems.NewBadRequest("userID is required", "")
	}

	ctx := c.Request().Context()
	s.logger.Info().Str("userID", userID).Msg("rider ws connected")

	for {
		_, msg, err := ws.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				s.logger.Info().Str("userID", userID).Msg("rider ws closed normally")
				break
			}
			s.logger.Err(err).Msg("error reading from websocket")
			break
		}
		s.logger.Info().Msgf("rider ws received: %s", msg)
	}
	return nil
}

func (s *server) handleDriverWS(c echo.Context) error {
	ws, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return err
	}
	defer ws.CloseNow()
	var queryParams driverWSQueryParams
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &queryParams); err != nil {
		return problems.NewBadRequest(err.Error(), "")
	}
	if err := c.Validate(&queryParams); err != nil {
		return err
	}

	ctx := context.Background()
	defer func() {
		s.driverClient.Client.UnregisterDriver(ctx, &driverv1.RegisterDriverRequest{
			DriverID:    queryParams.UserID,
			PackageSlug: queryParams.PackageSlug,
		})
	}()

	driverData, err := s.driverClient.Client.RegisterDriver(ctx, &driverv1.RegisterDriverRequest{
		PackageSlug: queryParams.PackageSlug,
		DriverID:    queryParams.UserID,
	},
	)
	if err != nil {
		return err
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
	}

	if err := wsjson.Write(ctx, ws, msg); err != nil {
		return err
	}

	for {
		_, msg, err := ws.Read(c.Request().Context())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				s.logger.Info().Str("userID", queryParams.UserID).Msg("rider ws closed normally")
				break
			}
			s.logger.Err(err).Msg("error reading from websocket")
			break
		}
		s.logger.Info().Msgf("driver ws received: %s", msg)
	}

	return nil
}
