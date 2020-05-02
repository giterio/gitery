package prototypes

import (
	"context"
)

// Auth ...
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JwtPayload ...
type JwtPayload struct {
	Iss string `json:"iss,omitempty"` // issuer
	Exp int64  `json:"exp,omitempty"` // expiration time
	Sub string `json:"sub,omitempty"` // subject
	Aud string `json:"aud,omitempty"` // audience
	Nbf int64  `json:"nbf,omitempty"` // Not Before
	Iat int64  `json:"iat"`           // Issued At
	Jti int64  `json:"jti,omitempty"` // JWT ID
	Pub struct {
		ID    *int   `json:"user_id"`
		Email string `json:"email"`
	} `json:"pub"`
}

// AuthService ...
type AuthService interface {
	Login(ctx context.Context, auth Auth) (token string, err error)
	Verify(ctx context.Context, token string) (payload JwtPayload, err error)
	// getUserPub(ctx context.Context, r http.Request) (userPub UserPub, err error)
}
