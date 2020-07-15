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
func (cs *CommentService) Fetch(ctx context.Context, id int) (comment *prototypes.Comment, err error) {
	comment = &prototypes.Comment{}
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
func (cs *CommentService) FetchDetail(ctx context.Context, id int) (comment *prototypes.Comment, err error) {
	// query comment from DB
	comment = &prototypes.Comment{}
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
		return
	}

	// query user information
	user := prototypes.User{}
	err = cs.DB.QueryRowContext(ctx, `
		SELECT id, email, hashed_pwd, nickname, created_at, updated_at
		FROM users
		WHERE id = $1
		`, *comment.UserID).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}
	comment.Author = &user
	return
}

// Create comment
func (cs *CommentService) Create(ctx context.Context, comment *prototypes.Comment) (err error) {
	// check if user exist
	var isUserExist bool
	err = cs.DB.QueryRowContext(ctx, `
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
	err = cs.DB.QueryRowContext(ctx, `
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
		err = cs.DB.QueryRowContext(ctx, `
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
	err = cs.DB.QueryRowContext(ctx, `
		INSERT INTO comments (content, user_id, post_id, parent_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_deleted, created_at, updated_at
		`, comment.Content, comment.UserID, comment.PostID, comment.ParentID).Scan(
		&comment.ID, &comment.IsDeleted, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}

	// query user information
	user := prototypes.User{}
	err = cs.DB.QueryRowContext(ctx, `
		SELECT id, email, hashed_pwd, nickname, created_at, updated_at
		FROM users
		WHERE id = $1
		`, *comment.UserID).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	}
	comment.Author = &user
	return
}

// Update a comment
func (cs *CommentService) Update(ctx context.Context, comment *prototypes.Comment) (err error) {
	err = cs.DB.QueryRowContext(ctx, `
		UPDATE comments
		SET content = $3, is_deleted = $4, updated_at = $5
		WHERE id = $1 AND user_id = $2
		RETURNING is_deleted, updated_at
		`, comment.ID, comment.UserID, comment.Content, false, time.Now()).Scan(&comment.IsDeleted, &comment.UpdatedAt)
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

// CommentVoteService ...
type CommentVoteService struct {
	DB *sql.DB
}

// Vote ...
func (cvs *CommentVoteService) Vote(ctx context.Context, userID int, commentID int, vote bool) (err error) {
	statement := `
		INSERT INTO comment_vote (user_id, comment_id, vote)
		VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING`
	stmt, err := cvs.DB.PrepareContext(ctx, statement)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userID, commentID, vote)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Cancel ...
func (cvs *CommentVoteService) Cancel(ctx context.Context, userID int, commentID int) (err error) {
	_, err = cvs.DB.ExecContext(ctx, `
		DELETE FROM comment_vote WHERE user_id = $1 AND comment_id = $2
		`, userID, commentID)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
