package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitery/internal/models"
	"gitery/internal/prototype"
	"gitery/internal/views"
)

// CommentHandler ...
type CommentHandler struct {
	Model prototype.CommentService
}

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodPut:
		err = h.handlePut(w, r)
	case http.MethodDelete:
		err = h.handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create a comment
// POST /comment/
func (h *CommentHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	comment := prototype.Comment{}
	json.Unmarshal(body, &comment)
	ctx := r.Context()
	err = h.Model.Create(ctx, &comment)
	if err != nil {
		return
	}
	err = views.Render(w, comment)
	return
}

// Update a comment
// PUT /comment/1
func (h *CommentHandler) handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	comment, err := h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// parse json from request body
	json.Unmarshal(body, &comment)
	err = h.Model.Update(ctx, &comment)
	if err != nil {
		return
	}
	err = views.Render(w, comment)
	return
}

// Delete a comment
// DELETE /comment/1
func (h *CommentHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		return
	}
	err = h.Model.Delete(ctx, id)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
