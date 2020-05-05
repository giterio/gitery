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
	// retrieve user matched given email
	err = as.DB.QueryRowContext(ctx, "select id, email, hashed_pwd, created_at, updated_at from users where email = $1", auth.Email).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = IdentityNonExistError(ctx, err)
		return
	}
	// check if password match the hash
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(auth.Password))
	if err != nil {
		err = InvalidPasswordError(ctx, err)
		return
	}
	userPub := prototypes.JwtUserPub{
		ID:    user.ID,
		Email: user.Email,
	}
	payload := prototypes.JwtPayload{
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Pub: userPub,
	}
	// generate JWT token for user
	token, err = jwt.Encode(payload, as.JwtSecret)
	return
}

// Verify if token is valid
func (as *AuthService) Verify(ctx context.Context, token string) (payload prototypes.JwtPayload, err error) {
	payload = prototypes.JwtPayload{}
	// retrieve public payload data from token
	err = jwt.Decode(token, as.JwtSecret, &payload)
	if err != nil {
		return
	}
	// check if token expired
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		err = AuthorizationError(ctx, err)
		return
	}
	return
}
