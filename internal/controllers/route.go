package controllers

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
	Path string
}

// BindContext is to bind the route with request's context
func (route *Route) BindContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), routeKey, route)
	return r.WithContext(ctx)
}

// Shift is to get the first parameter from route path and generate next sub-route
func (route *Route) Shift() (string, *Route) {
	param, subPath := shiftPath(route.Path)
	if subPath == "/" {
		return param, nil
	}
	return param, &Route{Path: subPath}
}

// ExtractRoute is to extract path parameter from context and generate next sub-route
func ExtractRoute(ctx context.Context) (string, *Route) {
	rv := ctx.Value(routeKey)
	if route, ok := rv.(*Route); ok {
		param, subRoute := route.Shift()
		return param, subRoute
	}
	return "", nil
}
