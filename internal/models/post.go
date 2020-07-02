package models

import (
	"context"
	"database/sql"
	"time"

	"gitery/internal/prototypes"
)

// PostService implement prototypes.PostService interface
type PostService struct {
	DB *sql.DB
}

// Fetch ...
func (ps *PostService) Fetch(ctx context.Context, id int) (post prototypes.Post, err error) {
	post = prototypes.Post{}
	err = ps.DB.QueryRowContext(ctx, `
		SELECT id, title, content, user_id, created_at, updated_at
		FROM posts
		WHERE id = $1
		`, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// FetchDetail is to fetch single post detail
func (ps *PostService) FetchDetail(ctx context.Context, id int) (post prototypes.Post, err error) {
	txn, err := ps.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// query post data
	post = prototypes.Post{}
	err = txn.QueryRowContext(ctx, `
		SELECT id, title, content, user_id, created_at, updated_at
		FROM posts
		WHERE id = $1
		`, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}

	// query user information
	user := prototypes.User{}
	err = txn.QueryRowContext(ctx, `
		SELECT id, email, hashed_pwd, nickname, created_at, updated_at
		FROM users
		WHERE id = $1
		`, *post.UserID).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}
	post.Author = &user

	// query tags related to the post
	tagRows, err := txn.QueryContext(ctx, `
		SELECT id, name FROM tags
		WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = $1)
		`, id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer tagRows.Close()

	post.Tags = []prototypes.Tag{}
	// Assemble tags with post structure
	for tagRows.Next() {
		tag := prototypes.Tag{}
		if err = tagRows.Scan(&tag.ID, &tag.Name); err != nil {
			err = TransactionError(ctx, err)
			return
		}
		post.Tags = append(post.Tags, tag)
	}

	// query comments related to the post
	commentRows, err := txn.QueryContext(ctx, `
		SELECT id, content, user_id, created_at, updated_at
		FROM comments
		WHERE post_id =$1
		`, id)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer commentRows.Close()

	post.Comments = []prototypes.Comment{}
	// Assemble comments with post structure
	for commentRows.Next() {
		comment := prototypes.Comment{PostID: &id}
		err = commentRows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			err = TransactionError(ctx, err)
			return
		}
		post.Comments = append(post.Comments, comment)
	}

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// FetchList is to get latest posts
func (ps *PostService) FetchList(ctx context.Context, limit int, offset int) (posts []prototypes.Post, err error) {
	if limit == 0 {
		limit = 10
	}

	txn, err := ps.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// postMap is used to assemble posts and comments efficiently
	postMap := map[int]*prototypes.Post{}

	// query all the posts of the user
	postRows, err := txn.QueryContext(ctx, `
		SELECT posts.id, posts.title, posts.user_id, posts.created_at, posts.updated_at,
		users.id AS user_id, users.email, users.nickname, users.created_at AS user_created_at, users.updated_at AS user_updated_at
		FROM posts LEFT JOIN users ON (posts.user_id = users.id)
		ORDER BY posts.updated_at DESC
		LIMIT $1 OFFSET $2
		`, limit, offset)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer postRows.Close()

	// fill the posts into list
	for postRows.Next() {
		post := prototypes.Post{Comments: []prototypes.Comment{}, Tags: []prototypes.Tag{}, Author: &prototypes.User{}}
		err = postRows.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt,
			&post.Author.ID, &post.Author.Email, &post.Author.Nickname, &post.Author.CreatedAt, &post.Author.UpdatedAt)
		if err != nil {
			return
		}
		postMap[*post.ID] = &post
	}

	// query all the tags related to the posts
	tagRows, err := txn.QueryContext(ctx, `
		SELECT tags.id, tags.name, post_tag.post_id
		FROM tags INNER JOIN post_tag
		ON tags.id = post_tag.tag_id AND post_tag.post_id IN (SELECT id FROM posts ORDER BY posts.updated_at DESC LIMIT $1 OFFSET $2)
		`, limit, offset)
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
		post.Tags = append(post.Tags, tag)
	}

	// convert postMap to post list
	posts = []prototypes.Post{}
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Create a new post
func (ps *PostService) Create(ctx context.Context, post *prototypes.Post) (err error) {
	statement := `
		INSERT INTO posts (title, content, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	stmt, err := ps.DB.PrepareContext(ctx, statement)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, post.Title, post.Content, post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}
	post.Comments = []prototypes.Comment{}
	return
}

// Update a post
func (ps *PostService) Update(ctx context.Context, post *prototypes.Post) (err error) {
	err = ps.DB.QueryRowContext(ctx, `
		UPDATE posts SET title = $3, content = $4, updated_at = $5
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
		`, post.ID, post.UserID, post.Title, post.Content, time.Now()).Scan(&post.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Delete a post
func (ps *PostService) Delete(ctx context.Context, post *prototypes.Post) (err error) {
	_, err = ps.DB.ExecContext(ctx, `
		DELETE FROM posts
		WHERE id = $1 AND user_id =$2`, post.ID, post.UserID)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
