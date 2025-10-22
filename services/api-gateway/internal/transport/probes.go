package transport

import "github.com/labstack/echo/v4"

func (s *server) readinessProbe(c echo.Context) error {
	return c.NoContent(200)
}

func (s *server) livenessProbe(c echo.Context) error {
	return c.NoContent(200)
}
