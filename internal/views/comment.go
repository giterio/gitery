package views

import (
	"context"
	"gitery/internal/prototypes"
	"net/http"
)

// CommentView ...
type CommentView struct {
	prototypes.Comment
	CreatedAt int64          `json:"createdAt"`
	UpdatedAt int64          `json:"updatedAt"`
	Author    *UserView      `json:"author,omitempty"`
	Comments  []*CommentView `json:"comments,omitempty"`
}

// BuildCommentView compose a CommentView from a Comment
func BuildCommentView(comment *prototypes.Comment) (commentView CommentView) {
	commentView = CommentView{
		Comment:   *comment,
		CreatedAt: comment.CreatedAt.Unix(),
		UpdatedAt: comment.UpdatedAt.Unix(),
	}
	if comment.Author != nil {
		author := BuildUserView(comment.Author)
		commentView.Author = &author
	}
	if len(comment.Comments) > 0 {
		commentView.Comments = []*CommentView{}
		for _, childComment := range comment.Comments {
			childCommentView := BuildCommentView(childComment)
			commentView.Comments = append(commentView.Comments, &childCommentView)
		}
	}

	return
}

// RenderComment ...
func RenderComment(ctx context.Context, w http.ResponseWriter, comment *prototypes.Comment) (err error) {
	commentView := BuildCommentView(comment)
	err = Render(ctx, w, commentView)
	return
}
