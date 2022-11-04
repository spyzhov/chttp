package transport

import (
	"net/http"
	"time"

	"github.com/spyzhov/chttp"
)

type TraceMiddleware struct {
	Name   string
	Logger Logger
}

func Trace(name string, logger Logger) chttp.Middleware {
	return (&TraceMiddleware{Name: name, Logger: logger}).Middleware
}

func (s *TraceMiddleware) Middleware(request *http.Request, next chttp.RoundTripper) (response *http.Response, err error) {
	defer func(start time.Time) {
		logger := getLogger(s.Logger).
			WithContext(request.Context()).
			WithField("name", s.Name).
			WithField("method", request.Method).
			WithField("host", request.Host).
			WithField("path", request.RequestURI).
			WithField("request_time", time.Since(start))

		if response != nil {
			logger = logger.
				WithField("status_code", response.StatusCode).
				WithField("content_length", response.ContentLength)
		}
		if err != nil {
			logger = logger.WithField("error", err)
		}
		logger.Printf("http client response")
	}(time.Now())

	return next(request)
}
