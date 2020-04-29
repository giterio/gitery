package views

import (
	"context"
	"net/http"

	"gitery/internal/prototypes"
)

// TokenView ...
type TokenView struct {
	prototypes.Token
	ExpiredAt int64 `json:"expired_at"`
}

// BuildTokenView ...
func BuildTokenView(token prototypes.Token) TokenView {
	return TokenView{
		Token:     token,
		ExpiredAt: token.ExpiredAt.Unix(),
	}
}

// RenderToken ...
func RenderToken(ctx context.Context, w http.ResponseWriter, token prototypes.Token) (err error) {
	tokenView := BuildTokenView(token)
	err = Render(ctx, w, tokenView)
	return
}
