package core

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core/db"
)

type Config struct {
	DbConn        *sql.DB
	ImageUpEvents chan string
}

func NewConfig(clickhouseAddr string) *Config {
	return &Config{
		DbConn:        db.Connect(clickhouseAddr),
		ImageUpEvents: make(chan string, 10000),
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
