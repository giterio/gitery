package models

import (
	"context"
	"database/sql"

	"gitery/internal/prototypes"
)

// UserService ...
type UserService struct {
	DB *sql.DB
}

// Fetch user information
func (us *UserService) Fetch(ctx context.Context, id int) (user prototypes.User, err error) {
	user = prototypes.User{}
	err = us.DB.QueryRowContext(ctx, "select id, email, hashed_pwd, created_at, updated_at from users where id = $1", id).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.CreatedAt, &user.UpdatedAt)
	return
}

// Create new user
func (us *UserService) Create(ctx context.Context, user *prototypes.User) (err error) {
	statement := "insert into users (email, hashed_pwd) values ($1, $2) returning id, created_at, updated_at"
	stmt, err := us.DB.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, user.Email, user.HashedPwd).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return
}

// Update a user
func (us *UserService) Update(ctx context.Context, user *prototypes.User) (err error) {
	err = us.DB.QueryRowContext(ctx, "update users set hashed_pwd = $2 where id = $1 returning updated_at",
		user.ID, user.HashedPwd).Scan(&user.UpdatedAt)
	return
}

// Delete a post
func (us *UserService) Delete(ctx context.Context, user *prototypes.User) (err error) {
	_, err = us.DB.ExecContext(ctx, "delete from users where email = $1 and hashed_pwd returning id", user.Email, user.HashedPwd)
	return
}
