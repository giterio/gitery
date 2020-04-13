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

// Shift ...
func (route *Route) Shift() (string, *Route) {
	param, subPath := shiftPath(route.Path)
	if subPath == "/" {
		return param, nil
	}
	return param, &Route{Path: subPath}
}

// BindContext ...
func (route *Route) BindContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), routeKey, route)
	return r.WithContext(ctx)
}

// Extract ...
func Extract(ctx context.Context) (string, *Route, error) {
	if rv := ctx.Value(routeKey); rv != nil {
		param, subRoute := rv.(*Route).Shift()
		return param, subRoute, nil
	}
	return "", nil, errors.New("Not Exist")
}
