package prototypes

import (
	"context"
	"time"
)

// Register ...
type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User ...
type User struct {
	ID        *int      `json:"id"`
	Email     string    `json:"email"`
	HashedPwd string    `json:"-"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"-"` // reconstruct in UserView
	UpdatedAt time.Time `json:"-"` // reconstruct in UserView
}

// UserService ...
type UserService interface {
	Fetch(ctx context.Context, id int) (user User, err error)
	Create(ctx context.Context, user *User) (err error)
	Update(ctx context.Context, user *User) (err error)
	Delete(ctx context.Context, login *Login) (err error)
}

// UserPostService ...
type UserPostService interface {
	Fetch(ctx context.Context, id int) (posts []Post, err error)
}
