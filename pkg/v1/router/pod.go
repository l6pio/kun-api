package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/service"
	"net/http"
)

func PodRouter(group *echo.Group) {
	group.GET("/overview", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		ret, err := service.GetPodsOverview(conf)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})

	group.GET("/timeline", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		ret, err := db.FindRunningPodTimeline(conf)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, ret)
	})
}
