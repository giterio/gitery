package prototypes

import (
	"context"
)

// Login ...
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JwtUserPub ...
type JwtUserPub struct {
	ID    *int   `json:"userID"`
	Email string `json:"email"`
}

// JwtPayload ...
type JwtPayload struct {
	Iss string     `json:"iss,omitempty"` // issuer
	Sub string     `json:"sub,omitempty"` // subject
	Aud string     `json:"aud,omitempty"` // audience
	Iat int64      `json:"iat"`           // Issued At
	Nbf int64      `json:"nbf,omitempty"` // Not Before
	Exp int64      `json:"exp,omitempty"` // expiration time
	Jti int64      `json:"jti,omitempty"` // JWT ID
	Pub JwtUserPub `json:"pub"`           // Public info of user
}

// AuthService ...
type AuthService interface {
	Login(ctx context.Context, login Login) (token string, user User, err error)
	Verify(ctx context.Context, token string) (payload JwtPayload, err error)
}
