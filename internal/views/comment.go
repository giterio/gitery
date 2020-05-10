package views

import (
	"context"
	"gitery/internal/prototypes"
	"net/http"
)

// CommentView ...
type CommentView struct {
	prototypes.Comment
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// BuildCommentView compose a CommentView from a Comment
func BuildCommentView(comment prototypes.Comment) CommentView {
	return CommentView{
		Comment:   comment,
		CreatedAt: comment.CreatedAt.Unix(),
		UpdatedAt: comment.UpdatedAt.Unix(),
	}
}

// RenderComment ...
func RenderComment(ctx context.Context, w http.ResponseWriter, comment prototypes.Comment) (err error) {
	commentView := BuildCommentView(comment)
	err = Render(ctx, w, commentView)
	return
}
