package models

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gitery/internal/prototypes"
	"gitery/internal/tools/jwt"
)

// AuthService ...
type AuthService struct {
	DB        *sql.DB
	JwtSecret string
}

// Login ...
func (as *AuthService) Login(ctx context.Context, auth prototypes.Auth) (token string, err error) {
	user := prototypes.User{}
	err = as.DB.QueryRowContext(ctx, "select id, email, hashed_pwd, created_at, updated_at from users where email = $1", auth.Email).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = IdentityNonExistError(ctx, err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(auth.Password))
	if err != nil {
		err = InvalidPasswordError(ctx, err)
		return
	}
	payload := prototypes.JwtPayload{
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Pub: struct {
			ID    *int   `json:"user_id"`
			Email string `json:"email"`
		}{
			ID:    user.ID,
			Email: user.Email,
		},
	}
	token, err = jwt.Encode(payload, as.JwtSecret)
	return
}

// Verify ...
func (as *AuthService) Verify(ctx context.Context, token string) (payload prototypes.JwtPayload, err error) {
	payload = prototypes.JwtPayload{}
	err = jwt.Decode(token, as.JwtSecret, &payload)
	if err != nil {
		return
	}
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		err = AuthorizationError(ctx, err)
		return
	}
	return
}
