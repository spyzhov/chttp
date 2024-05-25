package middleware

import (
	"net/http"
	"time"

	"github.com/spyzhov/chttp"
)

// Trace middleware adds short logs on each request.
func Trace(logger Logger) chttp.Middleware {
	return func(
		request *http.Request,
		next func(request *http.Request) (*http.Response, error),
	) (response *http.Response, err error) {
		defer func(start time.Time) {
			var path string
			if request.URL != nil {
				path = request.URL.Path
			}
			log := getLogger(logger).
				WithContext(request.Context()).
				WithField("method", request.Method).
				WithField("host", request.Host).
				WithField("path", path).
				WithField("request_time", time.Since(start))

			if response != nil {
				log = log.
					WithField("status_code", response.StatusCode).
					WithField("content_length", response.ContentLength)
			}
			if err != nil {
				log = log.WithField("error", err)
			}
			log.Printf("http client response")
		}(time.Now())

		return next(request)
	}
}
