package bots

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleBotSecond(c echo.Context) error {
	return c.String(http.StatusOK, "Handled bot second")
}
