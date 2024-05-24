package chttp

import (
	"net/http"
	"time"
)

type Option func(c *Client)

// WithMiddleware set middlewares for the client
func WithMiddleware(middleware ...Middleware) Option {
	return func(c *Client) {
		c.With(middleware...)
	}
}

// WithTimeout set timeout for the http.Client
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.HTTP.Timeout = timeout
	}
}

// WithCookieJar set the cookie jar for the http.Client
func WithCookieJar(jar http.CookieJar) Option {
	return func(c *Client) {
		c.HTTP.Jar = jar
	}
}

// WithCheckRedirect set the check redirect function for the http.Client
func WithCheckRedirect(check func(req *http.Request, via []*http.Request) error) Option {
	return func(c *Client) {
		c.HTTP.CheckRedirect = check
	}
}
