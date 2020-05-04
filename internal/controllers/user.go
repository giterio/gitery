package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// UserHandler ...
type UserHandler struct {
	Model           prototypes.UserService
	UserPostHandler UserPostHandler
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	// user is the resource to manipulate
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

// Retrieve a user
// GET /user
func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
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
	auth := prototypes.Auth{}
	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		err = models.BadRequestError(ctx, err)
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
	return
}
