package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// UserHandler ...
type UserHandler struct {
	Model           prototypes.UserService
	UserPostHandler *UserPostHandler
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	resource, nextRoute := models.CurrentRoute(r).Shift()
	if _, err := strconv.Atoi(resource); err == nil {
		if nextRoute.IsLast() {
			resource = ""
		} else {
			resource, _ = nextRoute.Shift()
		}
	}

	if resource != "" {
		switch resource {
		case "posts":
			h.UserPostHandler.ServeHTTP(w, r)
		default:
			e := models.ForbiddenError(ctx, nil)
			views.RenderError(ctx, w, e)
		}
		return
	}

	// user is the resource to manipulate
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

// Retrieve a user
// GET /user or /user/:id
func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	resource, _ := models.ShiftRoute(r)
	var id int
	if resource == "" {
		// get current user's id from JWT payload
		payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
		if !ok {
			err = models.AuthorizationError(ctx, err)
			return
		}
		id = *payload.Pub.ID
	} else {
		// parse id from url
		id, err = strconv.Atoi(resource)
		if err != nil {
			err = models.BadRequestError(ctx, err)
			return
		}
	}
	user, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderUser(ctx, w, user)
	return
}

// create a user
// POST /user
func (h *UserHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	register := prototypes.Register{}
	err = json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		err = models.ServerError(ctx, err)
		return
	}
	user := prototypes.User{
		Email:     register.Email,
		HashedPwd: string(hash),
	}
	err = h.Model.Create(ctx, &user)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderUser(ctx, w, user)
	return
}

// update user information
// Patch /user
func (h *UserHandler) handlePatch(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	user, err := h.Model.Fetch(ctx, *payload.Pub.ID)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	err = h.Model.Update(ctx, &user)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderUser(ctx, w, user)
	return
}

// delete current user
// DELETE /user
func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	payload, ok := ctx.Value(prototypes.UserKey).(prototypes.JwtPayload)
	if !ok {
		err = models.AuthorizationError(ctx, err)
		return
	}
	auth := prototypes.Auth{}
	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	if payload.Pub.Email != auth.Email {
		err = models.AuthorizationError(ctx, err)
		return
	}
	err = h.Model.Delete(ctx, &auth)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	views.RenderEmpty(ctx, w)
	return
}

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
		// parse id from url
		id, err = strconv.Atoi(resource)
		if err != nil {
			err = models.BadRequestError(ctx, err)
			return
		}
	}
	posts, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderPostList(ctx, w, posts)
	return
}
