package views

import (
	"context"
	"net/http"
)

// RenderEmpty writes a empty response with status code 200 to http connection
func RenderEmpty(ctx context.Context, w http.ResponseWriter) (err error) {
	return Render(ctx, w, nil)
}
