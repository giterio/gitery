package route

import (
	"context"
	"errors"
	"net/http"
	"path"
	"strings"
)

type contextKey int

const (
	routeKey contextKey = iota
	paramKey
)

// ShiftPath ...
func ShiftPath(p string) (head, tail string) {
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

// BasePath ...
func (route *Route) BasePath() string {
	return path.Base(route.Path)
}

// Next ...
func (route *Route) Next() (string, *Route) {
	currentPath, subPath := ShiftPath(route.Path)
	return currentPath, &Route{Path: subPath}
}

// InjectContext ...
func (route *Route) InjectContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), routeKey, route)
	return r.WithContext(ctx)
}

// Of ...
func Of(r *http.Request) (*Route, error) {
	if rv := r.Context().Value(routeKey); rv != nil {
		return rv.(*Route), nil
	}
	return nil, errors.New("Not Exist")
}
