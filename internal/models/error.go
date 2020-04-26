package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// Error represents an error with description and trace
type Error struct {
	StatusCode  int               `json:"status_code"`
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

func createError(ctx context.Context, statusCode int, errorCode int, description string, err error) Error {
	// bubble the err if it has been already formatted
	if e, ok := err.(Error); ok {
		return e
	}

	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
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
		StatusCode:  statusCode,
		ErrorCode:   errorCode,
		Description: description,
		Timestamp:   time.Now(),
		Trace:       trace,
	}
}

// ServerError means some server error are occurred.
func ServerError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, http.StatusInternalServerError, description, err)
}

// BadRequestError means the request URL is illegal or request body cannot be parsed
func BadRequestError(ctx context.Context) Error {
	description := "The request URL is illegal or request body cannot be parsed"
	return createError(ctx, http.StatusAccepted, http.StatusBadRequest, description, nil)
}

// AuthorizationError return 401 for unauthorized request
func AuthorizationError(ctx context.Context) Error {
	description := "Unauthorized, maybe invalid token."
	return createError(ctx, http.StatusAccepted, http.StatusUnauthorized, description, nil)
}

// ForbiddenError return 403 for unauthorized request
func ForbiddenError(ctx context.Context) Error {
	description := http.StatusText(http.StatusForbidden)
	return createError(ctx, http.StatusAccepted, http.StatusForbidden, description, nil)
}

// NotFoundError means resource is not found.
func NotFoundError(ctx context.Context) Error {
	description := http.StatusText(http.StatusNotFound)
	return createError(ctx, http.StatusAccepted, http.StatusNotFound, description, nil)
}

// TransactionError means there is something wrong on database.
func TransactionError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, 10001, description, err)
}

// IdentityNonExistError means email or username is not existent.
func IdentityNonExistError(ctx context.Context) Error {
	description := "Email or Username is not exist."
	return createError(ctx, http.StatusAccepted, 10011, description, nil)
}

// InvalidPasswordError means the password is invalid.
func InvalidPasswordError(ctx context.Context) Error {
	description := "Password invalid."
	return createError(ctx, http.StatusAccepted, 10012, description, nil)
}

// PasswordTooSimpleError means the password is too simple.
func PasswordTooSimpleError(ctx context.Context) Error {
	description := "Password too simple, at least 8 characters required."
	return createError(ctx, http.StatusAccepted, 10013, description, nil)
}
