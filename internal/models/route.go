package models

import (
	"context"
	"net/http"
	"path"
	"strings"

	"gitery/internal/prototypes"
)

// shiftPath is to get the first parameter from path and generate next sub-path
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// Route ...
type Route struct {
	Path string // remaining path to explore
}

// infestContext is to bind request's context with the route
func (route *Route) infestContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), prototypes.RouteKey, route)
	return r.WithContext(ctx)
}

// Shift is to get the first parameter from route path and generate next sub-route
func (route *Route) shift() (resource string, subRoute *Route) {
	resource, subPath := shiftPath(route.Path)
	subRoute = &Route{Path: subPath}
	return
}

// IsLast ...
func (route *Route) IsLast() bool {
	return route.Path == "/"
}

// CurrentRoute ...
func CurrentRoute(r *http.Request) (route *Route) {
	ctx := r.Context()
	route, ok := ctx.Value(prototypes.RouteKey).(*Route)
	if !ok {
		// create a Route with full path
		route = &Route{Path: r.URL.Path}
	}
	return
}

// ShiftRoute is to shift resource name from request and generate next sub-route
func ShiftRoute(r *http.Request) (resource string, rn *http.Request) {
	route := CurrentRoute(r)
	// extract the first parameter and generate a sub-route
	resource, subRoute := route.shift()
	// bind the sub-route with request's context
	rn = subRoute.infestContext(r)
	return
}
