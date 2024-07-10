package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mikeshawdev/contact-book-go-htmx/components"
	errorPages "github.com/mikeshawdev/contact-book-go-htmx/components/errors"
	"github.com/mikeshawdev/contact-book-go-htmx/middleware"
)

type PageData struct {
	PageName string
}

func main() {
	app := echo.New()
	app.Logger.SetLevel(log.INFO)

	app.Use(echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Generator: func() string {
			id, err := gonanoid.New()

			if err != nil {
				panic(err)
			}

			return id
		},
	}))

	app.Use(echoMiddleware.Logger())
	app.Use(echoMiddleware.Recover())
	app.Use(echoMiddleware.Secure())
	app.Use(echoMiddleware.Gzip())
	app.Use(echoMiddleware.RemoveTrailingSlash())
	app.Use(echoMiddleware.CSRF())
	app.Use(echoMiddleware.CORS())
	app.Use(middleware.MarkHtmxRequests)

	app.Static("/assets", "public/assets")

	app.GET("/", func(c echo.Context) error {
		return components.Render(c, http.StatusOK, components.Contacts())
	})

	app.GET("/new", func(c echo.Context) error {
		return components.Render(c, http.StatusOK, components.NewContact())
	})

	app.GET("/settings", func(c echo.Context) error {
		return components.Render(c, http.StatusOK, components.Settings())
	})

	app.RouteNotFound("*", func(c echo.Context) error {
		return components.Render(c, http.StatusNotFound, errorPages.NotFound())
	})

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)

		components.Render(c, http.StatusNotFound, errorPages.InternalServerError())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := app.Start(":1323"); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
}
