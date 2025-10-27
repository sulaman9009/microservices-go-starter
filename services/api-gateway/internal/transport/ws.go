package transport

import (
	"context"
	"fmt"
	"ride-sharing/services/api-gateway/internal/domain"
	"ride-sharing/services/api-gateway/internal/problems"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

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
				return nil
			}
			s.logger.Err(err).Msg("error reading from websocket")
			return nil
		}
		s.logger.Info().Str("userID", userID).Msgf("received: %s", msg)
	}
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
	fmt.Println("queryparams:", queryParams)
	if err := c.Validate(&queryParams); err != nil {
		return err
	}

	ctx := context.Background()
	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: domain.Driver{
			ID:             queryParams.UserID,
			Name:           "Sulaman Ahmad",
			ProfilePicture: util.GetRandomAvatar(2),
			CarPlate:       "ABCD-1234",
			PackageSlug:    queryParams.PackageSlug,
		},
	}

	if err := wsjson.Write(ctx, ws, msg); err != nil {
		return err
	}

	for {
		_, msg, err := ws.Read(c.Request().Context())
		if err != nil {
			s.logger.Err(err).Msg("error reading from websocket")
			break
		}
		fmt.Printf("received: %s\n", msg)
	}

	return nil
}
