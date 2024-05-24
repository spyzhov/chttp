package chttp

import (
	"net/http"
)

func (c *Client) transport() http.RoundTripper {
	base := c.base
	if base == nil {
		base = http.DefaultTransport
	}
	return &transport{
		Client:  c,
		Default: base,
	}
}

type transport struct {
	Client  *Client
	Default http.RoundTripper
}

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	middlewares := t.Client.getMiddlewares()
	var next func(request *http.Request) (*http.Response, error)
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
	return func(request *http.Request, _ func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		return t.Default.RoundTrip(request)
	}
}
