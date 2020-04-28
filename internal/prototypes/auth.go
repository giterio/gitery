package prototypes

import (
	"context"
)

// Auth ...
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthService ...
type AuthService interface {
	Login(ctx context.Context, auth *Auth) (token string, err error)
	Logout(ctx context.Context, token string) (err error)
}
