package web

import (
	"context"
	"gowebecho/web/auth"
	"gowebecho/web/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	e.HideBanner = true
	e.DisableHTTP2 = true

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session_secret"))))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "web/public/assets")

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
	e.GET("/", handlers.Root)
	e.GET("/auth", func(c echo.Context) error {
		prov, err := auth.NewOktaProvider()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create auth provider")
		}
		return prov.Authorize(c, "state", "nonce")
	})
}

func setupTemplates() {
	t = NewTemplate("web/public/views/*.tmpl")
	e.Renderer = t
}
