package models

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"gitery/internal/prototypes"
)

// UserService ...
type UserService struct {
	DB *sql.DB
}

// Fetch user information
func (us *UserService) Fetch(ctx context.Context, id int) (user prototypes.User, err error) {
	user = prototypes.User{}
	err = us.DB.QueryRowContext(ctx, "SELECT id, email, hashed_pwd, nickname, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Create new user
func (us *UserService) Create(ctx context.Context, user *prototypes.User) (err error) {
	statement := "INSERT INTO users (email, hashed_pwd, nickname) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
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
	err = us.DB.QueryRowContext(ctx, "UPDATE users set hashed_pwd = $2, nickname = $3 WHERE id = $1 RETURNING updated_at",
		user.ID, user.HashedPwd, user.Nickname).Scan(&user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Delete a post
func (us *UserService) Delete(ctx context.Context, auth *prototypes.Auth) (err error) {
	user := prototypes.User{}
	err = us.DB.QueryRowContext(ctx, "SELECT id, email, hashed_pwd, nickname, created_at, updated_at FROM users WHERE email = $1", auth.Email).Scan(
		&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = IdentityNonExistError(ctx, err)
		} else {
			err = TransactionError(ctx, err)
		}
		return
	}
	// check if hash matched password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(auth.Password))
	if err != nil {
		err = InvalidPasswordError(ctx, err)
		return
	}
	_, err = us.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", user.ID)
	return
}

// UserPostService supply service about user's posts
type UserPostService struct {
	DB *sql.DB
}

// Fetch user's all posts
func (ups *UserPostService) Fetch(ctx context.Context, id int) (posts []prototypes.Post, err error) {
	// postMap is used to assemble posts and comments efficiently
	postMap := map[int]*prototypes.Post{}
	// query all the posts of the user
	postRows, err := ups.DB.QueryContext(ctx, "SELECT id, title, content, created_at, updated_at FROM posts WHERE user_id =$1", id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	// fill the posts into postMap using post ID as the key
	for postRows.Next() {
		post := prototypes.Post{UserID: &id, Comments: []prototypes.Comment{}}
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return
		}
		postMap[*post.ID] = &post
	}
	// query all the comments related to the posts
	commentRows, err := ups.DB.QueryContext(ctx, "SELECT id, content, post_id, created_at, updated_at FROM comments WHERE post_id IN (SELECT id FROM posts WHERE user_id = $1)", id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	// Assemble comments with post structure
	for commentRows.Next() {
		comment := prototypes.Comment{UserID: &id}
		err = commentRows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return
		}
		post := postMap[*comment.PostID]
		post.Comments = append(post.Comments, comment)
	}
	// convert postMap to post list
	posts = []prototypes.Post{}
	for _, post := range postMap {
		posts = append(posts, *post)
	}
	return
}
