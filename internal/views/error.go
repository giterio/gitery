package views

import (
	"encoding/json"
	"gitery/internal/models"
	"net/http"
	"time"
)

// ErrorView is a custom error
type ErrorView struct {
	models.Error           // using anonymous member to get flat json structure
	Ok           bool      `json:"ok"`
	Timestamp    time.Time `json:"timestamp"`
}

// RenderError ...
func RenderError(w http.ResponseWriter, e models.Error) {
	errorView := ErrorView{
		Error:     e,
		Ok:        false,
		Timestamp: time.Now(),
	}
	output, err := json.MarshalIndent(errorView, "", "\t\t")
	if err != nil {
		return
	}
	_, err = w.Write(output)
	return
}
