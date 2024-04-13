package middleware

import (
	"github.com/labstack/echo/v4"
)

func BlockUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implement blocking logic
		//if /* condition to block */ {
		//	return c.String(http.StatusForbidden, "Blocked")
		//}
		return next(c)
	}
}
