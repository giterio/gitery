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
	switch r.Method {
	case http.MethodPost:
		err = h.handlePost(w, r)
	}
	if err != nil {
		ctx := r.Context()
		e := models.ServerError(ctx, err)
		views.RenderError(ctx, w, e)
	}
}

// handle login event
func (h *AuthHandler) handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	auth := prototypes.Auth{}
	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		err = models.BadRequestError(ctx, err)
		return
	}
	token, err := h.Model.Login(ctx, auth)
	if err != nil {
		return
	}
	err = views.RenderAuth(ctx, w, token)
	return
}
