package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

const (
	defaultErrorField = "body"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrRequestBody = errors.New("invalid request body")
	ErrInternal    = errors.New("internal server error")
)

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	e := toError(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// TODO: waht should do I with the error
	jsonResponse(w, e, e.Code)
}

type Error struct {
	Code   int                `json:"-"`
	Errors map[string][]error `json:"errors"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	m := make(map[string][]string)
	for field, errs := range e.Errors {
		for _, err := range errs {
			m[field] = append(m[field], err.Error())
		}
	}

	return json.Marshal(m)
}

func (e Error) Error() string {
	return ""
}

func toError(err error) (e Error) {
	switch t := err.(type) {
	case Error:
		e = t
	case validation.Errors:
		e = newValidationError(t)
	default:
		e = newInternalError(ErrInternal)
	}
	return
}

func newError(code int, err error) Error {
	return Error{
		Code: code,
		Errors: map[string][]error{
			defaultErrorField: {err},
		},
	}
}

func newInternalError(err error) Error {
	return newError(http.StatusInternalServerError, err)
}

func newValidationError(errs validation.Errors) (e Error) {
	e.Code = http.StatusUnprocessableEntity
	e.Errors = make(map[string][]error)
	for field, err := range errs {
		e.Errors[field] = []error{err}
	}

	return e
}
