package main

import (
	_ "embed"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/k8s"
	"l6p.io/kun/api/pkg/core/service"
	"l6p.io/kun/api/pkg/v1/router"
	"net/http"
	"os"
)

func main() {
	var mongodbAddr = flag.String("mongodb-addr", "localhost:32017", "The mongodb connection address")
	var mongodbUser = flag.String("mongodb-user", "root", "The mongodb username")
	var mongodbPass = flag.String("mongodb-pass", "rootpassword", "The mongodb password")
	var master = flag.String("master", "", "the address of the Kubernetes API server.")
	var kubeConfig = flag.String("kubeconfig", "", "path to a kubeconfig.")
	flag.Parse()

	conf, err := core.NewConfig(*mongodbAddr, *mongodbUser, *mongodbPass, *master, *kubeConfig)
	if err != nil {
		log.Fatalf("initialization of configuration failed: %v", err)
	}

	err = service.SaveGrypeConfigFile(os.Getenv("REGISTRY"))
	if err != nil {
		log.Fatalf("save Grype configuration file failed: %v", err)
	}

	server := echo.New()
	server.HideBanner = true

	server.Use(middleware.CORS())

	// Passing config into the router's context.
	server.Use(core.WithConfig(conf)...)

	core.AddValidator(server)

	apiV1 := server.Group("/api/v1")
	router.PingRouter(apiV1.Group("/ping"))
	router.PodRouter(apiV1.Group("/pod"))
	router.ImageRouter(apiV1.Group("/image"))
	router.ArtifactRouter(apiV1.Group("/artifact"))
	router.CveRouter(apiV1.Group("/cve"))
	router.GarbageRouter(apiV1.Group("/garbage"))

	server.HTTPErrorHandler = ErrorHandler

	if err := db.RemovePods(conf); err != nil {
		log.Fatal(err)
	}

	if err := db.CleanImagePods(conf); err != nil {
		log.Fatal(err)
	}

	go service.PeriodicallyUpdateVulnerabilityDatabase()
	go k8s.StartPodInformer(conf)

	if err := server.Start(":1323"); err != nil {
		log.Fatalf("server startup failed: %v", err)
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
