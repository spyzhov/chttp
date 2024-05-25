package chttp

import (
	"net/http"
)

// Middleware is an extended interface to the RoundTrip function of the http.RoundTripper interface.
type Middleware func(
	request *http.Request,
	next func(request *http.Request) (*http.Response, error),
) (*http.Response, error)
