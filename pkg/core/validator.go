package core

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func AddValidator(server *echo.Echo) {
	server.Validator = &CustomValidator{validator: validator.New()}
}
