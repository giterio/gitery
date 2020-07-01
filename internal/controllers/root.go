package controllers

import (
	"net/http"

	"gitery/internal/models"
	"gitery/internal/views"
)

// RootHandler is the root handler of comming request
type RootHandler struct {
	AuthHandler    *AuthHandler
	UserHandler    *UserHandler
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
	TagHandler     *TagHandler
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resource, r := models.ShiftRoute(r)
	switch resource {
	case "auth":
		h.AuthHandler.ServeHTTP(w, r)
		return
	case "user":
		h.UserHandler.ServeHTTP(w, r)
		return
	case "post":
		h.PostHandler.ServeHTTP(w, r)
		return
	case "comment":
		h.CommentHandler.ServeHTTP(w, r)
		return
	case "tag":
		h.TagHandler.ServeHTTP(w, r)
		return
	default:
		ctx := r.Context()
		err := models.NotFoundError(ctx, nil)
		views.RenderError(ctx, w, err)
	}
}
