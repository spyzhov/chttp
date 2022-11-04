package transport

import (
	"github.com/spyzhov/chttp"
)

func JSON() chttp.Middleware {
	return Headers(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
}
