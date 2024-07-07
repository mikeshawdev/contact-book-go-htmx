package middleware

import (
	"github.com/labstack/echo/v4"
)

func MarkHtmxRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("HX-Request") == "true" {
			c.Set("htmx-request", true)
		} else {
			c.Set("htmx-request", false)
		}

		return next(c)
	}
}
