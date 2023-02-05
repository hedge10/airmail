package api

import (
	"github.com/labstack/echo/v4"

	"github.com/hedge10/airmail/pkg/conf"
)

type AirmailContext struct {
	echo.Context
}


func RegisterHandler(echo *echo.Echo, config *conf.Config) {
	echo.POST("/mail/send", sendMail)
}
