package middleware

import (
	"net/http"

	"github.com/spyzhov/chttp"
)

// Headers is a chttp.Middleware constructor to add static headers to request.
func Headers(headers map[string]string, force bool) chttp.Middleware {
	return func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		for name, value := range headers {
			if force || request.Header.Get(name) == "" {
				request.Header.Set(name, value)
			}
		}
		return next(request)
	}
}
