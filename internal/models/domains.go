package models

import (
	"context"
)

// PostService ...
type PostService interface {
	Fetch(ctx context.Context, id int) (err error)
	Create(ctx context.Context) (err error)
	Update(ctx context.Context) (err error)
	Delete(ctx context.Context) (err error)
}
