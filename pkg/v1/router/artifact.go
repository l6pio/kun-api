package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"net/http"
)

func ArtifactRouter(group *echo.Group) {
	group.GET("/:p1/vulnerability", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)
		id := ctx.Param("p1")
		page := IntParam(ctx, "page")
		order := OrderParam(ctx, "order", "name")

		ret, err := db.FindVulnerabilityByArtifactId(conf, id, page, order)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})
}
