package views

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// DataView ...
type DataView struct {
	Data      interface{} `json:"data,omitempty"`
	Ok        bool        `json:"ok"`
	Timestamp time.Time   `json:"timestamp"`
}

// Render ...
func Render(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	if data == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	dataView := DataView{
		Data:      data,
		Ok:        true,
		Timestamp: time.Now(),
	}
	output, err := json.MarshalIndent(dataView, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
