package petstore

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/spyzhov/chttp"
	"github.com/spyzhov/chttp/middleware"
)

// PetClient provides everything about your Pets
type PetClient struct {
	path   string
	client *chttp.JSONClient
}

func NewPetClient(host string, apiKey string) *PetClient {
	client := chttp.NewJSON(&http.Client{Timeout: 10 * time.Second})
	client.With(middleware.JSON())
	client.With(middleware.Headers(map[string]string{
		"api_key": apiKey,
	}, true))

	return &PetClient{
		path:   host + "/pet",
		client: client,
	}
}

// Update calls HTTP PUT /pet Update an existing pet
func (c *PetClient) Update(ctx context.Context, pet Pet) (result Pet, err error) {
	uri, err := url(c.path, "", nil)
	if err != nil {
		return result, err
	}
	err = c.client.PUT(ctx, uri, pet, &result)
	return result, err
}

// Add calls HTTP POST /pet Update an existing pet
func (c *PetClient) Add(ctx context.Context, pet Pet) (result Pet, err error) {
	uri, err := url(c.path, "", nil)
	if err != nil {
		return result, err
	}
	err = c.client.POST(ctx, uri, pet, &result)
	return result, err
}

// FindByStatus calls HTTP POST /pet/findByStatus Finds Pets by status
func (c *PetClient) FindByStatus(ctx context.Context, status string) (result Pets, err error) {
	uri, err := url(c.path, "/findByStatus", map[string][]string{"status": {status}})
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// FindByTags calls HTTP POST /pet/findByTags Finds Pets by tags
func (c *PetClient) FindByTags(ctx context.Context, tags []string) (result Pets, err error) {
	uri, err := url(c.path, "/findByTags", map[string][]string{"tags[]": tags})
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// Find calls HTTP GET /pet/{petId} Find pet by ID
func (c *PetClient) Find(ctx context.Context, petId string) (result Pet, err error) {
	uri, err := url(c.path, "/"+petId, nil)
	if err != nil {
		return result, err
	}
	err = c.client.GET(ctx, uri, nil, &result)
	return result, err
}

// UpdateByID calls HTTP /pet/{petId} Updates a pet in the store with form data
func (c *PetClient) UpdateByID(ctx context.Context, petId string, name *string, status *string) (result Pet, err error) {
	params := make(map[string][]string)
	if name != nil {
		params["name"] = []string{*name}
	}
	if status != nil {
		params["status"] = []string{*status}
	}
	uri, err := url(c.path, "/"+petId, params)
	if err != nil {
		return result, err
	}
	err = c.client.POST(ctx, uri, nil, &result)
	return result, err
}

// Delete calls HTTP DELETE /pet/{petId} Deletes a pet
func (c *PetClient) Delete(ctx context.Context, petId string) (err error) {
	uri, err := url(c.path, "/"+petId, nil)
	if err != nil {
		return err
	}
	err = c.client.DELETE(ctx, uri, nil, nil)
	return err
}

// UploadImage calls HTTP POST /pet/{petId}/uploadImage uploads an image
func (c *PetClient) UploadImage(ctx context.Context, petId string, image io.Reader) (result ApiResponse, err error) {
	uri, err := url(c.path, "/"+petId+"/uploadImage", nil)
	if err != nil {
		return result, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, image)
	request.Header.Set("Content-Type", "application/octet-stream")
	response, err := c.client.Do(request)
	if err != nil {
		return result, err
	}
	err = c.client.UnmarshalHTTPResponse(response, err, &result)
	return result, err
}
