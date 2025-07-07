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

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: NewData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			FormData := newFormData()
			FormData.Values["name"] = name
			FormData.Values["email"] = email
			FormData.Errors["email"] = "email already exists"
			return c.Render(http.StatusUnprocessableEntity, "form", FormData)
		}

		contact := Contact{
			Name:  name,
			Email: email,
		}
		page.Data.Contacts = append(page.Data.Contacts, Contact{Name: name, Email: email})

		return c.Render(http.StatusOK, "outofband-contact", contact)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
