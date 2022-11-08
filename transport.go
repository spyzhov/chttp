package chttp

import (
	"net/http"
)

func (c *Client) transport(Default http.RoundTripper) http.RoundTripper {
	if Default == nil {
		Default = http.DefaultTransport
	}
	return &transport{
		Client:  c,
		Default: Default.RoundTrip,
	}
}

type transport struct {
	Client  *Client
	Default RoundTripper
}

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	middlewares := t.Client.getMiddlewares()
	var next RoundTripper
	next = func(request *http.Request) (*http.Response, error) {
		var middleware Middleware
		if len(middlewares) == 0 {
			middleware = t.getDefault()
		} else {
			middleware = middlewares[0]
			middlewares = middlewares[1:]
		}
		return middleware(request, next)
	}
	return next(request)
}

func (t *transport) getDefault() Middleware {
	return func(request *http.Request, _ RoundTripper) (*http.Response, error) {
		return t.Default(request)
	}
}
