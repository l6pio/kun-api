package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/img"
	"l6p.io/kun/api/pkg/v1/router/vo"
	"net/http"
)

func ImageRouter(group *echo.Group) {
	group.GET("", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "name")

		ret, err := img.List(conf, page, order)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, vo.Response{Result: ret})
	})

	group.POST("/up", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		key := new(struct {
			Image string `json:"image" validate:"required"`
		})
		if err := ctx.Bind(key); err != nil {
			return err
		}

		if err := ctx.Validate(key); err != nil {
			return err
		}

		conf.ImageUpEvents <- key.Image
		return ctx.NoContent(http.StatusOK)
	})
}
