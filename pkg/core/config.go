package core

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
	"l6p.io/kun/api/pkg/v1/router/vo/scan"
)

type Config struct {
	EsClient     *elastic.Client
	ScanRequests chan scan.Key
}

func NewConfig() *Config {
	client, err := elastic.NewClient(
		// TODO: Put the urls of ES into the configuration.
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Elasticsearch connection failed: %v", err)
	}
	log.Info("Elasticsearch is connected")

	return &Config{
		EsClient:     client,
		ScanRequests: make(chan scan.Key, 10000),
	}
}

func WithConfig(conf *Config) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				ctx.Set("config", conf)
				return next(ctx)
			}
		},
	}
}
