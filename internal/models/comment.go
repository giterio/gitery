package models

import (
	"context"
	"database/sql"
	"errors"
)

// Comment ...
type Comment struct {
	DB      *sql.DB `json:"-"`
	ID      int     `json:"id"`
	Content string  `json:"content"`
	Author  string  `json:"author"`
	PostID  *int    `json:"post_id"`
}

// Fetch single comment
func (comment *Comment) Fetch(ctx context.Context, id int) (err error) {
	err = comment.DB.QueryRowContext(ctx, "select id, content, author from comments where id = $1", id).Scan(&comment.ID, &comment.Content, &comment.Author)
	return
}

// Create ...
func (comment *Comment) Create(ctx context.Context) (err error) {
	if comment.PostID == nil {
		err = errors.New("Post not found")
		return
	}
	err = comment.DB.QueryRowContext(ctx, "insert into comments (content, author, post_id) values ($1, $2, $3) returning id", comment.Content, comment.Author, comment.PostID).Scan(&comment.ID)
	return
}

// Update a comment
func (comment *Comment) Update(ctx context.Context) (err error) {
	_, err = comment.DB.ExecContext(ctx, "update comments set content = $2, author = $3 where id = $1", comment.ID, comment.Content, comment.Author)
	return
}

// Delete a comment
func (comment *Comment) Delete(ctx context.Context) (err error) {
	_, err = comment.DB.ExecContext(ctx, "delete from comments where id = $1", comment.ID)
	return
}
