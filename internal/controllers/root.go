package controllers

import (
	"net/http"

	"gitery/internal/models"
)

// Root is the root handler of comming request
type Root struct {
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
}

func (h *Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resource, r := models.ShiftRoute(r)
	switch resource {
	case "post":
		h.PostHandler.ServeHTTP(w, r)
		return
	case "comment":
		h.CommentHandler.ServeHTTP(w, r)
		return
	default:
		http.NotFound(w, r)
	}
}
