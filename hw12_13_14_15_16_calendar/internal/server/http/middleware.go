package internalhttp

import (
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused, revive
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:revive
		// TODO
	})
}
