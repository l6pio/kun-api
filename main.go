package main

import (
	_ "embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/img"
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
	router.ImageRouter(apiV1.Group("/img"))
	router.CveRouter(apiV1.Group("/cve"))

	server.HTTPErrorHandler = ErrorHandler

	db.Init(conf.DbConn)

	// Read and process CVE scan requests
	WaitForScanRequests(conf)

	if err := server.Start(":1323"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}

func WaitForScanRequests(conf *core.Config) {
	go func() {
		for {
			key := <-conf.ImageUpEvents
			imageName := fmt.Sprintf("%s:%s", key.ImageRepo, key.ImageTag)

			report := cve.Scan(imageName)

			imageId := report.Source.Target.ImageID
			id, err := img.Status(conf.DbConn, imageId, imageName, img.StatusUp)
			if err != nil {
				log.Error(err)
				continue
			}

			pickId, err := img.PickId(conf.DbConn, imageId)

			if id != pickId {
				log.Infof("Image '%s' scan report does not need to be updated repeatedly")
				continue
			}

			exists, err := img.Exists(conf.DbConn, imageId)
			if err != nil {
				log.Error(err)
				continue
			}

			if exists {
				log.Infof("Image '%s' scan report already exists", imageId)
				continue
			}

			if len(report.Matches) == 0 {
				log.Info("No vulnerabilities found")
			}

			log.Info("Start saving scan results")

			err = cve.Insert(conf.DbConn, report)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Infof("Scan results of '%v' has been saved", imageId)
		}
	}()
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
