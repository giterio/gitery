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
	return
}

// Login ...
func (us *AuthService) logout(ctx context.Context, token string) (err error) {
	// TODO implement logout function
	return
}
