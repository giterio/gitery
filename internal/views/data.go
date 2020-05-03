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
	Timestamp int64       `json:"timestamp"`
}

// Render ...
func Render(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	dataView := DataView{
		Ok:        true,
		Timestamp: time.Now().Unix(),
	}
	if data != nil {
		dataView.Data = data
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
