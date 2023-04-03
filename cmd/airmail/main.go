package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"

	"github.com/hedge10/airmail/pkg/api"
	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/middleware"
	"github.com/hedge10/airmail/pkg/storage"
)

func main() {
	cfg, err := conf.New()
	if err != nil {
		log.Fatal(err)
	}

	var s *storage.Storage
	if cfg.UseStorage {
		s, _ = storage.Connect(cfg)
	}

	middlewares := alice.New(middleware.EnforceContentType)

	mux := http.NewServeMux()
	mux.Handle("/", middlewares.Then(api.IncomingMessageHandler(cfg, s)))

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	serve_err := http.ListenAndServe(address, mux)
	if serve_err != nil {
		log.Fatal(serve_err)
	}
}
