package main

import (
	"fmt"
	"html/template"
	"io"
	"lenkr/db"
	"lenkr/lib"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	godotenv.Load()
	var err error

	conn := db.ConnectToDb()
	defer conn.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Renderer = NewTemplates()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/admin", func(c echo.Context) error {
		links := db.GetLinks(conn)
		return c.Render(http.StatusOK, "links.html", links)
	})

	e.POST("/shorten", func(c echo.Context) error {
		url := c.FormValue("url")
		shortUrl := c.FormValue("shortUrl")

		if url == "" {
			return c.String(http.StatusBadRequest, "URL is required")
		}

		if shortUrl == "" {
			shortUrl, err = lib.GenerateRandomStringURLSafe(6)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Error generating random string")
			}
		}

		err := db.InsertLink(conn, url, shortUrl)
		if err != nil {
			fmt.Println("Error inserting link: ", err)
			return c.String(http.StatusInternalServerError, "Error inserting link")
		}

		return c.Render(http.StatusOK, "shortened-url-item", db.Link{ShortUrl: shortUrl, Url: url})
	})

	e.GET("/l/:shortUrl", func(c echo.Context) error {
		shortUrl := c.Param("shortUrl")
		link, err := db.GetLink(conn, shortUrl)
		if err != nil {
			fmt.Println("Error getting link: ", err)
			return c.String(http.StatusInternalServerError, "Error getting link")
		}

		err = c.Redirect(http.StatusMovedPermanently, link.Url)
		if err == nil {
			db.IncrementFetches(conn, link)
		}

		return err
	})

	port := getPort()
	e.Logger.Fatal(e.Start(":" + port))

}
