package models

import (
	"database/sql"
)

// Text ...
type Text interface {
	Fetch(id int) (err error)
	Create() (err error)
	Update() (err error)
	Delete() (err error)
}

// Post ...
type Post struct {
	Db      *sql.DB
	ID      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Fetch single post
func (post *Post) Fetch(id int) (err error) {
	err = post.Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.ID, &post.Content, &post.Author)
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := post.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.ID)
	return
}

// Update a post
func (post *Post) Update() (err error) {
	_, err = post.Db.Exec("update posts set content = $2, author = $3 where id = $1", post.ID, post.Content, post.Author)
	return
}

// Delete a post
func (post *Post) Delete() (err error) {
	_, err = post.Db.Exec("delete from posts where id = $1", post.ID)
	return
}
