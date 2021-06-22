package core

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/v1/router/vo/img"
)

type Config struct {
	DbConn        *sql.DB
	ImageUpEvents chan img.Key
}

func NewConfig() *Config {
	return &Config{
		DbConn:        db.Connect(),
		ImageUpEvents: make(chan img.Key, 10000),
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
