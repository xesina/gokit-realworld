package error

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-ozzo/ozzo-validation/v4"
	realworld "github.com/xesina/go-kit-realworld-example-app"
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
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("EncodeError with nil error")
	}
	e := toError(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
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
	switch {
	case errors.As(err, &Error{}):
		e = Error{}
		errors.As(err, &e)

	case errors.As(err, &validation.Errors{}):
		temp := validation.Errors{}
		errors.As(err, &temp)
		e = newValidationError(temp)

	case errors.As(err, &realworld.Error{}):
		temp := realworld.Error{}
		errors.As(err, &temp)
		e = newDomainError(temp)

	default:
		e = newInternalError(ErrInternal)
	}

	return
}

func NewError(code int, err error) Error {
	return Error{
		Code: code,
		Errors: map[string][]error{
			defaultErrorField: {err},
		},
	}
}

func newInternalError(err error) Error {
	return NewError(http.StatusInternalServerError, err)
}

func newValidationError(errs validation.Errors) (e Error) {
	e.Code = http.StatusUnprocessableEntity
	e.Errors = make(map[string][]error)
	for field, err := range errs {
		e.Errors[field] = []error{err}
	}

	return e
}

func newDomainError(err realworld.Error) (e Error) {
	e.Code = mapDomainErrorCode(err.Code)
	e.Errors = make(map[string][]error)
	e.Errors[defaultErrorField] = []error{err.Err}
	return e
}

func mapDomainErrorCode(code string) int {
	switch code {
	case realworld.EIncorrectPassword:
		return http.StatusForbidden
	case realworld.EConflict:
		return http.StatusUnprocessableEntity
	case realworld.ENotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
