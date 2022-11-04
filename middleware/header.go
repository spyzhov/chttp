package transport

import (
	"net/http"

	"github.com/spyzhov/chttp"
)

// Headers is a chttp.Middleware constructor to add static headers to request.
func Headers(headers map[string]string) chttp.Middleware {
	return func(request *http.Request, next chttp.RoundTripper) (*http.Response, error) {
		for name, value := range headers {
			request.Header.Set(name, value)
		}
		return next(request)
	}
}
