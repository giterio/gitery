package controllers

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
func (h *PostHandler) handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	resource, _ := models.ShiftRoute(r)
	ctx := r.Context()
	id, err := strconv.Atoi(resource)
	if err != nil {
		err = models.BadRequestError(ctx)
		return
	}
	post, err := h.Model.Fetch(ctx, id)
	if err != nil {
		err = models.TransactionError(ctx, err)
		return
	}
	err = views.RenderPost(ctx, w, post)
	return
}

