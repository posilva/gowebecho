package web

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type PageTemplate struct {
	templates *template.Template
}

// NewTemplate creates a new instance of a Template struct pointer
func NewTemplate(globPath string) *PageTemplate {
	return &PageTemplate{
		templates: template.Must(template.ParseGlob(globPath)),
	}
}

func (t *PageTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
