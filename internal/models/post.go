package models

import (
	"context"
	"database/sql"

	"gitery/internal/prototype"
)

// PostService implement prototype.PostService interface
type PostService struct {
	DB *sql.DB
}

// Fetch single post
func (ps *PostService) Fetch(ctx context.Context, id int) (post prototype.Post, err error) {
	post = prototype.Post{}
	post.Comments = []prototype.Comment{}
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
		comment := prototype.Comment{PostID: &id}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.Author, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	return
}

// Create a new post
func (ps *PostService) Create(ctx context.Context, post *prototype.Post) (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id, created_at, updated_at"
	stmt, err := ps.DB.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, post.Content, post.Author).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	post.Comments = []prototype.Comment{}
	return
}

// Update a post
func (ps *PostService) Update(ctx context.Context, post *prototype.Post) (err error) {
	err = ps.DB.QueryRowContext(ctx, "update posts set content = $2, author = $3 where id = $1 returning updated_at",
		post.ID, post.Content, post.Author).Scan(&post.UpdatedAt)
	return
}

// Delete a post
func (ps *PostService) Delete(ctx context.Context, id int) (err error) {
	_, err = ps.DB.ExecContext(ctx, "delete from posts where id = $1", id)
	return
}
