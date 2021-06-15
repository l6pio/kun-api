package core

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/olivere/elastic/v7"
)

type Config struct {
	EsClient *elastic.Client
}

func NewConfig() *Config {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Elasticsearch connection failed: %v", err)
	}
	log.Info("Elasticsearch is connected")

	return &Config{
		EsClient: client,
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
