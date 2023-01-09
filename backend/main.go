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
	// x requests per y sec
	// limiterStore := middleware.NewRateLimiterMemoryStoreWithConfig(
	// 	middleware.RateLimiterMemoryStoreConfig{Rate: x/y, Burst: 1},
	// )
	// 2.5 token/sec(5 token/2sec, 1 token/0.4sec) burst: 1
	// limiterStore := middleware.NewRateLimiterMemoryStore(2.5)
	limiterStore := middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: -1.0},
	)
	e.Use(middleware.RateLimiter(limiterStore))

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
