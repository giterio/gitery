package middlewares

import (
	"context"
	"net/http"
	"time"
)

// LoadContext put database and request in r.Context
func LoadContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
