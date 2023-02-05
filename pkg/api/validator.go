package api

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type(
	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)

	if err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
