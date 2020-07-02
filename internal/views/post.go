package views

import (
	"context"
	"gitery/internal/prototypes"
	"net/http"
)

// PostView is the response data structure for Post
type PostView struct {
	prototypes.Post
	CreatedAt int64         `json:"createdAt"`
	UpdatedAt int64         `json:"updatedAt"`
	Author    *UserView     `json:"author,omitempty"`
	Comments  []CommentView `json:"comments,omitempty"`
}

// BuildPostView compose PostView from a Post
func BuildPostView(post prototypes.Post) (postView PostView) {
	postView = PostView{
		Post:      post,
		CreatedAt: post.CreatedAt.Unix(),
		UpdatedAt: post.UpdatedAt.Unix(),
	}
	if post.Comments != nil {
		comments := []CommentView{}
		for _, comment := range post.Comments {
			commentView := BuildCommentView(*comment)
			comments = append(comments, commentView)
		}
		postView.Comments = comments
	}
	if post.Author != nil {
		author := BuildUserView(*post.Author)
		postView.Author = &author
	}
	return
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
