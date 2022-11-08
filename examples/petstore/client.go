package petstore

type Client struct {
	Pets *PetClient
}

func NewClient(host string, apiKey string) *Client {
	return &Client{
		Pets: NewPetClient(host, apiKey),
	}
}
