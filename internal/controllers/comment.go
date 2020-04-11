package controllers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"gitery/internal/models"
)

// CommentHandler ...
type CommentHandler struct {
	Model models.CommentService
}

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "POST":
		err = h.handlePost(w, r)
	case "PUT":
		err = h.handlePut(w, r)
	case "DELETE":
		err = h.handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create a post
// POST /post/
func (h *CommentHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, h.Model)
	ctx := r.Context()
	err = h.Model.Create(ctx)
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

// Update a post
// PUT /post/1
func (h *CommentHandler) handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	ctx := r.Context()
	err = h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// parse json from request body
	json.Unmarshal(body, h.Model)
	err = h.Model.Update(ctx)
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

// Delete a post
// DELETE /post/1
func (h *CommentHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	ctx := r.Context()
	err = h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	err = h.Model.Delete(ctx)
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
