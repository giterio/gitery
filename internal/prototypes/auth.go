package prototypes

import (
	"context"
)

// Auth ...
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserPub ...
type UserPub struct {
	ID    *int   `json:"user_id"`
	Email string `json:"email"`
}

// AuthService ...
type AuthService interface {
	Login(ctx context.Context, auth Auth) (token string, err error)
	Verify(ctx context.Context, token string) (userPub UserPub, err error)
	// getUserPub(ctx context.Context, r http.Request) (userPub UserPub, err error)
}
