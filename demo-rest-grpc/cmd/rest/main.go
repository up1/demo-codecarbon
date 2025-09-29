package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/hello", func(c echo.Context) error {
		msg := c.QueryParam("message")
		if msg == "" {
			msg = "hello world"
		}
		return c.JSON(http.StatusOK, map[string]string{"message": msg})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
