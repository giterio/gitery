package prototype

import (
	"context"
	"time"
)

// Post ...
type Post struct {
	ID        *int      `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostService ...
type PostService interface {
	Fetch(ctx context.Context, id int) (post Post, err error)
	Create(ctx context.Context, post *Post) (err error)
	Update(ctx context.Context, post *Post) (err error)
	Delete(ctx context.Context, id int) (err error)
}
