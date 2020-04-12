package controllers

import (
	"net/http"

	route "gitery/internal/utils"
)

// Router ...
type Router struct {
	PostHandler *PostHandler
	// add other sub-handlers
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := route.Route{Path: r.URL.Path}
	currentPath, subRoute := route.Next()
	r = subRoute.InjectContext(r)
	switch currentPath {
	case "post":
		h.PostHandler.ServeHTTP(w, r)
		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
