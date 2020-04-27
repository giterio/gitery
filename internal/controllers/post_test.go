package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitery/internal/prototypes"
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
	var post prototypes.Post
	json.Unmarshal(writer.Body.Bytes(), &post)
	if *post.ID != 1 {
		t.Errorf("Cannot retrieve JSON post")
	}
}

func TestHandlePost(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", &Root{
		PostHandler: &PostHandler{Model: &testdata.FakePostService{}},
	})

	writer := httptest.NewRecorder()
	jsonStr := strings.NewReader(`{"content":"Updated post","author":"Sau Sheong"}`)
	request, _ := http.NewRequest("POST", "/post/1", jsonStr)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var post prototypes.Post
	json.Unmarshal(writer.Body.Bytes(), &post)
	if post.Content != "Updated post" || post.Author != "Sau Sheong" {
		t.Errorf("Post not match, Content: %s, Author: %s", post.Content, post.Author)
	}
}
