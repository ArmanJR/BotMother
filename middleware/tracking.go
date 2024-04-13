package middleware

import (
	"github.com/labstack/echo/v4"
)

func Tracking(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implement tracking logic
		return next(c)
	}
}
