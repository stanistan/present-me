package errors

import (
	"fmt"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
)

// Error will represent an error/error Code.
type Error struct {
	Msg      string
	Cause    error
	HttpCode int
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s - %s", e.Msg, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

var WithStack = errors.WithStack

func WrapGithubErr(e error, msg string, args ...interface{}) error {
	if e == nil {
		return nil
	}

	err := &Error{
		Msg:      fmt.Sprintf(msg, args...),
		Cause:    e,
		HttpCode: 500,
	}

	var errorResponse *github.ErrorResponse
	if errors.As(e, &errorResponse) {
		err.HttpCode = errorResponse.Response.StatusCode
	}

	return err
}

func WrapErr(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return &Error{
		Msg:      "Unexpected Error",
		Cause:    err,
		HttpCode: 500,
	}
}

func Is404(err error) bool {
	return IsHttpCode(err, 404)
}

func IsHttpCode(err error, code int) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*Error); ok {
		return e.HttpCode == code
	}

	return false
}

func New(msg string) error {
	return errors.New(msg)
}
