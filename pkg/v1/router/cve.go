package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/cve"
	"l6p.io/kun/api/pkg/core/es"
	"l6p.io/kun/api/pkg/v1/router/vo/scan"
	"l6p.io/kun/api/pkg/v1/router/vo/search"
	"net/http"
)

func CveRouter(group *echo.Group) {
	group.POST("/scan", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		data := new(scan.Key)
		if err := ctx.Bind(data); err != nil {
			return err
		}

		if err := ctx.Validate(data); err != nil {
			return err
		}

		id, err := cve.Scan(conf, data.ImageRepo, data.ImageTag)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, scan.Response{ID: id})
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

		reports, err := es.SearchByImageID(conf, data.ImageID)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, search.Response{
			Request: data,
			Reports: reports,
		})
	})
}
