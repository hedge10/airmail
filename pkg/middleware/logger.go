package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func LogRequest() echo.MiddlewareFunc {
	log := logrus.New()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI: true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"host": values.Host,
				"method": values.Method,
				"uri": values.URI,
				"status": values.Status,
				"referrer": values.Referer,

			}).Info("Request")

			return nil
		},
	})
}
