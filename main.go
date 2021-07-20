package main

import (
	_ "embed"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/service"
	"l6p.io/kun/api/pkg/v1/router"
	"net/http"
)

func main() {
	var mongodbAddr string
	var mongodbUser string
	var mongodbPass string
	flag.StringVar(&mongodbAddr, "mongodb-addr", "localhost:32017", "The mongodb connection address")
	flag.StringVar(&mongodbUser, "mongodb-user", "root", "The mongodb username")
	flag.StringVar(&mongodbPass, "mongodb-pass", "rootpassword", "The mongodb password")
	flag.Parse()

	conf := core.NewConfig(mongodbAddr, mongodbUser, mongodbPass)

	server := echo.New()
	server.HideBanner = true

	server.Use(middleware.CORS())

	// Passing config into the router's context.
	server.Use(core.WithConfig(conf)...)

	core.AddValidator(server)

	apiV1 := server.Group("/api/v1")
	router.PingRouter(apiV1.Group("/ping"))
	router.ImageRouter(apiV1.Group("/image"))
	router.ArtifactRouter(apiV1.Group("/artifact"))
	router.CveRouter(apiV1.Group("/cve"))

	server.HTTPErrorHandler = ErrorHandler

	// Read and process CVE scan requests
	WaitForImageEvents(conf)

	if err := server.Start(":1323"); err != nil {
		log.Fatalf("server startup failed: %v", err)
	}
}

func WaitForImageEvents(conf *core.Config) {
	go func() {
		for {
			imageEvent := <-conf.ImageEvents
			if *imageEvent.Status == core.ImageUp {
				if err := ProcessImageUp(conf, imageEvent); err != nil {
					log.Error(err)
				}
			} else {
				if err := ProcessImageDown(conf, imageEvent); err != nil {
					log.Error(err)
				}
			}
		}
	}()
}

func ProcessImageUp(conf *core.Config, event core.ImageEvent) error {
	report := service.Scan(event.Image)

	if len(report.Matches) == 0 {
		log.Info("no vulnerabilities found")
	}

	imageId := report.Source.Target.ImageID
	if err := db.SaveImageStatus(conf, event.Timestamp, imageId, *event.Status); err != nil {
		log.Error(err)
	}

	log.Info("start saving scan results")
	service.Insert(conf, report)
	log.Infof("scan results of '%v' has been saved", imageId)
	return db.UpdateImageUsage(conf, imageId, core.ImageUp)
}

func ProcessImageDown(conf *core.Config, event core.ImageEvent) error {
	report := service.Scan(event.Image)

	imageId := report.Source.Target.ImageID
	if err := db.SaveImageStatus(conf, event.Timestamp, imageId, *event.Status); err != nil {
		log.Error(err)
	}
	return db.UpdateImageUsage(conf, imageId, core.ImageDown)
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
