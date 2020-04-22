package views

import (
	"time"
)

// ErrorView is a custom error
type ErrorView struct {
	StatusCode  int       `json:"status_code"`
	ErrorCode   int       `json:"error_code"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	trace       string
}
