package middlewares

import (
	"context"
	"net/http"
	"strings"

	"gitery/internal/controllers"
	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// Authentication ...
func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			ctx := r.Context()
			token := header[7:]
			h, ok := h.(*controllers.RootHandler)
			if !ok {
				err := models.ServerError(ctx, nil)
				views.RenderError(ctx, w, err)
			}
			// verify if token is valid
			payload, err := h.AuthHandler.Model.Verify(ctx, token)
			if err == nil {
				ctx = context.WithValue(ctx, prototypes.UserKey, *payload)
				r = r.WithContext(ctx)
			}
		}
		h.ServeHTTP(w, r)
	})
}
