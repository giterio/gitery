package views

import (
	"context"
	"net/http"

	"gitery/internal/prototypes"
)

// UserView ...
type UserView struct {
	prototypes.User
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

// BuildUserView compose UserView from a User
func BuildUserView(user prototypes.User) UserView {
	return UserView{
		User:      user,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

// RenderUser ...
func RenderUser(ctx context.Context, w http.ResponseWriter, user prototypes.User) (err error) {
	userView := BuildUserView(user)
	err = Render(ctx, w, userView)
	return
}
