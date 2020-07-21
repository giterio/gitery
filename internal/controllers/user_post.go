package controllers

import (
	"net/http"
	"strconv"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// UserPostHandler ...
type UserPostHandler struct {
	Model prototypes.UserPostService
}

func (h *UserPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
	default:
		err = models.ForbiddenError(ctx, nil)
	}
	if err != nil {
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// fetch current user's post
// GET /user/:id/posts or /user/posts
func (h *UserPostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	resource, _ := models.ShiftRoute(r)
	var id int
	if resource == "posts" {
		// get current user's id from JWT payload
		payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
		if !ok {
			err = models.AuthorizationError(ctx, err)
			return
		}
		id = *payload.Pub.ID
	} else {
		// parse id from URL
		id, err = strconv.Atoi(resource)
		if err != nil {
			err = models.BadRequestError(ctx, err)
			return
		}
	}
	// fetch user posts from DB
	posts, err := h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	err = views.RenderPostList(ctx, w, posts)
	return
}
