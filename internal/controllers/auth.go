package controllers

import (
	"encoding/json"
	"net/http"

	"gitery/internal/models"
	"gitery/internal/prototypes"
	"gitery/internal/views"
)

// AuthHandler ...
type AuthHandler struct {
	Model prototypes.AuthService
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	switch r.Method {
	// handle login event
	case http.MethodPost:
		err = h.handlePost(w, r)
	default:
		err = models.ForbiddenError(ctx, nil)
	}
	if err != nil {
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// handle login event
func (h *AuthHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	login := prototypes.Login{}
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	token, user, err := h.Model.Login(ctx, login)
	if err != nil {
		return
	}
	err = views.RenderAuth(ctx, w, token, user)
	return
}
