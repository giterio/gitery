package models

import (
	"context"
	"database/sql"
)

// PostLikeService ...
type PostLikeService struct {
	DB *sql.DB
}

// FetchLikes ...
func (pls *PostLikeService) FetchLikes(ctx context.Context, userID int) (likes []*int, err error) {
	likeRows, err := pls.DB.QueryContext(ctx, `
		SELECT post_id
		FROM post_like
		WHERE user_id = $1
		`, userID)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer likeRows.Close()

	likes = []*int{}
	for likeRows.Next() {
		var postID int
		err = likeRows.Scan(&postID)
		if err != nil {
			return
		}
		likes = append(likes, &postID)
	}
	return
}

// Like a post
func (pls *PostLikeService) Like(ctx context.Context, userID int, postID int) (err error) {
	statement := `
		INSERT INTO post_like (user_id, post_id)
		VALUES ($1, $2) ON CONFLICT (user_id, post_id) DO NOTHING`
	stmt, err := pls.DB.PrepareContext(ctx, statement)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userID, postID)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Unlike ...
func (pls *PostLikeService) Unlike(ctx context.Context, userID int, postID int) (err error) {
	_, err = pls.DB.ExecContext(ctx, `
		DELETE FROM post_like WHERE user_id = $1 AND post_id = $2
		`, userID, postID)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
