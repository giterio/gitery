package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// Error represents an error with description and trace
type Error struct {
	StatusCode  int               `json:"statusCode"`
	ErrorCode   int               `json:"errorCode"`
	Description string            `json:"description"`
	Trace       map[string]string `json:"trace"`
	Timestamp   time.Time         `json:"-"`
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
func BadRequestError(ctx context.Context, err error) Error {
	description := "The request URL is illegal or request body cannot be parsed"
	return createError(ctx, http.StatusAccepted, http.StatusBadRequest, description, err)
}

// AuthorizationError return 401 for unauthorized request
func AuthorizationError(ctx context.Context, err error) Error {
	description := "Unauthorized, maybe invalid token."
	return createError(ctx, http.StatusAccepted, http.StatusUnauthorized, description, err)
}

// ForbiddenError return 403 for unauthorized request
func ForbiddenError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusForbidden)
	return createError(ctx, http.StatusAccepted, http.StatusForbidden, description, err)
}

// NotFoundError means resource is not found.
func NotFoundError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusNotFound)
	return createError(ctx, http.StatusAccepted, http.StatusNotFound, description, err)
}

// ConflictError means resource is already exist.
func ConflictError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusConflict)
	return createError(ctx, http.StatusAccepted, http.StatusConflict, description, err)
}

// TransactionError means there is something wrong on database.
func TransactionError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, 10001, description, err)
}

// IdentityNonExistError means email or username is not existent.
func IdentityNonExistError(ctx context.Context, err error) Error {
	description := "Email or Username is not exist."
	return createError(ctx, http.StatusAccepted, 10011, description, err)
}

// InvalidPasswordError means the password is invalid.
func InvalidPasswordError(ctx context.Context, err error) Error {
	description := "Password invalid."
	return createError(ctx, http.StatusAccepted, 10012, description, err)
}

// Following are logic errors

// IllegalEmailFormatError means the email format is incorrect.
func IllegalEmailFormatError(ctx context.Context) Error {
	description := "Illegal email format"
	return createError(ctx, http.StatusAccepted, 10013, description, nil)
}

// IncorrectPasswordFormatError means the password is too simple.
func IncorrectPasswordFormatError(ctx context.Context) Error {
	description := "The password is a combination of uppercase and lowercase letters and numbers with a length of 8-32"
	return createError(ctx, http.StatusAccepted, 10013, description, nil)
}

// HandleDatabaseQueryError handles DB transaction error
func HandleDatabaseQueryError(ctx context.Context, err error) Error {
	if err == sql.ErrNoRows {
		return NotFoundError(ctx, err)
	}
	return TransactionError(ctx, err)
}
