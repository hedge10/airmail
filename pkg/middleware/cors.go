package middleware

import (
	"strings"

	"github.com/rs/cors"
)

func NewCorsConfig(origins string, headers string, debug bool) *cors.Cors {
	c := cors.New(cors.Options{
		AllowedHeaders: strings.Split(headers, ","),
		AllowedOrigins: strings.Split(origins, ","),
		Debug:          debug,
	})

	return c
}
