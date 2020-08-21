package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// CommentVoteHandler ...
// POST /comments/votes
type CommentVoteHandler struct {
	Model prototypes.CommentVoteService
}

func (h *CommentVoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *CommentVoteHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// retrieve postID from query
	var postID int
	q := r.URL.Query()
	postID, err = strconv.Atoi(q.Get("post_id"))
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// fetch user's votes on post's comments
	votes, err := h.Model.FetchVotes(ctx, *payload.Pub.ID, postID)
	if err != nil {
		return
	}
	err = views.Render(ctx, w, votes)
	return
}

// Vote up/down comment
// POST /comments/votes
func (h *CommentVoteHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// retrieve commentID from request body
	commentVote := prototypes.CommentVote{UserID: payload.Pub.ID}
	err = json.NewDecoder(r.Body).Decode(&commentVote)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// vote comment
	err = h.Model.Vote(ctx, &commentVote)
	if err != nil {
		return
	}
	err = views.RenderEmpty(ctx, w)
	return
}

// cancel vote
// POST /comments/votes
func (h *CommentVoteHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	// Check user auth
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}

	// retrieve commentID from request body
	param := struct {
		CommentID *int `json:"commentID"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}

	// cancel vote
	err = h.Model.Cancel(ctx, *payload.Pub.ID, *param.CommentID)
	if err != nil {
		return
	}
	err = views.RenderEmpty(ctx, w)
	return
}
