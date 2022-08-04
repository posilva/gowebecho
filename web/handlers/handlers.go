package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Root(c echo.Context) error {
	// Get session data

	return c.String(http.StatusOK, "it works")
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "world")
}

func Image(c echo.Context) error {

	return c.Render(http.StatusOK, "image", nil)
}

func Cookie(c echo.Context) error {
	cookie, err := c.Cookie("cookie_test")
	if err != nil {
		cookie = new(http.Cookie)
		cookie.Name = "cookie_test"
		cookie.Value = "new cookie data"
	}

	cookie.Value = "new " + cookie.Value
	c.SetCookie(cookie)
	return c.Render(http.StatusOK, "cookie", cookie.Value)
}

func Session(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["foo"] = "bar"
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func Hello(c echo.Context) error {
	n := c.Param("name")
	return c.Render(http.StatusOK, "hello", n)
}
