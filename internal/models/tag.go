package models

import (
	"context"
	"database/sql"
)

// TagService ...
type TagService struct {
	DB *sql.DB
}

// Create ...
func (ts *TagService) Create(ctx context.Context, postID int, tagName string) {
	return
}
