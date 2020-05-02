package models

import (
	"context"
	"database/sql"
	"encoding/json"
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
func (us *AuthService) Login(ctx context.Context, auth prototypes.Auth) (token string, err error) {
	user := prototypes.User{}
	err = us.DB.QueryRowContext(ctx, "select id, email, hashed_pwd, created_at, updated_at from users where email = $1", auth.Email).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(auth.Password))
	if err != nil {
		return
	}
	userPub := prototypes.UserPub{
		ID:    user.ID,
		Email: user.Email,
	}
	payload := jwt.Payload{
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Pub: userPub,
	}
	token, err = jwt.Encode(payload, us.JwtSecret)
	return
}

// Verify ...
func (us *AuthService) Verify(ctx context.Context, token string) (userPub prototypes.UserPub, err error) {
	payload, err := jwt.Decode(token, us.JwtSecret)
	userPubBytes, err := json.Marshal(payload.Pub)
	if err != nil {
		return
	}
	err = json.Unmarshal(userPubBytes, &userPub)
	return
}
