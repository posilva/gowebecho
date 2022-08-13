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

// WebService holds the data to  interact with the service
type WebService struct {
	e            *echo.Echo
	authProvider *auth.OktaProvider
}

var (
	Server WebService
)

func init() {
	ap := auth.NewOktaProvider()
	Server = WebService{
		e:            echo.New(),
		authProvider: ap,
	}
}

// Setup the web server
func (s *WebService) Setup() error {
	s.e.Renderer = NewTemplate("web/public/views/*.tmpl")

	s.e.HideBanner = true
	s.e.DisableHTTP2 = true
	s.e.Use(session.Middleware(sessions.NewCookieStore([]byte("session_secret"))))
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Static("/static", "web/public/assets")
	s.setupRoutes()
	return nil

}

// Serve handles server initialization with support for graceful shutdown
func (s *WebService) Serve(address string) {

	go func() {
		if err := s.e.Start(address); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatal("shutting down the server")
		}
	}()

	// create a limited buffered channel to receive signals
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Fatal(err)
	}

	s.e.Logger.Fatal(s.e.Start(address))
}

func (s *WebService) setupRoutes() {
	s.e.GET("/", handlers.Root)
	s.e.GET("/auth", func(c echo.Context) error {
		url, err := s.authProvider.AuthorizeURL(c, c.Request().URL.String(), "nonce")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get authorization url")
		}

		s.e.Logger.Warnf("url: %s", url)
		return c.Redirect(http.StatusMovedPermanently, url)
	})

	s.e.GET("/authorization-code/callback", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

}
