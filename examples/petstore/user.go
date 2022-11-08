package petstore

import (
	"context"
	"net/http"
	"time"

	"github.com/spyzhov/chttp"
	"github.com/spyzhov/chttp/middleware"
)

// UserClient operations about User
type UserClient struct {
	path   string
	client *chttp.JSONClient
}

func NewUserClient(host string, apiKey string) *UserClient {
	client := chttp.NewJSON(&http.Client{Timeout: 10 * time.Second})
	client.With(middleware.JSON())
	client.With(middleware.Headers(map[string]string{
		"api_key": apiKey,
	}, true))

	return &UserClient{
		path:   host + "/user",
		client: client,
	}
}

// Create calls HTTP POST /user Create user
func (c *UserClient) Create(ctx context.Context, user User) (result User, err error) {
	uri, err := url(c.path, "", nil)
	if err != nil {
		return result, err
	}
	err = c.client.POST(ctx, uri, user, &result)
	return result, err
}

// CreateWithList calls HTTP POST /user/createWithList Creates list of users with given input array
func (c *UserClient) CreateWithList(ctx context.Context, users []User) (result []User, err error) {
	uri, err := url(c.path, "/createWithList", nil)
	if err != nil {
		return result, err
	}
	err = c.client.POST(ctx, uri, users, &result)
	return result, err
}

// Login calls HTTP GET /user/login Logs User into the system
func (c *UserClient) Login(ctx context.Context, username, password string) (result string, err error) {
	uri, err := url(c.path, "/login", map[string][]string{"username": {username}, "password": {password}})
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// Logout calls HTTP GET /user/logout Logs out current logged-in user session
func (c *UserClient) Logout(ctx context.Context) (err error) {
	uri, err := url(c.path, "/logout", nil)
	if err != nil {
		return err
	}
	err = c.client.GET(ctx, uri, nil, nil)
	return err
}

// GetByName calls HTTP GET /user/{username} Get user by username
func (c *UserClient) GetByName(ctx context.Context, username string) (result User, err error) {
	uri, err := url(c.path, "/"+username, nil)
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// Update calls HTTP PUT /user/{username} Update user
func (c *UserClient) Update(ctx context.Context, username string, user User) (result User, err error) {
	uri, err := url(c.path, "/"+username, nil)
	if err != nil {
		return result, err
	}
	err = c.client.PUT(ctx, uri, user, &result)
	return result, err
}

// Delete calls HTTP DELETE /user/{username} Delete user
func (c *UserClient) Delete(ctx context.Context, username string) (err error) {
	uri, err := url(c.path, "/"+username, nil)
	if err != nil {
		return err
	}
	err = c.client.DELETE(ctx, uri, nil, nil)
	return err
}
