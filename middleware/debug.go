package middleware

import (
	"net/http"
	"net/http/httputil"

	"github.com/spyzhov/chttp"
)

// Debug is a constructor for Debug, that provides default transport
func Debug(active bool, logger Logger, options ...DebugOption) chttp.Middleware {
	cfg := &debugConfig{
		DumpRequestBody:  true,
		DumpResponseBody: false,
	}
	for _, option := range options {
		option(cfg)
	}
	return func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		if !active {
			return next(request)
		}
		var path string
		if request.URL != nil {
			path = request.URL.Path
		}
		log := getLogger(logger).
			WithContext(request.Context()).
			WithField("Method", request.Method).
			WithField("Host", request.Host).
			WithField("Path", path)

		data, err := httputil.DumpRequest(request, cfg.DumpRequestBody)
		if err != nil {
			log.WithField("error", err).Printf("error dumping request")
		} else {
			log.WithField("request", data).Printf("dump request")
		}

		response, err := next(request)
		if err != nil {
			log.WithField("error", err).Printf("error requesting")
		}
		if response != nil {
			data, err = httputil.DumpResponse(response, cfg.DumpResponseBody)
			if err != nil {
				log.WithField("error", err).Printf("error dumping response")
			} else {
				log.WithField("response", data).Printf("dump response")
			}
		}
		return response, err
	}
}

type debugConfig struct {
	DumpRequestBody  bool
	DumpResponseBody bool
}

type DebugOption func(c *debugConfig)

func WithDebugDumpRequestBody(dumpRequestBody bool) DebugOption {
	return func(c *debugConfig) {
		c.DumpRequestBody = dumpRequestBody
	}
}

func WithDebugDumpResponseBody(dumpResponseBody bool) DebugOption {
	return func(c *debugConfig) {
		c.DumpResponseBody = dumpResponseBody
	}
}
