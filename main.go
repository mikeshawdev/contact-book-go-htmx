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
	"github.com/mikeshawdev/contact-book-go-htmx/middleware"
	"github.com/mikeshawdev/contact-book-go-htmx/models"
	"github.com/mikeshawdev/contact-book-go-htmx/views"
	"github.com/mikeshawdev/contact-book-go-htmx/views/components"
	errorPages "github.com/mikeshawdev/contact-book-go-htmx/views/errors"
)

var contacts models.Contacts

func main() {
	app := echo.New()
	app.Logger.SetLevel(log.INFO)

	app.Use(echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Generator: func() string {
			return gonanoid.Must()
		},
	}))

	app.Use(echoMiddleware.Logger())
	app.Use(echoMiddleware.Recover())
	app.Use(echoMiddleware.Secure())
	app.Use(echoMiddleware.Gzip())
	app.Use(echoMiddleware.RemoveTrailingSlash())
	app.Use(echoMiddleware.CORS())
	app.Use(middleware.MarkHtmxRequests)

	app.Static("/assets", "public/assets")

	app.GET("/", func(c echo.Context) error {
		formData := models.QuickContactAddFormData{
			Name:  "",
			Email: "",
		}

		return views.Render(c, http.StatusOK, views.Contacts(contacts, formData))
	})

	app.GET("/new", func(c echo.Context) error {
		return views.Render(c, http.StatusOK, views.NewContact())
	})

	app.GET("/settings", func(c echo.Context) error {
		return views.Render(c, http.StatusOK, views.Settings())
	})

	app.POST("/contacts", func(c echo.Context) error {
		formData := models.QuickContactAddFormData{
			Name:  c.FormValue("name"),
			Email: c.FormValue("email"),
		}

		errors := formData.Validate()

		if len(errors) > 0 {
			return views.Render(c, http.StatusBadRequest, components.QuickContactAddForm(formData, errors))
		}

		contact := models.Contact{}.New(formData.Name, formData.Email)
		contacts = contacts.Add(contact)

		views.Render(c, http.StatusCreated, components.QuickContactAddForm(models.QuickContactAddFormData{}, nil))
		return views.Render(c, http.StatusCreated, components.OobContact(contact))
	})

	app.RouteNotFound("*", func(c echo.Context) error {
		return views.Render(c, http.StatusNotFound, errorPages.NotFound())
	})

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)

		views.Render(c, http.StatusNotFound, errorPages.InternalServerError())
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
