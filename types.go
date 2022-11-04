package chttp

import (
	"net/http"
)

type RoundTripper func(request *http.Request) (*http.Response, error)

type Middleware func(request *http.Request, next RoundTripper) (*http.Response, error)
