package middlewares

import (
	"context"
	"gitery/internal/controllers"
	"gitery/internal/prototypes"
	"net/http"
	"strings"
)

// Authentication ...
func Authentication(h *controllers.RootHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			ctx := r.Context()
			token := header[7:]
			userPub, err := h.AuthHandler.Model.Verify(ctx, token)
			if err == nil {
				r.WithContext(context.WithValue(ctx, prototypes.UserKey, userPub))
			}
		}
		h.ServeHTTP(w, r)
	})
}
