package testdata

import (
	"context"
	"gitery/internal/prototypes"
)

// FakePostService ...
type FakePostService struct{}

// Fetch ...
func (ps *FakePostService) Fetch(ctx context.Context, id int) (post prototypes.Post, err error) {
	post = prototypes.Post{ID: &id}
	return
}

// Create ...
func (ps *FakePostService) Create(ctx context.Context, post *prototypes.Post) (err error) {
	return
}

// Update ...
func (ps *FakePostService) Update(ctx context.Context, post *prototypes.Post) (err error) {
	return
}

// Delete ...
func (ps *FakePostService) Delete(ctx context.Context, post *prototypes.Post) (err error) {
	return
}
