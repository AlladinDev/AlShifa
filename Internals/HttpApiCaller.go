package internals

import "net/http"

type APIClient struct {
	api *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		api: &http.Client{},
	}
}
