package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
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

	auth := middleware.NewToken(cfg.AuthToken)

	middlewares := alice.New(auth.Validate, middleware.EnforceContentType)

	mux := http.NewServeMux()
	mux.Handle("/", middlewares.Then(api.IncomingMessageHandler(cfg)))

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	serve_err := http.ListenAndServe(address, mux)
	if serve_err != nil {
		log.Fatal(serve_err)
	}
}
