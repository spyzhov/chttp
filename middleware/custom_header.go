package middleware

import (
	"net/http"

	"github.com/spyzhov/chttp"
)

// CustomHeaders is a chttp.Middleware constructor to add custom header based on the request to request.
func CustomHeaders(headers func(request *http.Request) map[string]string) chttp.Middleware {
	return func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		for name, value := range headers(request) {
			request.Header.Set(name, value)
		}
		return next(request)
	}
}
