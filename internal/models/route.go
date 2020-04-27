package models

import (
	"context"
	"net/http"
	"path"
	"strings"
)

type contextKey int

const (
	routeKey contextKey = iota
	paramKey
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
	ctx := context.WithValue(r.Context(), routeKey, route)
	return r.WithContext(ctx)
}

// Shift is to get the first parameter from route path and generate next sub-route
func (route *Route) shift() (resource string, subRoute *Route) {
	resource, subPath := shiftPath(route.Path)
	if subPath == "/" {
		return
	}
	subRoute = &Route{Path: subPath}
	return
}

// ShiftRoute is to shift resource name from request and generate next sub-route
func ShiftRoute(r *http.Request) (resource string, rn *http.Request) {
	ctx := r.Context()
	rv := ctx.Value(routeKey)
	route, ok := rv.(*Route)
	if !ok {
		// create a Route with full path
		route = &Route{Path: r.URL.Path}
	}
	// extract the first parameter and generate a sub-route
	resource, subRoute := route.shift()
	// bind the sub-route with request's context
	rn = subRoute.infestContext(r)
	return
}
