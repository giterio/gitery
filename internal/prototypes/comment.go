package prototypes

import (
	"context"
	"time"
)

// Comment ...
type Comment struct {
	ID        *int      `json:"id"`
	Content   string    `json:"content"`
	PostID    *int      `json:"postID"`
	UserID    *int      `json:"userID"`
	ParentID  *int      `json:"parentID,omitempty"`
	CreatedAt time.Time `json:"-"` // reconstruct in CommentView
	UpdatedAt time.Time `json:"-"` // reconstruct in CommentView
	IsDeleted bool      `json:"isDeleted"`
	// Linked data
	Author   *User      `json:"-"` // reconstruct in PostView
	Comments []*Comment `json:"-"` // reconstruct in PostView
	VoteUp   int        `json:"voteUp"`
	VoteDown int        `json:"voteDown"`
}

// CommentService ...
type CommentService interface {
	Fetch(ctx context.Context, id int) (comment *Comment, err error)
	FetchDetail(ctx context.Context, id int) (comment *Comment, err error)
	Create(ctx context.Context, comment *Comment) (err error)
	Update(ctx context.Context, comment *Comment) (err error)
	Delete(ctx context.Context, comment *Comment) (err error)
}

// CommentVoteService ...
type CommentVoteService interface {
	Vote(ctx context.Context, userID int, commentID int, vote bool) (err error)
	Cancel(ctx context.Context, userID int, commentID int) (err error)
}
