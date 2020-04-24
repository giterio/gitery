package views

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// ErrorView is a custom error
type ErrorView struct {
	Ok          bool      `json:"ok"`
	StatusCode  int       `json:"status_code"`
	ErrorCode   int       `json:"error_code"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	trace       string
}

func createError(ctx context.Context, statusCode int, errorCode int, description string, err error) ErrorView {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("[ERROR %d] %s\n%s:%d %s", errorCode, description, file, line, funcName)
	if err != nil {
		trace = trace + "\n" + err.Error()
	}

	return ErrorView{
		Ok:          false,
		StatusCode:  statusCode,
		ErrorCode:   errorCode,
		Description: description,
		Timestamp:   time.Now(),
		trace:       trace,
	}
}
