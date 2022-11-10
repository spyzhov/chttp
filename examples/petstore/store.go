package petstore

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spyzhov/chttp"
	"github.com/spyzhov/chttp/middleware"
)

// StoreClient access to Petstore orders
type StoreClient struct {
	path   string
	client *chttp.JSONClient
}

func NewStoreClient(host string, apiKey string) *StoreClient {
	client := chttp.NewJSON(&http.Client{Timeout: 10 * time.Second})
	client.With(middleware.JSON())
	client.With(middleware.Headers(map[string]string{
		"api_key": apiKey,
	}, true))

	return &StoreClient{
		path:   host + "/store",
		client: client,
	}
}

// Inventory calls HTTP GET /store/inventory Returns pet inventories by status
func (c *StoreClient) Inventory(ctx context.Context) (result map[OrderStatus]int, err error) {
	uri, err := url(c.path, "", nil)
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// PlaceOrder calls HTTP POST /store/order Place an Order for a Pet
func (c *StoreClient) PlaceOrder(ctx context.Context, order Order) (result Order, err error) {
	uri, err := url(c.path, "/order", nil)
	if err != nil {
		return result, err
	}
	err = c.client.POST(ctx, uri, order, &result)
	return result, err
}

// FindOrder calls HTTP GET /store/order/{orderId} Find purchase Order by ID
func (c *StoreClient) FindOrder(ctx context.Context, orderId int) (result Order, err error) {
	uri, err := url(c.path, fmt.Sprintf("/order/%d", orderId), nil)
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// DeleteOrder calls HTTP DELETE /store/order/{orderId} Delete purchase Order by ID
func (c *StoreClient) DeleteOrder(ctx context.Context, orderId int) (err error) {
	uri, err := url(c.path, fmt.Sprintf("/order/%d", orderId), nil)
	if err != nil {
		return err
	}
	err = c.client.DELETE(ctx, uri, nil, nil)
	return err
}
