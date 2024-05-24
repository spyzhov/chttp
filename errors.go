package chttp

import (
	"fmt"
	"net/http"
)

var (
	ErrStatusCode = fmt.Errorf("wrong status code")
	ErrUnknown    = fmt.Errorf("unknown error")
)

type Error struct {
	Response *http.Response
	Body     []byte
	Base     error
}

func newError(response *http.Response, body []byte, err error) *Error {
	return &Error{
		Response: response,
		Body:     body,
		Base:     err,
	}
}

func (e *Error) IsStatusCode() bool {
	if e == nil {
		return false
	}
	if e.Base != nil {
		return false
	}
	if e.Response != nil {
		return e.Response.StatusCode >= http.StatusMultipleChoices
	}
	return false
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
	if e.Response != nil {
		return fmt.Errorf("%w, status_code=%d", ErrStatusCode, e.Response.StatusCode)
	}
	return ErrUnknown
}
