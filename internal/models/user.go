package models

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gitery/internal/prototypes"
)

// UserService ...
type UserService struct {
	DB *sql.DB
}

// Fetch user information
func (us *UserService) Fetch(ctx context.Context, id int) (user *prototypes.User, err error) {
	userCh := make(chan *prototypes.User)
	likesCh := make(chan []*int)
	errCh := make(chan error)

	go func() {
		user = &prototypes.User{}
		err = us.DB.QueryRowContext(ctx, `
		SELECT id, email, hashed_pwd, nickname, created_at, updated_at
		FROM users
		WHERE id = $1
		`, id).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			errCh <- HandleDatabaseQueryError(ctx, err)
			return
		}
		userCh <- user
	}()

	go func() {
		likeRows, err := us.DB.QueryContext(ctx, `
		SELECT post_id
		FROM post_like
		WHERE user_id = $1
		`, id)
		if err != nil {
			errCh <- TransactionError(ctx, err)
			return
		}
		defer likeRows.Close()

		likes := []*int{}
		for likeRows.Next() {
			var postID int
			err = likeRows.Scan(&postID)
			if err != nil {
				errCh <- err
				return
			}
			likes = append(likes, &postID)
		}
		likesCh <- likes
	}()

	var likes []*int
	for i := 0; i < 2; i++ {
		select {
		case err = <-errCh:
			return
		case user = <-userCh:
		case likes = <-likesCh:
		}
	}
	user.Likes = likes
	return
}

// Create new user
func (us *UserService) Create(ctx context.Context, user *prototypes.User) (err error) {
	statement := `
		INSERT INTO users (email, hashed_pwd, nickname)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	stmt, err := us.DB.PrepareContext(ctx, statement)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, user.Email, user.HashedPwd, user.Nickname).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Update a user
func (us *UserService) Update(ctx context.Context, user *prototypes.User) (err error) {
	err = us.DB.QueryRowContext(ctx, `
		UPDATE users
		set hashed_pwd = $2, nickname = $3, updated_at = $4
		WHERE id = $1
		RETURNING updated_at
		`, user.ID, user.HashedPwd, user.Nickname, time.Now()).Scan(&user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Delete a post
func (us *UserService) Delete(ctx context.Context, login *prototypes.Login) (err error) {
	user := prototypes.User{}
	err = us.DB.QueryRowContext(ctx, `
		SELECT id, email, hashed_pwd, nickname, created_at, updated_at
		FROM users
		WHERE email = $1
		`, login.Email).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = IdentityNonExistError(ctx, err)
		} else {
			err = TransactionError(ctx, err)
		}
		return
	}
	// check if hash matched password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(login.Password))
	if err != nil {
		err = InvalidPasswordError(ctx, err)
		return
	}
	_, err = us.DB.ExecContext(ctx, `
		UPDATE users
		set is_deleted = $2, updated_at = $4
		WHERE id = $1
		`, user.ID, true, time.Now())
	return
}
