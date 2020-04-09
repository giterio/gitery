package models

import (
	"database/sql"
	"errors"
)

// Comment ...
type Comment struct {
	Db      *sql.DB
	ID      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Post    *Post
}

// Create ...
func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("Post not found")
		return
	}
	err = comment.Db.QueryRow("insert into comments (content, author, post_id) values ($1, $2, $3) returning id",
		comment.Content, comment.Author, comment.Post.ID).Scan(&comment.ID)
	return
}
