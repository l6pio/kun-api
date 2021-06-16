package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve"
	"l6p.io/kun/api/pkg/v1/router/vo/scan"
	"l6p.io/kun/api/pkg/v1/router/vo/search"
	"net/http"
)

func CveRouter(group *echo.Group) {
	group.GET("", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		reports, err := cve.ListAll(conf)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, search.Response{
			Reports: reports,
		})
	})

	group.POST("/scan", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		key := new(scan.Key)
		if err := ctx.Bind(key); err != nil {
			return err
		}

		if err := ctx.Validate(key); err != nil {
			return err
		}

		conf.ScanRequests <- *key
		return ctx.NoContent(http.StatusOK)
	})

	group.POST("/search/by-image-id", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		data := new(search.ByImageID)
		if err := ctx.Bind(data); err != nil {
			return err
		}

		if err := ctx.Validate(data); err != nil {
			return err
		}

		reports, err := cve.FindByImageID(conf, data.ImageID)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, search.Response{
			Request: data,
			Reports: reports,
		})
	})
}
