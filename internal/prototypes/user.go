package prototypes

import (
	"context"
	"time"
)

// User ...
type User struct {
	ID        *int      `json:"id"`
	Email     string    `json:"email"`
	HashedPwd string    `json:"hashed_pwd,omitempty"`
	CreatedAt time.Time `json:"-"` // reconstruct in UserView
	UpdatedAt time.Time `json:"-"` // reconstruct in UserView
}

// UserService ...
type UserService interface {
	Fetch(ctx context.Context, id int) (user User, err error)
	Create(ctx context.Context, user *User) (token string, err error)
	Update(ctx context.Context, user *User) (err error)
	Delete(ctx context.Context, token string) (err error)
}
