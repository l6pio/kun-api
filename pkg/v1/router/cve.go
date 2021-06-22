package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve"
	"l6p.io/kun/api/pkg/v1/router/vo"
	"net/http"
)

func CveRouter(group *echo.Group) {
	group.GET("", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "severity")

		ret, err := cve.List(conf, page, order)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, vo.Response{Result: ret})
	})
}
