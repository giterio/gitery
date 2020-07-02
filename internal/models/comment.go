package models

import (
	"context"
	"database/sql"
	"time"

	"gitery/internal/prototypes"
)

// CommentService implement the prototypes.CommentService interface
type CommentService struct {
	DB *sql.DB
}

// Fetch ...
func (cs *CommentService) Fetch(ctx context.Context, id int) (comment prototypes.Comment, err error) {
	comment = prototypes.Comment{}
	err = cs.DB.QueryRowContext(ctx, `
		SELECT id, content, user_id, post_id, parent_id, is_deleted, created_at, updated_at
		FROM comments
		WHERE id = $1
		`, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentID,
		&comment.IsDeleted,
		&comment.CreatedAt,
		&comment.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// FetchDetail is to fetch single comment detail
func (cs *CommentService) FetchDetail(ctx context.Context, id int) (comment prototypes.Comment, err error) {
	txn, err := cs.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// query comment from DB
	comment = prototypes.Comment{}
	err = txn.QueryRowContext(ctx, `
		SELECT id, content, user_id, post_id, parent_id, is_deleted, created_at, updated_at
		FROM comments
		WHERE id = $1
		`, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentID,
		&comment.IsDeleted,
		&comment.CreatedAt,
		&comment.UpdatedAt)
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
		`, *comment.UserID).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}
	comment.Author = &user

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Create comment
func (cs *CommentService) Create(ctx context.Context, comment *prototypes.Comment) (err error) {
	txn, err := cs.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// check if user exist
	var isUserExist bool
	err = txn.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)
		`, comment.UserID).Scan(&isUserExist)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	} else if !isUserExist {
		err = NotFoundError(ctx, err)
		return
	}

	// check if post exist
	var isPostExist bool
	err = txn.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1)
		`, comment.PostID).Scan(&isPostExist)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	} else if !isPostExist {
		err = NotFoundError(ctx, err)
		return
	}

	if comment.ParentID != nil {
		var isCommentExist bool
		err = txn.QueryRowContext(ctx, `
			SELECT EXISTS (SELECT 1 FROM comments WHERE id = $1 AND post_id = $2)
			`, comment.ParentID, comment.PostID).Scan(&isCommentExist)
		if err != nil {
			err = HandleDatabaseQueryError(ctx, err)
			return
		} else if !isCommentExist {
			err = NotFoundError(ctx, err)
			return
		}
	}

	// insert new comments
	err = txn.QueryRowContext(ctx, `
		INSERT INTO comments (content, user_id, post_id, parent_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
		`, comment.Content, comment.UserID, comment.PostID, comment.ParentID).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Update a comment
func (cs *CommentService) Update(ctx context.Context, comment *prototypes.Comment) (err error) {
	err = cs.DB.QueryRowContext(ctx, `
		UPDATE comments
		SET content = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
		`, comment.ID, comment.UserID, comment.Content, time.Now()).Scan(&comment.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Delete a comment
func (cs *CommentService) Delete(ctx context.Context, comment *prototypes.Comment) (err error) {
	_, err = cs.DB.ExecContext(ctx, `
		UPDATE comments
		SET is_deleted = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2
		`, comment.ID, comment.UserID, true, time.Now())
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
