package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

func (d Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type Data struct {
	Contacts Contacts
}

func NewData() Data {
	return Data{
		Contacts: []Contact{
			Contact{Name: "John", Email: "jd@gmail.com"},
			Contact{Name: "Clara", Email: "cd@gmail.com"},
		},
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	data := NewData()
	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		data.Contacts = append(data.Contacts, Contact{Name: name, Email: email})

		return c.Render(http.StatusOK, "display", data)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
