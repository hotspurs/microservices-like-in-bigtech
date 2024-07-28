package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Chat!")
	})
	e.GET("/live", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	e.GET("/ready", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/auth", func(c echo.Context) error {
		res, err := http.Get("http://auth:8080")
		if err != nil {
		}
		fmt.Println(res)
		return c.String(res.StatusCode, "")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
