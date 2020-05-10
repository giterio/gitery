package views

import (
	"context"
	"net/http"
)

// AuthView is the data structure for login response
type AuthView struct {
	Token string `json:"token"`
}

// BuildTokenView compose the login response data structure
func BuildTokenView(token string) AuthView {
	return AuthView{
		Token: token,
	}
}

// RenderAuth writes login response to http connection
func RenderAuth(ctx context.Context, w http.ResponseWriter, token string) (err error) {
	tokenView := BuildTokenView(token)
	err = Render(ctx, w, tokenView)
	return
}
