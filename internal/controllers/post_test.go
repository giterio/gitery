package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitery/internal/domains"
	"gitery/test/testdata"
)

func TestGetPost(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", &Router{
		PostHandler: &PostHandler{Model: &testdata.FakePostService{}},
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var post domains.Post
	json.Unmarshal(writer.Body.Bytes(), &post)
	if *post.ID != 1 {
		t.Errorf("Cannot retrieve JSON post")
	}
}

// func TestPutPost(t *testing.T) {
// 	mux := http.NewServeMux()
// 	post := &FakePost{}
// 	mux.HandleFunc("/post/", handleRequest(post))

// 	writer := httptest.NewRecorder()
// 	json := strings.NewReader(`{"content":"Updated post","author":"Sau Sheong"}`)
// 	request, _ := http.NewRequest("PUT", "/post/1", json)
// 	mux.ServeHTTP(writer, request)

// 	if writer.Code != 200 {
// 		t.Error("Response code is %v", writer.Code)
// 	}

// 	if post.Content != "Updated post" {
// 		t.Error("Content is not correct", post.Content)
// 	}
// }
