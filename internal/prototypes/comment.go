package prototypes

import (
	"context"
	"time"
)

// Comment ...
type Comment struct {
	ID        *int      `json:"id"`
	Content   string    `json:"content"`
	PostID    *int      `json:"postId"`
	UserID    *int      `json:"userId"`
	CreatedAt time.Time `json:"-"` // reconstruct in CommentView
	UpdatedAt time.Time `json:"-"` // reconstruct in CommentView
}

// CommentService ...
type CommentService interface {
	Fetch(ctx context.Context, id int) (comment Comment, err error)
	Create(ctx context.Context, comment *Comment) (err error)
	Update(ctx context.Context, comment *Comment) (err error)
	Delete(ctx context.Context, comment *Comment) (err error)
}
