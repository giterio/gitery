package controllers

import (
	"encoding/json"
	"net/http"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// TagHandler ...
type TagHandler struct {
	Model prototypes.TagService
}

func (h *TagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodDelete:
		err = h.handleDelete(w, r)
	default:
		err = models.ForbiddenError(ctx, nil)
	}
	if err != nil {
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// Create a tag
// POST /tag/
func (h *TagHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	param := struct {
		TagName string `json:"tagName"`
		PostID  int    `json:"postID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// create new comment record in DB
	tag, err := h.Model.Assign(ctx, param.PostID, param.TagName)
	if err != nil {
		return
	}
	err = views.Render(ctx, w, tag)
	return
}

// Delete a comment
// DELETE /comment/1
func (h *TagHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	return
}
