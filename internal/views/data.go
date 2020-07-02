package views

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// DataView is common response structure with specific view as data
type DataView struct {
	Data      interface{} `json:"data,omitempty"`
	Ok        bool        `json:"ok"`
	Timestamp int64       `json:"timestamp"`
}

// Render compose response data structure and write the data to http(s) connection
func Render(ctx context.Context, w http.ResponseWriter, data interface{}) (err error) {
	dataView := DataView{
		Ok:        true,
		Timestamp: time.Now().Unix(),
	}
	if data != nil {
		dataView.Data = data
	}
	output, err := json.MarshalIndent(dataView, "", "\t")
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
