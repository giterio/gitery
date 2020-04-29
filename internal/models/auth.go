package models

import (
	"context"
	"database/sql"

	"gitery/internal/prototypes"
)

// AuthService ...
type AuthService struct {
	DB *sql.DB
}

// Login ...
func (us *AuthService) Login(ctx context.Context, auth *prototypes.Auth) (err error) {
	// TODO implement login function
	// token has been expired
	// if time.Now().After(tokenCreateAt.Add(time.Hour * 24 * 30)) {
	// 	err = AuthorizationError(ctx)
	// }
	return
}

func (us *AuthService) Verify(ctx context.Context token string) (err error) {
	// TODO implement verify function
	return
}

// Login ...
func (us *AuthService) logout(ctx context.Context, token string) (err error) {
	// TODO implement logout function
	return
}
