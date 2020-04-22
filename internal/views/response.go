package views

import (
	"encoding/json"
	"net/http"
	"time"
)

// ResponseView ...
type ResponseView struct {
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// RenderData ...
func RenderData(w http.ResponseWriter, data interface{}) (err error) {
	respView := ResponseView{Data: data, Timestamp: time.Now()}
	output, err := json.MarshalIndent(respView, "", "\t\t")
	if err != nil {
		return
	}
	_, err = w.Write(output)
	return
}
