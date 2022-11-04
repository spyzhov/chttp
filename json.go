package chttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JSONClient struct {
	*Client
}

func JSON(client *Client) *JSONClient {
	return &JSONClient{Client: client}
}

func NewJSON(client *http.Client) *JSONClient {
	return &JSONClient{Client: NewClient(client)}
}

func (c *JSONClient) Request(ctx context.Context, method string, url string, body interface{}, result interface{}) (err error) {
	var data []byte
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling body error: %w", err)
		}
	}
	res, err := c.Client.Request(ctx, method, url, data)
	if err != nil {
		return err
	}
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return NewError(res.StatusCode, data, fmt.Errorf("reading response body error: %w", err))
	}
	if res.StatusCode >= http.StatusMultipleChoices {
		return NewError(res.StatusCode, data, nil)
	}
	if len(data) > 0 {
		err = json.Unmarshal(data, &result)
		if err != nil {
			return NewError(res.StatusCode, data, fmt.Errorf("unmarshaling response error: %w", err))
		}
	}
	return nil
}

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

func (c *JSONClient) GET(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodGet, url, body, &result)
}

func (c *JSONClient) HEAD(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodHead, url, body, &result)
}

func (c *JSONClient) POST(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, url, body, &result)
}

func (c *JSONClient) PUT(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, url, body, &result)
}

func (c *JSONClient) PATCH(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, url, body, &result)
}

func (c *JSONClient) DELETE(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, url, body, &result)
}

func (c *JSONClient) CONNECT(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodConnect, url, body, &result)
}

func (c *JSONClient) OPTIONS(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodOptions, url, body, &result)
}

func (c *JSONClient) TRACE(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodTrace, url, body, &result)
}
