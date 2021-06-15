package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/es"
	"l6p.io/kun/api/pkg/v1/router"
	"net/http"
)

func main() {
	conf := core.NewConfig()

	server := echo.New()
	server.HideBanner = true

	server.Use(middleware.CORS())

	// Passing config into the router's context.
	server.Use(core.WithConfig(conf)...)

	core.AddValidator(server)

	apiV1 := server.Group("/api/v1")
	router.PingRouter(apiV1.Group("/ping"))
	router.CveRouter(apiV1.Group("/cve"))

	server.HTTPErrorHandler = ErrorHandler

	// Create a new index when the index does not exist.
	es.CreateIndex(conf)

	if err := server.Start(":1323"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}

func ErrorHandler(err error, ctx echo.Context) {
	log.Errorf("API handling error: %v", err)

	code := http.StatusBadRequest
	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
	}

	err = ctx.JSON(code, struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})
	if err != nil {
		log.Error(err)
	}
}
