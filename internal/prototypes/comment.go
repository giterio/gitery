package prototypes

import (
	"context"
	"time"
)

// Comment ...
type Comment struct {
	ID        *int       `json:"id"`
	Content   string     `json:"content"`
	PostID    *int       `json:"postID"`
	UserID    *int       `json:"userID"`
	ParentID  *int       `json:"parentID"`
	CreatedAt time.Time  `json:"-"` // reconstruct in CommentView
	UpdatedAt time.Time  `json:"-"` // reconstruct in CommentView
	IsDeleted bool       `json:"isDeleted"`
	Author    *User      `json:"-"` // reconstruct in PostView
	Comments  []*Comment `json:"-"` // reconstruct in PostView
}

// CommentService ...
type CommentService interface {
	Fetch(ctx context.Context, id int) (comment Comment, err error)
	FetchDetail(ctx context.Context, id int) (comment Comment, err error)
	Create(ctx context.Context, comment *Comment) (err error)
	Update(ctx context.Context, comment *Comment) (err error)
	Delete(ctx context.Context, comment *Comment) (err error)
}
