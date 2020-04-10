package models

import (
	"database/sql"
	"errors"
)

// Comment ...
type Comment struct {
	Db      *sql.DB `json:"-"`
	ID      int     `json:"id"`
	Content string  `json:"content"`
	Author  string  `json:"author"`
	PostID  *int    `json:"post_id"`
}

// Create ...
func (comment *Comment) Create() (err error) {
	if comment.PostID == nil {
		err = errors.New("Post not found")
		return
	}
	err = comment.Db.QueryRow("insert into comments (content, author, post_id) values ($1, $2, $3) returning id",
		comment.Content, comment.Author, comment.PostID).Scan(&comment.ID)
	return
}
