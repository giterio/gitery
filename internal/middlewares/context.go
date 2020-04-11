package middlewares

import (
	"context"
	"net/http"
	"time"
)

// WrapContext put database and request in r.Context
func WrapContext(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
		defer cancel()
		h(w, req.WithContext(ctx))
	}
}
