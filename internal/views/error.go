package views

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"gitery/internal/models"
)

// ErrorView is the response data structure for a custom error
type ErrorView struct {
	models.Error       // using anonymous member to get flat json structure
	Ok           bool  `json:"ok"`
	Timestamp    int64 `json:"timestamp"`
}

// RenderError writes the ErrorView response to http connection
func RenderError(ctx context.Context, w http.ResponseWriter, e models.Error) {
	errorView := ErrorView{
		Error:     e,
		Ok:        false,
		Timestamp: e.Timestamp.Unix(),
	}
	output, err := json.MarshalIndent(errorView, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
		return
	}

	w.WriteHeader(e.StatusCode)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}
}
