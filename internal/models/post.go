package models

import (
	"context"
	"database/sql"
	"time"

	"gitery/internal/prototypes"
)

// PostService implement prototypes.PostService interface
type PostService struct {
	DB *sql.DB
}

// Fetch single post
func (ps *PostService) Fetch(ctx context.Context, id int) (post prototypes.Post, err error) {
	post = prototypes.Post{}
	post.Comments = []prototypes.Comment{}
	err = ps.DB.QueryRowContext(ctx, "select id, title, content, user_id, created_at, updated_at from posts where id = $1", id).Scan(
		&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		err = NotFoundError(ctx, err)
		return
	}
	// query comments related to the post
	rows, err := ps.DB.QueryContext(ctx, "select id, content, user_id, created_at, updated_at from comments where post_id =$1", id)
	if err != nil {
		return
	}
	// Assemble comments with post structure
	for rows.Next() {
		comment := prototypes.Comment{PostID: &id}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	return
}

// FetchList is to get multiple posts most recently
func (ps *PostService) FetchList(ctx context.Context, limit int, offset int) (posts []prototypes.Post, err error) {
	if limit == 0 {
		limit = 10
	}
	posts = []prototypes.Post{}
	// query all the posts of the user
	postRows, err := ps.DB.QueryContext(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return
	}

	// fill the posts into list
	for postRows.Next() {
		post := prototypes.Post{Comments: []prototypes.Comment{}}
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Create a new post
func (ps *PostService) Create(ctx context.Context, post *prototypes.Post) (err error) {
	statement := "INSERT INTO posts (title, content, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at"
	stmt, err := ps.DB.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, post.Title, post.Content, post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	post.Comments = []prototypes.Comment{}
	return
}

// Update a post
func (ps *PostService) Update(ctx context.Context, post *prototypes.Post) (err error) {
	err = ps.DB.QueryRowContext(ctx, "UPDATE posts SET title = $3, content = $4, updated_at = $5 WHERE id = $1 AND user_id = $2 RETURNING updated_at",
		post.ID, post.UserID, post.Title, post.Content, time.Now()).Scan(&post.UpdatedAt)
	return
}

// Delete a post
func (ps *PostService) Delete(ctx context.Context, post *prototypes.Post) (err error) {
	_, err = ps.DB.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 AND user_id =$2", post.ID, post.UserID)
	return
}
