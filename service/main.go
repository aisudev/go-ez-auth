package main

import (
	"net/http"
	"service-api/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(middleware.AuthMiddleware())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Start(":9999")
}
