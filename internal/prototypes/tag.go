package prototypes

import "context"

// Tag ...
type Tag struct {
	ID   *int   `json:"id"`
	Name string `json:"name"`
}

// TagService ...
type TagService interface {
	Assign(ctx context.Context, userID int, postID int, tagName string) (tag *Tag, err error)
	Remove(ctx context.Context, userID int, postID int, tagID int) (err error)
}
