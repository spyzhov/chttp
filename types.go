package chttp

import (
	"net/http"
)

// RoundTripper is a RoundTrip function implementation of the http.RoundTripper interface.
type RoundTripper func(request *http.Request) (*http.Response, error)

// Middleware is an extended interface to the RoundTrip function of the http.RoundTripper interface.
type Middleware func(request *http.Request, next RoundTripper) (*http.Response, error)
