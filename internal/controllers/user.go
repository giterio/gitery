package controllers

import (
	"encoding/json"
	"net/http"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// UserHandler ...
type UserHandler struct {
	Model prototypes.UserService
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		err = h.handleGet(w, r)
	case http.MethodPost:
		err = h.handlePost(w, r)
	case http.MethodPatch:
		err = h.handlePatch(w, r)
	case http.MethodDelete:
		err = h.handleDelete(w, r)
	}
	if err != nil {
		ctx := r.Context()
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// Retrieve a user
// GET /user/1
func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	token := r.Header.Get("Authorization")
	ctx := r.Context()
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	user, err := h.Model.Fetch(ctx, token)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderUser(ctx, w, user)
	return
}

// create a user
func (h *UserHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	user := prototypes.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
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
func (h *UserHandler) handlePatch(w http.ResponseWriter, r *http.Request) (err error) {
	token := r.Header.Get("Authorization")
	ctx := r.Context()
	user, err := h.Model.Fetch(ctx, token)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err = models.BadRequestError(ctx)
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

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	user := prototypes.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	err = h.Model.Delete(ctx, &user)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	views.RenderEmpty(ctx, w)
	return
}
