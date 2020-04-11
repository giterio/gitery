package router

import (
	"fmt"
	"net/http"
)

// RootHandler ...
type RootHandler struct{}

func (RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO switch logic based on r.URL.Path
	// TODO default -> post handler
	fmt.Fprintf(w, "Hello!")
}
