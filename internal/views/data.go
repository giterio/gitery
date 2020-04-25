package views

import (
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
func Render(w http.ResponseWriter, data interface{}) (err error) {
	dataView := DataView{
		Data:      data,
		Ok:        true,
		Timestamp: time.Now(),
	}
	output, err := json.MarshalIndent(dataView, "", "\t\t")
	if err != nil {
		return
	}
	_, err = w.Write(output)
	return
}
