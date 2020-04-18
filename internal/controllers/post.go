package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitery/internal/domains"
)

// PostHandler ...
type PostHandler struct {
	Model domains.PostService
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
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

// Retrieve a post
// GET /post/1
func (h *PostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	param, _ := ExtractRoute(ctx)
	id, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(post, "", "\t\t")
	if err != nil {
		return
	}
	w.Write(output)
	return
}

// Create a post
// POST /post/
func (h *PostHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	post := domains.Post{}
	json.Unmarshal(body, &post)
	ctx := r.Context()
	err = h.Model.Create(ctx, &post)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(post, "", "\t\t")
	if err != nil {
		return
	}
	w.Write(output)
	return
}

// Update a post
// PUT /post/1
func (h *PostHandler) handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	param, _ := ExtractRoute(ctx)
	id, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	// parse json from request body
	json.Unmarshal(body, &post)
	err = h.Model.Update(ctx, &post)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(post, "", "\t\t")
	if err != nil {
		return
	}
	w.Write(output)
	return
}

// Delete a post
// DELETE /post/1
func (h *PostHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	param, _ := ExtractRoute(ctx)
	id, err := strconv.Atoi(param)
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
