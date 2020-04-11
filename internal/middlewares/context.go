package middlewares

import (
	"context"
	"net/http"
	"time"
)

// WrapContext put database and request in r.Context
func WrapContext(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, err := context.WithTimeout(req.Context(), 10*time.Second)
		if err != nil {
			return
		}
		h(w, req.WithContext(ctx))
	}
}
