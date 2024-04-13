package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckSecretToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headerValue := c.Request().Header.Get("secret_token")
		botToken := "bozbozi"
		if headerValue != botToken {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Forbidden",
			})
		}
		return next(c)
	}
}
