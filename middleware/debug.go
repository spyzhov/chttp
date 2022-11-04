package transport

import (
	"net/http"
	"net/http/httputil"

	"github.com/spyzhov/chttp"
)

// DebugMiddleware implements the http.RoundTripper interface.
// Describe the request in DebugMiddleware mode.
type DebugMiddleware struct {
	Debug  bool
	Logger Logger
}

// Debug is a constructor for Debug, that provides default transport
func Debug(active bool, logger Logger) chttp.Middleware {
	return (&DebugMiddleware{Debug: active, Logger: logger}).Middleware
}

func (s *DebugMiddleware) Middleware(request *http.Request, next chttp.RoundTripper) (*http.Response, error) {
	if !s.Debug {
		return next(request)
	}
	logger := getLogger(s.Logger).
		WithContext(request.Context()).
		WithField("Method", request.Method).
		WithField("Host", request.Host).
		WithField("RequestURI", request.RequestURI)

	data, err := httputil.DumpRequest(request, true)
	if err != nil {
		logger.WithField("error", err).Printf("error dumping request")
	} else {
		logger.WithField("request", data).Printf("dump request")
	}

	response, err := next(request)
	if err != nil {
		logger.WithField("error", err).Printf("error requesting")
	} else {
		data, err = httputil.DumpResponse(response, false)
		if err != nil {
			logger.WithField("error", err).Printf("error dumping response")
		} else {
			logger.WithField("response", data).Printf("dump response")
		}
	}
	return response, err
}
