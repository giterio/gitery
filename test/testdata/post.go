package testdata

import (
	"context"
	"gitery/internal/domains"
)

// FakePostService ...
type FakePostService struct{}

// Fetch ...
func (ps *FakePostService) Fetch(ctx context.Context, id int) (post domains.Post, err error) {
	post = domains.Post{ID: &id}
	return
}

// Create ...
func (ps *FakePostService) Create(ctx context.Context, post *domains.Post) (err error) {
	return
}

// Update ...
func (ps *FakePostService) Update(ctx context.Context, post *domains.Post) (err error) {
	return
}

// Delete ...
func (ps *FakePostService) Delete(ctx context.Context, id int) (err error) {
	return
}
