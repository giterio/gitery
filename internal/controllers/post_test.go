package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitery/internal/views"
	"gitery/test/testdata"
)

func TestHandleGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", &Root{
		PostHandler: &PostHandler{Model: &testdata.FakePostService{}},
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var dataView views.DataView
	err := json.Unmarshal(writer.Body.Bytes(), &dataView)
	if err != nil {
		t.Errorf("response body not parsable %s", err.Error())
	}
	postData, ok := dataView.Data.(map[string]interface{})
	if ok && postData["id"] == 1 {
		t.Errorf("Cannot retrieve JSON post")
	}
}

func TestHandlePost(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", &Root{
		PostHandler: &PostHandler{Model: &testdata.FakePostService{}},
	})

	writer := httptest.NewRecorder()
	jsonStr := strings.NewReader(`{"content":"Updated post","user_id": 1}`)
	request, _ := http.NewRequest("POST", "/post/1", jsonStr)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var dataView views.DataView
	err := json.Unmarshal(writer.Body.Bytes(), &dataView)
	if err != nil {
		t.Errorf("response body not parsable %s", err.Error())
	}
	postData, ok := dataView.Data.(map[string]interface{})
	if ok && postData["content"] == "Updated post" && postData["user_id"] == 1 {
		t.Errorf("Cannot retrieve JSON post")
	}
}
