package models

import (
	"context"
	"database/sql"

	"gitery/internal/prototypes"
)

// UserPostService supply service about user's posts
type UserPostService struct {
	DB *sql.DB
}

// Fetch user's all posts
func (ups *UserPostService) Fetch(ctx context.Context, id int) (posts []*prototypes.Post, err error) {
	// postMap is used to assemble posts and comments efficiently
	postMap := map[int]*prototypes.Post{}
	postList := []*prototypes.Post{}
	// query all the posts of the user
	postRows, err := ups.DB.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM posts
		WHERE user_id =$1
		`, id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer postRows.Close()

	// fill the posts into postMap using post ID as the key
	for postRows.Next() {
		post := prototypes.Post{UserID: &id, Comments: []*prototypes.Comment{}, Tags: []*prototypes.Tag{}}
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return
		}
		postMap[*post.ID] = &post
		postList = append(postList, &post)
	}

	// query all the tags related to the posts
	tagRows, err := ups.DB.QueryContext(ctx, `
		SELECT tags.id, tags.name, post_tag.post_id
		FROM tags INNER JOIN post_tag
		ON tags.id = post_tag.tag_id AND post_tag.post_id IN (SELECT id FROM posts WHERE user_id =$1)
		`, id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var postID int
		tag := prototypes.Tag{}
		err = tagRows.Scan(&tag.ID, &tag.Name, &postID)
		if err != nil {
			return
		}
		post := postMap[postID]
		post.Tags = append(post.Tags, &tag)
	}

	// convert postMap to post list
	posts = []*prototypes.Post{}
	for _, post := range postList {
		posts = append(posts, postMap[*post.ID])
	}
	return
}
