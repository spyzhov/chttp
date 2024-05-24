package chttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONClient is an HTTP client wrapper around the Client
// with automated marshaling of the request body and unmarshaling response body.
type JSONClient struct {
	*Client
}

// JSON wraps Client with the JSONClient.
func JSON(client *Client) *JSONClient {
	return &JSONClient{
		Client: client,
	}
}

// NewJSON creates a JSONClient with new Client based on given http.Client.
func NewJSON(client *http.Client, options ...Option) *JSONClient {
	return &JSONClient{
		Client: NewClient(client, options...),
	}
}

// Request prepares the request by marshaling request body and tries to unmarshal response with the JSONClient.UnmarshalHTTPResponse function.
func (c *JSONClient) Request(ctx context.Context, method string, url string, body interface{}, result interface{}) (err error) {
	data, err := marshal(body)
	if err != nil {
		return err
	}
	res, err := c.Client.Request(ctx, method, url, data)
	return c.UnmarshalHTTPResponse(res, err, &result)
}

// UnmarshalHTTPResponse tries to unmarshal the response body into the given result interface.
// Result should be reference type and not nil.
func (c *JSONClient) UnmarshalHTTPResponse(response *http.Response, httpErr error, result interface{}) (err error) {
	if httpErr != nil {
		return newError(response, nil, fmt.Errorf("requesting error: %w", err))
	}
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return newError(response, data, fmt.Errorf("reading response body error: %w", err))
	}
	if response.StatusCode >= http.StatusMultipleChoices {
		return newError(response, data, nil)
	}
	if len(data) > 0 && result != nil {
		err = json.Unmarshal(data, &result)
		if err != nil {
			return newError(response, data, fmt.Errorf("unmarshaling response error: %w", err))
		}
	}
	return nil
}

// Method returns a function implementation of the HTTP Method from the Client by its name.
// Returns Client.GET as the default method.
func (c *JSONClient) Method(method string) func(ctx context.Context, url string, body interface{}, result interface{}) error {
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
func (c *JSONClient) GET(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodGet, url, body, &result)
}

// HEAD is an alias to do the Request with the http.MethodHead method.
func (c *JSONClient) HEAD(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodHead, url, body, &result)
}

// POST is an alias to do the Request with the http.MethodPost method.
func (c *JSONClient) POST(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, url, body, &result)
}

// PUT is an alias to do the Request with the http.MethodPut method.
func (c *JSONClient) PUT(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, url, body, &result)
}

// PATCH is an alias to do the Request with the http.MethodPatch method.
func (c *JSONClient) PATCH(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, url, body, &result)
}

// DELETE is an alias to do the Request with the http.MethodDelete method.
func (c *JSONClient) DELETE(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, url, body, &result)
}

// CONNECT is an alias to do the Request with the http.MethodConnect method.
func (c *JSONClient) CONNECT(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodConnect, url, body, &result)
}

// OPTIONS is an alias to do the Request with the http.MethodOptions method.
func (c *JSONClient) OPTIONS(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodOptions, url, body, &result)
}

// TRACE is an alias to do the Request with the http.MethodTrace method.
func (c *JSONClient) TRACE(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodTrace, url, body, &result)
}

// Clone will clone an instance of the JSONClient without references to the old one.
func (c *JSONClient) Clone() *JSONClient {
	return JSON(c.Client.Clone())
}

func marshal(body interface{}) (data []byte, err error) {
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, newError(nil, nil, fmt.Errorf("marshaling request error: %w", err))
		}
	}
	return data, nil
}
