package middleware

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	CT_JSON = "application/json"
	CT_FORM = "application/x-www-form-urlencoded"
)

func EnforceContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		content_type := r.Header.Get("Content-Type")

		if content_type == "" {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
			return
		}

		if content_type != CT_JSON && content_type != CT_FORM {
			log.Info("Incoming request with unsupported Content-Type header")
			http.Error(w, fmt.Sprintf("Content-Type header must be '%s' or '%s'", CT_FORM, CT_JSON), http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
