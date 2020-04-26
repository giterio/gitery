package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"
)

// Error represents an error with description and trace
type Error struct {
	ErrorCode   int               `json:"error_code"`
	Description string            `json:"description"`
	Timestamp   time.Time         `json:"timestamp"`
	Trace       map[string]string `json:"trace"`
}

// Error ...
func (sessionError Error) Error() string {
	str, err := json.MarshalIndent(sessionError, "", "\t\t")
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func createError(ctx context.Context, errorCode int, description string, err error) Error {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	// trace := fmt.Sprintf("[ERROR %d] %s\n%s:%d %s", errorCode, description, file, line, funcName)
	trace := map[string]string{
		"name":        fmt.Sprintf("[ERROR %d]", errorCode),
		"description": description,
		"file":        fmt.Sprintf("%s:%d", file, line),
		"function":    funcName,
	}
	if err != nil {
		trace["error"] = err.Error()
	}

	return Error{
		ErrorCode:   errorCode,
		Description: description,
		Timestamp:   time.Now(),
		Trace:       trace,
	}
}

// BubbleError ...
func BubbleError(ctx context.Context, errorCode int, description string, err error) (e Error) {
	if e, ok := err.(Error); ok {
		return e
	}
	e = createError(ctx, errorCode, description, err)
	return e
}
