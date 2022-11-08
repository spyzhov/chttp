package petstore

type Pet struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
	Status string `json:"status"`
}

type Pets []Pet

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
