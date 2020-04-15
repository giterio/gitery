package domains

import (
	"context"
)

// Post ...
type Post struct {
	ID       *int      `json:"id"`
	Content  string    `json:"content"`
	Author   string    `json:"author"`
	Comments []Comment `json:"comments"`
}

// PostService ...
type PostService interface {
	Fetch(ctx context.Context, id int) (err error)
	Create(ctx context.Context) (err error)
	Update(ctx context.Context) (err error)
	Delete(ctx context.Context) (err error)
}
