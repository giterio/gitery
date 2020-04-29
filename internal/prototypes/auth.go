package prototypes

import (
	"context"
	"time"
)

// Auth ...
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token ...
type Token struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"-"`
}

// AuthService ...
type AuthService interface {
	Login(ctx context.Context, auth *Auth) (token string, err error)
	Verify(ctx context.Context, token string) (err error)
	Logout(ctx context.Context, token string) (err error)
}
