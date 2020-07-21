package prototypes

import (
	"context"
	"time"
)

// Post ...
type Post struct {
	// Structural data
	ID        *int      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    *int      `json:"userID"`
	CreatedAt time.Time `json:"-"` // reconstruct in PostView
	UpdatedAt time.Time `json:"-"` // reconstruct in PostView
	IsDeleted bool      `json:"-"`
	// Linked data
	Author   *User      `json:"-"` // reconstruct in PostView
	Tags     []*Tag     `json:"tags,omitempty"`
	Comments []*Comment `json:"-"` // reconstruct in PostView
	Likes    int        `json:"likes"`
}

// PostService ...
type PostService interface {
	Fetch(ctx context.Context, id int) (post *Post, err error)
	FetchDetail(ctx context.Context, id int) (post *Post, err error)
	FetchList(ctx context.Context, limit int, offset int, userID int) (posts []*Post, err error)
	Create(ctx context.Context, post *Post) (err error)
	Update(ctx context.Context, post *Post) (err error)
	Delete(ctx context.Context, post *Post) (err error)
}

// PostLikeService ...
type PostLikeService interface {
	FetchLikes(ctx context.Context, userID int) (likes []*int, err error)
	Like(ctx context.Context, userID int, postID int) (err error)
	Unlike(ctx context.Context, userID int, postID int) (err error)
}
