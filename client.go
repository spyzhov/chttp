package chttp

import (
	"bytes"
	"context"
	"net/http"
	"sync"
)

// Client is an HTTP client wrapper that injects middleware transport in the http.Client and provides a list of useful methods.
type Client struct {
	HTTP        *http.Client
	middlewares []Middleware
	base        http.RoundTripper
	mu          sync.RWMutex
}

// NewClient is a constructor for the Client. If no http.Client provided as an argument, a new client will be created.
func NewClient(client *http.Client, options ...Option) *Client {
	if client == nil {
		client = new(http.Client)
	}
	clone := *client

	result := &Client{
		HTTP:        &clone,
		middlewares: []Middleware{},
		base:        clone.Transport,
	}
	clone.Transport = result.transport()

	for _, opt := range options {
		opt(result)
	}

	return result
}

// With appends Middleware to the end of the list. The first added will be called first (FIFO).
func (c *Client) With(middleware ...Middleware) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.middlewares = append(c.middlewares, middleware...)
}

func (c *Client) getMiddlewares() []Middleware {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := make([]Middleware, len(c.middlewares))
	copy(clone, c.middlewares)
	return clone
}

// JSON creates a JSONClient wrapper with the given Client as a basic one.
func (c *Client) JSON() *JSONClient {
	return JSON(c)
}

// Do is the proxy to the same method in the http.Client.
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	return c.HTTP.Do(request)
}

// Request creates a http.Request with the given arguments and calls the Client.Do method.
func (c *Client) Request(ctx context.Context, method string, url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

// Method returns a function implementation of the HTTP Method from the Client by its name.
// Returns Client.GET as the default method.
func (c *Client) Method(method string) func(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
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
func (c *Client) GET(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, url, getBody(body))
}

// HEAD is an alias to do the Request with the http.MethodHead method.
func (c *Client) HEAD(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodHead, url, getBody(body))
}

// POST is an alias to do the Request with the http.MethodPost method.
func (c *Client) POST(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, url, getBody(body))
}

// PUT is an alias to do the Request with the http.MethodPut method.
func (c *Client) PUT(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, url, getBody(body))
}

// PATCH is an alias to do the Request with the http.MethodPatch method.
func (c *Client) PATCH(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, url, getBody(body))
}

// DELETE is an alias to do the Request with the http.MethodDelete method.
func (c *Client) DELETE(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, url, getBody(body))
}

// CONNECT is an alias to do the Request with the http.MethodConnect method.
func (c *Client) CONNECT(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodConnect, url, getBody(body))
}

// OPTIONS is an alias to do the Request with the http.MethodOptions method.
func (c *Client) OPTIONS(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodOptions, url, getBody(body))
}

// TRACE is an alias to do the Request with the http.MethodTrace method.
func (c *Client) TRACE(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodTrace, url, getBody(body))
}

// Clone will clone an instance of the Client without references to the old one.
func (c *Client) Clone() *Client {
	if c == nil {
		return nil
	}
	httpClient := *c.HTTP

	clone := &Client{
		HTTP:        &httpClient,
		middlewares: make([]Middleware, len(c.middlewares)),
		base:        c.base,
	}
	copy(clone.middlewares, c.middlewares)
	httpClient.Transport = clone.transport()

	return clone
}

func getBody(body [][]byte) []byte {
	if len(body) > 0 {
		return body[0]
	}
	return nil
}
