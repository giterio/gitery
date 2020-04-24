package models

import (
	"context"
	"database/sql"
	"errors"
	"gitery/internal/domains"
)

// CommentService ...
type CommentService struct {
	DB *sql.DB
}

// Fetch single comment
func (cs *CommentService) Fetch(ctx context.Context, id int) (comment domains.Comment, err error) {
	comment = domains.Comment{}
	err = cs.DB.QueryRowContext(ctx, "select id, content, author, post_id, created_at, updated_at from comments where id = $1", id).Scan(
		&comment.ID, &comment.Content, &comment.Author, &comment.PostID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Create ...
func (cs *CommentService) Create(ctx context.Context, comment *domains.Comment) (err error) {
	if comment.PostID == nil {
		err = errors.New("Post not found")
		return
	}
	err = cs.DB.QueryRowContext(ctx, "insert into comments (content, author, post_id) values ($1, $2, $3) returning id, created_at, updated_at",
		comment.Content, comment.Author, comment.PostID).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Update a comment
func (cs *CommentService) Update(ctx context.Context, comment *domains.Comment) (err error) {
	err = cs.DB.QueryRowContext(ctx, "update comments set content = $2, author = $3, updated_at = $4 where id = $1 returning updated_at",
		comment.ID, comment.Content, comment.Author).Scan(&comment.UpdatedAt)
	return
}

// Delete a comment
func (cs *CommentService) Delete(ctx context.Context, id int) (err error) {
	_, err = cs.DB.ExecContext(ctx, "delete from comments where id = $1", id)
	return
}
