package prototypes

import (
	"context"
	"time"
)

// Comment ...
type Comment struct {
	ID        *int      `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	PostID    *int      `json:"post_id"`
	CreatedAt time.Time `json:"-"` // reconstruct in CommentView
	UpdatedAt time.Time `json:"-"` // reconstruct in CommentView
}

// CommentService ...
type CommentService interface {
	Fetch(ctx context.Context, id int) (comment Comment, err error)
	Create(ctx context.Context, comment *Comment) (err error)
	Update(ctx context.Context, comment *Comment) (err error)
	Delete(ctx context.Context, id int) (err error)
}
