package middleware

import (
	"github.com/hedge10/airmail/pkg/conf"
	"github.com/labstack/echo/v4"
)

func Config(config *conf.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", config)
			return next(c)
		}
	}
}
