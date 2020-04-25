package controllers

import (
	"net/http"
)

// Router is the root handler of comming request
type Router struct {
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// create a Route with full path
	route := Route{Path: r.URL.Path}
	// extract the first parameter and generate a sub-route
	resource, subRoute := route.Shift()
	// bind the sub-route with request's context
	r = subRoute.BindContext(r)

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
