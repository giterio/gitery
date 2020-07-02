package models

import (
	"context"
	"database/sql"
	"log"

	"gitery/internal/prototypes"
)

// TagService ...
type TagService struct {
	DB *sql.DB
}

// Assign is to add tag for post
func (ts *TagService) Assign(ctx context.Context, userID int, postID int, tagName string) (tag prototypes.Tag, err error) {
	txn, err := ts.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// check if user is author of the post
	var isAuthor bool
	err = txn.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1 AND user_id = $2)
		`, postID, userID).Scan(&isAuthor)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	} else if !isAuthor {
		err = ForbiddenError(ctx, nil)
		return
	}

	// upsert tag
	tag = prototypes.Tag{}
	err = txn.QueryRowContext(ctx, `
	 	WITH ins AS (
		 	INSERT INTO tags (name) VALUES ($1)
		 	ON CONFLICT (name) DO NOTHING
		 	RETURNING id, name
	 	)
	 	SELECT id, name FROM ins
		UNION ALL
		SELECT id, name FROM tags
		WHERE name = $1 LIMIT 1
		`, tagName).Scan(&tag.ID, &tag.Name)
	if err != nil {
		if rollbackErr := txn.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
		err = HandleDatabaseQueryError(ctx, err)
		return
	}

	// create link of post and tag
	_, err = txn.ExecContext(ctx, `
		INSERT INTO post_tag (post_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT (post_id, tag_id)
		DO NOTHING
		`, postID, tag.ID)
	if err != nil {
		if rollbackErr := txn.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
		err = TransactionError(ctx, err)
		return
	}

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// Remove is to remove tag from post
func (ts *TagService) Remove(ctx context.Context, userID int, postID int, tagID int) (err error) {
	txn, err := ts.DB.Begin()
	if err != nil {
		err = ServerError(ctx, err)
		return
	}

	// check if user is author of the post
	var isAuthor bool
	err = txn.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1 AND user_id = $2)
		`, postID, userID).Scan(&isAuthor)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
		return
	} else if !isAuthor {
		err = ForbiddenError(ctx, nil)
		return
	}

	_, err = txn.ExecContext(ctx, `
		DELETE FROM post_tag WHERE post_id = $1 AND tag_id = $2
		`, postID, tagID)
	if err != nil {
		if rollbackErr := txn.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
		err = TransactionError(ctx, err)
		return
	}

	if err = txn.Commit(); err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
