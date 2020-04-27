package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// CommentHandler ...
type CommentHandler struct {
	Model prototypes.CommentService
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
		ctx := r.Context()
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// Create a comment
// POST /comment/
func (h *CommentHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	len := r.ContentLength
	body := make([]byte, len)
	_, err = r.Body.Read(body)
	if err != io.EOF {
		err = models.BadRequestError(ctx)
		return
	}
	comment := prototypes.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Create(ctx, &comment)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderComment(ctx, w, comment)
	return
}

// Update a comment
// PUT /comment/1
func (h *CommentHandler) handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	comment, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	_, err = r.Body.Read(body)
	if err != io.EOF {
		err = models.BadRequestError(ctx)
		return
	}
	// parse json from request body
	err = json.Unmarshal(body, &comment)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Update(ctx, &comment)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderComment(ctx, w, comment)
	return
}

// Delete a comment
// DELETE /comment/1
func (h *CommentHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
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
