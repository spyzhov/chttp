package chttp

import (
	"fmt"
)

type Error struct {
	StatusCode int
	Body       []byte
	Base       error
}

func NewError(code int, body []byte, err error) *Error {
	return &Error{
		StatusCode: code,
		Body:       body,
		Base:       err,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Unwrap().Error()
}

func (e *Error) Unwrap() error {
	if e.Base != nil {
		return e.Base
	}
	return fmt.Errorf("http error, status_code=%d", e.StatusCode)
}
