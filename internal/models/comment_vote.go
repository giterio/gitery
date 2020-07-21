package models

import (
	"context"
	"database/sql"
	"gitery/internal/prototypes"
)

// CommentVoteService ...
type CommentVoteService struct {
	DB *sql.DB
}

// FetchVotes ...
func (cvs *CommentVoteService) FetchVotes(ctx context.Context, userID int, postID int) (votes []*prototypes.CommentVote, err error) {
	voteRows, err := cvs.DB.QueryContext(ctx, `
		SELECT comment_id, vote
		FROM comment_vote
		WHERE user_id = $1 AND comment_id IN (SELECT comment_id FROM comments WHERE post_id = $2)
		`, userID, postID)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer voteRows.Close()

	votes = []*prototypes.CommentVote{}
	for voteRows.Next() {
		commentVote := prototypes.CommentVote{UserID: &userID}
		err = voteRows.Scan(&commentVote.CommentID, &commentVote.Vote)
		if err != nil {
			return
		}
		votes = append(votes, &commentVote)
	}
	return
}

// Vote ...
func (cvs *CommentVoteService) Vote(ctx context.Context, commentVote *prototypes.CommentVote) (err error) {
	statement := `
		INSERT INTO comment_vote (user_id, comment_id, vote)
		VALUES ($1, $2, $3) ON CONFLICT (user_id, comment_id) DO NOTHING`
	stmt, err := cvs.DB.PrepareContext(ctx, statement)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, commentVote.UserID, commentVote.CommentID, commentVote.Vote)
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
