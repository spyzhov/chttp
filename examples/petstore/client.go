package petstore

import (
	ulrp "net/url"
)

type Client struct {
	Pets  *PetClient
	Store *StoreClient
	User  *UserClient
}

func NewClient(host string, apiKey string) *Client {
	return &Client{
		Pets:  NewPetClient(host, apiKey),
		Store: NewStoreClient(host, apiKey),
		User:  NewUserClient(host, apiKey),
	}
}

func url(host, path string, query map[string][]string) (string, error) {
	result, err := ulrp.JoinPath(host, path)
	if err != nil {
		return "", err
	}
	uri, err := ulrp.Parse(result)
	if err != nil {
		return "", err
	}
	uri.RawQuery = ulrp.Values(query).Encode()
	return uri.String(), nil
}
