package main

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/hedge10/airmail/pkg/api"
	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/middleware"
)



func main() {
	cfg, err := conf.New()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Validator = &api.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Config(cfg))
	e.Use(middleware.LogRequest())
	e.Use(middleware.EnforceContentType())

	// Register our routes
	api.RegisterHandler(e, cfg)

	if err = e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatal("Server failed")
	}
}
