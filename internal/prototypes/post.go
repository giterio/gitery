package prototypes

import (
	"context"
	"time"
)

// Post ...
type Post struct {
	ID        *int      `json:"id"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"-"` // reconstruct in PostView
	UserID    *int      `json:"user_id"`
	CreatedAt time.Time `json:"-"` // reconstruct in PostView
	UpdatedAt time.Time `json:"-"` // reconstruct in PostView
}

// PostService ...
type PostService interface {
	Fetch(ctx context.Context, id int) (post Post, err error)
	Create(ctx context.Context, post *Post) (err error)
	Update(ctx context.Context, post *Post) (err error)
	Delete(ctx context.Context, id int) (err error)
}
