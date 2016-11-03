package RPRS

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

func upload(c echo.Context) error {
	// Read form fields
	repo := c.FormValue("repo")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("data")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s .</p>", file.Filename, repo))
}
func main() {

	initializeConfig()
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Static("public"))

	e.POST("/upload", upload)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Run(standard.New(viper.GetString("Port")))
}
