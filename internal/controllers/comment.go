package controllers

import (
	"encoding/json"
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
	ctx := r.Context()
	switch r.Method {
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

// Create a comment
// POST /comment/
func (h *CommentHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// retrieve comment data from request body
	comment := prototypes.Comment{}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// set user ID for comment
	comment.UserID = payload.Pub.ID
	// create new comment record in DB
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
func (h *CommentHandler) handlePatch(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// parse comment ID from URL
	resource, _ := models.ShiftRoute(r)
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// fetch comment from DB
	comment, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	// the comment requested to update does not belong to current user
	if *comment.UserID != *payload.Pub.ID {
		err = models.ForbiddenError(ctx, nil)
		return
	}
	// merge and update comment instance
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	// update comment in DB
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
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	// parse comment ID from URL
	resource, _ := models.ShiftRoute(r)
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	comment := prototypes.Comment{
		ID:     &id,
		UserID: payload.Pub.ID,
	}
	// delete comment from DB
	err = h.Model.Delete(ctx, &comment)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	views.RenderEmpty(ctx, w)
	return
}
