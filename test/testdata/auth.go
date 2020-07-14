package testdata

import (
	"context"
	"errors"
	"time"

	"gitery/internal/prototypes"
	"gitery/internal/tools/jwt"
)

// FakeAuthService ...
type FakeAuthService struct {
	Token     string
	JwtSecret string
}

// Login ...
func (as *FakeAuthService) Login(ctx context.Context, login *prototypes.Login) (token string, user *prototypes.User, err error) {
	token = as.Token
	ID := 1001
	user = &prototypes.User{
		ID: &ID,
	}
	return
}

// Verify ...
func (as *FakeAuthService) Verify(ctx context.Context, token string) (payload *prototypes.JwtPayload, err error) {
	payload = &prototypes.JwtPayload{}
	err = jwt.Decode(token, as.JwtSecret, payload)
	if err != nil {
		return
	}
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		err = errors.New("Invalid token: token expired")
		return
	}
	return
}
