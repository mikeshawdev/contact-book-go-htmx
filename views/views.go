package views

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, statusCode int, t []templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	for _, component := range t {
		if err := component.Render(ctx.Request().Context(), buf); err != nil {
			return err
		}
	}

	return ctx.HTML(statusCode, buf.String())
}
