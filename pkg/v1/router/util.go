package router

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func IntParam(ctx echo.Context, name string) int {
	param := ctx.QueryParam(name)

	ret := 0
	if param != "" {
		ret, _ = strconv.Atoi(param)
	}
	return ret
}

func OrderParam(ctx echo.Context, name string, defaultOrder string) string {
	param := ctx.QueryParam(name)

	if param == "" {
		return defaultOrder
	}
	return param
}
