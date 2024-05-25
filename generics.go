package chttp

import (
	"context"
	"net/http"
)

// GenericJSONClient is an HTTP client wrapper around the Client
// with automated marshaling of the request body and unmarshaling response body into a known structure.
type GenericJSONClient[Result any] struct {
	JSONClient
}

// NewGenericJSONClient wraps JSONClient with the GenericJSONClient.
func NewGenericJSONClient[Result any](json *JSONClient) *GenericJSONClient[Result] {
	if json == nil {
		json = NewJSON(nil)
	}
	return &GenericJSONClient[Result]{
		JSONClient: *json,
	}
}

// Request prepares the request by marshaling request body and tries to unmarshal response
// with the JSONClient.UnmarshalHTTPResponse function.
func (c *GenericJSONClient[Result]) Request(
	ctx context.Context,
	method string,
	url string,
	body interface{},
) (result Result, err error) {
	err = c.JSONClient.Request(ctx, method, url, body, &result)
	return result, err
}

// UnmarshalHTTPResponse tries to unmarshal the response body into the given result interface.
// Result should be reference type and not nil.
func (c *GenericJSONClient[Result]) UnmarshalHTTPResponse(
	response *http.Response,
	httpErr error,
) (result Result, err error) {
	return result, c.JSONClient.UnmarshalHTTPResponse(response, httpErr, &result)
}

// Method returns a function implementation of the HTTP Method from the Client by its name.
// Returns Client.GET as the default method.
func (c *GenericJSONClient[Result]) Method(
	method string,
) func(ctx context.Context, url string, body interface{}) (Result, error) {
	switch method {
	case http.MethodHead:
		return c.HEAD
	case http.MethodPost:
		return c.POST
	case http.MethodPut:
		return c.PUT
	case http.MethodPatch:
		return c.PATCH
	case http.MethodDelete:
		return c.DELETE
	case http.MethodConnect:
		return c.CONNECT
	case http.MethodOptions:
		return c.OPTIONS
	case http.MethodTrace:
		return c.TRACE
	case http.MethodGet:
		fallthrough
	default:
		return c.GET
	}
}

// GET is an alias to do the Request with the http.MethodGet method.
func (c *GenericJSONClient[Result]) GET(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodGet, url, body)
}

// HEAD is an alias to do the Request with the http.MethodHead method.
func (c *GenericJSONClient[Result]) HEAD(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodHead, url, body)
}

// POST is an alias to do the Request with the http.MethodPost method.
func (c *GenericJSONClient[Result]) POST(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodPost, url, body)
}

// PUT is an alias to do the Request with the http.MethodPut method.
func (c *GenericJSONClient[Result]) PUT(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodPut, url, body)
}

// PATCH is an alias to do the Request with the http.MethodPatch method.
func (c *GenericJSONClient[Result]) PATCH(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodPatch, url, body)
}

// DELETE is an alias to do the Request with the http.MethodDelete method.
func (c *GenericJSONClient[Result]) DELETE(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodDelete, url, body)
}

// CONNECT is an alias to do the Request with the http.MethodConnect method.
func (c *GenericJSONClient[Result]) CONNECT(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodConnect, url, body)
}

// OPTIONS is an alias to do the Request with the http.MethodOptions method.
func (c *GenericJSONClient[Result]) OPTIONS(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodOptions, url, body)
}

// TRACE is an alias to do the Request with the http.MethodTrace method.
func (c *GenericJSONClient[Result]) TRACE(ctx context.Context, url string, body interface{}) (Result, error) {
	return c.Request(ctx, http.MethodTrace, url, body)
}

// Clone will clone an instance of the GenericJSONClient without references to the old one.
func (c *GenericJSONClient[Result]) Clone() *GenericJSONClient[Result] {
	return NewGenericJSONClient[Result](JSON(c.Client.Clone()))
}
