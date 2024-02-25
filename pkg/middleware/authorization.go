package middleware

import (
	"net/http"
	"strings"
)

type authorization struct {
	Token string
}

func (auth *authorization) Validate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth_header := r.Header.Get("Authorization")
		tokenSplit := strings.Split(auth_header, "Bearer ")

		if len(tokenSplit) != 2 || len(tokenSplit[1]) == 0 {
			http.Error(w, "Invalid token length", http.StatusBadRequest)
			return
		}

		if tokenSplit[1] != auth.Token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NewToken(token string) *authorization {
	return &authorization{
		Token: token,
	}
}
