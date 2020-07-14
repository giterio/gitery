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
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodPatch:
		err = h.handlePatch(w, r)
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

// Retrieve a post
// GET /post/:id or /post?limit=10&offset=0
func (h *PostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// parse post ID from URL
	resource, _ := models.ShiftRoute(r)

	switch resource {
	// handle /post?limit=10&offset=0
	case "":
		// pre-declaration to avoid shadowing of variable err
		var limit, offset int
		var posts []*prototypes.Post
		q := r.URL.Query()
		limit, err = strconv.Atoi(q.Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err = strconv.Atoi(q.Get("offset"))
		if err != nil {
			offset = 0
		}
		posts, err = h.Model.FetchList(ctx, limit, offset)
		if err != nil {
			return
		}
		err = views.RenderPostList(ctx, w, posts)
		return

	// handle /post/:id
	default:
		// pre-declaration to avoid shadowing of variable err
		var id int
		var post *prototypes.Post
		id, err = strconv.Atoi(resource)
		if err != nil {
			err = models.BadRequestError(ctx, err)
			return
		}
		// fetch post from DB
		post, err = h.Model.FetchDetail(ctx, id)
		if err != nil {
			return
		}
		err = views.RenderPost(ctx, w, post)
		return
	}
}

// Create a post
// POST /post/
func (h *PostHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// retrieve post data from request body
	post := prototypes.Post{}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// set user ID for post
	post.UserID = payload.Pub.ID
	err = h.Model.Create(ctx, &post)
	if err != nil {
		return
	}
	err = views.RenderPost(ctx, w, &post)
	return
}

// Update a post
// PATCH /post/1
func (h *PostHandler) handlePatch(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// parse post ID from URL
	resource, _ := models.ShiftRoute(r)
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// fetch post from DB
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		return
	}
	// the post requested to update does not belong to current user
	if *post.UserID != *payload.Pub.ID {
		err = models.ForbiddenError(ctx, nil)
		return
	}
	// merge and update post instance
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// update post in DB
	err = h.Model.Update(ctx, post)
	if err != nil {
		return
	}
	err = views.RenderPost(ctx, w, post)
	return
}

// Delete a post
// DELETE /post/1
func (h *PostHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// parse post ID from URL
	resource, _ := models.ShiftRoute(r)
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	post := prototypes.Post{
		ID:     &id,
		UserID: payload.Pub.ID,
	}
	// delete post from DB
	err = h.Model.Delete(ctx, &post)
	if err != nil {
		return
	}
	views.RenderEmpty(ctx, w)
	return
}
