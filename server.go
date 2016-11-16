package main

import (
	"github.com/radepal/RPRS/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"net/http"
)

func main() {

	initializeConfig()
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.Static("public"))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// Login route
	e.POST("/login", controller.Login)
	// Restricted group
	r := e.Group("/rpm")
	r.Use(middleware.JWT([]byte(viper.GetString("Secret"))))
	r.PUT("/upload", controller.Upload)

	e.Run(standard.New(viper.GetString("Port")))
}
