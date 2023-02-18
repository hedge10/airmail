package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const (
	CT_JSON = "application/json"
	CT_FORM = "application/x-www-form-urlencoded"
)

func EnforceContentType() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			content_type := c.Request().Header.Get("Content-Type")

			if content_type == "" {
				log.Info("Incomig request with malformed header")
				return echo.NewHTTPError(http.StatusBadRequest, "Malformed Content-Type header")
			}

			if content_type != CT_JSON && content_type != CT_FORM {
				log.Info("Incomig request with unsupported header")
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Wrong Content-Type header supplied")
			}

			return next(c)
		}
	}
}
