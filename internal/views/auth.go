package views

import (
	"context"
	"net/http"
)

// AuthView ...
type AuthView struct {
	Token string `json:"token"`
}

// BuildTokenView ...
func BuildTokenView(token string) AuthView {
	return AuthView{
		Token: token,
	}
}

// RenderAuth ...
func RenderAuth(ctx context.Context, w http.ResponseWriter, token string) (err error) {
	tokenView := BuildTokenView(token)
	err = Render(ctx, w, tokenView)
	return
}
