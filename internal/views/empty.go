package views

import (
	"context"
	"net/http"
)

// RenderEmpty ...
func RenderEmpty(ctx context.Context, w http.ResponseWriter) (err error) {
	return Render(ctx, w, nil)
}
