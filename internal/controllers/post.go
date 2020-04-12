package controllers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"gitery/internal/models"
	route "gitery/internal/utils"
)

// PostHandler ...
type PostHandler struct {
	Model models.PostService
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = h.handleGet(w, r)
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

// Retrieve a post
// GET /post/1
func (h *PostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	route, err := route.Of(r)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(route.BasePath())
	if err != nil {
		return
	}
	ctx := r.Context()
	err = h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(h.Model, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// Create a post
// POST /post/
func (h *PostHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
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
func (h *PostHandler) handlePut(w http.ResponseWriter, r *http.Request) (err error) {
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
func (h *PostHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
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
