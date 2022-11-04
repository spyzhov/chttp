package chttp

import (
	"bytes"
	"context"
	"net/http"
	"sync"
)

type Client struct {
	http        *http.Client
	middlewares []Middleware
	mu          sync.Mutex
}

func NewClient(client *http.Client) *Client {
	if client == nil {
		client = new(http.Client)
	}
	clone := *client

	result := &Client{
		http:        &clone,
		middlewares: []Middleware{},
	}
	clone.Transport = result.transport(clone.Transport)

	return result
}

func (c *Client) With(middleware ...Middleware) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.middlewares = append(c.middlewares, middleware...)
}

func (c *Client) JSON() *JSONClient {
	return JSON(c)
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	return c.http.Do(request)
}

func (c *Client) Request(ctx context.Context, method string, url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

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

func (c *Client) GET(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, url, getBody(body))
}

func (c *Client) HEAD(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodHead, url, getBody(body))
}

func (c *Client) POST(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, url, getBody(body))
}

func (c *Client) PUT(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, url, getBody(body))
}

func (c *Client) PATCH(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, url, getBody(body))
}

func (c *Client) DELETE(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, url, getBody(body))
}

func (c *Client) CONNECT(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodConnect, url, getBody(body))
}

func (c *Client) OPTIONS(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodOptions, url, getBody(body))
}

func (c *Client) TRACE(ctx context.Context, url string, body ...[]byte) (*http.Response, error) {
	return c.Request(ctx, http.MethodTrace, url, getBody(body))
}

func getBody(body [][]byte) []byte {
	if len(body) > 0 {
		return body[0]
	}
	return nil
}
