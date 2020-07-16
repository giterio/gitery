package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"gitery/internal/prototypes"
)

// PostService implement prototypes.PostService interface
type PostService struct {
	DB *sql.DB
}

// Fetch ...
func (ps *PostService) Fetch(ctx context.Context, id int) (post *prototypes.Post, err error) {
	post = &prototypes.Post{}
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
func (ps *PostService) FetchDetail(ctx context.Context, id int) (post *prototypes.Post, err error) {
	postCh := make(chan *prototypes.Post)
	tagsCh := make(chan []*prototypes.Tag)
	likesCh := make(chan int)
	commentsCh := make(chan []*prototypes.Comment)
	errCh := make(chan error)

	go func() {
		// query post
		post = &prototypes.Post{}
		err = ps.DB.QueryRowContext(ctx, `
			SELECT id, title, content, user_id, created_at, updated_at
			FROM posts
			WHERE id = $1
			`, id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			errCh <- HandleDatabaseQueryError(ctx, err)
			return
		}

		// query user information
		user := prototypes.User{}
		err = ps.DB.QueryRowContext(ctx, `
			SELECT id, email, hashed_pwd, nickname, created_at, updated_at
			FROM users
			WHERE id = $1
			`, *post.UserID).Scan(&user.ID, &user.Email, &user.HashedPwd, &user.Nickname, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			errCh <- HandleDatabaseQueryError(ctx, err)
			return
		}
		post.Author = &user
		postCh <- post
	}()

	go func() {
		// query tags related to the post
		tagRows, err := ps.DB.QueryContext(ctx, `
			SELECT id, name
			FROM tags
			WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = $1)
			`, id)
		if err != nil {
			errCh <- TransactionError(ctx, err)
			return
		}
		defer tagRows.Close()

		tags := []*prototypes.Tag{}
		// Assemble tags with post structure
		for tagRows.Next() {
			tag := prototypes.Tag{}
			if err = tagRows.Scan(&tag.ID, &tag.Name); err != nil {
				errCh <- TransactionError(ctx, err)
				return
			}
			tags = append(tags, &tag)
		}
		tagsCh <- tags
	}()

	go func() {
		var likes int
		err := ps.DB.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM post_like
			WHERE post_id = $1
			`, id).Scan(&likes)
		if err != nil {
			errCh <- TransactionError(ctx, err)
			return
		}
		likesCh <- likes
	}()

	go func() {
		// query comments related to the post
		commentRows, err := ps.DB.QueryContext(ctx, `
			SELECT comments.id, comments.content, comments.user_id, comments.parent_id, comments.is_deleted, comments.created_at, comments.updated_at,
			users.id, users.email, users.nickname, users.created_at, users.updated_at,
			COUNT(CASE WHEN comment_vote.vote = true then 1 ELSE NULL END) as vote_up,
			COUNT(CASE WHEN comment_vote.vote = false then 1 ELSE NULL END) as vote_down
			FROM comments
			INNER JOIN users
			ON comments.post_id = $1 AND comments.user_id = users.id
			LEFT JOIN comment_vote
			ON comments.id = comment_vote.comment_id
			ORDER BY comments.created_at ASC
			`, id)
		if err != nil {
			errCh <- TransactionError(ctx, err)
			return
		}
		defer commentRows.Close()

		commentList := []*prototypes.Comment{}
		commentMap := map[int]*prototypes.Comment{}

		// Assemble commentList and commentMap
		for commentRows.Next() {
			comment := prototypes.Comment{PostID: &id, Author: &prototypes.User{}}
			err = commentRows.Scan(
				&comment.ID,
				&comment.Content,
				&comment.UserID,
				&comment.ParentID,
				&comment.IsDeleted,
				&comment.CreatedAt,
				&comment.UpdatedAt,
				&comment.Author.ID,
				&comment.Author.Email,
				&comment.Author.Nickname,
				&comment.Author.CreatedAt,
				&comment.Author.UpdatedAt,
				&comment.VoteUp,
				&comment.VoteDown,
			)
			if err != nil {
				errCh <- TransactionError(ctx, err)
				return
			}
			commentMap[*comment.ID] = &comment
			commentList = append(commentList, &comment)
		}

		// convert commentMap to tree-like structure
		for _, comment := range commentList {
			if comment.ParentID != nil {
				c := commentMap[*comment.ParentID]
				if c != nil {
					c.Comments = append(c.Comments, comment)
				}
			}
		}

		// filter redundant 1-deep node and format as list
		comments := []*prototypes.Comment{}
		for _, comment := range commentList {
			if comment.ParentID == nil {
				comments = append(comments, commentMap[*comment.ID])
			}
		}
		commentsCh <- comments
	}()

	// receive data from goroutin via channels
	var (
		tags     []*prototypes.Tag
		comments []*prototypes.Comment
		likes    int
	)
	for i := 0; i < 4; i++ {
		select {
		case err = <-errCh:
			return
		case post = <-postCh:
		case tags = <-tagsCh:
		case likes = <-likesCh:
		case comments = <-commentsCh:
		}
	}
	post.Comments = comments
	post.Tags = tags
	post.Likes = likes
	return
}

// FetchList is to get latest posts
func (ps *PostService) FetchList(ctx context.Context, limit int, offset int) (posts []*prototypes.Post, err error) {
	if limit == 0 {
		limit = 10
	}

	// postMap is used to assemble posts and comments efficiently
	postMap := map[int]*prototypes.Post{}
	postList := []*prototypes.Post{}
	postIDs := []int{}

	// query all the posts of the user
	postRows, err := ps.DB.QueryContext(ctx, `
		SELECT posts.id, posts.title, posts.user_id, posts.created_at, posts.updated_at,
		users.id, users.email, users.nickname, users.created_at, users.updated_at
		FROM posts INNER JOIN users
		ON posts.user_id = users.id AND posts.is_deleted = false
		ORDER BY posts.created_at DESC
		LIMIT $1 OFFSET $2
		`, limit, offset)
	if err != nil {
		err = TransactionError(ctx, err)
		return
	}
	defer postRows.Close()

	// fill the posts into list
	for postRows.Next() {
		post := prototypes.Post{Comments: []*prototypes.Comment{}, Tags: []*prototypes.Tag{}, Author: &prototypes.User{}}
		err = postRows.Scan(&post.ID, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt,
			&post.Author.ID, &post.Author.Email, &post.Author.Nickname, &post.Author.CreatedAt, &post.Author.UpdatedAt)
		if err != nil {
			return
		}
		postMap[*post.ID] = &post
		postList = append(postList, &post)
		postIDs = append(postIDs, *post.ID)
	}

	// query all the tags related to the posts
	tagRows, err := ps.DB.QueryContext(ctx, `
		SELECT tags.id, tags.name, post_tag.post_id
		FROM tags INNER JOIN post_tag
		ON tags.id = post_tag.tag_id AND post_tag.post_id IN ($1)
		`, pq.Array(postIDs))
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
		post.Tags = append(post.Tags, &tag)
	}

	// convert postMap to post list
	posts = []*prototypes.Post{}
	for _, post := range postList {
		posts = append(posts, postMap[*post.ID])
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
	post.Comments = []*prototypes.Comment{}
	return
}

// Update a post
func (ps *PostService) Update(ctx context.Context, post *prototypes.Post) (err error) {
	err = ps.DB.QueryRowContext(ctx, `
		UPDATE posts
		SET title = $3, content = $4, is_deleted = $5, updated_at = $6
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at
		`, post.ID, post.UserID, post.Title, post.Content, false, time.Now()).Scan(&post.UpdatedAt)
	if err != nil {
		err = HandleDatabaseQueryError(ctx, err)
	}
	return
}

// Delete a post
func (ps *PostService) Delete(ctx context.Context, post *prototypes.Post) (err error) {
	_, err = ps.DB.ExecContext(ctx, `
		UPDATE posts
		SET is_deleted = $3, updated_at = $4
		WHERE id = $1 AND user_id =$2
		`, post.ID, post.UserID, true, time.Now())
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}

// PostLikeService ...
type PostLikeService struct {
	DB *sql.DB
}

// Like a post
func (pl *PostLikeService) Like(ctx context.Context, userID int, postID int) (err error) {
	statement := `
		INSERT INTO post_like (user_id, post_id)
		VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`
	stmt, err := pl.DB.PrepareContext(ctx, statement)
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
func (pl *PostLikeService) Unlike(ctx context.Context, userID int, postID int) (err error) {
	_, err = pl.DB.ExecContext(ctx, `
		DELETE FROM post_like WHERE user_id = $1 AND post_id = $2
		`, userID, postID)
	if err != nil {
		err = TransactionError(ctx, err)
	}
	return
}
