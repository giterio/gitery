package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gitery/internal/prototypes"
)

// CommentService implement the prototypes.CommentService interface
type CommentService struct {
	DB *sql.DB
}

// Fetch single comment
func (cs *CommentService) Fetch(ctx context.Context, id int) (comment prototypes.Comment, err error) {
	comment = prototypes.Comment{}
	err = cs.DB.QueryRowContext(ctx, "SELECT id, content, user_id, post_id, created_at, updated_at FROM comments WHERE id = $1", id).Scan(
		&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Create comment
func (cs *CommentService) Create(ctx context.Context, comment *prototypes.Comment) (err error) {
	if comment.PostID == nil {
		err = errors.New("Post not found")
		return
	}
	err = cs.DB.QueryRowContext(ctx, "INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		comment.Content, comment.UserID, comment.PostID).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Update a comment
func (cs *CommentService) Update(ctx context.Context, comment *prototypes.Comment) (err error) {
	err = cs.DB.QueryRowContext(ctx, "UPDATE comments SET content = $3, updated_at = $4 WHERE id = $1 AND user_id = $2 RETURNING updated_at",
		comment.ID, comment.UserID, comment.Content, time.Now()).Scan(&comment.UpdatedAt)
	return
}

// Delete a comment
func (cs *CommentService) Delete(ctx context.Context, comment *prototypes.Comment) (err error) {
	_, err = cs.DB.ExecContext(ctx, "DELETE FROM comments WHERE id = $1 AND user_id = $2", comment.ID, comment.UserID)
	return
}
