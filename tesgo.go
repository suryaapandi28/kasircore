package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// route sederhana
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo + Go + Docker!")
	})

	// listen di port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
