package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e *echo.Echo
	t *PageTemplate
)

// Setup the web server
func Setup() error {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	setupTemplates()
	setupRoutes()
	return nil

}

//Serve handles server initialization with support for graceful shutdown
func Serve(address string) {

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// create a limited buffered channel to receive signals
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(address))
}

func setupRoutes() {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "it works")
	})

	e.GET("/index", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", "World")
	})

	e.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.Render(http.StatusOK, "hello", name)
	})

}

func setupTemplates() {
	t = NewTemplate("web/public/views/*.tmpl")
	e.Renderer = t
}
