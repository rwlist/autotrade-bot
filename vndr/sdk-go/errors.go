package chatexsdk

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrInternalServer      = errors.New("internal error")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrTooManyRequests     = errors.New("too many requests")
)

type ValidationError struct {
	msg    string
	errors map[string]interface{}
}

func NewValidationError(msg string, errors map[string]interface{}) ValidationError {
	return ValidationError{
		msg:    msg,
		errors: errors,
	}
}

func (e ValidationError) Error() string {
	return e.msg
}
