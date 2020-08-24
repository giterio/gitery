package prototypes

import (
	"context"
	"time"
)

// Comment ...
type Comment struct {
	// Structural data
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

// CommentVote ...
type CommentVote struct {
	CommentID *int `json:"commentID"`
	UserID    *int `json:"userID"`
	Vote      bool `json:"vote"`
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
	FetchVotes(ctx context.Context, userID int, postID int) (votes []*CommentVote, err error)
	Vote(ctx context.Context, commentVote *CommentVote) (err error)
	Cancel(ctx context.Context, userID int, commentID int) (err error)
}
