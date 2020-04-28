package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// PostHandler ...
type PostHandler struct {
	Model prototypes.PostService
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodPatch:
		err = h.handlePatch(w, r)
	case http.MethodDelete:
		err = h.handleDelete(w, r)
	}
	if err != nil {
		ctx := r.Context()
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// Retrieve a post
// GET /post/1
func (h *PostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderPost(ctx, w, post)
	return
}

// Create a post
// POST /post/
func (h *PostHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	post := prototypes.Post{}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Create(ctx, &post)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderPost(ctx, w, post)
	return
}

// Update a post
// PUT /post/1
func (h *PostHandler) handlePatch(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Update(ctx, &post)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderPost(ctx, w, post)
	return
}

// Delete a post
// DELETE /post/1
func (h *PostHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Delete(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
