package views

import (
	"context"
	"gitery/internal/prototypes"
	"net/http"
)

// AuthView is the data structure for login response
type AuthView struct {
	Token string    `json:"token"`
	User  *UserView `json:"profile"`
}

// BuildAuthView compose the login response data structure
func BuildAuthView(token string, user *prototypes.User) (authView AuthView) {
	authView = AuthView{
		Token: token,
	}
	if user != nil {
		userView := BuildUserView(user)
		authView.User = &userView
	}
	return
}

// RenderAuth writes login response to http connection
func RenderAuth(ctx context.Context, w http.ResponseWriter, token string, user *prototypes.User) (err error) {
	tokenView := BuildAuthView(token, user)
	err = Render(ctx, w, tokenView)
	return
}
