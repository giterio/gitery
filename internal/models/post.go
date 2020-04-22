package models

import (
	"context"
	"database/sql"
	"time"

	"gitery/internal/domains"
)

// PostService ...
type PostService struct {
	DB *sql.DB
}

// Fetch single post
func (ps *PostService) Fetch(ctx context.Context, id int) (post domains.Post, err error) {
	post = domains.Post{}
	post.Comments = []domains.Comment{}
	err = ps.DB.QueryRowContext(ctx, "select id, content, author, created_at, updated_at from posts where id = $1", id).Scan(
		&post.ID, &post.Content, &post.Author, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return
	}
	rows, err := ps.DB.QueryContext(ctx, "select id, content, author, created_at, updated_at from comments where post_id =$1", id)
	if err != nil {
		return
	}
	for rows.Next() {
		comment := domains.Comment{PostID: &id}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	return
}

// Create a new post
func (ps *PostService) Create(ctx context.Context, post *domains.Post) (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := ps.DB.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.ID)
	return
}

// Update a post
func (ps *PostService) Update(ctx context.Context, post *domains.Post) (err error) {
	_, err = ps.DB.ExecContext(ctx, "update posts set content = $2, author = $3, updated_at = $4 where id = $1", post.ID, post.Content, post.Author, time.Now())
	return
}

// Delete a post
func (ps *PostService) Delete(ctx context.Context, id int) (err error) {
	_, err = ps.DB.ExecContext(ctx, "delete from posts where id = $1", id)
	return
}
