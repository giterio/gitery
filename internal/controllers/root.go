package controllers

import (
	"gitery/internal/utils"
	"net/http"
)

// RootHandler ...
type RootHandler struct {
	PostHandler *PostHandler
	// add other sub-handlers
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)
	switch head {
	case "post":
		h.PostHandler.ServeHTTP(w, r)
		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
