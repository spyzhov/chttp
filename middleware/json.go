package middleware

import (
	"github.com/spyzhov/chttp"
)

// JSON is a middleware that adds a `Content-Type` and `Accept` headers with the `application/json` value.
func JSON() chttp.Middleware {
	return Headers(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}, false)
}
