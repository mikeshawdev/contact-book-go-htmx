package middleware

import (
	"github.com/labstack/echo/v4"
)

func MarkHtmxRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("htmx-request", c.Request().Header.Get("HX-Request") == "true")

		return next(c)
	}
}
