package gokit_realworld

import (
	"errors"
	"fmt"
)

// Application error codes.
const (
	// Action cannot be performed.
	EConflict = "conflict"
	// Internal error.
	EInternal = "internal"
	// Entity does not exist.
	ENotFound = "not_found"
	// Too many API requests.
	ERateLimit = "rate_limit"
	// User ID validation failed.
	EInvalidUserID = "invalid_user_id"
	// Username validation failed.
	EInvalidUsername   = "invalid_username"
	EIncorrectPassword = "incorrect_password"
)

type Error struct {
	Code string
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Code, e.Err)
}

func (e Error) Unwrap() error {
	return e.Err
}

// ErrorCode returns the code of the error, if available.
func ErrorCode(err error) string {
	var e Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

func InternalError(err error) error {
	return Error{
		Code: EInternal,
		Err:  fmt.Errorf("internal error: %w", err),
	}
}
