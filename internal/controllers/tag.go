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
// POST /tag
func (h *TagHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()

	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// struct to receive param
	param := struct {
		TagName string `json:"tagName"`
		PostID  *int   `json:"postID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// assign tag to post in DB
	tag, err := h.Model.Assign(ctx, *payload.Pub.ID, *param.PostID, param.TagName)
	if err != nil {
		return
	}
	err = views.Render(ctx, w, tag)
	return
}

// Remove a tag
// DELETE /tag
func (h *TagHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()

	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// struct to receive param
	param := struct {
		TagID  *int `json:"tagID"`
		PostID *int `json:"postID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// remove tag of post in DB
	err = h.Model.Remove(ctx, *payload.Pub.ID, *param.PostID, *param.TagID)
	if err != nil {
		return
	}

	views.RenderEmpty(ctx, w)
	return
}
