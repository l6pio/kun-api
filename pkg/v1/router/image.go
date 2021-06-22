package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/v1/router/vo/img"
	"net/http"
)

func ImageRouter(group *echo.Group) {
	group.POST("/up", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		key := new(img.Key)
		if err := ctx.Bind(key); err != nil {
			return err
		}

		if err := ctx.Validate(key); err != nil {
			return err
		}

		conf.ImageUpEvents <- *key
		return ctx.NoContent(http.StatusOK)
	})
}
