package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			Response{
				Message: "Hello World",
			},
		)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
