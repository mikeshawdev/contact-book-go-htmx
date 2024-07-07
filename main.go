package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mikeshawdev/contact-book-go-htmx/middleware"
	"github.com/mikeshawdev/contact-book-go-htmx/templates"
)

type PageData struct {
	PageName string
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Renderer = templates.TemplateRenderer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	e.Use(echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Generator: func() string {
			id, err := gonanoid.New()

			if err != nil {
				panic(err)
			}

			return id
		},
	}))

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Secure())
	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.RemoveTrailingSlash())
	e.Use(echoMiddleware.CSRF())
	e.Use(echoMiddleware.CORS())
	e.Use(middleware.MarkHtmxRequests)

	e.Static("/assets", "public/assets")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "contacts.html", PageData{PageName: "contacts"})
	})

	e.GET("/new", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new-contact.html", PageData{PageName: "new-contact"})
	})

	e.GET("/settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "settings.html", PageData{PageName: "settings"})
	})

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
