package router

import (
	"github.com/labstack/echo/v4"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/service"
	"net/http"
)

func GarbageRouter(group *echo.Group) {
	group.GET("/namespace", func(ctx echo.Context) error {
		conf := ctx.Get("config").(*core.Config)

		err := service.FindNamespaceGarbage(conf)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, "")
	})
}
