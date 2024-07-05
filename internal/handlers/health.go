package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping godoc
//
//	@Tags       Health
//	@Produce    plain
//	@Success    200 string  pong
//	@Router     /ping [get]
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
