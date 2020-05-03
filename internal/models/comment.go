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
	err = cs.DB.QueryRowContext(ctx, "select id, content, user_id, post_id, created_at, updated_at from comments where id = $1", id).Scan(
		&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Create ...
func (cs *CommentService) Create(ctx context.Context, comment *prototypes.Comment) (err error) {
	if comment.PostID == nil {
		err = errors.New("Post not found")
		return
	}
	err = cs.DB.QueryRowContext(ctx, "insert into comments (content, user_id, post_id) values ($1, $2, $3) returning id, created_at, updated_at",
		comment.Content, comment.UserID, comment.PostID).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	return
}

// Update a comment
func (cs *CommentService) Update(ctx context.Context, comment *prototypes.Comment) (err error) {
	err = cs.DB.QueryRowContext(ctx, "update comments set content = $1, updated_at = $2 where id = $3 and user_id = $4 returning updated_at",
		comment.Content, time.Now(), comment.ID, comment.UserID).Scan(&comment.UpdatedAt)
	return
}

// Delete a comment
func (cs *CommentService) Delete(ctx context.Context, comment *prototypes.Comment) (err error) {
	_, err = cs.DB.ExecContext(ctx, "delete from comments where id = $1 and user_id = $2", comment.ID, comment.UserID)
	return
}
