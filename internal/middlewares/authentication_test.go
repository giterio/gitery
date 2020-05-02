package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitery/internal/controllers"
	"gitery/internal/views"
	"gitery/test/testdata"
)

func TestAuthentication(t *testing.T) {
	// Payload structure:
	// {
	//   "iat": 1516239022,
	//   "pub": {
	//      "user_id": 1,
	//      "email": "Murphy@jwt.com"
	//   }
	// }
	fakeToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsInB1YiI6eyJ1c2VyX2lkIjoxLCJlbWFpbCI6Ik11cnBoeUBqd3QuY29tIn19.C0MqWff5aL5jxUA_kVNZL8Wh9N_zsPFjBnhcNMtKi2g"
	mux := http.NewServeMux()
	router := &controllers.RootHandler{
		AuthHandler: &controllers.AuthHandler{Model: &testdata.FakeAuthService{Token: fakeToken, JwtSecret: "this is screct"}},
	}
	mux.Handle("/", Authentication(router))

	writer := httptest.NewRecorder()
	jsonStr := strings.NewReader(`{"email":"Murphy@jwt.com","password":"fakePassword"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth", jsonStr)
	request.Header.Add("Authorization", "Bearer "+fakeToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response body is %v", string(writer.Body.Bytes()))
	}
	var dataView views.DataView
	err := json.Unmarshal(writer.Body.Bytes(), &dataView)
	if err != nil {
		t.Errorf("response body not parsable, error: %s", err.Error())
	}
	if postData, ok := dataView.Data.(map[string]interface{}); !ok {
		t.Errorf("Data structure is not expected")
	} else if postData["token"] != fakeToken {
		t.Errorf("Cannot retrieve JSON post")
	}
}
