package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"net/http"
)

func ImageRouter(group *echo.Group) {
	group.GET("", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "name")

		ret, err := db.ListAllImages(conf, page, order)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})

	group.GET("/:p1", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		id := ctx.Param("p1")

		ret, err := db.FindImageById(conf, id)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})

	group.GET("/:p1/artifact", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		id := ctx.Param("p1")
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "name")

		ret, err := db.FindArtifactByImageId(conf, id, page, order)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})

	group.GET("/:p1/vulnerability", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		id := ctx.Param("p1")
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "name")

		ret, err := db.FindVulnerabilityByImageId(conf, id, page, order)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})

	group.POST("/status", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		data := new(struct {
			Image  string `json:"image" validate:"required"`
			Status *int64 `json:"status" validate:"required"`
		})
		if err := ctx.Bind(data); err != nil {
			return err
		}

		if err := ctx.Validate(data); err != nil {
			return err
		}

		conf.ImageEvents <- core.ImageEvent{
			Type:  core.ImageEventType(*data.Status),
			Image: data.Image,
		}
		return ctx.NoContent(http.StatusOK)
	})
}
