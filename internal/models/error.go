package models

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// Error represents an error with description and trace
type Error struct {
	ErrorCode   int       `json:"error_code"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	trace       string
}

func createError(ctx context.Context, errorCode int, description string, err error) Error {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("[ERROR %d] %s\n%s:%d %s", errorCode, description, file, line, funcName)
	if err != nil {
		trace = trace + "\n" + err.Error()
	}

	return Error{
		ErrorCode:   errorCode,
		Description: description,
		Timestamp:   time.Now(),
		trace:       trace,
	}
}
