package controllers

import (
	"encoding/json"
	"net/http"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// PostLikeHandler ...
type PostLikeHandler struct {
	Model prototypes.PostLikeService
}

func (h *PostLikeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
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

// get likes of post
// GET /posts/likes
func (h *PostLikeHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// fetch user's likes
	likes, err := h.Model.FetchLikes(ctx, *payload.Pub.ID)
	if err != nil {
		return
	}
	err = views.Render(ctx, w, likes)
	return
}

// like post
// POST /posts/likes
func (h *PostLikeHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// retrieve postID from request body
	param := struct {
		PostID *int `json:"postID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// like post
	err = h.Model.Like(ctx, *payload.Pub.ID, *param.PostID)
	if err != nil {
		return
	}
	err = views.RenderEmpty(ctx, w)
	return
}

// like post
// POST /posts/likes
func (h *PostLikeHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// retrieve postID from request body
	param := struct {
		PostID *int `json:"postID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// like post
	err = h.Model.Unlike(ctx, *payload.Pub.ID, *param.PostID)
	if err != nil {
		return
	}
	err = views.RenderEmpty(ctx, w)
	return
}
