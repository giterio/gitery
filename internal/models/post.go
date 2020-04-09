package models

import (
	"database/sql"
)

// Post ...
type Post struct {
	DB       *sql.DB
	ID       int       `json:"id"`
	Content  string    `json:"content"`
	Author   string    `json:"author"`
	Comments []Comment `json:"comments"`
}

// Fetch single post
func (post *Post) Fetch(id int) (err error) {
	post.Comments = []Comment{}
	err = post.DB.QueryRow("select id, content, author from posts where id = $1",
		id).Scan(&post.ID, &post.Content, &post.Author)
	rows, err := post.DB.Query("select id, content, author from comments where post_id =$1",
		id)
	if err != nil {
		return
	}
	for rows.Next() {
		comment := Comment{Post: post}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := post.DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.ID)
	return
}

// Update a post
func (post *Post) Update() (err error) {
	_, err = post.DB.Exec("update posts set content = $2, author = $3 where id = $1",
		post.ID, post.Content, post.Author)
	return
}

// Delete a post
func (post *Post) Delete() (err error) {
	_, err = post.DB.Exec("delete from posts where id = $1", post.ID)
	return
}
