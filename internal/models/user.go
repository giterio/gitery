package models

import (
	"context"
	"database/sql"
	"gitery/internal/prototypes"
	"time"
)

// UserService ...
type UserService struct {
	DB *sql.DB
}

// Fetch user information
func (us *UserService) Fetch(ctx context.Context, token string) (user prototypes.User, err error) {
	user = prototypes.User{}
	var tokenCreateAt time.Time
	err = us.DB.QueryRowContext(ctx, "select user_id, created_at from auth where token = $1", token).Scan(
		&user.ID, &tokenCreateAt)
	// token has been expired
	if time.Now().After(tokenCreateAt.Add(time.Hour * 24 * 30)) {
		err = AuthorizationError(ctx)
	}
	err = us.DB.QueryRowContext(ctx, "select email, created_at, updated_at from users where id = $1", user.ID).Scan(
		&user.Email, &user.CreatedAt, &user.UpdatedAt)
	return
}

// Create new user
func (us *UserService) Create(ctx context.Context, user *prototypes.User) (err error) {
	statement := "insert into users (email, password) values ($1, $2) returning id, created_at, updated_at"
	stmt, err := us.DB.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return
}
