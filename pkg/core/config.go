package core

import (
	"github.com/labstack/echo/v4"
)

type Config struct {
	DatabaseName string
	MongoAddr    string
	MongoUser    string
	MongoPass    string
	ImageEvents  chan ImageEvent
}

type ImageStatus int

const (
	ImageUp   ImageStatus = 1
	ImageDown ImageStatus = 0
)

type ImageEvent struct {
	Timestamp int64        `json:"timestamp" validate:"required"`
	Image     string       `json:"image" validate:"required"`
	Status    *ImageStatus `json:"status" validate:"required"`
}

func NewConfig(mongoAddr, mongoUser, mongoPass string) *Config {
	return &Config{
		DatabaseName: "kun",
		MongoAddr:    mongoAddr,
		MongoUser:    mongoUser,
		MongoPass:    mongoPass,
		ImageEvents:  make(chan ImageEvent, 10000),
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
