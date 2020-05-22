package views

import (
	"context"
	"gitery/internal/prototypes"
	"net/http"
)

// PostView is the response data structure for Post
type PostView struct {
	prototypes.Post
	Comments  []CommentView `json:"comments,omitempty"`
	CreatedAt int64         `json:"createdAt"`
	UpdatedAt int64         `json:"updatedAt"`
}

// BuildPostView compose PostView from a Post
func BuildPostView(post prototypes.Post) PostView {
	comments := []CommentView{}
	if post.Comments != nil {
		for _, comment := range post.Comments {
			commentView := BuildCommentView(comment)
			comments = append(comments, commentView)
		}
	}
	return PostView{
		Post:      post,
		Comments:  comments,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
	}
}

// RenderPost writes the PostView response to http connection
func RenderPost(ctx context.Context, w http.ResponseWriter, post prototypes.Post) (err error) {
	postView := BuildPostView(post)
	err = Render(ctx, w, postView)
	return
}

// RenderPostList ...
func RenderPostList(ctx context.Context, w http.ResponseWriter, posts []prototypes.Post) (err error) {
	postListView := []PostView{}
	for _, post := range posts {
		postListView = append(postListView, BuildPostView(post))
	}
	err = Render(ctx, w, postListView)
	return
}
