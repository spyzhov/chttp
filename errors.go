package chttp

import (
	"fmt"
	"net/http"
)

type Error struct {
	Response *http.Response
	Body     []byte
	Base     error
}

func NewError(response *http.Response, body []byte, err error) *Error {
	return &Error{
		Response: response,
		Body:     body,
		Base:     err,
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
	if e.Response != nil {
		return fmt.Errorf("http error, status_code=%d", e.Response.StatusCode)
	}
	return fmt.Errorf("unknown error")
}
