package prototypes

import "context"

// User ...
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserService ...
type UserService interface {
	Login(ctx context.Context, user *User) (string token, err error)
	Create(ctx context.Context, user *User) (string token, err error)
	Update(ctx context.Context, user *User) (err error)
	Delete(ctx context.Context, string token) (err error)
}
